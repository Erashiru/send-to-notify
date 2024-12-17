package dto

import (
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"time"
)

type StoreDsMenuDto struct {
	ID                     string    `json:"menu_id"`
	Name                   string    `json:"name"`
	IsActive               bool      `json:"is_active"`
	IsDeleted              bool      `json:"is_deleted"`
	IsSync                 bool      `json:"is_sync"`
	SyncAttributes         bool      `json:"sync_attributes"`
	Delivery               string    `json:"delivery"`
	Timestamp              int       `json:"timestamp"`
	UpdatedAt              time.Time `json:"updated_at"`
	Status                 string    `json:"status"`
	IsDiscount             bool      `json:"is_discount"`
	IsProductOnStop        bool      `json:"is_product_on_stop"`
	HasWoltPromo           bool      `json:"has_wolt_promo"`
	IsPosMenu              bool      `json:"is_pos_menu"`
	EmptyProductPercentage int       `json:"empty_product_percentage"`
	CreationSource         string    `json:"creation_source"`
}

func FromModel(m models.StoreDSMenu, isPosMenu bool) StoreDsMenuDto {
	return StoreDsMenuDto{
		ID:                     m.ID,
		Name:                   m.Name,
		IsActive:               m.IsActive,
		IsDeleted:              m.IsDeleted,
		IsSync:                 m.IsSync,
		SyncAttributes:         m.SyncAttributes,
		Delivery:               m.Delivery,
		Timestamp:              m.Timestamp,
		UpdatedAt:              m.UpdatedAt,
		Status:                 m.Status,
		IsDiscount:             m.IsDiscount,
		IsProductOnStop:        m.IsProductOnStop,
		HasWoltPromo:           m.HasWoltPromo,
		IsPosMenu:              isPosMenu,
		EmptyProductPercentage: m.EmptyProductPercentage,
		CreationSource:         m.CreationSource,
	}
}
