package menu

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/config/menu"
	posIntegrationModels "github.com/kwaaka-team/orders-core/core/integration_api/models"
	"github.com/kwaaka-team/orders-core/core/menu/database"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/managers"
	"github.com/kwaaka-team/orders-core/pkg/menu/utils"

	"github.com/kwaaka-team/orders-core/core/menu/managers/validator"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/entity_changes_history"
	entityChangesHistoryModels "github.com/kwaaka-team/orders-core/service/entity_changes_history/models"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"

	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
	notifyClient "github.com/kwaaka-team/orders-core/pkg/que"
	storeCore "github.com/kwaaka-team/orders-core/pkg/store"
	storeDto "github.com/kwaaka-team/orders-core/pkg/store/dto"
)

type Client interface {
	GetMenu(ctx context.Context, externalStoreID string, deliveryService dto.DeliveryService) (models.Menu, error)
	GetPromos(ctx context.Context, externalStoreID string, deliveryService dto.DeliveryService) (models.Promo, error)
	GetStorePromos(ctx context.Context, query dto.GetPromosSelector) ([]dto.PromoDiscount, error)
	GetStores(ctx context.Context, deliveryService dto.DeliveryService) ([]storeModels.Store, error)
	GetMenuByID(ctx context.Context, menuID string) (models.Menu, error)
	GetMenuGroups(ctx context.Context, menuID string) ([]dto.MenuGroup, error)
	GetMenuStatus(ctx context.Context, storeId string, isDeleted bool) ([]storeDto.StoreDsMenuDto, error)
	GetPosDiscounts(ctx context.Context, storeID string, deliveryService dto.DeliveryService) (dto.PosDiscount, error)

	ListStoresByProduct(ctx context.Context, req dto.GetStoreByProductRequest) ([]storeModels.Store, int64, error)

	UpsertMenu(ctx context.Context, req dto.MenuGroupRequest, author string, upsertToAggrMenu bool) (string, error)
	UpsertDeliveryMenu(ctx context.Context, req dto.MenuGroupRequest, author string) (string, error)
	UpdateStopListStores(ctx context.Context, req dto.ProductStopList, author string) (dto.StoreProductStopLists, error)
	UpsertMenuByFields(ctx context.Context, fields models.UpdateFields, agg models.UpdateFieldsAggregators, req dto.MenuGroupRequest, author string) error

	UploadMenu(ctx context.Context, req dto.MenuUploadRequest) (string, error)
	CreateMenuUploadTransaction(ctx context.Context, req dto.MenuUploadTransaction) (string, error)
	VerifyUploadMenu(ctx context.Context, req dto.MenuUploadVerifyRequest) (dto.MenuUploadTransaction, error)
	GetProcessingMenuUploadTransactions(ctx context.Context, req dto.GetMenuUploadTransactions) ([]dto.MenuUploadTransaction, error)
	GetMenuUploadTransactions(ctx context.Context, req dto.GetMenuUploadTransactions) ([]models.MenuUploadTransaction, int64, error)

	ValidateMatching(ctx context.Context, req dto.MenuValidateRequest, sv3 *s3.S3) (string, error)
	ValidateVirtualStoreMatching(ctx context.Context, req dto.MenuValidateRequest) (string, error)
	ValidateAggProductErr(ctx context.Context, menuID string, storeID string, limit int) (aggregatorProduct []models.Product, posProduct []models.Product, total int, err error)

	RecoveryMenu(ctx context.Context, menuId, author string) error
	MergeMenus(ctx context.Context, restaurantID string, restaurantIDs []string, author string) (string, error)
	StopPositionsInVirtualStore(ctx context.Context, restaurantID, originalRestaurantID string) error
	RenewPositionsInVirtualStore(ctx context.Context, restaurantID, originalRestaurantID string) error

	DeleteProducts(ctx context.Context, req dto.DeleteProducts) error
	DeleteProductsFromDB(ctx context.Context, req dto.DeleteProducts) error
	UpdateMatchingProduct(ctx context.Context, req models.MatchingProducts) error

	AttributesStopList(ctx context.Context, storeId string, attributes []dto.StopListItem, author string) error
	UpdateMenu(ctx context.Context, req models.Menu) error
	UpdateMenuUploadTransaction(ctx context.Context, req dto.MenuUploadTransaction) error
	GetMenuUploadTransaction(ctx context.Context, req dto.MenuUploadTransaction) (dto.MenuUploadTransaction, error)
	GetAttributeForUpdate(ctx context.Context, query dto.AttributeSelector) (dto.AttributesUpdate, error)
	UpdateMenuName(ctx context.Context, query dto.UpdateMenuName) error
	DeleteAttributeGroupFromDB(ctx context.Context, req dto.DeleteAttributeGroup) error
	ValidateAttributeGroupName(ctx context.Context, menuId, name string) (bool, error)
	CreateAttributeGroup(ctx context.Context, menuID, attrGroupName string, min, max int) (string, error)
	AutoUploadMenuByPOS(ctx context.Context, req models.Menu) (string, error)

	AutoUpdateMenuPrices(ctx context.Context, storeId string) error
	AutoUpdateMenuDescriptions(ctx context.Context, storeId string) error

	PosIntegrationUpdateStopList(ctx context.Context, storeId string, request posIntegrationModels.StopListRequest, author string) error
	UpdateProductByFields(ctx context.Context, menuId string, productID string, req models.ProductUpdateRequest) error

	CreateGlovoSuperCollection(ctx context.Context, menuId string, superCollections dto.MenuSuperCollections) error
	InsertMenu(ctx context.Context, menu models.Menu) (string, error)
	CreateMenuByAggregatorApi(ctx context.Context, storeId string, aggregator string) (string, error)

	GetEmptyProducts(ctx context.Context, menuID string, page int64, limit int64) ([]models.Product, int, error)
	UpdateProductAvailableStatus(ctx context.Context, menuID, productID string, status bool) error

	AddRowToAttributeGroup(ctx context.Context, menuId string) error
}

