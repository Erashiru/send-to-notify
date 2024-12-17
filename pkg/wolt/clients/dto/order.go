package dto

import "time"

type AcceptOrderRequest struct {
	ID         string     `json:"-"`
	PickupTime *time.Time `json:"adjusted_pickup_time,omitempty"`
}

type RejectOrderRequest struct {
	ID     string `json:"-"`
	Reason string `json:"reason"`
}

type AcceptSelfDeliveryOrderOrderRequest struct {
	ID           string     `json:"-"`
	DeliveryTime *time.Time `json:"total_delivery_time,omitempty"`
}
