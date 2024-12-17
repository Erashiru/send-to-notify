package scheduledupdatestatus

import (
	"context"
	"errors"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	orderModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	glovoClient "github.com/kwaaka-team/orders-core/pkg/glovo/clients"
	glovoModels "github.com/kwaaka-team/orders-core/pkg/glovo/clients/dto"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/order/dto"
	"strconv"
	"time"
)

type Glovo struct {
	glovoCli glovoClient.Glovo
	orderCli *order.OrderCoreClient
}

func NewGlovoService(glovoCli glovoClient.Glovo, orderCLi *order.OrderCoreClient) (Glovo, error) {
	if glovoCli == nil {
		return Glovo{}, errors.New("glovo client is empty")
	}

	if orderCLi == nil {
		return Glovo{}, errors.New("orderCLi client is empty")
	}
	return Glovo{
		glovoCli: glovoCli,
		orderCli: orderCLi,
	}, nil
}

func (g Glovo) UpdateStatusToReady(ctx context.Context, store storeModels.Store) error {
	orderTimeFrom := time.Now().Add(-2 * time.Hour)
	orders, _, err := g.orderCli.GetOrdersWithFilters(ctx, dto.OrderSelector{
		DeliveryService: models.GLOVO.String(),
		OnlyActive:      true,
		OrderTimeFrom:   orderTimeFrom,
		StoreID:         store.ID,
	})

	if err != nil {
		return err
	}
	for _, v := range orders {
		if g.needToUpdateStatus(v, store.Settings.ScheduledStatusChange.SwitchInterval) {
			if err = g.updateInDBAggregator(ctx, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g Glovo) needToUpdateStatus(order orderModels.Order, timeNeedToBeUpdate int) bool {
	currentState := order.Status == "ACCEPTED" || order.Status == "COOKING_STARTED" || order.Status == "WAIT_SENDING"
	orderTime := order.OrderTime.Value.Time
	currentTime := time.Now().UTC()
	if currentState && int(currentTime.Sub(orderTime).Minutes()) > timeNeedToBeUpdate {
		return true
	}
	return false
}

func (g Glovo) updateInDBAggregator(ctx context.Context, order orderModels.Order) error {
	formatedOrderID, err := strconv.Atoi(order.OrderID)
	if err != nil {
		return err
	}
	if _, err = g.glovoCli.UpdateOrderStatus(ctx, glovoModels.OrderUpdateRequest{
		StoreID: order.StoreID,
		Status:  "READY_FOR_PICKUP",
		ID:      int64(formatedOrderID),
	}); err != nil {
		return err
	}
	if err = g.orderCli.UpdateOrderStatusByID(ctx, order.ID, order.PosType, "READY_FOR_PICKUP"); err != nil {
		return err
	}
	return nil
}
