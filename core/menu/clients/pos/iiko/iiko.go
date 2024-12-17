package iiko

import (
	"context"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/base"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	iikoConf "github.com/kwaaka-team/orders-core/pkg/iiko/clients"
	iikoClient "github.com/kwaaka-team/orders-core/pkg/iiko/clients/http"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
)

var emptyImage = "https://share-menu.kwaaka.com/images/64214e349abc941cc0ef5d84/parsed-from-json/products/c062521c-ec54-447f-a9ff-7c3c5f5ec25e.jpg"

var litresSections = map[string]struct{}{
	"Напитки":        {},
	"Супы (с 11:00)": {},
	"Суп":            {},
	"Напитки на основе фруктов": {},
	"Горячие напитки":           {},
	"Алкогольные напитки":       {},
	"ФРЕШ":                      {},
	"СУПЫ":                      {},
	"СМУЗИ":                     {},
	"НАПИТКИ":                   {},
	"ДОМАШНИЕ ЛИМОНАДЫ":         {},
	"ДЕТОКС КОКТЕЙЛИ":           {},
	"ПИВО":                      {},
}

type manager struct {
	cli          iikoConf.IIKO
	globalConfig menu.Configuration
	menuRepo     drivers.MenuRepository
}

func NewIIKOManager(
	conf menu.Configuration,
	menuRepo drivers.MenuRepository,
	store storeModels.Store) (base.Manager, error) {

	baseURL := conf.IIKOConfiguration.BaseURL
	if store.PosType == models.SYRVE.String() {
		baseURL = conf.SyrveConfiguration.BaseURL
	}

	if store.IikoCloud.CustomDomain != "" {
		baseURL = store.IikoCloud.CustomDomain
	}

	cli, err := iikoClient.New(&iikoConf.Config{
		Protocol: "http",
		BaseURL:  baseURL,
		ApiLogin: store.IikoCloud.Key,
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize IIKO Client")
		return nil, err
	}

	return &manager{
		cli:          cli,
		globalConfig: conf,
		menuRepo:     menuRepo,
	}, nil
}

func (man manager) existCombos(ctx context.Context, menuID string) (map[string]models.Combo, error) {
	if menuID == "" {
		return map[string]models.Combo{}, nil
	}

	combos, _, err := man.menuRepo.GetCombos(ctx, selector.EmptyMenuSearch().
		SetMenuID(menuID))
	if err != nil {
		return nil, err
	}

	comboExist := make(map[string]models.Combo, len(combos))
	for _, combo := range combos {
		comboExist[combo.SourceActionID] = combo
	}

	return comboExist, nil
}

func (man manager) GetAggMenu(ctx context.Context, store storeModels.Store) ([]models.Menu, error) {
	if store.RestaurantGroupID != "646258ad1db3ef4dcf23c174" {
		return nil, errors.New("method is not for this restaurant")
	}
	rsp, err := man.cli.GetMenu(ctx, store.IikoCloud.OrganizationID)
	if err != nil {
		return nil, err
	}

	return man.aggMenusFromPos(store, rsp), nil
}

func (man manager) GetMenu(ctx context.Context, store storeModels.Store) (models.Menu, error) {
	if !store.IikoCloud.IsExternalMenu {
		rsp, err := man.cli.GetMenu(ctx, store.IikoCloud.OrganizationID)
		if err != nil {
			return models.Menu{}, err
		}

		products, err := man.existProducts(ctx, store.MenuID)
		if err != nil {
			return models.Menu{}, err
		}
		var (
			combos     iikoModels.GetCombosResponse
			comboExist map[string]models.Combo
		)

		if store.IikoCloud.HasCombo {
			combos, err = man.cli.GetCombos(ctx, iikoModels.GetCombosRequest{
				OrganizationID: store.IikoCloud.OrganizationID,
			})
			if err != nil {
				return models.Menu{}, err
			}

			comboExist, err = man.existCombos(ctx, store.MenuID)
			if err != nil {
				return models.Menu{}, err
			}
		}

		return menuFromClient(rsp, store.Settings, products, combos, comboExist), nil
	}

	rsp, err := man.cli.GetExternalMenu(ctx, store.IikoCloud.OrganizationID, store.IikoCloud.ExternalMenuID, store.IikoCloud.PriceCategory)
	if err != nil {
		return models.Menu{}, err
	}

	existProducts, err := man.getExistProducts(ctx, store.MenuID)
	if err != nil {
		return models.Menu{}, err
	}
	menu, err := man.menuRepo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(store.MenuID))
	if err != nil {
		return models.Menu{}, err
	}

	var collection models.MenuCollection
	if len(menu.Collections) > 0 {
		collection = menu.Collections[0]
	}

	return externalMenuFromClient(rsp, store.Settings, existProducts, collection, store.IikoCloud.IgnoreExternalMenuProductsWithZeroNullPrice), err
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
		productExist[product.ProductID+product.ParentGroupID+product.SizeID] = product
	}
	return productExist, nil
}

