package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type StoreType interface {
	GetList(ctx context.Context) ([]models.StoreType, error)
}

type StoreTypeManager struct {
	storeTypeRps drivers.StoreTypeRepository
}

func NewStoreTypeManager(ds drivers.Datastore) StoreType {
	return &StoreTypeManager{
		storeTypeRps: ds.StoreTypeRepository(),
	}
}

func (s *StoreTypeManager) GetList(ctx context.Context) ([]models.StoreType, error) {
	return s.storeTypeRps.GetList(ctx)
}
