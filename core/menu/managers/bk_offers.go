package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
)

type BkOffersManager interface {
	List(ctx context.Context, query selector.BkOffers) ([]models.BkOffers, error)
}

type bkOffersMan struct {
	globalConfig menu.Configuration
	bkOffersRepo drivers.BkOffersRepository
}

func NewBkOffersManager(
	globalConfig menu.Configuration,
	bkOffersRepo drivers.BkOffersRepository) BkOffersManager {
	return &bkOffersMan{
		globalConfig: globalConfig,
		bkOffersRepo: bkOffersRepo,
	}
}

func (m *bkOffersMan) List(ctx context.Context, query selector.BkOffers) ([]models.BkOffers, error) {
	return m.bkOffersRepo.List(ctx, query)
}
