package managers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kwaaka-team/orders-core/config/menu"
	externalPosIntegrationModels "github.com/kwaaka-team/orders-core/core/integration_api/models"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator"
	aggregatorErrors "github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/errors"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/managers/validator"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
	pkgUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	entityChangesHistoryModels "github.com/kwaaka-team/orders-core/service/entity_changes_history/models"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/pkg/store"
	storeCoreModel "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/rs/zerolog/log"
)

// IS ? separate, menu, stoplist, product,etc ?
type MenuManager interface {
	GetMenu(ctx context.Context, posMenuID string, query selector.Menu) (models.Menu, error)
	GetMenuByID(ctx context.Context, query selector.Menu) (models.Menu, error)
	GetPosMenu(ctx context.Context, store storeModels.Store) (models.Menu, error)
	GetAggMenusFromPos(ctx context.Context, store storeModels.Store) ([]models.Menu, error)
	GetPromoProducts(ctx context.Context, query selector.Promo) (models.Promo, error)
	GetPromos(ctx context.Context, query selector.Promo) ([]models.Promo, error)
	GetMenuStatus(ctx context.Context, storeId string, isDeleted bool) ([]storeCoreModel.StoreDsMenuDto, error)
	GetAttributesForUpdate(ctx context.Context, query selector.Menu) (models.Attributes, models.AttributeGroups, int, error)

	UpsertMenu(ctx context.Context, query selector.Menu, history entityChangesHistoryModels.EntityChangesHistoryRequest, upsertToAggrMenu bool) (string, error)
	VerifyUploadMenu(ctx context.Context, transaction models.MenuUploadTransaction, store storeModels.Store) (models.MenuUploadTransaction, error)
	UpsertMenuByFields(ctx context.Context, query selector.Menu, fields models.UpdateFields, agg models.UpdateFieldsAggregators, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest) error

	UploadMenu(ctx context.Context, storeId, menuId, aggregatorName string, sv3 *s3.S3, userRole, userName string) (string, error)

	GetMenuGroups(ctx context.Context, query selector.Menu) (models.Groups, error)
	UpsertMenuByGroupID(ctx context.Context, query selector.Menu, history entityChangesHistoryModels.EntityChangesHistoryRequest) (string, error)

	UpdateStopListStores(ctx context.Context, req models.UpdateStopListProduct, history entityChangesHistoryModels.EntityChangesHistoryRequest) ([]models.StopListTransaction, error)
	UploadMenuGeneratedByPOS(ctx context.Context, storeID string, delivery string) error
	ValidatePosAndAggregator(ctx context.Context, menuID string, storeID string, sv3 *s3.S3) (string, error)
	ValidateVirtualStoreMenus(ctx context.Context, menuID string, storeID string) (string, error)
	ValidateAggAndPosMatching(ctx context.Context, menuID string, storeID string, limit int) (aggregatorProducts []models.Product, posProducts []models.Product, total int, err error)

	RecoveryMenu(ctx context.Context, req models.Menu, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest) error
	MergeMenus(ctx context.Context, restaurantID string, restaurantIDs []string, history entityChangesHistoryModels.EntityChangesHistoryRequest) (string, error)

	StopPositionsInVirtualStore(ctx context.Context, restaurantID, originalRestaurantID string) error
	RenewPositionsInVirtualStore(ctx context.Context, restaurantID, originalRestaurantID string) error

	DeleteProducts(ctx context.Context, menuId string, productsIds []string) error
	DeleteProductsFromDB(ctx context.Context, menuId string, productsIds []string) error
	UpdateMatchingProduct(ctx context.Context, req models.MatchingProducts) error

	AttributesStopList(ctx context.Context, storeId string, items []models.ItemStopList, history entityChangesHistoryModels.EntityChangesHistoryRequest) error
	UpdateMenu(ctx context.Context, menu models.Menu) error
	UpdateMenuName(ctx context.Context, query models.UpdateMenuName) error

	StopListSchedule(ctx context.Context, req models.StopListScheduler, history entityChangesHistoryModels.EntityChangesHistoryRequest) error

	DeleteAttributeGroupFromDB(ctx context.Context, menuId string, attrGroupExtId string) error
	ValidateAttributeGroupName(ctx context.Context, menuId, name string) (bool, error)
	CreateAttributeGroup(ctx context.Context, menuID, attrGroupName string, min, max int) (string, error)
	AutoUploadMenuByPOS(ctx context.Context, menu models.Menu, query selector.Menu) (string, error)

	AutoUpdateMenuPrices(ctx context.Context, storeId string) error
	AutoUpdateMenuDescriptions(ctx context.Context, storeId string) error

	PosIntegrationUpdateStopList(ctx context.Context, storeId string, request externalPosIntegrationModels.StopListRequest, history entityChangesHistoryModels.EntityChangesHistoryRequest) error

	UpdateProductByFields(ctx context.Context, menuId string, productID string, req models.ProductUpdateRequest) error
	CreateGlovoSuperCollection(ctx context.Context, menuId string, superCollections dto.MenuSuperCollections) error
	InsertMenu(ctx context.Context, menu models.Menu) (string, error)
	CreateMenuByAggregatorAPI(ctx context.Context, aggregator string, aggregatorStoreId string) (string, error)

	GetEmptyProducts(ctx context.Context, menuID string, pagination selector.Pagination) ([]models.Product, int, error)
	UpdateProductAvailableStatus(ctx context.Context, menuID, productID string, status bool) error

	AddRowToAttributeGroup(ctx context.Context, menuId string) error
}

type mnm struct {
	tx           drivers.TxStarter
	globalConfig menu.Configuration

	storeRepo         drivers.StoreRepository
	menuRepo          drivers.MenuRepository
	promoRepo         drivers.PromoRepository
	mspRepo           drivers.MSPositionsRepository
	stRepo            drivers.StopListTransactionRepository
	restGroupMenuRepo drivers.RestaurantGroupMenuRepository

	mutMan MenuUploadTransactionManager
	stm    StopListTransaction

	menuValidator   validator.Menu
	notifyCli       que.SQSInterface
	storeCli        store.Client
	bkOffersRepo    BkOffersManager
	menuServiceRepo menuServicePkg.MongoRepository
}

func NewMenuManager(
	globalConfig menu.Configuration,
	tx drivers.TxStarter,
	menuRepo drivers.MenuRepository,
	storeRepo drivers.StoreRepository,
	promoRepo drivers.PromoRepository,
	mutMan MenuUploadTransactionManager,
	stm StopListTransaction,
	menuValidator validator.Menu,
	notifyCli que.SQSInterface,
	mspRepo drivers.MSPositionsRepository,
	stRepo drivers.StopListTransactionRepository,
	storeCli store.Client,
	bkOffersRepo drivers.BkOffersRepository,
	menuServiceRepo menuServicePkg.MongoRepository,
	restGroupMenuRepo drivers.RestaurantGroupMenuRepository) MenuManager {

	return &mnm{
		tx:                tx,
		globalConfig:      globalConfig,
		menuRepo:          menuRepo,
		storeRepo:         storeRepo,
		promoRepo:         promoRepo,
		mutMan:            mutMan,
		stm:               stm,
		menuValidator:     menuValidator,
		notifyCli:         notifyCli,
		mspRepo:           mspRepo,
		stRepo:            stRepo,
		storeCli:          storeCli,
		bkOffersRepo:      bkOffersRepo,
		menuServiceRepo:   menuServiceRepo,
		restGroupMenuRepo: restGroupMenuRepo,
	}
}

func (m *mnm) UpdateProductByFields(ctx context.Context, menuId string, productID string, req models.ProductUpdateRequest) error {
	if err := m.menuRepo.UpdateProductByFields(ctx, menuId, productID, req); err != nil {
		return err
	}

	return nil
}

func (m *mnm) PosIntegrationUpdateStopList(ctx context.Context, storeId string, request externalPosIntegrationModels.StopListRequest, history entityChangesHistoryModels.EntityChangesHistoryRequest) error {
	store, err := m.storeRepo.Get(ctx, selector.Store{
		ID: storeId,
	})
	if err != nil {
		return err
	}

	if !store.ExternalPosIntegrationSettings.StopListIsOn {
		return fmt.Errorf("stoplist integration is off")
	}

	result := models.StopListTransaction{
		StoreID:   store.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	for _, menu := range store.Menus {
		if !menu.IsActive {
			continue
		}

		products, attributes, err := m.updateStopListAggregator(ctx, menu.ID, request.Positions.ToSliceOfString(), false, "", history)
		if err != nil {
			log.Trace().Err(err).Msgf("could not get/update menu %s", menu.ID)
			continue
		}

		if len(products) == 0 && len(attributes) == 0 {
			log.Trace().Msgf("no products && attributes changed")
			continue
		}

		trx, err := m.bulkUpdateAggregator(ctx, storeModels.AggregatorName(menu.Delivery), store, products.Unique(), attributes.Unique())
		if err != nil {
			log.Trace().Err(err).Msgf("could not bulk update store %s", store.ID)
			continue
		}

		if trx != nil {
			result.Transactions = append(result.Transactions, trx...)
		}

	}

	if result.Transactions == nil {
		log.Info().Msgf("something wrong, transactions is null")
		return nil
	}

	m.stm.Insert(context.Background(), []models.StopListTransaction{result})

	return nil
}

// todo interface segregation - StopList;
func (m *mnm) StopListSchedule(ctx context.Context, setting models.StopListScheduler, history entityChangesHistoryModels.EntityChangesHistoryRequest) error {

	if !setting.IsActive {
		return errors.New(fmt.Sprintf("not active rst for run scheduler by rst_id %v", setting.RstID))
	}

	rst, err := m.storeRepo.Get(ctx, selector.Store{
		ID: setting.RstID,
	})
	if err != nil {
		return err
	}

	setting.Available = setting.DefineAvailability()

	fmt.Printf("status available %v", setting.Available)

	var products = make([]models.Product, 0, len(setting.Products))

	for _, item := range setting.Products {
		products = append(products, models.Product{
			ExtID:       item.ID,
			IsAvailable: setting.Available,
		})
	}
	fmt.Println("stopListItems items", products)

	var transData []models.TransactionData
	for _, menu := range rst.Menus {

		if !menu.IsActive {
			continue
		}

		transData, err = m.updateProductsInAggregator(ctx,
			storeModels.AggregatorName(menu.Delivery),
			rst,
			products,
		)
		if err != nil {
			return err
		}

		_, err = m.updateAggregatorMenuProducts(ctx, menu.ID, products, setting.Available, history)
		if err != nil {
			return err
		}
	}

	result := models.StopListTransaction{
		StoreID:   setting.RstID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	//fmt.Println(updated, "updated ")

	if transData != nil {
		result.Transactions = append(result.Transactions, transData...)
	}
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		// save to DB
		m.stm.Insert(context.Background(), []models.StopListTransaction{result})
	}()

	go func() {
		defer wg.Done()
		// send to telegram
		m.sendToNotify(context.Background(), rst, result)
	}()
	wg.Wait()

	if result.Transactions == nil {
		log.Info().Msgf("something wrong, transactions is null")
		return models.ErrNoAnyStopListTransaction
	}

	return nil
}

func (m *mnm) appendMenus(newMenu models.Menu, mainMenu models.Menu, secondaryMenu models.Menu, restaurantID string) models.Menu {
	var (
		existProducts = make(map[string]models.Product, len(mainMenu.Products))
	)

	for _, product := range mainMenu.Products {
		existProducts[product.ExtID] = product
	}

	for _, product := range secondaryMenu.Products {
		product.ExtID = restaurantID + "_" + product.ExtID

		for index, defaultAttribute := range product.MenuDefaultAttributes {
			defaultAttribute.ExtID = restaurantID + "_" + defaultAttribute.ExtID
			product.MenuDefaultAttributes[index] = defaultAttribute
		}

		if val, ok := existProducts[product.ExtID]; ok {
			for _, defaultAttribute := range val.MenuDefaultAttributes {
				if defaultAttribute.ByAdmin {
					product.MenuDefaultAttributes = append(product.MenuDefaultAttributes, defaultAttribute)
				}
			}
		}

		for index, attributeID := range product.Attributes {
			product.Attributes[index] = restaurantID + "_" + attributeID
		}

		newMenu.Products = append(newMenu.Products, product)
	}

	for _, attribute := range secondaryMenu.Attributes {
		attribute.ExtID = restaurantID + "_" + attribute.ExtID
		newMenu.Attributes = append(newMenu.Attributes, attribute)
	}

	for _, attributeGroup := range secondaryMenu.AttributesGroups {
		for index, attributeID := range attributeGroup.Attributes {
			attributeGroup.Attributes[index] = restaurantID + "_" + attributeID
		}

		newMenu.AttributesGroups = append(newMenu.AttributesGroups, attributeGroup)
	}

	newMenu.Sections = append(newMenu.Sections, secondaryMenu.Sections...)

	newMenu.Collections = append(newMenu.Collections, secondaryMenu.Collections...)

	newMenu.SuperCollections = append(newMenu.SuperCollections, secondaryMenu.SuperCollections...)

	return newMenu
}

func (m *mnm) MergeMenus(ctx context.Context, restaurantID string, restaurantIDs []string, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest) (string, error) {
	mainStore, err := m.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{
		ID: restaurantID,
	})
	if err != nil {
		log.Info().Msgf("main restaurant with id %s not found", restaurantID)
		return "", err
	}

	mainMenu := models.Menu{
		Name: mainStore.Name + "menu",
		CreatedAt: coreModels.Time{
			Time: time.Now().UTC(),
		},
	}

	posMenu, err := m.GetMenuByID(ctx, selector.Menu{
		ID: mainStore.MenuID,
	})
	if err == nil {
		mainMenu = posMenu
	}

	newMenu := models.Menu{
		ID:        mainMenu.ID,
		Name:      mainMenu.Name,
		CreatedAt: mainMenu.CreatedAt,
		UpdatedAt: coreModels.Time{
			Time: time.Now().UTC(),
		},
	}

	for _, id := range restaurantIDs {
		secondaryStore, err := m.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{
			ID: id,
		})
		if err != nil {
			log.Info().Msgf("secondary restaurant with id %s not found", id)
			return "", err
		}

		secondaryMenu, err := m.GetMenuByID(ctx, selector.Menu{
			ID: secondaryStore.MenuID,
		})
		if err != nil {
			log.Info().Msgf("secondary pos menu with id %s not found for restaurant name %s", secondaryStore.MenuID, secondaryStore.Name)
			return "", err
		}

		newMenu = m.appendMenus(newMenu, mainMenu, secondaryMenu, secondaryStore.ID)
	}

	id, err := m.createOrUpdateMenu(ctx, selector.EmptyMenuSearch().SetMenuID(newMenu.ID), newMenu, mainStore, entityChangesHistoryRequest)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()

	if err := m.storeCli.Update(ctx, storeCoreModel.UpdateStore{
		ID:        &mainStore.ID,
		MenuID:    &id,
		UpdatedAt: &now,
	}); err != nil {
		return "", err
	}

	log.Info().Msgf("menu id = %s", id)

	return id, nil
}

