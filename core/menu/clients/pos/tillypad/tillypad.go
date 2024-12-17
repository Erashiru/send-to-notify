package tillypad

import (
	"context"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/base"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/tillypad/yandexDeliveryProtocolTillypad"
	"github.com/kwaaka-team/orders-core/pkg/tillypad/yandexDeliveryProtocolTillypad/clients"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type manager struct {
	globalConfig menu.Configuration
	menuRepo     drivers.MenuRepository
	cli          clients.Client
}

func NewTillypadManager(globalConfig menu.Configuration, menuRepo drivers.MenuRepository, store storeModels.Store) (base.Manager, error) {

	cli, err := yandexDeliveryProtocolTillypad.NewTillypadClient(clients.Config{
		BaseURL:      globalConfig.TillypadConfiguration.BaseUrl,
		Protocol:     "http",
		ClientId:     store.TillyPad.ClientId,
		ClientSecret: store.TillyPad.ClientSecret,
		PathPrefix:   store.TillyPad.PathPrefix,
	})
	if err != nil {
		log.Trace().Err(err).Msg("can't initialize Tillypad Client")
		return nil, err
	}

	return &manager{
		globalConfig: globalConfig,
		menuRepo:     menuRepo,
		cli:          cli,
	}, nil

}

func (man manager) GetMenu(ctx context.Context, store storeModels.Store) (models.Menu, error) {

	menu, err := man.cli.GetMenu(ctx, store.TillyPad.PointId)
	if err != nil {
		log.Info().Msgf("get menu from tillypad client point id:%s error:%s ", store.TillyPad.PointId, err)
		return models.Menu{}, err
	}

	return man.menuFromClient(menu, store), nil
}

func (man manager) getExistProducts(ctx context.Context, menuId string) (map[string]models.Product, error) {

	if menuId == "" {
		return nil, nil
	}

	products, _, err := man.menuRepo.ListProducts(ctx, selector.EmptyMenuSearch().SetMenuID(menuId))
	if err != nil {
		return nil, err
	}

	productExist := make(map[string]models.Product, len(products))
	for _, product := range products {
		productExist[product.ProductID] = product
	}

	return productExist, nil
}

func (man manager) GetAggMenu(ctx context.Context, store storeModels.Store) ([]models.Menu, error) {
	return nil, errors.New("method not implemented")
}
