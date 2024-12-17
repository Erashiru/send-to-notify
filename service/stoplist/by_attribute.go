package stoplist

import (
	"context"
	"fmt"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/rs/zerolog/log"
)

func (s *ServiceImpl) UpdateStopListByAttributeID(ctx context.Context, isAvailable bool, storeID string, attributeID string) error {
	store, err := s.storeService.GetByID(ctx, storeID)
	if err != nil {
		return err
	}

	attributeIDs := []string{attributeID}

	posMenu, err := s.getPosMenu(ctx, store)
	if err != nil {
		return err
	}

	posAttributes, err := s.getAttributesByIDs(*posMenu, attributeIDs)
	if err != nil {
		return err
	}

	posAttributes = s.filterAttributes(posAttributes)

	if len(posAttributes) == 0 {
		log.Info().Msgf("stoplist servie - pos attributes for updating are not found")
		return nil
	}

	if err = s.stopAttributePosMenu(ctx, isAvailable, posMenu, posAttributes); err != nil {
		return err
	}

	for _, menu := range store.Menus {
		if err = s.stopAttributesAggregatorMenu(ctx, isAvailable, store, attributeIDs, menu.Delivery); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) stopAttributePosMenu(ctx context.Context, isAvailable bool, posMenu *menuModels.Menu, posAttributes []menuModels.Attribute) error {
	if len(posAttributes) == 0 {
		return nil
	}
	for i := range posAttributes {
		posAttributes[i].IsAvailable = isAvailable
	}

	if err := s.updateStopListByAttributesInDatabase(ctx, posMenu.ID, posAttributes); err != nil {
		return err
	}

	return nil

}

func (s *ServiceImpl) getAttributesByIDs(menu menuModels.Menu, attributeIDs []string) ([]menuModels.Attribute, error) {

	targetAttributes := make(map[string]interface{})
	for _, attributeID := range attributeIDs {
		targetAttributes[attributeID] = struct{}{}
	}

	var notDeletedOnly = make([]menuModels.Attribute, 0)
	for _, attribute := range menu.Attributes {
		if attribute.IsDeleted {
			continue
		}
		if _, ok := targetAttributes[attribute.ExtID]; !ok {
			continue
		}
		notDeletedOnly = append(notDeletedOnly, attribute)
	}

	return notDeletedOnly, nil
}

func (s *ServiceImpl) stopAttributesAggregatorMenu(ctx context.Context, isAvailable bool, store coreStoreModels.Store, attributeIDs []string, deliveryService string) error {
	isMenuExists := s.menuService.IsMenuExists(store, deliveryService)
	if !isMenuExists {
		return fmt.Errorf("menu %s for delivery service %s is not exists", store.ID, deliveryService)
	}
	aggregatorMenu, err := s.menuService.GetAggregatorMenuIfExists(ctx, store, deliveryService)
	if err != nil {
		return err
	}

	attributes, err := s.getAttributesByIDs(aggregatorMenu, attributeIDs)
	if err != nil {
		return err
	}

	attributes = s.filterAttributes(attributes)

	if len(attributes) == 0 {
		return nil
	}
	for i := range attributes {
		attributes[i].IsAvailable = isAvailable
	}

	if err = s.updateStopListByAttributesInAggregator(ctx, store, attributes, deliveryService, nil); err != nil {
		return err
	}

	if err = s.updateStopListByAttributesInDatabase(ctx, aggregatorMenu.ID, attributes); err != nil {
		return err
	}

	return nil
}
