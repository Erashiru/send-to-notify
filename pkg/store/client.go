package store

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/storecore/config"
	"github.com/kwaaka-team/orders-core/core/storecore/database"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	managers2 "github.com/kwaaka-team/orders-core/core/storecore/managers"
	selector2 "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	palomaModels "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
	"github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/rs/zerolog/log"
)

type Client interface {
	FindStore(ctx context.Context, query dto.StoreSelector) (models.Store, error)
	FindStores(ctx context.Context, query dto.StoreSelector) ([]models.Store, error)
	FindCallCenterStores(ctx context.Context) ([]models.CallCenterRestaurant, error)
	FindDirectStores(ctx context.Context) ([]models.DirectRestaurant, error)
	CreateStore(ctx context.Context, store dto.CreateStoreRequest) (string, error)
	Update(ctx context.Context, req dto.UpdateStore) error
	ManageStore(ctx context.Context, req dto.StoreManagementRequest) ([]dto.StoreManagementResponse, error)
	FindApiTokenStores(ctx context.Context, req dto.ApiTokenSelector) ([]models.Store, error)
	DeleteStore(ctx context.Context, storeID string) error
	FindUserStores(ctx context.Context, query models.UserStore) ([]models.UserStore, error)
	GetPoints(ctx context.Context, storeID string) ([]palomaModels.Point, error)
	FindStoreGroup(ctx context.Context, query selector2.StoreGroup) (models.StoreGroup, error)
	UpdateStoreGroup(ctx context.Context, query dto.UpdateStoreGroup) error
	CreateTapRestaurant(ctx context.Context, req models.TapRestaurant) (string, error)
	GetTapRestaurantList(ctx context.Context, query selector2.TapRestaurant) ([]models.TapRestaurant, int, error)
	GetTapRestaurantByID(ctx context.Context, id string) (models.TapRestaurant, error)
	GetTapRestaurantByName(ctx context.Context, query selector2.TapRestaurant) (models.TapRestaurant, error)
	UpdateTapRestaurant(ctx context.Context, req models.UpdateTapRestaurant) error
	DeleteTapRestaurant(ctx context.Context, id string) error
	UpdateYandexConfig(ctx context.Context, storeID string, yandexConfig models.UpdateStoreYandexConfig) error
	CreateYandexConfig(ctx context.Context, storeID string, yandexConfig models.YandexConfig) error
	Get3plRestaurantStatus(ctx context.Context, storeID string) (bool, error)
	Update3plRestaurantStatus(ctx context.Context, query models.Update3plRestaurantStatus, indriveStoreID string) error
	UpdateRestaurants3PlPolygons(ctx context.Context, restaurantID string, polygons []models.Polygon) error
	UpdateRestaurants3PlDynamicPolygon(ctx context.Context, restaurantID string, isDynamic bool, cpo float64) error
	UpdateDispatchDeliveryStatus(ctx context.Context, query models.UpdateDispatchDeliveryAvailable) error
	UpdateWoltBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error
	UpdateDirectBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error
	UpdateCookingTimeWolt(ctx context.Context, restaurantID string, cookingTime int) error
	GetStoresIDsAndNamesByGroupId(ctx context.Context, groupID string) ([]models.StoreIdAndName, error)
	AppendMenuToStoreMenus(ctx context.Context, storeId string, menu models.StoreDSMenu) error
	AddAddressCoordinates(ctx context.Context, storeID string, long, lat float64) error
	SetTwoGisLink(ctx context.Context, twoGisLink, restID string) error
	UpdateRestaurantCharge(ctx context.Context, req models.UpdateRestaurantCharge, restID string) error
}

type Store struct {
	storeManager         managers2.Store
	storeGroupManager    managers2.StoreGroup
	userStoreManager     managers2.UserStore
	storeTypeManager     managers2.StoreType
	apiTokenManager      managers2.ApiToken
	virtualStoreManager  managers2.VirtualStore
	tapRestaurantManager managers2.TapRestaurant
}

func NewClient(cfg dto.Config) (Client, error) {
	opts, err := config.LoadConfig(context.Background())
	if err != nil {
		return nil, err
	}

	ds, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSDB,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create datastore %s: %v", opts.DSName, err)
	}

	if err = ds.Connect(cfg.MongoCli); err != nil {
		return nil, fmt.Errorf("cannot connect to datastore: %s", err)
	}

	storeManager, err := managers2.NewStoreManager(opts, ds.StoreRepository())
	if err != nil {
		return nil, err
	}

	return &Store{
		storeManager:         storeManager,
		storeGroupManager:    managers2.NewStoreGroupManager(ds),
		userStoreManager:     managers2.NewUserStoreManager(ds),
		storeTypeManager:     managers2.NewStoreTypeManager(ds),
		apiTokenManager:      managers2.NewApiTokenManager(opts, ds),
		virtualStoreManager:  managers2.NewVirtualStoreManager(opts, ds),
		tapRestaurantManager: managers2.NewTapRestaurantManager(ds),
	}, nil
}

