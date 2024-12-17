package models

import "time"

type PaymentSystemCustomer struct {
	ExternalID              string     `bson:"_id,omitempty"`
	PaymentSystemCustomerID string     `bson:"id"`
	CreatedAt               time.Time  `bson:"created_at"`
	Status                  string     `bson:"status"`
	Email                   string     `bson:"email,omitempty"`
	Phone                   string     `bson:"phone,omitempty"`
	PhoneCheckDate          string     `bson:"phone_check_date,omitempty"`
	Bin                     string     `bson:"bin,omitempty"`
	Name                    string     `bson:"name,omitempty"`
	Channel                 string     `bson:"channel"`
	Accounts                []Accounts `bson:"accounts"`
	CheckoutURL             string     `bson:"checkout_url"`
	AccessToken             string     `bson:"access_token"`
	CustomerAccessToken     string     `bson:"customer_access_token"`
	UpdatedAt               time.Time  `bson:"updated_at"`
	CustomerRestaurants     []string   `bson:"restaurants"`
	Cards                   []Card     `bson:"cards"`
}

type Accounts struct {
	ID         string      `bson:"id"`
	ShopID     string      `bson:"shop_id"`
	CustomerID string      `bson:"customer_id,omitempty"`
	Status     string      `bson:"status"`
	Name       string      `bson:"name,omitempty"`
	Amount     int64       `bson:"amount"`
	Currency   string      `bson:"currency"`
	Resources  []Resources `bson:"resources"`
	CreatedAt  time.Time   `bson:"created_at"`
	ExternalID string      `bson:"external_id,omitempty"`
}

type Resources struct {
	ID        string `bson:"id,omitempty"`
	Iban      string `bson:"iban,omitempty"`
	IsDefault bool   `bson:"is_default,omitempty"`
}

type Card struct {
	ID             string    `bson:"id"`
	Status         string    `bson:"status"`
	CreatedAt      time.Time `bson:"created_at"`
	PanMasked      string    `bson:"pan_masked"`
	ExpiryDate     string    `bson:"expiry_date"`
	PaymentSystem  string    `bson:"payment_system"`
	Emitter        string    `bson:"emitter"`
	CvcRequired    bool      `bson:"cvc_required"`
	MasterPassCard bool      `bson:"master_pass_card"`
	CustomerID     string    `bson:"-"`
	CustomerStatus string    `bson:"-"`
}
