package iiko

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"strings"

	"github.com/google/uuid"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

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
			resProduct.SizeID = product.Prices[0].ID
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
		UpdatedAt:           coreModels.TimeNow(),
		Weight:              req.Weight,
		FatAmount:           req.FatAmount,
		ProteinsAmount:      req.ProteinsAmount,
		CarbohydratesAmount: req.CarbohydratesAmount,
		EnergyAmount:        req.EnergyAmount,

		FatFullAmount:           req.FatFullAmount,
		ProteinsFullAmount:      req.ProteinsFullAmount,
		CarbohydratesFullAmount: req.CarbohydratesFullAmount,
		EnergyFullAmount:        req.EnergyFullAmount,
	}

	if req.IsDeleted {
		product.IsAvailable = false
		product.IsDeleted = true
	}
	if !req.IsInMenu() {
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
