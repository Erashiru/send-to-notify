package clients

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/pkg/burgerking/clients/models"
)

type Config struct {
	Protocol string
	Address  string
}

type BK interface {
	SendOrder(ctx context.Context, order models2.Order) (models2.OrderResponse, error)
	CancelOrder(ctx context.Context, order models2.CancelOrderRequest) error
	SendMenu(ctx context.Context, req models2.SendMenuRequest) error
	SetSchedule(ctx context.Context, storeID string, req models2.SetScheduleRequest) error

	GetClosing(ctx context.Context, storeID string) (models2.ClosingResponse, error)
	UpdateClosing(ctx context.Context, storeID string) (models2.ClosingResponse, error)
	DeleteClosing(ctx context.Context, storeID string) error
	StopLists(ctx context.Context, req models2.StopListRequest) error
}
