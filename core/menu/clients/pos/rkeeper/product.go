package rkeeper

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"strconv"
	"strings"

	"github.com/google/uuid"
	rkeeperModels "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/dto"
	"github.com/rs/zerolog/log"
)

func productsToModel(req rkeeperModels.Menu, productExist map[string]string, store storeModels.Store) models.Products {

	results := make(models.Products, 0, len(req.Products))

	for _, product := range req.Products {
		if store.RKeeper.IgnoreUpsertProductWithPrice0 {
			price, err := strconv.ParseFloat(product.Price, 64)
			if err != nil {
				return models.Products{}
			}
			if price == 0 {
				continue
			}
		}
		res, err := productToModel(req, product, productExist, store.Settings)
		if err != nil {
			log.Err(err).Msgf("rkeeper cli err: get product %s", product.ID)
			continue
		}

		results = append(results, res)
	}

	return results
}

func productToModel(
	req rkeeperModels.Menu,
	product rkeeperModels.Product,
	productExist map[string]string,
	setting storeModels.Settings) (models.Product, error) {

	res := models.Product{
		ProductID:     product.ID,
		ExtID:         uuid.NewString(),
		ParentGroupID: product.CategoryID,
		Section:       product.CategoryID,
		IsAvailable:   true,
		ExtName:       product.Name,
		Name: []models.LanguageDescription{
			{
				Value: product.Name,
			},
		},
		Description: []models.LanguageDescription{
			{
				Value: product.Description,
			},
		},
		ImageURLs:   product.ImageUrls,
		MeasureUnit: product.Measure.Unit,
		ProductsCreatedAt: models.ProductsCreatedAt{
			Value:     coreModels.TimeNow(),
			Timezone:  setting.TimeZone.TZ,
			UTCOffset: setting.TimeZone.UTCOffset,
		},
		UpdatedAt: coreModels.TimeNow(),
	}

	price, err := strconv.ParseFloat(product.Price, 64)
	if err != nil {
		return models.Product{}, err
	}

	res.Price = []models.Price{
		{
			Value: price,
		},
	}

	if product.Measure.Value == "" {
		res.Weight = 0
	}

	measureValue, err := strconv.ParseFloat(product.Measure.Value, 64)
	if err != nil {
		log.Err(err).Msgf("product %s measure value %s", product.ID, product.Measure)
		measureValue = 0
	}

	res.Weight = measureValue

	// TODO: check this point
	// checking ext_id has in db
	key := strings.TrimSpace(res.ProductID + res.ParentGroupID)
	extID, ok := productExist[key]
	if ok {
		res.ExtID = extID
	}

	if product.SchemeId != "" {
		attributeGroups, attributes := collectAttributes(product.SchemeId, req)
		res.AttributesGroups = attributeGroups
		res.Attributes = attributes
	}

	return res, nil
}

func collectAttributes(schemeId string, menu rkeeperModels.Menu) ([]string, []string) {

	existAttributeGroup := make(map[string]struct{}, len(menu.IngredientsSchemes))

	attributeGroups := make([]string, 0, len(menu.IngredientsGroups))

	for _, scheme := range menu.IngredientsSchemes {
		if schemeId == scheme.ID {
			if len(scheme.IngredientsGroups) != 0 {
				for _, ingredientGroup := range scheme.IngredientsGroups {
					existAttributeGroup[ingredientGroup.ID] = struct{}{}
					attributeGroups = append(attributeGroups, ingredientGroup.ID)
				}
			}
		}
	}

	attributes := make([]string, 0, len(menu.Ingredients))

	for _, attributeGroup := range menu.IngredientsGroups {
		if _, ok := existAttributeGroup[attributeGroup.ID]; ok {
			attributes = append(attributes, attributeGroup.Ingredients...)
		}
	}

	return attributeGroups, attributes
}
