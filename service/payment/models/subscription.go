package models

import "time"

type PaymentSystemSubscription struct {
	ID             string    `bson:"_id,omitempty"`
	SubscriptionID string    `bson:"id"`
	CreatedAtPS    time.Time `bson:"created_at_ps"`
	Amount         int       `bson:"amount"`
	Currency       string    `bson:"currency"`
	Description    string    `bson:"description"`
	ExtraInfo      string    `bson:"extra_info"`
	Payer          Payer     `bson:"payer"`
	Schedule       Schedule  `bson:"schedule"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"update_at"`
	PaymentSystem  string    `bson:"payment_system"`
}

type Payer struct {
	Type          string `json:"type"`
	PanMasked     string `json:"pan_masked"`
	ExpiryDate    string `json:"expiry_date"`
	Holder        string `json:"holder"`
	PaymentSystem string `json:"payment_system"`
	Emitter       string `json:"emitter"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	CustomerID    string `json:"customer_id"`
	CardID        string `json:"card_id"`
}
type Schedule struct {
	Status  string `json:"status"`
	NextPay string `json:"next_pay"`
	Step    int    `json:"step"`
	Unit    string `json:"unit"`
}
