package order

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
)

func (s *ServiceImpl) splitVirtualStoreOrder(ctx context.Context, store storeModels.Store, req interface{}, agg aggregator.Aggregator, deliveryService string) (models.Order, error) {
	systemOrderRequest, err := agg.GetSystemCreateOrderRequestByAggregatorRequest(req, store)
	if err != nil {
		return models.Order{}, err
	}
	systemOrderRequest.PosType = store.PosType
	systemOrderRequest.RestaurantName = store.Name
	systemOrderRequest.RestaurantID = store.ID
	systemOrderRequest.IsParentOrder = true
	order, err := s.saveOrderToDb(ctx, systemOrderRequest)
	if err != nil {
		return order, err
	}

	childRestaurantOrders, err := agg.SplitVirtualStoreOrder(req, store)
	if err != nil {
		log.Err(err).Msg("split aggregator order error")
		return models.Order{}, err
	}

	for _, childOrder := range childRestaurantOrders {
		aggStoreID, err := agg.GetStoreIDFromAggregatorOrderRequest(childOrder)
		if err != nil {
			log.Err(err).Msg("get aggregator store id from child order error")
			continue
		}
		_, err = s.CreateOrder(ctx, aggStoreID, deliveryService, childOrder, "")
		if err != nil {
			log.Err(err).Msg("create child order error")
			continue
		}
	}

	isAutoAcceptOn, err := s.storeService.IsAutoAccept(store, deliveryService)
	if err != nil {
		return order, err
	}

	if isAutoAcceptOn {
		order.Status = models.ACCEPTED.String()

		aggStatus := agg.MapSystemStatusToAggregatorStatus(order, models.ACCEPTED, store)

		if err = agg.UpdateOrderInAggregator(ctx, order, store, aggStatus); err != nil {
			return models.Order{}, err
		}

		if err = s.repository.UpdateOrderStatusByID(ctx, order.ID, models.ACCEPTED.String()); err != nil {
			return models.Order{}, err
		}
	}

	return order, nil
}

func (s *ServiceImpl) splitVirtualStoreItemID(id string, sep string) (string, string, error) {
	result := strings.Split(id, sep)
	if len(result) != 2 && len(result) != 3 {
		return "", "", errors.New("not valid signature, len will be 2 with _")
	}

	return result[0], result[1], nil
}
