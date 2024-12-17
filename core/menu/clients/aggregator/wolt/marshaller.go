package wolt

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	woltModels "github.com/kwaaka-team/orders-core/pkg/wolt/clients/dto"
	"sort"
	"strings"
)

const WoltTaxRate int = 12

func toWoltMenu(store storeModels.Store, menu models.Menu) woltModels.Menu {
	attributeGroupMap := make(map[string]models.AttributeGroup, len(menu.AttributesGroups))
	for _, attributeGroup := range menu.AttributesGroups {
		attributeGroupMap[attributeGroup.ExtID] = attributeGroup
	}
	attributeMap := make(map[string]models.Attribute, len(menu.Attributes))
	for _, attribute := range menu.Attributes {
		attributeMap[attribute.ExtID] = attribute
	}
	sectionMap := make(map[string]models.Section, len(menu.Sections))
	sort.Sort(menu.Sections)
	for _, section := range menu.Sections {
		sectionMap[section.ExtID] = section
	}

	categoriesMap := make(map[string]woltModels.Category)

	for _, product := range menu.Products {
		if !isValidProduct(product) {
			continue
		}

		tempItemNames := toItemNames(product.Name)
		tempItemDescription := toItemNames(product.Description)
		if !validateDescription(tempItemDescription) {
			tempItemDescription = tempItemNames
		}

		tempImageUrl := ""
		if len(product.ImageURLs) != 0 {
			tempImageUrl = product.ImageURLs[0]
		}

		tempOptions := make([]woltModels.OptionItem, 0, len(product.AttributesGroups))
		for _, attributeGroupID := range product.AttributesGroups {
			attributeGroupMapVal, ok := attributeGroupMap[attributeGroupID]
			if !ok {
				continue
			}
			if attributeGroupMapVal.ExtID == "" {
				continue
			}

			tempOption := woltModels.OptionItem{
				Name: addLanguageName(attributeGroupMapVal.NamesByLanguage, woltModels.ItemName{
					Lang:  store.Settings.LanguageCode,
					Value: attributeGroupMapVal.Name,
				}),
				ExternalData: attributeGroupMapVal.ExtID,
				Type:         toTypeSelection(attributeGroupMapVal),
				Values:       toOptionValues(attributeGroupMapVal, attributeMap, store.Settings.LanguageCode, product.ProductInformation),
			}
			if tempOption.Type == models.MULTICHOICE.String() {
				tempOption.SelectionRange = &woltModels.SelectionRange{
					Min: attributeGroupMapVal.Min,
					Max: attributeGroupMapVal.Max,
				}
			}

			tempOptions = append(tempOptions, tempOption)
		}

		regulatoryInformations := make([]woltModels.RegulatoryInformationValues, 0, len(product.ProductInformation.RegulatoryInformation))
		if product.ProductInformation.RegulatoryInformation != nil && len(product.ProductInformation.RegulatoryInformation) > 0 {
			for _, regInfo := range product.ProductInformation.RegulatoryInformation {
				req := woltModels.RegulatoryInformationValues{
					Name:  regInfo.Name,
					Value: regInfo.Value,
				}
				regulatoryInformations = append(regulatoryInformations, req)
			}
		}

		tempItem := woltModels.MenuItem{
			DeliveryMethods:    []string{models.HOMEDELIVERY.String(), models.TAKEAWAY.String()},
			Enabled:            product.IsAvailable,
			Price:              float32(product.Price[0].Value),
			ExternalData:       product.ExtID,
			ImageUrl:           tempImageUrl,
			SalesTaxPercentage: WoltTaxRate,
			AlcoholPercentage:  product.AlcoholPercentage,
			Name:               tempItemNames,
			Description:        tempItemDescription,
			Options:            tempOptions,
			ProductInformation: woltModels.ProductInformation{
				RegulatoryInformation: regulatoryInformations,
			},
		}
		sectionMapVal, ok := sectionMap[product.Section]
		if !ok {
			continue
		}
		temp, ok := categoriesMap[product.Section]
		if !ok {
			categoriesMap[product.Section] = woltModels.Category{
				ID:    product.Section,
				Items: []woltModels.MenuItem{tempItem},
				Name: addLanguageName(sectionMapVal.NamesByLanguage, woltModels.ItemName{
					Lang:  store.Settings.LanguageCode,
					Value: sectionMapVal.Name,
				}),
				Description: addLanguageName(sectionMapVal.Description, woltModels.ItemName{
					Lang:  store.Settings.LanguageCode,
					Value: sectionMapVal.Name,
				}),
			}
			continue
		}
		temp.Items = append(categoriesMap[product.Section].Items, tempItem)
		categoriesMap[product.Section] = temp
	}

	categories := make([]woltModels.Category, 0, len(menu.Sections))
	for _, section := range menu.Sections {
		if val, ok := categoriesMap[section.ExtID]; ok {
			categories = append(categories, val)
		}
	}

	return woltModels.Menu{
		Currency:        store.Settings.Currency,
		PrimaryLanguage: store.Settings.LanguageCode,
		Categories:      categories,
	}
}

