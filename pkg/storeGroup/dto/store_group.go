package dto

import (
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type StoreGroup struct {
	ID           string   `bson:"_id,omitempty"`
	IsTopPartner bool     `bson:"is_top_partner"`
	Name         string   `bson:"name"`
	StoreIDs     []string `bson:"restaurant_ids"`
}

func FromStoreGroupModel(req models.StoreGroup) StoreGroup {
	return StoreGroup{
		ID:           req.ID,
		Name:         req.Name,
		StoreIDs:     req.StoreIds,
		IsTopPartner: req.IsTopPartner,
	}
}
