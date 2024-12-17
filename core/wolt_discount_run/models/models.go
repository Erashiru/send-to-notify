package models

type DiscountRunRequest struct {
	MenuID       string `json:"menu_id"`
	RestaurantID string `json:"restaurant_id"`
}