func toTypeSelection(attributeGroup models.AttributeGroup) string {
	if attributeGroup.Max == 1 && attributeGroup.Min == 1 {
		return models.SINGLECHOICE.String()
	}
	return models.MULTICHOICE.String()
}

func toOptionValues(attributeGroup models.AttributeGroup, attributesMap map[string]models.Attribute, lang string, productInfo models.ProductInformation) []woltModels.ValueItem {
	res := make([]woltModels.ValueItem, 0, len(attributeGroup.Attributes))

	for _, attributeID := range attributeGroup.Attributes {
		attribute, ok := attributesMap[attributeID]
		if !ok {
			continue
		}

		regulatoryInforamation := make([]woltModels.RegulatoryInformationValues, 0, len(productInfo.RegulatoryInformation))
		if productInfo.RegulatoryInformation != nil && len(productInfo.RegulatoryInformation) > 0 {
			for _, regInfo := range productInfo.RegulatoryInformation {
				req := woltModels.RegulatoryInformationValues{
					Name:  regInfo.Name,
					Value: regInfo.Value,
				}
				regulatoryInforamation = append(regulatoryInforamation, req)
			}
		}

		tempValue := woltModels.ValueItem{
			Name: addLanguageName(attribute.NamesByLanguage, woltModels.ItemName{
				Lang:  lang,
				Value: attribute.Name,
			}),
			Enabled:      attribute.IsAvailable,
			Default:      false,
			Price:        float32(attribute.Price),
			ExternalData: attribute.ExtID,
			ProductInformation: woltModels.ProductInformation{
				RegulatoryInformation: regulatoryInforamation,
			},
		}
		res = append(res, tempValue)
	}

	if len(attributeGroup.AttributeMinMax) != 0 {
		for i := range attributeGroup.AttributeMinMax {
			attribute := attributeGroup.AttributeMinMax[i]

			for j := range res {
				if attribute.ExtId == res[j].ExternalData {
					res[j].SelectionRange = toOptionValuesSelectionRange(attribute.Min, attribute.Max, toTypeSelection(attributeGroup))
				}
			}
		}
	} else {
		for i := range attributeGroup.Attributes {
			attribute := attributeGroup.Attributes[i]

			for j := range res {
				if attribute == res[j].ExternalData {
					res[j].SelectionRange = toOptionValuesSelectionRange(0, attributeGroup.Max, toTypeSelection(attributeGroup))
				}
			}
		}
	}

	if len(res) > 0 {
		res[0].Default = true
	}
	return res
}

func toOptionValuesSelectionRange(min, max int, selectionType string) *woltModels.SelectionRange {
	if selectionType == models.SINGLECHOICE.String() {
		return nil
	}
	return &woltModels.SelectionRange{
		Min: min,
		Max: max,
	}
}

func toItemNames(names []models.LanguageDescription) []woltModels.ItemName {
	res := make([]woltModels.ItemName, 0, len(names))
	for _, name := range names {
		tempItemName := woltModels.ItemName{
			Lang:  name.LanguageCode,
			Value: name.Value,
		}
		if tempValue := strings.Trim(name.Value, " "); tempValue == "" {
			tempItemName.Value = ""
		}
		res = append(res, tempItemName)
	}
	return res
}

