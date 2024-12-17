package jowi

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	jowiDto "github.com/kwaaka-team/orders-core/pkg/jowi/client/dto"
	"strconv"
)

func getProducts(req jowiDto.ResponseCourse, settings storeModels.Settings) models.Products {

	products := make(models.Products, 0, len(req.Courses))

	for _, product := range req.Courses {
		resProduct := productToModel(product, settings)
		products = append(products, resProduct)
	}

	return products
}

func productToModel(req jowiDto.Course, setting storeModels.Settings) models.Product {

	price, err := strconv.ParseFloat(req.PriceForOnlineOrder, 64)
	if err != nil {
		price = 0
	}

	weight, err := strconv.ParseFloat(req.Weight, 64)
	if err != nil {
		weight = 0
	}

	product := models.Product{
		ExtID:            req.Id,
		ProductID:        req.Id,
		Section:          req.CourseCategoryId, // req.Section is product categories in iiko, but here we used linked groups
		ExtName:          req.Title,
		Code:             req.PackageCode,
		ImageURLs:        []string{req.ImageUrl},
		Weight:           weight,
		MeasureUnit:      req.UnitName,
		IsAvailable:      req.OnlineOrder,
		IsIncludedInMenu: req.OnlineOrder,
		ProductsCreatedAt: models.ProductsCreatedAt{
			Value:     coreModels.TimeNow(),
			Timezone:  setting.TimeZone.TZ,
			UTCOffset: setting.TimeZone.UTCOffset,
		},
		Name: []models.LanguageDescription{
			{
				Value:        req.Title,
				LanguageCode: setting.LanguageCode,
			},
		},
		Price: []models.Price{
			{
				Value:        price,
				CurrencyCode: setting.Currency,
			},
		},
		UpdatedAt: coreModels.TimeNow(),
	}

	return product
}
