package poster

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	posterModels "github.com/kwaaka-team/orders-core/pkg/poster/clients/models"
	"github.com/rs/zerolog/log"
	"strconv"
)

func fromHiddenToBool(hidden string) bool {
	return hidden != "1"
}

func fromSpotsVisibleToAvailable(storeSpotID string, posterProductSpots []posterModels.GetProductsResponseStop) bool {
	isAvailable := false
	for _, productSpot := range posterProductSpots {
		if productSpot.SpotId != storeSpotID {
			continue
		}
		isAvailable = productSpot.Visible != "0"
	}
	return isAvailable
}

func toEntities(posterProducts []posterModels.GetProductsResponseBody, store storeModels.Store, productsExist map[string]models.Product) ([]models.Product, []models.AttributeGroup, []models.Attribute, error) {
	var (
		products                = make([]models.Product, 0, len(posterProducts))
		attributeGroups         = make([]models.AttributeGroup, 0, 4)
		attributes              = make([]models.Attribute, 0, 4)
		uniqueAttributeGroups   = make(map[string]struct{})
		uniqueAttributes        = make(map[string]struct{})
		uniqueAvailableProducts = make(map[string]bool)
	)

	for _, item := range posterProducts {
		productAvailable := fromSpotsVisibleToAvailable(store.Poster.SpotId, item.Spots)
		product := models.Product{
			ExtID:        item.ProductId,
			ProductID:    item.ProductId,
			PosID:        item.ProductId,
			IngredientID: item.IngredientId,
			Section:      item.MenuCategoryId,
			Name: []models.LanguageDescription{
				{
					Value:        item.ProductName,
					LanguageCode: store.Settings.LanguageCode,
				},
			},
			IsAvailable:      productAvailable,
			IsIncludedInMenu: fromHiddenToBool(item.Hidden),
			ImageURLs:        []string{item.Photo},
		}
		uniqueAvailableProducts[item.ProductId] = productAvailable

		if item.Price.Field1 != "" {
			price, err := strconv.Atoi(item.Price.Field1)
			if err != nil {
				log.Err(err).Msgf("price for product %s is not exist", item.ProductName)
				continue
			}
			product.Price = []models.Price{
				{
					Value:        float64(price) / 100,
					CurrencyCode: store.Settings.Currency,
				},
			}
		}

		//teh-karta - have group-modifications
		for _, modifierGroup := range item.GroupModifications {
			log.Info().Msgf("GroupModifications  %v", modifierGroup.DishModificationGroupId)
			attributeGroupID := strconv.Itoa(modifierGroup.DishModificationGroupId)
			_, ok := uniqueAttributeGroups[attributeGroupID]
			if !ok {
				attributeGroup := models.AttributeGroup{
					PosID: attributeGroupID,
					ExtID: attributeGroupID,
					Name:  modifierGroup.Name,
				}

				attributeGroup.Min = modifierGroup.NumMin
				attributeGroup.Max = modifierGroup.NumMax

				for _, modifier := range modifierGroup.Modifications {
					attributeID := strconv.Itoa(modifier.DishModificationId)
					ingredientID := strconv.Itoa(modifier.IngredientId)
					_, exist := uniqueAttributes[attributeID]
					if !exist {
						attributeAvailable := true
						available, existsInProducts := uniqueAvailableProducts[ingredientID]
						if existsInProducts {
							attributeAvailable = available
						}
						attribute := models.Attribute{
							ExtID:        ingredientID,
							PosID:        ingredientID,
							IngredientID: attributeID,
							Name:         modifier.Name,
							IsAvailable:  attributeAvailable,
						}
						attribute.Price = float64(modifier.Price)
						attributes = append(attributes, attribute)
						uniqueAttributes[attributeID] = struct{}{}
					}

					//attributeIDs = append(attributeIDs, id)
					attributeGroup.Attributes = append(attributeGroup.Attributes, ingredientID)
				}

				attributeGroups = append(attributeGroups, attributeGroup)
				uniqueAttributeGroups[attributeGroupID] = struct{}{}
			}

			product.AttributesGroups = append(product.AttributesGroups, attributeGroupID)
		}

		if posProduct, posProductExists := productsExist[product.ExtID]; posProductExists {
			if posProduct.MenuDefaultAttributes != nil {
				for _, defaultAttribute := range posProduct.MenuDefaultAttributes {
					if defaultAttribute.ByAdmin {
						product.MenuDefaultAttributes = append(product.MenuDefaultAttributes, defaultAttribute)
					}
				}
			}
			product.CookingTime = posProduct.CookingTime
		}
		products = append(products, product)

		//size product; len(product) create
		for _, modifier := range item.Modifications {
			log.Info().Msgf("Modifications len Sizes len mods %v  mod_id% v", len(item.Modifications), modifier.ModificatorID)

			var sizeProduct = product

			sizeProduct.SizeID = modifier.ModificatorID
			sizeProduct.ExtID = product.ExtID + "_" + modifier.ModificatorID

			price, err := strconv.ParseFloat(modifier.ModificatorSelfprice, 64)
			if err != nil {
				return nil, nil, nil, err
			}

			sizeProduct.Price = []models.Price{{
				Value:        price / 100,
				CurrencyCode: store.Settings.Currency,
			}}

			if posSizeProduct, posSizeProductExists := productsExist[sizeProduct.ExtID]; posSizeProductExists {
				if posSizeProduct.MenuDefaultAttributes != nil {
					for _, defaultAttribute := range posSizeProduct.MenuDefaultAttributes {
						if defaultAttribute.ByAdmin {
							sizeProduct.MenuDefaultAttributes = append(sizeProduct.MenuDefaultAttributes, defaultAttribute)
						}
					}
				}
			}

			products = append(products, sizeProduct)
		}

	}

	log.Info().Msgf("len AG %+v len attr %v ", len(attributeGroups), len(attributes))

	return products, attributeGroups, attributes, nil
}

func (man manager) existProducts(ctx context.Context, menuID string) (map[string]models.Product, error) {

	if menuID == "" {
		return map[string]models.Product{}, nil
	}

	// get products from main menu if exist
	products, _, err := man.menuRepo.ListProducts(ctx, selector.EmptyMenuSearch().
		SetMenuID(menuID))
	if err != nil {
		return nil, err
	}

	// add to hash map
	productExist := make(map[string]models.Product, len(products))
	for _, product := range products {
		// cause has cases if product_id && parent_id same, size_id different
		productExist[product.ProductID] = product
	}

	return productExist, nil
}
