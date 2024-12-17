package menu

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
)

type express24 struct{}

func newExpress24() *express24 {
	return &express24{}
}

func (w *express24) Validate(ctx context.Context, menu models.Menu) error {
	return nil
}
