package store

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type express24Service struct {
}

func newExpress24Service() (*express24Service, error) {
	return &express24Service{}, nil
}

func (s *express24Service) GetStoreSchedulePrefix() string {
	return ""
}

func (s *express24Service) IgnoreStatusUpdate(store models.Store) bool {
	return store.Express24.IgnoreStatusUpdate
}

func (s *express24Service) GetOrderCodePrefix(ctx context.Context, store models.Store) (string, error) {
	return store.Express24.OrderCodePrefix, nil
}

func (s *express24Service) IsSendToPos(store models.Store) (bool, error) {
	return store.Express24.SendToPos, nil
}

func (s *express24Service) IsAutoAccept(store models.Store) (bool, error) {
	return false, nil
}
func (s *express24Service) IsPostAutoAccept(store models.Store) (bool, error) {
	return false, nil
}

func (s *express24Service) IsMarketplace(store models.Store) (bool, error) {
	return store.Express24.IsMarketplace, nil
}

func (s *express24Service) GetPaymentTypes(store models.Store, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error) {
	return store.Express24.PaymentTypes, nil
}

func (s *express24Service) GetStoreExternalIds(store models.Store) ([]string, error) {
	return store.Express24.StoreID, nil
}

func (s *express24Service) UpdateSchedule(ctx context.Context, req models.UpdateStoreSchedule) error {
	return nil
}

func (s *express24Service) IsSecretValid(store models.Store, secret string) (bool, error) {
	return true, nil
}