func (s *Store) FindStoreGroup(ctx context.Context, query selector2.StoreGroup) (models.StoreGroup, error) {
	storeGroup, err := s.storeGroupManager.FindStoreGroup(ctx, query)
	if err != nil {
		return models.StoreGroup{}, err
	}

	return storeGroup, nil
}

func (s *Store) CreateStore(ctx context.Context, req dto.CreateStoreRequest) (string, error) {
	store := models.Store{
		Name: req.Name,
		Address: models.StoreAddress{
			City:   req.Address.City.Name,
			Street: req.Address.Street,
			Coordinates: models.Coordinates{
				Longitude: req.Address.Coordinates.Longitude,
				Latitude:  req.Address.Coordinates.Latitude,
			},
		},
		Settings: models.Settings{
			TimeZone: models.TimeZone{
				TZ:        req.Address.City.Timezone.Tz,
				UTCOffset: req.Address.City.Timezone.UtcOffset,
			},
			Currency:     req.Currency,
			LanguageCode: req.LanguageCode,
			PriceSource:  models.DELIVERY_SERVICE,
		},
		RestaurantGroupID: req.StoreGroupId,
		QRMenu: models.StoreQRMenuConfig{
			NoTable: req.StoreQRMenuConfig.NoTable,
		},
		LegalEntityId:  req.LegalEntityId,
		SalesManagerId: req.SalesManagerId,
		Telegram: models.StoreTelegramConfig{
			CancelChatID: req.Telegram.CancelChatID,
		},
	}

	for _, contact := range req.Contacts {
		store.Contacts = append(store.Contacts, models.Contact{
			FullName: contact.FullName,
			Position: contact.Position,
			Phone:    contact.Phone,
			Comment:  contact.Comment,
		})
	}

	for _, link := range req.Links {
		store.ExternalLinks = append(store.ExternalLinks, models.Link{
			Name:      link.Name,
			Url:       link.Url,
			ImageLink: link.ImageLink,
		})
	}

	storeId, err := s.storeManager.CreateStore(ctx, store)

	if err != nil {
		return "", err
	}

	err = s.storeManager.Update(ctx, dto.UpdateStore{
		ID:    &storeId,
		Token: &storeId,
	}.ToModel())
	if err != nil {
		return "", err
	}

	restGroup, err := s.storeGroupManager.FindStoreGroupById(ctx, req.StoreGroupId)
	if err != nil {
		return "", err
	}

	var storeIdsNew = []string{storeId}
	storeIdsNew = append(storeIdsNew, restGroup.StoreIds...)

	_, err = s.storeGroupManager.UpdateStores(ctx, dto.UpdateStoreGroup{
		ID:       &req.StoreGroupId,
		StoreIds: storeIdsNew}.ToModel(),
	)
	if err != nil {
		return "", err
	}

	userStores := make([]models.UserStore, 0, len(req.Usernames))
	for _, user := range req.Usernames {
		userStores = append(userStores, models.UserStore{
			Username:     user,
			StoreId:      storeId,
			StoreGroupId: req.StoreGroupId,
		})
	}

	err = s.userStoreManager.Create(ctx, userStores)
	if err != nil {
		return "", err
	}
	return storeId, nil
}

func (s *Store) FindStore(ctx context.Context, query dto.StoreSelector) (models.Store, error) {
	deliveryService, externalDeliveryService := s.compareDeliveryService(query.DeliveryService)
	store, err := s.storeManager.FindStore(ctx, selector2.NewEmptyStoreSearch().
		SetID(query.ID).
		SetToken(query.Token).
		SetClientSecret(query.ClientSecret).
		SetExternalStoreID(query.ExternalStoreID).
		SetDeliveryService(deliveryService).
		SetExternalDeliveryService(externalDeliveryService).
		SetPosType(query.PosType).
		SetHash(query.Hash).
		SetPosOrganizationID(query.PosOrganizationID).
		SetAggregatorMenuID(query.AggregatorMenuID).
		SetAggregatorMenuIDs(query.AggregatorMenuIDs).
		SetIsActiveMenu(query.IsActiveMenu).
		SetExpress24StoreId(query.Express24StoreId).
		SetPosterAccountNumber(query.PosterAccountNumber).
		SetHasVirtualStore(query.HasVirtualStore).
		SetStoreGroupId(query.GroupID).
		SetTalabatRemoteBranchId(query.TalabatRemoteBranchId).
		SetYarosStoreId(query.YarosStoreId).
		SetIsChildStore(query.IsChildStore),
	)

	if err != nil {
		log.Trace().Err(err).Msg("finding store")
		return models.Store{}, err
	}

	return store, nil
}

