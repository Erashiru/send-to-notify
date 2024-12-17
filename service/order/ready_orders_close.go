package order

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models"
	"time"
)

type OrderStatusReadyToClose interface {
	UpdateOrdersWithTwoMoreHours(ctx context.Context) error
}

func (s *ServiceImpl) UpdateOrdersWithTwoMoreHours(ctx context.Context) error {

	nonActiveOrderStatuses = append(nonActiveOrderStatuses, models.COOKING_COMPLETE.String())

	orders, err := s.repository.GetOrdersByStatusesAndPosType(ctx, models.Kwaaka.String(), nonActiveOrderStatuses)
	if err != nil {
		return err
	}

	for i := range orders {
		if !s.isStatusMoreThatTwoHours(orders[i]) {
			continue
		}
		if err := s.closeOrder(ctx, orders[i]); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) isStatusMoreThatTwoHours(order models.Order) bool {

	currentTime := time.Now()

	twoHourDuration := 2 * time.Hour

	return currentTime.Sub(order.StatusesHistory[len(order.StatusesHistory)-1].Time) >= twoHourDuration && order.StatusesHistory[len(order.StatusesHistory)-1].Name == order.Status
}

func (s *ServiceImpl) closeOrder(ctx context.Context, order models.Order) error {

	if order.Status == models.COOKING_COMPLETE.String() {
		if err := s.repository.UpdateOrderStatusByID(ctx, order.ID, models.CLOSED.String()); err != nil {
			return err
		}

		if err := s.repository.AddStatusToHistory(ctx, order.OrderID, models.CLOSED.String()); err != nil {
			return err
		}
	} else {
		if err := s.repository.UpdateOrderStatusByID(ctx, order.ID, order.Status); err != nil {
			return err
		}
	}
	return nil
}
