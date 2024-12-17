package store

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type woltStoreService struct {
}

func newWoltStoreService() (*woltStoreService, error) {
	return &woltStoreService{}, nil
}

func (s *woltStoreService) GetStoreSchedulePrefix() string {
	return "store_schedule.wolt_schedule"
}

func (s *woltStoreService) IgnoreStatusUpdate(store models.Store) bool {
	return store.Wolt.IgnoreStatusUpdate
}

func (s *woltStoreService) GetOrderCodePrefix(ctx context.Context, store models.Store) (string, error) {
	return store.Wolt.OrderCodePrefix, nil
}

func (s *woltStoreService) IsSendToPos(store models.Store) (bool, error) {
	return store.Wolt.SendToPos, nil
}

func (s *woltStoreService) IsAutoAccept(store models.Store) (bool, error) {
	return store.Wolt.AutoAcceptOn, nil
}

func (s *woltStoreService) IsPostAutoAccept(store models.Store) (bool, error) {
	return store.Wolt.PostAutoAcceptOn, nil
}

func (s *woltStoreService) IsMarketplace(store models.Store) (bool, error) {
	return store.Wolt.IsMarketplace, nil
}

func (s *woltStoreService) GetPaymentTypes(store models.Store, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error) {
	return store.Wolt.PaymentTypes, nil
}

func (s *woltStoreService) GetStoreExternalIds(store models.Store) ([]string, error) {
	return store.Wolt.StoreID, nil
}

func (s *woltStoreService) UpdateSchedule(ctx context.Context, req models.UpdateStoreSchedule) error {
	return nil
}

func (s *woltStoreService) IsSecretValid(store models.Store, secret string) (bool, error) {
	return true, nil
}
