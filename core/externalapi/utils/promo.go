package utils

import (
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
)

func ParsePromos(req coreMenuModels.Promo) models.Promo {
	response := models.Promo{
		PromoItems: []models.PromoItem{},
	}

	for _, item := range req.ProductGifts {
		response.PromoItems = append(response.PromoItems, models.PromoItem{
			Id:      item.ProductId,
			PromoId: item.PromoId,
		})
	}

	return response
}