func (s *Store) FindStores(ctx context.Context, query dto.StoreSelector) ([]models.Store, error) {
	deliveryService, externalDeliveryService := s.compareDeliveryService(query.DeliveryService)

	stores, err := s.storeManager.FindStores(ctx, selector2.NewEmptyStoreSearch().
		SetID(query.ID).
		SetToken(query.Token).
		SetClientSecret(query.ClientSecret).
		SetExternalStoreID(query.ExternalStoreID).
		SetDeliveryService(deliveryService).
		SetExternalDeliveryService(externalDeliveryService).
		SetPosType(query.PosType).
		SetHash(query.Hash).
		SetAggregatorMenuID(query.AggregatorMenuID).
		SetAggregatorMenuIDs(query.AggregatorMenuIDs).
		SetIsActiveMenu(query.IsActiveMenu).
		SetStoreIDs(query.IDs).
		SetHasVirtualStore(query.HasVirtualStore).
		SetStoreGroupId(query.GroupID).
		SetScheduledStatusChange(query.HasScheduledStatusChange).
		SetCity(query.City).
		SetPosterAccountNumber(query.PosterAccountNumber).
		SetOrderAutoClose(query.OrderAutoClose),
	)
	if err != nil {
		log.Trace().Err(err).Msg("finding stores")
		return nil, err
	}

	return stores, nil
}

func (s *Store) FindCallCenterStores(ctx context.Context) ([]models.CallCenterRestaurant, error) {
	stores, err := s.storeManager.FindCallCenterStores(ctx)
	if err != nil {
		log.Trace().Err(err).Msg("finding call center stores")
		return nil, err
	}
	return stores, nil
}

func (s *Store) FindDirectStores(ctx context.Context) ([]models.DirectRestaurant, error) {
	stores, err := s.storeManager.FindDirectStores(ctx)
	if err != nil {
		log.Trace().Err(err).Msg("finding direct stores")
		return nil, err
	}
	return stores, nil
}

func (s *Store) ManageStore(ctx context.Context, req dto.StoreManagementRequest) ([]dto.StoreManagementResponse, error) {
	response, err := s.storeManager.ManageStore(ctx, req.ToModel())
	if err != nil {
		return nil, err
	}

	return dto.FromStoreManagementModels(response), nil
}

func (s *Store) Update(ctx context.Context, req dto.UpdateStore) error {
	if err := s.storeManager.Update(ctx, req.ToModel()); err != nil {
		return err
	}

	return nil
}

func (s *Store) FindApiTokenStores(ctx context.Context, req dto.ApiTokenSelector) ([]models.Store, error) {
	stores, err := s.apiTokenManager.FindStores(ctx, selector2.NewEmptyApiTokenSearch().SetToken(req.ApiToken))

	if err != nil {
		log.Trace().Err(err).Msgf("finding apiToken stores %v", err)
		return nil, err
	}

	return stores, nil
}

func (s *Store) compareDeliveryService(deliveryService string) (string, string) {
	switch deliveryService {
	case "glovo", "wolt", "qr_menu", "moysklad", "chocofood", "express24", "deliveroo":
		return deliveryService, ""
	default:
		return "", deliveryService
	}
}

func (s *Store) DeleteStore(ctx context.Context, storeID string) error {
	return s.storeManager.DeleteStore(ctx, storeID)
}

