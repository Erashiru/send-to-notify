package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/core/externalapi/utils"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/order/dto"
	"github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
)

type OrderClient interface {
	UpdateOrder(ctx context.Context, order models.Order, orderID, service, clientSecret string) error
	GetOrder(ctx context.Context, orderID, service string) (models.Order, error)
	CancelOrder(ctx context.Context, req models.CancelOrderRequest, id, service, clientSecret string) error
	GetOrderStatus(ctx context.Context, orderID, service string) (models.OrderStatusResponse, error)
	GetOrders(ctx context.Context, query dto.OrderSelector) ([]models.Order, int, error)
}

type OrderClientManager struct {
	orderCli order.Client
	storeCli store.Client
}

func NewOrderClientManager(orderCli order.Client, storeCli store.Client) OrderClient {
	return &OrderClientManager{
		orderCli: orderCli,
		storeCli: storeCli,
	}
}

func (manager *OrderClientManager) GetOrders(ctx context.Context, query dto.OrderSelector) ([]models.Order, int, error) {
	orders, total, err := manager.orderCli.GetOrdersWithFilters(ctx, query)
	if err != nil {
		log.Trace().Err(err).Msg("can't get orders from db")
		return nil, 0, err
	}

	return utils.ParseOrders(orders), total, nil
}

func (manager *OrderClientManager) GetOrder(ctx context.Context, orderID, service string) (models.Order, error) {
	order, err := manager.orderCli.GetOrder(ctx, dto.OrderSelector{
		ID: orderID,
	})
	if err != nil {
		log.Trace().Err(err).Msg("Can't get order from db")
		return models.Order{}, err
	}

	return utils.ParseOrder(order), nil
}

func (manager *OrderClientManager) UpdateOrder(ctx context.Context, req models.Order, orderID, service, clientSecret string) error {
	store, err := manager.storeCli.FindStore(ctx, storeModels.StoreSelector{
		DeliveryService: service,
		ClientSecret:    clientSecret,
	})
	if err != nil {
		log.Trace().Err(err).Msg("Can't find store by external store id")
		return err
	}

	order, err := req.ToModel(store, service)
	if err != nil {
		return err
	}

	order.ID = orderID

	err = manager.orderCli.UpdateOrder(ctx, order)
	if err != nil {
		log.Trace().Err(err).Msg("Can't update order")
		return err
	}

	return nil
}

func (manager *OrderClientManager) CancelOrder(ctx context.Context, req models.CancelOrderRequest, id, service, clientSecret string) error {
	_, err := manager.storeCli.FindStore(ctx, storeModels.StoreSelector{
		DeliveryService: service,
		ClientSecret:    clientSecret,
	})
	if err != nil {
		return err
	}

	order, err := manager.orderCli.GetOrder(ctx, dto.OrderSelector{
		ID: id,
	})
	if err != nil {
		return err
	}

	if order.DeliveryService != service {
		return errors.New("order is not exist")
	}

	if err := manager.orderCli.CancelOrder(ctx, coreModels.CancelOrder{
		ID:      id,
		OrderID: req.EatsId,
		Comment: req.Comment,
	}); err != nil {
		log.Trace().Err(err).Msg("Can't cancel order")
		return err
	}

	return nil
}

func (manager *OrderClientManager) GetOrderStatus(ctx context.Context, orderID, service string) (models.OrderStatusResponse, error) {
	order, err := manager.orderCli.GetOrder(ctx, dto.OrderSelector{
		ID: orderID,
	})
	if err != nil {
		log.Trace().Err(err).Msg("Can't get order from db")
		return models.OrderStatusResponse{}, err
	}

	return utils.ParseOrderStatus(order), nil
}
