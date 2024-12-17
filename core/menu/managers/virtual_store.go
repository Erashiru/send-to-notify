package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

// StopPositionsInVirtualStore - function to stop positions in virtual store (restaurantID  - actual _id of restaurant)
func (m *mnm) StopPositionsInVirtualStore(ctx context.Context, restaurantID, originalRestaurantID string) error {
	store, err := m.storeRepo.Get(ctx, selector.EmptyStoreSearch().SetID(restaurantID))
	if err != nil {
		return err
	}

	result := models.StopListTransaction{
		StoreID:   store.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	for _, menuInfo := range store.Menus {
		if !menuInfo.IsActive {
			continue
		}

		menu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuInfo.ID))
		if err != nil {
			return err
		}

		var (
			products   = make(models.Products, 0)
			attributes = make(models.Attributes, 0)
		)

		for _, product := range menu.Products {
			if strings.Contains(product.ExtID, originalRestaurantID) {
				product.IsAvailable = false
				products = append(products, product)
			}
		}

		for _, attribute := range menu.Attributes {
			if strings.Contains(attribute.ExtID, originalRestaurantID) {
				attribute.IsAvailable = false
				attributes = append(attributes, attribute)
			}
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

	m.stm.Insert(context.Background(), []models.StopListTransaction{result})

	return nil
}

// RenewPositionsInVirtualStore - function to renew positions in virtual store (restaurantID - _id of virtual store in db; originalRestaurantID - actual _id of restaurant)
func (m *mnm) RenewPositionsInVirtualStore(ctx context.Context, restaurantID, originalRestaurantID string) error {
	// getting virtual store
	virtualStore, err := m.storeRepo.Get(ctx, selector.EmptyStoreSearch().SetID(restaurantID))
	if err != nil {
		return err
	}

	// getting real store
	originalStore, err := m.storeRepo.Get(ctx, selector.EmptyStoreSearch().SetID(originalRestaurantID))
	if err != nil {
		return err
	}

	// getting real stores' pos menu
	originalPosMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(originalStore.MenuID))
	if err != nil {
		return err
	}

	var aggregatorMenuID string

	for _, menu := range virtualStore.Menus {
		if menu.IsActive && menu.Delivery == models.GLOVO.String() {
			aggregatorMenuID = menu.ID
			break
		}
	}

	// getting virtual stores' aggregator menu (glovo)
	aggregatorMenu, err := m.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(aggregatorMenuID))
	if err != nil {
		return err
	}

	result := models.StopListTransaction{
		StoreID:   virtualStore.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	var (
		products   = make(models.Products, 0, len(originalPosMenu.Products))
		attributes = make(models.Attributes, 0, len(originalPosMenu.Attributes)+len(originalPosMenu.Products))
	)

	// going through the products pos menu of a real restaurant
	for _, product := range originalPosMenu.Products {
		for _, aggregatorProduct := range aggregatorMenu.Products {
			if strings.Contains(aggregatorProduct.ExtID, product.ExtID) {
				aggregatorProduct.IsAvailable = product.IsAvailable
				products = append(products, aggregatorProduct)
			}
		}

		for _, aggregatorAttribute := range aggregatorMenu.Attributes {
			if strings.Contains(aggregatorAttribute.ExtID, product.ExtID) {
				aggregatorAttribute.IsAvailable = product.IsAvailable
				attributes = append(attributes, aggregatorAttribute)
			}
		}
	}

	// going through the attributes pos menu of a real restaurant
	for _, attribute := range originalPosMenu.Attributes {
		for _, aggregatorAttribute := range aggregatorMenu.Attributes {
			if strings.Contains(aggregatorAttribute.ExtID, attribute.ExtID) {
				aggregatorAttribute.IsAvailable = attribute.IsAvailable
				attributes = append(attributes, aggregatorAttribute)
			}
		}
	}

	for _, menu := range virtualStore.Menus {
		if !menu.IsActive {
			continue
		}

		if len(products) == 0 && len(attributes) == 0 {
			log.Trace().Msgf("no products && attributes changed")
			continue
		}

		trx, err := m.bulkUpdateAggregator(ctx, storeModels.AggregatorName(menu.Delivery), virtualStore, products.Unique(), attributes.Unique())
		if err != nil {
			log.Trace().Err(err).Msgf("could not bulk update store %s", virtualStore.ID)
			continue
		}

		if trx != nil {
			result.Transactions = append(result.Transactions, trx...)
		}
	}

	m.stm.Insert(context.Background(), []models.StopListTransaction{result})

	return nil
}
