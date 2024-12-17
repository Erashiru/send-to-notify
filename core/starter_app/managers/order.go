package managers

import (
	"context"
	"fmt"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/starter_app/models"
	orderCli "github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/rs/zerolog/log"
)

type Order interface {
	CancelOrder(ctx context.Context, order models.Order) error
}

type orderImpl struct {
	orderCli orderCli.Client
}

func NewOrder(orderCli orderCli.Client) Order {
	return &orderImpl{
		orderCli: orderCli,
	}
}

func (o *orderImpl) CancelOrder(ctx context.Context, order models.Order) error {
	if order.Status != models.Canceled.String() {
		log.Info().Msgf("order status not canceled, order id: %s, status: %s", order.GlobalId, order.Status)
		return nil
	}

	if err := o.orderCli.CancelOrderInPos(ctx, coreModels.CancelOrderInPos{
		OrderID:         order.GlobalId,
		DeliveryService: models.STARTERAPP.String(),
		CancelReason: coreModels.CancelReason{
			Reason: fmt.Sprintf("cancel by delivery service, status: %s", order.Status),
		},
	}); err != nil {
		log.Err(err).Msgf("cancel order error, order id: %s", order.GlobalId)
		return nil
	}

	log.Info().Msgf("successfully cancelled order, order id: %s", order.GlobalId)
	return nil
}
