package stoplist

import (
	"context"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/rs/zerolog/log"
)

func (s *ServiceImpl) ActualizeStopListByStoreID2(ctx context.Context, storeID string) error {
	store, err := s.storeService.GetByID(ctx, storeID)
	if err != nil {
		return err
	}

	if err = s.actualizeStopList2(ctx, store); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) actualizeStopList2(ctx context.Context, store storeModels.Store) error {
	posService, err := s.posFactory.GetPosService(models.Pos(store.PosType), store)
	if err != nil {
		return err
	}

	stopListItems, err := posService.GetStopList(ctx)
	if err != nil {
		return err
	}

	posMenu, err := s.menuService.FindById(ctx, store.MenuID)
	if err != nil {
		return err
	}

	posProducts, err := s.getChangedProducts(stopListItems, *posMenu)
	if err != nil {
		return err
	}

	posAttributes, err := s.getChangedAttributes(stopListItems, *posMenu)
	if err != nil {
		return err
	}

	log.Info().Msgf("actualize stop list for store_id = %s, store_name = %s", store.ID, store.Name)
	if err := s.actualizeStopListPosMenu2(ctx, store, posProducts, posAttributes); err != nil {
		return err
	}

	for _, menu := range store.Menus {
		if err := s.actualizeStopListAggregatorMenu2(ctx, store, *posMenu, menu, posProducts, posAttributes, stopListItems); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) actualizeStopListAggregatorMenu2(ctx context.Context, store storeModels.Store, posMenu coreMenuModels.Menu, menu storeModels.StoreDSMenu, posMenuProducts []coreMenuModels.Product, posMenuAttributes []coreMenuModels.Attribute, posStopListItems coreMenuModels.StopListItems) error {
	if !menu.IsActive {
		return nil
	}

	aggregatorMenu, err := s.menuService.GetAggregatorMenuIfExists(ctx, store, menu.Delivery)
	if err != nil {
		return err
	}

	aggregatorProducts := s.menuService.MapPosProductsToAggregatorProducts(posMenu, aggregatorMenu, posMenuProducts)
	aggregatorProducts = s.setAvailabilityToAggregatorProductsFromPosProducts(aggregatorProducts, posMenuProducts)

	aggregatorAttributes := s.menuService.MapPosAttributesToAggregatorAttributes(posMenu, aggregatorMenu, posMenuAttributes)
	aggregatorAttributes = s.setAvailabilityToAggregatorAttributesFromPosAttributes(aggregatorAttributes, posMenuAttributes)

	menuID := menu.ID

	if err = s.updateStopListByProductIDsInDatabase(ctx, menuID, aggregatorProducts); err != nil {
		return err
	}
	if err = s.updateStopListByProductIDInAggregator(ctx, store, aggregatorProducts, menu.Delivery, posStopListItems); err != nil {
		return err
	}

	if err = s.updateStopListByAttributesInDatabase(ctx, menuID, aggregatorAttributes); err != nil {
		return err
	}
	if err = s.updateStopListByAttributesInAggregator(ctx, store, aggregatorAttributes, menu.Delivery, posStopListItems); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) setAvailabilityToAggregatorProductsFromPosProducts(aggregatorProducts []coreMenuModels.Product, posMenuProducts []coreMenuModels.Product) []coreMenuModels.Product {
	availabilityMap := s.toMapPosProductsAvailability(posMenuProducts)

	for i := range aggregatorProducts {
		aggregatorProductSystemID := s.menuService.GetSystemIDFromAggregatorProduct(aggregatorProducts[i])
		availability, ok := availabilityMap[aggregatorProductSystemID]
		if !ok {
			continue
		}
		aggregatorProducts[i].IsAvailable = availability
	}

	return aggregatorProducts
}
func (s *ServiceImpl) setAvailabilityToAggregatorAttributesFromPosAttributes(aggregatorAttributes []coreMenuModels.Attribute, posMenuAttributes []coreMenuModels.Attribute) []coreMenuModels.Attribute {
	availabilityMap := s.toMapPosAttributesAvailability(posMenuAttributes)

	for i := range aggregatorAttributes {
		aggregatorProductPosID := s.menuService.GetSystemIDFromAggregatorAttribute(aggregatorAttributes[i])
		availability, ok := availabilityMap[aggregatorProductPosID]
		if !ok {
			continue
		}
		aggregatorAttributes[i].IsAvailable = availability
	}

	return aggregatorAttributes
}

func (s *ServiceImpl) actualizeStopListPosMenu2(ctx context.Context, store storeModels.Store, products []coreMenuModels.Product, attributes []coreMenuModels.Attribute) error {

	menuID := store.MenuID
	//if err := s.menuService.UpdateStopList(ctx, menuID, stopListItems.Products()); err != nil {
	//	return nil, err
	//}

	if len(products) != 0 {
		if err := s.updateStopListByProductIDsInDatabase(ctx, menuID, products); err != nil {
			return err
		}
	}

	if len(attributes) != 0 {
		if err := s.updateStopListByAttributesInDatabase(ctx, menuID, attributes); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) getChangedProducts(stopListItems coreMenuModels.StopListItems, posMenu coreMenuModels.Menu) ([]coreMenuModels.Product, error) {
	stopListItemsMap := s.toMap(stopListItems)

	result := make([]coreMenuModels.Product, 0)

	for i := range posMenu.Products {
		product := posMenu.Products[i]
		posID := s.menuService.GetPosIDFromPosProduct(product)

		var targetAvailability bool
		if _, ok := stopListItemsMap[posID]; ok {
			targetAvailability = false
		} else {
			targetAvailability = true
		}

		currentAvailability := product.IsAvailable

		if targetAvailability == currentAvailability {
			continue
		}

		product.IsAvailable = targetAvailability
		result = append(result, product)
	}

	return result, nil
}

func (s *ServiceImpl) getChangedAttributes(stopListItems coreMenuModels.StopListItems, posMenu coreMenuModels.Menu) ([]coreMenuModels.Attribute, error) {
	stopListItemsMap := s.toMap(stopListItems)

	result := make([]coreMenuModels.Attribute, 0)

	for i := range posMenu.Attributes {
		attribute := posMenu.Attributes[i]
		posID := s.menuService.GetPosIDFromPosAttribute(attribute)

		var targetAvailability bool
		if _, ok := stopListItemsMap[posID]; ok {
			targetAvailability = false
		} else {
			targetAvailability = true
		}

		currentAvailability := attribute.IsAvailable

		if targetAvailability == currentAvailability {
			continue
		}

		attribute.IsAvailable = targetAvailability
		result = append(result, attribute)
	}

	return result, nil
}

func (s *ServiceImpl) toMap(stopListItems coreMenuModels.StopListItems) map[string]coreMenuModels.StopListItem {
	result := make(map[string]coreMenuModels.StopListItem)
	for i := range stopListItems {
		stopListItem := stopListItems[i]
		result[stopListItem.ProductID] = stopListItem
	}
	return result
}

func (s *ServiceImpl) toMapPosProductsAvailability(posProducts []coreMenuModels.Product) map[string]bool {
	result := make(map[string]bool)

	for i := range posProducts {
		posProduct := posProducts[i]
		systemID := s.menuService.GetSystemIDFromPosProduct(posProduct)
		result[systemID] = posProduct.IsAvailable
	}

	return result
}
func (s *ServiceImpl) toMapPosAttributesAvailability(posAttributes []coreMenuModels.Attribute) map[string]bool {
	result := make(map[string]bool)

	for i := range posAttributes {
		posAttribute := posAttributes[i]
		systemID := s.menuService.GetSystemIDFromPosAttribute(posAttribute)
		result[systemID] = posAttribute.IsAvailable
	}

	return result
}
