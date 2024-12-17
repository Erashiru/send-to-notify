package order

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
)

type CreateOrderRouter struct {
	storeService    store.Service
	serviceWithPos  CreationService
	serviceEmptyPos CreationService
}

func NewCreateOrderRouter(storeService store.Service, serviceWithPos CreationService, serviceEmptyPos CreationService) (*CreateOrderRouter, error) {
	if storeService == nil {
		return nil, errors.New("storeService is nil")
	}
	if serviceWithPos == nil {
		return nil, errors.New("service is nil")
	}
	if serviceEmptyPos == nil {
		return nil, errors.New("service is nil")
	}
	return &CreateOrderRouter{
		storeService:    storeService,
		serviceWithPos:  serviceWithPos,
		serviceEmptyPos: serviceEmptyPos,
	}, nil
}

func (s *CreateOrderRouter) CreateOrder(ctx context.Context, externalStoreID, deliveryService string, aggReq interface{}, storeSecret string) (models.Order, error) {
	st, err := s.storeService.GetByExternalIdAndDeliveryService(ctx, externalStoreID, deliveryService)
	if err != nil {
		return models.Order{}, err
	}

	if st.PosType == models.Kwaaka.String() {
		return s.serviceEmptyPos.CreateOrder(ctx, externalStoreID, deliveryService, aggReq, storeSecret)
	} else {
		return s.serviceWithPos.CreateOrder(ctx, externalStoreID, deliveryService, aggReq, storeSecret)
	}
}
