package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/posist/clients/models"
)

type Posist interface {
	GetMenu(ctx context.Context, customerKey, tabId string) (models.Menu, error)
	GetStopList(ctx context.Context, customerKey string) ([]models.Item, error)
	CreateOrder(ctx context.Context, customerKey string, order models.Order) error
	GetTabs(ctx context.Context, customerKey string) ([]models.Tab, error)
	GetOrderStatus(ctx context.Context, orderId string) (models.OrderStatusResponse, error)
}

type Config struct {
	Protocol    string
	BaseURL     string
	AuthBasic   string
	CustomerKey string
}
