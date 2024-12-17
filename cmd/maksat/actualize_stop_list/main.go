package main

import (
	"context"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/stoplist"
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

	stopListService, err := stoplist.CreateStopListServiceByWebhook(db, cfg, 1)
	if err != nil {
		return err
	}
	_ = stopListService

	start := time.Now()

	var storeIDs []string

	for i := range storeIDs {
		storeID := storeIDs[i]
		st, err := storeService.GetByID(ctx, storeID)
		if err != nil {
			return err
		}
		if err := stopListService.ActualizeStopListByStoreID(ctx, st.ID); err != nil {
			return err
		}

	}

	end := time.Now()
	log.Printf("ActualizeStopListByPosType took %v", end.Sub(start))

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
