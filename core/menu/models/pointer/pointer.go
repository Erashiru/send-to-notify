package pointer

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
)

func OfBool(b bool) *bool {
	return &b
}

func OfFloat64(f float64) *float64 {
	return &f
}

func OfProduct(product models.Product) *models.Product {
	return &product
}

func OfString(b string) *string {
	return &b
}
