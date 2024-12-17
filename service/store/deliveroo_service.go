package store

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type deliverooService struct {
}

func newDeliverooService() (*deliverooService, error) {
	return &deliverooService{}, nil
}

func (s *deliverooService) GetStoreSchedulePrefix() string {
	return ""
}

func (s *deliverooService) IgnoreStatusUpdate(store models.Store) bool {
	return store.Deliveroo.IgnoreStatusUpdate
}

func (s *deliverooService) GetOrderCodePrefix(ctx context.Context, store models.Store) (string, error) {
	return store.Deliveroo.OrderCodePrefix, nil
}

func (s *deliverooService) IsSendToPos(store models.Store) (bool, error) {
	return true, nil
}

func (s *deliverooService) IsAutoAccept(store models.Store) (bool, error) {
	return false, nil
}
func (s *deliverooService) IsPostAutoAccept(store models.Store) (bool, error) {
	return false, nil
}

func (s *deliverooService) IsMarketplace(store models.Store) (bool, error) {
	return false, nil
}

func (s *deliverooService) GetPaymentTypes(store models.Store, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error) {
	return models.DeliveryServicePaymentType{}, nil
}

func (s *deliverooService) GetStoreExternalIds(store models.Store) ([]string, error) {
	return []string{store.Deliveroo.StoreID}, nil
}

func (s *deliverooService) UpdateSchedule(ctx context.Context, req models.UpdateStoreSchedule) error {
	return nil
}

func (s *deliverooService) IsSecretValid(store models.Store, secret string) (bool, error) {
	return true, nil
}
