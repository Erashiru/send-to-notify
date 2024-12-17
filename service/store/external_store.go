package store

import (
	"context"
	"fmt"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	models3 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
)

var deliveryServiceNotFound = errors.New("delivery service not found")

var supportedDeliveryServices = map[models2.AggregatorName]interface{}{
	models2.YANDEX: struct{}{},
	models2.EMENU:  struct{}{},
}

type externalStoreService struct {
	deliveryService string
}

func newExternalStoreService(deliveryService string) (*externalStoreService, error) {
	if _, ok := supportedDeliveryServices[models2.AggregatorName(deliveryService)]; !ok {
		return nil, errors.Wrap(deliveryServiceNotFound, fmt.Sprintf("delivery service %s not found", deliveryService))
	}

	return &externalStoreService{
		deliveryService,
	}, nil
}

func (s *externalStoreService) GetStoreSchedulePrefix() string {
	return ""
}

func (s *externalStoreService) IgnoreStatusUpdate(store models.Store) bool {
	for _, item := range store.ExternalConfig {
		if item.Type != s.deliveryService {
			continue
		}
		return item.IgnoreStatusUpdate
	}

	return false
}

func (s *externalStoreService) GetOrderCodePrefix(ctx context.Context, store models.Store) (string, error) {
	for _, item := range store.ExternalConfig {
		if item.Type != s.deliveryService {
			continue
		}
		return item.OrderCodePrefix, nil
	}

	return "", errors.Wrap(deliveryServiceNotFound, "order code prefix")
}

func (s *externalStoreService) IsSendToPos(store models.Store) (bool, error) {
	for _, item := range store.ExternalConfig {
		if item.Type != s.deliveryService {
			continue
		}
		return item.SendToPos, nil
	}

	return false, errors.Wrap(deliveryServiceNotFound, "is send to pos")
}

func (s *externalStoreService) IsAutoAccept(store models.Store) (bool, error) {
	for _, item := range store.ExternalConfig {
		if item.Type != s.deliveryService {
			continue
		}
		return item.AutoAcceptOn, nil
	}

	return false, errors.Wrap(deliveryServiceNotFound, "is auto accept")
}
func (s *externalStoreService) IsPostAutoAccept(store models.Store) (bool, error) {
	for _, item := range store.ExternalConfig {
		if item.Type != s.deliveryService {
			continue
		}
		return item.PostAutoAcceptOn, nil
	}

	return false, errors.Wrap(deliveryServiceNotFound, "is auto accept")
}

func (s *externalStoreService) IsMarketplace(store models.Store) (bool, error) {
	for _, item := range store.ExternalConfig {
		if item.Type != s.deliveryService {
			continue
		}
		return item.IsMarketplace, nil
	}

	return false, errors.Wrap(deliveryServiceNotFound, "is marketplace")
}

func (s *externalStoreService) GetPaymentTypes(store models.Store, paymentInfo models3.PosPaymentInfo) (models.DeliveryServicePaymentType, error) {
	for _, item := range store.ExternalConfig {
		if item.Type != s.deliveryService {
			continue
		}
		return item.PaymentTypes, nil
	}

	return models.DeliveryServicePaymentType{}, errors.Wrap(deliveryServiceNotFound, "payment types")
}

func (s *externalStoreService) GetStoreExternalIds(store models.Store) ([]string, error) {
	for _, item := range store.ExternalConfig {
		if item.Type != s.deliveryService {
			continue
		}
		return item.StoreID, nil
	}

	return nil, errors.Wrap(deliveryServiceNotFound, "store external ids")

}

func (s *externalStoreService) UpdateSchedule(ctx context.Context, req models.UpdateStoreSchedule) error {
	return nil
}

func (s *externalStoreService) IsSecretValid(store models.Store, secret string) (bool, error) {

	for i := range store.ExternalConfig {
		if store.ExternalConfig[i].Type != s.deliveryService {
			continue
		}
		return store.ExternalConfig[i].ClientSecret == secret, nil
	}
	return false, errors.Wrap(deliveryServiceNotFound, "secret validation")
}
