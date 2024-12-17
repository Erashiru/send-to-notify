package iiko

import (
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func menuFromClient(req iikoModels.GetMenuResponse, store storeModels.Settings, products map[string]models.Product, combo iikoModels.GetCombosResponse, comboExist map[string]models.Combo) models.Menu {

	menu := models.Menu{
		Groups:   groupsToModel(req.Groups),
		Products: getProducts(req, store, products),
	}

	attributes, attributeGroups := getAttributes(req)
	menu.Attributes = attributes
	menu.AttributesGroups = attributeGroups

	superCollections, collections, sections := getCollections(menu.Groups)

	menu.Sections = sections
	menu.Collections = collections
	menu.SuperCollections = superCollections

	if len(combo.ComboSpecifications) != 0 {
		var (
			comboProducts        models.Products
			comboAttributeGroups models.AttributeGroups
			comboAttributes      models.Attributes
		)

		menu.Combos, comboProducts, comboAttributeGroups, comboAttributes = getCombos(combo, comboExist, menu.Products)

		menu.Products = append(menu.Products, comboProducts...)
		menu.AttributesGroups = append(menu.AttributesGroups, comboAttributeGroups...)
		menu.Attributes = append(menu.Attributes, comboAttributes...)
	}

	return menu
}

func externalMenuFromClient(req iikoModels.GetExternalMenuResponse, settings storeModels.Settings, productsExist map[string]models.Product, collection models.MenuCollection, ignoreProductsWithZeroNullPrice bool) models.Menu {
	products := make(models.Products, 0, len(req.ItemCategories))
	attributesMap := make(map[string]models.Attribute)
	attributeGroupsMap := make(map[string]models.AttributeGroup)
	sections := make(models.Sections, 0, len(req.ItemCategories))
	//createdUUIDList := make(map[string]string)
	if collection.ExtID == "" {
		collection.ExtID = uuid.New().String()
		collection.Name = "Меню"

	}
	for i, category := range req.ItemCategories {
		for _, item := range category.Items {
			product := models.Product{
				ProductID:        item.ItemID,
				ExtID:            item.ItemID,
				PosID:            item.ItemID,
				ExtName:          item.Name,
				Section:          category.ID,
				IsSync:           true,
				IsAvailable:      true,
				IsDeleted:        false,
				IsDisabled:       false,
				IsIncludedInMenu: true,
				Name: []models.LanguageDescription{
					{
						Value:        item.Name,
						LanguageCode: settings.LanguageCode,
					},
				},
				Description: []models.LanguageDescription{
					{
						Value:        item.Description,
						LanguageCode: settings.LanguageCode,
					},
				},
			}

			for _, itemSize := range item.ItemSizes {
				if itemSize.Prices != nil && len(itemSize.Prices) != 0 {
					product.Price = []models.Price{
						{
							Value:        itemSize.Prices[0].Price,
							CurrencyCode: settings.Currency,
						},
					}
					if itemSize.ButtonImageURL != "" {
						product.ImageURLs = []string{itemSize.ButtonImageURL}
					}
				}

				// check if nutrition fields are exist in NutritionPerHundredGrams as there's no info about them in iiko external menu endpoint
				// if one of the variants work if statements needs to be removed
				product.CarbohydratesAmount = itemSize.NutritionPerHundredGrams.CarbohydratesAmount
				if itemSize.NutritionPerHundredGrams.CarbohydratesAmount == 0 {
					product.CarbohydratesAmount = itemSize.NutritionPerHundredGrams.Carbs
				}
				product.ProteinsAmount = itemSize.NutritionPerHundredGrams.ProteinsAmount
				if itemSize.NutritionPerHundredGrams.ProteinsAmount == 0 {
					product.ProteinsAmount = itemSize.NutritionPerHundredGrams.Proteins
				}
				product.EnergyAmount = itemSize.NutritionPerHundredGrams.EnergyAmount
				if itemSize.NutritionPerHundredGrams.EnergyAmount == 0 {
					product.EnergyAmount = itemSize.NutritionPerHundredGrams.Energy
				}
				product.FatAmount = itemSize.NutritionPerHundredGrams.FatAmount
				if itemSize.NutritionPerHundredGrams.FatAmount == 0 {
					product.FatAmount = itemSize.NutritionPerHundredGrams.Fats
				}

				product.Weight = itemSize.PortionWeightGrams
				product.MeasureUnit = convertMeasureUnit(itemSize.MeasureUnitType)
				for _, itemModifierGroup := range itemSize.ItemModifierGroups {
					attributes := make([]string, 0, len(itemModifierGroup.Items))
					for _, item := range itemModifierGroup.Items {
						//uniqID := item.ItemID + "--" + itemModifierGroup.ItemGroupID + "--" + strconv.Itoa(item.Restrictions.MinQuantity) + "--" + strconv.Itoa(item.Restrictions.MaxQuantity)
						//if _, ok := createdUUIDList[uniqID]; !ok {
						//	guid := uuid.New().String()
						//	attributesMap[guid] = createAttribute(item, itemModifierGroup, guid)
						//	attributes = append(attributes, guid)
						//	createdUUIDList[uniqID] = guid
						//} else {
						//	attributes = append(attributes, createdUUIDList[uniqID])
						//
						//}
						var price float64
						if item.Prices != nil && len(item.Prices) != 0 {
							price = item.Prices[0].Price
						}
						attributesMap[item.ItemID] = models.Attribute{
							ExtID:                item.ItemID,
							PosID:                item.ItemID,
							ExtName:              item.Name,
							Name:                 item.Name,
							IsAvailable:          true,
							IncludedInMenu:       true,
							Min:                  item.Restrictions.MinQuantity,
							Max:                  item.Restrictions.MaxQuantity,
							ParentAttributeGroup: itemModifierGroup.ItemGroupID,
							AttributeGroupExtID:  itemModifierGroup.ItemGroupID,
							HasAttributeGroup:    true,
							AttributeGroupName:   itemModifierGroup.Name,
							AttributeGroupMax:    itemModifierGroup.Restrictions.MaxQuantity,
							AttributeGroupMin:    itemModifierGroup.Restrictions.MinQuantity,
							Price:                price,
						}

						attributes = append(attributes, item.ItemID)

						if item.Restrictions.ByDefault > 0 {
							var price float64
							if item.Prices != nil && len(item.Prices) != 0 {
								price = item.Prices[0].Price
							}
							product.MenuDefaultAttributes = append(product.MenuDefaultAttributes, models.MenuDefaultAttributes{
								ExtID:         item.ItemID,
								Price:         int(price),
								DefaultAmount: item.Restrictions.ByDefault,
								Name:          item.Name,
							})
							itemModifierGroup.IsDefault = true
						}
					}
					if itemModifierGroup.ItemGroupID == "" && !itemModifierGroup.IsDefault {
						for _, modifierItem := range itemModifierGroup.Items {
							var price float64
							if modifierItem.Prices != nil && len(modifierItem.Prices) != 0 {
								price = modifierItem.Prices[0].Price
							}
							product.MenuDefaultAttributes = append(product.MenuDefaultAttributes, models.MenuDefaultAttributes{
								ExtID:         modifierItem.ItemID,
								Price:         int(price),
								DefaultAmount: modifierItem.Restrictions.ByDefault,
								Name:          modifierItem.Name,
							})
						}
						itemModifierGroup.IsDefault = true
					}
					if itemModifierGroup.IsDefault {
						continue
					}

					internalAttributeGroupId := uuid.New().String()
					attributeGroupsMap[internalAttributeGroupId] = models.AttributeGroup{
						ExtID:      internalAttributeGroupId,
						PosID:      internalAttributeGroupId,
						Name:       itemModifierGroup.Name,
						Max:        itemModifierGroup.Restrictions.MaxQuantity,
						Min:        itemModifierGroup.Restrictions.MinQuantity,
						Attributes: attributes,
					}
					product.AttributesGroups = append(product.AttributesGroups, internalAttributeGroupId)
				}
			}

			if ignoreProductsWithZeroNullPrice && product.Price[0].Value == 0 {
				continue
			}

			products = append(products, product)
		}
		section := models.Section{
			ExtID:        category.ID,
			Name:         category.Name,
			SectionOrder: i,
			Description: []models.LanguageDescription{
				{
					Value:        category.Description,
					LanguageCode: settings.LanguageCode,
				},
			},
			Collection: collection.ExtID,
		}
		sections = append(sections, section)
	}

	menuAttributes := make(models.Attributes, 0, len(attributesMap))
	menuAttributeGroups := make(models.AttributeGroups, 0, len(attributeGroupsMap))

	for _, ag := range attributeGroupsMap {
		menuAttributeGroups = append(menuAttributeGroups, ag)
	}
	for _, a := range attributesMap {
		menuAttributes = append(menuAttributes, a)
	}

	return models.Menu{
		Collections:      models.MenuCollections{collection},
		Sections:         sections,
		Products:         products,
		AttributesGroups: menuAttributeGroups,
		Attributes:       menuAttributes,
	}
}

func groupsToModel(req []iikoModels.Group) []models.Group {
	groups := make([]models.Group, 0, len(req))
	for _, group := range req {
		groups = append(groups, groupToModel(group))
	}
	return groups
}

func groupToModel(group iikoModels.Group) models.Group {
	menuGroup := models.Group{
		ID:              group.ID.String(),
		Name:            group.Name,
		Description:     group.Description,
		Order:           group.Order,
		InMenu:          group.InMenu,
		IsGroupModifier: group.IsGroupModifier,
	}

	if len(group.Images) != 0 {
		menuGroup.Images = group.Images
	}

	if group.ParentGroup.Valid {
		menuGroup.ParentGroup = group.ParentGroup.UUID.String()
	}

	return menuGroup
}

func createAttribute(item iikoModels.ModifierGroupItems, itemModifierGroup iikoModels.ItemModifierGroups, extID string) models.Attribute {
	var price float64
	if item.Prices != nil && len(item.Prices) != 0 {
		price = item.Prices[0].Price
	}
	return models.Attribute{
		ExtID:                extID,
		PosID:                item.ItemID,
		ExtName:              item.Name,
		Name:                 item.Name,
		IsAvailable:          true,
		IncludedInMenu:       true,
		Min:                  item.Restrictions.MinQuantity,
		Max:                  item.Restrictions.MaxQuantity,
		ParentAttributeGroup: itemModifierGroup.ItemGroupID,
		AttributeGroupExtID:  itemModifierGroup.ItemGroupID,
		HasAttributeGroup:    true,
		AttributeGroupName:   itemModifierGroup.Name,
		AttributeGroupMax:    itemModifierGroup.Restrictions.MaxQuantity,
		AttributeGroupMin:    itemModifierGroup.Restrictions.MinQuantity,
		Price:                price,
	}
}

func getCollections(groups models.Groups) (models.MenuSuperCollections, models.MenuCollections, models.Sections) {

	superCollections := make(models.MenuSuperCollections, 0, len(groups))
	collections := make(models.MenuCollections, 0, len(groups))
	sections := make(models.Sections, 0, len(groups))

	for _, group := range groups {
		if group.ParentGroup == "" || group.ParentGroup == group.ID {
			superCollections = append(superCollections, models.MenuSuperCollection{
				ExtID:                group.ID,
				Name:                 group.Name,
				ImageUrl:             group.Image(),
				SuperCollectionOrder: len(superCollections) + 1,
			})
			continue
		}

		sections = append(sections, models.Section{
			ExtID:        group.ID,
			Name:         group.Name,
			ImageUrl:     group.Image(),
			Collection:   group.ParentGroup,
			IsDeleted:    group.InMenu,
			SectionOrder: len(sections) + 1,
		})

		collections = append(collections, models.MenuCollection{
			ExtID:           group.ParentGroup,
			Name:            group.Name,
			ImageURL:        group.Image(),
			CollectionOrder: len(collections) + 1,
		})
	}

	return superCollections, collections, sections
}

func convertMeasureUnit(str string) string {
	switch str {
	case "MILLILITER":
		return "мл"
	case "GRAM":
		return "г"
	case "LITER":
		return "л"
	case "KILOGRAM":
		return "кг"
	default:
		return ""
	}
}
