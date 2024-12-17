package dto

type CreateCustomerRequest struct {
	ExternalId     string `json:"external_id,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Fingerprint    string `json:"fingerprint,omitempty"`
	PhoneCheckDate string `json:"phone_check_date,omitempty"`
	Channel        string `json:"channel"`
}

type CreateCustomerResponse struct {
	Customer            Customer `json:"customer"`
	CustomerAccessToken string   `json:"customer_access_token"`
}

type GetCustomerByIdResponse struct {
	ID          string     `json:"id"`
	CreatedAt   string     `json:"created_at"`
	Status      string     `json:"status"`
	ExternalID  string     `json:"external_id"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Accounts    []Accounts `json:"accounts"`
	CheckoutURL string     `json:"checkout_url"`
	AccessToken string     `json:"access_token"`
}

type Resources struct {
	ID        string `json:"id,omitempty"`
	Iban      string `json:"iban,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
}
type Accounts struct {
	ID         string      `json:"id"`
	ShopID     string      `json:"shop_id"`
	CustomerID string      `json:"customer_id,omitempty"`
	Status     string      `json:"status"`
	Name       string      `json:"name,omitempty"`
	Amount     int64       `json:"amount"`
	Currency   string      `json:"currency"`
	Resources  []Resources `json:"resources"`
	CreatedAt  string      `json:"created_at"`
	ExternalID string      `json:"external_id,omitempty"`
}
type Customer struct {
	ID          string     `json:"id"`
	CreatedAt   string     `json:"created_at"`
	Status      string     `json:"status"`
	ExternalID  string     `json:"external_id,omitempty"`
	Email       string     `json:"email,omitempty"`
	Phone       string     `json:"phone,omitempty"`
	Accounts    []Accounts `json:"accounts"`
	CheckoutURL string     `json:"checkout_url"`
	AccessToken string     `json:"access_token"`
}

type GetCustomerCardsResponse struct {
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
