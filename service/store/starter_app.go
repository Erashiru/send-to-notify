package store

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type starterAppService struct {
}

func newStarterAppService() (*starterAppService, error) {
	return &starterAppService{}, nil
}

func (s *starterAppService) GetStoreSchedulePrefix() string {
	return ""
}

func (s *starterAppService) IgnoreStatusUpdate(store models.Store) bool {
	return store.StarterApp.IgnoreStatusUpdate
}

func (s *starterAppService) GetOrderCodePrefix(ctx context.Context, store models.Store) (string, error) {
	return "", nil
}

func (s *starterAppService) IsSendToPos(store models.Store) (bool, error) {
	return store.StarterApp.SendToPos, nil
}

func (s *starterAppService) IsAutoAccept(store models.Store) (bool, error) {
	return false, nil
}
func (s *starterAppService) IsPostAutoAccept(store models.Store) (bool, error) {
	return false, nil
}

func (s *starterAppService) IsMarketplace(store models.Store) (bool, error) {
	return store.StarterApp.IsMarketPlace, nil
}

func (s *starterAppService) GetPaymentTypes(store models.Store, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error) {
	return store.StarterApp.PaymentTypes, nil
}

func (s *starterAppService) GetStoreExternalIds(store models.Store) ([]string, error) {
	return []string{store.StarterApp.ShopID}, nil
}

func (s *starterAppService) UpdateSchedule(ctx context.Context, req models.UpdateStoreSchedule) error {
	return nil
}

func (s *starterAppService) IsSecretValid(store models.Store, secret string) (bool, error) {
	return true, nil
}
