package managers

import (
	"context"
	"encoding/json"
	"fmt"
	models2 "github.com/kwaaka-team/orders-core/core/foodband/models"
	"github.com/kwaaka-team/orders-core/core/foodband/resources/http/v1/dto"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/domain/logger"
	"go.uber.org/zap"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	menuCoreModel "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	storeCore "github.com/kwaaka-team/orders-core/pkg/store"
	storeCoreModel "github.com/kwaaka-team/orders-core/pkg/store/dto"
)

type Menu interface {
	UploadMenu(ctx context.Context, req models2.UploadMenuReq) (string, error)
	AttributeStopList(ctx context.Context, req models2.StopListReq) error
	GetMenuUploadStatus(ctx context.Context, req models2.GetMenuUploadStatusReq) (menuCoreModel.MenuUploadTransaction, error)
}

type menuImplementation struct {
	menuCli  menuCore.Client
	storeCli storeCore.Client
	sv3      *s3.S3
	logger   *zap.SugaredLogger
}

func NewMenuManager(menuCli menuCore.Client, storeCli storeCore.Client, sess *session.Session, logger *zap.SugaredLogger) Menu {
	sv3 := s3.New(sess)
	return &menuImplementation{
		menuCli:  menuCli,
		storeCli: storeCli,
		sv3:      sv3,
		logger:   logger,
	}
}

// Finds active menu from store
func (man *menuImplementation) updateDsMenu(ctx context.Context, store coreStoreModels.Store, deliveryService string, newMenu coreMenuModels.Menu) (string, error) {

	for _, menu := range store.Menus {
		if menu.Delivery == deliveryService && menu.IsActive {
			newMenu.ID = menu.ID
			break
		}
	}

	if newMenu.ID == "" {
		return "", dto.ErrNotFoundActiveDsMenu
	}

	err := man.menuCli.UpdateMenu(ctx, newMenu)
	if err != nil {
		return "", err
	}
	return newMenu.ID, nil
}

// Find POS menu and update with new items
func (man *menuImplementation) UpdatePosMenu(ctx context.Context, store coreStoreModels.Store, newMenu coreMenuModels.Menu) error {
	//Find POS menu
	posMenu, err := man.menuCli.GetMenuByID(ctx, store.MenuID)
	if err != nil {
		return err
	}
	posAttributeMap := make(map[string]coreMenuModels.Attribute, len(posMenu.Attributes))
	posAttributeGroupMap := make(map[string]coreMenuModels.AttributeGroup, len(posMenu.AttributesGroups))
	posProductsMap := make(map[string]coreMenuModels.Product, len(posMenu.Products))
	posCollectionMap := make(map[string]coreMenuModels.MenuCollection)
	posSuperCollectionMap := make(map[string]coreMenuModels.MenuSuperCollection, len(posMenu.SuperCollections))
	posSectionMap := make(map[string]coreMenuModels.Section, len(posMenu.Sections))

	for _, v := range posMenu.Attributes {
		if v.PosID == "" {
			continue
		}
		posAttributeMap[v.PosID] = v
	}
	for _, v := range posMenu.AttributesGroups {
		if v.ExtID == "" {
			continue
		}
		posAttributeGroupMap[v.ExtID] = v
	}
	for _, v := range posMenu.Products {
		if v.PosID == "" {
			continue
		}
		posProductsMap[v.PosID] = v
	}
	for _, v := range posMenu.Sections {
		if v.ExtID == "" {
			continue
		}
		posSectionMap[v.ExtID] = v
	}
	for _, v := range posMenu.Collections {
		if v.ExtID == "" {
			continue
		}
		posCollectionMap[v.ExtID] = v
	}
	for _, v := range posMenu.SuperCollections {
		if v.ExtID == "" {
			continue
		}
		posSuperCollectionMap[v.ExtID] = v
	}

	for _, attribute := range newMenu.Attributes {
		if attribute.PosID == "" {
			continue
		}
		if _, ok := posAttributeMap[attribute.PosID]; !ok {
			posAttributeMap[attribute.PosID] = attribute
			posMenu.Attributes = append(posMenu.Attributes, attribute)
		}
	}

	for _, product := range newMenu.Products {
		if product.PosID == "" {
			continue
		}
		if _, ok := posProductsMap[product.PosID]; !ok {
			posProductsMap[product.PosID] = product
			posMenu.Products = append(posMenu.Products, product)
		}
	}

	for _, attributeGroup := range newMenu.AttributesGroups {
		if attributeGroup.ExtID == "" {
			continue
		}
		if _, ok := posAttributeGroupMap[attributeGroup.ExtID]; !ok {
			posAttributeGroupMap[attributeGroup.ExtID] = attributeGroup
			posMenu.AttributesGroups = append(posMenu.AttributesGroups, attributeGroup)
		}
	}

	for _, section := range newMenu.Sections {
		if section.ExtID == "" {
			continue
		}
		if _, ok := posSectionMap[section.ExtID]; !ok {
			posSectionMap[section.ExtID] = section
			posMenu.Sections = append(posMenu.Sections, section)
		}
	}

	for _, collection := range newMenu.Collections {
		if collection.ExtID == "" {
			continue
		}
		if _, ok := posCollectionMap[collection.ExtID]; !ok {
			posCollectionMap[collection.ExtID] = collection
			posMenu.Collections = append(posMenu.Collections, collection)
		}
	}

	for _, superCollection := range newMenu.SuperCollections {
		if superCollection.ExtID == "" {
			continue
		}
		if _, ok := posSuperCollectionMap[superCollection.ExtID]; !ok {
			posSuperCollectionMap[superCollection.ExtID] = superCollection
			posMenu.SuperCollections = append(posMenu.SuperCollections, superCollection)
		}
	}

	return man.menuCli.UpdateMenu(ctx, posMenu)
}

