package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	"time"
)

type KwaakaAdmin interface {
	AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error
	AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error
	RejectOrder(ctx context.Context, orderID, reason string) error
	MarkOrder(ctx context.Context, orderID string) error
	ConfirmPreOrder(ctx context.Context, orderID string) error
	DeliveredOrder(ctx context.Context, orderID string) error
	UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error
}

type KwaakaAdminManager struct {
}

func NewKwaakaAdminManager() (*KwaakaAdminManager, error) {
	return &KwaakaAdminManager{}, nil
}

func (k KwaakaAdminManager) AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error {
	return nil
}

func (k KwaakaAdminManager) AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error {
	return nil
}
func (k KwaakaAdminManager) RejectOrder(ctx context.Context, orderID, reason string) error {
	return nil
}
func (k KwaakaAdminManager) MarkOrder(ctx context.Context, orderID string) error {
	return nil
}
func (k KwaakaAdminManager) ConfirmPreOrder(ctx context.Context, orderID string) error {
	return nil
}
func (k KwaakaAdminManager) DeliveredOrder(ctx context.Context, orderID string) error {
	return nil
}
func (k KwaakaAdminManager) UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error {
	return nil
}
