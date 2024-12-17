package models

type StoreJowiConfig struct {
	ApiKey       string `bson:"api_key" json:"api_key"`
	RestaurantID string `bson:"restaurant_id" json:"restaurant_id"`
}
