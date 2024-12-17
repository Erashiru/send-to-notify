package models

type Order struct {
	StoreID             string    `json:"store_id"`
	OrderID             string    `json:"order_id"`
	Products            []Product `json:"products"`
	OrderCode           string    `json:"order_code"`
	PickUpCode          string    `json:"pick_up_code"`
	UTCOffsetMinutes    string    `json:"utc_offset_minutes"`
	EstimatedPickupTime string    `json:"estimated_pickup_time"`
}

type CancelOrderRequest struct {
	OrderID         string `json:"order_id"`
	StoreID         string `json:"store_id"`
	CancelReason    string `json:"cancel_reason,omitempty"`
	PaymentStrategy string `json:"payment_strategy,omitempty"`
}
