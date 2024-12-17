package pos

import (
	"context"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"math"
	"strings"
)

func (iikoSvc *iikoService) getExistProducts(ctx context.Context, products []models.Product) (map[string]models.Product, error) {
	// add to hash map
	productExist := make(map[string]models.Product, len(products))
	for _, product := range products {
		// cause has cases if product_id && parent_id same, size_id different
		productExist[product.ProductID] = product
	}
	return productExist, nil
}

func (iikoSvc *iikoService) existCombos(ctx context.Context, combos []models.Combo) (map[string]models.Combo, error) {
	comboExist := make(map[string]models.Combo, len(combos))
	for _, combo := range combos {
		comboExist[combo.SourceActionID] = combo
	}

	return comboExist, nil
}

func (iikoSvc *iikoService) existProducts(ctx context.Context, products []models.Product) (map[string]models.Product, error) {
	// add to hash map
	productExist := make(map[string]models.Product, len(products))
	for _, product := range products {
		// cause has cases if product_id && parent_id same, size_id different
		productExist[product.ProductID+product.ParentGroupID+product.SizeID] = product
	}
	return productExist, nil
}

func (iikoSvc *iikoService) menuFromClient(req iikoModels.GetMenuResponse, store storeModels.Settings, products map[string]models.Product, combo iikoModels.GetCombosResponse, comboExist map[string]models.Combo) models.Menu {
	menu := models.Menu{
		Groups:   groupsToModel(req.Groups),
		Products: getProducts(req, store, products),
	}

	// get all attributes and set attribute groups ID from IIKO terminal
	attributes, attributeGroups := getAttributes(req)
	menu.Attributes = attributes
	menu.AttributesGroups = attributeGroups

	// get all collections
	// fixme: have to test this case
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
	// we can add to stop list:)
	// menu.StopLists = getStopLists(menu.Products)

	return menu
}

