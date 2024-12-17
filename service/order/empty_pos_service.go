package order

import (
	"context"
	"fmt"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/custom"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	newPos "github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type EmptyPosService struct {
	storeService      store.Service
	aggregatorFactory aggregator.Factory
	repository        Repository
}

func NewEmptyPosService(
	storeService store.Service,
	aggregatorFactory aggregator.Factory,
	repository Repository) (*EmptyPosService, error) {
	if storeService == nil {
		return nil, errors.New("store service is nil")
	}
	if aggregatorFactory == nil {
		return nil, errors.New("aggregator factory is nil")
	}
	if repository == nil {
		return nil, errors.New("repository is nil")
	}

	return &EmptyPosService{
		storeService:      storeService,
		aggregatorFactory: aggregatorFactory,
		repository:        repository,
	}, nil
}

func (s *EmptyPosService) CreateOrder(ctx context.Context, externalStoreID, deliveryService string, aggReq interface{}, storeSecret string) (models.Order, error) {

	st, err := s.storeService.GetByExternalIdAndDeliveryService(ctx, externalStoreID, deliveryService)
	if err != nil {
		return models.Order{}, err
	}
	if st.PosType != models.Kwaaka.String() {
		return models.Order{}, errors.New("pos type validation error")
	}

	utils.Beautify("getting store body", st)

	isSecretValid, err := s.storeService.IsSecretValid(st, deliveryService, storeSecret)
	if err != nil {
		return models.Order{}, err
	}

	if !isSecretValid {
		return models.Order{}, errors.New("store secret is not valid")
	}

	agg, err := s.aggregatorFactory.GetAggregator(deliveryService, st)
	if err != nil {
		return models.Order{}, err
	}

	if st.Settings.HasVirtualStore {
		return s.splitVirtualStoreOrder(ctx, st, aggReq, agg, deliveryService)
	}

	req, err := agg.GetSystemCreateOrderRequestByAggregatorRequest(aggReq, st)
	if err != nil {
		return models.Order{}, err
	}

	utils.Beautify("aggregator request to system request", req)

	aggregatorRequestBody, err := utils.GetJsonFormatFromModel(aggReq)
	if err != nil {
		return models.Order{}, err
	}
	req = fillRequestData(req, st, aggregatorRequestBody)

	order, err := s.saveOrderToDb(ctx, req)
	if err != nil {
		return order, err
	}

	if order.Type == "PREORDER" {
		order, err = s.waitSendingOrder(ctx, order, st)
		if err != nil {
			return order, err
		}
	}

	isAutoAcceptOn, err := s.storeService.IsAutoAccept(st, deliveryService)
	if err != nil {
		return order, err
	}

	if isAutoAcceptOn {
		status := models.ACCEPTED
		if err = s.updateStatus(ctx, agg, order, st, status); err != nil {
			return models.Order{}, err
		} else {
			order.Status = status.String()
		}
	}

	return order, nil
}

func (s *EmptyPosService) updateStatus(ctx context.Context, agg aggregator.Aggregator, order models.Order, st storeModels.Store, status models.PosStatus) error {
	aggStatus := agg.MapSystemStatusToAggregatorStatus(order, status, st)

	if err := agg.UpdateOrderInAggregator(ctx, order, st, aggStatus); err != nil {
		return err
	}

	if err := s.repository.UpdateOrderStatusByID(ctx, order.ID, status.String()); err != nil {
		return err
	}

	return nil
}

func (s *EmptyPosService) saveOrderToDb(ctx context.Context, req models.Order) (models.Order, error) {
	order, err := s.repository.InsertOrder(ctx, req)
	if err != nil {
		log.Err(err).Msgf("orders core, insert order error")

		if errors.Is(err, errs.ErrAlreadyExist) {
			log.Info().Msg("Order already exist, skipping...")
			req.FailReason.Code = newPos.ORDER_ALREADY_EXIST_CODE
			req.FailReason.Message = newPos.ORDER_ALREADY_EXIST
			return models.Order{}, errors.Wrap(validator.ErrPassed, fmt.Sprintf("order %s passed", req.OrderID))
		}

		return s.failOrder(ctx, order, err)
	}

	return order, nil
}

func (s *EmptyPosService) failOrder(ctx context.Context, req models.Order, err error) (models.Order, error) {

	var errs custom.Error

	log.Trace().Err(err).Msgf("%v", validator.ErrFailed)

	req.Status = string(models.STATUS_FAILED)

	req.StatusesHistory = append(req.StatusesHistory, models.OrderStatusUpdate{
		Name: string(models.STATUS_FAILED),
		Time: models.TimeNow().Time,
	})
	if updateErr := s.repository.UpdateOrder(ctx, req); updateErr != nil {
		return req, validator.ErrFailed
	}

	if err != nil {
		log.Trace().Err(err).Msg("Error while saving order")
		return req, err
	}

	errs.Append(err, validator.ErrFailed)

	return req, errs
}

func (s *EmptyPosService) waitSendingOrder(ctx context.Context, order models.Order, store storeModels.Store) (models.Order, error) {
	log.Info().Msgf("order waiting sending, status: %v", string(models.STATUS_WAIT_SENDING))

	systemStatus := models.WAIT_SENDING

	if err := s.repository.UpdateOrderStatusByID(ctx, order.ID, systemStatus.String()); err != nil {
		return order, err
	}
	order.Status = systemStatus.String()

	return order, nil
}

func (s *EmptyPosService) splitVirtualStoreOrder(ctx context.Context, store storeModels.Store, req interface{}, agg aggregator.Aggregator, deliveryService string) (models.Order, error) {
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
		status := models.ACCEPTED
		if err = s.updateStatus(ctx, agg, order, store, status); err != nil {
			return models.Order{}, err
		} else {
			order.Status = status.String()
		}
	}

	return order, nil
}
