package rkeeper_xml

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_menu_modifiers_response"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func (man manager) modifiersToModel(menuModifiers get_menu_modifiers_response.RK7QueryResult, settings storeModels.Settings, mappingAttributes map[string]string) (models.Attributes, map[string]struct{}) {
	results := make(models.Attributes, 0, len(mappingAttributes))

	nonDeliveryModifiers := make(map[string]struct{})

	for _, item := range menuModifiers.RK7Reference.Items.Item {
		if strings.Contains(item.Name, "зал") {
			nonDeliveryModifiers[item.Ident] = struct{}{}
			continue
		}

		res, err := man.modifierToModel(item, settings, mappingAttributes)
		if err != nil {
			log.Err(err).Msgf("rkeeper7 xml cli err: get attribute %s", item.Ident)
			continue
		}

		results = append(results, res)
	}

	return results, nonDeliveryModifiers
}

func (man manager) modifierToModel(item get_menu_modifiers_response.Item, setting storeModels.Settings, mappingAttributes map[string]string) (models.Attribute, error) {
	res := models.Attribute{
		ExtID:          item.Ident,
		ExtName:        item.Name,
		Name:           item.Name,
		UpdatedAt:      coreModels.TimeNow().Time,
		IncludedInMenu: true,
	}

	if val, ok := mappingAttributes[item.Ident]; !ok {
		res.IsAvailable = false
	} else {
		price, err := strconv.Atoi(val)
		if err != nil {
			return models.Attribute{}, err
		}

		price = price / 100

		res.Price = float64(price)
		res.IsAvailable = true
	}

	return res, nil
}
