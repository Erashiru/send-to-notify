package rkeeper_xml

import (
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	rkeeperXMLModels "github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_menu_modifiers_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_groups"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_schema_details"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_schemas"
	"github.com/rs/zerolog/log"
	"strconv"
)

func (man manager) menuFromClient(items rkeeperXMLModels.MenuRK7QueryResult,
	modifiers get_menu_modifiers_response.RK7QueryResult,
	settings storeModels.Settings,
	mappingProducts map[string]string,
	mappingAttributes map[string]string,
	modifierGroups get_modifier_groups.RK7QueryResult,
	modifierSchemas get_modifier_schemas.RK7QueryResult,
	modifierSchemaDetails get_modifier_schema_details.RK7QueryResult,
	existProducts map[string]models.Product,
	entityPriceMap map[string]string) models.Menu {

	attributes, nonDeliveryModifiers := man.modifiersToModel(modifiers, settings, mappingAttributes)

	attributeGroups, newModifierSchemaDetails := man.modifierGroupsToModel(modifierGroups, modifierSchemaDetails, nonDeliveryModifiers)

	relationShip := man.buildRelationships(modifierSchemas, newModifierSchemaDetails)

	menu := models.Menu{
		Name:             models.RKEEPER7XML.String(),
		ExtName:          models.MAIN.String(),
		Description:      "rkeeper pos menu",
		AttributesGroups: attributeGroups,
		Products:         man.productsToModel(items, settings, mappingProducts, relationShip, existProducts, entityPriceMap),
		Attributes:       attributes,
		CreatedAt:        coreModels.TimeNow(),
		UpdatedAt:        coreModels.TimeNow(),
	}
	return menu
}

func (man manager) buildRelationships(modifierSchemas get_modifier_schemas.RK7QueryResult, modifierSchemaDetails get_modifier_schema_details.RK7QueryResult) map[string][]string {
	relationship := make(map[string][]string)

	for _, detail := range modifierSchemaDetails.RK7Reference.Items.Item {
		if val, ok := relationship[detail.ModiScheme]; ok {
			relationship[detail.ModiScheme] = append(val, detail.ModiGroup)
			continue
		}

		relationship[detail.ModiScheme] = []string{detail.ModiGroup}
	}

	return relationship
}

func (man manager) modifierGroupsToModel(modifierGroups get_modifier_groups.RK7QueryResult, modifierSchemaDetails get_modifier_schema_details.RK7QueryResult, nonDeliveryModifiers map[string]struct{}) (models.AttributeGroups, get_modifier_schema_details.RK7QueryResult) {
	unique := map[string]get_modifier_groups.Item{}

	for _, modifierGroup := range modifierGroups.RK7Reference.Items.Item {
		unique[modifierGroup.Ident] = modifierGroup
	}

	var attributeGroups = make(models.AttributeGroups, 0, len(modifierGroups.RK7Reference.Items.Item))

	for index, detail := range modifierSchemaDetails.RK7Reference.Items.Item {
		if val, ok := unique[detail.ModiGroup]; !ok {
			continue
		} else {
			ids := make([]string, 0, len(val.Childs.Child))

			for _, child := range val.Childs.Child {
				if _, ok := nonDeliveryModifiers[child.ChildIdent]; ok {
					continue
				}

				ids = append(ids, child.ChildIdent)
			}

			minLimit, err := strconv.Atoi(detail.DownLimit)
			if err != nil {
				log.Info().Msgf("atoi down limit error for %s", detail.Name)
			}

			maxLimit, err := strconv.Atoi(detail.UpLimit)
			if err != nil {
				log.Info().Msgf("atoi up limit error for %s", detail.Name)
			}

			if maxLimit == 0 {
				maxLimit = 1
			}

			extId := uuid.New().String()

			attributeGroups = append(attributeGroups, models.AttributeGroup{
				ExtID:      extId,
				Name:       val.Name,
				Attributes: ids,
				Min:        minLimit,
				Max:        maxLimit,
			})

			modifierSchemaDetails.RK7Reference.Items.Item[index].ModiGroup = extId
		}
	}

	return attributeGroups, modifierSchemaDetails
}
