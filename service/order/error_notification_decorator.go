package order

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp/clients"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ErrorNotificationDecorator struct {
	service         CreationService
	storeService    store.Service
	telegramService TelegramService
	whatsappService clients.Whatsapp
}

func NewErrorNotificationDecorator(service CreationService, storeService store.Service, telegramService TelegramService, whatsappService clients.Whatsapp) (*ErrorNotificationDecorator, error) {
	if service == nil {
		return nil, errors.New("service is nil")
	}
	if storeService == nil {
		return nil, errors.New("storeService is nil")
	}
	if telegramService == nil {
		return nil, errors.New("telegramService is nil")
	}
	if whatsappService == nil {
		return nil, errors.New("whatsappService is nil")
	}
	return &ErrorNotificationDecorator{
		service:         service,
		storeService:    storeService,
		telegramService: telegramService,
		whatsappService: whatsappService,
	}, nil
}

func (s *ErrorNotificationDecorator) CreateOrder(ctx context.Context, externalStoreID, deliveryService string, aggReq interface{}, storeSecret string) (models.Order, error) {
	order, errOrder := s.service.CreateOrder(ctx, externalStoreID, deliveryService, aggReq, storeSecret)
	if errOrder == nil {
		return order, nil
	}

	if !errors.Is(errOrder, errWithNotification) {
		return order, errOrder
	}

	st, err := s.storeService.GetByExternalIdAndDeliveryService(ctx, externalStoreID, deliveryService)
	if err != nil {
		log.Err(err).Msgf("get store by external id and delivery service error: %s", externalStoreID)
		return order, errOrder
	}

	if tgErr := s.telegramService.SendMessageToQueue(telegram.CreateOrder, order, st, errOrder.Error(), "", "", models2.Product{}); tgErr != nil {
		log.Err(tgErr).Msg("send telegram message error")
	}

	return order, errOrder
}
