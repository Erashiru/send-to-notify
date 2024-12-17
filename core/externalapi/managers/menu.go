package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/externalapi/database/drivers"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/core/externalapi/utils"
	helperUtils "github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/menu"
	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
	"github.com/kwaaka-team/orders-core/pkg/store"
	storeDto "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/rs/zerolog/log"
)

type MenuClient interface {
	GetMenu(ctx context.Context, storeID, service, clientSecret string) (models.Menu, error)
	GetStores(ctx context.Context, service, clientSecret string) (models.GetStoreResponse, error)
	GetPromos(ctx context.Context, storeID, service, clientSecret string) (models.Promo, error)
	FindStore(ctx context.Context, storeID, service string) (coreStoreModels.Store, error)
	GetRetailMenu(ctx context.Context, storeID, service, clientSecret string) (models.RetailMenu, error)
}

type MenuClientManager struct {
	ds       drivers.DataStore
	menuCli  menu.Client
	storeCli store.Client
}

func NewMenuClientManager(menuCli menu.Client, storeCli store.Client) MenuClient {
	return &MenuClientManager{
		menuCli:  menuCli,
		storeCli: storeCli,
	}
}

func (manager *MenuClientManager) FindStore(ctx context.Context, storeID, service string) (coreStoreModels.Store, error) {
	store, err := manager.storeCli.FindStore(ctx, storeDto.StoreSelector{
		DeliveryService: service,
		ExternalStoreID: storeID,
	})
	if err != nil {
		return coreStoreModels.Store{}, err
	}

	return store, nil
}

func (manager *MenuClientManager) GetMenu(ctx context.Context, storeID, service, clientSecret string) (models.Menu, error) {
	_, err := manager.getStore(ctx, storeID, service, clientSecret)
	if err != nil {
		log.Trace().Err(err).Msg("Can't get store")
		return models.Menu{}, err
	}

	req, err := manager.menuCli.GetMenu(ctx, storeID, dto.DeliveryService(service))
	if err != nil {
		log.Trace().Err(err).Msg("Can't get menu")
		return models.Menu{}, err
	}

	result := utils.ParseMenu(req)

	helperUtils.Beautify("external response menu body", result)

	return result, nil
}

func (manager *MenuClientManager) getStore(ctx context.Context, storeID, service, clientSecret string) (models.Place, error) {
	store, err := manager.storeCli.FindStore(ctx, storeDto.StoreSelector{
		DeliveryService: service,
		ExternalStoreID: storeID,
		ClientSecret:    clientSecret,
	})
	if err != nil {
		log.Trace().Err(err).Msg("Can't get store")
		return models.Place{}, err
	}

	return models.Place{
		Id:      storeID,
		Title:   store.Name,
		Address: store.Address.Street,
	}, nil
}

func (manager *MenuClientManager) GetStores(ctx context.Context, service, clientSecret string) (models.GetStoreResponse, error) {
	stores, err := manager.storeCli.FindStores(ctx, storeDto.StoreSelector{
		DeliveryService: service,
		ClientSecret:    clientSecret,
	})
	if err != nil {
		log.Trace().Err(err).Msg("Can't get stores")
		return models.GetStoreResponse{}, err
	}

	return utils.ParseStores(stores, service), nil
}

func (manager *MenuClientManager) GetPromos(ctx context.Context, storeID, service, clientSecret string) (models.Promo, error) {
	_, err := manager.getStore(ctx, storeID, service, clientSecret)
	if err != nil {
		log.Trace().Err(err).Msg("Can't get store")
		return models.Promo{}, err
	}

	promo, err := manager.menuCli.GetPromos(ctx, storeID, dto.DeliveryService(service))
	if err != nil {
		log.Trace().Err(err).Msg("Can't get promos")
		return models.Promo{}, err
	}

	return utils.ParsePromos(promo), err
}

func (manager *MenuClientManager) GetRetailMenu(ctx context.Context, storeID, service, clientSecret string) (models.RetailMenu, error) {
	_, err := manager.getStore(ctx, storeID, service, clientSecret)
	if err != nil {
		log.Trace().Err(err).Msg("Can't get store")
		return models.RetailMenu{}, err
	}

	req, err := manager.menuCli.GetMenu(ctx, storeID, dto.DeliveryService(service))
	if err != nil {
		log.Trace().Err(err).Msg("Can't get menu")
		return models.RetailMenu{}, err
	}

	result := utils.ParseRetailMenu(req)

	helperUtils.Beautify("external response menu body", result)

	return result, nil
}