func (m *mnm) ValidateVirtualStoreMenus(ctx context.Context, menuID string, storeID string) (string, error) {
	hasVirtualStore := true

	var msg string

	store, err := m.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{
		ID:              storeID,
		HasVirtualStore: &hasVirtualStore,
	})
	if err != nil {
		return "", err
	}

	if menuID == "" {
		// TODO: for all aggregators
		for _, menu := range store.Menus {
			if menu.IsActive && menu.Delivery == models.GLOVO.String() {
				menuID = menu.ID
				break
			}
		}
	}

	menu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuID))
	if err != nil {
		return "", err
	}

	var uniqueRestaurants = make(map[string]Entity)

	for _, product := range menu.Products {
		arr := strings.Split(product.ExtID, "_")
		if len(arr) > 0 {
			if _, err := primitive.ObjectIDFromHex(arr[0]); err != nil {
				msg += fmt.Sprintf("product with id %s is wrong, restaurant id is empty\n", product.ExtID)
				continue
			}

			if _, ok := uniqueRestaurants[arr[0]]; !ok {
				uniqueRestaurants[arr[0]] = Entity{}
			}
		}
	}

	for id := range uniqueRestaurants {
		restaurant, err := m.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{
			ID: id,
		})
		if err != nil {
			return "", err
		}

		posMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(restaurant.MenuID))
		if err != nil {
			return "", err
		}

		entity := Entity{
			Products:   make(map[string]models.Product),
			Attributes: make(map[string]models.Attribute),
		}

		for _, product := range posMenu.Products {
			entity.Products[product.ExtID] = product
		}

		for _, attribute := range posMenu.Attributes {
			entity.Attributes[attribute.ExtID] = attribute
		}

		uniqueRestaurants[id] = entity
	}

	for _, product := range menu.Products {
		arr := strings.Split(product.ExtID, "_")
		if len(arr) > 1 {
			entity, ok := uniqueRestaurants[arr[0]]
			if !ok {
				return "", fmt.Errorf("can not find restaurant %s", arr[0])
			}

			if _, exist := entity.Products[arr[1]]; exist {
				continue
			}

			msg += fmt.Sprintf("product with id %s doesn't exist in real pos menu\n", product.ExtID)
		}
	}

	for _, attribute := range menu.Products {
		arr := strings.Split(attribute.ExtID, "_")
		if len(arr) > 1 {
			entity, ok := uniqueRestaurants[arr[0]]
			if !ok {
				return "", fmt.Errorf("can not find restaurant %s", arr[0])
			}

			if _, exist := entity.Attributes[arr[1]]; exist {
				continue
			}

			if _, exist := entity.Products[arr[1]]; exist {
				continue
			}

			msg += fmt.Sprintf("attribute with id %s doesn't exist in real pos menu\n", attribute.ExtID)
		}
	}

	return msg, nil
}

type Entity struct {
	Products   map[string]models.Product
	Attributes map[string]models.Attribute
}

func (m *mnm) ValidatePosAndAggregator(ctx context.Context, menuID string, storeID string, sv3 *s3.S3) (string, error) {
	store, err := m.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{
		ID: storeID,
	})
	if err != nil {
		log.Trace().Err(err).Msgf("validate matching, find store error, id: %s", storeID)
		return "", err
	}

	//_, err = m.UpsertMenu(ctx, selector.EmptyMenuSearch().SetStoreID(store.ID))
	//if err != nil {
	//	log.Trace().Err(err).Msgf("upsert menu error, id: %s", storeID)
	//	return "", err
	//}

	posMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(store.MenuID))
	if err != nil {
		log.Trace().Err(err).Msgf("get POS menu error, id: %s", store.MenuID)
		return "", err
	}

	positions := make(map[string]models.ProductStatus)
	posAttributeGroupsMap := make(map[string]models.AttributeGroup)

	// write product status(del, includ, attrgroupIDs, defaults)
	for _, product := range posMenu.Products {
		var defaults []string

		for _, defaultAttribute := range product.MenuDefaultAttributes {
			if defaultAttribute.ByAdmin {
				defaults = append(defaults, defaultAttribute.ExtID)
			}
		}

		positions[product.ExtID] = models.ProductStatus{
			IsDeleted:         product.IsDeleted,
			IsIncludedInMenu:  product.IsIncludedInMenu,
			AttributeGroupIDs: product.AttributesGroups,
			Defaults:          defaults,
		}
	}

	for _, attribute := range posMenu.Attributes {
		positions[attribute.ExtID] = models.ProductStatus{
			Name:             attribute.Name,
			IsDeleted:        attribute.IsDeleted,
			IsIncludedInMenu: attribute.IncludedInMenu,
		}
	}

	for _, attributeGroup := range posMenu.AttributesGroups {
		posAttributeGroupsMap[attributeGroup.ExtID] = attributeGroup
	}

	body := fmt.Sprintf("<b>Название ресторана: %s\n\n", store.Name)

Loop:
	for _, menu := range store.Menus {
		if !menu.IsActive {
			continue
		}

		restMessage := body + fmt.Sprintf("Сервис доставки: %s\nНазвание меню: %s\nID меню: %s\n\n</b>", menu.Delivery, menu.Name, menu.ID)

		aggregatorMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menu.ID))
		if err != nil {
			continue
		}

		switch aggregatorMenu.Delivery {
		case "glovo":
			if !store.Glovo.SendToPos {
				continue
			}
		case "wolt":
			if !store.Wolt.SendToPos {
				continue
			}
		case "qr_menu":
			if !store.QRMenu.IsIntegrated {
				continue
			}
		case "chocofood":
			if !store.Chocofood.SendToPos {
				continue
			}
		case "yandex":
			for _, external := range store.ExternalConfig {
				if external.Type == "yandex" {
					if !external.SendToPos {
						continue Loop
					}
				}
			}
		}

		products, attributes, minMaxReports, err := m.validateAggregatorMenu(ctx, positions, aggregatorMenu, posAttributeGroupsMap, storeID)
		if err != nil {
			continue
		}

		report := models.ValidateReport{
			ID:                   store.ID,
			RestaurantName:       store.Name,
			Delivery:             menu.Delivery,
			MenuID:               menu.ID,
			Products:             products,
			AttributeGroupMinMax: minMaxReports,
		}

		for _, groupReport := range attributes {
			report.AttributeGroups = append(report.AttributeGroups, groupReport)
		}

		date := fmt.Sprintf("%v", time.Now().Format("2006.01.02"))

		if len(products) != 0 || len(attributes) != 0 || len(minMaxReports) != 0 {
			log.Info().Msgf("error for %s, msg: %s", store.Name, restMessage)

			//urlEscape := url.QueryEscape(path.Join(store.Name, menu.ID))

			var storeName = reduceSpaces(removeNonAlphaNumericSymbols(store.Name))
			link := strings.TrimSpace(fmt.Sprintf("s3://%v/validation/%s/%v", os.Getenv(models.S3_BUCKET), date, path.Join(storeName, menu.ID)))

			//imageLink := strings.TrimSpace(fmt.Sprintf("https://share-menu.kwaaka.com/validation/%s/%v.json", date, urlEscape))

			t, _ := json.Marshal(report)

			_, err := sv3.PutObject(&s3.PutObjectInput{
				Bucket:      aws.String(os.Getenv(models.S3_BUCKET)),
				Key:         aws.String(name(link) + ".json"),
				Body:        strings.NewReader(string(t)),
				ContentType: aws.String("application/json"),
			})
			if err != nil {
				log.Err(err).Msgf("s3 load error")
				continue
			}

			//restMessage += imageLink + "\n\n-------------------------------------------------\n\n"
			//if err := m.notifyCli.SendMessage("stoplist-telegram", restMessage, "-957120845"); err != nil {
			//	log.Err(err).Msgf("send to telegram error")
			//}

			continue
		}

		log.Info().Msgf("success for %s", store.Name)
	}

	return "", nil
}

func removeNonAlphaNumericSymbols(s string) string {
	var result string
	for _, char := range s {
		if isAlphanumeric(char) || char == ' ' || isCyrillic(char) {
			result += string(char)
		}
	}

	return result
}

func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

func isCyrillic(char rune) bool {
	return (char >= 'А' && char <= 'я') || char == 'Ё' || char == 'ё'
}

// Replace consecutive spaces with a single one
func reduceSpaces(s string) string {
	words := strings.Fields(s)
	result := strings.Join(words, " ")

	return result
}

func name(img string) string {
	raw := strings.TrimPrefix(string(img), "s3://")
	i := strings.Index(raw, "/")
	if i == -1 {
		return string(img)
	}
	return raw[i+1:]
}

