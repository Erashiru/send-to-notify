package managers

import (
	"context"
	drivers2 "github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	models2 "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
)

type StoreGroup interface {
	FindAll(ctx context.Context) ([]models2.StoreGroup, error)
	FindStoreGroup(ctx context.Context, query selector.StoreGroup) (models2.StoreGroup, error)
	FindStoreGroupById(ctx context.Context, id string) (models2.StoreGroup, error)
	UpdateStores(ctx context.Context, storeGroup models2.UpdateStoreGroup) (int64, error)
	CreateStoreGroup(ctx context.Context, storeGroup models2.StoreGroup) (string, error)
}

type StoreGroupManager struct {
	storeGroupRepository drivers2.StoreGroupRepository
}

func NewStoreGroupManager(ds drivers2.Datastore) StoreGroup {
	return &StoreGroupManager{
		storeGroupRepository: ds.StoreGroupRepository(),
	}
}

func (s *StoreGroupManager) FindStoreGroupById(ctx context.Context, id string) (models2.StoreGroup, error) {
	storeGroup, err := s.storeGroupRepository.Get(ctx, selector.NewEmptyStoreGroupSearch().SetID(id))
	if err != nil {
		return models2.StoreGroup{}, err
	}

	return storeGroup, nil
}

func (s *StoreGroupManager) CreateStoreGroup(ctx context.Context, storeGroup models2.StoreGroup) (string, error) {

	storeGroupDB, err := s.storeGroupRepository.Get(ctx, selector.StoreGroup{
		Name: storeGroup.Name,
	})
	if storeGroupDB.Name != "" {
		return "", errors.New("store group already exists name: " + storeGroupDB.Name)
	}

	if err != nil && !errors.Is(err, drivers2.ErrNotFound) {
		return "", err
	}

	storeId, err := s.storeGroupRepository.Create(ctx, storeGroup)

	return storeId, err
}

func (s *StoreGroupManager) UpdateStores(ctx context.Context, storeGroup models2.UpdateStoreGroup) (int64, error) {
	count, err := s.storeGroupRepository.UpdateByFields(ctx, storeGroup)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *StoreGroupManager) FindStoreGroup(ctx context.Context, query selector.StoreGroup) (models2.StoreGroup, error) {
	storeGroup, err := s.storeGroupRepository.Get(ctx, query)
	if err != nil {
		return models2.StoreGroup{}, err
	}

	return storeGroup, nil
}

func (s *StoreGroupManager) FindAll(ctx context.Context) ([]models2.StoreGroup, error) {

	storeGroups, err := s.storeGroupRepository.All(ctx)
	if err != nil {
		return nil, err
	}

	return storeGroups, nil
}
