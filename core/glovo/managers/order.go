package managers

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/glovo/models"
	"github.com/rs/zerolog/log"

	coreModels "github.com/kwaaka-team/orders-core/core/models"
	orderCli "github.com/kwaaka-team/orders-core/pkg/order"
	storeCli "github.com/kwaaka-team/orders-core/pkg/store"
)

type Order interface {
	CancelOrder(ctx context.Context, req models2.CancelOrderRequest) error
}

type orderImplementation struct {
	orderCli orderCli.Client
	storeCli storeCli.Client
}

func NewOrder(orderCli orderCli.Client, storeCli storeCli.Client) Order {
	return &orderImplementation{
		orderCli: orderCli,
		storeCli: storeCli,
	}
}

func (man *orderImplementation) CancelOrder(ctx context.Context, req models2.CancelOrderRequest) error {
	err := man.orderCli.CancelOrderInPos(ctx, coreModels.CancelOrderInPos{
		OrderID:         req.OrderID,
		PaymentStrategy: req.PaymentStrategy,
		CancelReason: coreModels.CancelReason{
			Reason: req.CancelReason,
		},
	})

	if err != nil {
		log.Err(err).Msgf("glovo cancel order id %s", req.OrderID)
		return err
	}

	return nil
}
