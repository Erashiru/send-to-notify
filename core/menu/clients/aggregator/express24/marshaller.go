package express24

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	express24Models "github.com/kwaaka-team/orders-core/pkg/express24/clients/dto"
	"github.com/rs/zerolog/log"
)

func toProducts(req models.Products) []express24Models.Product {
	products := make([]express24Models.Product, 0, len(req))

	for i := range req {
		var price int
		if len(req[i].Price) > 0 {
			price = int(req[i].Price[0].Value)
		}

		log.Info().Msgf("PRODUCT_ID: %v, AVAILABLE: %v", req[i].ExtID, req[i].IsAvailable)
		products = append(products, express24Models.Product{
			ExternalId:  req[i].ExtID,
			Quantity:    models.BASEQUANTITY,
			IsAvailable: toIsAvailable(req[i].IsAvailable),
			Price:       price,
		})
	}
	return products
}

func toIsAvailable(isAvailable bool) int {
	if isAvailable {
		return 1
	}
	return 0
}
