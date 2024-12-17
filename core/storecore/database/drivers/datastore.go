package drivers

import (
	"context"
	selector2 "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	models2 "github.com/kwaaka-team/orders-core/core/storecore/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Datastore interface {
	Base

	StoreRepository() StoreRepository
	StoreGroupRepository() StoreGroupRepository
	UserStoreRepository() UserStoreRepository
	StoreTypeRepository() StoreTypeRepository
	ApiTokensRepository() ApiTokensRepository
	VirtualRepository() VirtualRepository
	TapRestaurantRepository() TapRestaurantRepository
}

type Base interface {
	// Name - возвращает название DataStore.
	Name() string

	// Ping - проверка на работоспособность.
	Ping() error

	// Close - закрывает соединение с DataStore.
	Close(ctx context.Context) error

	// Connect - устанавливает соединение с DataStore.
	Connect(cli *mongo.Client) error
}

type StoreGroupRepository interface {
	Get(ctx context.Context, query selector2.StoreGroup) (models2.StoreGroup, error)
	All(ctx context.Context) ([]models2.StoreGroup, error)
	List(ctx context.Context, query selector2.StoreGroup) ([]models2.StoreGroup, error)
	UpdateByFields(ctx context.Context, storeGroup models2.UpdateStoreGroup) (int64, error)
	Create(ctx context.Context, storeGroup models2.StoreGroup) (string, error)
}

type StoreRepository interface {
	Get(ctx context.Context, query selector2.Store) (models2.Store, error)
	List(ctx context.Context, query selector2.Store) ([]models2.Store, error)
	FindCallCenterStores(ctx context.Context) ([]models2.CallCenterRestaurant, error)
	FindDirectStores(ctx context.Context) ([]models2.DirectRestaurant, error)
	Create(ctx context.Context, store models2.Store) (string, error)
	//Update(ctx context.Context, store models2.Store) error
	UpdateStoreByFields(ctx context.Context, store models2.UpdateStore) error
	DeleteStore(ctx context.Context, storeId string) error
	UpdateYandexConfig(ctx context.Context, storeID string, yandexConfig models2.UpdateStoreYandexConfig) error
	CreateYandexConfig(ctx context.Context, storeID string, yandexConfig models2.YandexConfig) error
	Get3plRestaurantStatus(ctx context.Context, storeID string) (bool, error)
	Update3plRestaurantStatus(ctx context.Context, query models2.Update3plRestaurantStatus, indriveStoreID string) error
	UpdateDispatchDeliveryStatus(ctx context.Context, query models2.UpdateDispatchDeliveryAvailable) error
	UpdateWoltBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error
	UpdateDirectBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error
	UpdateCookingTimeWolt(ctx context.Context, restaurantID string, cookingTime int) error
	GetStoresIDsAndNamesByGroupId(ctx context.Context, groupID string) ([]models2.StoreIdAndName, error)
	AppendMenuToStoreMenus(ctx context.Context, storeId string, menu models2.StoreDSMenu) error
	AddAddressCoordinates(ctx context.Context, storeID string, long, lat float64) error
	SetTwoGisLink(ctx context.Context, twoGisLink, restID string) error
	UpdateRestaurantPolygons(ctx context.Context, restaurantID string, polygons []models2.Polygon) error
	UpdateDynamicPolygon(ctx context.Context, restaurantID string, isDynamic bool, cpo float64) error
	UpdateRestaurantCharge(ctx context.Context, req models2.UpdateRestaurantCharge, restID string) error
}

type UserStoreRepository interface {
	Insert(ctx context.Context, userStores []models2.UserStore) error
	FindUsers(ctx context.Context, user selector2.User) ([]models2.UserStore, error)
	Delete(ctx context.Context, user selector2.User) error
	UpdateUserOrderNotifications(ctx context.Context, username, fcmToken string, stores []string) error
}

type StoreTypeRepository interface {
	GetList(ctx context.Context) ([]models2.StoreType, error)
}

type ApiTokensRepository interface {
	GetStores(ctx context.Context, query selector2.ApiToken) ([]models2.Store, error)
}

type VirtualRepository interface {
	GetVirtualStore(ctx context.Context, query selector2.VirtualStore) (models2.VirtualStore, error)
}

type TapRestaurantRepository interface {
	Create(ctx context.Context, req models2.TapRestaurant) (string, error)
	GetList(ctx context.Context, query selector2.TapRestaurant) ([]models2.TapRestaurant, int, error)
	GetByID(ctx context.Context, id string) (models2.TapRestaurant, error)
	GetByQuery(ctx context.Context, query selector2.TapRestaurant) (models2.TapRestaurant, error)
	Update(ctx context.Context, req models2.UpdateTapRestaurant) error
	Delete(ctx context.Context, id string) error
}