var _ Client = &menuImpl{}

type menuImpl struct {
	menuManager                  managers.MenuManager
	storeManager                 managers.StoreManager
	storeCli                     storeCore.Client
	menuUploadTransactionManager managers.MenuUploadTransactionManager
}

func New(cfg dto.Config) (Client, error) {

	opts, err := menu.LoadConfig(context.Background(), cfg.SecretEnv, cfg.Region)
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

	stopListMan := managers.NewStopListTransactionManager(ds.StopListTransactionRepository())
	menuUploadTransactionManager := managers.NewMenuUploadTransactionManager(opts, ds.MenuUploadTransactionRepository(), validator.NewMenuUploadTransactionValidator())
	bkOffersManager := managers.NewBkOffersManager(opts, ds.BkOffersRepository())
	storeCli, err := storeCore.NewClient(storeDto.Config{})
	if err != nil {
		log.Trace().Err(err).Msgf("Creating store core error")
		return nil, err
	}

	sqsCli := notifyClient.NewSQS(sqs.NewFromConfig(opts.AwsConfig))

	mongoClient := ds.Client()
	db := mongoClient.Database(opts.DSDB)

	entityChangesHistoryRepo, err := entity_changes_history.NewEntityChangesHistoryMongoRepository(db)
	if err != nil {
		return nil, err
	}

	menuServiceRepo, err := menuServicePkg.NewMenuMongoRepository(db)
	if err != nil {
		return nil, err
	}

	return &menuImpl{
		storeCli:                     storeCli,
		menuUploadTransactionManager: menuUploadTransactionManager,
		storeManager:                 managers.NewStoreManager(opts, ds.StoreRepository(), ds.MenuRepository(entityChangesHistoryRepo), validator.NewStoreValidator()),
		menuManager:                  managers.NewMenuManager(opts, ds, ds.MenuRepository(entityChangesHistoryRepo), ds.StoreRepository(), ds.PromoRepository(), menuUploadTransactionManager, stopListMan, validator.NewMenuValidator(), sqsCli, ds.MSPositionsRepository(), ds.StopListTransactionRepository(), storeCli, bkOffersManager, *menuServiceRepo, ds.RestGroupMenuRepository()),
	}, nil
}

