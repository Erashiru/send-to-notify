package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/errors"

	"github.com/kwaaka-team/orders-core/core/models/selector"
	deliverooConf "github.com/kwaaka-team/orders-core/pkg/deliveroo"
	deliverooCli "github.com/kwaaka-team/orders-core/pkg/deliveroo/clients"
	deliverooModels "github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/dto"
	deliverooHttpCli "github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/http"
	"github.com/rs/zerolog/log"
	"time"
)

type Deliveroo struct {
	cli deliverooHttpCli.Client
}

func (d Deliveroo) AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error {
	return errors.ErrUnsupportedMethod
}

func (d Deliveroo) AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error {
	return errors.ErrUnsupportedMethod
}

func (d Deliveroo) RejectOrder(ctx context.Context, orderID, reason string) error {
	err := d.cli.UpdateOrderStatus(ctx, deliverooModels.UpdateOrderStatusRequest{
		Status:       "rejected",
		RejectReason: reason,
	}, orderID)

	log.Trace().Err(err).Msgf("Reject order response order_id=%v", orderID)

	if err != nil {
		log.Trace().Err(err).Msgf("reject deliveroo order error, order_id=%v", orderID)
		return err
	}

	log.Info().Msgf("success rejected deliveroo order, order_id=%v", orderID)

	return nil
}

func (d Deliveroo) MarkOrder(ctx context.Context, orderID string) error {
	return errors.ErrUnsupportedMethod

}

func (d Deliveroo) ConfirmPreOrder(ctx context.Context, orderID string) error {
	return errors.ErrUnsupportedMethod

}

func (d Deliveroo) UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error {
	err := d.cli.UpdateOrderStatus(ctx, deliverooModels.UpdateOrderStatusRequest{
		Status: req.OrderStatus,
	}, req.OrderID)

	log.Trace().Err(err).Msgf("Update order response order_id=%v", req.OrderID)

	if err != nil {
		log.Trace().Err(err).Msgf("update deliveroo order, order_id=%v, status=%v", req.OrderID, req.OrderStatus)
		return err
	}

	log.Info().Msgf("success update deliveroo order, order_id=%v, status=%v", req.OrderID, req.OrderStatus)

	return nil
}

func (d Deliveroo) DeliveredOrder(ctx context.Context, orderID string) error {
	return errors.ErrUnsupportedMethod

}

func NewManager(username, password, baseUrl string) (Deliveroo, error) {

	cli, err := deliverooConf.NewDeliverooClient(&deliverooCli.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		Username: username,
		Password: password,
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize Deliveroo client ")
		return Deliveroo{}, err
	}
	return Deliveroo{
		cli: *cli,
	}, nil
}
