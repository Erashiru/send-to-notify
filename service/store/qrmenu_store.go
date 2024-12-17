package store

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type qrMenuStoreService struct {
}

func newQrMenuStoreService() (*qrMenuStoreService, error) {
	return &qrMenuStoreService{}, nil
}

func (s *qrMenuStoreService) GetStoreSchedulePrefix() string {
	return ""
}

func (s *qrMenuStoreService) GetOrderCodePrefix(ctx context.Context, store models.Store) (string, error) {
	return "", nil
}

func (s *qrMenuStoreService) IsSendToPos(store models.Store) (bool, error) {
	return store.QRMenu.SendToPos, nil
}

func (s *qrMenuStoreService) IsAutoAccept(store models.Store) (bool, error) {
	return false, nil
}
func (s *qrMenuStoreService) IsPostAutoAccept(store models.Store) (bool, error) {
	return false, nil
}

func (s *qrMenuStoreService) IgnoreStatusUpdate(store models.Store) bool {
	return store.QRMenu.IgnoreStatusUpdate
}

func (s *qrMenuStoreService) IsMarketplace(store models.Store) (bool, error) {
	return store.QRMenu.IsMarketplace, nil
}

func (s *qrMenuStoreService) GetPaymentTypes(store models.Store, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error) {
	return store.QRMenu.PaymentTypes, nil
}

func (s *qrMenuStoreService) GetStoreExternalIds(store models.Store) ([]string, error) {
	return store.QRMenu.StoreID, nil
}

func (s *qrMenuStoreService) UpdateSchedule(ctx context.Context, req models.UpdateStoreSchedule) error {
	return nil
}

func (s *qrMenuStoreService) IsSecretValid(store models.Store, secret string) (bool, error) {
	return true, nil
}
