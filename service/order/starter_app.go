package order

import (
	"context"
	"errors"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	starterAppModels "github.com/kwaaka-team/orders-core/core/starter_app/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
)

func (s *ServiceImpl) addExtIDAndNamesAndCookingTime(req interface{}, store storeModels.Store) (starterAppModels.Order, error) {
	r, ok := req.(starterAppModels.Order)
	if !ok {
		return starterAppModels.Order{}, errors.New("starter app addExtID casting error")
	}

	menu, err := s.menuService.GetAggregatorMenuIfExists(context.Background(), store, models.STARTERAPP.String())
	if err != nil {
		return starterAppModels.Order{}, err
	}

	mProducts := make(map[string]menuModels.Product, len(menu.Products))
	for _, product := range menu.Products {
		mProducts[product.StarterAppID] = product
	}

	mAttributes := make(map[string]menuModels.Attribute, len(menu.Attributes))
	for _, attribute := range menu.Attributes {
		mAttributes[attribute.StarterAppID] = attribute
	}

	for i := range r.OrderItems {
		if r.OrderItems[i].ExtID != "" {
			continue
		}
		product, ok := mProducts[r.OrderItems[i].MealId]
		if !ok {
			return starterAppModels.Order{}, errors.New("starter app order item not found in menu by MealId and StarterAppId")
		}
		productName := product.ExtName
		if product.Name != nil && len(product.Name) != 0 {
			productName = product.Name[0].Value
		}
		r.OrderItems[i].ExtID = product.ExtID
		r.OrderItems[i].Name = productName

		for j := range r.OrderItems[i].Modifiers {
			attribute, ok := mAttributes[r.OrderItems[i].Modifiers[j].ModifierId]
			if !ok {
				return starterAppModels.Order{}, errors.New("starter app order item modifier not found in menu by ModifierId and StarterAppId")
			}
			r.OrderItems[i].Modifiers[j].ExtId = attribute.ExtID
			r.OrderItems[i].Modifiers[j].Name = attribute.Name
		}

		cookingTime := store.StarterApp.CookingTime
		if cookingTime == 0 {
			cookingTime = 30
		}
		if product.CookingTime > cookingTime {
			cookingTime = product.CookingTime
		}
	}

	return r, nil
}
