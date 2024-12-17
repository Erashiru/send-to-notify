package models

type AnotherBillParams struct {
	ID           string   `bson:"_id,omitempty"`
	RestaurantID string   `bson:"restaurant_id"`
	Deliveries   []string `bson:"deliveries,omitempty"`
}
