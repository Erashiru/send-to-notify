package stoplist

import (
	"context"
	"fmt"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
)

func (s *ServiceImpl) GetStopListByDeliveryService(ctx context.Context, externalStoreId, deliveryService string,
	storeSecret string) ([]menuModels.Product, []menuModels.Attribute, error) {

	store, err := s.storeService.GetByExternalIdAndDeliveryService(ctx, externalStoreId, deliveryService)
	if err != nil {
		return nil, nil, err
	}

	if err = s.checkSecretAndMenu(store, deliveryService, storeSecret); err != nil {
		return nil, nil, err
	}

	aggMenu, err := s.menuService.GetAggregatorMenuIfExists(ctx, store, deliveryService)
	if err != nil {
		return nil, nil, err
	}

	stopListProducts := make([]menuModels.Product, 0)
	for i := range aggMenu.Products {
		product := aggMenu.Products[i]
		if !s.isProductOnStop(product) {
			continue
		}
		stopListProducts = append(stopListProducts, product)
	}

	stopListAttributes := make([]menuModels.Attribute, 0)
	for i := range aggMenu.Attributes {
		attribute := aggMenu.Attributes[i]
		if !s.isAttributeOnStop(attribute) {
			continue
		}
		stopListAttributes = append(stopListAttributes, attribute)
	}

	return stopListProducts, stopListAttributes, nil
}

func (s *ServiceImpl) isProductOnStop(product menuModels.Product) bool {
	return !product.IsAvailable
}

func (s *ServiceImpl) isAttributeOnStop(attribute menuModels.Attribute) bool {
	return !attribute.IsAvailable
}

func (s *ServiceImpl) checkSecretAndMenu(store storeModels.Store, deliveryService, storeSecret string) error {
	isValid, err := s.storeService.IsSecretValid(store, deliveryService, storeSecret)
	if err != nil {
		return err
	}
	if !isValid {
		return fmt.Errorf("store with name %s has invalid secret", store.Name)
	}

	isMenuExists := s.menuService.IsMenuExists(store, deliveryService)
	if !isMenuExists {
		return fmt.Errorf("store with name %s doesn't has %s active menu", store.Name, deliveryService)
	}

	return nil
}

func (s *ServiceImpl) GetRetailRemains(ctx context.Context, externalStoreId, deliveryService string,
	storeSecret string) (menuModels.StopListProducts, menuModels.StopListAttributes, error) {

	store, err := s.storeService.GetByExternalIdAndDeliveryService(ctx, externalStoreId, deliveryService)
	if err != nil {
		return nil, nil, err
	}

	isValid, err := s.storeService.IsSecretValid(store, deliveryService, storeSecret)
	if err != nil {
		return nil, nil, err
	}
	if !isValid {
		return nil, nil, fmt.Errorf("store with name %s has invalid secret", store.Name)
	}

	stopListItems, err := s.getStoplistItems(ctx, store)
	if err != nil {
		return nil, nil, err
	}

	stopListProducts := make(menuModels.StopListProducts, 0)
	for _, item := range stopListItems {
		if item.Balance <= 0 {
			continue
		}
		stoplistProduct := menuModels.StopListProduct{
			ExtID: item.ProductID,
			Stock: fmt.Sprintf("%f", item.Balance),
		}
		stopListProducts = append(stopListProducts, stoplistProduct)
	}

	// TODO: implement remains for attributes too

	return stopListProducts, make(menuModels.StopListAttributes, 0), nil
}

func (s *ServiceImpl) getStoplistItems(ctx context.Context, store storeModels.Store) (menuModels.StopListItems, error) {
	posService, err := s.posFactory.GetPosService(models.Pos(store.PosType), store)
	if err != nil {
		return nil, err
	}
	stopListItems, err := posService.GetStopList(ctx)
	if err != nil {
		return nil, err
	}

	return stopListItems, nil
}
