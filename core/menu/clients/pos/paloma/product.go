package paloma

import (
	"context"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	palomaModels "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
	"strconv"
)

const (
	intTrue    = 1
	stringTrue = "1"
)

func toEntities(req palomaModels.Menu, settings storeModels.Settings, productsExist map[string]models.Product, aggregatorPriceTypeId string) ([]models.Product, []models.Group, []models.AttributeGroup, []models.Attribute) {
	products := make([]models.Product, 0, 10)
	groups := make([]models.Group, 0, len(req.ItemGroups))
	attributeGroups := make([]models.AttributeGroup, 0, 4)
	attributes := make([]models.Attribute, 0, 4)

	var mapAttributeGroup = make(map[string]struct{})
	var mapAttribute = make(map[string]struct{})

	for _, itemGroup := range req.ItemGroups {
		for _, item := range itemGroup.Items {
			product, attributeGroups_, attributes_ := toEntity(item, settings, productsExist, aggregatorPriceTypeId)

			// unique attribute groups
			for _, attributeGroup := range attributeGroups_ {
				if _, ok := mapAttributeGroup[attributeGroup.ExtID]; ok {
					continue
				}

				mapAttributeGroup[attributeGroup.ExtID] = struct{}{}
				attributeGroups = append(attributeGroups, attributeGroup)
			}

			// unique attributes
			for _, attribute := range attributes_ {
				if _, ok := mapAttribute[attribute.ExtID]; ok {
					continue
				}

				mapAttribute[attribute.ExtID] = struct{}{}
				attributes = append(attributes, attribute)
			}

			product.ParentGroupID = strconv.Itoa(itemGroup.ObjectId)
			products = append(products, product)
		}

		groups = append(groups, models.Group{
			ID:     strconv.Itoa(itemGroup.ObjectId),
			Name:   itemGroup.Name,
			Images: []string{itemGroup.Image},
			InMenu: true,
		})
	}

	return products, groups, attributeGroups, attributes
}

func toEntity(req palomaModels.Item, settings storeModels.Settings, productsExist map[string]models.Product, aggregatorPriceTypeId string) (models.Product, []models.AttributeGroup, []models.Attribute) {
	extID := uuid.New().String()

	posProduct, ok := productsExist[strconv.Itoa(req.ObjectId)]
	if ok {
		extID = posProduct.ExtID
	}

	price := req.Price

	if aggregatorPriceTypeId != "" {
		for _, priceType := range req.OtherPrices.PriceTypes {
			if priceType.Id == aggregatorPriceTypeId {
				priceFloat, err := strconv.ParseFloat(priceType.Price, 64)
				if err != nil {
					continue
				}

				price = float64(int(priceFloat))
			}
		}
	}

	product := models.Product{
		ExtID:     extID,
		ProductID: strconv.Itoa(req.ObjectId),
		Name: []models.LanguageDescription{
			{
				Value:        req.Name,
				LanguageCode: settings.LanguageCode,
			},
		},
		Description: []models.LanguageDescription{
			{
				Value:        req.Description,
				LanguageCode: settings.LanguageCode,
			},
		},
		Price: []models.Price{
			{
				Value:        price,
				CurrencyCode: settings.Currency,
			},
		},
		ImageURLs:   []string{req.Image},
		IsAvailable: true,
	}

	if req.IUseInMenu == intTrue {
		product.IsIncludedInMenu = true
	}

	if req.MarkDeleted == intTrue {
		product.IsDeleted = true
	}

	for _, defaultAttribute := range posProduct.MenuDefaultAttributes {
		if defaultAttribute.ByAdmin {
			product.MenuDefaultAttributes = append(product.MenuDefaultAttributes, defaultAttribute)
		}
	}

	if len(req.ComplexGroups) > 0 {
		product.IsCombo = true
	}

	var attributeGroups = make([]models.AttributeGroup, 0, len(req.ModifierGroups))
	var attributes = make([]models.Attribute, 0, 4)

	for _, complexGroup := range req.ComplexGroups {
		attributeGroup := models.AttributeGroup{
			ExtID: strconv.Itoa(complexGroup.ObjectId),
			Name:  complexGroup.Name,
			Min:   complexGroup.MinCount,
			Max:   complexGroup.MaxCount,
		}

		product.AttributesGroups = append(product.AttributesGroups, strconv.Itoa(complexGroup.ObjectId))

		for _, modifier := range complexGroup.ComplexItems {
			attribute := models.Attribute{
				ExtID:                strconv.Itoa(modifier.ObjectId),
				Name:                 modifier.Name,
				Price:                modifier.Price,
				IsAvailable:          true,
				ParentAttributeGroup: strconv.Itoa(complexGroup.ObjectId),
			}

			if modifier.IUseInMenu == stringTrue {
				attribute.IncludedInMenu = true
			}

			if modifier.MarkDeleted == stringTrue {
				attribute.IsDeleted = true
			}

			attributes = append(attributes, attribute)
			attributeGroup.Attributes = append(attributeGroup.Attributes, strconv.Itoa(modifier.ObjectId))
		}

		attributeGroups = append(attributeGroups, attributeGroup)
	}

	for _, modifierGroup := range req.ModifierGroups {
		attributeGroup := models.AttributeGroup{
			ExtID: strconv.Itoa(modifierGroup.ObjectId),
			Name:  modifierGroup.Name,
		}

		product.AttributesGroups = append(product.AttributesGroups, strconv.Itoa(modifierGroup.ObjectId))

		for _, modifier := range modifierGroup.Modifiers {
			attribute := models.Attribute{
				ExtID:                strconv.Itoa(modifier.ObjectId),
				Name:                 modifier.Name,
				Price:                modifier.Price,
				IsAvailable:          true,
				ParentAttributeGroup: strconv.Itoa(modifierGroup.ObjectId),
			}

			if modifier.IUseInMenu == intTrue {
				attribute.IncludedInMenu = true
			}

			if modifier.MarkDeleted == intTrue {
				attribute.IsDeleted = true
			}

			attributes = append(attributes, attribute)
			attributeGroup.Attributes = append(attributeGroup.Attributes, strconv.Itoa(modifier.ObjectId))
		}

		attributeGroups = append(attributeGroups, attributeGroup)
	}

	return product, attributeGroups, attributes
}

func (man manager) existProducts(ctx context.Context, menuID string) (map[string]models.Product, error) {

	if menuID == "" {
		return map[string]models.Product{}, nil
	}

	// get products from main menu if exist
	products, _, err := man.menuRepo.ListProducts(ctx, selector.EmptyMenuSearch().
		SetMenuID(menuID))
	if err != nil {
		return nil, err
	}

	// add to hash map
	productExist := make(map[string]models.Product, len(products))
	for _, product := range products {
		// cause has cases if product_id && parent_id same, size_id different
		productExist[product.ProductID] = product
	}

	return productExist, nil
}
