package dto

type UpsertMenuRequest struct {
	StoreId string `json:"store_id"`
}

type BusyModeRequest struct {
	RestaurantID   string `json:"restaurant_id"`
	BusyMode       bool   `json:"busy_mode"`
	BusyModeMinute int    `json:"busy_mode_minute"`
}
