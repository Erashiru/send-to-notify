package storeStatus

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/store/repository/storeclosedtime"
	"time"
)

type DatastoreObserver struct {
	DatastoreClient storeclosedtime.Repository
}

func (s DatastoreObserver) Notify(ctx context.Context, restaurant models.Store, externalStoreID, deliveryService string, storeIsOpened bool) error {
	err := s.DatastoreClient.Insert(ctx, models.StoreActiveTime{
		RestaurantID:    restaurant.ID,
		StoreID:         externalStoreID,
		DeliveryService: deliveryService,
		StartTime:       time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s DatastoreObserver) NotifyStatusReport(ctx context.Context, restaurant models.Store, durations []models.OpenTimeDuration) error {
	return nil
}

func (s DatastoreObserver) NotifyStatusChange(ctx context.Context, status string, phone string) error {
	return nil
}