func (cli *menuImpl) PosIntegrationUpdateStopList(ctx context.Context, storeId string, request posIntegrationModels.StopListRequest, author string) error {
	if err := cli.menuManager.PosIntegrationUpdateStopList(ctx, storeId, request, entityChangesHistoryModels.EntityChangesHistoryRequest{
		Author:   author,
		TaskType: "pkg/menu/client.go - PosIntegrationUpdateStopList",
	}); err != nil {
		return err
	}

	return nil
}

func (cli *menuImpl) UpdateProductByFields(ctx context.Context, menuId string, productID string, req models.ProductUpdateRequest) error {
	if err := cli.menuManager.UpdateProductByFields(ctx, menuId, productID, req); err != nil {
		return err
	}

	return nil
}

func (cli *menuImpl) MergeMenus(ctx context.Context, restaurantID string, restaurantIDs []string, author string) (string, error) {
	id, err := cli.menuManager.MergeMenus(ctx, restaurantID, restaurantIDs, entityChangesHistoryModels.EntityChangesHistoryRequest{
		Author:   author,
		TaskType: "pkg/menu/client.go - MergeMenus",
	})
	if err != nil {
		log.Trace().Err(err).Msgf("merge menus error: id=%s", id)
		return "", err
	}

	return id, nil
}

func (cli *menuImpl) AutoUpdateMenuPrices(ctx context.Context, storeId string) error {
	err := cli.menuManager.AutoUpdateMenuPrices(ctx, storeId)
	if err != nil {
		log.Trace().Err(err).Msgf("auto-update menu prices error: %s", storeId)
		return err
	}

	return nil
}

func (cli *menuImpl) AutoUpdateMenuDescriptions(ctx context.Context, storeId string) error {
	err := cli.menuManager.AutoUpdateMenuDescriptions(ctx, storeId)
	if err != nil {
		log.Trace().Err(err).Msgf("auto-update menu descriptions error: %s", storeId)
		return err
	}

	return nil
}

func (cli *menuImpl) ValidateVirtualStoreMatching(ctx context.Context, req dto.MenuValidateRequest) (string, error) {
	message, err := cli.menuManager.ValidateVirtualStoreMenus(ctx, req.ID, req.StoreID)
	if err != nil {
		log.Trace().Err(err).Msgf("validate matching error: %s", req.StoreID)
		return "", err
	}

	return message, err
}

func (cli *menuImpl) ValidateMatching(ctx context.Context, req dto.MenuValidateRequest, sv3 *s3.S3) (string, error) {
	message, err := cli.menuManager.ValidatePosAndAggregator(ctx, req.ID, req.StoreID, sv3)
	if err != nil {
		log.Trace().Err(err).Msgf("validate matching error: %s", req.StoreID)
		return "", err
	}

	return message, err
}

func (cli *menuImpl) GetMenu(ctx context.Context, externalStoreID string, deliveryService dto.DeliveryService) (models.Menu, error) {

	store, err := cli.storeCli.FindStore(ctx, storeDto.StoreSelector{
		ExternalStoreID: externalStoreID,
		DeliveryService: deliveryService.String(),
	})
	if err != nil {
		return models.Menu{}, err
	}

	menuID, err := utils.ActiveMenu(store.Menus, deliveryService)
	if err != nil {
		return models.Menu{}, err
	}

	menu, err := cli.menuManager.GetMenu(ctx, store.MenuID, selector.EmptyMenuSearch().
		SetMenuID(menuID))
	if err != nil {
		return models.Menu{}, err
	}

	return menu, nil
}

