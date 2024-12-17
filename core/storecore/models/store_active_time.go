package models

import (
	"time"
)

type StoreActiveTime struct {
	ID              string    `bson:"_id,omitempty"`
	RestaurantID    string    `bson:"restaurant_id"`
	StoreID         string    `bson:"store_id"`
	DeliveryService string    `bson:"delivery_service"`
	StartTime       time.Time `bson:"start_time"`
	EndTime         time.Time `bson:"end_time,omitempty"`
}

type FilterStoreActiveTime struct {
	RestaurantID    string `bson:"restaurant_id"`
	StoreID         string `bson:"store_id"`
	DeliveryService string `bson:"delivery_service"`
}
