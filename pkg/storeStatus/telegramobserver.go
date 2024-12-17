package storeStatus

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	orderModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	orderService "github.com/kwaaka-team/orders-core/service/order"
)

type TelegramObserver struct {
	TelegramClient orderService.TelegramService
}

func (s TelegramObserver) Notify(ctx context.Context, restaurant models.Store, externalStoreID, deliveryService string, storeIsOpened bool) error {
	switch {
	case storeIsOpened:
		return s.TelegramClient.SendMessageToQueue(telegram.StoreClosed, orderModels.Order{}, restaurant, "", telegram.StoreIsOpened, deliveryService, models2.Product{})
	default:
		return s.TelegramClient.SendMessageToQueue(telegram.StoreClosed, orderModels.Order{}, restaurant, "", "", deliveryService, models2.Product{})
	}
}

func (s TelegramObserver) NotifyStatusReport(ctx context.Context, restaurant models.Store, durations []models.OpenTimeDuration) error {
	return s.TelegramClient.SendMessageToQueue(telegram.StoreStatusReport, orderModels.Order{}, restaurant, "", telegram.ConstructStoreStatusReportToNotify(restaurant, durations), "", models2.Product{})
}

func (s TelegramObserver) NotifyStatusChange(ctx context.Context, status string, phone string) error {
	// TODO: implement telegram notifier
	return nil
}
