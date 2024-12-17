package external

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	externalModels "github.com/kwaaka-team/orders-core/pkg/externalapi/clients/dto"
)

func (m mnm) BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error) {
	var items = make([]externalModels.Product, 0, len(products))

	for _, product := range products {
		items = append(items, m.toItem(ctx, product))
	}

	for _, item := range items {
		if err := m.cli.UpdateProductStopList(ctx, item, m.webhookProductStoplist); err != nil {
			return "", err
		}
	}

	return "", nil
}