// Download and serialize menu by URL
func (man *menuImplementation) downloadMenu(ctx context.Context, menuURL string) (models2.Menu, []string) {
	var details []string

	response, err := http.Get(menuURL)
	if err != nil {
		details = append(details, fmt.Sprintf("Error downloading menu: %s", err.Error()))
		return models2.Menu{}, details
	}

	defer response.Body.Close()

	menu := models2.Menu{}

	if err := json.NewDecoder(response.Body).Decode(&menu); err != nil {
		details = append(details, fmt.Sprintf("Menu body is not valid: %s", err.Error()))
		return models2.Menu{}, details
	}

	return menu, details
}

// Create Menu upload transactions in DB
func (man *menuImplementation) createMUT(ctx context.Context, transaction menuCoreModel.MenuUploadTransaction) (string, error) {
	transactionID, err := man.menuCli.CreateMenuUploadTransaction(ctx, transaction)

	if err != nil {
		return "", err
	}

	return transactionID, nil
}

func (man *menuImplementation) isStoreIntegrated(store coreStoreModels.Store, deliveryService string) bool {
	switch deliveryService {
	case "glovo":
		return len(store.Glovo.StoreID) > 0
	case "chocofood":
		return len(store.Chocofood.StoreID) > 0
	case "wolt":
		return len(store.Wolt.StoreID) > 0
	case "qr_menu":
		return store.QRMenu.IsIntegrated
	}

	for _, external := range store.ExternalConfig {
		if external.Type == deliveryService {
			return len(external.StoreID) > 0
		}
	}

	return false
}

