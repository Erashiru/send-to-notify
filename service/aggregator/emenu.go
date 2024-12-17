package aggregator

import (
	"context"
	"fmt"
	externalApiModels "github.com/kwaaka-team/orders-core/core/externalapi/models"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models3 "github.com/kwaaka-team/orders-core/core/wolt/models"
	"github.com/kwaaka-team/orders-core/pkg/externalapi/clients"
	externalModels "github.com/kwaaka-team/orders-core/pkg/externalapi/clients/dto"
	externalHttp "github.com/kwaaka-team/orders-core/pkg/externalapi/clients/http"
	"github.com/pkg/errors"
)

type emenuService struct {
	*externalService
	deliveryServiceName      models.Aggregator
	client                   clients.Client
	webhookURL               string
	webhookProductStopList   string
	webhookAttributeStopList string
	storeID                  string
}

func newEmenuService(store storeModels.Store) (*emenuService, error) {
	externalService, err := newExternalService()
	if err != nil {
		return nil, errors.Wrap(constructorError, "emenuService constructor error")
	}

	aggregatorName := models.EMENU

	var extCfg *storeModels.StoreExternalConfig
	for _, config := range store.ExternalConfig {
		if config.Type != aggregatorName.String() {
			continue
		}
		extCfg = &config
	}
	if extCfg == nil {
		return nil, errors.Wrap(constructorError, fmt.Sprintf("aggregator %s is not found for store %s", aggregatorName, store.Name))
	}

	authToken := extCfg.AuthToken
	webhookURL := extCfg.WebhookURL

	webhookProductStopList := extCfg.WebhookProductStoplist
	if webhookProductStopList == "" {
		webhookProductStopList = webhookURL
	}

	webhookAttributeStopList := extCfg.WebhookAttributeStoplist
	if webhookAttributeStopList == "" {
		webhookAttributeStopList = webhookURL
	}

	storeID := extCfg.StoreID[0]

	externalClient, err := externalHttp.NewClient(&clients.Config{
		AuthToken: authToken,
	})
	if err != nil {
		return nil, errors.Wrap(err, constructorError.Error())
	}

	return &emenuService{
		externalService, aggregatorName, externalClient, webhookURL,
		webhookProductStopList, webhookAttributeStopList, storeID,
	}, nil
}

func (s *emenuService) OpenStore(ctx context.Context, aggregatorStoreId string) error {
	return errors.New("method not implemented")
}

func (s *emenuService) GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error) {
	return false, errors.New("method not implemented")
}

func (s *emenuService) GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error) {
	return storeModels.AggregatorSchedule{}, errors.New("method not implemented")
}

func (s *emenuService) MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	return s.mapSystemStatusToAggregatorStatus(s.deliveryServiceName, order, posStatus, store)
}

func (s *emenuService) UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error {
	return s.client.UpdateOrderWebhook(ctx, externalModels.Order{
		OrderID: order.OrderID,
		Status:  aggregatorStatus,
	}, s.webhookURL)
}

func (s *emenuService) IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error) {
	for _, item := range store.ExternalConfig {
		if item.Type != models.EMENU.String() {
			continue
		}
		return item.IsMarketplace, nil
	}

	return false, errors.New("is marketplace: emenu object not found in store")
}

func (s *emenuService) SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	return nil, nil
}

func (s *emenuService) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	order, ok := req.(externalApiModels.Order)
	if !ok {
		return "", errors.New("casting error")
	}

	return order.RestaurantId, nil
}

func (s *emenuService) GetSystemCreateOrderRequestByAggregatorRequest(req interface{}, store storeModels.Store) (models.Order, error) {
	return s.getSystemCreateOrderRequestByAggregatorRequest(req, store, s.deliveryServiceName.String())
}

func (s *emenuService) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error) {
	modifiersRequest := make([]externalModels.Modifier, 0, len(attributes))

	for _, attribute := range attributes {
		modifiersRequest = append(modifiersRequest, s.toModifier(attribute, attribute.IsAvailable))
	}

	utils.Beautify("modifiers request to update", modifiersRequest)

	for _, modifier := range modifiersRequest {
		if err := s.client.UpdateModifierStopList(ctx, modifier, s.webhookAttributeStopList); err != nil {
			return "", err
		}
	}

	return "", nil
}

func (s *emenuService) toModifier(attribute menuModels.Attribute, isAvailable bool) externalModels.Modifier {
	return externalModels.Modifier{
		StoreID:    s.storeID,
		ModifierID: attribute.ExtID,
		Price:      attribute.Price,
		Available:  isAvailable,
	}
}

func (s *emenuService) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error) {

	var items = make([]externalModels.Product, 0, len(products))

	for _, product := range products {
		items = append(items, s.toItem(product, isAvailable))
	}

	for _, item := range items {
		if err := s.client.UpdateProductStopList(ctx, item, s.webhookProductStopList); err != nil {
			return "", err
		}
	}

	return "", nil

}

func (s *emenuService) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error) {

	var items = make([]externalModels.Product, 0, len(products))

	for _, product := range products {
		items = append(items, s.toItem(product, product.IsAvailable))
	}

	for _, item := range items {
		if err := s.client.UpdateProductStopList(ctx, item, s.webhookProductStopList); err != nil {
			return "", err
		}
	}

	return "", nil

}

func (s *emenuService) toItem(product menuModels.Product, isAvailable bool) externalModels.Product {
	if len(product.Price) == 0 {
		return externalModels.Product{}
	}

	return externalModels.Product{
		StoreID:   s.storeID,
		ProductID: product.ExtID,
		Price:     product.Price[0].Value,
		Available: isAvailable,
	}
}

func (s *emenuService) GetAggregatorOrder(ctx context.Context, orderID string) (models3.Order, error) {
	return models3.Order{}, nil
}

func (s *emenuService) SendOrderErrorNotification(ctx context.Context, req interface{}) error {
	return nil
}

func (s *emenuService) SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error {
	return nil
}
