package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/config"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models/utils"
)

type VirtualStore interface {
	GetVirtualStore(ctx context.Context, query selector.VirtualStore) (models.VirtualStore, error)
}

type VirtualStoreManager struct {
	ds               drivers.Datastore
	virtualStoreRepo drivers.VirtualRepository
	globalConfig     config.Configuration
}

func NewVirtualStoreManager(globalConfig config.Configuration, ds drivers.Datastore) VirtualStore {
	return &VirtualStoreManager{
		ds:               ds,
		virtualStoreRepo: ds.VirtualRepository(),
		globalConfig:     globalConfig,
	}
}

func (vs VirtualStoreManager) GetVirtualStore(ctx context.Context, query selector.VirtualStore) (models.VirtualStore, error) {
	virtualStore, err := vs.virtualStoreRepo.GetVirtualStore(ctx, query)
	if err != nil {
		utils.Beautify("can not find virtual store by selector", query)
		return models.VirtualStore{}, err
	}

	return virtualStore, nil
}
