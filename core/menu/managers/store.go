package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/managers/validator"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
)

type StoreManager interface {
	GetStore(ctx context.Context, store selector.Store) (storeModels.Store, error)
	GetStores(ctx context.Context, store selector.Store) ([]storeModels.Store, error)
	ListStoresByProduct(ctx context.Context, query selector.Menu) ([]storeModels.Store, int64, error)
}

type storeImpl struct {
	globalConfig   menu.Configuration
	storeRepo      drivers.StoreRepository
	menuRepo       drivers.MenuRepository
	storeValidator validator.Store
}

func NewStoreManager(globalConfig menu.Configuration, storeRepo drivers.StoreRepository, menuRepo drivers.MenuRepository, validator validator.Store) StoreManager {
	return &storeImpl{
		globalConfig:   globalConfig,
		storeRepo:      storeRepo,
		menuRepo:       menuRepo,
		storeValidator: validator,
	}
}

func (sm *storeImpl) GetStores(ctx context.Context, query selector.Store) ([]storeModels.Store, error) {
	if err := sm.storeValidator.ValidateDelivery(ctx, query); err != nil {
		return nil, err
	}

	stores, _, err := sm.storeRepo.List(ctx, query)
	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (sm *storeImpl) GetStore(ctx context.Context, query selector.Store) (storeModels.Store, error) {
	if err := sm.storeValidator.ValidateExternalID(ctx, query); err != nil {
		return storeModels.Store{}, err
	}

	store, err := sm.storeRepo.Get(ctx, query)
	if err != nil {
		return storeModels.Store{}, err
	}

	return store, nil
}

func (sm *storeImpl) ListStoresByProduct(ctx context.Context, query selector.Menu) ([]storeModels.Store, int64, error) {

	if err := sm.storeValidator.ValidateListStores(ctx, query); err != nil {
		return nil, 0, err
	}

	menus, err := sm.menuRepo.GetMenuIDs(ctx, query)
	if err != nil {
		if errors.Is(err, drivers.ErrNotFound) {
			return []storeModels.Store{}, 0, nil
		}
		return nil, 0, err
	}

	stores, total, err := sm.storeRepo.List(ctx, selector.EmptyStoreSearch().
		SetAggregatorMenuIDs(menus))
	if err != nil {
		return nil, 0, err
	}
	isProductAvailable := *query.IsProductAvailable
	stores = SetProductsStopInMenu(menus, stores, isProductAvailable)

	return stores, total, nil
}

// To set product is stop or not by looping store aaray and matching menus where that produsts is
func SetProductsStopInMenu(menus []string, stores []storeModels.Store, state bool) []storeModels.Store {
	m := make(map[string]bool)
	for _, v := range menus {
		m[v] = true
	}

	for i := range stores {
		for j, val := range stores[i].Menus {
			if m[val.ID] {
				stores[i].Menus[j].IsProductOnStop = !state
				continue
			}
			stores[i].Menus[j].IsProductOnStop = state
		}
	}

	return stores
}