func (iikoSvc *iikoService) externalMenuFromClient(req iikoModels.GetExternalMenuResponse, settings storeModels.Settings, productsExist map[string]models.Product, collection models.MenuCollection, ignoreProductsWithZeroNullPrice bool) models.Menu {
	products := make(models.Products, 0, len(req.ItemCategories))
	attributesMap := make(map[string]models.Attribute)
	attributeGroupsMap := make(map[string]models.AttributeGroup)
	sections := make(models.Sections, 0, len(req.ItemCategories))

	//TODO: does external menu include info about КБЖУ? if yes, where it is?

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
				product.CarbohydratesAmount = itemSize.NutritionPerHundredGrams.Carbs
				product.ProteinsAmount = itemSize.NutritionPerHundredGrams.Proteins
				product.EnergyAmount = itemSize.NutritionPerHundredGrams.Energy
				product.FatAmount = itemSize.NutritionPerHundredGrams.Fats
				product.Weight = itemSize.PortionWeightGrams
				product.MeasureUnit = "г"
				for _, itemModifierGroup := range itemSize.ItemModifierGroups {
					attributes := make([]string, 0, len(itemModifierGroup.Items))
					for _, item := range itemModifierGroup.Items {
						if _, ok := attributesMap[item.ItemID]; !ok {
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

func getCollections(groups models.Groups) (models.MenuSuperCollections, models.MenuCollections, models.Sections) {

	// fixme: what about super collections?

	superCollections := make(models.MenuSuperCollections, 0, len(groups))
	collections := make(models.MenuCollections, 0, len(groups))
	sections := make(models.Sections, 0, len(groups))

	// fixme: there is not ok here
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
			ExtID:           group.ID,
			Name:            group.Name,
			ImageURL:        group.Image(),
			CollectionOrder: len(collections) + 1,
		})
	}

	return superCollections, collections, sections
}

func getProducts(req iikoModels.GetMenuResponse, settings storeModels.Settings, productExist map[string]models.Product) models.Products {

	// set `sizes` to hash, we need it if price of products more 1
	sizes := make(map[string]string, len(req.Sizes))
	for _, size := range req.Sizes {
		sizes[size.ID] = size.Name
	}

	// set products
	products := make(models.Products, 0, len(req.Products))

	// sameExtID allows check if iiko have same id with different parent id
	sameExtID := make(map[string]struct{}, len(req.Products))

	for _, product := range req.Products {

		// product Type may be all type except 'modifier'
		if product.Type == iikoModels.MODIFIER || !product.IsInMenu() {
			continue
		}

		// set product field except attributes Groups
		// here ext_id will be same product_id by default
		resProduct := productToModel(product, settings)

		// checking ext_id has in db
		prod, ok := productExist[resProduct.ProductID+resProduct.ParentGroupID+resProduct.SizeID]
		if ok {
			resProduct.ExtID = prod.ExtID

			for _, def := range prod.MenuDefaultAttributes {
				if def.ByAdmin {
					resProduct.MenuDefaultAttributes = append(resProduct.MenuDefaultAttributes, def)
				}
			}
		}

		// have to check that new product ext_id not repeated before,
		// if it is, have to generate new ext_id
		if _, ok := sameExtID[resProduct.ExtID]; ok {
			resProduct.ExtID = uuid.NewString()
		}
		sameExtID[resProduct.ExtID] = struct{}{}

		if len(product.Prices) == 1 {
			products = append(products, resProduct)
			continue
		}

		// get price from size if len > 2 && create new product
		for _, price := range product.Prices {
			if price.ID == "" || price.Price.Value == 0 {
				continue
			}
			products = append(products, productSizes(price, resProduct, sizes, productExist, settings))
			continue
		}
	}

	return products
}

// Product is a non-modifier type in IIKO terminal
func productToModel(req iikoModels.Product, setting storeModels.Settings) models.Product {

	product := models.Product{
		ExtID:            req.ID,
		ProductID:        req.ID,
		Section:          req.ParentGroup, // req.Section is product categories in iiko, but here we used linked groups
		GroupID:          req.GroupID,
		ParentGroupID:    req.ParentGroup,
		ExtName:          req.Name,
		IsAvailable:      req.IsInMenu(),
		IsIncludedInMenu: req.IsInMenu(),
		Code:             req.Article,
		ImageURLs:        req.Images,

		FatAmount:               req.FatAmount,
		EnergyAmount:            req.EnergyAmount,
		ProteinsAmount:          req.ProteinsAmount,
		CarbohydratesAmount:     req.CarbohydratesAmount,
		FatFullAmount:           req.FatFullAmount,
		EnergyFullAmount:        req.EnergyFullAmount,
		ProteinsFullAmount:      req.ProteinsFullAmount,
		CarbohydratesFullAmount: req.CarbohydratesFullAmount,

		ProductsCreatedAt: models.ProductsCreatedAt{
			Value:     coreModels.TimeNow(),
			Timezone:  setting.TimeZone.TZ,
			UTCOffset: setting.TimeZone.UTCOffset,
		},
		Name: []models.LanguageDescription{
			{
				Value:        req.Name,
				LanguageCode: setting.LanguageCode,
			},
		},
		Description: []models.LanguageDescription{
			{
				Value:        req.Description,
				LanguageCode: setting.LanguageCode,
			},
		},
		Price: []models.Price{
			{
				Value:        req.Price(),
				CurrencyCode: setting.Currency,
			},
		},
		UpdatedAt: coreModels.TimeNow(),
	}

	if req.IsDeleted {
		product.IsAvailable = false
		product.IsDeleted = true
	}

	attributes := make([]string, 0, len(req.Modifiers))
	defaults := make([]string, 0, len(req.Modifiers))
	defaultAttributes := make([]models.MenuDefaultAttributes, 0, len(req.Modifiers))
	attributesGroups := make([]string, 0, len(req.GroupModifiers))

	for _, modifier := range req.Modifiers {
		if modifier.DefaultAmount > 0 {
			defaultAttributes = append(defaultAttributes, defaultAttributeToModel(modifier))
		}
		attributes = append(attributes, modifier.ID)
	}

	for _, groupModifier := range req.GroupModifiers {

		// add free attributes from attribute groups
		for _, mod := range groupModifier.ChildModifiers {

			// add to default if amounts not null
			if mod.DefaultAmount > 0 && mod.FreeOfChargeAmount > 0 {
				defaults = append(defaults, mod.ID)
				defaultAttributes = append(defaultAttributes, models.MenuDefaultAttributes{
					ExtID:         mod.ID,
					DefaultAmount: mod.DefaultAmount,
				})
			}

			attributes = append(attributes, mod.ID)
		}

		attributesGroups = append(attributesGroups, groupModifier.ID)
	}

	product.Attributes = attributes
	product.AttributesGroups = attributesGroups
	product.MenuDefaultAttributes = defaultAttributes
	product.DefaultAttributes = defaults

	return product
}

func defaultAttributeToModel(req iikoModels.Modifier) models.MenuDefaultAttributes {
	return models.MenuDefaultAttributes{
		ExtID:         req.ID,
		DefaultAmount: req.DefaultAmount,
		ByAdmin:       false,
	}
}

func productSizes(price iikoModels.Price,
	product models.Product,
	sizes map[string]string,
	productExist map[string]models.Product,
	setting storeModels.Settings) models.Product {

	product.ExtID = uuid.NewString()

	// checking ext_id has in db
	if prod, ok := productExist[product.ProductID+product.ParentGroupID+price.ID]; ok {
		product.ExtID = prod.ExtID
	}

	product.Price = []models.Price{
		{
			Value:        price.Price.Value,
			CurrencyCode: setting.Currency,
		},
	}

	product.SizeID = price.ID

	if val, ok := sizes[price.ID]; ok {
		product.ExtName = strings.TrimSpace(product.ExtName + " " + val)
		product.Name = []models.LanguageDescription{
			{
				Value:        product.ExtName,
				LanguageCode: setting.LanguageCode,
			},
		}

	}

	return product
}

func getAttributes(req iikoModels.GetMenuResponse) (models.Attributes, []models.AttributeGroup) {

	// get attributes & attributes Groups and set to unique group
	attributes := make(models.Attributes, 0, len(req.Products))
	modifierGroups := make(map[string]models.AttributeGroup, len(req.Products))

	for _, modifier := range req.Products {

		// attributes is Modifier type in IIKO terminal
		if modifier.Type != iikoModels.MODIFIER {
			continue
		}

		attribute := attributeToModel(modifier)

		if modifier.GroupID != "" {

			attribute.HasAttributeGroup = true

			// check
			if modifierGroup, ok := modifierGroups[modifier.GroupID]; ok {

				modifierGroup.Attributes = append(modifierGroup.Attributes, attribute.ExtID)
				attribute.AttributeGroupExtID = modifierGroup.ExtID
				modifierGroups[modifier.GroupID] = modifierGroup // set again

			} else {
				modifierGroups[modifier.GroupID] = models.AttributeGroup{
					ExtID:      modifier.GroupID,
					Min:        math.MaxInt,
					Attributes: []string{attribute.ExtID},
				}
				attribute.AttributeGroupExtID = modifier.GroupID
			}

			attribute.ParentAttributeGroup = modifier.ParentGroup

			for _, groupIIKO := range req.Groups {
				if groupIIKO.ID.String() == modifier.GroupID {
					attribute.AttributeGroupName = groupIIKO.Name
				}
			}

		}
		attributes = append(attributes, attribute)

	}

	for _, product := range req.Products {
		// product Type may be all type except 'modifier'
		if product.IsDeleted || product.Type == iikoModels.MODIFIER || !product.IsInMenu() {
			continue
		}

		// update modifier groups
		for _, groupModifier := range product.GroupModifiers {
			if modifier, ok := modifierGroups[groupModifier.ID]; ok {
				modifier.Min = int(math.Min(float64(modifier.Min), float64(groupModifier.MinAmount)))
				modifier.Max = int(math.Max(float64(modifier.Max), float64(groupModifier.MaxAmount)))
				modifierGroups[groupModifier.ID] = modifier
			}
		}
	}

	// set Attribute Groups
	attributeGroups := make([]models.AttributeGroup, 0, len(modifierGroups))
	for _, v := range modifierGroups {
		attributeGroups = append(attributeGroups, v)
	}

	return attributes, attributeGroups
}

// Attribute is a modifier type in IIKO terminal
func attributeToModel(req iikoModels.Product) models.Attribute {

	attribute := models.Attribute{
		ExtID:                req.ID,
		Code:                 req.Article,
		ExtName:              req.Name,
		Name:                 strings.ToLower(req.Name),
		Price:                req.Price(),
		IsAvailable:          true,
		IncludedInMenu:       req.IsInMenu(),
		AttributeGroupExtID:  req.GroupID,
		ParentAttributeGroup: req.ParentGroup,
	}

	if req.IsDeleted {
		attribute.IsAvailable = false
		attribute.IsDeleted = true
	}

	return attribute
}

func getCombos(req iikoModels.GetCombosResponse, isExist map[string]models.Combo, menuProducts models.Products) ([]models.Combo, models.Products, models.AttributeGroups, models.Attributes) {
	productInfos := make(map[string]models.Product, len(menuProducts))

	for _, product := range menuProducts {
		productInfos[product.ExtID] = product
	}

	combos := make([]models.Combo, 0, len(req.ComboSpecifications))
	products := make(models.Products, 0, len(req.ComboSpecifications))
	attributeGroups := make(models.AttributeGroups, 0, len(req.ComboSpecifications))
	attributes := make(models.Attributes, 0, len(req.ComboSpecifications))

	// combos
	for _, comboSpecification := range req.ComboSpecifications {
		// groups
		attributeGroupIDs := make([]string, 0, len(comboSpecification.Groups))
		groups := make([]models.ComboGroup, 0, len(comboSpecification.Groups))

		for _, group := range comboSpecification.Groups {
			// products
			comboProducts := make([]models.ComboProduct, 0, len(group.Products))
			attributeIDs := make([]string, 0, len(group.Products))

			for _, product := range group.Products {
				val, ok := productInfos[product.ProductId]
				if !ok {
					comboProducts = append(comboProducts, models.ComboProduct{
						ProductId: product.ProductId,
						PriceModificationAmount: models.Price{
							Value: product.PriceModificationAmount,
						},
						IsExistInMenu: false,
					})
					continue
				}

				var name string

				if len(val.Name) != 0 {
					name = val.Name[0].Value
				}

				comboProducts = append(comboProducts, models.ComboProduct{
					ProductId: product.ProductId,
					PriceModificationAmount: models.Price{
						Value: product.PriceModificationAmount,
					},
					Name:          name,
					IsExistInMenu: true,
				})

				attributeIDs = append(attributeIDs, product.ProductId)

				attributes = append(attributes, models.Attribute{
					ExtID:             product.ProductId,
					Name:              name,
					ExtName:           name,
					Price:             product.PriceModificationAmount,
					IsAvailable:       val.IsAvailable,
					IsDeleted:         val.IsDeleted,
					HasAttributeGroup: true,
					IncludedInMenu:    true,
					IsComboAttribute:  true,
				})

			}

			groups = append(groups, models.ComboGroup{
				Id:          group.Id,
				Name:        group.Name,
				IsMainGroup: group.IsMainGroup,
				Products:    comboProducts,
			})

			attributeGroupIDs = append(attributeGroupIDs, group.Id)

			attributeGroups = append(attributeGroups, models.AttributeGroup{
				ExtID:        group.Id,
				Name:         group.Name,
				Min:          0,
				Max:          1,
				Attributes:   attributeIDs,
				IsComboGroup: true,
			})
		}

		id := uuid.New().String()
		var programID string

		if existCombo, ok := isExist[comboSpecification.SourceActionId]; ok {
			id = existCombo.ID
			programID = existCombo.ProgramID
		}

		combo := models.Combo{
			ID:   id,
			Name: comboSpecification.Name,
			Price: models.Price{
				Value: comboSpecification.PriceModification,
			},
			IsActive:       comboSpecification.IsActive,
			SourceActionID: comboSpecification.SourceActionId,
			ComboGroup:     groups,
			ProgramID:      programID,
		}

		product := models.Product{
			ExtID:     id,
			ProductID: id,
			Name: []models.LanguageDescription{
				{
					Value: comboSpecification.Name,
				},
			},
			Price: []models.Price{
				{
					Value: comboSpecification.PriceModification,
				},
			},
			ProductsCreatedAt: models.ProductsCreatedAt{
				Value: coreModels.TimeNow(),
			},
			IsCombo:          true,
			AttributesGroups: attributeGroupIDs,
			IsAvailable:      true,
			IsIncludedInMenu: true,
			UpdatedAt:        coreModels.TimeNow(),
		}

		products = append(products, product)
		combos = append(combos, combo)
	}

	return combos, products.Unique(), attributeGroups.Unique(), attributes.Unique()
}
