package main

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/core/config"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
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

	st, err := storeService.GetByID(ctx, "645cd660d11b8d74217b9b3d")
	if err != nil {
		return err
	}

	if err = findDuplicatesInMenu(ctx, menuService, st.MenuID); err != nil {
		return err
	}
	for i := range st.Menus {
		menu := st.Menus[i]
		if err = findDuplicatesInMenu(ctx, menuService, menu.ID); err != nil {
			return err
		}
	}

	fmt.Println("finish")

	return nil
}

func findDuplicatesInMenu(ctx context.Context, menuService *menu.Service, menuID string) error {
	menu, err := menuService.FindById(ctx, menuID)
	if err != nil {
		return err
	}
	fmt.Printf("finding duplicates in menu with id %s, delivery %s", menu.ID, menu.Delivery)
	fmt.Println()
	items := make(map[string]struct{})
	duplicates := make([]menuModels.Product, 0)
	for i := range menu.Products {
		product := menu.Products[i]
		id := product.ExtID
		if _, ok := items[id]; ok {
			duplicates = append(duplicates, product)
		}
		items[id] = struct{}{}
	}

	if len(duplicates) > 0 {
		fmt.Println("Duplicates found:")
		for i := range duplicates {
			fmt.Println(duplicates[i].Name)
		}
	}
	fmt.Println()

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
