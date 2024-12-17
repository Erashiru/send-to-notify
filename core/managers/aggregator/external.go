package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"time"

	externalClient "github.com/kwaaka-team/orders-core/pkg/externalapi"
	externalConfig "github.com/kwaaka-team/orders-core/pkg/externalapi/clients"
	externalModels "github.com/kwaaka-team/orders-core/pkg/externalapi/clients/dto"

	"github.com/kwaaka-team/orders-core/core/models/selector"
	"github.com/rs/zerolog/log"
)

type External interface {
	UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error
	AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error
	AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error
	RejectOrder(ctx context.Context, orderID, reason string) error
	MarkOrder(ctx context.Context, orderID string) error
	ConfirmPreOrder(ctx context.Context, orderID string) error
	DeliveredOrder(ctx context.Context, orderID string) error
}

type ExternalManager struct {
	externalClient externalConfig.Client
	webhookURL     string
}

func NewExternalManager(store coreStoreModels.Store, delivery string) (External, error) {
	var (
		authToken  string
		webhookURL string
	)

	for _, config := range store.ExternalConfig {
		if config.Type == delivery && config.WebhookURL != "" {
			authToken = config.AuthToken
			webhookURL = config.WebhookURL
			break
		}
	}

	if webhookURL == "" {
		log.Trace().Err(errors.ErrNoWebhookSubscription).Msgf("%s has not webhook subscription for %s delivery", store.Name, delivery)
		return nil, errors.ErrNoWebhookSubscription
	}

	externalClient, err := externalClient.NewWebhookClient(&externalConfig.Config{
		Protocol:  "http",
		AuthToken: authToken,
	})
	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize External Manager")
		return nil, err
	}

	return &ExternalManager{
		externalClient: externalClient,
		webhookURL:     webhookURL,
	}, nil
}

func (em ExternalManager) UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error {
	log.Info().Msgf("update order status, webhook url=%s, order_id=%s, order_status=%s", em.webhookURL, req.OrderID, req.OrderStatus)

	if err := em.externalClient.UpdateOrderWebhook(ctx, externalModels.Order{
		OrderID: req.OrderID,
		Status:  req.OrderStatus,
	}, em.webhookURL); err != nil {
		utils.Beautify("update order by webhook error", req)
		log.Trace().Err(err).Msgf("update order by webhook error")
		return err
	}

	return nil
}

func (em ExternalManager) AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error {
	return errors.ErrUnsupportedMethod
}
func (em ExternalManager) AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error {
	return errors.ErrUnsupportedMethod
}
func (em ExternalManager) RejectOrder(ctx context.Context, orderID, reason string) error {
	return errors.ErrUnsupportedMethod
}
func (em ExternalManager) MarkOrder(ctx context.Context, orderID string) error {
	return errors.ErrUnsupportedMethod
}
func (em ExternalManager) ConfirmPreOrder(ctx context.Context, orderID string) error {
	return errors.ErrUnsupportedMethod
}

func (em ExternalManager) DeliveredOrder(ctx context.Context, orderID string) error {
	return errors.ErrUnsupportedMethod
}
