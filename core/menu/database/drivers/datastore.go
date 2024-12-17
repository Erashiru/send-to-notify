package drivers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
	"github.com/kwaaka-team/orders-core/service/entity_changes_history"
	entityChangesHistoryModels "github.com/kwaaka-team/orders-core/service/entity_changes_history/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataStore interface {
	Base
	TxStarter

	MenuRepository(entityChangesHistoryRepo entity_changes_history.Repository) MenuRepository
	StoreRepository() StoreRepository
	MenuUploadTransactionRepository() MenuUploadTransactionRepository
	StopListTransactionRepository() StopListTransactionRepository
	SequencesRepository() SequencesRepository
	PromoRepository() PromoRepository
	MSPositionsRepository() MSPositionsRepository
	BkOffersRepository() BkOffersRepository
	RestGroupMenuRepository() RestaurantGroupMenuRepository
}

// Base представляет базовый интерфейс для работы с DataStore.
type Base interface {
	// Name - возвращает название DataStore.
	Name() string

	// Ping - проверка на работоспособность.
	Ping() error

	// Close - закрывает соединение с DataStore.
	Close(ctx context.Context) error

	// Connect - устанавливает соединение с DataStore.
	Connect(client *mongo.Client) error

	DataBase() *mongo.Database
	// Client - получение mongo клиента
	Client() *mongo.Client
}

type TxCallback func(error) error

type TxStarter interface {
	StartSession(ctx context.Context) (context.Context, TxCallback, error)
}

type MenuRepository interface {
	List(ctx context.Context, query selector.Menu) ([]models.Menu, error)

	Get(ctx context.Context, query selector.Menu) (models.Menu, error)
	GetMenuIDs(ctx context.Context, query selector.Menu) ([]string, error)
	GetIDByName(ctx context.Context, name string) (string, error)

	Insert(ctx context.Context, menu models.Menu) (string, error)
	Update(ctx context.Context, menu models.Menu, history entityChangesHistoryModels.EntityChangesHistory) error
	Upsert(ctx context.Context, req models.Menu) (models.Menu, error)
	Delete(ctx context.Context, menuID string) error
	GetGroups(ctx context.Context, query selector.Menu) (models.Groups, error)
	UpdateMenuName(ctx context.Context, query models.UpdateMenuName) error
	CreateGlovoSuperCollection(ctx context.Context, menuId string, superCollections dto.MenuSuperCollections) error

	AddRowToAttributeGroup(ctx context.Context, menuId string, attributeMinMax []models.AttributeIdMinMax, attributeGroupID string) error

	UpdateAttributeStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error
	UpdateProductStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error
	UpdateAttributeGroupStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error
	UpdateSectionStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error
	UpdateCollectionStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error
	UpdateSuperCollectionStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error
	UpdateProductStarterAppOfferIDByExtID(ctx context.Context, menuID, extID, starterAppOfferID string) error
	UpdateAttributeStarterAppOfferIDByExtID(ctx context.Context, menuID, extID, starterAppOfferID string) error

	SectionRepository
	ProductRepository
	AttributeRepository
	ComboRepository
}

type RestaurantGroupMenuRepository interface {
	GetMenuByRestGroupId(ctx context.Context, restGroupId string) (models.RestGroupMenu, error)
	UpdateOrCreateMenu(ctx context.Context, restGroupId string, newMenu models.RestGroupMenu) error
}

type ComboRepository interface {
	GetCombos(ctx context.Context, query selector.Menu) ([]models.Combo, int64, error)
}

type StoreRepository interface {
	Get(ctx context.Context, query selector.Store) (storeModels.Store, error)
	List(ctx context.Context, query selector.Store) ([]storeModels.Store, int64, error)
	Update(ctx context.Context, store storeModels.Store) error
	ListStoresByTalabatRestautantID(ctx context.Context, restaurantID string) ([]storeModels.Store, error)
}

type SectionRepository interface {
	GetMenuSection(ctx context.Context, query selector.Menu) ([]*models.Section, error)
	GetMenuSections(ctx context.Context, query selector.Menu) ([]*models.Section, error)
	UpdateSection(ctx context.Context, menuID string, section models.Section) error
}

type ProductRepository interface {
	ListProducts(ctx context.Context, query selector.Menu) ([]models.Product, int64, error)
	GetProduct(ctx context.Context, query selector.Menu) (models.Product, error)
	GetProductsByIDs(ctx context.Context, query selector.Menu, ids []string) ([]models.Product, error)
	GetPromoProducts(ctx context.Context, query selector.Menu) ([]models.Product, error)
	DeleteProducts(ctx context.Context, menuId string, productsIds []string) error
	DeleteProductsFromDB(ctx context.Context, menuId string, productsIds []string) error
	UpdateProductForMatching(ctx context.Context, req models.MatchingProducts) error
	UpdateProductByFields(ctx context.Context, menuId string, productID string, req models.ProductUpdateRequest) error
	GetEmptyProducts(ctx context.Context, menuID string, pagination selector.Pagination) ([]models.Product, int, error)
	UpdateProductAvailableStatus(ctx context.Context, menuID, productID string, status bool) error
}

type AttributeRepository interface {
	GetAttributes(ctx context.Context, query selector.Menu) (models.Attributes, int, error)
	GetAttributeGroups(ctx context.Context, query selector.Menu) (models.AttributeGroups, error)
	DeleteAttributeGroupFromDB(ctx context.Context, menuId string, attrGroupExtId string) error
	ValidateAttributeGroupName(ctx context.Context, menuId, name string) (bool, error)
	CreateAttributeGroup(ctx context.Context, menuID string, attribute models.Attribute) (string, error)
}

type MenuUploadTransactionRepository interface {
	Get(ctx context.Context, query selector.MenuUploadTransaction) (models.MenuUploadTransaction, error)
	List(ctx context.Context, query selector.MenuUploadTransaction) ([]models.MenuUploadTransaction, int64, error)
	Insert(ctx context.Context, req models.MenuUploadTransaction) (string, error)
	Update(ctx context.Context, req models.MenuUploadTransaction) error
	Delete(ctx context.Context, id string) error
}

type StopListTransactionRepository interface {
	Insert(ctx context.Context, req models.StopListTransaction) (string, error)
}

type SequencesRepository interface {
	NextSequenceValue(ctx context.Context, name string) (int, error)
}

type PromoRepository interface {
	GetPromos(ctx context.Context, query selector.Promo) (models.Promo, error)
	FindPromos(ctx context.Context, query selector.Promo) ([]models.Promo, error)
}

type MSPositionsRepository interface {
	GetPositions(ctx context.Context, query selector.MoySklad) ([]models.MoySkladPosition, error)
	RemovePosition(ctx context.Context, query selector.MoySklad) error
	CreatePosition(ctx context.Context, position models.MoySkladPosition) error
}

type BkOffersRepository interface {
	List(ctx context.Context, query selector.BkOffers) ([]models.BkOffers, error)
}
