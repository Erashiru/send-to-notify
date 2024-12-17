package clients

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/pkg/poster/clients/models"
)

type Config struct {
	Protocol string
	BaseURL  string
	Insecure bool

	Token string
}

type Poster interface {
	GetProducts(ctx context.Context) (models2.GetProductsResponse, error)
	GetSpots(ctx context.Context) (models2.GetSpotsResponse, error)
	CreateOrder(ctx context.Context, req models2.CreateOrderRequest) (models2.CreateOrderResponse, error)
	GetOrder(ctx context.Context, id string) (models2.CreateOrderResponse, error)
	GetOrders(ctx context.Context, req models2.GetOrdersRequest) (models2.GetOrdersResponse, error)
	GetStopList(ctx context.Context) (models2.GetStopListResponse, error)
	GetIngredients(ctx context.Context) (models2.GetIngridientsResponse, error)
	GetProduct(ctx context.Context, productId string) (models2.GetProductsResponseBody, error)
}
