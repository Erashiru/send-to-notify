package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	"time"
)

type QRMenu interface {
	AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error
	AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error
	RejectOrder(ctx context.Context, orderID, reason string) error
	MarkOrder(ctx context.Context, orderID string) error
	ConfirmPreOrder(ctx context.Context, orderID string) error
	DeliveredOrder(ctx context.Context, orderID string) error
	UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error
}

type QRMenuManager struct {
}

func NewQRMenuManager() (*QRMenuManager, error) {
	return &QRMenuManager{}, nil
}

func (k QRMenuManager) AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error {
	return nil
}

func (k QRMenuManager) AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error {
	return nil
}
func (k QRMenuManager) RejectOrder(ctx context.Context, orderID, reason string) error {
	return nil
}
func (k QRMenuManager) MarkOrder(ctx context.Context, orderID string) error {
	return nil
}
func (k QRMenuManager) ConfirmPreOrder(ctx context.Context, orderID string) error {
	return nil
}
func (k QRMenuManager) DeliveredOrder(ctx context.Context, orderID string) error {
	return nil
}
func (k QRMenuManager) UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error {
	return nil
}
