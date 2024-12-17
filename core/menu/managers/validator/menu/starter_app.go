package menu

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
)

type starterApp struct{}

func newstarterApp() *starterApp {
	return &starterApp{}
}

func (w *starterApp) Validate(ctx context.Context, menu models.Menu) error {
	return nil
}