func (m *mnm) validateAggregatorMenu(ctx context.Context, positions map[string]models.ProductStatus, aggregatorMenu models.Menu, posAttributeGroupsMap map[string]models.AttributeGroup, storeID string) ([]models.ProductReport, map[string]models.AttributeGroupReport, []models.ProductMinMaxReport, error) {
	var isSync = true
	var (
		products              = make([]models.ProductReport, 0)
		attributeGroupsReport = make(map[string]models.AttributeGroupReport)
		attributeGroups       = make(map[string]models.AttributeGroupStatus, len(aggregatorMenu.AttributesGroups))
		attributesMap         = make(map[string]models.AttributeReport)
		productMinMaxReport   = make([]models.ProductMinMaxReport, 0)
	)

	for index, attribute := range aggregatorMenu.Attributes {
		attributesMap[attribute.ExtID] = models.AttributeReport{
			ID:       attribute.ExtID,
			Name:     attribute.Name,
			Position: index,
		}
	}

	for index, attributeGroup := range aggregatorMenu.AttributesGroups {
		attributeGroups[attributeGroup.ExtID] = models.AttributeGroupStatus{
			Position:   index,
			Attributes: attributeGroup.Attributes,
			Name:       attributeGroup.Name,
			Min:        attributeGroup.Min,
			Max:        attributeGroup.Max,
		}
	}

	for index, aggregatorProduct := range aggregatorMenu.Products {
		if aggregatorProduct.IsDeleted {
			continue
		}

		// logic for check min and max

		existingAttributes := make(map[string]struct{})
		var doubleAttributeInProduct bool

		for _, attributeGroupID := range aggregatorProduct.AttributesGroups {
			if aggregatorAttributeGroup, ok := attributeGroups[attributeGroupID]; ok {
				for _, attribute := range aggregatorAttributeGroup.Attributes {
					existingAttributes[attribute] = struct{}{}
				}
			}
		}

		if posProduct, ok := positions[aggregatorProduct.ExtID]; ok {
			for _, defaultAttribute := range posProduct.Defaults {
				// if product has double attribute
				if _, ok := existingAttributes[defaultAttribute]; ok {
					doubleAttributeInProduct = true
				}
				existingAttributes[defaultAttribute] = struct{}{}
			}

			var attributeGroupsMinMaxReport []models.AttributeMinMaxReport

			for _, attributeGroupID := range posProduct.AttributeGroupIDs {
				if group, exist := posAttributeGroupsMap[attributeGroupID]; exist {
					min := 0
					posAttributesName := make([]string, 0, len(group.Attributes))
					aggregatorAttributesName := make([]string, 0)

					if group.Min > 0 {
						for _, posAttributeID := range group.Attributes {
							if _, exist := existingAttributes[posAttributeID]; exist {
								min++
								if val, ok := attributesMap[posAttributeID]; ok {
									aggregatorAttributesName = append(aggregatorAttributesName, val.Name)
								}
							}

							if val, ok := positions[posAttributeID]; ok {
								posAttributesName = append(posAttributesName, val.Name)
							}
						}

						if min < group.Min {
							attrGroupsMinMaxReport := models.AttributeMinMaxReport{
								ID:                      group.ExtID,
								Name:                    group.Name,
								Min:                     group.Min,
								Max:                     group.Max,
								PosAttributeName:        posAttributesName,
								AggregatorAttributeName: aggregatorAttributesName,
								CurrentMin:              min,
							}

							attributeGroupsMinMaxReport = append(attributeGroupsMinMaxReport, attrGroupsMinMaxReport)
						}
					}
				}
			}

			if len(attributeGroupsMinMaxReport) != 0 {
				prodMinMaxReport := models.ProductMinMaxReport{
					ID:                     aggregatorProduct.ExtID,
					Position:               index,
					AttributeMinMaxReports: attributeGroupsMinMaxReport,
				}

				if len(aggregatorProduct.Name) != 0 {
					prodMinMaxReport.Name = aggregatorProduct.Name[0].Value
				}

				log.Info().Msgf("has bad restrictions of min and max in attribute groups")

				productMinMaxReport = append(productMinMaxReport, prodMinMaxReport)
			}
		}

		// aggregator attribute groups in product
		for _, groupID := range aggregatorProduct.AttributesGroups {

			var flag bool

			var attributesReport = make([]models.AttributeReport, 0, 2)

			if gr, hasMainReport := attributeGroupsReport[groupID]; hasMainReport {
				productReport := models.ProductReport{
					ID:       aggregatorProduct.ExtID,
					Position: index,
				}

				if len(aggregatorProduct.Name) != 0 {
					productReport.Name = aggregatorProduct.Name[0].Value
				}

				gr.Products = append(gr.Products, productReport)
				attributeGroupsReport[groupID] = gr
				continue
			}

			// aggregator attribute groups array
			status, ok := attributeGroups[groupID]
			if ok {

				// aggregator attributes array in attribute group
				for _, attributeID := range status.Attributes {

					// if attribute not exist in POS
					if _, exist := positions[attributeID]; !exist {
						attribute, hasAttribute := attributesMap[attributeID]
						if !hasAttribute {
							continue
						}

						flag = true
						attributesReport = append(attributesReport, attribute)
					}
				}
			}

			if flag {

				_, hasReport := attributeGroupsReport[groupID]
				if !hasReport {
					report := models.AttributeGroupReport{
						ID:         groupID,
						Name:       status.Name,
						Min:        status.Min,
						Max:        status.Max,
						Position:   status.Position,
						Status:     "Эти атрибуты в атрибут группе не существуют в POS меню, в in_products указаны продукты с неправильной атрибут группой",
						Attributes: attributesReport,
						Products: []models.ProductReport{
							{
								ID:       aggregatorProduct.ExtID,
								Position: index,
							},
						},
					}

					if len(aggregatorProduct.Name) != 0 {
						report.Name = aggregatorProduct.Name[0].Value
					}

					attributeGroupsReport[groupID] = report
					continue
				}

				//productReport := models.ProductReport{
				//	ID:       aggregatorProduct.ExtID,
				//	Position: index,
				//}
				//
				//if len(aggregatorProduct.Name) != 0 {
				//	productReport.Name = aggregatorProduct.Name[0].Value
				//}
				//
				//group.Products = append(group.Products, productReport)
				//
				//attributeGroupsReport[groupID] = group
			}
		}

		id := aggregatorProduct.ExtID
		if aggregatorProduct.PosID != "" {
			id = aggregatorProduct.PosID
		}

		status, ok := positions[id]
		if !ok {
			isSync = false
			aggregatorMenu.Products[index].IsSync = false
			cur := models.ProductReport{
				ID:       id,
				Position: index,
				Status:   "pos product doesn't exist in POS",
				Solution: []string{"1. Зайти в админку", "2. Изменить айди продукта на другой", "3. Если данного продукта все еще нету для метчинга, удалить продукт в админке", "4. Обратиться к ресторану и передать информацию о том, что продукта не существует в выгрузке айко"},
			}

			if len(aggregatorProduct.Name) != 0 {
				cur.Name = aggregatorProduct.Name[0].Value
			}

			products = append(products, cur)
			continue
		}

		if status.IsDeleted {
			isSync = false
			aggregatorMenu.Products[index].IsSync = false
			cur := models.ProductReport{
				ID:       id,
				Position: index,
				Status:   "pos product is deleted in POS",
				Solution: []string{"1. Зайти в админку", "2. Изменить айди продукта на другой", "3. Если данного продукта все еще нету для метчинга, удалить продукт в админке", "4. Обратиться к ресторану и передать информацию о том, что продукт удален в выгрузке айко"},
			}

			if len(aggregatorProduct.Name) != 0 {
				cur.Name = aggregatorProduct.Name[0].Value
			}

			products = append(products, cur)
			continue
		}

		if !status.IsIncludedInMenu {
			isSync = false
			aggregatorMenu.Products[index].IsSync = false
			cur := models.ProductReport{
				ID:       id,
				Position: index,
				Status:   "pos product included in menu false in POS",
				Solution: []string{"1. Зайти в админку", "2. Изменить айди продукта на другой", "3. Если данного продукта все еще нету для метчинга, удалить продукт в админке", "4. Обратиться к ресторану и передать информацию о том, что продукт не включен в меню в выгрузке айко"},
			}

			if len(aggregatorProduct.Name) != 0 {
				cur.Name = aggregatorProduct.Name[0].Value
			}

			products = append(products, cur)
			continue
		}

		if !aggregatorMenu.Products[index].IsSync {
			aggregatorMenu.Products[index].IsSync = true
		}

		if doubleAttributeInProduct && !m.isIgnoreStoreList(storeID) {
			cur := models.ProductReport{
				ID:       aggregatorProduct.ExtID,
				Position: index,
				Status:   "product has a duplicate default attribute in the attribute group",
				Solution: []string{"1. Зайти в админку", "2. Найти продукт по имени или айди", "3. Удалить либо скрытый аттрибут, либо аттрибут группу которая дублирует этот скрытый аттрибут"},
			}
			if len(aggregatorProduct.Name) != 0 {
				cur.Name = aggregatorProduct.Name[0].Value
			}
			products = append(products, cur)
		}

	}

	//for index, aggregatorAttribute := range aggregatorMenu.Attributes {
	//	if aggregatorAttribute.IsDeleted {
	//		continue
	//	}
	//
	//	id := aggregatorAttribute.ExtID
	//	if aggregatorAttribute.PosID != "" {
	//		id = aggregatorAttribute.PosID
	//	}
	//
	//	status, ok := positions[id]
	//	if !ok {
	//		isSync = false
	//		aggregatorMenu.Attributes[index].IsSync = false
	//		attributes = append(attributes, models.AttributeReport{
	//			ID:       aggregatorAttribute.ExtID,
	//			Name:     aggregatorAttribute.Name,
	//			Position: index,
	//			Status:   "pos attribute doesn't exist in POS",
	//			Solution: []string{"1. Если атрибут привязан к какой-то атрибут группе, необходимо пересоздать атрибут группу у соответствующих продуктов"},
	//		})
	//		continue
	//	}
	//
	//	if status.IsDeleted {
	//		isSync = false
	//		aggregatorMenu.Attributes[index].IsSync = false
	//		attributes = append(attributes, models.AttributeReport{
	//			ID:       aggregatorAttribute.ExtID,
	//			Name:     aggregatorAttribute.Name,
	//			Position: index,
	//			Status:   "pos attribute is deleted in POS",
	//			Solution: []string{"1. Если атрибут привязан к какой-то атрибут группе, необходимо пересоздать атрибут группу у соответствующих продуктов"},
	//		})
	//		continue
	//	}
	//
	//	if !aggregatorMenu.Attributes[index].IsSync {
	//		aggregatorMenu.Attributes[index].IsSync = true
	//	}
	//}

	if !isSync {
		aggregatorMenu.IsSync = false
	}

	if err := m.menuRepo.Update(ctx, aggregatorMenu, entityChangesHistoryModels.EntityChangesHistory{
		CallFunction: "validateAggregatorMenu",
		Author:       "cron",
	}); err != nil {
		return products, attributeGroupsReport, productMinMaxReport, err
	}

	return products, attributeGroupsReport, productMinMaxReport, nil
}

func (m *mnm) UploadMenuGeneratedByPOS(ctx context.Context, storeID string, delivery string) error {
	store, err := m.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{
		ID: storeID,
	})
	if err != nil {
		return err
	}

	menu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(store.MenuID))
	if err != nil {
		return err
	}

	aggregatorMenu := generateAggregatorMenu(menu, delivery)

	id, err := m.menuRepo.Insert(ctx, aggregatorMenu)
	if err != nil {
		return err
	}

	var menus = make([]storeCoreModel.UpdateStoreDSMenu, 0, len(store.Menus)+1)

	for _, menu := range store.Menus {
		cur := menu

		menus = append(menus, storeCoreModel.UpdateStoreDSMenu{
			MenuID:    &cur.ID,
			Name:      &cur.Name,
			IsActive:  &cur.IsActive,
			IsDeleted: &cur.IsDeleted,
			Delivery:  &cur.Delivery,
			UpdatedAt: &cur.UpdatedAt,
		})
	}

	menus = append(menus, storeCoreModel.UpdateStoreDSMenu{
		MenuID:    &id,
		Name:      &aggregatorMenu.Name,
		IsActive:  &aggregatorMenu.IsActive,
		IsDeleted: &aggregatorMenu.IsDeleted,
		Delivery:  &aggregatorMenu.Delivery,
		UpdatedAt: &aggregatorMenu.UpdatedAt.Time,
	})

	if err := m.storeCli.Update(ctx, storeCoreModel.UpdateStore{
		ID:    &store.ID,
		Menus: menus,
	}); err != nil {
		return err
	}

	return nil
}

func generateAggregatorMenu(menu models.Menu, delivery string) models.Menu {
	var products = make([]models.Product, 0, len(menu.Products))
	var attributes = make([]models.Attribute, 0, len(menu.Attributes))
	var sections = make([]models.Section, 0, len(menu.Sections))

	for index, group := range menu.Groups {
		sections = append(sections, models.Section{
			ExtID:        group.ID,
			Name:         group.Name,
			SectionOrder: index + 1,
			Collection:   "1",
		})
	}

	for _, product := range menu.Products {
		cur := models.Product{
			ExtID:            uuid.New().String(),
			PosID:            product.ExtID,
			Name:             product.Name,
			Price:            product.Price,
			Description:      product.Description,
			Section:          product.ParentGroupID,
			Attributes:       product.Attributes,
			IsAvailable:      product.IsAvailable,
			AttributesGroups: product.AttributesGroups,
		}

		if product.Price[0].Value == 0 {
			product.Price[0].Value = 1
		}

		products = append(products, cur)
	}

	for _, attribute := range menu.Attributes {
		cur := models.Attribute{
			ExtID:       attribute.ExtID,
			PosID:       attribute.PosID,
			Name:        attribute.Name,
			Price:       attribute.Price,
			IsAvailable: attribute.IsAvailable,
		}

		attributes = append(attributes, cur)
	}

	collection := models.MenuCollection{
		ExtID:           "1",
		Name:            "Menu",
		CollectionOrder: 1,
	}

	return models.Menu{
		Name:             menu.Name + " " + delivery,
		Delivery:         delivery,
		IsActive:         true,
		Products:         products,
		Attributes:       attributes,
		Sections:         sections,
		Collections:      []models.MenuCollection{collection},
		AttributesGroups: menu.AttributesGroups,
		UpdatedAt:        coreModels.TimeNow(),
	}
}

func (m *mnm) GetMenu(ctx context.Context, posMenuID string, query selector.Menu) (models.Menu, error) {

	if err := m.menuValidator.ValidateID(ctx, query); err != nil {
		return models.Menu{}, err
	}

	menu, err := m.menuRepo.Get(ctx, query)
	if err != nil {
		return models.Menu{}, err
	}

	return menu, nil
}

func (m *mnm) GetMenuByID(ctx context.Context, query selector.Menu) (models.Menu, error) {

	if err := m.menuValidator.ValidateID(ctx, query); err != nil {
		return models.Menu{}, err
	}

	menu, err := m.menuRepo.Get(ctx, query)
	if err != nil {
		return models.Menu{}, err
	}

	return menu, nil
}

func (m *mnm) GetPromoProducts(ctx context.Context, query selector.Promo) (models.Promo, error) {

	promo, err := m.promoRepo.GetPromos(ctx, query)
	if err != nil {
		return models.Promo{}, err
	}

	return promo, nil
}

func (m *mnm) GetPromos(ctx context.Context, query selector.Promo) ([]models.Promo, error) {

	promos, err := m.promoRepo.FindPromos(ctx, query)
	if err != nil {
		return nil, err
	}

	return promos, nil
}

// GetPosMenu get menu from pos terminal and setting to kwaaka struct.
func (m *mnm) GetPosMenu(ctx context.Context, store storeModels.Store) (models.Menu, error) {

	posManager, err := pos.NewPosManager(m.globalConfig, m.menuRepo, store)
	if err != nil {
		return models.Menu{}, err
	}

	// send hash map with products id + parent_id to get menu method
	menu, err := posManager.GetMenu(ctx, store)
	if err != nil {
		return models.Menu{}, err
	}
	return menu, nil
}

