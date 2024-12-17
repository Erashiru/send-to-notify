package scheduledupdatestatus

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	glovoClient "github.com/kwaaka-team/orders-core/pkg/glovo/clients"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/pkg/errors"
)

type Base interface {
	UpdateStatusToReady(ctx context.Context, store storeModels.Store) error
}

type DeliveryServices struct {
	glovoService Glovo
}

func NewDeliveryService(glovoCli glovoClient.Glovo, orderCli *order.OrderCoreClient) (DeliveryServices, error) {
	glovoService, err := NewGlovoService(glovoCli, orderCli)
	if err != nil {
		return DeliveryServices{}, err
	}
	return DeliveryServices{
		glovoService: glovoService,
	}, nil
}

func (b DeliveryServices) factory(deliveryService string) (Base, error) {
	switch deliveryService {
	case models.GLOVO.String():
		return b.glovoService, nil
	default:
		return nil, errors.New("that delivery service not supported")
	}
}

func (b DeliveryServices) UpdateStatusToReady(ctx context.Context, store storeModels.Store) error {
	if len(store.Settings.ScheduledStatusChange.DeliveryServices) == 0 {
		return errors.New("no delivery service to change status to ready")
	}
	for _, v := range store.Settings.ScheduledStatusChange.DeliveryServices {
		deliveryService, err := b.factory(v)
		if err != nil {
			return err
		}
		if err = deliveryService.UpdateStatusToReady(ctx, store); err != nil {
			return err
		}
	}
	return nil
}
