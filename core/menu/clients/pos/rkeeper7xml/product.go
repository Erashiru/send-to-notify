package rkeeper_xml

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	rkeeperXMLModels "github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models"
	"github.com/rs/zerolog/log"
	"strconv"
)

func (man manager) getProductAttributeGroupsIDs(modiSchemeID string, relationShip map[string][]string) []string {
	if ids, ok := relationShip[modiSchemeID]; ok {
		return ids
	}

	return nil
}

func (man manager) productsToModel(req rkeeperXMLModels.MenuRK7QueryResult, settings storeModels.Settings, mappingProducts map[string]string, relationship map[string][]string, existProducts map[string]models.Product, entityPriceMap map[string]string) models.Products {
	results := make(models.Products, 0, len(mappingProducts))

	for _, product := range req.RK7Reference.Items.Item {
		attributeGroupsIDs := man.getProductAttributeGroupsIDs(product.ModiScheme, relationship)

		res, err := man.productToModel(product, settings, mappingProducts, attributeGroupsIDs, existProducts, entityPriceMap)
		if err != nil {
			log.Err(err).Msgf("rkeeper7 xml cli err: get product %s", product.Ident)
			continue
		}

		results = append(results, res)
	}

	return results
}

func (man manager) productToModel(item rkeeperXMLModels.Item, setting storeModels.Settings, mapping map[string]string, attributeGroupIDs []string, existProducts map[string]models.Product, entityPriceMap map[string]string) (models.Product, error) {
	res := models.Product{
		ProductID: item.Ident, // TODO: id?
		ExtID:     item.Ident, // TODO: id?
		ExtName:   item.Name,
		Name: []models.LanguageDescription{
			{
				Value: item.Name,
			},
		},
		ProductsCreatedAt: models.ProductsCreatedAt{
			Value:     coreModels.TimeNow(),
			Timezone:  setting.TimeZone.TZ,
			UTCOffset: setting.TimeZone.UTCOffset,
		},
		AttributesGroups: attributeGroupIDs,
		IsIncludedInMenu: true,
		UpdatedAt:        coreModels.TimeNow(),
	}

	if val, ok := mapping[item.Ident]; !ok {
		res.IsAvailable = false

		price, exist := entityPriceMap[item.Ident]
		if !exist {
			res.Price = []models.Price{
				{
					Value:        0,
					CurrencyCode: setting.Currency,
				},
			}
		} else {
			priceInt, err := strconv.Atoi(price)
			if err != nil {
				res.Price = []models.Price{
					{
						Value:        0,
						CurrencyCode: setting.Currency,
					},
				}
			} else {
				res.Price = []models.Price{
					{
						Value:        float64(priceInt / 100),
						CurrencyCode: setting.Currency,
					},
				}
			}
		}

	} else {
		price, err := strconv.Atoi(val)
		if err != nil {
			return models.Product{}, err
		}

		price = price / 100

		res.Price = []models.Price{
			{
				Value:        float64(price),
				CurrencyCode: setting.Currency,
			},
		}
		res.IsAvailable = true
	}

	defaultAttributes := make([]models.MenuDefaultAttributes, 0)

	product, ok := existProducts[item.Ident]
	if ok {
		for _, def := range product.MenuDefaultAttributes {
			if def.ByAdmin {
				defaultAttributes = append(defaultAttributes, def)
			}
		}
	}

	res.MenuDefaultAttributes = defaultAttributes

	return res, nil
}