func (m *mnm) GetAggMenusFromPos(ctx context.Context, store storeModels.Store) ([]models.Menu, error) {
	posManager, err := pos.NewPosManager(m.globalConfig, m.menuRepo, store)
	if err != nil {
		return nil, err
	}

	menus, err := posManager.GetAggMenu(ctx, store)
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (m *mnm) GetMenuGroups(ctx context.Context, query selector.Menu) (models.Groups, error) {
	return m.menuRepo.GetGroups(ctx, query)
}

// GetMenuStatus allows to get upload status. Checking in DB menu upload transactions collection,
// if not exist validate menu
func (m *mnm) GetMenuStatus(ctx context.Context, storeId string, isDeleted bool) ([]storeCoreModel.StoreDsMenuDto, error) {
	store, err := m.storeRepo.Get(ctx, selector.EmptyStoreSearch().SetID(storeId))
	if err != nil {
		return nil, err
	}

	if isDeleted {
		menus, err := m.getDeletedMenusStatus(ctx, store)
		if err != nil {
			return nil, err
		}
		return menus, nil
	}

	for i := range store.Menus {

		if store.Menus[i].HasWoltPromo {
			store.Menus[i].Status = models.HAS_PROMO.String()
			continue
		}

		trx, err := m.mutMan.Get(ctx, selector.EmptyMenuUploadTransactionSearch().
			SetStoreID(storeId).
			SetService(models.AggregatorName(store.Menus[i].Delivery)).
			SetSorting("created_at.value", -1))

		if err != nil {
			log.Err(err).Msgf("menu upload transaction err: menu_id %s", store.Menus[i].ID)
			store.Menus[i].Status = m.validateMenuToUpload(ctx, store.Menus[i].ID).String()
			continue
		}

		store.Menus[i].Status = trx.Status

		tr, err := trx.ExtTransactions.GetByMenu(store.Menus[i].ID)
		if err != nil {
			log.Err(err).Msgf("menu upload transaction: get menu_id %s", store.Menus[i].ID)
			store.Menus[i].Status = setStatus(trx.Status).String()
			continue
		}

		// check in aggregator
		if tr.Status == models.PROCESSING.String() {

			status, err := m.verifyMenu(ctx, store, storeModels.AggregatorName(trx.Service), tr)
			if err != nil {
				log.Err(err).Msgf("verify aggregator menu %s", store.Menus[i].ID)
			}

			if status != "" {
				tr.Status = status.String()
			}

		}

		store.Menus[i].Status = setStatus(trx.Status).String()

		store.Menus[i].UpdatedAt = trx.UpdatedAt.Value.Time

		store.Menus[i].EmptyProductPercentage, err = m.getEmptyProductsPercentage(ctx, store.Menus[i].ID)
		if err != nil {
			log.Err(err).Msgf("empty product per cent %s", store.Menus[i].ID)
		}
	}

	result := make([]storeCoreModel.StoreDsMenuDto, 0, len(store.Menus)+1)
	for i := range store.Menus {
		item := store.Menus[i]
		result = append(result, storeCoreModel.FromModel(item, false))
	}

	result = append(result, storeCoreModel.FromModel(storeModels.StoreDSMenu{
		ID:       store.MenuID,
		Name:     models.MAIN.String(),
		IsActive: true,
		Status:   models.NOT_PROCESSED.String(),
		Delivery: models.MAIN.String(),
	}, true))

	return result, err
}

func (m *mnm) getDeletedMenusStatus(ctx context.Context, store storeModels.Store) ([]storeCoreModel.StoreDsMenuDto, error) {
	var menus = make([]storeModels.StoreDSMenu, 0, len(store.Menus))

	for _, menu := range store.Menus {
		if menu.IsDeleted {
			menus = append(menus, menu)
		}
	}

	for i := range menus {
		if !menus[i].IsDeleted {
			continue
		}

		trx, err := m.mutMan.Get(ctx, selector.EmptyMenuUploadTransactionSearch().
			SetStoreID(store.ID).
			SetService(models.AggregatorName(menus[i].Delivery)).
			SetMenuID(menus[i].ID).
			SetSorting("created_at.value", -1))
		if err != nil {
			log.Err(err).Msgf("menu upload transaction err: menu_id %s", menus[i].ID)
			menus[i].Status = m.validateMenuToUpload(ctx, menus[i].ID).String()
			continue
		}

		menus[i].Status = trx.Status

		_, err = trx.ExtTransactions.GetByMenu(menus[i].ID)
		if err != nil {
			log.Err(err).Msgf("menu upload transaction: get menu_id %s", menus[i].ID)
			menus[i].Status = setStatus(trx.Status).String()
			continue
		}

		store.Menus[i].Status = setStatus(trx.Status).String()
		menus[i].UpdatedAt = trx.UpdatedAt.Value.Time
	}

	result := make([]storeCoreModel.StoreDsMenuDto, 0, len(store.Menus)+1)
	for i := range menus {
		item := menus[i]
		result = append(result, storeCoreModel.FromModel(item, false))
	}

	result = append(result, storeCoreModel.FromModel(storeModels.StoreDSMenu{
		ID:       store.MenuID,
		Name:     models.MAIN.String(),
		IsActive: true,
		Status:   models.NOT_PROCESSED.String(),
		Delivery: models.MAIN.String(),
	}, true))

	return result, nil
}

func (m *mnm) getEmptyProductsPercentage(ctx context.Context, menuID string) (int, error) {

	emptyProducts, _, err := m.GetEmptyProducts(ctx, menuID, selector.Pagination{Page: 0, Limit: 10})
	if err != nil {
		return 0, err
	}
	if len(emptyProducts) == 0 {
		return 100, nil
	}
	menu, err := m.GetMenuByID(ctx, selector.EmptyMenuSearch().SetMenuID(menuID))
	if err != nil {
		return 0, err
	}

	percentageEmptyProducts := ((len(menu.Products) - len(emptyProducts)) * 100) / len(menu.Products)

	return percentageEmptyProducts, nil
}

func (m *mnm) UpsertMenuByGroupID(ctx context.Context, query selector.Menu, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest) (string, error) {

	tx, cb, err := m.tx.StartSession(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		err = cb(err)
	}()

	store, err := m.storeRepo.Get(tx, selector.EmptyStoreSearch().
		SetID(query.StoreID))
	if err != nil {
		// TODO
		return "", err
	}

	posMenu, err := m.GetPosMenu(tx, store)
	if err != nil {
		// TODO
		return "", err
	}

	menuID, err := m.createOrUpdateMenu(tx, query, posMenu, store, entityChangesHistoryRequest)
	if err != nil {
		return "", err
	}

	if query.HasMenuID() {
		return menuID, err
	}

	store.Menus = append(store.Menus, storeModels.StoreDSMenu{
		ID:        menuID,
		Name:      store.Name,
		UpdatedAt: coreModels.TimeNow().Time,
	})

	// update stores
	if err = m.storeRepo.Update(tx, store); err != nil {
		return "", err
	}

	return menuID, err

}

func (m *mnm) getPosMenuProductPrices(posMenu models.Menu) map[string]float64 {
	posMenuProductsPrices := map[string]float64{}

	for _, product := range posMenu.Products {
		if len(product.Price) < 1 {
			continue
		}

		posMenuProductsPrices[product.ExtID] = product.Price[0].Value
	}

	return posMenuProductsPrices
}

type ChangedProductsInfo struct {
	Name     string  `json:"name"`
	OldPrice float64 `json:"old_price"`
	NewPrice float64 `json:"new_price"`
}

func (m *mnm) changeAggregatorProductsPrices(products []models.Product, posMenuProductsPrices map[string]float64, changedProductsInfo map[string]ChangedProductsInfo, productIdsMap map[string]bool) []models.Product {
	for i := range products {
		if products[i].DiscountPrice.Value != 0 {
			continue
		}

		id := products[i].ExtID
		if products[i].PosID != "" {
			id = products[i].PosID
		}

		if _, ok := productIdsMap[id]; ok {
			continue
		}

		price, ok := posMenuProductsPrices[id]
		if !ok {
			continue
		}

		if price == 0 {
			continue
		}

		for j := range products[i].Price {
			if products[i].Price[j].Value != price {
				var productName string

				if len(products[i].Name) > 0 {
					productName = products[i].Name[0].Value
				}

				changedProductsInfo[id] = ChangedProductsInfo{
					Name:     productName,
					OldPrice: products[i].Price[j].Value,
					NewPrice: price,
				}

				products[i].Price[j].Value = price
			}
		}
	}

	return products
}

func (m *mnm) prepareMessageToNotify(changedProductsInfo map[string]ChangedProductsInfo, store storeModels.Store) string {
	message := fmt.Sprintf("Автообновление цен для %s\n", store.Name)

	i := 1

	for _, value := range changedProductsInfo {
		message += fmt.Sprintf("%d. Название продукта: %s, старая цена: %v, новая цена: %v\n", i, value.Name, value.OldPrice, value.NewPrice)
		i++
	}

	message += "\n"

	return message
}

func (m *mnm) updatePricesInAggregatorMenus(ctx context.Context, store storeModels.Store, posMenu models.Menu, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest) error {
	if store.Settings.PriceSource != coreModels.POSPriceSource {
		return nil
	}

	posMenuProductsPrices := m.getPosMenuProductPrices(posMenu)

	changedProductsInfo := make(map[string]ChangedProductsInfo, 0)

	for _, menuDs := range store.Menus {
		if !menuDs.IsActive {
			continue
		}

		aggregatorMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuDs.ID))
		if err != nil {
			log.Err(err).Msgf("get aggregator menu by id %s", menuDs.ID)
			return err
		}

		promos, err := m.promoRepo.FindPromos(ctx, selector.EmptyPromoSearch().SetStoreID(store.ID).SetDeliveryService(menuDs.Delivery))
		if err != nil {
			log.Err(err).Msgf("get promos error: %v", err)
			return err
		}

		productIdsMap := make(map[string]bool)

		for _, promo := range promos {
			for _, id := range promo.ProductIds {
				productIdsMap[id] = true
			}
		}

		updatedProductsWithPrice := m.changeAggregatorProductsPrices(aggregatorMenu.Products, posMenuProductsPrices, changedProductsInfo, productIdsMap)

		aggregatorMenu.Products = updatedProductsWithPrice

		if err = m.menuRepo.Update(ctx, aggregatorMenu, entityChangesHistoryModels.EntityChangesHistory{
			CallFunction: "updatePricesInAggregatorMenus",
			Author:       entityChangesHistoryRequest.Author,
			TaskType:     entityChangesHistoryRequest.TaskType,
		}); err != nil {
			log.Err(err).Msgf("update aggregator menu by id %s", menuDs.ID)
			return err
		}
	}

	message := m.prepareMessageToNotify(changedProductsInfo, store)

	if err := m.notifyCli.SendMessage(m.globalConfig.Telegram, message, m.globalConfig.NotificationConfiguration.AutoUpdatePriceTelegramChatId, store.Telegram.TelegramBotToken); err != nil {
		return err
	}

	return nil
}

func (m *mnm) UpsertMenu(ctx context.Context, query selector.Menu, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest, upsertToAggrMenu bool) (string, error) {

	tx, cb, err := m.tx.StartSession(ctx)
	if err != nil {
		return "", err
	}
	defer func() {
		err = cb(err)
	}()

	store, err := m.storeRepo.Get(tx, selector.EmptyStoreSearch().
		SetID(query.StoreID).
		SetToken(query.Token))
	if err != nil {
		// TODO
		return "", err
	}

	if store.ValidationSettings.ForbiddenUpsert {
		return "", nil
	}

	if store.RestaurantGroupID == "646258ad1db3ef4dcf23c174" {
		menus, err := m.GetAggMenusFromPos(ctx, store)
		if err != nil {
			return "", err
		}

		var updateMenus []storeCoreModel.UpdateStoreDSMenu
		var yandexMenu *models.Menu

		// Look for a GLOVO menu and create a YANDEX menu from it
		for _, menu := range menus {
			if menu.Delivery == models.GLOVO.String() {
				duplicate := menu
				duplicate.Delivery = models.YANDEX.String()
				yandexMenu = &duplicate
				break
			}
		}

		// If a YANDEX menu was created, append it to the menus
		if yandexMenu != nil {
			menus = append(menus, *yandexMenu)
		}

		// Process all menus, inserting and preparing update structures
		for _, menu := range menus {
			cur := menu
			id, err := m.InsertMenu(ctx, menu)
			if err != nil {
				return "", err
			}

			cur.ID = id

			updateMenus = append(updateMenus, storeCoreModel.UpdateStoreDSMenu{
				MenuID:    &cur.ID,
				Name:      &cur.Name,
				IsActive:  &cur.IsActive,
				IsDeleted: &cur.IsDeleted,
				Delivery:  &cur.Delivery,
			})
		}

		// Update the store with the new menus
		if err := m.storeCli.Update(ctx, storeCoreModel.UpdateStore{
			ID:    &store.ID,
			Menus: updateMenus,
		}); err != nil {
			return "", err
		}
	}

	posMenu, err := m.GetPosMenu(tx, store)
	if err != nil {
		// TODO
		return "", err
	}

	if store.MenuID != "" {
		query = query.SetMenuID(store.MenuID)
	}

	menuID, err := m.createOrUpdateMenu(tx, query, posMenu, store, entityChangesHistoryRequest)
	if err != nil {
		return "", err
	}

	if !query.HasMenuID() {
		store.MenuID = menuID
		// update stores
		if err = m.storeRepo.Update(tx, store); err != nil {
			return "", err
		}
	}

	if store.IikoCloud.IsExternalMenu && upsertToAggrMenu {
		if err := m.insertPosMenuInfoToAggregatorMenus(ctx, store, posMenu, entityChangesHistoryRequest); err != nil {
			return "", err
		}
	}

	if err = m.updatePricesInAggregatorMenus(ctx, store, posMenu, entityChangesHistoryRequest); err != nil {
		return "", err
	}

	if err = m.updateDeletedProductsInAggregators(ctx, store, posMenu); err != nil {
		return "", err
	}

	return menuID, nil
}

