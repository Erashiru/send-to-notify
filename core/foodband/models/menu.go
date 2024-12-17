package models

import (
	"fmt"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
)

type UploadMenuReq struct {
	StoreID         string
	DeliveryService string
	MenuURL         string
}

type GetMenuUploadStatusReq struct {
	StoreID         string
	DeliveryService string
	TransactionID   string
}

type Attribute struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	SelectedByDefault bool    `json:"selected_by_default"`
	PriceImpact       float64 `json:"price_impact"`
	Available         bool    `json:"available"`
}

type AttributeGroup struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Min               int      `json:"min"`
	Max               int      `json:"max"`
	Collapse          bool     `json:"collapse"`
	MultipleSelection bool     `json:"multiple_selection"`
	Attributes        []string `json:"attributes"`
}

type Restrictions struct {
	IsAlcoholic bool `json:"is_alcoholic"`
	IsTobacco   bool `json:"is_tobacco"`
}

type Product struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Section         string       `json:"section"`
	Price           float64      `json:"price"`
	ImageURL        string       `json:"image_url"`
	ExtraImageURLs  []string     `json:"extra_image_urls"`
	Description     string       `json:"description"`
	Weight          float64      `json:"weight"`
	MeasureUnit     string       `json:"measure_unit"`
	Available       bool         `json:"available"`
	AttributeGroups []string     `json:"attributes_groups"`
	Restrictions    Restrictions `json:"restrictions"`
}

type Section struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	SectionOrder int    `json:"section_order"`
	Collection   string `json:"collection"`
}

type Collection struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	CollectionOrder int    `json:"collection_order"`
	ImageURL        string `json:"image_url"`
	SuperCollection string `json:"supercollection"`
}

type SuperCollection struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	SupercollectionOrder int    `json:"supercollection_order"`
	ImageURL             string `json:"image_url"`
}

type Menu struct {
	SuperCollections []SuperCollection `json:"supercollections"`
	Collections      []Collection      `json:"collections"`
	Sections         []Section         `json:"sections"`
	Products         []Product         `json:"products"`
	AttributeGroups  []AttributeGroup  `json:"attribute_groups"`
	Attributes       []Attribute       `json:"attributes"`
}