func (u *Store) FindUserStores(ctx context.Context, query models.UserStore) ([]models.UserStore, error) {
	res, err := u.userStoreManager.FindUsers(ctx, selector2.NewEmptyUserSearch().
		SetID(query.ID).
		SetStoreID(query.StoreId).
		SetStoreGroupID(query.StoreGroupId).
		SetUsername(query.Username).
		SetSendNotification(query.SendNotification),
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Store) GetPoints(ctx context.Context, storeID string) ([]palomaModels.Point, error) {
	return s.storeManager.GetPoints(ctx, storeID)
}

func (s *Store) UpdateStoreGroup(ctx context.Context, query dto.UpdateStoreGroup) error {
	_, err := s.storeGroupManager.UpdateStores(ctx, query.ToModel())
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateTapRestaurant(ctx context.Context, req models.TapRestaurant) (string, error) {
	return s.tapRestaurantManager.CreateTapRestaurant(ctx, req)
}

func (s *Store) GetTapRestaurantList(ctx context.Context, query selector2.TapRestaurant) ([]models.TapRestaurant, int, error) {
	return s.tapRestaurantManager.GetTapRestaurantList(ctx, query)
}

func (s *Store) GetTapRestaurantByID(ctx context.Context, id string) (models.TapRestaurant, error) {
	return s.tapRestaurantManager.GetTapRestaurant(ctx, id)
}

func (s *Store) GetTapRestaurantByName(ctx context.Context, query selector2.TapRestaurant) (models.TapRestaurant, error) {
	return s.tapRestaurantManager.GetTapRestaurantByName(ctx, query)
}

func (s *Store) UpdateTapRestaurant(ctx context.Context, req models.UpdateTapRestaurant) error {
	return s.tapRestaurantManager.UpdateTapRestaurant(ctx, req)
}

func (s *Store) DeleteTapRestaurant(ctx context.Context, id string) error {
	return s.tapRestaurantManager.DeleteTapRestaurant(ctx, id)
}

func (s *Store) UpdateYandexConfig(ctx context.Context, storeID string, yandexConfig models.UpdateStoreYandexConfig) error {
	return s.storeManager.UpdateYandexConfig(ctx, storeID, yandexConfig)
}

func (s *Store) CreateYandexConfig(ctx context.Context, storeID string, yandexConfig models.YandexConfig) error {
	return s.storeManager.CreateYandexConfig(ctx, storeID, yandexConfig)
}

func (s *Store) Get3plRestaurantStatus(ctx context.Context, storeID string) (bool, error) {
	return s.storeManager.Get3plRestaurantStatus(ctx, storeID)
}

func (s *Store) Update3plRestaurantStatus(ctx context.Context, query models.Update3plRestaurantStatus, indriveStoreID string) error {
	return s.storeManager.Update3plRestaurantStatus(ctx, query, indriveStoreID)
}

func (s *Store) UpdateRestaurants3PlPolygons(ctx context.Context, restaurantID string, polygons []models.Polygon) error {
	return s.storeManager.Update3PlPolygons(ctx, restaurantID, polygons)
}

func (s *Store) UpdateRestaurants3PlDynamicPolygon(ctx context.Context, restaurantID string, isDynamic bool, cpo float64) error {
	return s.storeManager.Update3PlDynamic(ctx, restaurantID, isDynamic, cpo)
}

func (s *Store) UpdateDispatchDeliveryStatus(ctx context.Context, query models.UpdateDispatchDeliveryAvailable) error {
	return s.storeManager.UpdateDispatchDeliveryStatus(ctx, query)
}

func (s *Store) UpdateWoltBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error {
	return s.storeManager.UpdateWoltBusyMode(ctx, storeID, busyMode, busyModeTime)
}

func (s *Store) UpdateDirectBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error {
	return s.storeManager.UpdateDirectBusyMode(ctx, storeID, busyMode, busyModeTime)
}

func (s *Store) UpdateCookingTimeWolt(ctx context.Context, restaurantID string, cookingTime int) error {
	return s.storeManager.UpdateCookingTimeWolt(ctx, restaurantID, cookingTime)
}

func (s *Store) GetStoresIDsAndNamesByGroupId(ctx context.Context, groupID string) ([]models.StoreIdAndName, error) {
	return s.storeManager.GetStoresIDsAndNamesByGroupId(ctx, groupID)
}

func (s *Store) AppendMenuToStoreMenus(ctx context.Context, storeId string, menu models.StoreDSMenu) error {
	return s.storeManager.AppendMenuToStoreMenus(ctx, storeId, menu)
}

func (s *Store) AddAddressCoordinates(ctx context.Context, storeID string, long, lat float64) error {
	return s.storeManager.AddAddressCoordinates(ctx, storeID, long, lat)
}

func (s *Store) SetTwoGisLink(ctx context.Context, twoGisLink, restID string) error {
	return s.storeManager.SetTwoGisLink(ctx, twoGisLink, restID)
}

func (s *Store) UpdateRestaurantCharge(ctx context.Context, req models.UpdateRestaurantCharge, restID string) error {
	return s.storeManager.UpdateRestaurantCharge(ctx, req, restID)
}
