package models


type CancelOrderRequest struct {
	OrderID                 string           `json:"order_id"`
	StoreID 				string           `json:"store_id"`
	CancelReason            string           `json:"cancel_reason"`
	PaymentStrategy         string           `json:"payment_strategy"`
}