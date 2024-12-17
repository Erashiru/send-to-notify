package store

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type glovoStoreService struct{}

func newGlovoStoreService() (*glovoStoreService, error) {
	return &glovoStoreService{}, nil
}

func (s *glovoStoreService) GetStoreSchedulePrefix() string {
	return "store_schedule.glovo_schedule"
}

func (s *glovoStoreService) IgnoreStatusUpdate(store models.Store) bool {
	return store.Glovo.IgnoreStatusUpdate
}

func (s *glovoStoreService) IsSendToPos(store models.Store) (bool, error) {
	return store.Glovo.SendToPos, nil
}

func (s *glovoStoreService) GetOrderCodePrefix(ctx context.Context, store models.Store) (string, error) {
	return store.Glovo.OrderCodePrefix, nil
}

func (s *glovoStoreService) IsAutoAccept(store models.Store) (bool, error) {
	return store.Glovo.AutoAcceptOn, nil
}
func (s *glovoStoreService) IsPostAutoAccept(store models.Store) (bool, error) {
	return store.Glovo.PostAutoAcceptOn, nil
}

func (s *glovoStoreService) IsMarketplace(store models.Store) (bool, error) {
	return store.Glovo.IsMarketplace, nil
}

func (s *glovoStoreService) GetPaymentTypes(store models.Store, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error) {
	return store.Glovo.PaymentTypes, nil
}

func (s *glovoStoreService) GetStoreExternalIds(store models.Store) ([]string, error) {
	return store.Glovo.StoreID, nil
}

func (s *glovoStoreService) IsSecretValid(store models.Store, secret string) (bool, error) {
	return true, nil
}
