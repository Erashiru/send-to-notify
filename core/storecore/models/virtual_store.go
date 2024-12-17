package models

import (
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
)

type OrderUpdateEvent struct {
	RestaurantIDs []string `bson:"restaurant_ids"`
}

type VirtualStore struct {
	ID                       string   `bson:"id" json:"id"`
	RestaurantIds            []string `bson:"restaurant_ids" json:"restaurant_ids"`
	DeliveryService          string   `bson:"delivery_service" json:"delivery_service"`
	StoreIds                 []string `bson:"store_ids" json:"store_ids"`
	VirtualStoreRestaurantID string   `bson:"virtual_store_restaurant_id" json:"virtual_store_restaurant_id"`
	StoreType                string   `bson:"store_type" json:"store_type"`
	Name                     string   `bson:"name" json:"name"`
	ClientSecret             string   `bson:"client_secret" json:"client_secret"`
}

type OrderInfo struct {
	RestaurantName string                `bson:"restaurant_name" json:"restaurant_name"`
	RestaurantID   string                `bson:"restaurant_id" json:"restaurant_id"`
	Order          coreOrderModels.Order `bson:"order" json:"order"`
}
