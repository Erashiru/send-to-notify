package storeGroup

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/storecore/config"
	"github.com/kwaaka-team/orders-core/core/storecore/database"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/pkg/storeGroup/dto"
)

type Client interface {
	FindStoreGroup(ctx context.Context, query selector.StoreGroup) (dto.StoreGroup, error)
	FindStoreGroupById(ctx context.Context, id string) (dto.StoreGroup, error)
	FindAll(ctx context.Context) ([]dto.StoreGroup, error)
	CreateStoreGroup(ctx context.Context, query dto.StoreGroup) (string, error)
}

type StoreGroup struct {
	storeGroupManager managers.StoreGroup
}

func NewClient(cfg dto.Config) (Client, error) {
	opts, err := config.LoadConfig(context.Background())
	if err != nil {
		return nil, err
	}

	ds, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSDB,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create datastore %s: %v", opts.DSName, err)
	}

	if err = ds.Connect(cfg.MongoCli); err != nil {
		return nil, fmt.Errorf("cannot connect to datastore: %s", err)
	}

	return &StoreGroup{
		storeGroupManager: managers.NewStoreGroupManager(ds),
	}, nil
}

func (s *StoreGroup) FindStoreGroupById(ctx context.Context, id string) (dto.StoreGroup, error) {
	storeGroup, err := s.storeGroupManager.FindStoreGroupById(ctx, id)
	if err != nil {
		return dto.StoreGroup{}, err
	}

	return dto.FromStoreGroupModel(storeGroup), nil
}

func (s *StoreGroup) CreateStoreGroup(ctx context.Context, query dto.StoreGroup) (string, error) {
	storeGroup := dto.ToStoreGroupModel(query)

	storeGroupID, err := s.storeGroupManager.CreateStoreGroup(ctx, storeGroup)
	if err != nil {
		return "", err
	}
	return storeGroupID, nil
}

func (s *StoreGroup) FindStoreGroup(ctx context.Context, group selector.StoreGroup) (dto.StoreGroup, error) {
	storeGroup, err := s.storeGroupManager.FindStoreGroup(ctx, group)
	if err != nil {
		return dto.StoreGroup{}, err
	}

	return dto.FromStoreGroupModel(storeGroup), nil
}

func (s *StoreGroup) FindAll(ctx context.Context) ([]dto.StoreGroup, error) {
	storeGroups, err := s.storeGroupManager.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responseStoreGroup := make([]dto.StoreGroup, 0, len(storeGroups))

	for _, storeGroup := range storeGroups {
		responseStoreGroup = append(responseStoreGroup, dto.FromStoreGroupModel(storeGroup))
	}
	return responseStoreGroup, nil
}
