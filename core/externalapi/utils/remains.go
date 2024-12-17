package utils

import (
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/rs/zerolog/log"
	"strconv"
)

func ConvertStopListToRemains(stopList coreMenuModels.StopListResponse) models.RemainsResponse {
	response := models.RemainsResponse{
		Items: []models.RemainsItem{},
	}

	for _, product := range stopList.Products {
		stock, err := strconv.ParseFloat(product.Stock, 64)
		if err != nil {
			log.Err(err).Msg("Error parsing stock in yandex")
			continue
		}
		response.Items = append(response.Items, models.RemainsItem{
			Id:    product.ExtID,
			Stock: stock,
		})
	}

	return response
}
