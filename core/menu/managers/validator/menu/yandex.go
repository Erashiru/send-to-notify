package menu

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
)

type yandex struct{}

func newYandex() *yandex {
	return &yandex{}
}

func (g *yandex) Validate(ctx context.Context, menu models.Menu) error {

	return nil
}
