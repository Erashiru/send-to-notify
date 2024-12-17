package models

type EventStoplist struct {
	VirtualStoreID string `json:"virtual_restaurant_id"`
	RestaurantID   string `json:"real_restaurant_id"`
	Action         string `json:"action"`
}

const ClosingAction = "closing"
