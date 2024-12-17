package managers

import (
	"context"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/wolt/managers/validator"
	models2 "github.com/kwaaka-team/orders-core/core/wolt/models"
	orderCli "github.com/kwaaka-team/orders-core/pkg/order"
	storeCli "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/rs/zerolog/log"
)

type Order interface {
	CancelOrder(ctx context.Context, order models2.OrderNotification) (string, error)
}

type orderImplementation struct {
	orderCli       orderCli.Client
	storeCli       storeCli.Client
	orderValidator validator.Order
	baseURL        string
}

func NewOrder(orderCli orderCli.Client, storeCli storeCli.Client, baseURL string) Order {
	return &orderImplementation{
		orderCli:       orderCli,
		storeCli:       storeCli,
		orderValidator: validator.NewOrder(),
		baseURL:        baseURL,
	}
}

func (man *orderImplementation) CancelOrder(ctx context.Context, webhook models2.OrderNotification) (string, error) {
	if webhook.Body.Status != models2.Canceled.String() {
		log.Info().Msgf("successfully created order, id=%s", webhook.Body.Id)
		return "Successfully created order", nil
	}

	if err := man.orderCli.CancelOrderInPos(ctx, coreModels.CancelOrderInPos{
		OrderID:         webhook.Body.Id,
		DeliveryService: models2.WOLT.String(),
		CancelReason: coreModels.CancelReason{
			Reason: models2.USER_ERROR,
		},
		PaymentStrategy: models2.PAY_NOTHING,
	}); err != nil {
		log.Err(err).Msgf("cancel order error, order_id %s", webhook.Body.Id)
		return "", nil
	}
	log.Info().Msgf("successfully cancelled order, id=%s", webhook.Body.Id)
	return "Successfully canceled order", nil

}
