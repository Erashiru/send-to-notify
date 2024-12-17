package drivers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/auth/models"
	"github.com/kwaaka-team/orders-core/core/auth/models/selector"
)

type DataStore interface {
	Base

	AuthRepository() UserRepository
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
	Connect() error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUser(ctx context.Context, query selector.User) (models.User, error)
	UpdateUserInfo(ctx context.Context, user models.User) error
}
