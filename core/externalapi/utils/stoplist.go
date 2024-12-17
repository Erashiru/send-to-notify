package utils

import (
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
)

func ParseStopList(stopList coreMenuModels.StopListResponse) models.StopListResponse {
	response := models.StopListResponse{
		Items:     []models.StopListItem{},
		Modifiers: []models.StopListModifier{},
	}

	for _, product := range stopList.Products {
		stoplistItem := models.StopListItem{
			ItemId: product.ExtID,
		}
		response.Items = append(response.Items, stoplistItem)
	}

	for _, attribute := range stopList.Attributes {
		response.Modifiers = append(response.Modifiers, models.StopListModifier{
			ModifierId: attribute.AttributeID,
		})
	}

	return response
}
