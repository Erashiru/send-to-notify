package yaros

import (
	"context"
	"errors"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/base"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	yarosCli "github.com/kwaaka-team/orders-core/pkg/yaros"
	yarosConf "github.com/kwaaka-team/orders-core/pkg/yaros/clients"
)

type manager struct {
	cfg      menu.Configuration
	cli      yarosConf.Yaros
	menuRepo drivers.MenuRepository
}

func NewManager(
	cfg menu.Configuration,
	menuRepo drivers.MenuRepository,
	store storeModels.Store) (base.Manager, error) {

	cli, err := yarosCli.NewClient(&yarosConf.Config{
		Protocol: "http",
		BaseURL:  cfg.YarosConfiguration.BaseURL,
		Username: store.Yaros.Username,
		Password: store.Yaros.Password,
		RestID:   store.Yaros.StoreId,
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

func (man manager) GetAggMenu(ctx context.Context, store storeModels.Store) ([]models.Menu, error) {
	return nil, errors.New("method not implemented")
}

func (m manager) GetMenu(ctx context.Context, store storeModels.Store) (models.Menu, error) {
	items, err := m.cli.GetItems(ctx, store.Yaros.StoreId)
	if err != nil {
		return models.Menu{}, err
	}
	categories, err := m.cli.GetCategories(ctx, store.Yaros.StoreId)
	if err != nil {
		return models.Menu{}, err
	}
	return menuFromClient(items, categories)
}
