package poster

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	posterModels "github.com/kwaaka-team/orders-core/pkg/poster/clients/models"
)

func menuFromClient(posterProducts posterModels.GetProductsResponse, store storeModels.Store, productsExist map[string]models.Product) (models.Menu, error) {
	products, attributeGroups, attributes, err := toEntities(posterProducts.Response, store, productsExist)
	if err != nil {
		return models.Menu{}, err
	}

	menu := models.Menu{
		Name:             models.POSTER.String(),
		ExtName:          models.MAIN.String(),
		CreatedAt:        coreModels.TimeNow(),
		UpdatedAt:        coreModels.TimeNow(),
		Products:         products,
		AttributesGroups: attributeGroups,
		Attributes:       attributes,
	}

	return menu, nil
}
