package clients

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/pkg/moysklad/models"
)

type Config struct {
	Protocol           string
	BaseURL            string
	Insecure           bool
	Username, Password string
}

type MoySklad interface {
	GetOrders(ctx context.Context) (models2.Order, error)
	GetMenu(ctx context.Context, req models2.GetMenuRequest) (models2.Menu, error)
	CreateSupplierOrder(ctx context.Context, req models2.SupplierOrder) (string, error)
	AddProductSupplier(ctx context.Context, position models2.Position) (string, error)
	DeleteProductSupplier(ctx context.Context, position models2.Position) error
}
