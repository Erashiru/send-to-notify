package main

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/stoplist"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Err(err).Msgf("fatal error")
		return
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

	menuRepo, err := menu.NewMenuMongoRepository(db)
	if err != nil {
		return err
	}
	menuService, err := menu.NewMenuService(menuRepo, nil, nil, nil)
	if err != nil {
		return err
	}

	storeService, err := createStoreService(db)
	if err != nil {
		return err
	}

	stopListService, err := stoplist.CreateStopListServiceByWebhook(db, cfg, 1)
	if err != nil {
		return err
	}
	_ = stopListService

	workerPool, err := NewWorkerPool(10, stopListService, menuService, menuRepo)
	if err != nil {
		return err
	}

	inputStores, err := getStoresByStoreGroupIDs(ctx, storeService)
	if err != nil {
		return err
	}
	if len(inputStores) == 0 {
		err = errors.New("stores length is 0")
		log.Err(err).Msgf("")
		return err
	}
	before := time.Now()

	workerPool.Start(ctx)
	workerPool.Submit(inputStores)

	log.Printf("reset is_disabled status, stores count = %d, time took %v", len(inputStores), time.Since(before))
	log.Printf("success")

	return nil
}

func getStoresByStoreGroupIDs(ctx context.Context, storeService store.Service) ([]coreStoreModels.Store, error) {
	var storeGroupIDs []string
	_ = storeGroupIDs
	storeGroupIDs = generateStoreGroups()

	inputStores := make([]coreStoreModels.Store, 0, len(storeGroupIDs))
	for i := range storeGroupIDs {
		storeGroupID := storeGroupIDs[i]
		stores, err := storeService.GetStoresByStoreGroupID(ctx, storeGroupID)
		if err != nil {
			log.Err(err).Msgf("finding store error, store_id = %s", storeGroupID)
			continue
		}
		if len(stores) == 0 {
			return nil, fmt.Errorf("stores not found by storeGroupID %s", storeGroupID)
		}
		inputStores = append(inputStores, stores...)
	}

	return inputStores, nil
}

func generateStoreGroups() []string {
	var storeGroupIDs []string

	return storeGroupIDs
}

func process(ctx context.Context, st coreStoreModels.Store, stopListService stoplist.Service,
	menuService *menu.Service, menuRepo menu.Repository) error {

	if err := resetDisabledStatus(ctx, menuService, menuRepo, st); err != nil {
		return err
	}

	if err := stopListService.ActualizeStopListByStoreID(ctx, st.ID); err != nil {
		return err
	}

	return nil
}

func resetDisabledStatus(ctx context.Context, menuService *menu.Service, menuRepo menu.Repository, st coreStoreModels.Store) error {
	if err := resetDisabledStatusByMenuID(ctx, menuService, menuRepo, st.MenuID); err != nil {
		return err
	}
	for i := range st.Menus {
		aggMenu := st.Menus[i]
		if err := resetDisabledStatusByMenuID(ctx, menuService, menuRepo, aggMenu.ID); err != nil {
			return err
		}
	}
	return nil
}
func resetDisabledStatusByMenuID(ctx context.Context, menuService *menu.Service, menuRepo menu.Repository, menuID string) error {
	menu, err := menuService.FindById(ctx, menuID)
	if err != nil {
		return err
	}

	if err = resetDisabledStatusProducts(ctx, menuRepo, menu); err != nil {
		return err
	}
	if err = resetDisabledStatusAttributes(ctx, menuRepo, menu); err != nil {
		return err
	}

	return nil
}
func resetDisabledStatusProducts(ctx context.Context, menuRepo menu.Repository, menu *models.Menu) error {
	productIDs := make([]string, 0)
	for i := range menu.Products {
		product := menu.Products[i]
		if !product.IsDisabled {
			continue
		}
		productIDs = append(productIDs, product.ExtID)
	}
	if len(productIDs) == 0 {
		return nil
	}
	return menuRepo.BulkUpdateProductsDisabledStatus(ctx, menu.ID, productIDs, false)
}

func resetDisabledStatusAttributes(ctx context.Context, menuRepo menu.Repository, menu *models.Menu) error {
	attributeIDs := make([]string, 0)
	for i := range menu.Attributes {
		attribute := menu.Attributes[i]
		if !attribute.IsDisabled {
			continue
		}
		attributeIDs = append(attributeIDs, attribute.ExtID)
	}
	if len(attributeIDs) == 0 {
		return nil
	}
	return menuRepo.BulkUpdateAttributesDisabledStatus(ctx, menu.ID, attributeIDs, false)
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
