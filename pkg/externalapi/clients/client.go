package clients

import (
	"context"
	dto2 "github.com/kwaaka-team/orders-core/pkg/externalapi/clients/dto"
)

type Client interface {
	UpdateOrderWebhook(ctx context.Context, order dto2.Order, path string) error
	UpdateProductStopList(ctx context.Context, product dto2.Product, path string) error
	UpdateModifierStopList(ctx context.Context, modifier dto2.Modifier, path string) error
}

type Config struct {
	Protocol  string
	Insecure  bool
	AuthToken string
}
