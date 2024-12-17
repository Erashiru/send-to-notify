package managers

import (
	"context"
	aggregator2 "github.com/kwaaka-team/orders-core/core/storecore/clients/aggregator"
	"github.com/kwaaka-team/orders-core/core/storecore/config"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	models2 "github.com/kwaaka-team/orders-core/core/storecore/models"
	palomaCli "github.com/kwaaka-team/orders-core/pkg/paloma"
	palomaConf "github.com/kwaaka-team/orders-core/pkg/paloma/clients"
	"github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
	"github.com/pkg/errors"
)

var (
	errAggregatorNotFound = errors.New("aggregator not found")
)

type Store interface {
	FindStore(ctx context.Context, query selector.Store) (models2.Store, error)
	FindStores(ctx context.Context, query selector.Store) ([]models2.Store, error)
	FindCallCenterStores(ctx context.Context) ([]models2.CallCenterRestaurant, error)
	FindDirectStores(ctx context.Context) ([]models2.DirectRestaurant, error)
	CreateStore(ctx context.Context, store models2.Store) (string, error)
	//UpdateStore(ctx context.Context, store models2.Store) error
	Update(ctx context.Context, req models2.UpdateStore) error
	ManageStore(ctx context.Context, req models2.StoreManagement) ([]models2.StoreManagementResponse, error)
	DeleteStore(ctx context.Context, storeId string) error
	GetPoints(ctx context.Context, storeId string) ([]models.Point, error)
	UpdateYandexConfig(ctx context.Context, storeID string, yandexConfig models2.UpdateStoreYandexConfig) error
	CreateYandexConfig(ctx context.Context, storeID string, yandexConfig models2.YandexConfig) error
	Get3plRestaurantStatus(ctx context.Context, storeID string) (bool, error)
	Update3plRestaurantStatus(ctx context.Context, query models2.Update3plRestaurantStatus, indriveStoreID string) error
	Update3PlPolygons(ctx context.Context, restaurantID string, polygons []models2.Polygon) error
	UpdateDispatchDeliveryStatus(ctx context.Context, query models2.UpdateDispatchDeliveryAvailable) error
	UpdateWoltBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error
	UpdateDirectBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error
	UpdateCookingTimeWolt(ctx context.Context, restaurantID string, cookingTime int) error
	GetStoresIDsAndNamesByGroupId(ctx context.Context, groupID string) ([]models2.StoreIdAndName, error)
	AppendMenuToStoreMenus(ctx context.Context, storeId string, menu models2.StoreDSMenu) error
	AddAddressCoordinates(ctx context.Context, storeID string, long, lat float64) error
	SetTwoGisLink(ctx context.Context, twoGisLink, restID string) error
	Update3PlDynamic(ctx context.Context, restaurantID string, isDynamic bool, cpo float64) error
	UpdateRestaurantCharge(ctx context.Context, req models2.UpdateRestaurantCharge, restID string) error
}

type StoreManager struct {
	storeRepository drivers.StoreRepository
	glovoService    *aggregator2.GlovoService
	woltService     *aggregator2.WoltService
	config          config.Configuration
}

func NewStoreManager(globalConfig config.Configuration, storeRepository drivers.StoreRepository) (Store, error) {

	wolt := aggregator2.NewWoltService(storeRepository, globalConfig.WoltConfiguration.BaseURL)

	glovo, err := aggregator2.NewGlovoService(storeRepository, globalConfig.GlovoConfiguration)
	if err != nil {
		return nil, err
	}

	return &StoreManager{
		storeRepository: storeRepository,
		woltService:     wolt,
		glovoService:    glovo,
		config:          globalConfig,
	}, nil
}

func (s *StoreManager) CreateStore(ctx context.Context, store models2.Store) (string, error) {
	storeId, err := s.storeRepository.Create(ctx, store)
	if err != nil {
		return "", err
	}

	return storeId, nil
}

func (s *StoreManager) ManageStore(ctx context.Context, req models2.StoreManagement) ([]models2.StoreManagementResponse, error) {

	aggregatorService, err := s.getAggregator(models2.AggregatorName(req.DeliveryService))
	if err != nil {
		return nil, err
	}

	responses := make([]models2.StoreManagementResponse, 0, len(req.StoreInfo))
	for _, reqStore := range req.StoreInfo {

		var response models2.StoreManagementResponse
		if reqStore.StoreStatus {
			response, err = aggregatorService.OpenStore(ctx, reqStore.StoreID, reqStore.RestaurantId)
			if err != nil {
				return nil, err
			}
		} else {
			response, err = aggregatorService.CloseStore(ctx, reqStore.StoreID, reqStore.RestaurantId)
			if err != nil {
				return nil, err
			}
		}

		responses = append(responses, response)
	}

	aggregatorService.UpdateStoreStatus(ctx, responses)

	return responses, nil
}

