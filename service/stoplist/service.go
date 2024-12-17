package stoplist

import (
	"context"

	"github.com/kwaaka-team/orders-core/core/config"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/rs/zerolog/log"
)

type Service interface {
	UpdateStopListByPosProductID(ctx context.Context, isAvailable bool, storeID string, productID string) error
	UpdateStopListByAttributeID(ctx context.Context, isAvailable bool, storeID string, attributeID string) error
	UpdateStopListBySectionID(ctx context.Context, isAvailable bool, storeGroupID string, deliveryToSectionIDs map[string][]string) error
	ActualizeStopListByStoreID(ctx context.Context, storeID string) error
	ActualizeStopListByToken(ctx context.Context, token string) error
	ActualizeStopListByPosType(ctx context.Context, posType string) error
	ActualizeStopListByPosTypes(ctx context.Context, posTypes []string) error
	GetStopListByDeliveryService(ctx context.Context, externalStoreId, deliveryService string, storeSecret string) ([]menuModels.Product, []menuModels.Attribute, error)
	ActualizeStoplistbyYarosStoreID(ctx context.Context, storeID string) error
	GetRetailRemains(ctx context.Context, externalStoreId, deliveryService string, storeSecret string) (menuModels.StopListProducts, menuModels.StopListAttributes, error)
	UpdateStopListForValidateStoreMenus(ctx context.Context, storeID, deliveryService string, productDetails []menu.ProductDetail) error
	AddYandexTransaction(ctx context.Context, storeID, deliveryService string, products menuModels.StopListProducts, attributes menuModels.StopListAttributes) error
}

type ServiceImpl struct {
	storeService      store.Service
	storeGroupService storegroup.Service
	menuService       *menu.Service
	aggregatorFactory aggregator.Factory
	posFactory        pos.Factory
	repo              Repository
	concurrencyLevel  int
	woltCfg           config.WoltConfiguration
	notifyCli         que.SQSInterface
	stopListType
}

