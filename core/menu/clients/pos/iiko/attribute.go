package iiko

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"math"
	"strings"

	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func getAttributes(req iikoModels.GetMenuResponse) (models.Attributes, []models.AttributeGroup) {

	// get attributes & attributes Groups and set to unique group
	attributes := make(models.Attributes, 0, len(req.Products))
	modifierGroups := make(map[string]models.AttributeGroup, len(req.Products))
	iikoModifierGroups := make(map[string]iikoModels.Group)

	for _, modifier := range req.Products {

		// attributes is Modifier type in IIKO terminal
		if modifier.Type != iikoModels.MODIFIER {
			continue
		}

		attribute := attributeToModel(modifier)

		if modifier.GroupID != "" {

			attribute.HasAttributeGroup = true

			// check
			if modifierGroup, ok := modifierGroups[modifier.GroupID]; ok {

				modifierGroup.Attributes = append(modifierGroup.Attributes, attribute.ExtID)
				attribute.AttributeGroupExtID = modifierGroup.ExtID
				modifierGroups[modifier.GroupID] = modifierGroup // set again

			} else {
				modifierGroups[modifier.GroupID] = models.AttributeGroup{
					ExtID:      modifier.GroupID,
					Min:        math.MaxInt,
					Attributes: []string{attribute.ExtID},
				}
				attribute.AttributeGroupExtID = modifier.GroupID
			}

			attribute.ParentAttributeGroup = modifier.ParentGroup

			for _, groupIIKO := range req.Groups {
				iikoModifierGroups[groupIIKO.ID.String()] = groupIIKO
				if groupIIKO.ID.String() == modifier.GroupID {
					attribute.AttributeGroupName = groupIIKO.Name
				}
			}

		}
		attributes = append(attributes, attribute)

	}

	for _, product := range req.Products {
		// product Type may be all type except 'modifier'
		if product.IsDeleted || product.Type == iikoModels.MODIFIER || !product.IsInMenu() {
			continue
		}

		// update modifier groups
		for _, groupModifier := range product.GroupModifiers {

			if modifier, ok := modifierGroups[groupModifier.ID]; ok {
				modifier.Min = int(math.Min(float64(modifier.Min), float64(groupModifier.MinAmount)))
				modifier.Max = int(math.Max(float64(modifier.Max), float64(groupModifier.MaxAmount)))
				modifierGroups[groupModifier.ID] = modifier
			}
		}
	}

	// set Attribute Groups
	attributeGroups := make([]models.AttributeGroup, 0, len(modifierGroups))
	for _, v := range modifierGroups {
		v.Name = iikoModifierGroups[v.ExtID].Name
		attributeGroups = append(attributeGroups, v)
	}

	return attributes, attributeGroups
}

func attributeToModel(req iikoModels.Product) models.Attribute {

	attribute := models.Attribute{
		ExtID:                req.ID,
		Code:                 req.Article,
		ExtName:              req.Name,
		Name:                 strings.ToLower(req.Name),
		Price:                req.Price(),
		IsAvailable:          true,
		IncludedInMenu:       req.IsInMenu(),
		AttributeGroupExtID:  req.GroupID,
		ParentAttributeGroup: req.ParentGroup,
	}

	if req.IsDeleted {
		attribute.IsAvailable = false
		attribute.IsDeleted = true
	}

	return attribute
}
