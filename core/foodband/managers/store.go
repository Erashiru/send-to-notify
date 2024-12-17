package managers

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/foodband/models"
	"github.com/kwaaka-team/orders-core/core/foodband/resources/http/v1/dto"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/domain/logger"
	storeCore "github.com/kwaaka-team/orders-core/pkg/store"
	storeCoreModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Store interface {
	GetApiTokenStores(ctx context.Context, apiToken, storeId string) ([]dto.Store, string, error)
	ManageStoreInAggregator(ctx context.Context, req models.ManageAggregatorStoreRequest) error
}

type storeImplementation struct {
	storeCli storeCore.Client
	logger   *zap.SugaredLogger
}

func NewStoreManager(storeCli storeCore.Client, logger *zap.SugaredLogger) Store {
	return &storeImplementation{
		storeCli: storeCli,
		logger:   logger,
	}
}

func (man *storeImplementation) ManageStoreInAggregator(ctx context.Context, req models.ManageAggregatorStoreRequest) error {
	if req.DeliveryService == "" || req.PosIntegrationStoreID == "" {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: errors.New("invalid deliveryService or storeID"),
		})
		return errors.New("invalid deliveryService or storeID")
	}

	store, err := man.storeCli.FindStore(ctx, storeCoreModels.StoreSelector{
		ExternalStoreID: req.PosIntegrationStoreID,
		DeliveryService: "foodband",
	})
	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: err,
		})
		return err
	}

	storeIDs := store.GetAggregatorStoreIDs(req.DeliveryService)
	if len(storeIDs) == 0 {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: errors.New(fmt.Sprintf("not found store ids for %s", req.DeliveryService)),
		})
		return errors.New(fmt.Sprintf("not found store ids for %s", req.DeliveryService))
	}

	storeInfos := make([]storeCoreModels.StoreInfo, 0, len(storeIDs))
	for _, v := range storeIDs {
		storeInfos = append(storeInfos, storeCoreModels.StoreInfo{
			RestaurantId: store.ID,
			StoreId:      v,
			StoreStatus:  req.IsOpen,
		})
	}

	_, err = man.storeCli.ManageStore(ctx, storeCoreModels.StoreManagementRequest{
		DeliveryService: req.DeliveryService,
		StoreInfos:      storeInfos,
	})

	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: err,
		})
		return err
	}

	man.logger.Info(logger.LoggerInfo{
		System:   "foodband response",
		Response: fmt.Sprintf("manage store successful, restaurantID: %s, DS: %s, storeID: %s, isOpen: %v", store.ID, req.DeliveryService, req.PosIntegrationStoreID, req.IsOpen),
	})
	return nil
}

func (man *storeImplementation) GetApiTokenStores(ctx context.Context, apiToken, storeId string) ([]dto.Store, string, error) {
	if apiToken == "" {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: nil,
		})
		return nil, "", fmt.Errorf("invalid authorization token")
	}
	stores, err := man.storeCli.FindApiTokenStores(ctx, storeCoreModels.ApiTokenSelector{ApiToken: apiToken})
	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: err,
		})
		return nil, "", err
	}

	var res []dto.Store
	var restaurantId string
	for _, s := range stores {
		if s.PosType != "foodband" {
			continue
		}
		fbStore := fromStoreCore(s)
		if fbStore.ID == "" {
			continue
		}
		res = append(res, fbStore)

		if storeId == fbStore.ID {
			restaurantId = s.ID
		}
	}
	man.logger.Info(logger.LoggerInfo{
		System: "foodband response",
		Response: []interface{}{
			restaurantId,
			res,
		},
	})

	return res, restaurantId, nil
}

func fromStoreCore(store coreStoreModels.Store) dto.Store {
	fbStoreID := ""
	for _, ext := range store.ExternalConfig {
		if ext.Type != "foodband" {
			continue
		}
		if ext.StoreID[0] == "" {
			continue
		}
		fbStoreID = ext.StoreID[0]
		break
	}

	res := dto.Store{
		ID:               fbStoreID,
		Name:             store.Name,
		PosType:          store.PosType,
		DeliveryServices: fromDeliveryServices(store),
	}

	return res
}

func fromDeliveryServices(store coreStoreModels.Store) []string {
	var deliveriServices []string

	if len(store.Glovo.StoreID) != 0 && store.Glovo.StoreID[0] != "" {
		deliveriServices = append(deliveriServices, "glovo")
	}
	if len(store.Wolt.StoreID) != 0 && store.Wolt.StoreID[0] != "" {
		deliveriServices = append(deliveriServices, "wolt")
	}
	if len(store.Chocofood.StoreID) != 0 && store.Chocofood.StoreID[0] != "" {
		deliveriServices = append(deliveriServices, "chocofood")
	}
	if store.MoySklad.UserName != "" && store.MoySklad.Password != "" {
		deliveriServices = append(deliveriServices, "moysklad")
	}

	for _, ext := range store.ExternalConfig {
		if ext.ServiceType != "aggregator" && ext.Type != "yandex" && ext.Type != "emenu" {
			continue
		}
		if len(ext.StoreID) != 0 && ext.StoreID[0] != "" {
			deliveriServices = append(deliveriServices, ext.Type)
		}
	}

	return deliveriServices
}
