package selector

import (
	"github.com/kwaaka-team/orders-core/core/menu/models/pointer"
	"time"
)

type MoySklad struct {
	ID           string    `json:"id"`
	OrderID      string    `bson:"order_id" json:"order_id"`
	PositionID   string    `bson:"position_id" json:"position_id"`
	MenuID       string    `bson:"menu_id" json:"menu_id"`
	RestaurantID string    `bson:"restaurant_id" json:"restaurant_id"`
	IsDeleted    *bool     `bson:"is_deleted" json:"is_deleted"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
	ProductID    string    `bson:"product_id" json:"product_id"`
	MsID         string    `bson:"ms_id" json:"ms_id"`
	Code         string    `bson:"code" json:"code"`
	Available    bool      `bson:"is_available" json:"available"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

func EmptyMoySkladSearch() MoySklad {
	return MoySklad{}
}

func (m MoySklad) SetRestaurantID(id string) MoySklad {
	m.RestaurantID = id
	return m
}

func (m MoySklad) HasRestaurantID() bool {
	return m.RestaurantID != ""
}
func (m MoySklad) SetIsDeleted(value bool) MoySklad {
	m.IsDeleted = pointer.OfBool(value)
	return m
}

func (m MoySklad) HasDeleted() *bool {
	return m.IsDeleted
}