func (s *StoreManager) FindCallCenterStores(ctx context.Context) ([]models2.CallCenterRestaurant, error) {
	stores, err := s.storeRepository.FindCallCenterStores(ctx)
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (s *StoreManager) FindDirectStores(ctx context.Context) ([]models2.DirectRestaurant, error) {
	stores, err := s.storeRepository.FindDirectStores(ctx)
	if err != nil {
		return nil, err
	}
	return stores, nil
}

//func (s *StoreManager) UpdateStore(ctx context.Context, store models2.Store) error {
//	if err := s.storeRepository.Update(ctx, store); err != nil {
//		return err
//	}
//
//	return nil
//}

func (s *StoreManager) FindStore(ctx context.Context, query selector.Store) (models2.Store, error) {
	store, err := s.storeRepository.Get(ctx, query)
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

func (s *StoreManager) FindStores(ctx context.Context, query selector.Store) ([]models2.Store, error) {
	stores, err := s.storeRepository.List(ctx, query)
	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (s *StoreManager) Update(ctx context.Context, req models2.UpdateStore) error {
	if err := s.storeRepository.UpdateStoreByFields(ctx, req); err != nil {
		return err
	}

	return nil
}

func (s *StoreManager) DeleteStore(ctx context.Context, storeId string) error {
	return s.storeRepository.DeleteStore(ctx, storeId)
}

func (s *StoreManager) getAggregator(aggregatorName models2.AggregatorName) (Aggregator, error) {
	switch aggregatorName {
	case models2.WOLT:
		return s.woltService, nil
	case models2.GLOVO:
		return s.glovoService, nil
	}

	return nil, errAggregatorNotFound
}

func (s *StoreManager) GetPoints(ctx context.Context, storeId string) ([]models.Point, error) {
	store, err := s.FindStore(ctx, selector.Store{
		ID: storeId,
	})
	if err != nil {
		return nil, err
	}

	palomaClient, err := palomaCli.New(&palomaConf.Config{
		Protocol: "http",
		BaseURL:  s.config.PalomaConfiguration.BaseURL,
		ApiKey:   store.Paloma.ApiKey,
		Class:    s.config.PalomaConfiguration.Class,
	})
	if err != nil {
		return nil, err
	}

	res, err := palomaClient.GetPoints(ctx, store.Paloma.ApiKey)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *StoreManager) UpdateYandexConfig(ctx context.Context, storeID string, yandexConfig models2.UpdateStoreYandexConfig) error {
	return s.storeRepository.UpdateYandexConfig(ctx, storeID, yandexConfig)
}

func (s *StoreManager) CreateYandexConfig(ctx context.Context, storeID string, yandexConfig models2.YandexConfig) error {
	return s.storeRepository.CreateYandexConfig(ctx, storeID, yandexConfig)
}

func (s *StoreManager) Get3plRestaurantStatus(ctx context.Context, storeID string) (bool, error) {
	return s.storeRepository.Get3plRestaurantStatus(ctx, storeID)
}

func (s *StoreManager) Update3plRestaurantStatus(ctx context.Context, query models2.Update3plRestaurantStatus, indriveStoreID string) error {
	return s.storeRepository.Update3plRestaurantStatus(ctx, query, indriveStoreID)
}

func (s *StoreManager) Update3PlPolygons(ctx context.Context, restaurantID string, polygons []models2.Polygon) error {
	return s.storeRepository.UpdateRestaurantPolygons(ctx, restaurantID, polygons)
}

func (s *StoreManager) Update3PlDynamic(ctx context.Context, restaurantID string, isDynamic bool, cpo float64) error {
	return s.storeRepository.UpdateDynamicPolygon(ctx, restaurantID, isDynamic, cpo)
}

func (s *StoreManager) UpdateDispatchDeliveryStatus(ctx context.Context, query models2.UpdateDispatchDeliveryAvailable) error {
	return s.storeRepository.UpdateDispatchDeliveryStatus(ctx, query)
}

func (s *StoreManager) UpdateWoltBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error {
	return s.storeRepository.UpdateWoltBusyMode(ctx, storeID, busyMode, busyModeTime)
}

func (s *StoreManager) UpdateDirectBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error {
	return s.storeRepository.UpdateDirectBusyMode(ctx, storeID, busyMode, busyModeTime)
}

func (s *StoreManager) UpdateCookingTimeWolt(ctx context.Context, restaurantID string, cookingTime int) error {
	return s.storeRepository.UpdateCookingTimeWolt(ctx, restaurantID, cookingTime)
}

func (s *StoreManager) GetStoresIDsAndNamesByGroupId(ctx context.Context, groupID string) ([]models2.StoreIdAndName, error) {
	return s.storeRepository.GetStoresIDsAndNamesByGroupId(ctx, groupID)
}

func (s *StoreManager) AppendMenuToStoreMenus(ctx context.Context, storeId string, menu models2.StoreDSMenu) error {
	return s.storeRepository.AppendMenuToStoreMenus(ctx, storeId, menu)
}

func (s *StoreManager) AddAddressCoordinates(ctx context.Context, storeID string, long, lat float64) error {
	return s.storeRepository.AddAddressCoordinates(ctx, storeID, long, lat)
}

func (s *StoreManager) SetTwoGisLink(ctx context.Context, twoGisLink, restID string) error {
	return s.storeRepository.SetTwoGisLink(ctx, twoGisLink, restID)
}

func (s *StoreManager) UpdateRestaurantCharge(ctx context.Context, req models2.UpdateRestaurantCharge, restID string) error {
	return s.storeRepository.UpdateRestaurantCharge(ctx, req, restID)
}