func (m *mnm) UpsertMenuByFields(ctx context.Context, query selector.Menu, fields models.UpdateFields, agg models.UpdateFieldsAggregators, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest) error {

	_, err := m.UpsertMenu(ctx, query, entityChangesHistoryRequest, true)
	if err != nil {
		log.Err(err).Msgf("error: upsertMenu")
		return err
	}

	store, err := m.storeRepo.Get(ctx, selector.EmptyStoreSearch().SetID(query.StoreID))
	if err != nil {
		log.Err(err).Msgf("error: Get Store for: %s", query.StoreID)
		return err
	}

	posMenu, err := m.GetPosMenu(ctx, store)
	if err != nil {
		log.Err(err).Msgf("error: Get pos menu for store id: %s", query.StoreID)
		return err
	}

	err = m.updateAggMenuByFields(ctx, store, posMenu, fields, agg, entityChangesHistoryRequest)
	if err != nil {
		log.Err(err).Msgf("error:  updateAggMenuByFields for store id: %s", query.StoreID)
		return err
	}
	return nil
}

func (m *mnm) UploadMenu(ctx context.Context, storeId, menuId, aggregatorName string, sv3 *s3.S3, userRole, userName string) (string, error) {

	store, err := m.storeRepo.Get(ctx,
		selector.EmptyStoreSearch().
			SetID(storeId),
	)
	if err != nil {
		if errors.Is(err, drivers.ErrNotFound) {
			return "", fmt.Errorf("upload menu error: %w menu_id %s in store %s", drivers.ErrNotFound, menuId, storeId)
		}
		return "", err
	}

	menu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuId))
	if err != nil {
		if errors.Is(err, drivers.ErrNotFound) {
			return "", fmt.Errorf("upload menu error: %w menu_id %s ", drivers.ErrNotFound, menuId)
		}
		return "", err
	}

	if err = m.menuValidator.ValidateMenu(ctx, menu); err != nil {
		return "", err
	}

	if aggregatorName == "" && menu.Delivery != "" {
		aggregatorName = menu.Delivery
	}

	bkOffers, err := m.bkOffersRepo.List(ctx, selector.EmptyBkOffersSearch())
	if err != nil {
		if !errors.Is(err, drivers.ErrNotFound) {
			return "", fmt.Errorf("upload menu error: %w menu_id %s in store %s", err, menuId, storeId)
		}
	}

	trx, err := m.uploadAggregatorMenu(ctx, storeModels.AggregatorName(aggregatorName), store, menuId, menu, bkOffers, sv3, userRole)
	if err != nil {
		return "", fmt.Errorf("upload menu error: %s menu_id %s in store_id %s", err, menuId, store.ID)
	}

	menuDBVersionUrl, err := m.updateMenuDBVersionToS3(ctx, store.ID, m.globalConfig.S3_BUCKET.KwaakaMenuFilesBucket, aggregatorName, m.globalConfig.S3_BUCKET.ShareMenuBaseUrl, menu, sv3)
	if err != nil {
		log.Err(err).Msgf("update menu db verions in S3 error")
	}

	mut := models.MenuUploadTransaction{
		StoreID:          store.ID,
		ExtTransactions:  trx,
		Service:          menu.Delivery,
		UserName:         userName,
		MenuDBVersionUrl: menuDBVersionUrl,
	}
	if mut.ExtTransactions.HasNotSuccessProcessingStatus() {
		mut.Status = models.ERROR.String()
	} else if mut.ExtTransactions.HasProcessingStatus() {
		mut.Status = models.PROCESSING.String()
	} else {
		mut.Status = models.SUCCESS.String()
	}

	trId, err := m.mutMan.Create(ctx, mut)

	if err != nil {
		log.Err(err).Msgf("could not create transaction in store id %s", storeId)
		return "", fmt.Errorf("create upload menu error: %s menu_id %s in store_id %s", err, menuId, store.ID)
	}

	aggrManager, err := m.getAggregatorManager(ctx, store, storeModels.AggregatorName(aggregatorName))
	if err != nil {
		return trId, err
	}

	mutValid, err := aggrManager.ValidateMenu(ctx, models.MenuValidateRequest{
		MenuUploadTransaction: mut,
		Menu:                  menu,
		OffersBK:              bkOffers,
	})
	if err != nil {
		mutValid.ID = trId
		mutValid.CreatedAt.Value = coreModels.TimeNow()
		if err := m.mutMan.Update(ctx, mutValid); err != nil {
			return trId, err
		}
		return trId, err
	}

	return trId, nil
}

func (m *mnm) updateMenuDBVersionToS3(ctx context.Context, storeID, bucketName, deliveryService, shareMenuURL string, menu interface{}, sv3 *s3.S3) (string, error) {
	menuDBVersionUrl, err := pkgUtils.UploadMenuDBVersionToS3(storeID, bucketName, deliveryService, shareMenuURL, menu, sv3)
	if err != nil {
		log.Err(err).Msgf("could not save menu DB version in S3 in store %s", storeID)
		return "", fmt.Errorf("save menu in S3 error: %w in store_id %s", err, storeID)
	}

	menuUpTrnsctns, _, err := m.mutMan.List(ctx, selector.EmptyMenuUploadTransactionSearch().
		SetStoreID(storeID).
		SetService(models.AggregatorName(deliveryService)).
		SetSorting("created_at.value", -1).
		SetLimit(5))
	if err != nil {
		return "", err
	}

	var menuToDeleteUrl string
	menuToDeleteCount := 4
	for i, tr := range menuUpTrnsctns {
		if i != menuToDeleteCount {
			continue
		}
		menuToDeleteUrl = tr.MenuDBVersionUrl
	}

	if menuToDeleteUrl != "" {
		urlArr := strings.Split(menuToDeleteUrl, "kwaaka.com/")

		_, err = sv3.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(urlArr[1]),
		})
		if err != nil {
			log.Err(err).Msgf("could't delete menu DB version from S3 in store %s", storeID)
			return "", fmt.Errorf("delete old version menu from S3 error: %w in store_id %s", err, storeID)
		}
	}

	return menuDBVersionUrl, nil
}

func (m *mnm) updateChocofoodMenu(ctx context.Context, menuId, extStoreID string, aggregatorManager aggregator.Base) error {

	chocofoodMenu, err := aggregatorManager.GetMenu(ctx, extStoreID)
	if err != nil {
		return err
	}

	attributesMap := make(map[string]models.Attribute, len(chocofoodMenu.Attributes))
	for _, attribute := range chocofoodMenu.Attributes {
		attributesMap[attribute.ExtID] = attribute
	}

	productsMap := make(map[string]models.Product, len(chocofoodMenu.Products))
	for _, product := range chocofoodMenu.Products {
		productsMap[product.ExtID] = product
	}

	menu, err := m.menuRepo.Get(ctx, selector.Menu{
		ID: menuId,
	})

	if err != nil {
		return err
	}

	for idx, product := range menu.Products {
		chocofoodProduct, ok := productsMap[product.ExtID]

		if !ok {
			log.Warn().Msgf("Product %s %s not found in Chocofood menu", product.ExtID, product.Name[0].Value)
			continue
		}

		menu.Products[idx].ChocofoodFoodId = chocofoodProduct.ChocofoodFoodId
	}

	for idx, attribute := range menu.Attributes {
		chocofoodAttribute, ok := attributesMap[attribute.ExtID]

		if !ok {
			log.Warn().Msgf("Attribute %s %s not found in Chocofood menu", attribute.ExtID, attribute.Name)
			continue
		}

		menu.Attributes[idx].ChocofoodFoodId = chocofoodAttribute.ChocofoodFoodId
	}

	err = m.menuRepo.Update(ctx, menu, entityChangesHistoryModels.EntityChangesHistory{
		CallFunction: "updateChocofoodMenu",
		Author:       "cron",
	})

	if err != nil {
		return err
	}

	return nil
}

func (m *mnm) VerifyUploadMenu(ctx context.Context, transaction models.MenuUploadTransaction, store storeModels.Store) (models.MenuUploadTransaction, error) {

	if len(transaction.ExtTransactions) < 1 {
		return transaction, errors.New("no external transactions")
	}

	aggregatorManager, err := m.getAggregatorManager(ctx, store, storeModels.AggregatorName(transaction.Service))
	if err != nil {
		return transaction, err
	}

	menuUploadStatuses := make([]models.Status, 0, len(transaction.ExtTransactions))

	for idx, extTransaction := range transaction.ExtTransactions {
		status, err := aggregatorManager.VerifyMenu(ctx, extTransaction)
		if err != nil {
			log.Err(err).Msgf("error verifying transaction %s: %+v", extTransaction.ID, extTransaction)
		}
		transaction.ExtTransactions[idx].Status = status.String()
		menuUploadStatuses = append(menuUploadStatuses, status)

		if status.String() == models.SUCCESS.String() {

			if transaction.Service == models.CHOCOFOOD.String() {
				log.Info().Msgf("Start updating Chocofood menu %s", extTransaction.MenuID)
				err := m.updateChocofoodMenu(ctx, extTransaction.MenuID, extTransaction.ExtStoreID, aggregatorManager)
				if err != nil {
					log.Err(err).Msgf("Cant match chocofood products in store: %s", extTransaction.ExtStoreID)
					return transaction, err
				}
			}
		}

	}

	for _, mutStatus := range models.TransactionStatuses {
		for _, badStatus := range menuUploadStatuses {
			if badStatus.String() == mutStatus.String() {
				transaction.Status = mutStatus.String()
			}
		}
	}

	err = m.mutMan.Update(ctx, transaction)

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (m *mnm) setMenuByGroup(posMenu models.Menu, groupID string) (models.Menu, error) {

	menu := models.Menu{
		UpdatedAt: coreModels.TimeNow(),
	}

	if groupID == "" {
		return posMenu, nil
	}

	groups := make(map[string]struct{}, len(posMenu.Groups)+1)
	if !menu.Groups.IsExist(groupID) {
		groups[groupID] = struct{}{}
	}

	for _, group := range posMenu.Groups {
		if groupID == "" {
			groups[group.ID] = struct{}{}
			continue
		}
		if group.ParentGroup == "" && group.ID != groupID {
			continue
		}

		if group.ParentGroup == groupID || group.ID == groupID {
			groups[group.ID] = struct{}{}
		}
	}

	products := make(models.Products, 0, len(posMenu.Products))
	for _, product := range posMenu.Products {
		if _, ok := groups[product.ParentGroupID]; ok {
			products = append(products, product)
		}
	}

	attributes := make(models.Attributes, 0, len(posMenu.Attributes))
	for _, attribute := range posMenu.Attributes {
		if _, ok := groups[attribute.ParentAttributeGroup]; ok {
			attributes = append(attributes, attribute)
		}
	}

	menu.Products = products
	menu.Attributes = attributes
	menu.AttributesGroups = posMenu.AttributesGroups
	menu.Sections = posMenu.Sections
	menu.Groups = posMenu.Groups
	menu.StopLists = posMenu.StopLists

	return menu, nil
}

func (m *mnm) setCookingTime(ctx context.Context, systemPosMenuID string, newPosMenuProducts []models.Product) ([]models.Product, error) {

	systemPosMenu, err := m.menuRepo.Get(ctx, selector.Menu{ID: systemPosMenuID})
	if err != nil {
		return nil, err
	}

	mapProduct := make(map[string]int32)

	for _, product := range systemPosMenu.Products {
		mapProduct[product.ExtID] = product.CookingTime
	}

	for _, product := range newPosMenuProducts {
		product.CookingTime = mapProduct[product.ExtID]
	}

	return newPosMenuProducts, nil
}

func (m *mnm) createOrUpdateMenu(ctx context.Context, query selector.Menu, posMenu models.Menu, store storeModels.Store, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest) (string, error) {

	resMenu, err := m.setMenuByGroup(posMenu, query.GroupID)
	if err != nil {
		return "", err
	}

	if query.HasMenuID() {

		resMenuProducts, err := m.setCookingTime(ctx, store.MenuID, resMenu.Products)
		if err != nil {
			return "", err
		}

		menu, err := m.GetMenuByID(ctx, selector.EmptyMenuSearch().
			SetMenuID(query.MenuID()))
		if err != nil {
			// TODO
			return "", err
		}

		menu.Name = posMenu.Name
		menu.CreatedAt = posMenu.CreatedAt
		menu.UpdatedAt = posMenu.UpdatedAt
		menu.Delivery = posMenu.Description
		menu.Products = m.WotMenuSaveProductLanguagesAndProductInformation(resMenuProducts, menu.Products)
		menu.Attributes = m.SaveAttributeLanguagesInMenu(resMenu.Attributes, menu.Attributes)
		menu.AttributesGroups = m.SaveAttributeGroupLanguagesByExtIDInMenu(posMenu.AttributesGroups, menu.AttributesGroups)
		menu.Sections = m.SaveSectionLanguagesInMenu(posMenu.Sections, menu.Sections)
		menu.StopLists = posMenu.StopLists
		menu.Combos = posMenu.Combos

		if len(posMenu.Collections) > 0 {
			menu.Collections = posMenu.Collections
		}

		if err = m.menuRepo.Update(ctx, menu, entityChangesHistoryModels.EntityChangesHistory{
			CallFunction: "createOrUpdateMenu",
			Author:       entityChangesHistoryRequest.Author,
			TaskType:     entityChangesHistoryRequest.TaskType,
		}); err != nil {
			return "", err
		}

		return menu.ID, nil
	}

	menuID, err := m.menuRepo.Insert(ctx, posMenu)
	if err != nil {
		return "", err
	}

	return menuID, nil
}

func (m *mnm) verifyMenu(ctx context.Context, store storeModels.Store, aggregatorName storeModels.AggregatorName, tr models.ExtTransaction) (models.Status, error) {

	aggregatorMan, err := m.getAggregatorManager(ctx, store, aggregatorName)
	if err != nil {
		return "", err
	}

	status, err := aggregatorMan.VerifyMenu(ctx, tr)
	if err != nil {
		return status, err
	}

	return status, nil
}

func (m *mnm) validateMenuToUpload(ctx context.Context, menuId string) models.Status {

	menu, err := m.GetMenuByID(ctx, selector.EmptyMenuSearch().SetMenuID(menuId))
	if err != nil {
		log.Err(err).Msgf("get menu by id %s", menuId)
		return models.NOT_READY
	}

	if err = m.menuValidator.ValidateMenu(ctx, menu); err != nil {
		log.Err(err).Msgf("get menu by id %s", menuId)
		return models.NOT_READY
	}

	return models.READY
}

func setStatus(status string) models.Status {
	if !models.Status(status).ValidStatus([]models.Status{
		models.PARTIALLY_PROCESSED,
		models.NOT_PROCESSED,
		models.SUCCESS,
		models.PROCESSING,
		models.FAILED,
	}) {
		return models.NOT_PROCESSED
	}

	return models.Status(status)
}

func (m *mnm) uploadAggregatorMenu(ctx context.Context, aggregatorName storeModels.AggregatorName, store storeModels.Store, menuId string, menu models.Menu, offers []models.BkOffers, sv3 *s3.S3, userRole string) ([]models.ExtTransaction, error) {

	aggregatorMan, err := m.getAggregatorManager(ctx, store, aggregatorName)
	if err != nil {
		return nil, err
	}

	extStoreIDs := store.GetAggregatorStoreIDs(aggregatorName.String())
	if len(extStoreIDs) == 0 {
		return nil, validator.ErrEmptyStores
	}

	transactions := make([]models.ExtTransaction, 0, len(extStoreIDs))

	for _, storeAggrID := range extStoreIDs {

		transaction := models.TransactionData{
			StoreID:  storeAggrID,
			Delivery: aggregatorName.String(),
		}

		rsp, err := aggregatorMan.UploadMenu(ctx, menuId, storeAggrID, menu, store, offers, sv3, userRole)
		if err != nil {
			log.Trace().Err(err).Msgf("could not upload menu store: %s with external store_id from aggregator %s, aggregatorStoreID %s menuId %s", store.ID, aggregatorName, storeAggrID, menuId)
			transaction.Status = models.ERROR
			if errors.Is(err, aggregatorErrors.ErrNoPermissionForPublishWoltMenu) {
				return transactions, err
			}
		}

		transactions = append(transactions, rsp)
	}

	return transactions, nil

}

func (m *mnm) getStopListName(ctx context.Context, menuId string) ([]string, error) {

	menu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuId))
	if err != nil {
		return []string{}, err
	}

	stopsExist := make(map[string]struct{}, len(menu.StopLists))
	for _, id := range menu.StopLists {
		stopsExist[id] = struct{}{}
	}

	stopLists := make([]string, 0, len(menu.StopLists))

	for _, product := range menu.Products {
		if _, ok := stopsExist[product.ProductID]; ok {
			stopLists = append(stopLists, product.Name[0].Value)
		}
	}

	for _, attribute := range menu.Attributes {
		if _, ok := stopsExist[attribute.ExtID]; ok {
			stopLists = append(stopLists, attribute.Name)
		}
	}

	return stopLists, nil

}

