package store

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type talabatStoreService struct {
}

func newTalabatStoreService() (*talabatStoreService, error) {
	return &talabatStoreService{}, nil
}

func (s *talabatStoreService) GetStoreSchedulePrefix() string {
	return ""
}

func (s *talabatStoreService) IgnoreStatusUpdate(store models.Store) bool {
	return store.Talabat.IgnoreStatusUpdate
}

func (s *talabatStoreService) GetOrderCodePrefix(ctx context.Context, store models.Store) (string, error) {
	return store.Talabat.OrderCodePrefix, nil
}

func (s *talabatStoreService) IsSendToPos(store models.Store) (bool, error) {
	return store.Talabat.SendToPos, nil
}

func (s *talabatStoreService) IsAutoAccept(store models.Store) (bool, error) {
	return false, nil
}
func (s *talabatStoreService) IsPostAutoAccept(store models.Store) (bool, error) {
	return false, nil
}

func (s *talabatStoreService) IsMarketplace(store models.Store) (bool, error) {
	return store.Talabat.IsMarketplace, nil
}

func (s *talabatStoreService) GetPaymentTypes(store models.Store, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error) {
	return store.Talabat.PaymentTypes, nil
}

func (s *talabatStoreService) GetStoreExternalIds(store models.Store) ([]string, error) {
	return store.Talabat.BranchID, nil
}

func (s *talabatStoreService) UpdateSchedule(ctx context.Context, req models.UpdateStoreSchedule) error {
	return nil
}

func (s *talabatStoreService) IsSecretValid(store models.Store, secret string) (bool, error) {
	return true, nil
}
