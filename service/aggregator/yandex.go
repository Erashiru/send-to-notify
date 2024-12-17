package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/config"
	externalApiModels "github.com/kwaaka-team/orders-core/core/externalapi/models"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models3 "github.com/kwaaka-team/orders-core/core/wolt/models"
	yandexClient "github.com/kwaaka-team/orders-core/pkg/yandex"
	yandexConfig "github.com/kwaaka-team/orders-core/pkg/yandex/clients"
	yandexModels "github.com/kwaaka-team/orders-core/pkg/yandex/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type yandexService struct {
	*externalService
	deliveryServiceName models.Aggregator
	cli                 yandexConfig.Yandex
}

func newYandexService(cfg config.YandexConfiguration) (*yandexService, error) {
	externalService, err := newExternalService()
	if err != nil {
		return nil, errors.Wrap(constructorError, "yandexService constructor error")
	}
	cli, err := yandexClient.NewClient(&yandexConfig.Config{
		Protocol:     "http",
		BaseURL:      cfg.BaseURL,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	})

	return &yandexService{
		externalService, models.YANDEX, cli,
	}, nil
}

func (s *yandexService) GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error) {
	return false, errors.New("method not implemented")
}

func (s *yandexService) OpenStore(ctx context.Context, aggregatorStoreId string) error {
	return errors.New("method not implemented")
}

func (s *yandexService) GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error) {
	return storeModels.AggregatorSchedule{}, errors.New("method not implemented")
}

func (s *yandexService) MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	return s.mapSystemStatusToAggregatorStatus(s.deliveryServiceName, order, posStatus, store)
}

func (s *yandexService) UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error {
	return nil
}

func (s *yandexService) SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	return s.splitVirtualStoreOrder(req, store)
}

func (s *yandexService) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	order, ok := req.(externalApiModels.Order)
	if !ok {
		return "", errors.New("casting error")
	}

	return order.RestaurantId, nil
}

func (s *yandexService) GetSystemCreateOrderRequestByAggregatorRequest(req interface{}, store storeModels.Store) (models.Order, error) {
	return s.getSystemCreateOrderRequestByAggregatorRequest(req, store, s.deliveryServiceName.String())
}

func (s *yandexService) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error) {
	if err := s.cli.MenuImportInitiation(ctx, yandexModels.MenuInitiationRequest{
		RestaurantID:  aggregatorStoreID,
		OperationType: yandexModels.StoplistOperationType,
	}); err != nil {
		log.Info().Msgf("error updating stop list by products: %v", err)
		return "", nil
	}

	return "", nil
}

func (s *yandexService) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error) {
	if err := s.cli.MenuImportInitiation(ctx, yandexModels.MenuInitiationRequest{
		RestaurantID:  aggregatorStoreID,
		OperationType: yandexModels.StoplistOperationType,
	}); err != nil {
		log.Info().Msgf("error updating stop list by products bulk: %v", err)
		return "", nil
	}

	return "", nil
}

func (s *yandexService) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error) {
	if err := s.cli.MenuImportInitiation(ctx, yandexModels.MenuInitiationRequest{
		RestaurantID:  aggregatorStoreID,
		OperationType: yandexModels.StoplistOperationType,
	}); err != nil {
		log.Info().Msgf("error updating stop list by attributes bulk: %v", err)
		return "", nil
	}

	return "", nil
}

func (s *yandexService) IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error) {
	return restaurantSelfDelivery, nil
}

func (s *yandexService) GetAggregatorOrder(ctx context.Context, orderID string) (models3.Order, error) {
	return models3.Order{}, nil
}

func (s *yandexService) SendOrderErrorNotification(ctx context.Context, req interface{}) error {
	return nil
}

func (s *yandexService) SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error {
	if err := s.cli.MenuImportInitiation(ctx, yandexModels.MenuInitiationRequest{
		RestaurantID:  aggregatorStoreID,
		OperationType: yandexModels.StoplistOperationType,
	}); err != nil {
		log.Info().Msgf("error sending stop list update notification: %v", err)
		return err
	}
	return nil
}
