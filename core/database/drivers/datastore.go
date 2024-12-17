package drivers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/kwaaka-team/orders-core/core/models/selector"

	"github.com/kwaaka-team/orders-core/core/models"
)

type DataStore interface {
	Base

	OrderRepository() OrderRepository
	BKOfferRepository() BKOfferRepository
	AnalyticsRepository() AnalyticsRepository
}

// Base представляет базовый интерфейс для работы с DataStore.
type Base interface {
	// Name - возвращает название DataStore.
	Name() string

	// Ping - проверка на работоспособность.
	Ping() error

	// Close - закрывает соединение с DataStore.
	Close(ctx context.Context) error

	// Connect - устанавливает соединение с DataStore.
	Connect(cli *mongo.Client) error

	Client() *mongo.Client
}

type OrderRepository interface {
	GetOrder(ctx context.Context, query selector.Order) (models.Order, error)
	GetOrders(ctx context.Context, query selector.Order) ([]models.Order, int, error)
	InsertOrder(ctx context.Context, req models.Order) (models.Order, error)
	UpdateOrder(ctx context.Context, req models.Order) error
	CancelOrder(ctx context.Context, req models.CancelOrder) (models.Order, error)
	GetOrderStatus(ctx context.Context, id string) (models.Order, error)
	UpdateOrderStatus(ctx context.Context, query selector.Order, status, errorDescription string) error
	UpdateOrderStatusByID(ctx context.Context, id, posName, status string) error
	SetPaidStatus(ctx context.Context, orderID string) error
	GetAllOrders(ctx context.Context, query selector.Order) ([]models.Order, error)
}

type BKOfferRepository interface {
	GetActiveOffers(ctx context.Context) ([]models.BKOffer, error)
}

type AnalyticsRepository interface {
}
