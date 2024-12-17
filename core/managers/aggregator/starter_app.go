package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	"github.com/kwaaka-team/orders-core/pkg/starterapp"
	"github.com/kwaaka-team/orders-core/pkg/starterapp/clients"
	"github.com/kwaaka-team/orders-core/pkg/starterapp/clients/dto"
	"github.com/rs/zerolog/log"
	"time"
)

type StarterApp interface {
	UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error
	RejectOrder(ctx context.Context, orderID, reason string) error
	AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error
	AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error
	MarkOrder(ctx context.Context, orderID string) error
	ConfirmPreOrder(ctx context.Context, orderID string) error
	DeliveredOrder(ctx context.Context, orderID string) error
}

type StarterAppManager struct {
	starterAppCli clients.StarterApp
}

func NewStarterAppManager(globalCfg config.Configuration, apiKey string) (StarterApp, error) {
	client, err := starterapp.NewStarterAppClient(&clients.Config{
		Protocol: "http",
		BaseURL:  globalCfg.StarterAppConfiguration.BaseUrl,
		ApiKey:   apiKey,
	})
	if err != nil {
		log.Trace().Err(err).Msgf("can't initialize starter app client")
		return nil, err
	}

	return &StarterAppManager{
		starterAppCli: client,
	}, nil
}
func (m *StarterAppManager) UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error {
	if err := m.starterAppCli.ChangeOrderStatus(ctx, dto.ChangeOrderStatusRequest{
		Status: req.OrderStatus,
	}, req.OrderID); err != nil {
		return err
	}
	return nil
}

func (m *StarterAppManager) RejectOrder(ctx context.Context, orderID, reason string) error {
	if err := m.starterAppCli.SendOrderErrorNotification(ctx, dto.SendOrderErrorNotificationRequest{
		IsOrderSent: false,
		IsPosError:  true,
		Error: dto.SendOrderErrorNotificationError{
			Message:  reason,
			Request:  "",
			Response: "",
		},
	}, orderID); err != nil {
		return err
	}
	return nil
}

func (m *StarterAppManager) AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error {
	return nil
}
func (m *StarterAppManager) AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error {
	return nil
}

func (m *StarterAppManager) MarkOrder(ctx context.Context, orderID string) error {
	return nil
}
func (m *StarterAppManager) ConfirmPreOrder(ctx context.Context, orderID string) error {
	return nil
}

func (m *StarterAppManager) DeliveredOrder(ctx context.Context, orderID string) error {
	return nil
}
