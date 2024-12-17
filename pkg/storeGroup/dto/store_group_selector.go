package dto

import (
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type StoreGroupSelector struct {
	ID       string
	Name     string
	StoreIDs []string
}

type CreateStoreGroupRequest struct {
	Name     string   `json:"name"`
	StoreIDs []string `json:"store_ids"`
}

func ToStoreGroupModel(query StoreGroup) models.StoreGroup {
	return models.StoreGroup{
		Name:     query.Name,
		StoreIds: query.StoreIDs,
	}
}
