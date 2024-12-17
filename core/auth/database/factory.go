package database

import (
	drivers2 "github.com/kwaaka-team/orders-core/core/auth/database/datastore/drivers"
	"github.com/kwaaka-team/orders-core/core/auth/database/datastore/drivers/mongodb"
	"github.com/pkg/errors"
)

func New(cfg drivers2.DataStoreConfig) (drivers2.DataStore, error) {
	switch cfg.DataStoreName {
	case "mongo":
		return mongodb.New(cfg)
	}

	return nil, errors.New("datastore create error")
}
