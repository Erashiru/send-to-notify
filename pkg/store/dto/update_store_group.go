package dto

import (
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type UpdateStoreGroup struct {
	ID         *string  `bson:"id"`
	Name       *string  `bson:"name"`
	StoreIds   []string `bson:"restaurant_ids"`
	RetryCount *int     `bson:"retry_count"`
	ColumnView *bool    `bson:"column_view"`
}

func (s UpdateStoreGroup) ToModel() models.UpdateStoreGroup {
	return models.UpdateStoreGroup{
		ID:         s.ID,
		Name:       s.Name,
		StoreIds:   s.StoreIds,
		RetryCount: s.RetryCount,
		ColumnView: s.ColumnView,
	}
}
