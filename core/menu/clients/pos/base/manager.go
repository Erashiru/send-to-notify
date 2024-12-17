package base

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
)

type Manager interface {
	GetMenu(ctx context.Context, store storeModels.Store) (models.Menu, error)
	GetAggMenu(ctx context.Context, store storeModels.Store) ([]models.Menu, error)
}
