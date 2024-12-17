package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/errors"
	"strconv"
	"time"

	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	glovoClient "github.com/kwaaka-team/orders-core/pkg/glovo"
	glovoConfig "github.com/kwaaka-team/orders-core/pkg/glovo/clients"
	glovoModels "github.com/kwaaka-team/orders-core/pkg/glovo/clients/dto"
	"github.com/rs/zerolog/log"
)

type Glovo interface {
	UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error
	AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error
	AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error
	RejectOrder(ctx context.Context, orderID, reason string) error
	MarkOrder(ctx context.Context, orderID string) error
	ConfirmPreOrder(ctx context.Context, orderID string) error
	DeliveredOrder(ctx context.Context, orderID string) error
}

type GlovoManager struct {
	glovoClient glovoConfig.Glovo
}

func NewGlovoManager(globalConfig config.Configuration) (Glovo, error) {
	glovoClient, err := glovoClient.NewGlovoClient(&glovoConfig.Config{
		Protocol: "http",
		BaseURL:  globalConfig.GlovoConfiguration.BaseURL,
		ApiKey:   globalConfig.GlovoConfiguration.Token,
	})
	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Glovo Manager.")
		return nil, err
	}
	return &GlovoManager{
		glovoClient: glovoClient,
	}, nil
}

func (g GlovoManager) UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error {
	orderID, err := strconv.ParseInt(req.OrderID, 10, 64)
	if err != nil {
		return err
	}
	res, err := g.glovoClient.UpdateOrderStatus(ctx, glovoModels.OrderUpdateRequest{
		ID:      orderID,
		Status:  req.OrderStatus,
		StoreID: req.StoreID,
	})

	utils.Beautify("Update order response", res)

	if err != nil {
		log.Trace().Err(err).Msgf("update glovo order, order_id=%v, status=%v", req.OrderID, req.OrderStatus)
		return err
	}

	log.Info().Msgf("success update glovo order, order_id=%v, status=%v", req.OrderID, req.OrderStatus)

	return nil
}

func (g GlovoManager) AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error {
	return errors.ErrUnsupportedMethod
}

func (g GlovoManager) AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error {
	return errors.ErrUnsupportedMethod
}
func (g GlovoManager) RejectOrder(ctx context.Context, orderID, reason string) error {
	return errors.ErrUnsupportedMethod
}
func (g GlovoManager) MarkOrder(ctx context.Context, orderID string) error {
	return errors.ErrUnsupportedMethod
}
func (g GlovoManager) ConfirmPreOrder(ctx context.Context, orderID string) error {
	return errors.ErrUnsupportedMethod
}

func (g GlovoManager) DeliveredOrder(ctx context.Context, orderID string) error {
	return errors.ErrUnsupportedMethod
}
