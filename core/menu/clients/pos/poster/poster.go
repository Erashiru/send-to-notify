package poster

import (
	"context"
	"errors"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/base"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	posterCli "github.com/kwaaka-team/orders-core/pkg/poster"
	posterConf "github.com/kwaaka-team/orders-core/pkg/poster/clients"
)

type manager struct {
	cfg      menu.Configuration
	cli      posterConf.Poster
	menuRepo drivers.MenuRepository
}

func NewManager(
	cfg menu.Configuration,
	menuRepo drivers.MenuRepository,
	store storeModels.Store) (base.Manager, error) {

	cli, err := posterCli.NewClient(&posterConf.Config{
		Protocol: "http",
		BaseURL:  "https://joinposter.com", //from secrets
		Token:    store.Poster.Token,
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

func (man manager) GetMenu(ctx context.Context, store storeModels.Store) (models.Menu, error) {
	menu, err := man.cli.GetProducts(ctx)
	if err != nil {
		return models.Menu{}, err
	}
	products, err := man.existProducts(ctx, store.MenuID)
	if err != nil {
		return models.Menu{}, err
	}

	return menuFromClient(menu, store, products)
}
