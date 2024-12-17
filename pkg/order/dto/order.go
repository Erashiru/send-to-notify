package dto

import (
	"time"

	coreModels "github.com/kwaaka-team/orders-core/core/models"
)

type OrderSelector struct {
	ID                   string                     `json:"id"`
	OrderID              string                     `json:"order_id"`
	PosOrderID           string                     `json:"pos_order_id"`
	DeliveryService      string                     `json:"delivery_service"`
	ExternalStoreID      string                     `json:"external_store_id"`
	OrderCode            string                     `json:"order_code"`
	Restaurants          []string                   `json:"restaurants"`
	OnlyActive           bool                       `json:"only_active"`
	Status               string                     `json:"status"`
	OrderTimeFrom        time.Time                  `json:"start_date"`
	OrderTimeTo          time.Time                  `json:"end_date"`
	ReadingTime          coreModels.TransactionTime `json:"reading_time"`
	IsParentOrder        bool                       `json:"is_parent_order"`
	StoreID              string                     `json:"restaurant_id"`
	PosType              string                     `json:"pos_type"`
	IsPickedUpByCustomer bool                       `json:"is_picked_up_by_customer"`
	DeliveryArray        []string                   `json:"delivery_array"`
	coreModels.Customer
	Pagination
	Sorting
}

type ActiveOrderSelector struct {
	StoreID         string `json:"store_id"`
	DeliveryService string `json:"delivery_service"`
	PosType         string `json:"pos_type"`
	OrderCode       string `json:"order_code"`
}

type Pagination struct {
	Limit int64
	Page  int64
}

type Sorting struct {
	Param string
	Dir   int8
}

type UpdateOrder struct {
	OrderID      string
	PosID        string
	ReadingTime  coreModels.TransactionTime
	RestaurantID string
}
