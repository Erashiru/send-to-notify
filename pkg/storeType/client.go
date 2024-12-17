package store_type

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/storecore/config"
	"github.com/kwaaka-team/orders-core/core/storecore/database"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/storeType/dto"
)

type Client interface {
	GetStoreType(ctx context.Context) ([]models.StoreType, error)
}

type StoreType struct {
	storeTypeMan managers.StoreType
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

	return &StoreType{
		storeTypeMan: managers.NewStoreTypeManager(ds),
	}, nil
}

func (s *StoreType) GetStoreType(ctx context.Context) ([]models.StoreType, error) {
	res, err := s.storeTypeMan.GetList(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}
