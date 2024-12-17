package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/jowi/models"
	"github.com/kwaaka-team/orders-core/core/jowi/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/order"
	"github.com/kwaaka-team/orders-core/pkg/store"
)

type JowiManager interface {
	UpdateOrder(ctx context.Context, updateOrder models.Event) error
}

type Jowi struct {
	storeCli store.Client
	orderCli order.Client
}

func NewJowiManager(storeCli store.Client, orderCli order.Client) JowiManager {
	return &Jowi{
		storeCli: storeCli,
		orderCli: orderCli,
	}
}

func (jowi *Jowi) UpdateOrder(ctx context.Context, updateOrder models.Event) error {
	req, err := updateOrder.ToOrderModel()
	if err != nil {
		return err
	}

	utils.Beautify("order request body to orders-core from jowi-handler", req)

	if err := jowi.orderCli.UpdateOrder(ctx, req); err != nil {
		return err
	}

	return nil
}
