package database

import (
	"fmt"

	"github.com/kwaaka-team/orders-core/core/externalapi/database/drivers"
	"github.com/kwaaka-team/orders-core/core/externalapi/database/drivers/mongo"
)

const (
	mongoDatastore = "mongo"
)

func New(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	if conf.DataStoreName == mongoDatastore {
		return mongo.New(conf)
	}
	return nil, fmt.Errorf("invalid datastore name: %s", conf.DataStoreName)
}
