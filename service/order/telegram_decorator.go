package order

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	storeServicePkg "github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type TelegramServiceDecorator struct {
	service           CreationService
	telegramService   TelegramService
	storeService      storeServicePkg.Service
	targetPosTypes    map[string]struct{}
	targetAggregators map[string]struct{}
}

func NewTelegramServiceDecorator(
	service CreationService,
	telegramService TelegramService,
	storeService storeServicePkg.Service,
	targetAggregators []models.Aggregator,
	targetPosTypes ...models.Pos,
) (*TelegramServiceDecorator, error) {
	if service == nil {
		return nil, errors.New("order service is nil")
	}

	if telegramService == nil {
		return nil, errors.New("telegram service is nil")
	}

	if storeService == nil {
		return nil, errors.New("store service is nil")
	}

	targetPosTypesMap := make(map[string]struct{}, 0)
	for i := range targetPosTypes {
		posType := targetPosTypes[i]
		targetPosTypesMap[posType.String()] = struct{}{}
	}

	targetAggregatorsMap := make(map[string]struct{}, 0)

	for i := range targetAggregators {
		targetAggregatorsMap[targetAggregators[i].String()] = struct{}{}
	}

	return &TelegramServiceDecorator{
		service:           service,
		telegramService:   telegramService,
		storeService:      storeService,
		targetPosTypes:    targetPosTypesMap,
		targetAggregators: targetAggregatorsMap,
	}, nil
}

func (s *TelegramServiceDecorator) CreateOrder(ctx context.Context, externalStoreID, deliveryService string, aggReq interface{}, storeSecret string) (models.Order, error) {
	order, err := s.service.CreateOrder(ctx, externalStoreID, deliveryService, aggReq, storeSecret)

	if tgErr := s.sendTelegramMessage(ctx, order, err); tgErr != nil {
		log.Err(tgErr).Msg("send telegram message error")
	}

	return order, err
}

func (s *TelegramServiceDecorator) sendTelegramMessage(ctx context.Context, order models.Order, orderErr error) error {
	store, err := s.storeService.GetByID(ctx, order.RestaurantID)
	if err != nil {
		return err
	}

	if !s.isToSendNotification(store, order) {
		return nil
	}

	if orderErr != nil {
		if telegramErr := s.telegramService.SendMessageToRestaurant(telegram.CreateOrder, order, store, orderErr.Error()); telegramErr != nil {
			return telegramErr
		}
		if telegramErr := s.telegramService.SendMessageToQueue(telegram.CreateOrder, order, store, orderErr.Error(), "", "", models2.Product{}); telegramErr != nil {
			return telegramErr
		}
		return nil
	}

	if telegramErr := s.telegramService.SendMessageToQueue(telegram.SuccessCreateOrder, order, store, "", "", "", models2.Product{}); telegramErr != nil {
		return telegramErr
	}
	return nil
}

func (s *TelegramServiceDecorator) isToSendNotification(store storeModels.Store, order models.Order) bool {
	if _, ok := s.targetPosTypes[store.PosType]; ok {
		return true
	}

	if _, ok := s.targetAggregators[order.DeliveryService]; ok {
		return true
	}

	return false
}
