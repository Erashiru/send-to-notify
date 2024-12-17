package store

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type kwaakaAdminStoreService struct {
}

func newKwaakaAdminStoreService() (*kwaakaAdminStoreService, error) {
	return &kwaakaAdminStoreService{}, nil
}

func (s *kwaakaAdminStoreService) GetStoreSchedulePrefix() string {
	return ""
}

func (s *kwaakaAdminStoreService) GetOrderCodePrefix(ctx context.Context, store models.Store) (string, error) {
	return "", nil
}

func (s *kwaakaAdminStoreService) IsSendToPos(store models.Store) (bool, error) {
	return store.KwaakaAdmin.SendToPos, nil
}

func (s *kwaakaAdminStoreService) IsAutoAccept(store models.Store) (bool, error) {
	return false, nil
}
func (s *kwaakaAdminStoreService) IsPostAutoAccept(store models.Store) (bool, error) {
	return false, nil
}

func (s *kwaakaAdminStoreService) IgnoreStatusUpdate(store models.Store) bool {
	return false
}

func (s *kwaakaAdminStoreService) IsMarketplace(store models.Store) (bool, error) {
	return true, nil
}

func (s *kwaakaAdminStoreService) GetPaymentTypes(store models.Store, paymentInfo models2.PosPaymentInfo) (models.DeliveryServicePaymentType, error) {
	return models.DeliveryServicePaymentType{
		DELAYED: models.PaymentType{
			PaymentTypeID:   paymentInfo.PaymentTypeID,
			PaymentTypeKind: paymentInfo.PaymentTypeKind,
		},
	}, nil
}

func (s *kwaakaAdminStoreService) GetStoreExternalIds(store models.Store) ([]string, error) {
	return store.KwaakaAdmin.StoreID, nil
}

func (s *kwaakaAdminStoreService) UpdateSchedule(ctx context.Context, req models.UpdateStoreSchedule) error {
	return nil
}

func (s *kwaakaAdminStoreService) IsSecretValid(store models.Store, secret string) (bool, error) {
	return true, nil
}
