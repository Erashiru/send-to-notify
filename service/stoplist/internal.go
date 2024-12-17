package stoplist

import (
	"context"
	"fmt"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
)

func (s *ServiceImpl) stopProductPosMenu(ctx context.Context, isAvailable bool, posMenu *menuModels.Menu, posProducts []menuModels.Product) error {
	if len(posProducts) == 0 {
		return nil
	}

	for i := range posProducts {
		posProducts[i].IsAvailable = isAvailable
	}

	if err := s.updateStopListByProductIDsInDatabase(ctx, posMenu.ID, posProducts); err != nil {
		return err
	}

	return nil

}

func (s *ServiceImpl) getPosMenu(ctx context.Context, store coreStoreModels.Store) (*menuModels.Menu, error) {
	menu, err := s.menuService.FindById(ctx, store.MenuID)
	if err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *ServiceImpl) stopProductsAggregatorMenu(ctx context.Context, isAvailable bool, store coreStoreModels.Store, posMenu *menuModels.Menu, posProducts []menuModels.Product, deliveryService string) error {
	isMenuExists := s.menuService.IsMenuExists(store, deliveryService)
	if !isMenuExists {
		return fmt.Errorf("menu %s for delivery service %s is not exists", store.ID, deliveryService)
	}
	aggregatorMenu, err := s.menuService.GetAggregatorMenuIfExists(ctx, store, deliveryService)
	if err != nil {
		return err
	}

	aggProducts := s.menuService.MapPosProductsToAggregatorProducts(*posMenu, aggregatorMenu, posProducts)

	aggProducts, err = s.filterNotDeletedOnlyAggregatorProducts(aggregatorMenu, aggProducts)
	if err != nil {
		return err
	}

	aggProducts = s.filterProducts(aggProducts)

	if len(aggProducts) == 0 {
		return nil
	}

	for i := range aggProducts {
		aggProducts[i].IsAvailable = isAvailable
	}

	if err = s.updateStopListByProductIDInAggregator(ctx, store, aggProducts, deliveryService, nil); err != nil {
		return err
	}

	if err = s.updateStopListByProductIDsInDatabase(ctx, aggregatorMenu.ID, aggProducts); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) filterNotDeletedOnlyPosProducts(posMenu menuModels.Menu, posProducts []menuModels.Product) ([]menuModels.Product, error) {
	return s.getNotDeletedProducts(posMenu, posProducts, s.menuService.GetSystemIDFromPosProduct)
}

func (s *ServiceImpl) filterNotDeletedOnlyAggregatorProducts(aggMenu menuModels.Menu, aggProducts []menuModels.Product) ([]menuModels.Product, error) {
	return s.getNotDeletedProducts(aggMenu, aggProducts, s.menuService.GetSystemIDFromAggregatorProduct)
}

func (s *ServiceImpl) getNotDeletedProducts(menu menuModels.Menu, products []menuModels.Product, productSystemIDExtractFunc func(aggProduct menuModels.Product) string) ([]menuModels.Product, error) {
	targetProducts := make(map[string]interface{})
	for _, product := range products {
		systemID := productSystemIDExtractFunc(product)
		targetProducts[systemID] = struct{}{}
	}

	var notDeletedOnly = make([]menuModels.Product, 0)
	for _, product := range menu.Products {
		if product.IsDeleted {
			continue
		}
		systemID := productSystemIDExtractFunc(product)
		if _, ok := targetProducts[systemID]; !ok {
			continue
		}
		notDeletedOnly = append(notDeletedOnly, product)
	}

	return notDeletedOnly, nil
}

func (s *ServiceImpl) getPosProductsByIDs(posMenu menuModels.Menu, posProductIDs []string) ([]menuModels.Product, error) {

	targetProducts := make(map[string]interface{})
	for _, productID := range posProductIDs {
		targetProducts[productID] = struct{}{}
	}

	var result = make([]menuModels.Product, 0)
	for _, product := range posMenu.Products {
		posId := s.menuService.GetPosIDFromPosProduct(product)
		if _, ok := targetProducts[posId]; !ok {
			continue
		}
		result = append(result, product)
	}

	return result, nil
}

func (s *ServiceImpl) disabledByValidationProductsInPosMenu(ctx context.Context, posMenu *menuModels.Menu, posProducts []menuModels.Product, disabledByValidation bool) error {
	if len(posProducts) == 0 {
		return nil
	}

	for i := range posProducts {
		posProducts[i].DisabledByValidation = disabledByValidation
	}

	if err := s.updateProductsDisabledByValidationInDatabase(ctx, posMenu.ID, posProducts); err != nil {
		return err
	}

	return nil
}
