package validator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/glovo/models"
)

type Order interface {
	ValidateCreateOrder(ctx context.Context, req models.Order) (models.Order, error)
}

type orderImpl struct{}

var _ Order = (*orderImpl)(nil)

func NewOrder() Order {
	return &orderImpl{}
}

func (o orderImpl) ValidateCreateOrder(ctx context.Context, req models.Order) (models.Order, error) {
	return models.Order{}, nil
}
