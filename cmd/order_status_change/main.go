package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/core/database"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	notifyClient "github.com/kwaaka-team/orders-core/pkg/que"
	orderServicePkg "github.com/kwaaka-team/orders-core/service/order"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"os"
	"strings"
)

const (
	baseUrlEnv = "BASE_URL"
)

type app struct {
	bot             *tgbotapi.BotAPI
	toSend          tgbotapi.MessageConfig
	telegramService orderServicePkg.TelegramServiceImpl
	store           *store.ServiceImpl
	restyClient     *resty.Client
}

type reviewInfo struct {
	OrderID      string  `json:"order_id"`
	RestaurantID string  `json:"restaurant_id"`
	Review       string  `json:"review"`
	Rating       float32 `json:"rating,omitempty"`
}

func main() {
	lambda.Start(Run)
}

func Run(ctx context.Context, request json.RawMessage) error {
	opts, err := general.LoadConfig(ctx)
	if err != nil {
		return err
	}
	db, err := initDB(opts)
	if err != nil {
		return err
	}

	baseURL := os.Getenv(baseUrlEnv)
	client := resty.New().
		SetBaseURL(baseURL)
	storeRepository, err := store.NewStoreMongoRepository(db)
	if err != nil {
		return err
	}
	storeFactory, err := store.NewService(storeRepository)
	if err != nil {
		return err
	}

	sqsCli := notifyClient.NewSQS(sqs.NewFromConfig(opts.AwsConfig))
	telegramRepo := telegram.NewTelegramRepo(db.Client().Database(opts.DSDB))
	telegramService, err := orderServicePkg.NewTelegramService(sqsCli, opts.QueConfiguration.Telegram, opts.NotificationConfiguration, telegramRepo)
	if err != nil {
		return err
	}

	var update tgbotapi.Update
	if err := json.Unmarshal(request, &update); err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	bot, err := tgbotapi.NewBotAPI(opts.OrderBotToken)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	var chatID int64
	if update.Message != nil {
		chatID = update.Message.Chat.ID
		log.Info().Msgf("Received message from chat ID %d: %s\n", update.Message.Chat.ID, update.Message.Text)
	} else if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
		log.Info().Msgf("Received message from chat ID %d: %s\n", update.CallbackQuery.From.ID, update.CallbackQuery.Message.Text)
	} else {
		log.Error().Msgf("invalid update type")
		return errors.New("invalid update type")
	}
	msg := tgbotapi.NewMessage(chatID, "")

	app := &app{
		bot:             bot,
		toSend:          msg,
		telegramService: telegramService,
		store:           storeFactory,
		restyClient:     client,
	}

	if update.Message != nil {
		app.handleMessage(ctx, update.Message)
	} else if update.CallbackQuery != nil {
		app.handleCallback(ctx, update.CallbackQuery)
	}

	return nil
}

func initDB(opts general.Configuration) (*mongo.Database, error) {
	ds, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSDB,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create datastore %s: %v", opts.DSName, err)
	}

	if err = ds.Connect(ds.Client()); err != nil {
		return nil, fmt.Errorf("cannot connect to datastore: %s", err)
	}
	mongoClient := ds.Client()

	return mongoClient.Database(opts.DSDB), nil
}

func (a *app) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if msg.IsCommand() {
		a.processCommand(ctx, msg)
		return
	}

	a.processMessage(ctx, msg)
}

func (a *app) processMessage(ctx context.Context, msg *tgbotapi.Message) {
	status, err := a.telegramService.GetUserStatus(ctx, msg.Chat.ID)
	if err != nil {
		log.Error().Err(err).Msg("error getting user status")
	}

	switch status {
	case telegram.WAIT_REVIEW:
		a.toSend.Text = "Спасибо большое за отзыв, мы это ценим!"
		if err := a.telegramService.UpdateUserStatus(ctx, msg.From.ID, telegram.REVIEWED); err != nil {
			log.Err(err).Msg("error updating user status")
			return
		}
		a.sendMessage()
		a.SetTelegramReview(ctx, msg.From.ID, msg.Text)
	}
}

func (a *app) processCommand(ctx context.Context, msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		a.toSend.Text = "Добро пожаловать в кваака бот!\n\nЗаказывайте блюда легко и доступно, а также отслеживайте их приготовение прямо здесь!"
		err := a.telegramService.SaveTelegramUser(ctx, msg.From.FirstName, msg.Chat.ID)
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}

	a.sendMessage()
}

