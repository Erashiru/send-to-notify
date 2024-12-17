package clients

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
)

type Paloma interface {
	CreateOrder(ctx context.Context, pointID string, req models2.Order) (models2.OrderResponse, error)
	GetMenu(ctx context.Context, pointID string) (models2.Menu, error)
	GetStopList(ctx context.Context, pointID string) (models2.StopList, error)
	GetPoints(ctx context.Context, authKey string) ([]models2.Point, error)
	GetOrderStatus(ctx context.Context, orderID string) (models2.OrderResponse, error)
	CancelOrder(ctx context.Context, orderID string) (models2.OrderResponse, error)
}

type Config struct {
	Protocol string
	BaseURL  string
	ApiKey   string
	Class    string
}
