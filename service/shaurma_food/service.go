package shaurma_food

import (
	"context"
	"fmt"
	customeErrors "github.com/kwaaka-team/orders-core/core/errors"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	posService "github.com/kwaaka-team/orders-core/service/pos"
	storeServicePkg "github.com/kwaaka-team/orders-core/service/store"
	"github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type Service struct {
	storeService      storeServicePkg.Service
	menuService       *menuServicePkg.Service
	posService        posService.Factory
	storeGroupService storegroup.Service
}

func NewService(storeService storeServicePkg.Service,
	menuService *menuServicePkg.Service,
	posService posService.Factory,
	storeGroupService storegroup.Service) (*Service, error) {

	return &Service{
		storeService:      storeService,
		menuService:       menuService,
		posService:        posService,
		storeGroupService: storeGroupService,
	}, nil
}

func (s *Service) SetAggregatorsMenuPricesFromExternalMenu(ctx context.Context, restaurantID, apiKey, organizationID, terminalID, externalMenuID, priceCategory string, ignoreExternalMenuProductsWithZeroNullPrice bool) error {
	shaurmaFoodRestiIDs, err := s.getShaurmaFoodRestaurantIDs(ctx)
	if err != nil {
		return err
	}
	if !s.isShaurmaFood(restaurantID, shaurmaFoodRestiIDs) {
		return errors.New("restaurant id is not for ShaurmaFood")
	}

	store, err := s.storeService.GetByID(ctx, restaurantID)
	if err != nil {
		return err
	}

	store.IikoCloud = models.StoreIikoConfig{
		Key:            apiKey,
		OrganizationID: organizationID,
		TerminalID:     terminalID,
		ExternalMenuID: externalMenuID,
		PriceCategory:  priceCategory,
		IgnoreExternalMenuProductsWithZeroNullPrice: ignoreExternalMenuProductsWithZeroNullPrice,
		IsExternalMenu: true,
	}

	externalPosMenu, err := s.getExternalMenu(ctx, store)
	if err != nil {
		return err
	}

	posProductsMap := s.getProductsPriceMap(externalPosMenu)
	posAttributesMap := s.getAttributesPriceMap(externalPosMenu)

	for _, storeMenu := range store.Menus {
		if !s.validMenuDeliveryAndStatus(storeMenu) {
			continue
		}

		aggrMenu, err := s.menuService.FindById(ctx, storeMenu.ID)
		if err != nil {
			return errors.Wrapf(err, "menuID: %s", storeMenu.ID)
		}

		menuWithNewPrices := s.setPosPricesToAggregatorMenu(posProductsMap, posAttributesMap, externalPosMenu, *aggrMenu)

		if err := s.menuService.UpdateMenuEntities(ctx, menuWithNewPrices.ID, menuWithNewPrices); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) validMenuDeliveryAndStatus(menu models.StoreDSMenu) bool {
	if !menu.IsActive || menu.IsDeleted {
		return false
	}

	return true
}

func (s *Service) getExternalMenu(ctx context.Context, store models.Store) (models2.Menu, error) {
	posService, err := s.posService.GetPosService(coreModels.Pos(store.PosType), store)
	if err != nil {
		return models2.Menu{}, err
	}

	externalMenu, err := posService.GetMenu(ctx, store, models2.Menu{})
	if err != nil {
		return models2.Menu{}, err
	}

	return externalMenu, nil
}

func (s *Service) setPosPricesToAggregatorMenu(posProductsMap, posAttributesMap map[string]float64, extPosMenu, aggrMenu models2.Menu) models2.Menu {

	switch aggrMenu.Delivery {
	case string(models2.YANDEX):
		posDefaultsPriceAttachedToProductID := s.getDefaultsPriceForProduct(extPosMenu)
		for _, product := range aggrMenu.Products {
			if defaultsSum, ok := posDefaultsPriceAttachedToProductID[product.ExtID]; ok {
				product.Price[0].Value += defaultsSum
			}
		}

	default:
		for i := range aggrMenu.Products {
			id := aggrMenu.Products[i].ExtID
			if aggrMenu.Products[i].PosID != "" {
				id = aggrMenu.Products[i].PosID
			}

			price, ok := posProductsMap[id]
			if !ok {
				continue
			}

			if price == 0 {
				continue
			}

			if len(aggrMenu.Products[i].Price) == 0 {
				continue
			}

			aggrMenu.Products[i].Price[0].Value = price
		}
	}

	for i := range aggrMenu.Attributes {
		price, ok := posAttributesMap[aggrMenu.Attributes[i].ExtID]
		if !ok {
			continue
		}
		aggrMenu.Attributes[i].Price = price
	}

	return aggrMenu
}

func (s *Service) getDefaultsPriceForProduct(menu models2.Menu) map[string]float64 {

	tableware := map[string]float64{
		"b4d6052f-426e-4ae1-9794-d417300ad097": 220,
		"5348df90-25ed-401b-b7df-43291f12c619": 220,
		"1c9416d5-39d5-44fc-b201-14bf170b3c67": 260,
		"b73a79b0-9ccc-45a6-b42d-83a17002f01c": 220,
		"f5e4569b-bd43-4292-a5e2-34d611b67bef": 220,
		"5c3bde4a-74e5-436b-8441-f1183f3730ca": 220,
		"bcd31b5d-4e86-4704-aedc-047ab8022fb4": 220,
		"0740a8b6-d478-47eb-afe0-b79a0dedccb3": 220,
		"e0132778-c8dd-4532-821f-875dc5e317d0": 220,
		"e8af2509-7c7d-452f-b372-699cbb4571d4": 940,
		"63ceca04-f39b-4f8f-a342-25e46748b2b8": 1320,
		"cb3abb55-2da9-48b8-a4e6-84761984ab18": 940,
		"2a01923e-3843-48b7-819e-87c620de9924": 550,
		"656dd1ab-59dc-418b-ae9b-0d214fa54895": 220,
		"7104f56e-4dbd-446b-9820-8ec9a4376d9f": 220,
		"cbd5d847-a915-4754-b16d-855d82159311": 220,
		"439d3b44-00fa-4793-813d-9f9968848228": 330,
		"33d92641-8dfd-4506-ab31-c61527cd2b5b": 550,
	}

	defaultsForProductSum := make(map[string]float64)

	for _, product := range menu.Products {
		if product.Attributes == nil {
			continue
		}
		var price float64
		for _, attr := range product.Attributes {
			if defPrice, ok := tableware[attr]; ok {
				price += defPrice
				defaultsForProductSum[product.ExtID] = price
			}
		}
	}
	return defaultsForProductSum
}

func (s *Service) getProductsPriceMap(menu models2.Menu) map[string]float64 {
	res := make(map[string]float64)

	for i := range menu.Products {
		res[menu.Products[i].ExtID] = menu.Products[i].Price[0].Value
	}

	return res
}

func (s *Service) getAttributesPriceMap(menu models2.Menu) map[string]float64 {
	res := make(map[string]float64)

	for i := range menu.Attributes {
		res[menu.Attributes[i].ExtID] = menu.Attributes[i].Price
	}

	return res
}

func (s *Service) getShaurmaFoodRestaurantIDs(ctx context.Context) ([]string, error) {
	shaurmaFoodGroupIDs := []string{"642ab068d5ad369ab4647d44", "6604030be5cb51f5698fafd5"}
	shaurmaFoodRestaurantIDs := make([]string, 0, 1)

	for _, id := range shaurmaFoodGroupIDs {
		group, err := s.storeGroupService.GetStoreGroupByID(ctx, id)
		if err != nil {
			continue
		}
		shaurmaFoodRestaurantIDs = append(shaurmaFoodRestaurantIDs, group.StoreIds...)
	}

	return shaurmaFoodRestaurantIDs, nil
}

func (s *Service) isShaurmaFood(restiID string, shaurmaFoodRestiIDs []string) bool {
	for _, shaurmaFoodRestiID := range shaurmaFoodRestiIDs {
		if restiID == shaurmaFoodRestiID {
			return true
		}
	}
	return false
}

func (s *Service) UpdateProductsFromMainMenu(ctx context.Context, restId, aggregator string) error {
	shaurmaFoodRestiIDs, err := s.getShaurmaFoodRestaurantIDs(ctx)
	if err != nil {
		return err
	}
	if !s.isShaurmaFood(restId, shaurmaFoodRestiIDs) {
		return errors.New("restaurant id is not for ShaurmaFood")
	}

	rest, err := s.storeService.GetByID(ctx, restId)
	if err != nil {
		return err
	}

	var multiError customeErrors.Error

	var mainMenuID string

	for _, menu := range rest.Menus {
		if menu.Delivery == aggregator && menu.IsActive {
			mainMenuID = menu.ID
			break
		}
	}

	mainMenu, err := s.menuService.FindById(ctx, mainMenuID)
	if err != nil {
		return err
	}

	mainMenuProducts := make(map[string]models2.Product)
	mainMenuAttributes := make(map[string]models2.Attribute)

	for _, product := range mainMenu.Products {
		if product.PosID == "" {
			continue
		}
		mainMenuProducts[product.PosID] = product
	}

	for _, attribute := range mainMenu.Attributes {
		if attribute.ExtID == "" {
			continue
		}
		mainMenuAttributes[attribute.ExtID] = attribute
	}

	otherStores, err := s.storeService.GetStoresByStoreGroupID(ctx, rest.RestaurantGroupID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	multiErrCh := make(chan error, len(otherStores))

	for _, store := range otherStores {
		if restId == store.ID {
			continue
		}

		wg.Add(1)
		go func(st models.Store) {
			defer wg.Done()

			var menuForUpdateID string
			for _, otherMenu := range st.Menus {

				if otherMenu.Delivery == aggregator && otherMenu.IsActive {
					menuForUpdateID = otherMenu.ID
					break
				}
			}

			if menuForUpdateID == "" {
				return
			}

			menuForUpdate, err := s.menuService.FindById(ctx, menuForUpdateID)

			if err != nil {
				multiErrCh <- errors.Wrap(err, fmt.Sprintf("error with this menu id: %s, in this restaurant %s", menuForUpdateID, st.ID))
				return
			}

			var req []models2.UpdateProductImageAndDescription
			var attributeReq []models2.UpdateAttributePrice

			for _, productForUpdate := range menuForUpdate.Products {
				mainProduct, ok := mainMenuProducts[productForUpdate.PosID]
				if !ok {
					continue
				}

				var updateProduct models2.UpdateProductImageAndDescription
				updateProduct.Description = []models2.LanguageDescription{}
				updateProduct.PosID = mainProduct.PosID

				if mainProduct.ImageURLs != nil {
					updateProduct.ImageURLs = mainProduct.ImageURLs
				}

				if mainProduct.Description != nil && len(mainProduct.Description) != 0 {
					if mainProduct.Description[0].Value != "" {
						updateProduct.Description = append(updateProduct.Description, mainProduct.Description[0])
					}
				}

				if mainProduct.Weight != 0 {
					updateProduct.Weight = mainProduct.Weight
					updateProduct.MeasureUnit = mainProduct.MeasureUnit
				}

				updateProduct.Price = mainProduct.Price

				req = append(req, models2.UpdateProductImageAndDescription{
					PosID:       updateProduct.PosID,
					ImageURLs:   updateProduct.ImageURLs,
					Description: updateProduct.Description,
					Weight:      updateProduct.Weight,
					MeasureUnit: updateProduct.MeasureUnit,
					Price:       updateProduct.Price,
				})
			}

			if err = s.menuService.UpdateProductsImageAndDescription(ctx, menuForUpdateID, req); err != nil {
				multiErrCh <- errors.Wrap(err, fmt.Sprintf("error with this menu id: %s, in this restaurant %s", menuForUpdateID, st.ID))
			}

			for _, attributeForUpdate := range menuForUpdate.Attributes {
				mainAttribute, ok := mainMenuAttributes[attributeForUpdate.ExtID]
				if !ok {
					continue
				}
				var updateAttribute models2.UpdateAttributePrice
				updateAttribute.ExtID = mainAttribute.ExtID
				updateAttribute.Price = mainAttribute.Price

				attributeReq = append(attributeReq, updateAttribute)
			}

			if err = s.menuService.UpdateAttributesPrice(ctx, menuForUpdateID, attributeReq); err != nil {
				multiErrCh <- errors.Wrap(err, fmt.Sprintf("error with this menu id: %s, in this restaurant %s", menuForUpdateID, st.ID))
			}

			time.Sleep(1 * time.Second)

		}(store)
	}

	wg.Wait()
	close(multiErrCh)

	for err := range multiErrCh {
		multiError.Append(err)
	}

	if multiError.ErrorOrNil() != nil {
		return multiError
	}

	return nil
}
