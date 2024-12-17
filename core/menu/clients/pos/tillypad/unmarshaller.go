package tillypad

import (
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storecoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
)

func (man manager) menuFromClient(posMenu models.Menu, store storecoreModels.Store) models2.Menu {

	convertedProducts, attributeGroups, attributes := man.productsToModel(posMenu.Items)

	menu := models2.Menu{
		Sections:         man.categoriesToModel(posMenu.Categories),
		Products:         convertedProducts,
		AttributesGroups: attributeGroups,
		Attributes:       attributes,

		UpdatedAt: coreModels.TimeNow(),
	}

	return menu
}

func (man manager) categoriesToModel(posMenuCategories []models.Category) models2.Sections {

	sections := make(models2.Sections, 0, len(posMenuCategories))

	for _, section := range posMenuCategories {
		sections = append(sections, man.categoryToModel(section))
	}

	return sections
}

func (man manager) categoryToModel(posMenuCategory models.Category) models2.Section {
	return models2.Section{
		ExtID: posMenuCategory.Id,
		Name:  posMenuCategory.Name,
	}
}

func (man manager) productsToModel(posMenuProducts []models.Item) (models2.Products, models2.AttributeGroups, models2.Attributes) {

	var (
		allAttributeGroups models2.AttributeGroups
		allAttributes      models2.Attributes
	)

	products := make(models2.Products, 0, len(posMenuProducts))

	for _, posMenuProduct := range posMenuProducts {
		product, attributeGroups, attributes := man.productToModel(posMenuProduct)
		products = append(products, product)
		allAttributeGroups = append(allAttributeGroups, attributeGroups...)
		allAttributes = append(allAttributes, attributes...)
	}

	return products, allAttributeGroups, allAttributes
}

func (man manager) productToModel(posMenuProduct models.Item) (models2.Product, models2.AttributeGroups, models2.Attributes) {

	var (
		attrGroupExtIds []string
	)

	attributeGroups, attributes := man.attributeGroupsToModel(posMenuProduct.ModifierGroups)
	for _, attrGroup := range attributeGroups {
		attrGroupExtIds = append(attrGroupExtIds, attrGroup.ExtID)
	}

	imageUrls := []string{}
	for _, image := range posMenuProduct.Images {
		imageUrls = append(imageUrls, image.Url)
	}

	return models2.Product{
		ExtID:   posMenuProduct.Id,
		Section: posMenuProduct.CategoryId,
		Name: []models2.LanguageDescription{
			{
				Value: posMenuProduct.Name,
			},
		},
		Description: []models2.LanguageDescription{
			{
				Value: posMenuProduct.Description,
			},
		},
		Price: []models2.Price{
			{
				Value: posMenuProduct.Price,
			},
		},
		Weight:           float64(posMenuProduct.Measure),
		MeasureUnit:      posMenuProduct.MeasureUnit,
		AttributesGroups: attrGroupExtIds,
		ImageURLs:        imageUrls,
	}, attributeGroups, attributes
}

func (man manager) attributeGroupsToModel(posMenuAttributeGroups []models.ModifierGroup) (models2.AttributeGroups, models2.Attributes) {

	var (
		attributeGroup models2.AttributeGroup
		attributes     models2.Attributes
		allAttributes  models2.Attributes
	)

	attributeGroups := make(models2.AttributeGroups, 0, len(posMenuAttributeGroups))

	for _, posAttributeGroup := range posMenuAttributeGroups {
		attributeGroup, attributes = man.attributeGroupToModel(posAttributeGroup)
		allAttributes = append(allAttributes, attributes...)
		attributeGroups = append(attributeGroups, attributeGroup)
	}

	return attributeGroups, allAttributes
}

func (man manager) attributeGroupToModel(posMenuAttributeGroup models.ModifierGroup) (models2.AttributeGroup, models2.Attributes) {

	attributes := man.attributesToModel(posMenuAttributeGroup.Id, posMenuAttributeGroup.Modifiers)
	attributesExtId := []string{}
	for _, attribute := range attributes {
		attributesExtId = append(attributesExtId, attribute.ExtID)
	}

	return models2.AttributeGroup{
		ExtID:      posMenuAttributeGroup.Id,
		Name:       posMenuAttributeGroup.Name,
		Attributes: attributesExtId,
		Min:        posMenuAttributeGroup.MinSelectedModifiers,
		Max:        posMenuAttributeGroup.MaxSelectedModifiers,
	}, attributes

}

func (man manager) attributesToModel(posMenuAttributeGroupId string, posMenuAttributes []models.Modifier) models2.Attributes {

	attributes := make(models2.Attributes, 0, len(posMenuAttributes))

	for _, posAttribute := range posMenuAttributes {
		attribute := man.attributeToModel(posAttribute)
		attribute.ParentAttributeGroup = posMenuAttributeGroupId
		attributes = append(attributes, attribute)
	}

	return attributes
}

func (man manager) attributeToModel(posMenuAttribute models.Modifier) models2.Attribute {

	return models2.Attribute{
		ExtID: posMenuAttribute.Id,
		Name:  posMenuAttribute.Name,
		Price: posMenuAttribute.Price,
		Min:   posMenuAttribute.MinAmount,
		Max:   posMenuAttribute.MaxAmount,
	}

}