func (m *mnm) RecoveryMenu(ctx context.Context, menu models.Menu, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest) error {
	menu.IsDeleted = false
	return m.menuRepo.Update(ctx, menu, entityChangesHistoryModels.EntityChangesHistory{
		CallFunction: "RecoveryMenu",
		Author:       entityChangesHistoryRequest.Author,
		TaskType:     entityChangesHistoryRequest.TaskType,
	})
}

func (m *mnm) DeleteProducts(ctx context.Context, menuId string, productsIds []string) error {
	return m.menuRepo.DeleteProducts(ctx, menuId, productsIds)
}

func (m *mnm) DeleteProductsFromDB(ctx context.Context, menuId string, productsIds []string) error {
	return m.menuRepo.DeleteProductsFromDB(ctx, menuId, productsIds)
}

func (m *mnm) GetAttributesForUpdate(ctx context.Context, query selector.Menu) (models.Attributes, models.AttributeGroups, int, error) {

	attributes, total, err := m.menuRepo.GetAttributes(ctx, query)
	if err != nil {
		return nil, nil, total, err
	}

	attributesGroups, err := m.menuRepo.GetAttributeGroups(ctx, query)
	if err != nil {
		return nil, nil, total, err
	}

	return attributes, attributesGroups, total, nil
}

func (m *mnm) UpdateMenu(ctx context.Context, menu models.Menu) error {
	return m.menuRepo.Update(ctx, menu, entityChangesHistoryModels.EntityChangesHistory{})
}

func (m *mnm) AddRowToAttributeGroup(ctx context.Context, menuId string) error {

	attributesGroups, err := m.menuRepo.GetAttributeGroups(ctx, selector.EmptyMenuSearch().SetMenuID(menuId))
	if err != nil {
		return err
	}

	var attributeMinMax []models.AttributeIdMinMax
	var tempAttributeMinMax models.AttributeIdMinMax

	for i := range attributesGroups {
		attributeGroup := attributesGroups[i]
		for j := range attributeGroup.Attributes {
			attribute := attributeGroup.Attributes[j]

			tempAttributeMinMax.ExtId = attribute
			tempAttributeMinMax.Min = 0
			tempAttributeMinMax.Max = attributeGroup.Max

			attributeMinMax = append(attributeMinMax, tempAttributeMinMax)
		}
		if err := m.menuRepo.AddRowToAttributeGroup(ctx, menuId, attributeMinMax, attributeGroup.ExtID); err != nil {
			return err
		}
		tempAttributeMinMax = models.AttributeIdMinMax{}
		attributeMinMax = []models.AttributeIdMinMax{}
	}

	return nil
}

func (m *mnm) UpdateMenuName(ctx context.Context, query models.UpdateMenuName) error {
	return m.menuRepo.UpdateMenuName(ctx, query)
}

func (m *mnm) DeleteAttributeGroupFromDB(ctx context.Context, menuId string, attrGroupExtId string) error {
	return m.menuRepo.DeleteAttributeGroupFromDB(ctx, menuId, attrGroupExtId)
}

func (m *mnm) ValidateAttributeGroupName(ctx context.Context, menuId, name string) (bool, error) {
	return m.menuRepo.ValidateAttributeGroupName(ctx, menuId, name)
}

func (m *mnm) CreateAttributeGroup(ctx context.Context, menuID, attrGroupName string, min, max int) (string, error) {

	attribute := models.Attribute{
		ExtID: uuid.New().String(),
		Name:  attrGroupName,
		Max:   max,
		Min:   min,
	}

	return m.menuRepo.CreateAttributeGroup(ctx, menuID, attribute)
}

func (m *mnm) AutoUploadMenuByPOS(ctx context.Context, menu models.Menu, query selector.Menu) (string, error) {
	if query.HasMenuID() {
		err := m.menuRepo.Update(ctx, menu, entityChangesHistoryModels.EntityChangesHistory{})
		if err != nil {
			return "", err
		}
		return menu.ID, nil
	}
	menuId, err := m.menuRepo.Insert(ctx, menu)
	if err != nil {
		return "", err
	}

	return menuId, nil
}

func (m *mnm) AutoUpdateMenuPrices(ctx context.Context, storeId string) error {
	store, err := m.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{ID: storeId})
	if err != nil {
		return err
	}
	posMenu, err := m.GetMenuByID(ctx, selector.EmptyMenuSearch().SetMenuID(store.MenuID))
	if err != nil {
		return err
	}

	posMenuPrices := make(map[string]float64)
	for _, product := range posMenu.Products {
		if len(product.Price) > 0 {
			posMenuPrices[product.ExtID] = product.Price[0].Value
		}
	}

	for i := range store.Menus {
		if store.Menus[i].IsActive {
			aggrMenu, err := m.GetMenuByID(ctx, selector.EmptyMenuSearch().SetMenuID(store.Menus[i].ID))
			if err != nil {
				log.Trace().Err(err).Msgf("couldn't aggregator menu with menu id:: %s", aggrMenu.ID)
				return err
			}

			for j := range aggrMenu.Products {
				id := aggrMenu.Products[j].ExtID
				if aggrMenu.Products[j].PosID != "" {
					id = aggrMenu.Products[j].PosID
				}
				if price, ok := posMenuPrices[id]; ok {
					if len(aggrMenu.Products[j].Price) > 0 {
						aggrMenu.Products[j].Price[0].Value = price
					}
				}
			}

			err = m.UpdateMenu(ctx, aggrMenu)
			if err != nil {
				log.Trace().Err(err).Msgf("couldn't update aggregator menu in database with menu id: %s", aggrMenu.ID)
				return err
			}
		}
	}
	return nil
}
func (m *mnm) AutoUpdateMenuDescriptions(ctx context.Context, storeId string) error {
	store, err := m.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{ID: storeId})
	if err != nil {
		return err
	}
	posMenu, err := m.GetMenuByID(ctx, selector.EmptyMenuSearch().SetMenuID(store.MenuID))
	if err != nil {
		return err
	}

	posMenuDesc := make(map[string]models.LanguageDescription)
	for _, product := range posMenu.Products {
		if len(product.Description) > 0 {
			posMenuDesc[product.ExtID] = product.Description[0]
		}
	}

	for i := range store.Menus {
		if store.Menus[i].IsActive {
			aggrMenu, err := m.GetMenuByID(ctx, selector.EmptyMenuSearch().SetMenuID(store.Menus[i].ID))
			if err != nil {
				log.Trace().Err(err).Msgf("couldn't aggregator menu with menu id:: %s", aggrMenu.ID)
				return err
			}

			for j := range aggrMenu.Products {
				id := aggrMenu.Products[j].ExtID
				if aggrMenu.Products[j].PosID != "" {
					id = aggrMenu.Products[j].PosID
				}
				if description, ok := posMenuDesc[id]; ok {
					if len(aggrMenu.Products[j].Description) > 0 {
						aggrMenu.Products[j].Description[0] = description
					}
				}
			}

			err = m.UpdateMenu(ctx, aggrMenu)
			if err != nil {
				log.Trace().Err(err).Msgf("couldn't update aggregator menu in database with menu id: %s", aggrMenu.ID)
				return err
			}
		}
	}
	return nil
}