func (cli *menuImpl) GetMenuByID(ctx context.Context, menuID string) (models.Menu, error) {
	menu, err := cli.menuManager.GetMenuByID(ctx,
		selector.EmptyMenuSearch().
			SetMenuID(menuID))
	if err != nil {
		return models.Menu{}, err
	}

	return menu, nil
}

func (cli *menuImpl) GetPromos(ctx context.Context, externalStoreID string, deliveryService dto.DeliveryService) (models.Promo, error) {

	store, err := cli.storeCli.FindStore(ctx, storeDto.StoreSelector{
		ExternalStoreID: externalStoreID,
		DeliveryService: deliveryService.String(),
	})
	if err != nil {
		return models.Promo{}, err
	}

	promoProducts, err := cli.menuManager.GetPromoProducts(ctx, selector.EmptyPromoSearch().
		SetStoreID(store.ID).SetDeliveryService(deliveryService.String()))
	if err != nil {
		return models.Promo{}, err
	}

	return promoProducts, nil
}

func (cli *menuImpl) GetStorePromos(ctx context.Context, query dto.GetPromosSelector) ([]dto.PromoDiscount, error) {
	promos, err := cli.menuManager.GetPromos(ctx, selector.EmptyPromoSearch().
		SetStoreID(query.StoreID).
		SetDeliveryService(query.DeliveryService).
		SetExternalStoreID(query.ExternalStoreID).
		SetProductIDs(query.ProductIDs).
		SetIsActive(query.IsActive))

	if err != nil {
		return nil, err
	}

	return dto.FromPromoDiscounts(promos), nil
}

