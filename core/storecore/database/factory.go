package database

import (
	"fmt"
	drivers2 "github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers/mongodb"
)

const (
	mongoDatastore = "mongo"
)

var databasesCache = make(map[string]map[drivers2.DataStoreConfig]drivers2.Datastore)

func New(conf drivers2.DataStoreConfig) (drivers2.Datastore, error) {
	if db := getFromCache(conf); db != nil {
		return db, nil
	}
	db, err := create(conf)

	if err == nil {
		cache(conf, db)
	}

	return db, err
}

func cache(conf drivers2.DataStoreConfig, db drivers2.Datastore) {
	if _, ok := databasesCache[conf.DataStoreName]; !ok {
		databasesCache[conf.DataStoreName] = make(map[drivers2.DataStoreConfig]drivers2.Datastore)
	}
	if _, ok := databasesCache[conf.DataStoreName][conf]; !ok {
		databasesCache[conf.DataStoreName][conf] = db
	}
}

func getFromCache(conf drivers2.DataStoreConfig) drivers2.Datastore {
	var mapConfigVsDS map[drivers2.DataStoreConfig]drivers2.Datastore
	var ok bool
	if mapConfigVsDS, ok = databasesCache[conf.DataStoreName]; !ok {
		return nil
	}
	var cachedDataSource drivers2.Datastore
	if cachedDataSource, ok = mapConfigVsDS[conf]; !ok {
		return nil
	}
	return cachedDataSource
}

func create(conf drivers2.DataStoreConfig) (drivers2.Datastore, error) {
	if conf.DataStoreName == mongoDatastore {
		return mongodb.New(conf)
	}
	return nil, fmt.Errorf("invalid datastore name: %s", conf.DataStoreName)
}
