package clients

import (
	"context"
	dto2 "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/dto"
)

type Config struct {
	Protocol string
	BaseURL  string
	ApiKey   string
}

type RKeeper interface {
	GetMenu(ctx context.Context, objectId int) (dto2.MenuResponse, error)
	GetStopList(ctx context.Context, objectId int) (dto2.StopListResponse, error)

	CreateOrderTask(ctx context.Context, taskGUID string) (dto2.CreateOrderTaskResponse, error)
	GetOrderTask(ctx context.Context, taskGUID string) (dto2.GetOrderTaskResponse, error)
	GetOrder(ctx context.Context, orderGUID string, objectID int) (dto2.SyncResponse, error)
	CreateOrder(ctx context.Context, objectID int, order dto2.Order) (dto2.SyncResponse, error)

	CancelOrder(ctx context.Context, objectID int, orderGUID string) (dto2.SyncResponse, error)
	CancelOrderTask(ctx context.Context, taskGUID string) (dto2.CancelOrderResponse, error)
	UpdateMenu(ctx context.Context, objectID int) (dto2.SyncResponse, error)
	PayOrder(ctx context.Context, objectID, amount int, orderId, currency string) (dto2.SyncResponse, error)
	GetMenuByParams(ctx context.Context, objectID, priceTypeID int) (dto2.MenuResponse, error)
}
