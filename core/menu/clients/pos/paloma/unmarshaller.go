package paloma

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	palomaModels "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
)

func menuFromClient(req palomaModels.Menu, settings storeModels.Settings, productsExist map[string]models.Product, aggregatorPriceTypeId string) models.Menu {

	menu := models.Menu{
		Name:        models.PALOMA.String(),
		ExtName:     models.MAIN.String(),
		Description: "paloma pos menu",
		CreatedAt:   coreModels.TimeNow(),
		UpdatedAt:   coreModels.TimeNow(),
	}

	products, groups, attributeGroups, attributes := toEntities(req, settings, productsExist, aggregatorPriceTypeId)

	menu.Products = products
	menu.Groups = groups
	menu.AttributesGroups = attributeGroups
	menu.Attributes = attributes

	return menu
}
