package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/config"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	models2 "github.com/kwaaka-team/orders-core/core/storecore/models"
	woltCli "github.com/kwaaka-team/orders-core/pkg/wolt"
	"github.com/kwaaka-team/orders-core/pkg/wolt/clients"
	"github.com/kwaaka-team/orders-core/pkg/wolt/clients/dto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type WoltService struct {
	repo    drivers.StoreRepository
	baseUrl string
}

func NewWoltService(repo drivers.StoreRepository, baseUrl string) *WoltService {
	return &WoltService{
		repo:    repo,
		baseUrl: baseUrl,
	}
}

func (w WoltService) UpdateStoreStatus(ctx context.Context, responses []models2.StoreManagementResponse) {
	for _, response := range responses {
		if !response.Success {
			continue
		}
		storeUpdate := models2.UpdateStore{
			ID: &response.RestaurantId,
			Wolt: &models2.UpdateStoreWoltConfig{
				IsOpen: &response.IsOpen,
			},
		}
		if err := w.repo.UpdateStoreByFields(ctx, storeUpdate); err != nil {
			log.Info().Msgf("restaurant_id = %s, store_id = %s not updated in DB: %s", response.RestaurantId, response.StoreID, err.Error())
		}
	}
}

func (w WoltService) OpenStore(ctx context.Context, aggregatorStoreID, systemStoreID string) (models2.StoreManagementResponse, error) {
	cli, err := w.createClient(ctx, systemStoreID)
	if err != nil {
		return models2.StoreManagementResponse{}, err
	}

	if err = cli.ManageStore(ctx, dto.IsStoreOpen{
		AvailableStore: models2.ONLINE,
		VenueId:        aggregatorStoreID,
	}); err != nil {
		log.Trace().Err(err).Msg("can't open storeInfo.")
		return models2.StoreManagementResponse{
			ErrMessage:      err.Error(),
			RestaurantId:    systemStoreID,
			StoreID:         aggregatorStoreID,
			IsOpen:          true,
			DeliveryService: models2.WOLT.String(),
		}, err
	}

	return models2.StoreManagementResponse{
		Success:         true,
		RestaurantId:    systemStoreID,
		StoreID:         aggregatorStoreID,
		IsOpen:          true,
		DeliveryService: models2.WOLT.String(),
	}, err
}

func (w WoltService) CloseStore(ctx context.Context, aggregatorStoreID, systemStoreID string) (models2.StoreManagementResponse, error) {
	cli, err := w.createClient(ctx, systemStoreID)
	if err != nil {
		return models2.StoreManagementResponse{}, err
	}

	if err = cli.ManageStore(ctx, dto.IsStoreOpen{
		AvailableStore: models2.OFFLINE,
		VenueId:        aggregatorStoreID,
	}); err != nil {
		log.Trace().Err(err).Msg("can't close storeInfo.")
		return models2.StoreManagementResponse{
			ErrMessage:      err.Error(),
			RestaurantId:    systemStoreID,
			StoreID:         aggregatorStoreID,
			IsOpen:          false,
			DeliveryService: models2.WOLT.String(),
		}, err
	}
	return models2.StoreManagementResponse{
		Success:         true,
		RestaurantId:    systemStoreID,
		StoreID:         aggregatorStoreID,
		IsOpen:          false,
		DeliveryService: models2.WOLT.String(),
	}, err
}

func (w WoltService) createClient(ctx context.Context, storeID string) (clients.Wolt, error) {
	store, err := w.findStore(ctx, selector.NewEmptyStoreSearch().SetID(storeID))
	if err != nil {
		return nil, err
	}
	return woltCli.NewWoltClient(&clients.Config{
		Protocol: config.PROTOCOL,
		BaseURL:  w.baseUrl,
		ApiKey:   store.Wolt.ApiKey,
		Username: store.Wolt.MenuUsername,
		Password: store.Wolt.MenuPassword,
	})
}

func (w WoltService) findStore(ctx context.Context, query selector.Store) (models2.Store, error) {
	store, err := w.repo.Get(ctx, query)
	if err != nil {
		return models2.Store{}, err
	}

	if query.HasExternalDeliveryService() && query.HasExternalStoreID() {
		for _, external := range store.ExternalConfig {
			if external.Type == query.ExternalDeliveryService {
				for _, storeID := range external.StoreID {
					if query.ExternalStoreID == storeID {
						return store, nil
					}
				}
			}
		}

		return models2.Store{}, errors.New("restaurant is not exists")
	}

	return store, nil
}
