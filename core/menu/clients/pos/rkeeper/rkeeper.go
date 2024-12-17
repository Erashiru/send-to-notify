package rkeeper

import (
	"context"
	"errors"
	"fmt"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/base"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/rs/zerolog/log"
	"strings"
	"time"

	rkeeperCli "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite"
	rkeeperConf "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients"
)

type manager struct {
	cfg      menu.Configuration
	cli      rkeeperConf.RKeeper
	menuRepo drivers.MenuRepository
}

func NewManager(
	cfg menu.Configuration,
	menuRepo drivers.MenuRepository, store storeModels.Store) (base.Manager, error) {
	var apyKey string
	switch {
	case store.RKeeper.ApiKey != "":
		apyKey = store.RKeeper.ApiKey
	default:
		apyKey = cfg.RkeeperConfiguration.ApiKey
	}

	cli, err := rkeeperCli.NewRKeeperClient(&rkeeperConf.Config{
		Protocol: "http",
		BaseURL:  cfg.RkeeperConfiguration.BaseURL,
		ApiKey:   apyKey,
	})
	if err != nil {
		return nil, err
	}

	return &manager{
		cli:      cli,
		cfg:      cfg,
		menuRepo: menuRepo,
	}, nil
}

func (man manager) updateMenu(ctx context.Context, store storeModels.Store) error {
	_, err := man.cli.UpdateMenu(ctx, store.RKeeper.ObjectId)
	if err != nil {
		return fmt.Errorf("couldn't update menu %w", err)
	}

	log.Info().Msgf("success update menu for rkeeper object id %d", store.RKeeper.ObjectId)

	return nil
}

func (man manager) GetAggMenu(ctx context.Context, store storeModels.Store) ([]models.Menu, error) {
	return nil, errors.New("method not implemented")
}

func (man manager) GetMenu(ctx context.Context, store storeModels.Store) (models.Menu, error) {
	if havePriceTypeID(store.RKeeper.PriceTypeID) {
		for i := 0; i < 5; i++ {
			time.Sleep(10 * time.Second)
			menu, err := man.GetMenuByParams(ctx, store)
			if err != nil {
				return models.Menu{}, err
			}
			return menu, nil
		}
	}

	if err := man.updateMenu(ctx, store); err != nil {
		return models.Menu{}, err
	}

	menu, err := man.cli.GetMenu(ctx, store.RKeeper.ObjectId)
	if err != nil {
		return models.Menu{}, err
	}

	// exist products in DB
	products, posProducts, err := man.existProducts(ctx, store.MenuID)
	if err != nil {
		return models.Menu{}, err
	}

	stopList, err := man.cli.GetStopList(ctx, store.RKeeper.ObjectId)
	if err != nil {
		return models.Menu{}, fmt.Errorf("couldn't get stoplist %w", err)
	}

	return menuFromClient(menu.TaskResponse.Menu, products, posProducts, stopList, store), nil
}

func (man manager) existProducts(ctx context.Context, menuID string) (map[string]string, map[string]models.Product, error) {

	if menuID == "" {
		return map[string]string{}, map[string]models.Product{}, nil
	}

	// get products from main menu if exist
	products, _, err := man.menuRepo.ListProducts(ctx, selector.EmptyMenuSearch().
		SetMenuID(menuID))
	if err != nil {
		return nil, map[string]models.Product{}, nil
	}

	// add to hash map
	productExist := make(map[string]string, len(products))
	posProducts := make(map[string]models.Product, len(products))
	for _, product := range products {
		key := strings.TrimSpace(product.ProductID + product.ParentGroupID)

		// cause has cases if product_id && parent_id same, size_id different
		productExist[key] = product.ExtID
		posProducts[product.ProductID] = product
	}

	return productExist, posProducts, nil
}

func havePriceTypeID(priceTypeID int) bool {
	if priceTypeID == 0 {
		return false
	}
	return true
}

func (man manager) GetMenuByParams(ctx context.Context, store storeModels.Store) (models.Menu, error) {
	menu, err := man.cli.GetMenuByParams(ctx, store.RKeeper.ObjectId, store.RKeeper.PriceTypeID)
	if err != nil {
		return models.Menu{}, err
	}

	// exist products in DB
	products, posProducts, err := man.existProducts(ctx, store.MenuID)
	if err != nil {
		return models.Menu{}, err
	}

	stopList, err := man.cli.GetStopList(ctx, store.RKeeper.ObjectId)
	if err != nil {
		return models.Menu{}, fmt.Errorf("couldn't get stoplist %w", err)
	}

	return menuFromClient(menu.TaskResponse.Menu, products, posProducts, stopList, store), nil
}
