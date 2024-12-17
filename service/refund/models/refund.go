package models

import "time"

type Refund struct {
	ID            string    `json:"_id" bson:"_id"`
	Amount        int       `json:"amount" bson:"amount"`
	Reason        string    `json:"reason" bson:"reason"`
	OrderID       string    `json:"order_id" bson:"order_id"`
	PaymentID     string    `json:"payment_id" bson:"payment_id"`
	PaymentSystem string    `json:"payment_system" bson:"payment_system"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}
