package stoplist

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	models2 "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/rs/zerolog/log"
)

func (s *ServiceImpl) UpdateStopListBySectionID(ctx context.Context, isAvailable bool, storeGroupID string, deliveryToSectionIDs map[string][]string) error {
	log.Info().Msgf("Updating stop list by section id for restaurant group %s", storeGroupID)

	stores, err := s.storeService.GetStoresByStoreGroupID(ctx, storeGroupID)
	if err != nil {
		return err
	}

	for _, store := range stores {
		log.Info().Msgf("Store: %s", store.Name)

		for _, menu := range store.Menus {
			if !menu.IsActive {
				continue
			}

			sectionIDs, ok := deliveryToSectionIDs[menu.Delivery]
			if !ok {
				continue
			}

			products, err := s.getProductsBySectionIDsAndUpdateAvailableStatus(ctx, store, menu.ID, sectionIDs, isAvailable)
			if err != nil {
				log.Err(err).Msgf("Get products by section ids error")
				continue
			}

			aggUpdateProducts := s.filterProductsWithDisableByValidationStatus(products)
			log.Info().Msgf("Attempting to update stoplist for section IDs %+v in menu %s to available %t", sectionIDs, menu.ID, isAvailable)

			if err := s.updateStopListByProductIDInAggregator(ctx, store, aggUpdateProducts, menu.Delivery, nil); err != nil {
				log.Err(err).Msgf("[Aggregator error] Failed to update stop list by sectionIDs(%+v) for delivery %s in store: %s", sectionIDs, menu.Delivery, store.Name)
				continue
			}

			if err := s.updateStopListByProductIDsInDatabase(ctx, menu.ID, products); err != nil {
				log.Err(err).Msgf("[Database error] Failed to update stop list by sectionIDs(%+v) for delivery %s in store: %s", sectionIDs, menu.Delivery, store.Name)
			}

			log.Info().Msgf("Successfully updated stoplist by sectionIDs(%+v) for delivery %s in store %s", sectionIDs, menu.Delivery, store.Name)
		}
	}

	return nil
}

func (s *ServiceImpl) getProductsBySectionIDsAndUpdateAvailableStatus(ctx context.Context, store models2.Store, menuID string, sectionIDs []string, isAvailable bool) (models.Products, error) {
	var (
		products       = make([]models.Product, 0, 20)
		uniqueSections = make(map[string]bool)
	)

	aggregatorMenu, err := s.menuService.FindById(ctx, menuID)
	if err != nil {
		log.Err(err).Msgf("failed to get aggregator menu by id %s", menuID)
		return nil, err
	}

	for _, sectionId := range sectionIDs {
		uniqueSections[sectionId] = true
	}

	for i := range aggregatorMenu.Products {
		if !uniqueSections[aggregatorMenu.Products[i].Section] {
			continue
		}
		aggregatorMenu.Products[i].IsAvailable = isAvailable
		aggregatorMenu.Products[i].IsDisabled = !isAvailable

		products = append(products, aggregatorMenu.Products[i])
	}

	return products, nil
}

func (s *ServiceImpl) filterProductsWithDisableByValidationStatus(products []models.Product) []models.Product {
	if len(products) == 0 {
		return products
	}

	res := make([]models.Product, 0, len(products))

	for _, product := range products {
		if product.DisabledByValidation {
			continue
		}
		res = append(res, product)
	}

	return res
}
