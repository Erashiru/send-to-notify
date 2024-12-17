package main

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	os.Setenv("SECRET_ENV", "ProdEnvs")
	os.Setenv("REGION", "eu-west-1")
	os.Setenv("SENTRY", "ProdSentry")

	ctx := context.Background()

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		return err
	}

	db, err := cmd.CreateMongo(ctx, cfg.DSURL, cfg.DSDB)
	if err != nil {
		return err
	}

	menuService, err := createMenuService(db)
	if err != nil {
		return err
	}
	_ = menuService

	storeService, err := createStoreService(db)
	if err != nil {
		return err
	}

	stores, err := storeService.FindAllStores(ctx)
	if err != nil {
		return err
	}

	workerPool, err := NewWorkerPool(50, menuService)
	if err != nil {
		return err
	}

	before := time.Now()

	storesDisabledMap := make(map[string][]coreStoreModels.Store)
	workerPool.Start(ctx)
	workerPool.Submit(stores)
	log.Printf("finding stores by pos, stores count = %d, time took %v", len(stores), time.Since(before))

	resultStores := workerPool.GetResult()
	for i := range resultStores {
		st := resultStores[i]
		if storesOfGroup, ok := storesDisabledMap[st.RestaurantGroupID]; ok {
			storesOfGroup = append(storesOfGroup, st)
			storesDisabledMap[st.RestaurantGroupID] = storesOfGroup
		} else {
			storesOfGroup = make([]coreStoreModels.Store, 0)
			storesOfGroup = append(storesOfGroup, st)
			storesDisabledMap[st.RestaurantGroupID] = storesOfGroup
		}
	}

	groupsCounter := 0
	storesCounter := 0
	for k, storesOfGroup := range storesDisabledMap {
		groupsCounter++
		fmt.Printf("group_id = %s \n", k)
		for i := range storesOfGroup {
			st := storesOfGroup[i]
			fmt.Printf("storeName = %s \n", st.Name)
			storesCounter++
		}
		fmt.Println()
	}

	fmt.Printf("store groups count = %d, stores count = %d \n", groupsCounter, storesCounter)

	return nil
}

func createStoreService(db *mongo.Database) (store.Service, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	repo, err := store.NewStoreMongoRepository(db)
	if err != nil {
		return nil, err
	}
	s, err := store.NewService(repo)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func createMenuService(db *mongo.Database) (*menu.Service, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	menuRepo, err := menu.NewMenuMongoRepository(db)
	if err != nil {
		return nil, err
	}
	menuService, err := menu.NewMenuService(menuRepo, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return menuService, nil
}
