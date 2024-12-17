package dto

type ChangeSubscriptionStatusRequest struct {
	Status string `json:"status"`
}
type CreateSubscriptionRequest struct {
	CustomerId  string `json:"customer_id"`
	CardId      string `json:"card_id"`
	Amount      int    `json:"amount"`
	Currency    string `json:"currency,omitempty"`
	Description string `json:"description,omitempty"`
	ExtraInfo   string `json:"extra_info,omitempty"`
	NextPay     string `json:"next_pay"`
	Step        int    `json:"step"`
	Unit        string `json:"unit"`
}

type CreateSubscriptionResponse struct {
	ID          string    `json:"id"`
	CreatedAt   string    `json:"created_at"`
	Amount      int       `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description,omitempty"`
	ExtraInfo   ExtraInfo `json:"extra_info,omitempty"`
	Payer       Payer     `json:"payer"`
	Schedule    Schedule  `json:"schedule"`
}

type Payer struct {
	Type          string `json:"type"`
	PanMasked     string `json:"pan_masked,omitempty"`
	ExpiryDate    string `json:"expiry_date,omitempty"`
	Holder        string `json:"holder,omitempty"`
	PaymentSystem string `json:"payment_system,omitempty"`
	Emitter       string `json:"emitter,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	CustomerID    string `json:"customer_id,omitempty"`
	CardID        string `json:"card_id,omitempty"`
}
type Schedule struct {
	Status  string `json:"status"`
	NextPay string `json:"next_pay"`
	Step    int    `json:"step"`
	Unit    string `json:"unit"`
}