func (m *mnm) insertPosMenuInfoToAggregatorMenus(ctx context.Context, store storeModels.Store, posMenu models.Menu, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest) error {
	for _, menu := range store.Menus {
		if !menu.IsActive || menu.IsDeleted {
			continue
		}
		aggrMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menu.ID))
		if err != nil {
			return err
		}

		newAggrMenu := posMenu

		if aggrMenu.Delivery == models.WOLT.String() {
			newAggrMenu = m.WoltMenuSaveLanguages(newAggrMenu, aggrMenu)
			newAggrMenu.Products = m.SaveProductInformationInMenu(newAggrMenu.Products, aggrMenu.Products)
		}

		newAggrMenu.ID = aggrMenu.ID
		newAggrMenu.Name = aggrMenu.Name
		newAggrMenu.Description = aggrMenu.Description
		newAggrMenu.Delivery = aggrMenu.Delivery

		if aggrMenu.Delivery == models.GLOVO.String() && len(aggrMenu.Collections) > 1 {
			newAggrMenu.Collections = aggrMenu.Collections
			newAggrMenu.Sections = aggrMenu.Sections
		}

		err = m.menuRepo.Update(ctx, newAggrMenu, entityChangesHistoryModels.EntityChangesHistory{
			CallFunction: "insertPosMenuInfoToAggregatorMenus",
			Author:       entityChangesHistoryRequest.Author,
			TaskType:     entityChangesHistoryRequest.TaskType,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mnm) updateAggMenuByFields(ctx context.Context, store storeModels.Store, posMenu models.Menu, fields models.UpdateFields, agg models.UpdateFieldsAggregators, entityChangesHistoryRequest entityChangesHistoryModels.EntityChangesHistoryRequest) error {

	mapAgg := map[string]bool{
		"wolt":   agg.Wolt,
		"glovo":  agg.Glovo,
		"yandex": agg.Yandex,
	}

	for _, menuObj := range store.Menus {
		if !menuObj.IsActive || menuObj.IsDeleted {
			continue
		}

		updateAgg, ok := mapAgg[menuObj.Delivery]
		if !updateAgg || !ok {
			continue
		}

		log.Info().Msgf("start to update %s menu with id %s and match with new pos menu by fields: %#v for store: %s", menuObj.Delivery, menuObj.ID, fields, store.ID)

		aggMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuObj.ID))
		if err != nil {
			return err
		}

		var (
			products               = make(map[string]models.Product)
			updatedProducts        = make([]models.Product, 0, len(aggMenu.Products))
			attributeGroups        = make(map[string]models.AttributeGroup)
			updatedAttributeGroups = make([]models.AttributeGroup, 0, len(aggMenu.AttributesGroups))
			attributes             = make(map[string]models.Attribute)
			updatedAttributes      = make([]models.Attribute, 0, len(aggMenu.Attributes))
		)

		if len(aggMenu.Products) > 0 && len(posMenu.Products) > 0 {
			products = m.updateAggMenuProductsByFields(posMenu, aggMenu, fields)
		}

		if len(aggMenu.AttributesGroups) > 0 && len(posMenu.AttributesGroups) > 0 {
			attributeGroups = m.updateAggMenuAttributeGroupsByFields(posMenu, aggMenu, fields)
		}

		if len(aggMenu.Attributes) > 0 && len(posMenu.Attributes) > 0 {
			attributes = m.updateAggMenuAttributeByFields(posMenu, aggMenu, fields)
		}

		for _, product := range products {
			updatedProducts = append(updatedProducts, product)
		}
		for _, attributeGroup := range attributeGroups {
			updatedAttributeGroups = append(updatedAttributeGroups, attributeGroup)
		}
		for _, attribute := range attributes {
			updatedAttributes = append(updatedAttributes, attribute)
		}

		aggMenu.Products = updatedProducts
		aggMenu.AttributesGroups = updatedAttributeGroups
		aggMenu.Attributes = updatedAttributes

		err = m.menuRepo.Update(ctx, aggMenu, entityChangesHistoryModels.EntityChangesHistory{
			CallFunction: "updateAggMenuByFields",
			Author:       entityChangesHistoryRequest.Author,
			TaskType:     entityChangesHistoryRequest.TaskType,
		})
		log.Info().Msgf("success updateAggMenuByFields for menu_id: %s", menuObj.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *mnm) updateAggMenuProductsByFields(posMenu models.Menu, aggMenu models.Menu, fields models.UpdateFields) map[string]models.Product {
	var aggMenuProductsMap = make(map[string]models.Product)

	for _, product := range aggMenu.Products {
		aggMenuProductsMap[product.ExtID] = product
	}

	for _, posMenuProduct := range posMenu.Products {

		aggMenuProduct, ok := aggMenuProductsMap[posMenuProduct.ExtID]
		if !ok {
			continue
		}

		if fields.ProductName && len(aggMenuProduct.Name) > 0 && len(posMenuProduct.Name) > 0 {
			for i, aggMenuProductName := range aggMenuProduct.Name {
				for _, posMenuProductName := range posMenuProduct.Name {
					if (len(aggMenuProduct.Name) == 1 || posMenuProductName.LanguageCode == aggMenuProductName.LanguageCode) &&
						posMenuProductName.Value != aggMenuProductName.Value &&
						len(posMenuProductName.Value) != 0 {
						aggMenuProduct.Name[i].Value = posMenuProductName.Value
					}
				}
			}
		}

		if fields.ProductPrice && len(aggMenuProduct.Price) > 0 && len(posMenuProduct.Price) > 0 {
			for i, aggMenuProductPrice := range aggMenuProduct.Price {
				for _, posMenuProductPrice := range posMenuProduct.Price {
					if posMenuProductPrice.CurrencyCode == aggMenuProductPrice.CurrencyCode &&
						posMenuProductPrice.Value != aggMenuProductPrice.Value &&
						posMenuProductPrice.Value != 0 {
						aggMenuProduct.Price[i].Value = posMenuProductPrice.Value
					}
				}
			}
		}

		if fields.ProductDescription && len(aggMenuProduct.Description) > 0 && len(posMenuProduct.Description) > 0 {
			for i, aggMenuProductDescription := range aggMenuProduct.Description {
				for _, posMenuProductDescription := range posMenuProduct.Description {
					if (len(aggMenuProduct.Description) == 1 || aggMenuProductDescription.LanguageCode == posMenuProductDescription.LanguageCode) &&
						aggMenuProductDescription.Value != posMenuProductDescription.Value &&
						len(posMenuProductDescription.Value) != 0 {
						aggMenuProduct.Description[i].Value = posMenuProductDescription.Value
					}
				}
			}
		}

		if fields.ProductImage && len(aggMenuProduct.ImageURLs) > 0 && len(posMenuProduct.ImageURLs) > 0 {
			aggMenuProduct.ImageURLs[0] = posMenuProduct.ImageURLs[0]
		}

		aggMenuProductsMap[posMenuProduct.ExtID] = aggMenuProduct
	}
	log.Info().Msgf("success updateAggMenuProductsByFields for agg menuID: %s", aggMenu.ID)
	return aggMenuProductsMap
}

func (m *mnm) updateAggMenuAttributeGroupsByFields(posMenu models.Menu, aggMenu models.Menu, fields models.UpdateFields) map[string]models.AttributeGroup {

	var (
		aggMenuAttributeGroupMap = make(map[string]models.AttributeGroup)
		posMenuAttributeGroup    = make(map[string]models.AttributeGroup) //key =  AttributeGroup.Attributes[0]
	)

	for _, posAttributeGroup := range posMenu.AttributesGroups {
		posMenuAttributeGroup[posAttributeGroup.Attributes[0]] = posAttributeGroup
	}

	for _, aggAttributeGroup := range aggMenu.AttributesGroups {

		posAttributeGroup, ok := posMenuAttributeGroup[aggAttributeGroup.Attributes[0]]
		if !ok {
			aggMenuAttributeGroupMap[aggAttributeGroup.ExtID] = aggAttributeGroup
			continue
		}

		if fields.AttributeGroupName && len(posAttributeGroup.Name) > 0 && posAttributeGroup.Name != aggAttributeGroup.Name {
			aggAttributeGroup.Name = posAttributeGroup.Name
		}

		if fields.AttributeGroupMinMax {

			if posAttributeGroup.Max != aggAttributeGroup.Max {
				aggAttributeGroup.Max = posAttributeGroup.Max
			}
			if posAttributeGroup.Min != aggAttributeGroup.Min {
				aggAttributeGroup.Min = posAttributeGroup.Min
			}
		}

		aggMenuAttributeGroupMap[aggAttributeGroup.ExtID] = aggAttributeGroup
	}
	log.Info().Msgf("success updateAggMenuAttributeGroupsByFields for agg menu_id: %s", aggMenu.ID)
	return aggMenuAttributeGroupMap
}

func (m *mnm) updateAggMenuAttributeByFields(posMenu models.Menu, aggMenu models.Menu, fields models.UpdateFields) map[string]models.Attribute {

	var aggMenuAttributesMap = make(map[string]models.Attribute)

	for _, aggAttribute := range aggMenu.Attributes {
		aggMenuAttributesMap[aggAttribute.ExtID] = aggAttribute
	}

	for _, posAttribute := range posMenu.Attributes {

		aggAttribute, ok := aggMenuAttributesMap[posAttribute.ExtID]
		if !ok {
			aggMenuAttributesMap[posAttribute.ExtID] = aggAttribute
			continue
		}
		if fields.AttributeName && len(posAttribute.Name) > 0 && posAttribute.Name != aggAttribute.Name {
			aggAttribute.Name = posAttribute.Name
		}
		if fields.AttributePrice && posAttribute.Price > 0 && posAttribute.Price != aggAttribute.Price {
			aggAttribute.Price = posAttribute.Price
		}

		aggMenuAttributesMap[posAttribute.ExtID] = aggAttribute

	}
	log.Info().Msgf("success updateAggMenuAttributeByFields for agg menu_id: %s", aggMenu.ID)
	return aggMenuAttributesMap

}

func (m *mnm) updateDeletedProductsInAggregators(ctx context.Context, store storeModels.Store, posMenu models.Menu) error {

	deletedPOSProducts := make(map[string]models.Product)
	for _, product := range posMenu.Products {
		if product.IsDeleted || (store.PosType == models.IIKO.String() && !product.IsIncludedInMenu) {
			deletedPOSProducts[product.ExtID] = product
		}
	}

	for _, menu := range store.Menus {
		if !menu.IsActive {
			continue
		}

		aggrMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menu.ID))
		if err != nil {
			return err
		}

		productsForUpdate := m.getProductsIDForUpdate(aggrMenu.Products, deletedPOSProducts)

		if productsForUpdate != nil {
			log.Info().Msgf("update deleted POS products in aggregator %v", productsForUpdate)

			if err := m.updateDelProductsInDB(ctx, menu.ID, productsForUpdate); err != nil {
				return err
			}

			if err := m.updateDelProductsInAggregator(ctx, store, productsForUpdate, menu.Delivery); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *mnm) getProductsIDForUpdate(aggregatorMenuProducts []models.Product, deletedPOSProducts map[string]models.Product) []models.Product {
	var productIDsForUpdate []models.Product
	for _, product := range aggregatorMenuProducts {
		if _, ok := deletedPOSProducts[product.PosID]; ok {

			if product.IsAvailable {
				product.IsAvailable = false
				productIDsForUpdate = append(productIDsForUpdate, product)
			}

		}
	}
	return productIDsForUpdate
}

func (m *mnm) updateDelProductsInDB(ctx context.Context, menuID string, products []models.Product) error {
	var productsForUpdate []string
	for _, product := range products {
		productsForUpdate = append(productsForUpdate, product.ExtID)
	}

	if err := m.menuServiceRepo.BulkUpdateProductsAvailability(ctx, menuID, productsForUpdate, false); err != nil {
		return err
	}

	return nil
}

func (m *mnm) updateDelProductsInAggregator(ctx context.Context, store storeModels.Store, products []models.Product, deliveryService string) error {
	externalStoreIDs := store.GetAggregatorStoreIDs(deliveryService)

	aggregatorManager, err := m.getAggregatorManager(ctx, store, storeModels.AggregatorName(deliveryService))
	if err != nil {
		return err
	}

	for _, aggrStoreID := range externalStoreIDs {
		_, err := aggregatorManager.BulkUpdate(ctx, store.ID, aggrStoreID, products, nil, store)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mnm) CreateGlovoSuperCollection(ctx context.Context, menuId string, superCollections dto.MenuSuperCollections) error {
	return m.menuRepo.CreateGlovoSuperCollection(ctx, menuId, superCollections)
}

func (m *mnm) CreateMenuByAggregatorAPI(ctx context.Context, aggregator string, storeId string) (string, error) {
	store, err := m.storeRepo.Get(ctx, selector.EmptyStoreSearch().SetID(storeId))
	if err != nil {
		return "", err
	}
	aggMan, err := m.getAggregatorManager(ctx, store, storeModels.AggregatorName(aggregator))
	if err != nil {
		return "", err
	}

	storeIds := store.GetAggregatorStoreIDs(aggregator)

	var newAggMenu models.Menu

	for _, storeId := range storeIds {
		newAggMenu, err = aggMan.GetMenu(ctx, storeId)
		if err != nil {
			return "", err
		}
	}

	id, err := m.menuRepo.Insert(ctx, newAggMenu)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (m *mnm) isIgnoreStoreList(storeID string) bool {
	storeMap := make(map[string]bool)
	storeMap["637b25c5c587b2336a8cfa2f"] = true
	storeMap["637b3992327bda69ad8f8958"] = true
	storeMap["637b279bc587b2336a8cfa3b"] = true
	storeMap["637b1680fcfe40c06630c19b"] = true
	storeMap["63e373e7e2400cb139baf0a2"] = true
	storeMap["2553f72f6eadb44a0b680a50"] = true
	storeMap["6353f72f6eadb44a0b680a20"] = true
	storeMap["637b1a197a76ddafe546ae40"] = true
	storeMap["637f4b8b05aaa0aa0fe6a170"] = true
	storeMap["637b22020f6be3f87c8d8513"] = true
	storeMap["637b3886327bda69ad8f894c"] = true
	storeMap["637f4aad05aaa0aa0fe6a16b"] = true
	storeMap["637b20d90f6be3f87c8d850a"] = true
	storeMap["6400685952bbea97f99c0cdf"] = true
	storeMap["641d2ad4b06d88ca9be819bd"] = true
	storeMap["645c894879917d2290a0a9f1"] = true
	storeMap["64993e5a7c94740df1dadd91"] = true
	storeMap["64c8bcfc0e3b0c38c41eeb1c"] = true
	storeMap["65bcb04e06cf4f070500ba98"] = true
	return storeMap[storeID]
}

func (m *mnm) InsertMenu(ctx context.Context, menu models.Menu) (string, error) {
	return m.menuRepo.Insert(ctx, menu)
}

func (m *mnm) validateProductError(ctx context.Context, menuID string, storeID string) ([]models.Product, error) {

	var productErrorMsg models.ProductWithErr
	var productsWithErr []models.ProductWithErr
	var result []models.Product

	store, err := m.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{ID: storeID})
	if err != nil {
		return nil, err
	}

	posMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(store.MenuID))
	if err != nil {
		return nil, err
	}

	positions := make(map[string]models.ProductStatus)
	posAttributeGroupsMap := make(map[string]models.AttributeGroup)

	// write product status(del, includ, attrgroupIDs, defaults)
	for _, product := range posMenu.Products {
		var defaults []string

		for _, defaultAttribute := range product.MenuDefaultAttributes {
			if defaultAttribute.ByAdmin {
				defaults = append(defaults, defaultAttribute.ExtID)
			}
		}

		positions[product.ExtID] = models.ProductStatus{
			IsDeleted:         product.IsDeleted,
			IsIncludedInMenu:  product.IsIncludedInMenu,
			AttributeGroupIDs: product.AttributesGroups,
			Defaults:          defaults,
		}
	}

	for _, attribute := range posMenu.Attributes {
		positions[attribute.ExtID] = models.ProductStatus{
			Name:             attribute.Name,
			IsDeleted:        attribute.IsDeleted,
			IsIncludedInMenu: attribute.IncludedInMenu,
		}
	}

	for _, attributeGroup := range posMenu.AttributesGroups {
		posAttributeGroupsMap[attributeGroup.ExtID] = attributeGroup
	}

	aggMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuID))
	if err != nil {
		return nil, err
	}

	products, attributes, minMaxReports, err := m.validateAggregatorMenu(ctx, positions, aggMenu, posAttributeGroupsMap, storeID)
	if err != nil {
		return nil, err
	}

	for i := range products {
		product := products[i]
		productErrorMsg = models.ProductWithErr{
			ID:         product.ID,
			Name:       product.Name,
			ProductErr: product.Status,
			Solution:   product.Solution,
		}
		productsWithErr = append(productsWithErr, productErrorMsg)
	}

	for _, attribute := range attributes {
		for i := range attribute.Products {
			product := attribute.Products[i]
			productErrorMsg = models.ProductWithErr{
				ID:               product.ID,
				Name:             product.Name,
				ProductErr:       attribute.Status,
				AttributeGroupID: attribute.ID,
			}
			productsWithErr = append(productsWithErr, productErrorMsg)
		}
	}

	for i := range minMaxReports {
		product := minMaxReports[i]
		for attributeID := range product.AttributeMinMaxReports {
			attribute := product.AttributeMinMaxReports[attributeID]
			productErrorMsg = models.ProductWithErr{
				ID:               product.ID,
				Name:             product.Name,
				ProductErr:       "Min/Max (атрибут группы)/атрибута  не совпадают. Pos attribute group id: " + attribute.ID,
				AttributeGroupID: attribute.ID,
			}
			productsWithErr = append(productsWithErr, productErrorMsg)
		}
	}

	aggProducts, _, err := m.menuRepo.ListProducts(ctx, selector.EmptyMenuSearch().SetMenuID(menuID))
	if err != nil {
		return nil, err
	}

	for i := range productsWithErr {
		productErrObj := productsWithErr[i]
		for j := range aggProducts {
			productObj := aggProducts[j]
			if productErrObj.ID == productObj.ExtID {
				productObj.ProductErr = productErrObj.ProductErr

				result = append(result, productObj)
			}

		}
	}

	return result, nil
}

