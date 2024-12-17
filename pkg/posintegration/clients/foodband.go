package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/domain/foodband"
)

type Config struct {
	CreateOrderUrl string
	CancelOrderUrl string
	ApiToken       string
	RetryMaxCount  int
}

type FOODBAND interface {
	CreateOrder(ctx context.Context, createOrderBody foodband.CreateOrderRequest) (int, error)
	CancelOrder(ctx context.Context, cancelOrderBody foodband.CancelOrderRequest) error
}
