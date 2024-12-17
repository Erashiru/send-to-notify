package models

type CustomerCards struct {
	ID            string `json:"id"`
	CustomerID    string `json:"customer_id"`
	CreatedAt     string `json:"created_at"`
	PanMasked     string `json:"pan_masked"`
	ExpiryDate    string `json:"expiry_date"`
	Holder        string `json:"holder"`
	PaymentSystem string `json:"payment_system"`
	Emitter       string `json:"emitter"`
	CvcRequired   bool   `json:"cvc_required"`
	Error         Error  `json:"error"`
	Action        Action `json:"action"`
}
