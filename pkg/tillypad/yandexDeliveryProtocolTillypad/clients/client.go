package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
)

type Config struct {
	BaseURL      string
	Protocol     string
	ClientId     string
	ClientSecret string
	PathPrefix   string
}

type Client interface {
	GetOrder(ctx context.Context, orderId string) (models.Order, error)
	GetOrderStatus(ctx context.Context, orderId string) (models.OrderStatusResponse, error)
	GetMenu(ctx context.Context, storeId string) (models.Menu, error)
	GetStores(ctx context.Context) (models.GetStoreResponse, error)
	GetPromos(ctx context.Context, storeId string) (models.Promo, error)
	GetAccessToken(ctx context.Context, clientId, clientSecret string) (string, error)
	GetAvailability(ctx context.Context, storeId string) (models.StopListResponse, error)
	CreateOrder(ctx context.Context, order models.Order) (models.CreationResult, error)
}
