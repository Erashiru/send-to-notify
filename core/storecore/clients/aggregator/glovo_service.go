package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/config"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	models2 "github.com/kwaaka-team/orders-core/core/storecore/models"
	glovoCli "github.com/kwaaka-team/orders-core/pkg/glovo"
	"time"

	glovo "github.com/kwaaka-team/orders-core/pkg/glovo/clients"
	glovoDto "github.com/kwaaka-team/orders-core/pkg/glovo/clients/dto"
	"github.com/rs/zerolog/log"
)

type GlovoService struct {
	repo drivers.StoreRepository
	cli  glovo.Glovo
}

func NewGlovoService(repo drivers.StoreRepository, cfg config.GlovoConfiguration) (*GlovoService, error) {
	cli, err := glovoCli.NewGlovoClient(&glovo.Config{
		Protocol: config.PROTOCOL,
		BaseURL:  cfg.BaseURL,
		ApiKey:   cfg.Token,
	})
	if err != nil {
		log.Trace().Err(err).Msg("can't initialize wolt client.")
		return nil, err
	}

	return &GlovoService{
		repo: repo,
		cli:  cli,
	}, nil
}

func (g GlovoService) UpdateStoreStatus(ctx context.Context, responses []models2.StoreManagementResponse) {
	for _, response := range responses {
		if !response.Success {
			continue
		}
		storeUpdate := models2.UpdateStore{
			ID: &response.RestaurantId,
			Glovo: &models2.UpdateStoreGlovoConfig{
				IsOpen: &response.IsOpen,
			},
		}
		if err := g.repo.UpdateStoreByFields(ctx, storeUpdate); err != nil {
			log.Info().Msgf("restaurant_id = %s, store_id = %s not updated in DB: %s", response.RestaurantId, response.StoreID, err.Error())
		}
	}
}

func (g GlovoService) OpenStore(ctx context.Context, storeID, restaurantId string) (models2.StoreManagementResponse, error) {
	err := g.cli.OpenStore(ctx, glovoDto.StoreManageRequest{
		StoreID: storeID,
		Until:   time.Now().AddDate(1, 0, 0),
	})
	if err != nil {
		log.Trace().Err(err).Msg("can't open store.")
		return models2.StoreManagementResponse{
			ErrMessage:      err.Error(),
			RestaurantId:    restaurantId,
			StoreID:         storeID,
			IsOpen:          true,
			DeliveryService: models2.GLOVO.String(),
		}, err
	}
	return models2.StoreManagementResponse{
		ErrMessage:      "",
		Success:         true,
		RestaurantId:    restaurantId,
		StoreID:         storeID,
		IsOpen:          true,
		DeliveryService: models2.GLOVO.String(),
	}, nil
}

func (g GlovoService) CloseStore(ctx context.Context, storeID, restaurantId string) (models2.StoreManagementResponse, error) {
	err := g.cli.CloseStore(ctx, glovoDto.StoreManageRequest{
		StoreID: storeID,
		Until:   time.Now().AddDate(1, 0, 0),
	})
	if err != nil {
		log.Trace().Err(err).Msg("can't close store.")
		return models2.StoreManagementResponse{
			ErrMessage:      err.Error(),
			RestaurantId:    restaurantId,
			StoreID:         storeID,
			IsOpen:          false,
			DeliveryService: models2.GLOVO.String(),
		}, err
	}
	return models2.StoreManagementResponse{
		ErrMessage:      "",
		Success:         true,
		RestaurantId:    restaurantId,
		StoreID:         storeID,
		IsOpen:          false,
		DeliveryService: models2.GLOVO.String(),
	}, nil
}
