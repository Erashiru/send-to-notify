package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/ytimes/clients/models"
)

type Config struct {
	BaseUrl string
	Token   string
}

type Client interface {
	GetPoints(ctx context.Context) (models.PointInfo, error)
	CreateOrder(ctx context.Context, req models.Order) (models.CreateOrderResponse, error)
	GetMenu(ctx context.Context, pointGuid string) (models.Menu, error)
	GetSupplementList(ctx context.Context, pointGuid string) (models.SupplementList, error)
}
