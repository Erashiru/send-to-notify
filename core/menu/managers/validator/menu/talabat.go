package menu

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
)

type talabat struct{}

func newTalabat() *talabat {
	return &talabat{}
}

func (g *talabat) Validate(ctx context.Context, menu models.Menu) error {

	return nil
}