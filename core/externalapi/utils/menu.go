package utils

import (
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"time"
)

func ParseMenu(req coreMenuModels.Menu) models.Menu {
	menu := models.Menu{
		Categories: []models.Category{},
		Items:      []models.Item{},
	}

	var zone, _ = time.LoadLocation("Asia/Almaty") // TODO From Config
	menu.LastChange = req.UpdatedAt.In(zone).Format(timeFormat)

	var attributesMap = make(map[string]coreMenuModels.Attribute, len(req.Attributes))
	var modifierGroupMap = make(map[string]models.ModifierGroup, len(req.AttributesGroups))

	for _, attribute := range req.Attributes {
		attributesMap[attribute.ExtID] = attribute
	}

	for _, attributeGroup := range req.AttributesGroups {

		modifierGroup := models.ModifierGroup{
			Id:                   attributeGroup.ExtID,
			Name:                 attributeGroup.Name,
			MinSelectedModifiers: attributeGroup.Min,
			MaxSelectedModifiers: attributeGroup.Max,
		}

		for _, id := range attributeGroup.Attributes {
			attribute, ok := attributesMap[id]
			if ok {
				modifier := models.Modifier{
					Id:        attribute.ExtID,
					Name:      attribute.Name,
					Price:     attribute.Price,
					MinAmount: attribute.Min,
					MaxAmount: attribute.Max,
				}

				modifierGroup.Modifiers = append(modifierGroup.Modifiers, modifier)
			}
		}

		modifierGroupMap[modifierGroup.Id] = modifierGroup
	}

	for _, collection := range req.Collections {
		category := models.Category{
			Id:        collection.ExtID,
			Name:      collection.Name,
			SortOrder: collection.CollectionOrder,
		}

		if collection.ImageURL != "" {
			category.Images = []models.CategoryImage{{Url: collection.ImageURL, UpdatedAt: collection.ImageUpdatedAt.In(zone).Format(timeFormat)}}
		}

		menu.Categories = append(menu.Categories, category)
	}

	for _, section := range req.Sections {
		category := models.Category{
			Id:        section.ExtID,
			ParentId:  section.Collection,
			Name:      section.Name,
			SortOrder: section.SectionOrder,
		}

		if section.ImageUrl != "" {
			category.Images = []models.CategoryImage{{Url: section.ImageUrl, UpdatedAt: section.ImageUpdatedAt.In(zone).Format(timeFormat)}}
		}

		menu.Categories = append(menu.Categories, category)
	}

	for _, product := range req.Products {
		//measureUnit, weight := measureMapping(product.MeasureUnit, product.Weight)

		if product.IsDeleted || !product.IsSync {
			continue
		}

		item := models.Item{
			Id:          product.ExtID,
			CategoryId:  product.Section,
			MeasureUnit: product.MeasureUnit,
			Measure:     int(product.Weight),
		}

		if item.MeasureUnit == "" {
			item.MeasureUnit = "г"
		}

		if item.Measure == 0 {
			item.Measure = 10
		}

		if product.ImageURLs != nil {
			for _, image := range product.ImageURLs {
				item.Images = append(item.Images, models.ItemImage{
					Url:  image,
					Hash: image,
				})
			}
		}

		if len(product.Name) != 0 {
			item.Name = product.Name[0].Value
		}

		if len(product.Price) != 0 {
			item.Price = product.Price[0].Value
		}

		if len(product.Description) != 0 {
			item.Description = product.Description[0].Value
		}

		for _, id := range product.AttributesGroups {
			modifierGroup, ok := modifierGroupMap[id]
			if ok {
				item.ModifierGroups = append(item.ModifierGroups, modifierGroup)
			}
		}

		nutrients := models.Nutrients{
			IsDeactivated: true,
		}

		if product.FatFullAmount != 0 {
			nutrients.Fat = product.FatFullAmount
		}

		if product.CarbohydratesFullAmount != 0 {
			nutrients.Carbohydrates = product.CarbohydratesFullAmount
		}

		if product.ProteinsFullAmount != 0 {
			nutrients.Proteins = product.ProteinsFullAmount
		}

		if product.EnergyFullAmount != 0 {
			nutrients.Calories = product.EnergyFullAmount
		}

		if nutrients.Fat != 0 || nutrients.Carbohydrates != 0 || nutrients.Proteins != 0 || nutrients.Calories != 0 {
			item.Nutrients = nutrients
		}

		if product.Halal {
			item.AdditionalDescriptions.Badges = append(item.AdditionalDescriptions.Badges, models.FoodSpecifics{
				Category: "food_specifics",
				Value:    "halal",
			})
		}

		menu.Items = append(menu.Items, item)
	}

	return menu
}

func ParseRetailMenu(req coreMenuModels.Menu) models.RetailMenu {
	menu := models.RetailMenu{
		Categories:        []models.Categories{},
		NomenclatureItems: []models.NomenclatureItem{},
	}

	for _, collection := range req.Collections {
		category := models.Categories{
			Id:       collection.ExtID,
			Name:     collection.Name,
			ParentId: nil, // TODO: refactor tree-like structure of retail menu
		}

		menu.Categories = append(menu.Categories, category)
	}

	for _, section := range req.Sections {
		category := models.Categories{
			Id:       section.ExtID,
			Name:     section.Name,
			ParentId: nil, // TODO: refactor tree-like structure of retail menu
		}

		menu.Categories = append(menu.Categories, category)
	}

	for _, product := range req.Products {
		if product.IsDeleted || !product.IsSync {
			continue
		}

		description := models.ItemDescription{}
		if len(product.Description) != 0 {
			description.General = product.Description[0].Value
		}

		images := []models.Images{}
		if product.ImageURLs != nil {
			for _, image := range product.ImageURLs {
				images = append(images, models.Images{
					Url: image,
				})
			}
		}

		unit, value := convertMeasureUnitAndValue(product.MeasureUnit, product.Weight)
		measure := models.Measure{
			Unit:  unit,
			Value: value,
		}

		nomenclatureItem := models.NomenclatureItem{
			Barcode:       convertBarcode(product.Barcode),
			CategoryId:    product.Section,
			Description:   description,
			Id:            product.ExtID,
			Images:        images,
			IsCatchWeight: product.IsCatchWeight,
			Measure:       measure,
			VendorCode:    product.VendorCode,
		}

		if len(product.Name) != 0 {
			nomenclatureItem.Name = product.Name[0].Value
		}

		if len(product.Price) != 0 {
			nomenclatureItem.Price = product.Price[0].Value
		}

		menu.NomenclatureItems = append(menu.NomenclatureItems, nomenclatureItem)
	}

	return menu
}

func convertBarcode(menuBarcode coreMenuModels.Barcode) models.Barcode {
	return models.Barcode{
		Type:           menuBarcode.Type,
		Value:          menuBarcode.Value,
		WeightEncoding: menuBarcode.WeightEncoding,
		Values:         menuBarcode.Values,
	}
}

func convertMeasureUnitAndValue(menuMeasureUnit string, menuMeasureValue float64) (string, int) {
	switch menuMeasureUnit {
	case "мл":
		return "MLT", int(menuMeasureValue)
	case "г":
		return "GRM", int(menuMeasureValue)
	case "л":
		return "MLT", int(menuMeasureValue * 1000)
	case "кг":
		return "GRM", int(menuMeasureValue * 1000)
	default:
		return "", 0
	}
}