func (m *mnm) ValidateAggAndPosMatching(ctx context.Context, menuID string, storeID string, limit int) (aggregatorProducts []models.Product, posProducts []models.Product, total int, err error) {

	aggregatorProducts, err = m.validateProductError(ctx, menuID, storeID)
	if err != nil {
		return nil, nil, 0, err
	}

	store, err := m.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{ID: storeID})
	if err != nil {
		return nil, nil, 0, err
	}

	posMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(store.MenuID))
	if err != nil {
		return nil, nil, 0, err
	}

	var dict = make(map[string]*models.Product)
	var groups = make(map[string]*models.AttributeGroup)

	for i := 0; i < len(posMenu.Products); i++ {
		temp := posMenu.Products[i]
		dict[posMenu.Products[i].ExtID] = &temp
	}

	for i := 0; i < len(posMenu.AttributesGroups); i++ {
		temp := posMenu.AttributesGroups[i]
		groups[posMenu.AttributesGroups[i].ExtID] = &temp
	}

	for i := 0; i < len(aggregatorProducts); i++ {
		if aggregatorProducts[i].IsSync {
			if val, ok := dict[aggregatorProducts[i].ExtID]; ok {
				//added new poles to match pos is deleted in included in aggregator
				aggregatorProducts[i].IsIncludedInMenu = val.IsIncludedInMenu
				aggregatorProducts[i].IsPosIsDeleted = val.IsDeleted

				for _, groupID := range val.AttributesGroups {
					if group, ok := groups[groupID]; ok {
						if group.Min > 0 {
							aggregatorProducts[i].HasRequiredAttributes = true
							break
						}
					}
				}
			}
		}

	}

	posProducts = posMenu.Products

	return
}

func (m *mnm) GetEmptyProducts(ctx context.Context, menuID string, pagination selector.Pagination) ([]models.Product, int, error) {
	return m.menuRepo.GetEmptyProducts(ctx, menuID, pagination)
}

func (m *mnm) UpdateProductAvailableStatus(ctx context.Context, menuID, productID string, status bool) error {
	return m.menuRepo.UpdateProductAvailableStatus(ctx, menuID, productID, status)
}

// for save name, description array languages in wolt menu
func (m *mnm) WoltMenuSaveLanguages(newMenu, aggrMenu models.Menu) models.Menu {

	newMenu.Products = m.SaveProductLanguagesInMenu(newMenu.Products, aggrMenu.Products)
	newMenu.Sections = m.SaveSectionLanguagesInMenu(newMenu.Sections, aggrMenu.Sections)
	newMenu.AttributesGroups = m.SaveAttributeGroupLanguagesByNameInMenu(newMenu.AttributesGroups, aggrMenu.AttributesGroups)
	newMenu.Attributes = m.SaveAttributeLanguagesInMenu(newMenu.Attributes, aggrMenu.Attributes)

	return newMenu
}

func (m *mnm) SaveProductLanguagesInMenu(newMenuProducts, oldMenuProducts models.Products) models.Products {
	oldMenuProductMap := make(map[string]models.Product, len(oldMenuProducts))

	for _, product := range oldMenuProducts {
		oldMenuProductMap[product.ExtID] = product
	}

	for i := range newMenuProducts {
		if _, ok := oldMenuProductMap[newMenuProducts[i].ExtID]; !ok {
			continue
		}

		saveNameMap := make(map[string]models.LanguageDescription, len(newMenuProducts[i].Name))
		saveDescriptionMap := make(map[string]models.LanguageDescription, len(newMenuProducts[i].Description))

		for _, n := range newMenuProducts[i].Name {
			saveNameMap[n.LanguageCode] = n
		}
		for _, description := range newMenuProducts[i].Description {
			saveDescriptionMap[description.LanguageCode] = description
		}

		newMenuProducts[i].Name = oldMenuProductMap[newMenuProducts[i].ExtID].Name
		newMenuProducts[i].Description = oldMenuProductMap[newMenuProducts[i].ExtID].Description

		for j := range newMenuProducts[i].Name {
			if _, ok := saveNameMap[newMenuProducts[i].Name[j].LanguageCode]; !ok {
				continue
			}
			newMenuProducts[i].Name[j].Value = saveNameMap[newMenuProducts[i].Name[j].LanguageCode].Value
		}

		for j := range newMenuProducts[i].Description {
			if _, ok := saveDescriptionMap[newMenuProducts[i].Description[j].LanguageCode]; !ok {
				continue
			}
			newMenuProducts[i].Description[j].Value = saveDescriptionMap[newMenuProducts[i].Description[j].LanguageCode].Value
		}

	}

	return newMenuProducts
}

func (m *mnm) SaveSectionLanguagesInMenu(newMenuSections, oldMenuSections models.Sections) models.Sections {
	oldMenuSectionMap := make(map[string]models.Section, len(oldMenuSections))

	for _, section := range oldMenuSections {
		oldMenuSectionMap[section.ExtID] = section
	}

	for i := range newMenuSections {
		if _, ok := oldMenuSectionMap[newMenuSections[i].ExtID]; !ok {
			continue
		}

		newMenuSections[i].NamesByLanguage = oldMenuSectionMap[newMenuSections[i].ExtID].NamesByLanguage

		saveDescriptionMap := make(map[string]models.LanguageDescription, len(newMenuSections[i].Description))
		for _, description := range newMenuSections[i].Description {
			saveDescriptionMap[description.LanguageCode] = description
		}

		newMenuSections[i].Description = oldMenuSectionMap[newMenuSections[i].ExtID].Description

		for j := range newMenuSections[i].Description {
			if _, ok := saveDescriptionMap[newMenuSections[i].Description[j].LanguageCode]; !ok {
				continue
			}
			newMenuSections[i].Description[j].Value = saveDescriptionMap[newMenuSections[i].Description[j].LanguageCode].Value
		}
	}

	return newMenuSections
}

func (m *mnm) SaveAttributeGroupLanguagesByNameInMenu(newMenuAttributeGroups, oldMenuAttributeGroups models.AttributeGroups) models.AttributeGroups {
	oldMenuAttributeGroupMap := make(map[string]models.AttributeGroup, len(oldMenuAttributeGroups))

	for _, attrGroup := range oldMenuAttributeGroups {
		if _, ok := oldMenuAttributeGroupMap[attrGroup.Name]; ok && oldMenuAttributeGroupMap[attrGroup.Name].NamesByLanguage == nil && attrGroup.NamesByLanguage != nil {
			oldMenuAttributeGroupMap[attrGroup.Name] = attrGroup
			continue
		}
		oldMenuAttributeGroupMap[attrGroup.Name] = attrGroup
	}

	for i := range newMenuAttributeGroups {
		if _, ok := oldMenuAttributeGroupMap[newMenuAttributeGroups[i].Name]; !ok {
			continue
		}

		newMenuAttributeGroups[i].NamesByLanguage = oldMenuAttributeGroupMap[newMenuAttributeGroups[i].Name].NamesByLanguage

		saveDescriptionMap := make(map[string]models.LanguageDescription, len(newMenuAttributeGroups[i].Description))
		for _, description := range newMenuAttributeGroups[i].Description {
			saveDescriptionMap[description.LanguageCode] = description
		}

		newMenuAttributeGroups[i].Description = oldMenuAttributeGroupMap[newMenuAttributeGroups[i].Name].Description

		for j := range newMenuAttributeGroups[i].Description {
			if _, ok := saveDescriptionMap[newMenuAttributeGroups[i].Description[j].LanguageCode]; !ok {
				continue
			}
			newMenuAttributeGroups[i].Description[j].Value = saveDescriptionMap[newMenuAttributeGroups[i].Description[j].LanguageCode].Value
		}
	}

	return newMenuAttributeGroups
}

func (m *mnm) SaveAttributeGroupLanguagesByExtIDInMenu(newMenuAttributeGroups, oldMenuAttributeGroups models.AttributeGroups) models.AttributeGroups {
	oldMenuAttributeGroupMap := make(map[string]models.AttributeGroup, len(oldMenuAttributeGroups))

	for _, attrGroup := range oldMenuAttributeGroups {
		oldMenuAttributeGroupMap[attrGroup.ExtID] = attrGroup
	}

	for i := range newMenuAttributeGroups {
		if _, ok := oldMenuAttributeGroupMap[newMenuAttributeGroups[i].ExtID]; !ok {
			continue
		}

		newMenuAttributeGroups[i].NamesByLanguage = oldMenuAttributeGroupMap[newMenuAttributeGroups[i].ExtID].NamesByLanguage

		saveDescriptionMap := make(map[string]models.LanguageDescription, len(newMenuAttributeGroups[i].Description))
		for _, description := range newMenuAttributeGroups[i].Description {
			saveDescriptionMap[description.LanguageCode] = description
		}

		newMenuAttributeGroups[i].Description = oldMenuAttributeGroupMap[newMenuAttributeGroups[i].Name].Description

		for j := range newMenuAttributeGroups[i].Description {
			if _, ok := saveDescriptionMap[newMenuAttributeGroups[i].Description[j].LanguageCode]; !ok {
				continue
			}
			newMenuAttributeGroups[i].Description[j].Value = saveDescriptionMap[newMenuAttributeGroups[i].Description[j].LanguageCode].Value
		}
	}

	return newMenuAttributeGroups
}

func (m *mnm) SaveAttributeLanguagesInMenu(newMenuAttributes, oldMenuAttributes models.Attributes) models.Attributes {
	oldMenuAttributeMap := make(map[string]models.Attribute, len(oldMenuAttributes))

	for _, attribute := range oldMenuAttributes {
		oldMenuAttributeMap[attribute.ExtID] = attribute
	}

	for i := range newMenuAttributes {
		if _, ok := oldMenuAttributeMap[newMenuAttributes[i].ExtID]; !ok {
			continue
		}

		newMenuAttributes[i].NamesByLanguage = oldMenuAttributeMap[newMenuAttributes[i].ExtID].NamesByLanguage

		saveDescriptionMap := make(map[string]models.LanguageDescription, len(newMenuAttributes[i].Description))
		for _, description := range newMenuAttributes[i].Description {
			saveDescriptionMap[description.LanguageCode] = description
		}

		newMenuAttributes[i].Description = oldMenuAttributeMap[newMenuAttributes[i].ExtID].Description

		for j := range newMenuAttributes[i].Description {
			if _, ok := saveDescriptionMap[newMenuAttributes[i].Description[j].LanguageCode]; !ok {
				continue
			}
			newMenuAttributes[i].Description[j].Value = saveDescriptionMap[newMenuAttributes[i].Description[j].LanguageCode].Value
		}
	}

	return newMenuAttributes
}

func (m *mnm) WotMenuSaveProductLanguagesAndProductInformation(newMenuProducts, oldMenuProducts models.Products) models.Products {

	newMenuProducts = m.SaveProductLanguagesInMenu(newMenuProducts, oldMenuProducts)
	newMenuProducts = m.SaveProductInformationInMenu(newMenuProducts, oldMenuProducts)

	return newMenuProducts
}

func (m *mnm) SaveProductInformationInMenu(newMenuProducts, oldMenuProducts models.Products) models.Products {
	oldMenuProductsMap := make(map[string]models.Product, len(oldMenuProducts))

	for _, product := range oldMenuProducts {
		oldMenuProductsMap[product.ExtID] = product
	}

	for i := range newMenuProducts {
		if _, ok := oldMenuProductsMap[newMenuProducts[i].ExtID]; !ok {
			continue
		}

		if oldMenuProductsMap[newMenuProducts[i].ExtID].ProductInformation.RegulatoryInformation != nil && len(oldMenuProductsMap[newMenuProducts[i].ExtID].ProductInformation.RegulatoryInformation) > 0 {
			newMenuProducts[i].ProductInformation = oldMenuProductsMap[newMenuProducts[i].ExtID].ProductInformation
		}
	}

	return newMenuProducts
}
