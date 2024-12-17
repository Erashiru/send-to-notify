package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/yaros/models"
)

type Config struct {
	Protocol string
	BaseURL  string
	Username string
	Password string
	RestID   string
}

type Yaros interface {
	GetItems(ctx context.Context, restID string) (models.GetItemsResponse, error)
	GetCategories(ctx context.Context, restID string) (models.GetCategoriesResponse, error)
	CreateOrder(ctx context.Context, restID string, order models.OrderRequest) (models.OrderResponse, error)
	UpdateOrder(ctx context.Context, restID string, update models.OrderRequest) (models.OrderResponse, error)
	GetStopList(ctx context.Context, restID string) (models.StopListResponse, error)
	GetOrders(ctx context.Context, restID, infoSystem, department string) (models.OrderResponse, error)
}
