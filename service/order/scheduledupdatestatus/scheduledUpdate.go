package scheduledupdatestatus

import (
	"context"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/pkg/store/dto"
)

type ScheduledUpdate interface {
	UpdateToReady(ctx context.Context, store storeModels.Store) error
}

type ScheduledUpdateStatus struct {
	storeCli        store.Client
	deliveryService DeliveryServices
}

func New(storeCli store.Client, deliveryService DeliveryServices) ScheduledUpdateStatus {
	return ScheduledUpdateStatus{
		storeCli:        storeCli,
		deliveryService: deliveryService,
	}
}

func (s ScheduledUpdateStatus) UpdateToReady(ctx context.Context) error {
	stores, err := s.storeCli.FindStores(ctx, dto.StoreSelector{
		HasScheduledStatusChange: true,
	})
	if err != nil {
		return err
	}
	for _, store := range stores {
		if err = s.deliveryService.UpdateStatusToReady(ctx, store); err != nil {
			return err
		}
	}

	return nil
}