func (s *ServiceImpl) updateStopListByProductIDInAggregator(ctx context.Context, store storeModels.Store,
	products []menuModels.Product, deliveryService string, posStopListItems menuModels.StopListItems) error {
	if len(products) == 0 {
		return nil
	}

	externalStoreIDs, err := s.storeService.GetStoreExternalIds(store, deliveryService)
	if err != nil {
		return err
	}

	aggregatorService, err := s.aggregatorFactory.GetAggregator(deliveryService, store)
	if err != nil {
		return err
	}

	slp := s.toStopListProducts(products)
	transaction := menuModels.StopListTransaction{
		StoreID:          store.ID,
		Products:         slp,
		PosStopListItems: posStopListItems,
	}

	for _, storeID := range externalStoreIDs {
		transactionData := menuModels.TransactionData{
			StoreID:  storeID,
			Delivery: deliveryService,
			Products: slp,
		}
		trID, err := aggregatorService.UpdateStopListByProductsBulk(ctx, storeID, products, s.IsSendRemains(store))
		if err != nil {
			transactionData.Status = menuModels.ERROR
			transactionData.Message = err.Error()
		} else if deliveryService == models.YANDEX.String() {
			transactionData.Status = menuModels.PROCESSING
		} else {
			transactionData.Status = menuModels.SUCCESS
		}

		transactionData.ID = trID
		transaction.Transactions = append(transaction.Transactions, transactionData)
	}

	if err = s.repo.InsertStopListTransaction(ctx, transaction); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) updateStopListByProductIDsInDatabase(ctx context.Context, menuID string, products []menuModels.Product) error {
	if err := s.updateStopListAvailableStatusByProductIDsInDatabase(ctx, menuID, products); err != nil {
		return err
	}
	if err := s.updateDisabledStatusByProductIDsInDatabase(ctx, menuID, products); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) updateStopListAvailableStatusByProductIDsInDatabase(ctx context.Context, menuID string, products []menuModels.Product) error {
	var (
		productIdsWithAvailabilityFalse = make([]string, 0, len(products))
		productIdsWithAvailabilityTrue  = make([]string, 0, len(products))
	)

	for _, product := range products {
		if product.IsAvailable {
			productIdsWithAvailabilityTrue = append(productIdsWithAvailabilityTrue, product.ExtID)
		} else {
			productIdsWithAvailabilityFalse = append(productIdsWithAvailabilityFalse, product.ExtID)
		}
	}

	if len(productIdsWithAvailabilityTrue) != 0 {
		if err := s.menuService.UpdateProductsAvailabilityStatus(ctx, menuID, productIdsWithAvailabilityTrue, true); err != nil {
			return err
		}
	}

	if len(productIdsWithAvailabilityFalse) != 0 {
		if err := s.menuService.UpdateProductsAvailabilityStatus(ctx, menuID, productIdsWithAvailabilityFalse, false); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) toStopListProducts(products []menuModels.Product) menuModels.StopListProducts {

	res := make([]menuModels.StopListProduct, 0, len(products))

	for _, product := range products {
		stoplistProduct := menuModels.StopListProduct{
			ExtID:       product.ExtID,
			IsAvailable: product.IsAvailable,
		}

		if len(product.Name) != 0 {
			stoplistProduct.Name = product.Name[0].Value
		}

		if len(product.Price) != 0 {
			stoplistProduct.Price = product.Price[0].Value
		}

		res = append(res, stoplistProduct)
	}

	return res
}

func (s *ServiceImpl) IsSendRemains(store storeModels.Store) bool {
	restGroupsToSend := []string{
		// Safia
		"6683b1c1a58d7e792c955111",
		// Safia Kazakhstan
		"66e953265bc5291653cbc717",
	}

	for _, group := range restGroupsToSend {
		if store.RestaurantGroupID == group {
			return true
		}
	}
	return false
}

func (s *ServiceImpl) updateProductsDisabledByValidationInDatabase(ctx context.Context, menuID string, products []menuModels.Product) error {
	var (
		productIDsDisabledByValidationTrue  = make([]string, 0, len(products))
		productIDsDisabledByValidationFalse = make([]string, 0, len(products))
	)

	for _, product := range products {
		if product.DisabledByValidation {
			productIDsDisabledByValidationTrue = append(productIDsDisabledByValidationTrue, product.ExtID)
		} else {
			productIDsDisabledByValidationFalse = append(productIDsDisabledByValidationFalse, product.ExtID)
		}

	}

	if len(productIDsDisabledByValidationTrue) != 0 {
		log.Info().Msgf("update status on true in disabled by validation for products: %v", productIDsDisabledByValidationTrue)
		if err := s.menuService.UpdateProductsDisabledByValidation(ctx, menuID, productIDsDisabledByValidationTrue, true); err != nil {
			return err
		}
	}

	if len(productIDsDisabledByValidationFalse) != 0 {
		log.Info().Msgf("update status on false in disabled by validation for products: %v", productIDsDisabledByValidationFalse)
		if err := s.menuService.UpdateProductsDisabledByValidation(ctx, menuID, productIDsDisabledByValidationFalse, false); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) AddYandexTransaction(ctx context.Context, storeID, deliveryService string, products menuModels.StopListProducts, attributes menuModels.StopListAttributes) error {

	store, err := s.storeService.GetByExternalIdAndDeliveryService(ctx, storeID, deliveryService)
	if err != nil {
		return err
	}

	productTransaction, err := s.repo.GetLastTransactionByStoreIDForProducts(ctx, store.ID, deliveryService)
	if err != nil {
		log.Err(err).Msgf("get last transaction by store ID: %s for products: %+v", store.ID, products)
		return err
	}
	if err := s.repo.UpdateByTransactionID(ctx, productTransaction.ID, menuModels.SUCCESS.String(), products, nil); err != nil {
		log.Err(err).Msgf("update by transaction id: %s", productTransaction.ID)
		return err
	}

	attributeTransaction, err := s.repo.GetLastTransactionByStoreIDForAttributes(ctx, store.ID, deliveryService)
	if err != nil {
		log.Err(err).Msgf("get last transaction by store ID: %s for attributes: %+v", store.ID, attributes)
		return err
	}
	if err := s.repo.UpdateByTransactionID(ctx, attributeTransaction.ID, menuModels.SUCCESS.String(), nil, attributes); err != nil {
		log.Err(err).Msgf("update by transaction id: %s", attributeTransaction.ID)
		return err
	}
	log.Info().Msgf("successfully update product transaction id: %s and attribute transaction id: %s", productTransaction.ID, attributeTransaction.ID)
	return nil
}
