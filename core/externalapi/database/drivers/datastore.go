package drivers

import (
	"context"

	"github.com/kwaaka-team/orders-core/core/externalapi/models"
)

type DataStore interface {
	Base

	AuthClientRepository() AuthClientRepository
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

type AuthClientRepository interface {
	FindByIDAndSecret(ctx context.Context, clientID string, clientSecret string) (models.AuthClient, error)
	FindByID(ctx context.Context, clientID string) (models.AuthClient, error)
	CreateAuthClient(ctx context.Context, req models.AuthClient) (string, error)
	AuthClientExist(ctx context.Context, clientID string, clientSecret string) error
	GetListID(ctx context.Context) ([]models.AuthClient, error)
	GetAuthClientByID(ctx context.Context, AuthId string) (models.AuthClient, error)
}
