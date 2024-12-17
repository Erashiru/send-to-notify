package external

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	externalModels "github.com/kwaaka-team/orders-core/pkg/externalapi/clients/dto"
)

func (m mnm) toItem(ctx context.Context, product models.Product) externalModels.Product {
	if len(product.Price) == 0 {
		return externalModels.Product{}
	}

	return externalModels.Product{
		StoreID:   m.storeID,
		ProductID: product.ExtID,
		Price:     product.Price[0].Value,
		Available: product.IsAvailable,
	}
}

func (m mnm) toModifier(ctx context.Context, attribute models.Attribute) externalModels.Modifier {
	return externalModels.Modifier{
		StoreID:    m.storeID,
		ModifierID: attribute.ExtID,
		Price:      attribute.Price,
		Available:  attribute.IsAvailable,
	}
}