func (man manager) getExistProducts(ctx context.Context, menuID string) (map[string]models.Product, error) {

	if menuID == "" {
		return nil, nil
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

func (man manager) aggMenusFromPos(store storeModels.Store, posMenu iikoModels.GetMenuResponse) []models.Menu {
	woltMenu := models.Menu{
		Delivery:  models.WOLT.String(),
		IsActive:  true,
		Name:      store.Name + " generated by pos menu",
		IsDeleted: false,
	}

	glovoMenu := models.Menu{
		Delivery:  models.GLOVO.String(),
		IsActive:  true,
		Name:      store.Name + " generated by pos menu",
		IsDeleted: false,
	}

	mapGroups := make(map[string][]iikoModels.Group)
	for _, group := range posMenu.Groups {
		mapGroups[group.ParentGroup.UUID.String()] = append(mapGroups[group.ParentGroup.UUID.String()], group)
	}

	mapProducts := make(map[string][]iikoModels.Product)
	for _, product := range posMenu.Products {
		mapProducts[product.ParentGroup] = append(mapProducts[product.ParentGroup], product)
	}

	for _, group := range posMenu.Groups {
		if group.Name == "Меню Wolt" {
			var sections models.Sections
			var products models.Products
			for _, section := range mapGroups[group.ID.String()] {
				sections = append(sections, models.Section{
					ExtID:        section.ID.String(),
					Name:         section.Name,
					SectionOrder: section.Order,
					IsDeleted:    !section.InMenu,
				})

				for _, req := range mapProducts[section.ID.String()] {

					if req.Type == iikoModels.MODIFIER || !req.IsInMenu() {
						continue
					}

					product := models.Product{
						ExtID:            req.ID,
						PosID:            req.ID,
						Section:          req.ParentGroup, // req.Section is product categories in iiko, but here we used linked groups
						GroupID:          req.GroupID,
						ParentGroupID:    req.ParentGroup,
						ExtName:          req.Name,
						IsAvailable:      req.IsInMenu(),
						IsIncludedInMenu: req.IsInMenu(),
						Code:             req.Article,
						ImageURLs:        req.Images,
						ProductsCreatedAt: models.ProductsCreatedAt{
							Value:     coreModels.TimeNow(),
							Timezone:  store.Settings.TimeZone.TZ,
							UTCOffset: store.Settings.TimeZone.UTCOffset,
						},
						Name: []models.LanguageDescription{
							{
								Value:        req.Name,
								LanguageCode: store.Settings.LanguageCode,
							},
						},
						Description: []models.LanguageDescription{
							{
								Value:        req.Description,
								LanguageCode: store.Settings.LanguageCode,
							},
						},
						Price: []models.Price{
							{
								Value:        req.Price(),
								CurrencyCode: store.Settings.Currency,
							},
						},
						UpdatedAt:           coreModels.TimeNow(),
						Weight:              req.Weight,
						MeasureUnit:         "г",
						FatAmount:           req.FatAmount,
						ProteinsAmount:      req.ProteinsAmount,
						CarbohydratesAmount: req.CarbohydratesAmount,
						EnergyAmount:        req.EnergyAmount,

						FatFullAmount:           req.FatFullAmount,
						ProteinsFullAmount:      req.ProteinsFullAmount,
						CarbohydratesFullAmount: req.CarbohydratesFullAmount,
						EnergyFullAmount:        req.EnergyFullAmount,
						IsSync:                  true,
					}

					if req.IsDeleted {
						product.IsAvailable = false
						product.IsDeleted = true
					}

					if len(product.ImageURLs) == 0 {
						product.ImageURLs = []string{emptyImage}
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

					products = append(products, product)
				}
			}
			woltMenu.Sections = sections
			woltMenu.Products = products
		}

		if group.Name == "Меню Glovo" {
			var collections models.MenuCollections
			var sections models.Sections
			var products models.Products
			for _, collection := range mapGroups[group.ID.String()] {
				collections = append(collections, models.MenuCollection{
					ExtID:           collection.ID.String(),
					Name:            collection.Name,
					CollectionOrder: collection.Order,
					IsDeleted:       !collection.InMenu,
				})

				sectionID := uuid.New().String()

				sections = append(sections, models.Section{
					ExtID:        sectionID,
					Name:         collection.Name,
					SectionOrder: collection.Order,
					IsDeleted:    !collection.InMenu,
					Collection:   collection.ID.String(),
				})

				for _, req := range mapProducts[collection.ID.String()] {
					if req.Type == iikoModels.MODIFIER || !req.IsInMenu() {
						continue
					}

					product := models.Product{
						ExtID:            req.ID,
						PosID:            req.ID,
						Section:          sectionID, // req.Section is product categories in iiko, but here we used linked groups
						GroupID:          req.GroupID,
						ParentGroupID:    req.ParentGroup,
						ExtName:          req.Name,
						IsAvailable:      req.IsInMenu(),
						IsIncludedInMenu: req.IsInMenu(),
						Code:             req.Article,
						ImageURLs:        req.Images,
						ProductsCreatedAt: models.ProductsCreatedAt{
							Value:     coreModels.TimeNow(),
							Timezone:  store.Settings.TimeZone.TZ,
							UTCOffset: store.Settings.TimeZone.UTCOffset,
						},
						Name: []models.LanguageDescription{
							{
								Value:        req.Name,
								LanguageCode: store.Settings.LanguageCode,
							},
						},
						Description: []models.LanguageDescription{
							{
								Value:        req.Description,
								LanguageCode: store.Settings.LanguageCode,
							},
						},
						Price: []models.Price{
							{
								Value:        req.Price(),
								CurrencyCode: store.Settings.Currency,
							},
						},
						UpdatedAt:           coreModels.TimeNow(),
						Weight:              req.Weight * 1000,
						MeasureUnit:         "г",
						FatAmount:           req.FatAmount,
						ProteinsAmount:      req.ProteinsAmount,
						CarbohydratesAmount: req.CarbohydratesAmount,
						EnergyAmount:        req.EnergyAmount,

						FatFullAmount:           req.FatFullAmount,
						ProteinsFullAmount:      req.ProteinsFullAmount,
						CarbohydratesFullAmount: req.CarbohydratesFullAmount,
						EnergyFullAmount:        req.EnergyFullAmount,
						IsSync:                  true,
					}

					if req.IsDeleted {
						product.IsAvailable = false
						product.IsDeleted = true
					}

					if len(product.ImageURLs) == 0 {
						product.ImageURLs = []string{emptyImage}
					}

					if _, ok := litresSections[collection.Name]; ok {
						product.MeasureUnit = "мл"
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
					products = append(products, product)
				}
			}

			glovoMenu.Sections = sections
			glovoMenu.Collections = collections
			glovoMenu.Products = products
		}
	}

	attributes := make(models.Attributes, 0, len(posMenu.Products))
	modifierGroups := make(map[string]models.AttributeGroup, len(posMenu.Products))
	iikoModifierGroups := make(map[string]iikoModels.Group)

	for _, modifier := range posMenu.Products {

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

			for _, groupIIKO := range posMenu.Groups {
				iikoModifierGroups[groupIIKO.ID.String()] = groupIIKO
				if groupIIKO.ID.String() == modifier.GroupID {
					attribute.AttributeGroupName = groupIIKO.Name
				}
			}

		}
		attributes = append(attributes, attribute)

	}

	for _, product := range posMenu.Products {
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
		v.Name = iikoModifierGroups[v.ExtID].Name
		attributeGroups = append(attributeGroups, v)
	}
	woltMenu.Attributes = attributes
	woltMenu.AttributesGroups = attributeGroups
	glovoMenu.Attributes = attributes
	glovoMenu.AttributesGroups = attributeGroups

	return []models.Menu{woltMenu, glovoMenu}
}