// Validate request menu and collect all errors to detais array
func (menu Menu) Validate(deliveryService string, store coreStoreModels.Store) (coreMenuModels.Menu, []string) {
	var details []string

	menuSections := make([]coreMenuModels.Section, 0, len(menu.Sections))
	menuSectionsMap := make(map[string]coreMenuModels.Section, len(menu.Sections))

	menuProducts := make([]coreMenuModels.Product, 0, len(menu.Products))
	menuProductsMap := make(map[string]coreMenuModels.Product, len(menu.Products))

	menuAttributes := make([]coreMenuModels.Attribute, 0, len(menu.Attributes))
	menuAttributesMap := make(map[string]coreMenuModels.Attribute, len(menu.Attributes))

	menuSuperCollections := make([]coreMenuModels.MenuSuperCollection, 0, len(menu.SuperCollections))
	menuSuperCollectionsMap := make(map[string]coreMenuModels.MenuSuperCollection, len(menu.SuperCollections))

	menuCollections := make([]coreMenuModels.MenuCollection, 0, len(menu.Collections))
	menuCollectionsMap := make(map[string]coreMenuModels.MenuCollection, len(menu.Collections))

	menuAttributeGroups := make([]coreMenuModels.AttributeGroup, 0, len(menu.AttributeGroups))
	menuAttributeGroupsMap := make(map[string]coreMenuModels.AttributeGroup, len(menu.AttributeGroups))

	if len(menu.Products) == 0 {
		details = append(details, "products can not be empty")
	}

	if len(menu.Sections) == 0 {
		details = append(details, "sections can not be empty")
	}

	// Valudate supercollections
	for idx, supercollection := range menu.SuperCollections {
		if supercollection.ID == "" {
			details = append(details, fmt.Sprintf("supercollections[%v].id is required", idx))
		}

		if supercollection.Name == "" {
			details = append(details, fmt.Sprintf("supercollections[%v].name is required", idx))
		}

		menuSuperCollection := coreMenuModels.MenuSuperCollection{
			ExtID:                supercollection.ID,
			Name:                 supercollection.Name,
			SuperCollectionOrder: supercollection.SupercollectionOrder,
			ImageUrl:             supercollection.ImageURL,
		}

		menuSuperCollections = append(menuSuperCollections, menuSuperCollection)
		menuSuperCollectionsMap[supercollection.ID] = menuSuperCollection
	}

	// Valudate collections
	for idx, collection := range menu.Collections {
		if collection.ID == "" {
			details = append(details, fmt.Sprintf("collections[%v].id is required", idx))
		}

		if collection.Name == "" {
			details = append(details, fmt.Sprintf("collections[%v].name is required", idx))
		}

		_, ok := menuSuperCollectionsMap[collection.SuperCollection]
		if !ok && collection.SuperCollection != "" {
			details = append(details, fmt.Sprintf("collections[%v].supercollection [%s] not found in supercollections", idx, collection.SuperCollection))
		}

		menuCollection := coreMenuModels.MenuCollection{
			ExtID:           collection.ID,
			Name:            collection.Name,
			CollectionOrder: collection.CollectionOrder,
			ImageURL:        collection.ImageURL,
			SuperCollection: collection.SuperCollection,
		}

		menuCollections = append(menuCollections, menuCollection)
		menuCollectionsMap[collection.ID] = menuCollection
	}

	// Validate Sections
	for idx, section := range menu.Sections {
		if section.ID == "" {
			details = append(details, fmt.Sprintf("sections[%v].id is required", idx))
		}

		if section.Name == "" {
			details = append(details, fmt.Sprintf("sections[%v].name is required", idx))
		}

		_, ok := menuCollectionsMap[section.Collection]
		if !ok && section.Collection != "" {
			details = append(details, fmt.Sprintf("sections[%v].collection [%s] not found in collections", idx, section.Collection))
		}

		menuSection := coreMenuModels.Section{
			ExtID:        section.ID,
			Name:         section.Name,
			SectionOrder: section.SectionOrder,
			Collection:   section.Collection,
		}

		menuSectionsMap[section.ID] = menuSection
		menuSections = append(menuSections, menuSection)
	}

	// Validate Attributes
	for idx, attribute := range menu.Attributes {
		if attribute.ID == "" {
			details = append(details, fmt.Sprintf("attributes[%v].id is required", idx))
		}

		if attribute.Name == "" {
			details = append(details, fmt.Sprintf("attributes[%v].name is required", idx))
		}

		menuAttribute := coreMenuModels.Attribute{
			ExtID:       attribute.ID,
			PosID:       attribute.ID,
			ExtName:     attribute.Name,
			Name:        attribute.Name,
			Default:     attribute.SelectedByDefault,
			Price:       attribute.PriceImpact,
			IsAvailable: attribute.Available,
		}

		menuAttributes = append(menuAttributes, menuAttribute)
		menuAttributesMap[attribute.ID] = menuAttribute
	}

	// Valudate Attribute groups
	for idx, attributeGroup := range menu.AttributeGroups {
		if attributeGroup.ID == "" {
			details = append(details, fmt.Sprintf("attribute_groups[%v].id is required", idx))
		}

		if attributeGroup.Name == "" {
			details = append(details, fmt.Sprintf("attribute_groups[%v].name is required", idx))
		}

		if len(attributeGroup.Attributes) == 0 {
			details = append(details, fmt.Sprintf("attribute_groups[%v].attribues can not be empty", idx))
		}

		if attributeGroup.Min < 0 {
			details = append(details, fmt.Sprintf("attribute_groups[%v].min can not be less than 0", idx))
		}

		if attributeGroup.Max < 1 {
			details = append(details, fmt.Sprintf("attribute_groups[%v].max can not be less than 1", idx))
		}

		for attributeIdx, attribute := range attributeGroup.Attributes {
			_, ok := menuAttributesMap[attribute]

			if !ok || attribute == "" {
				details = append(details, fmt.Sprintf("attribute_groups[%v].attributes[%v] [%s] not found in attributes or can not be empty", idx, attributeIdx, attribute))
			}
		}

		menuAttributeGroup := coreMenuModels.AttributeGroup{
			ExtID:          attributeGroup.ID,
			PosID:          attributeGroup.ID,
			Name:           attributeGroup.Name,
			Max:            attributeGroup.Max,
			Min:            attributeGroup.Min,
			Collapse:       attributeGroup.Collapse,
			MultiSelection: attributeGroup.MultipleSelection,
			Attributes:     attributeGroup.Attributes,
			IsSync:         true,
		}

		menuAttributeGroups = append(menuAttributeGroups, menuAttributeGroup)
		menuAttributeGroupsMap[attributeGroup.ID] = menuAttributeGroup
	}

	// Validate products
	for idx, product := range menu.Products {
		if product.ID == "" {
			details = append(details, fmt.Sprintf("products[%v].id is required", idx))
		}

		if product.Name == "" {
			details = append(details, fmt.Sprintf("products[%v].name is required", idx))
		}

		if product.Price <= 0 {
			details = append(details, fmt.Sprintf("products[%v].price is required or must me greated than 0", idx))
		}

		_, ok := menuSectionsMap[product.Section]
		if !ok || product.Section == "" {
			details = append(details, fmt.Sprintf("products[%v].section [%s] not found in sections", idx, product.Section))
		}

		if product.ImageURL != "" {
			product.ExtraImageURLs = append(product.ExtraImageURLs, product.ImageURL)
		}

		var productAttributes []string

		for attributeGroupIdx, attributeGroup := range product.AttributeGroups {
			menuAttributeGroup, ok := menuAttributeGroupsMap[attributeGroup]
			if !ok || attributeGroup == "" {
				details = append(details, fmt.Sprintf("products[%v].attribute_groups[%v] [%s] not found in attributes_groups or can not be empty", idx, attributeGroupIdx, attributeGroup))
			}

			productAttributes = append(productAttributes, menuAttributeGroup.Attributes...)
		}

		menuProduct := coreMenuModels.Product{
			ExtID:            product.ID,
			PosID:            product.ID,
			ProductID:        product.ID,
			Section:          product.Section,
			AttributesGroups: product.AttributeGroups,
			Attributes:       productAttributes,
			Name: []coreMenuModels.LanguageDescription{
				{
					Value:        product.Name,
					LanguageCode: store.Settings.LanguageCode,
				},
			},
			ExtName:   product.Name,
			ImageURLs: product.ExtraImageURLs,
			Description: []coreMenuModels.LanguageDescription{
				{
					Value:        product.Description,
					LanguageCode: store.Settings.LanguageCode,
				},
			},
			IsAvailable: product.Available,
			Price: []coreMenuModels.Price{
				{
					Value:        product.Price,
					CurrencyCode: store.Settings.Currency,
				},
			},
			IsSync:      true,
			Weight:      product.Weight,
			MeasureUnit: product.MeasureUnit,
		}

		menuProductsMap[product.ID] = menuProduct
		menuProducts = append(menuProducts, menuProduct)
	}

	res := coreMenuModels.Menu{
		Delivery:         deliveryService,
		Attributes:       menuAttributes,
		AttributesGroups: menuAttributeGroups,
		Products:         menuProducts,
		Sections:         menuSections,
		Collections:      menuCollections,
		SuperCollections: menuSuperCollections,
	}

	return res, details
}