// Upload menu manager
// Returns menu upload transaction ID
func (man *menuImplementation) UploadMenu(ctx context.Context, req models2.UploadMenuReq) (string, error) {
	store, err := man.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{
		ExternalStoreID: req.StoreID,
		DeliveryService: "foodband",
	})
	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: err,
		})
		return "", err
	}

	// Check if store is integrated with given delivery service
	if !man.isStoreIntegrated(store, req.DeliveryService) {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: fmt.Sprintf("store is not integrated with %s", req.DeliveryService),
		})
		return "", dto.ErrStoreNotIntegated
	}

	// Initialize menu upload transaction
	transaction := menuCoreModel.MenuUploadTransaction{
		StoreID: store.ID,
		Service: req.DeliveryService,
		Status:  models2.PROCESSING,
		MenuURL: req.MenuURL,
	}

	// Download menu from menuURL
	newMenu, details := man.downloadMenu(ctx, req.MenuURL)
	if len(details) > 0 {
		transaction.Details = details
		transaction.Status = models2.FETCH_MENU_SERVER_ERROR

		tx, err := man.createMUT(ctx, transaction)
		if err != nil {
			man.logger.Error(logger.LoggerInfo{
				System:   "foodband response error",
				Response: []interface{}{err, transaction},
			})
			return "", err
		}
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: tx,
		})
		return tx, nil
	}

	man.logger.Info(logger.LoggerInfo{
		System:   "foodband downloaded menu",
		Response: newMenu,
	})

	// Validate downloaded menu and try to convert it to DB menu format
	menu, details := newMenu.Validate(req.DeliveryService, store)
	if len(details) > 0 {
		transaction.Details = details
		transaction.Status = models2.FETCH_MENU_INVALID_PAYLOAD
		tx, err := man.createMUT(ctx, transaction)
		if err != nil {
			man.logger.Error(logger.LoggerInfo{
				System:   "foodband response error",
				Response: []interface{}{err, transaction},
			})
			return "", err
		}
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: tx,
		})
		return tx, nil
	}

	man.logger.Info(logger.LoggerInfo{
		System:   "foodband parsed menu",
		Response: menu,
	})

	if err := man.UpdatePosMenu(ctx, store, menu); err != nil {
		transaction.Details = append(transaction.Details, err.Error())
		transaction.Status = models2.FETCH_MENU_INVALID_PAYLOAD
		tx, err := man.createMUT(ctx, transaction)
		if err != nil {
			man.logger.Error(logger.LoggerInfo{
				System:   "foodband response error",
				Response: []interface{}{err, transaction},
			})
			return "", err
		}
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: tx,
		})
		return tx, nil
	}

	aggregatorMenuID, err := man.updateDsMenu(ctx, store, req.DeliveryService, menu)
	if err != nil {
		transaction.Details = append(transaction.Details, err.Error())
		transaction.Status = models2.FETCH_MENU_INVALID_PAYLOAD
		tx, err := man.createMUT(ctx, transaction)
		if err != nil {
			man.logger.Error(logger.LoggerInfo{
				System:   "foodband response error",
				Response: []interface{}{err, transaction},
			})
			return "", err
		}
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: tx,
		})
		return tx, nil
	}

	trId, err := man.menuCli.UploadMenu(ctx, menuCoreModel.MenuUploadRequest{
		StoreId:      store.ID,
		MenuId:       aggregatorMenuID,
		DeliveryName: req.DeliveryService,
		Sv3:          man.sv3,
	})
	if err != nil {
		if trId == "" {
			man.logger.Error(logger.LoggerInfo{
				System:   "foodband response error",
				Response: err,
			})
			return "", err
		}
		man.logger.Info(logger.LoggerInfo{
			System:   "foodband response",
			Response: trId,
		})
		return trId, nil
	}
	man.logger.Info(logger.LoggerInfo{
		System:   "foodband response",
		Response: trId,
	})
	return trId, nil
}

func (man *menuImplementation) AttributeStopList(ctx context.Context, req models2.StopListReq) error {
	if err := validateStopListInput(req); err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: err,
		})
		return err
	}

	err := man.menuCli.AttributesStopList(ctx, req.RestaurantID, []menuCoreModel.StopListItem{
		{
			ID:          req.ID,
			Price:       req.Price,
			IsAvailable: req.IsAvailable,
		},
	}, "")

	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: err,
		})
		return err
	}

	man.logger.Info(logger.LoggerInfo{
		System:   "foodband response",
		Response: fmt.Sprintf("attribute stoplist success, productID %s", req.ID),
	})
	return nil
}

func (man *menuImplementation) GetMenuUploadStatus(ctx context.Context, req models2.GetMenuUploadStatusReq) (menuCoreModel.MenuUploadTransaction, error) {
	if req.StoreID == "" || req.DeliveryService == "" || req.TransactionID == "" {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: fmt.Sprintf("invalid get menu upload status input, storeID %s, ds %s, trID %s", req.StoreID, req.DeliveryService, req.TransactionID),
		})
		return menuCoreModel.MenuUploadTransaction{}, fmt.Errorf("invalid incoming params")
	}

	store, err := man.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{
		ExternalStoreID: req.StoreID,
		DeliveryService: "foodband",
	})
	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: err,
		})
		return menuCoreModel.MenuUploadTransaction{}, dto.ErrStoreNotFound
	}

	res, err := man.menuCli.GetMenuUploadTransaction(ctx, menuCoreModel.MenuUploadTransaction{
		ID:      req.TransactionID,
		Service: req.DeliveryService,
		StoreID: store.ID,
	})
	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: err,
		})
		return menuCoreModel.MenuUploadTransaction{}, err
	}

	man.logger.Info(logger.LoggerInfo{
		System:   "foodband response",
		Response: res,
	})
	return res, nil
}

func validateStopListInput(req models2.StopListReq) error {
	if req.RestaurantID == "" {
		return fmt.Errorf("invalid restaurant id")
	}
	if req.ID == "" {
		return fmt.Errorf("invalid item id")
	}

	return nil
}
