package stoplist

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuService "github.com/kwaaka-team/orders-core/service/menu"
	"github.com/rs/zerolog/log"
)

func (s *ServiceImpl) UpdateStopListByPosProductID(ctx context.Context, isAvailable bool, storeID string, posProductID string) error {
	store, err := s.storeService.GetByID(ctx, storeID)
	if err != nil {
		return err
	}

	posProductIDs := []string{posProductID}

	posMenu, err := s.getPosMenu(ctx, store)
	if err != nil {
		return err
	}

	posProducts, err := s.getPosProductsByIDs(*posMenu, posProductIDs)
	if err != nil {
		return err
	}

	posProducts, err = s.filterNotDeletedOnlyPosProducts(*posMenu, posProducts)
	if err != nil {
		return err
	}

	posProducts = s.filterProducts(posProducts)

	if len(posProducts) == 0 {
		log.Info().Msgf("stoplist servie - pos products for updating are not found")
		return nil
	}

	if err = s.stopProductPosMenu(ctx, isAvailable, posMenu, posProducts); err != nil {
		return err
	}

	for i := range store.Menus {
		aggMenu := store.Menus[i]
		if err = s.stopProductsAggregatorMenu(ctx, isAvailable, store, posMenu, posProducts, aggMenu.Delivery); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) UpdateStopListForValidateStoreMenus(ctx context.Context, storeID, deliveryService string, productDetails []menuService.ProductDetail) error {
	log.Info().Msgf("update stoplist for validate store menus, store id: %s, delivery service: %s", storeID, deliveryService)

	productIDsForLock := make([]string, 0, len(productDetails))
	for _, p := range productDetails {
		productIDsForLock = append(productIDsForLock, p.ID)
	}

	store, err := s.storeService.GetByID(ctx, storeID)
	if err != nil {
		return err
	}

	aggregatorMenu, err := s.menuService.GetAggregatorMenuIfExists(ctx, store, deliveryService)
	if err != nil {
		return err
	}

	reportM := make(map[string]struct{}, len(productDetails))
	for _, p := range productDetails {
		reportM[p.ID] = struct{}{}
	}

	// remove from the stoplist products that are disabledByValidation: true, but they are not in the report
	var productIDsForUnlock []string
	for _, product := range aggregatorMenu.Products {
		if _, ok := reportM[product.ExtID]; !ok && product.DisabledByValidation {
			productIDsForUnlock = append(productIDsForUnlock, product.ExtID)
		}
	}

	if err := s.updateProductsInPosMenuForValidateStoreMenus(ctx, store, productIDsForLock, productIDsForUnlock); err != nil {
		return err
	}

	if err := s.updateProductInAggrForValidateStoreMenus(ctx, store, deliveryService, productIDsForLock, productIDsForUnlock); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) updateProductsInPosMenuForValidateStoreMenus(ctx context.Context, store coreStoreModels.Store, productIDsForLock, productIDsForUnlock []string) error {
	posMenu, err := s.getPosMenu(ctx, store)
	if err != nil {
		return err
	}

	posProductsLock, err := s.filterPosProducts(*posMenu, productIDsForLock)
	if err != nil {
		return err
	}

	posProductsUnlock, err := s.filterPosProducts(*posMenu, productIDsForUnlock)
	if err != nil {
		return err
	}

	if len(posProductsLock) == 0 && len(posProductsUnlock) == 0 {
		log.Info().Msgf("stoplist servie - pos products for updating are not found")
		return nil
	}

	if err := s.disabledByValidationProductsInPosMenu(ctx, posMenu, posProductsLock, true); err != nil {
		return err
	}

	if err := s.disabledByValidationProductsInPosMenu(ctx, posMenu, posProductsUnlock, false); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) filterPosProducts(posMenu models.Menu, productIDs []string) ([]models.Product, error) {
	posProducts, err := s.getPosProductsByIDs(posMenu, productIDs)
	if err != nil {
		return nil, err
	}

	posProducts, err = s.filterNotDeletedOnlyPosProducts(posMenu, posProducts)
	if err != nil {
		return nil, err
	}

	posProducts = s.filterProducts(posProducts)

	return posProducts, nil
}

func (s *ServiceImpl) updateProductInAggrForValidateStoreMenus(ctx context.Context, store coreStoreModels.Store, deliveryService string, productIDsForLock, productIDsForUnlock []string) error {
	isMenuExists := s.menuService.IsMenuExists(store, deliveryService)
	if !isMenuExists {
		return fmt.Errorf("menu for delivery service %s in store %s is not exists", deliveryService, store.ID)
	}

	aggregatorMenu, err := s.menuService.GetAggregatorMenuIfExists(ctx, store, deliveryService)
	if err != nil {
		return err
	}

	aggrProductsLock, err := s.filterAggrProducts(aggregatorMenu, productIDsForLock)
	if err != nil {
		return err
	}

	aggrProductsUnlock, err := s.filterAggrProducts(aggregatorMenu, productIDsForUnlock)
	if err != nil {
		return err
	}

	if len(aggrProductsLock) == 0 && len(aggrProductsUnlock) == 0 {
		log.Info().Msgf("stoplist servie - aggregator products for updating are not found")
		return nil
	}

	if err := s.updateProductsInAggrByStatusForValidateStoreMenus(ctx, store, aggregatorMenu.ID, deliveryService, aggrProductsLock, true); err != nil {
		return err
	}

	if err := s.updateProductsInAggrByStatusForValidateStoreMenus(ctx, store, aggregatorMenu.ID, deliveryService, aggrProductsUnlock, false); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) filterAggrProducts(aggrMenu models.Menu, productIDs []string) ([]models.Product, error) {
	aggProducts := s.menuService.GetAggrMenuProductsByIDs(aggrMenu, productIDs)

	aggProducts, err := s.filterNotDeletedOnlyAggregatorProducts(aggrMenu, aggProducts)
	if err != nil {
		return nil, err
	}

	aggProducts = s.filterProducts(aggProducts)

	return aggProducts, nil
}

func (s *ServiceImpl) updateProductsInAggrByStatusForValidateStoreMenus(ctx context.Context, store coreStoreModels.Store, menuID, deliveryService string, aggrProducts []models.Product, disabledByValidation bool) error {
	for i := range aggrProducts {
		aggrProducts[i].DisabledByValidation = disabledByValidation
	}

	aggUpdateProducts := s.productsUpdateInAggrForValidateStoreMenus(aggrProducts, disabledByValidation)

	log.Info().Msgf("update aggregator products in aggregator: %+v", aggUpdateProducts)

	if err := s.updateStopListByProductIDInAggregator(ctx, store, aggUpdateProducts, deliveryService, nil); err != nil {
		return err
	}

	// not save isAvailable status in DB, only disabledByValidation
	if err := s.updateProductsDisabledByValidationInDatabase(ctx, menuID, aggrProducts); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) productsUpdateInAggrForValidateStoreMenus(products []models.Product, disabledByValidation bool) []models.Product {
	if len(products) == 0 {
		return products
	}

	switch disabledByValidation {
	case false:
		res := make([]models.Product, 0, len(products))
		for i := range products {
			if products[i].IsDisabled {
				continue
			}
			products[i].IsAvailable = true
			res = append(res, products[i])
		}

		return res

	default:
		for i := range products {
			products[i].IsAvailable = false
		}

		return products
	}
}