func addLanguageName(names []models.LanguageDescription, defaultName woltModels.ItemName) []woltModels.ItemName {
	res := []woltModels.ItemName{defaultName}

	if names == nil || len(names) == 0 {
		return res
	}

	for i := range names {
		if names[i].LanguageCode == defaultName.Lang {
			continue
		}

		res = append(res, woltModels.ItemName{
			Lang:  names[i].LanguageCode,
			Value: names[i].Value,
		})
	}

	return res
}

func validateDescription(descriptions []woltModels.ItemName) bool {
	if len(descriptions) == 0 {
		return false
	}
	for _, description := range descriptions {
		if description.Lang == "" || description.Value == "" {
			return false
		}
	}
	return true
}

func isValidProduct(product models.Product) bool {
	if product.ExtID == "" {
		return false
	}
	if product.IsDeleted {
		return false
	}
	if len(product.Price) == 0 {
		return false
	}

	return true
}

func fromWoltToSystemMenu(woltMenu woltModels.Menu) models.Menu {
	attributeGroups := make([]models.AttributeGroup, 0)
	attributes := make([]models.Attribute, 0)
	sections := make([]models.Section, 0)
	products := make([]models.Product, 0)

	for _, category := range woltMenu.Categories {
		for _, menuItem := range category.Items {
			productNames := make([]models.LanguageDescription, 0)
			for _, itemName := range menuItem.Name {
				productNames = append(productNames, models.LanguageDescription{
					LanguageCode: itemName.Lang,
					Value:        itemName.Value,
				})
			}

			productDescription := make([]models.LanguageDescription, 0)
			for _, itemDescription := range menuItem.Description {
				productDescription = append(productDescription, models.LanguageDescription{
					LanguageCode: itemDescription.Lang,
					Value:        itemDescription.Value,
				})
			}

			product := models.Product{
				ExtID:             menuItem.ExternalData,
				Name:              productNames,
				Description:       productDescription,
				IsAvailable:       menuItem.Enabled,
				Price:             []models.Price{{Value: float64(menuItem.Price), CurrencyCode: woltMenu.Currency}}, // Assuming single price value
				ImageURLs:         []string{menuItem.ImageUrl},                                                       // Assuming single image URL
				AttributesGroups:  make([]string, 0),                                                                 // Adjust based on your model structure
				Section:           category.ID,
				IsDeleted:         false, // Assuming all products are not deleted
				AlcoholPercentage: menuItem.AlcoholPercentage,
			}

			for _, option := range menuItem.Options {
				attributeGroup := models.AttributeGroup{
					ExtID:      option.ExternalData,
					Name:       option.Name[0].Value,
					Min:        option.SelectionRange.Min,
					Max:        option.SelectionRange.Max,
					Attributes: make([]string, 0),
				}

				for _, valueItem := range option.Values {
					attribute := models.Attribute{
						ExtID:               valueItem.ExternalData,
						Name:                valueItem.Name[0].Value,
						IsAvailable:         valueItem.Enabled,
						Price:               float64(valueItem.Price),
						Min:                 valueItem.SelectionRange.Min,
						Max:                 valueItem.SelectionRange.Max,
						Default:             valueItem.Default,
						AttributeGroupExtID: option.ExternalData,
						AttributeGroupName:  option.Name[0].Value,
						AttributeGroupMax:   option.SelectionRange.Max,
						AttributeGroupMin:   option.SelectionRange.Min,
					}

					attributeGroup.Attributes = append(attributeGroup.Attributes, attribute.ExtID)
					attributes = append(attributes, attribute)
				}

				attributeGroups = append(attributeGroups, attributeGroup)
			}

			products = append(products, product)
		}

		section := models.Section{
			ExtID: category.ID,
			Name:  category.Name[0].Value,
		}

		for _, description := range category.Description {
			section.Description = append(section.Description, models.LanguageDescription{
				LanguageCode: description.Lang,
				Value:        description.Value,
			})
		}

		sections = append(sections, section)
	}

	menu := models.Menu{
		AttributesGroups: attributeGroups,
		Attributes:       attributes,
		Sections:         sections,
		Products:         products,
		Delivery:         models.WOLT.String(),
		IsActive:         true,
		StopLists:        []string{},
	}

	return menu
}
