package managers

import (
	"context"
	drivers2 "github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type TapRestaurant interface {
	CreateTapRestaurant(ctx context.Context, req models.TapRestaurant) (string, error)
	GetTapRestaurantList(ctx context.Context, query selector.TapRestaurant) ([]models.TapRestaurant, int, error)
	GetTapRestaurant(ctx context.Context, id string) (models.TapRestaurant, error)
	GetTapRestaurantByName(ctx context.Context, query selector.TapRestaurant) (models.TapRestaurant, error)
	UpdateTapRestaurant(ctx context.Context, req models.UpdateTapRestaurant) error
	DeleteTapRestaurant(ctx context.Context, id string) error
}

type TapRestaurantManager struct {
	tapRepository drivers2.TapRestaurantRepository
}

func NewTapRestaurantManager(ds drivers2.Datastore) TapRestaurant {
	return &TapRestaurantManager{
		tapRepository: ds.TapRestaurantRepository(),
	}
}

func (tm *TapRestaurantManager) CreateTapRestaurant(ctx context.Context, req models.TapRestaurant) (string, error) {
	return tm.tapRepository.Create(ctx, req)
}

func (tm *TapRestaurantManager) GetTapRestaurantList(ctx context.Context, query selector.TapRestaurant) ([]models.TapRestaurant, int, error) {
	return tm.tapRepository.GetList(ctx, query)
}

func (tm *TapRestaurantManager) GetTapRestaurant(ctx context.Context, id string) (models.TapRestaurant, error) {
	return tm.tapRepository.GetByID(ctx, id)
}

func (tm *TapRestaurantManager) GetTapRestaurantByName(ctx context.Context, query selector.TapRestaurant) (models.TapRestaurant, error) {
	return tm.tapRepository.GetByQuery(ctx, query)
}

func (tm *TapRestaurantManager) UpdateTapRestaurant(ctx context.Context, req models.UpdateTapRestaurant) error {
	return tm.tapRepository.Update(ctx, req)
}

func (tm *TapRestaurantManager) DeleteTapRestaurant(ctx context.Context, id string) error {
	return tm.tapRepository.Delete(ctx, id)
}
