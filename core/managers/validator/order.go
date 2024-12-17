package validator

import (
	"context"
	"github.com/pkg/errors"

	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/custom"
)

type Order interface {
	ValidateOrder(ctx context.Context, order models.Order) error
}

type orderImpl struct{}

var _ Order = orderImpl{}

func NewOrderValidator() Order {
	return &orderImpl{}
}

func (o orderImpl) ValidateOrder(ctx context.Context, order models.Order) error {
	var errs custom.Error

	if order.Type != "INSTANT" && order.Type != "PREORDER" {
		errs.Append(errors.New("type should be INSTANT or PREORDER"))
	}

	if models.Aggregator(order.DeliveryService).String() == "" {
		errs.Append(errors.New("delivery service is unknown"))
	}

	if order.StoreID == "" {
		errs.Append(errors.New("delivery store id is empty"))
	}

	if order.OrderTime.Value.IsZero() {
		errs.Append(errors.New("order time is null"))
	}

	if order.EstimatedPickupTime.Value.IsZero() {
		errs.Append(errors.New("estimated pick up time is null"))
	}

	if order.PaymentMethod != models.PAYMENT_METHOD_CASH && order.PaymentMethod != models.PAYMENT_METHOD_DELAYED {
		errs.Append(errors.New("payment method should be CASH or DELAYED"))
	}

	if order.EstimatedTotalPrice.Value == 0 {
		errs.Append(errors.New("estimated total price should be great 0"))
	}

	if order.TotalCustomerToPay.Value == 0 {
		errs.Append(errors.New("total customer to pay price should be great 0"))
	}

	if len(order.Products) == 0 {
		errs.Append(errors.New("products length is equal 0"))
	}

	return errs.ErrorOrNil()
}