func (cli *menuImpl) GetStores(ctx context.Context, deliveryService dto.DeliveryService) ([]storeModels.Store, error) {

	stores, err := cli.storeManager.GetStores(ctx, selector.EmptyStoreSearch().
		SetDeliveryService(deliveryService.String()))
	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (cli *menuImpl) GetPosDiscounts(ctx context.Context, storeID string, deliveryService dto.DeliveryService) (dto.PosDiscount, error) {
	promoProducts, err := cli.menuManager.GetPromoProducts(ctx,
		selector.EmptyPromoSearch().SetStoreID(storeID).SetDeliveryService(deliveryService.String()))
	if err != nil {
		return dto.PosDiscount{}, err
	}

	return dto.FromDiscounts(promoProducts), nil
}

func (cli *menuImpl) UpsertMenu(ctx context.Context, req dto.MenuGroupRequest, author string, upsertToAggrMenu bool) (string, error) {
	res, err := cli.menuManager.UpsertMenu(ctx, selector.EmptyMenuSearch().
		SetGroupID(req.GroupID).
		SetStoreID(req.StoreID).
		SetToken(req.Token), entityChangesHistoryModels.EntityChangesHistoryRequest{
		Author:   author,
		TaskType: "pkg/menu/client.go - UpsertMenu",
	}, upsertToAggrMenu)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (cli *menuImpl) UpsertMenuByFields(ctx context.Context, fields models.UpdateFields, agg models.UpdateFieldsAggregators, req dto.MenuGroupRequest, author string) error {
	err := cli.menuManager.UpsertMenuByFields(ctx, selector.EmptyMenuSearch().SetGroupID(req.GroupID).SetStoreID(req.StoreID).SetToken(req.Token),
		fields, agg, entityChangesHistoryModels.EntityChangesHistoryRequest{
			Author:   author,
			TaskType: "pkg/menu/client.go - UpsertMenuByFields",
		})
	if err != nil {
		log.Err(err).Msg("error pkg/menu/client.go - UpsertMenuByFields")
		return err
	}

	return nil
}

func (cli *menuImpl) UpsertDeliveryMenu(ctx context.Context, req dto.MenuGroupRequest, author string) (string, error) {

	res, err := cli.menuManager.UpsertMenuByGroupID(ctx, selector.EmptyMenuSearch().
		SetMenuID(req.MenuID).
		SetGroupID(req.GroupID).
		SetStoreID(req.StoreID), entityChangesHistoryModels.EntityChangesHistoryRequest{
		Author:   author,
		TaskType: "pkg/menu/client.go - UpsertDeliveryMenu",
	})
	if err != nil {
		return "", err
	}

	return res, nil
}

func (cli *menuImpl) GetMenuGroups(ctx context.Context, menuID string) ([]dto.MenuGroup, error) {

	res, err := cli.menuManager.GetMenuGroups(ctx, selector.EmptyMenuSearch().
		SetMenuID(menuID))
	if err != nil {
		return nil, err
	}

	return dto.FromMenuGroups(res), nil
}

func (cli *menuImpl) ListStoresByProduct(ctx context.Context, req dto.GetStoreByProductRequest) ([]storeModels.Store, int64, error) {

	res, total, err := cli.storeManager.ListStoresByProduct(ctx, selector.EmptyMenuSearch().
		SetProductExtID(req.ExtID).
		SetProductIsAvailable(req.IsAvailable))
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

func (cli *menuImpl) UpdateStopListStores(ctx context.Context, req dto.ProductStopList, author string) (dto.StoreProductStopLists, error) {

	transactionResults, err := cli.menuManager.UpdateStopListStores(ctx, models.UpdateStopListProduct{
		ProductID: req.ProductID,
		SetToStop: req.SetToStop,
		Data:      req.Data.ToModels(),
	}, entityChangesHistoryModels.EntityChangesHistoryRequest{
		TaskType: "pkg/menu/client.go - UpdateStopListStores",
		Author:   author,
	})
	if err != nil {
		return nil, err
	}

	return dto.FromStopListTransactions(transactionResults), nil
}

func (cli *menuImpl) GetMenuStatus(ctx context.Context, storeId string, isDeleted bool) ([]storeDto.StoreDsMenuDto, error) {
	res, err := cli.menuManager.GetMenuStatus(ctx, storeId, isDeleted)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (cli *menuImpl) UploadMenu(ctx context.Context, req dto.MenuUploadRequest) (string, error) {

	res, err := cli.menuManager.UploadMenu(ctx, req.StoreId, req.MenuId, req.DeliveryName, req.Sv3, req.UserRole, req.UserName)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (cli *menuImpl) CreateMenuUploadTransaction(ctx context.Context, req dto.MenuUploadTransaction) (string, error) {
	transactionID, err := cli.menuUploadTransactionManager.Create(ctx, dto.FromMenuUploadTransaction(req))

	if err != nil {
		log.Err(err).Msg("creation failed")
	}

	return transactionID, err
}

func (cli *menuImpl) VerifyUploadMenu(ctx context.Context, req dto.MenuUploadVerifyRequest) (dto.MenuUploadTransaction, error) {
	transaction, err := cli.menuUploadTransactionManager.Get(ctx, selector.MenuUploadTransactionSearch().SetID(req.TransactionId))
	if err != nil {
		return dto.ToMenuUploadTransaction(transaction), err
	}

	store, err := cli.storeManager.GetStore(ctx, selector.EmptyStoreSearch().SetID(transaction.StoreID))
	if err != nil {
		return dto.ToMenuUploadTransaction(transaction), err
	}

	transaction, err = cli.menuManager.VerifyUploadMenu(ctx, transaction, store)
	if err != nil {
		return dto.ToMenuUploadTransaction(transaction), err
	}

	return dto.ToMenuUploadTransaction(transaction), nil
}

func (cli *menuImpl) GetProcessingMenuUploadTransactions(ctx context.Context, req dto.GetMenuUploadTransactions) ([]dto.MenuUploadTransaction, error) {
	transactions, _, err := cli.menuUploadTransactionManager.List(ctx, selector.EmptyMenuUploadTransactionSearch().
		SetService(models.AggregatorName(req.DeliveryService)).
		SetStoreID(req.StoreId).
		SetCreatedFrom(time.Now().UTC().Add(-time.Hour)).
		SetStatus(models.PROCESSING.String()))

	if err != nil {
		return dto.ToMenuUploadTransactions(transactions), err
	}
	return dto.ToMenuUploadTransactions(transactions), nil
}

func (cli *menuImpl) GetMenuUploadTransactions(ctx context.Context, req dto.GetMenuUploadTransactions) ([]models.MenuUploadTransaction, int64, error) {
	menuUpTrans, total, err := cli.menuUploadTransactionManager.List(ctx, selector.EmptyMenuUploadTransactionSearch().
		SetStoreID(req.StoreId).
		SetSorting("created_at.value", -1).
		SetPage(req.Page).SetLimit(req.Limit))
	if err != nil {
		return nil, 0, err
	}

	return menuUpTrans, total, nil
}

func (cli *menuImpl) RecoveryMenu(ctx context.Context, menuId, author string) error {

	menu, err := cli.menuManager.GetMenuByID(ctx,
		selector.EmptyMenuSearch().
			SetMenuID(menuId))
	if err != nil {
		return err
	}

	return cli.menuManager.RecoveryMenu(ctx, menu, entityChangesHistoryModels.EntityChangesHistoryRequest{
		Author:   author,
		TaskType: "pkg/menu/client.go - RecoveryMenu",
	})
}

func (m *menuImpl) StopPositionsInVirtualStore(ctx context.Context, restaurantID, originalRestaurantID string) error {
	if err := m.menuManager.StopPositionsInVirtualStore(ctx, restaurantID, originalRestaurantID); err != nil {
		return err
	}

	return nil
}

func (m *menuImpl) RenewPositionsInVirtualStore(ctx context.Context, restaurantID, originalRestaurantID string) error {
	if err := m.menuManager.RenewPositionsInVirtualStore(ctx, restaurantID, originalRestaurantID); err != nil {
		return err
	}

	return nil
}

func (cli *menuImpl) AttributesStopList(ctx context.Context, storeId string, attributes []dto.StopListItem, author string) error {
	req := make([]models.ItemStopList, 0, len(attributes))

	for _, attribute := range attributes {
		req = append(req, attribute.ToModel())
	}

	return cli.menuManager.AttributesStopList(ctx, storeId, req, entityChangesHistoryModels.EntityChangesHistoryRequest{
		Author:   author,
		TaskType: "pkg/menu/client.go - AttributesStopList",
	})
}

func (m *menuImpl) DeleteProducts(ctx context.Context, req dto.DeleteProducts) error {
	return m.menuManager.DeleteProducts(ctx, req.MenuID, req.ProductIds)
}

func (m *menuImpl) DeleteProductsFromDB(ctx context.Context, req dto.DeleteProducts) error {
	return m.menuManager.DeleteProductsFromDB(ctx, req.MenuID, req.ProductIds)
}

func (cli *menuImpl) GetAttributeForUpdate(ctx context.Context, query dto.AttributeSelector) (dto.AttributesUpdate, error) {
	var res dto.AttributesUpdate

	attr, groups, total, err := (cli.menuManager.GetAttributesForUpdate(ctx, selector.EmptyMenuSearch().
		SetMenuID(query.MenuID).
		SetLimit(query.Limit).
		SetPage(query.Page)))

	if err != nil {
		return res, err
	}

	res.Attributes = attr
	res.AttributeGroups = groups
	res.Total = total
	return res, nil
}

func (m *menuImpl) UpdateMenu(ctx context.Context, menu models.Menu) error {
	return m.menuManager.UpdateMenu(ctx, menu)
}

func (m *menuImpl) AddRowToAttributeGroup(ctx context.Context, menuId string) error {
	return m.menuManager.AddRowToAttributeGroup(ctx, menuId)
}

func (m *menuImpl) UpdateMenuUploadTransaction(ctx context.Context, req dto.MenuUploadTransaction) error {
	return m.menuUploadTransactionManager.Update(ctx, dto.FromMenuUploadTransaction(req))
}

func (m *menuImpl) GetMenuUploadTransaction(ctx context.Context, req dto.MenuUploadTransaction) (dto.MenuUploadTransaction, error) {
	extTransactionID := ""
	for _, extTransaction := range req.ExtTransactions {
		if extTransaction.ID != "" {
			extTransactionID = extTransaction.ID
			break
		}
	}

	res, err := m.menuUploadTransactionManager.Get(ctx, selector.MenuUploadTransaction{
		ID:               req.ID,
		Service:          models.AggregatorName(req.Service),
		StoreID:          req.StoreID,
		ExtTransactionID: extTransactionID,
	})
	if err != nil {
		return dto.MenuUploadTransaction{}, err
	}

	return dto.ToMenuUploadTransaction(res), nil
}

func (m *menuImpl) UpdateMatchingProduct(ctx context.Context, req models.MatchingProducts) error {
	return m.menuManager.UpdateMatchingProduct(ctx, req)
}

func (m *menuImpl) UpdateMenuName(ctx context.Context, query dto.UpdateMenuName) error {
	return m.menuManager.UpdateMenuName(ctx, query.ToModel())
}

func (m *menuImpl) DeleteAttributeGroupFromDB(ctx context.Context, req dto.DeleteAttributeGroup) error {
	return m.menuManager.DeleteAttributeGroupFromDB(ctx, req.MenuId, req.AttributeGroupExtId)
}

func (m *menuImpl) ValidateAttributeGroupName(ctx context.Context, menuId, name string) (bool, error) {
	return m.menuManager.ValidateAttributeGroupName(ctx, menuId, name)
}

func (m *menuImpl) CreateAttributeGroup(ctx context.Context, menuID, attrGroupName string, min, max int) (string, error) {
	return m.menuManager.CreateAttributeGroup(ctx, menuID, attrGroupName, min, max)
}

func (m *menuImpl) AutoUploadMenuByPOS(ctx context.Context, req models.Menu) (string, error) {
	menuId, err := m.menuManager.AutoUploadMenuByPOS(ctx, req, selector.Menu{
		ID: req.ID,
	})
	if err != nil {
		return "", err
	}
	return menuId, nil
}

func (m *menuImpl) CreateGlovoSuperCollection(ctx context.Context, menuId string, superCollections dto.MenuSuperCollections) error {
	return m.menuManager.CreateGlovoSuperCollection(ctx, menuId, superCollections)
}

func (m *menuImpl) InsertMenu(ctx context.Context, menu models.Menu) (string, error) {
	return m.menuManager.InsertMenu(ctx, menu)
}

func (m *menuImpl) ValidateAggProductErr(ctx context.Context, menuID string, storeID string, limit int) (aggregatorProduct []models.Product, posProduct []models.Product, total int, err error) {
	aggregatorProduct, posProduct, total, err = m.menuManager.ValidateAggAndPosMatching(ctx, menuID, storeID, limit)
	if err != nil {
		return nil, nil, 0, err
	}
	return
}

func (m *menuImpl) GetEmptyProducts(ctx context.Context, menuID string, page int64, limit int64) ([]models.Product, int, error) {

	pagination := selector.Pagination{
		Page:  page,
		Limit: limit,
	}

	return m.menuManager.GetEmptyProducts(ctx, menuID, pagination)
}

func (m *menuImpl) UpdateProductAvailableStatus(ctx context.Context, menuID, productID string, status bool) error {
	return m.menuManager.UpdateProductAvailableStatus(ctx, menuID, productID, status)
}

func (m *menuImpl) CreateMenuByAggregatorApi(ctx context.Context, storeId string, aggregator string) (string, error) {
	return m.menuManager.CreateMenuByAggregatorAPI(ctx, aggregator, storeId)
}