func (a *app) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	edit := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "Сделано!")
	if _, err := a.bot.Send(edit); err != nil {
		log.Err(err).Msg("error creating new edit message text")
		return
	}

	data := callback.Data
	parts := strings.SplitN(data, ":", 2)
	if len(parts) != 2 {
		log.Error().Msg("invalid callback (len != 2)")
	}
	callbackType := parts[0]
	switch callbackType {
	case "review":
		a.processReview(ctx, parts[1], callback.From.ID)
	}
}

func (a *app) processReview(ctx context.Context, rating string, chatID int64) {
	switch rating {
	case "1", "2", "3":
		a.processBadReview(ctx, rating, chatID)
	case "4", "5":
		if a.processGoodReview(ctx, rating, chatID) {
			return
		}
	default:
		log.Error().Msg("invalid review rating")
	}

	if err := a.telegramService.UpdateUserStatus(ctx, chatID, telegram.WAIT_REVIEW); err != nil {
		log.Err(err).Msg("error updating user status")
		return
	}
	a.sendMessage()
}

func (a *app) processBadReview(ctx context.Context, rating string, chatID int64) {
	a.toSend.Text = "Пожалуйста, напишите что именно вам не понравилось, это позволит нам улучшить сервис"
	a.toSend.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "Мне не понравилось, что ",
		Selective:             true,
	}
	a.SaveTelegramReviewRating(ctx, chatID, rating)
}

func (a *app) processGoodReview(ctx context.Context, rating string, chatID int64) bool {
	restID, err := a.telegramService.GetReviewingRestaurantID(ctx, chatID)
	if err != nil {
		log.Err(err).Msg("error getting reviewing restaurant id")
		return true
	}
	link, err := a.store.GetTwoGisReviewLink(ctx, restID)
	if err != nil {
		log.Err(err).Msg("error getting 2gis review link")
		return true
	}

	if link != "" {
		a.toSend.Text = fmt.Sprintf("Благодарим за оценку. Поделитесь отзывом о нашем ресторане в 2гис:\n\n%s", link)
		if err := a.telegramService.UpdateUserStatus(ctx, chatID, telegram.REVIEWED); err != nil {
			log.Err(err).Msg("error updating user status")
			return true
		}
		a.sendMessage()
		a.SaveTelegramReviewRating(ctx, chatID, rating)

		return true
	}

	a.toSend.Text = "Благодарим за оценку. Поделитесь что именно вам понравилось"
	a.toSend.ReplyMarkup = tgbotapi.ForceReply{
		ForceReply:            true,
		InputFieldPlaceholder: "Мне понравилось, что ",
		Selective:             true,
	}

	return false
}

func (a *app) SaveTelegramReviewRating(ctx context.Context, chatID int64, ratingStr string) {
	orderID, err := a.telegramService.GetReviewingOrderID(ctx, chatID)
	if err != nil {
		log.Err(err).Msg("error getting reviewing order id")
		return
	}
	if err := a.telegramService.SaveTelegramReviewRating(ctx, orderID, ratingStr); err != nil {
		log.Err(err).Msg("error setting telegram review")
		return
	}
}

func (a *app) SetTelegramReview(ctx context.Context, chatID int64, description string) {
	path := "/qr/3pl/add-review"

	orderID, err := a.telegramService.GetReviewingOrderID(ctx, chatID)
	if err != nil {
		log.Err(err).Msg("error getting reviewing order id")
		return
	}
	restaurantID, err := a.telegramService.GetReviewingRestaurantID(ctx, chatID)
	if err != nil {
		log.Err(err).Msg("error getting reviewing restaurant id")
		return
	}
	rating, err := a.telegramService.GetTelegramReviewRatingFromOrder(ctx, orderID)
	if err != nil {
		log.Err(err).Msg("error getting telegram review rating from order")
		return
	}

	review := reviewInfo{
		OrderID:      orderID,
		RestaurantID: restaurantID,
		Review:       description,
		Rating:       float32(rating),
	}

	resp, err := a.restyClient.R().
		SetContext(ctx).
		SetBody(&review).
		EnableTrace().
		Post(path)
	if err != nil {
		log.Err(err).Msg("error sending review set request to admin client")
		return
	}
	if resp.StatusCode() != http.StatusNoContent {
		log.Error().Msg("invalid response status setting telegram review from admin client")
		return
	}
}

func (a *app) sendMessage() {
	if a.toSend.Text != "" {
		if _, err := a.bot.Send(a.toSend); err != nil {
			log.Err(err).Msgf("error sending telegram message: %s, err: %v", a.toSend.Text, err)
		}
	}
}
