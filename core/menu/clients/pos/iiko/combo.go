package iiko

import (
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func getCombos(req iikoModels.GetCombosResponse, isExist map[string]models.Combo, menuProducts models.Products) ([]models.Combo, models.Products, models.AttributeGroups, models.Attributes) {
	productInfos := make(map[string]models.Product, len(menuProducts))

	for _, product := range menuProducts {
		productInfos[product.ExtID] = product
	}

	combos := make([]models.Combo, 0, len(req.ComboSpecifications))
	products := make(models.Products, 0, len(req.ComboSpecifications))
	attributeGroups := make(models.AttributeGroups, 0, len(req.ComboSpecifications))
	attributes := make(models.Attributes, 0, len(req.ComboSpecifications))

	// combos
	for _, comboSpecification := range req.ComboSpecifications {
		// groups
		attributeGroupIDs := make([]string, 0, len(comboSpecification.Groups))
		groups := make([]models.ComboGroup, 0, len(comboSpecification.Groups))

		for _, group := range comboSpecification.Groups {
			// products
			comboProducts := make([]models.ComboProduct, 0, len(group.Products))
			attributeIDs := make([]string, 0, len(group.Products))

			for _, product := range group.Products {
				val, ok := productInfos[product.ProductId]
				if !ok {
					comboProducts = append(comboProducts, models.ComboProduct{
						ProductId: product.ProductId,
						PriceModificationAmount: models.Price{
							Value: product.PriceModificationAmount,
						},
						IsExistInMenu: false,
					})
					continue
				}

				var name string

				if len(val.Name) != 0 {
					name = val.Name[0].Value
				}

				comboProducts = append(comboProducts, models.ComboProduct{
					ProductId: product.ProductId,
					PriceModificationAmount: models.Price{
						Value: product.PriceModificationAmount,
					},
					Name:          name,
					IsExistInMenu: true,
				})

				attributeIDs = append(attributeIDs, product.ProductId)

				attributes = append(attributes, models.Attribute{
					ExtID:             product.ProductId,
					Name:              name,
					ExtName:           name,
					Price:             product.PriceModificationAmount,
					IsAvailable:       val.IsAvailable,
					IsDeleted:         val.IsDeleted,
					HasAttributeGroup: true,
					IncludedInMenu:    true,
					IsComboAttribute:  true,
				})

			}

			groups = append(groups, models.ComboGroup{
				Id:          group.Id,
				Name:        group.Name,
				IsMainGroup: group.IsMainGroup,
				Products:    comboProducts,
			})

			attributeGroupIDs = append(attributeGroupIDs, group.Id)

			attributeGroups = append(attributeGroups, models.AttributeGroup{
				ExtID:        group.Id,
				Name:         group.Name,
				Min:          0,
				Max:          1,
				Attributes:   attributeIDs,
				IsComboGroup: true,
			})
		}

		id := uuid.New().String()
		var programID string

		if existCombo, ok := isExist[comboSpecification.SourceActionId]; ok {
			id = existCombo.ID
			programID = existCombo.ProgramID
		}

		combo := models.Combo{
			ID:   id,
			Name: comboSpecification.Name,
			Price: models.Price{
				Value: comboSpecification.PriceModification,
			},
			IsActive:       comboSpecification.IsActive,
			SourceActionID: comboSpecification.SourceActionId,
			ComboGroup:     groups,
			ProgramID:      programID,
		}

		product := models.Product{
			ExtID:     id,
			ProductID: id,
			Name: []models.LanguageDescription{
				{
					Value: comboSpecification.Name,
				},
			},
			Price: []models.Price{
				{
					Value: comboSpecification.PriceModification,
				},
			},
			ProductsCreatedAt: models.ProductsCreatedAt{
				Value: coreModels.TimeNow(),
			},
			IsCombo:          true,
			AttributesGroups: attributeGroupIDs,
			IsAvailable:      true,
			IsIncludedInMenu: true,
			UpdatedAt:        coreModels.TimeNow(),
		}

		products = append(products, product)
		combos = append(combos, combo)
	}

	return combos, products.Unique(), attributeGroups.Unique(), attributes.Unique()
}
