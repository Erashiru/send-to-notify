package order

import (
	"context"
	"github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

type TelegramService interface {
	SendMessageToQueue(notificationType telegram.NotificationType, order models.Order, store storeModels.Store, err, msg, extraMsg string, product models2.Product) error
	SendMessageToRestaurant(notificationType telegram.NotificationType, order models.Order, store storeModels.Store, err string) error
}

type TelegramServiceImpl struct {
	sqsCli             que.SQSInterface
	telegramQueue      string
	notificationConfig general.NotificationConfiguration
	repo               *telegram.Repository
}

func NewTelegramService(sqsCli que.SQSInterface, queueName string, notificationConfig general.NotificationConfiguration, r *telegram.Repository) (TelegramServiceImpl, error) {
	if sqsCli == nil {
		return TelegramServiceImpl{}, errors.Wrap(errConstructor, "sqsCli is nil")
	}
	if queueName == "" {
		return TelegramServiceImpl{}, errors.Wrap(errConstructor, "queueName is nil")
	}
	if notificationConfig.OrderBotToken == "" {
		return TelegramServiceImpl{}, errors.Wrap(errConstructor, "OrderBotToken is nil")
	}
	return TelegramServiceImpl{
		sqsCli:             sqsCli,
		telegramQueue:      queueName,
		notificationConfig: notificationConfig,
		repo:               r,
	}, nil
}

func (s TelegramServiceImpl) SendMessageToRestaurant(notificationType telegram.NotificationType, order models.Order,
	store storeModels.Store, err string) error {

	var (
		message   string
		chatIDs   []string
		queueName = s.telegramQueue
		token     = store.Telegram.TelegramBotToken
	)

	switch notificationType {
	case telegram.CreateOrder, telegram.UpdateOrder:
		log.Trace().Msgf("ORDER ERROR (RESTAURANT): %s, ORDER ID: %s", err, order.OrderID)
		message = telegram.ConstructRestaurantMessageToNotify(order, store, err)
		chatIDs = append(chatIDs, "-4020753511") // TEST: "КЕХ - ошибки" chat for testing error specific messages
	}

	chatIDs = s.deleteEmptyChatIDs(chatIDs)

	log.Info().Msgf("(SendMessageToRestaurant) telegram creds: queue - %s  chat_ids - %v ", queueName, chatIDs)

	for _, chatID := range chatIDs {
		if err := s.sqsCli.SendMessage(queueName, message, chatID, token); err != nil {
			log.Err(err).Msgf("queue name: %v, order_id %v, notification type: %v, delivery service %v", queueName, order.OrderID, notificationType, order.DeliveryService)
		}
	}

	return nil
}

func (s TelegramServiceImpl) deleteEmptyChatIDs(chatIDs []string) []string {
	var newChatIDs []string
	for _, chatID := range chatIDs {
		if chatID != "" {
			newChatIDs = append(newChatIDs, chatID)
		}
	}
	return newChatIDs
}

func (s TelegramServiceImpl) SendMessageToQueue(notificationType telegram.NotificationType, order models.Order,
	store storeModels.Store, err, msg, extraMsg string, product models2.Product) error {

	var (
		message     string
		chatIDs     []string
		serviceName = telegram.Telegram
		queueName   = s.telegramQueue
		token       = store.Telegram.TelegramBotToken
	)

	switch order.DeliveryService {
	case models.QRMENU.String(), models.KWAAKA_ADMIN.String():
		if notificationType != telegram.Refund && notificationType != telegram.Compensation && notificationType != telegram.ThirdPartyError {
			chatIDs = append(chatIDs, s.notificationConfig.KwaakaDirectTelegramChatId)
		}
	}

	switch notificationType {
	case telegram.SuccessCreateOrder:
		message = telegram.ConstructSuccessMessage(order, store)
		chatIDs = append(chatIDs, store.Telegram.CreateOrderChatID)
	case telegram.CancelOrder:
		message = telegram.ConstructCancelMessageToNotify(order, store.Name, order.CancelReason.Reason)
		chatIDs = append(chatIDs, store.Telegram.CancelChatID)
	case telegram.CreateOrder, telegram.UpdateOrder:
		log.Trace().Msgf("ORDER ERROR: %v, ORDER ID %s", err, order.OrderID)
		message = telegram.ConstructOrderMessageToNotify(serviceName, order, store, err)
		chatIDs = append(chatIDs, s.notificationConfig.TelegramChatID, "-4020753511") // TEST: "КЕХ - ошибки" chat for testing error specific messages
		if store.Telegram.CreateOrderChatID != "" && s.notDirect(store.Telegram.CreateOrderChatID) {
			chatIDs = append(chatIDs, store.Telegram.CreateOrderChatID)
		}
	case telegram.StoreClosed:
		message = telegram.ConstructStoreClosedToNotify(store, msg, extraMsg)
		chatIDs = append(chatIDs, "-1002038506041")
		if store.Telegram.StoreStatusChatId != "" {
			chatIDs = append(chatIDs, store.Telegram.StoreStatusChatId)
		}
	case telegram.StoreStatusReport:
		message = msg
		chatIDs = append(chatIDs, "-1002038506041")
		if store.Telegram.StoreStatusChatId != "" {
			chatIDs = append(chatIDs, store.Telegram.StoreStatusChatId)
		}
	case telegram.OrderStat:
		message = msg
		chatIDs = append(chatIDs, s.notificationConfig.OrderStatChatID)
	case telegram.OrderStatusChange:
		message = msg
		token = s.notificationConfig.OrderBotToken
		userChatID, err := s.repo.GetUserChatID(context.Background())
		if err != nil {
			return err
		}
		chatIDs = append(chatIDs, userChatID)
	case telegram.ThirdPartyError:
		message = msg
		chatIDs = append(chatIDs, s.notificationConfig.KwaakaDirect3plNotificationsChatID)
		chatIDs = append(chatIDs, store.Kwaaka3PL.ChatID)
	case telegram.NoCourier:
		message = telegram.ConstructNoCourierMessage(order, store)
		chatIDs = []string{s.notificationConfig.KwaakaDirectNoCourierTelegramChatID}
	case telegram.Refund:
		message = telegram.ConstructRefundMessage(order, store, msg, extraMsg)
		chatIDs = append(chatIDs, s.notificationConfig.KwaakaDirectRefundChatID)
	case telegram.Compensation:
		message = telegram.ConstructCompensationMessage(order, store, err, msg, extraMsg)
		chatIDs = append(chatIDs, s.notificationConfig.KwaakaDirectKwaakaAdminCompensationChatID)
	case telegram.CancelDeliveryFromDispatcherPage:
		message = telegram.ConstructCancelDeliveryFromDispatcherPage(order, store)
		chatIDs = append(chatIDs, s.notificationConfig.KwaakaDirectTelegramChatId)
	case telegram.NoDeliveryDispatcher:
		message = telegram.ConstructNoDeliveryDispatcherMessage(order, store)
		chatIDs = append(chatIDs, s.notificationConfig.KwaakaDirect3plNotificationsChatID)
	case telegram.AutoUpdatePublicateMenu:
		message = msg
		chatIDs = append(chatIDs, s.notificationConfig.AutoUpdatePublicateNotificationChatID)
	case telegram.PutProductToStopListWithErrSolution:
		message = telegram.ConstructPutProductToStopListWithErrSolutionMessage(store, product, err)
		if len(message) > 0 {
			chatIDs = append(chatIDs, s.notificationConfig.PutProductToStopListWithErrSolutionChatID)
		}

	default:
		message = msg
		chatIDs = append(chatIDs, "-4167242286")
	}

	log.Info().Msgf("telegram creds: queue %v  chat_id %v ", queueName, chatIDs)
	for _, chatID := range s.uniqueChatIDs(chatIDs) {
		log.Info().Msgf("here telegram send with type: %s for store: %s and chat id: %s", notificationType.String(), store.ID, store.Kwaaka3PL.ChatID)
		if chatID == "" {
			log.Info().Msgf("continue because chat id is empty")
			continue
		}
		if err := s.sqsCli.SendMessage(queueName, message, chatID, token); err != nil {
			log.Err(err).Msgf("queue name: %v, order_id %v, notification type: %v, delivery service %v", queueName, order.OrderID, notificationType, order.DeliveryService)
			return err
		}
	}

	return nil
}

func (s TelegramServiceImpl) uniqueChatIDs(chatIDs []string) []string {
	uniqueMap := make(map[string]struct{})
	var result []string

	for _, chatID := range chatIDs {
		if _, exists := uniqueMap[chatID]; !exists {
			uniqueMap[chatID] = struct{}{}
			result = append(result, chatID)
		}
	}

	return result
}

func (s TelegramServiceImpl) notDirect(chatID string) bool {
	// Aula, Kazbek Saraishyk,Kazbek Bokeikhana
	directChatIDs := []string{"-4265935199", "-4270681463", "-4220435967"}
	for _, id := range directChatIDs {
		if chatID == id {
			return false
		}
	}
	return true
}

func (s TelegramServiceImpl) GetTelegramReviewRatingFromOrder(ctx context.Context, orderID string) (float32, error) {
	return s.repo.GetTelegramReviewRatingFromOrder(ctx, orderID)
}

func (s TelegramServiceImpl) SaveTelegramReviewRating(ctx context.Context, orderID, rating string) error {
	ratingFloat, err := strconv.ParseFloat(rating, 32)
	if err != nil {
		var syntaxError *strconv.NumError
		if errors.As(err, &syntaxError) {
			ratingFloat = 0
		} else {
			return err
		}
	}
	return s.repo.SaveTelegramReviewRating(ctx, orderID, float32(ratingFloat))
}

func (s TelegramServiceImpl) GetReviewingOrderID(ctx context.Context, chatID int64) (string, error) {
	chatIDStr := strconv.Itoa(int(chatID))
	return s.repo.GetReviewingOrderID(ctx, chatIDStr)
}

func (s TelegramServiceImpl) GetReviewingRestaurantID(ctx context.Context, chatID int64) (string, error) {
	chatIDStr := strconv.Itoa(int(chatID))
	return s.repo.GetReviewingRestaurantID(ctx, chatIDStr)
}

func (s TelegramServiceImpl) UpdateReviewingOrderInfo(ctx context.Context, chatID int64, orderID, restID string) error {
	chatIDStr := strconv.Itoa(int(chatID))
	return s.repo.UpdateReviewingOrderInfo(ctx, chatIDStr, orderID, restID)
}

func (s TelegramServiceImpl) GetUserStatus(ctx context.Context, chatID int64) (string, error) {
	chatIDStr := strconv.Itoa(int(chatID))
	return s.repo.GetUserStatus(ctx, chatIDStr)
}

func (s TelegramServiceImpl) UpdateUserStatus(ctx context.Context, chatID int64, status string) error {
	chatIDStr := strconv.Itoa(int(chatID))
	return s.repo.UpdateUserStatus(ctx, chatIDStr, status)
}

func (s TelegramServiceImpl) SaveTelegramUser(ctx context.Context, firstName string, chatID int64) error {
	chatIDStr := strconv.Itoa(int(chatID))
	return s.repo.InsertUser(ctx, firstName, chatIDStr)
}

func (s TelegramServiceImpl) GetUserChatID(ctx context.Context) (*int64, error) {
	chatIDStr, err := s.repo.GetUserChatID(ctx)
	if err != nil {
		return nil, err
	} else if chatIDStr == "" {
		return nil, nil
	}
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return &chatID, nil
}
