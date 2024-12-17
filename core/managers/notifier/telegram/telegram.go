package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	"github.com/kwaaka-team/orders-core/core/models"
	storecoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/order"
	"golang.org/x/net/context"
)

type TelegramNotifier struct {
	bot     *tgbotapi.BotAPI
	service order.TelegramServiceImpl
}

func NewTelegramNotifier(token string, svc order.TelegramServiceImpl) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true

	return &TelegramNotifier{
		bot:     bot,
		service: svc,
	}, nil
}

func (tn *TelegramNotifier) Notify(ctx context.Context, status string, order models.Order, storeGroup storecoreModels.StoreGroup, store storecoreModels.Store) error {
	chatID, err := tn.service.GetUserChatID(ctx)
	if err != nil {
		return err
	}

	if err := tn.service.UpdateReviewingOrderInfo(ctx, *chatID, order.ID, order.RestaurantID); err != nil {
		return err
	}

	defaultMessage := "Неизвестный статус заказа, пожалуйста, обратитесь в службу поддержки"
	msg := tgbotapi.NewMessage(*chatID, defaultMessage)
	telegram.ConstructOrderStatusChangeMessage(&msg, status)

	_, err = tn.bot.Send(msg)
	return err
}
