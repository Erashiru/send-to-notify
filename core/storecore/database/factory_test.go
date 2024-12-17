package database

import (
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"testing"
)

func TestFactoryNew_SingleConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panic")
		}
	}()

	url := "url"
	dsName := mongoDatastore
	dbName := "dbName"

	cfg1 := drivers.DataStoreConfig{
		URL:           url,
		DataStoreName: dsName,
		DataBaseName:  dbName,
	}
	ds1, err := New(cfg1)
	if err != nil {
		t.Fatal(err)
	}

	cfg2 := drivers.DataStoreConfig{
		URL:           url,
		DataStoreName: dsName,
		DataBaseName:  dbName,
	}
	ds2, err := New(cfg2)
	if err != nil {
		t.Fatal(err)
	}

	cfg3 := drivers.DataStoreConfig{
		URL:           url,
		DataStoreName: dsName,
		DataBaseName:  dbName,
	}
	ds3, err := New(cfg3)
	if err != nil {
		t.Fatal(err)
	}

	if ds1 != ds2 {
		t.Fatal("error during creation datastore in factory")
	}

	if ds2 != ds3 {
		t.Fatal("error during creation datastore in factory")
	}

	if ds1 != ds3 {
		t.Fatal("error during creation datastore in factory")
	}

}

func TestFactoryNew_CacheDifferentConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panic")
		}
	}()

	url := "url"
	dsName := mongoDatastore

	cfg1 := drivers.DataStoreConfig{
		URL:           url,
		DataStoreName: dsName,
		DataBaseName:  "dbName1",
	}
	ds1, err := New(cfg1)
	if err != nil {
		t.Fatal(err)
	}

	cfg2 := drivers.DataStoreConfig{
		URL:           url,
		DataStoreName: dsName,
		DataBaseName:  "dbName2",
	}
	ds2, err := New(cfg2)
	if err != nil {
		t.Fatal(err)
	}

	if ds1 == ds2 {
		t.Fatal("error during creation datastore in factory")
	}

}
