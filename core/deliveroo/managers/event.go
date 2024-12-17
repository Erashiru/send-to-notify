package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/config"
	models2 "github.com/kwaaka-team/orders-core/core/deliveroo/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/pkg/deliveroo/clients"
	deliverooModels "github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/dto"
	deliverooHttpCli "github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/http"
	menuCoreCli "github.com/kwaaka-team/orders-core/pkg/menu"
	menuCoreModels "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	"time"

	orderCoreCli "github.com/kwaaka-team/orders-core/pkg/order"
	storeCore "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/pkg/store/dto"
)

type Event struct {
	storeCoreCli storeCore.Client
	orderCoreCli orderCoreCli.Client
	menuCoreCli  menuCoreCli.Client
}

func NewEvent(
	storeCoreCli storeCore.Client,
	orderCoreCli orderCoreCli.Client,
	menuCoreCli menuCoreCli.Client,
) Event {
	return Event{
		storeCoreCli: storeCoreCli,
		orderCoreCli: orderCoreCli,
		menuCoreCli:  menuCoreCli,
	}
}

func (man *Event) OrderEvent(ctx context.Context, req models2.OrderEvent) (string, error) {
	switch req.EventType {
	case models2.OrderCreate:
		err := man.createOrder(ctx, req.EventBody.Order)
		if err != nil {
			return "", err
		}
	case models2.OrderUpdate:
		if req.EventBody.Order.Status == "cancelled" {
			err := man.cancelOrder(ctx, req.EventBody.Order)
			if err != nil {
				return "", err
			}
		}
		err := man.updateOrder(ctx, req.EventBody.Order)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func (man *Event) createOrder(ctx context.Context, req models2.Order) error {
	store, err := man.storeCoreCli.FindStore(ctx, dto.StoreSelector{
		ExternalStoreID: req.LocationID,
		DeliveryService: models2.DELIVEROO.String(),
	})
	if err != nil {
		return err
	}
	conf, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}
	deliverooCli, err := deliverooHttpCli.NewClient(&clients.Config{
		BaseURL:  conf.DeliverooConfiguration.BaseURL,
		Username: store.Deliveroo.Username,
		Password: store.Deliveroo.Password,
	})
	if err != nil {
		return err
	}

	if !store.Deliveroo.SendToPos {
		err = deliverooCli.CreateSyncStatus(ctx, deliverooModels.CreateSyncStatusRequest{
			Status:     coreModels.FailedDeliveroo.String(),
			Reason:     coreModels.LocationNotSupported.String(),
			OccurredAt: time.Now().String(),
		}, req.ID)
		if err != nil {
			return err
		}
	}

	err = deliverooCli.CreateSyncStatus(ctx, deliverooModels.CreateSyncStatusRequest{
		Status:     coreModels.AcceptedDeliveroo.String(),
		OccurredAt: time.Now().String(),
	}, req.ID)

	if err != nil {
		return err
	}
	return nil
}

func (man *Event) updateOrder(ctx context.Context, req models2.Order) error {
	return nil
}

func (man *Event) cancelOrder(ctx context.Context, req models2.Order) error {
	_, err := man.storeCoreCli.FindStore(ctx, dto.StoreSelector{
		ID: req.LocationID,
	})
	if err != nil {
		return err
	}

	err = man.orderCoreCli.CancelOrderInPos(ctx, coreModels.CancelOrderInPos{
		OrderID:         req.ID,
		DeliveryService: string(models2.DELIVEROO),
		CancelReason:    coreModels.CancelReason{},
	})
	if err != nil {
		return err
	}
	return nil
}

func (man *Event) MenuEvent(ctx context.Context, req models2.MenuEvent) error {
	switch req.EventType {
	case models2.MenuUpload:
		var status string
		if req.EventBody.MenuUploadResult.Errors != nil {
			status = "have some errors" //need to define
		}

		var details string //need to define
		if str, ok := req.EventBody.MenuUploadResult.Errors.(string); ok {
			details = str
		}

		_, err := man.menuCoreCli.CreateMenuUploadTransaction(ctx, menuCoreModels.MenuUploadTransaction{
			ID:      req.EventBody.MenuUploadResult.MenuId,
			StoreID: req.EventBody.MenuUploadResult.BrandId,
			Status:  status,
			Service: models2.DELIVEROO.String(),
			CreatedAt: coreModels.TransactionTime{
				Value: coreModels.Time{
					Time: time.Now(),
				},
			},
			Details: []string{details},
		})
		if err != nil {
			return nil
		}

		return nil
	default:
		return nil
	}
}
