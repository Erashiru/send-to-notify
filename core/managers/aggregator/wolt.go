package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/errors"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"time"

	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	woltClient "github.com/kwaaka-team/orders-core/pkg/wolt"
	woltConfig "github.com/kwaaka-team/orders-core/pkg/wolt/clients"
	woltModels "github.com/kwaaka-team/orders-core/pkg/wolt/clients/dto"
	"github.com/rs/zerolog/log"
)

type Wolt interface {
	AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error
	AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error
	RejectOrder(ctx context.Context, orderID, reason string) error
	MarkOrder(ctx context.Context, orderID string) error
	ConfirmPreOrder(ctx context.Context, orderID string) error
	DeliveredOrder(ctx context.Context, orderID string) error
	UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error
}

type WoltManager struct {
	woltClient woltConfig.Wolt
}

func NewWoltManager(globalConfig config.Configuration, store coreStoreModels.Store) (Wolt, error) {
	woltClient, err := woltClient.NewWoltClient(&woltConfig.Config{
		Protocol: "http",
		BaseURL:  globalConfig.WoltConfiguration.BaseURL,
		Username: store.Wolt.MenuUsername,
		Password: store.Wolt.MenuPassword,
		ApiKey:   store.Wolt.ApiKey,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Wolt Manager.")
		return nil, err
	}

	return &WoltManager{
		woltClient: woltClient,
	}, nil
}

func NewWoltManager2(baseUrl, username, password, apiKey string) (*WoltManager, error) {
	// Initialize new Wolt client
	woltClient, err := woltClient.NewWoltClient(&woltConfig.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		Username: username,
		Password: password,
		ApiKey:   apiKey,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Wolt Manager.")
		return nil, err
	}

	return &WoltManager{
		woltClient: woltClient,
	}, nil
}

func (w WoltManager) DeliveredOrder(ctx context.Context, orderID string) error {
	if err := w.woltClient.DeliveredOrder(ctx, orderID); err != nil {
		log.Trace().Err(err).Msgf("delivered order, order_id=%v", orderID)
		return err
	}

	return nil
}

func (w WoltManager) AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error {
	req := woltModels.AcceptSelfDeliveryOrderOrderRequest{
		ID: orderID,
	}
	if deliveryTime != nil {
		req.DeliveryTime = deliveryTime
	}
	if err := w.woltClient.AcceptSelfDeliveryOrder(ctx, req); err != nil {
		log.Trace().Err(err).Msgf("accept selfDelivery order, order_id=%v, deliveryTime=%v", orderID, deliveryTime)
		return err
	}

	log.Info().Msgf("success accept selfDelivery order, order_id=%v, deliveryTime=%v", orderID, deliveryTime)

	return nil
}

func (w WoltManager) AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error {
	req := woltModels.AcceptOrderRequest{
		ID: orderID,
	}
	if pickUpTime != nil {
		req.PickupTime = pickUpTime
	}
	if err := w.woltClient.AcceptOrder(ctx, req); err != nil {
		log.Trace().Err(err).Msgf("accept order, order_id=%v, pick_up_time=%v", orderID, pickUpTime)
		return err
	}

	log.Info().Msgf("success accept order, order_id=%v, pick_up_time=%v", orderID, pickUpTime)

	return nil
}

func (w WoltManager) RejectOrder(ctx context.Context, orderID, reason string) error {
	if err := w.woltClient.RejectOrder(ctx, woltModels.RejectOrderRequest{
		ID:     orderID,
		Reason: reason,
	}); err != nil {
		log.Trace().Err(err).Msgf("reject order, order_id=%v, reason=%v", orderID, reason)
		return err
	}

	log.Info().Msgf("success reject order, order_id=%v, reason=%v", orderID, reason)

	return nil
}

func (w WoltManager) MarkOrder(ctx context.Context, orderID string) error {
	if err := w.woltClient.MarkOrder(ctx, orderID); err != nil {
		log.Trace().Err(err).Msgf("mark order, order_id=%v", orderID)
		return err
	}

	log.Info().Msgf("success mark order, order_id=%v", orderID)

	return nil
}

func (w WoltManager) ConfirmPreOrder(ctx context.Context, orderID string) error {
	if err := w.woltClient.ConfirmPreOrder(ctx, orderID); err != nil {
		log.Trace().Err(err).Msgf("confirm pre-order, order_id=%v", orderID)
		return err
	}

	log.Info().Msgf("success confirm pre-order, order_id=%v", orderID)

	return nil
}

func (w WoltManager) UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error {
	return errors.ErrUnsupportedMethod
}
