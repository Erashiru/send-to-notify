package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/yandex/models"
)

type Config struct {
	Protocol     string
	BaseURL      string
	ClientID     string
	ClientSecret string
}

type Yandex interface {
	MenuImportInitiation(ctx context.Context, req models.MenuInitiationRequest) error
	Close(ctx context.Context)
}
