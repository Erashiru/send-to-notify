package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/models"
	"time"

	"github.com/kwaaka-team/orders-core/core/models/selector"
	talabatClient "github.com/kwaaka-team/orders-core/pkg/talabat"
	talabatConfig "github.com/kwaaka-team/orders-core/pkg/talabat/clients"
	talabatModels "github.com/kwaaka-team/orders-core/pkg/talabat/models"
	"github.com/rs/zerolog/log"
)

type Talabat interface {
	AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error
	AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error
	RejectOrder(ctx context.Context, orderID, reason string) error
	MarkOrder(ctx context.Context, orderID string) error
	ConfirmPreOrder(ctx context.Context, orderID string) error
	DeliveredOrder(ctx context.Context, orderID string) error
	UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error
}

type TalabatManager struct {
	talabatClient talabatConfig.TalabatMW
}

func NewTalabatManager(baseUrl, username, password string) (*TalabatManager, error) {
	talabatCli, err := talabatClient.NewMiddlewareClient(&talabatConfig.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Talabat Manager.")
		return nil, err
	}

	return &TalabatManager{
		talabatClient: talabatCli,
	}, nil
}

func (t TalabatManager) DeliveredOrder(ctx context.Context, orderID string) error {
	if err := t.talabatClient.OrderPickedUp(ctx, talabatModels.OrderPickedUpRequest{
		OrderToken: orderID,
		Status:     models.OrderPickedUp.String(),
	}); err != nil {
		log.Trace().Err(err).Msgf("pickedUp talabat order, order_token=%v", orderID)
		return err
	}

	log.Info().Msgf("success pickedUp talabat order, order_token=%v", orderID)

	return nil
}

func (t TalabatManager) AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error {
	return nil
}

func (t TalabatManager) AcceptOrder(ctx context.Context, orderID string, pickUpTimePointer *time.Time) error {
	pickUpTime := *pickUpTimePointer
	if err := t.talabatClient.AcceptOrder(ctx, talabatModels.AcceptOrderRequest{
		OrderToken:     orderID,
		RemoteOrderId:  orderID,
		Status:         models.OrderAccepted.String(),
		AcceptanceTime: pickUpTime.Format(time.RFC3339),
	}); err != nil {
		log.Trace().Err(err).Msgf("accept talabat order, order_token=%v", orderID)
		return err
	}

	log.Info().Msgf("success accept talabat order, order_token=%v", orderID)

	return nil
}

func (t TalabatManager) RejectOrder(ctx context.Context, orderID, reason string) error {
	if err := t.talabatClient.RejectOrder(ctx, talabatModels.RejectOrderRequest{
		OrderToken: orderID,
		Reason:     orderID,
		Status:     models.OrderRejected.String(),
		Message:    reason,
	}); err != nil {
		log.Trace().Err(err).Msgf("reject talabat order, order_token=%v", orderID)
		return err
	}

	log.Info().Msgf("success reject talabat order, order_token=%v", orderID)

	return nil
}

func (t TalabatManager) MarkOrder(ctx context.Context, orderID string) error {
	if err := t.talabatClient.MarkOrderPrepared(ctx, orderID); err != nil {
		log.Trace().Err(err).Msgf("mark order prepared, order_token=%v", orderID)
		return err
	}

	log.Info().Msgf("success mark order prepared, order_token=%v", orderID)

	return nil
}

func (t TalabatManager) ConfirmPreOrder(ctx context.Context, orderID string) error {
	return errors.ErrUnsupportedMethod
}

func (t TalabatManager) UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error {
	return errors.ErrUnsupportedMethod
}
