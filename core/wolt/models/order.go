package models

import "time"

type DeliveryService string

const (
	WOLT        DeliveryService = "wolt"
	PAY_NOTHING string          = "PAY_NOTHING"
	USER_ERROR  string          = "USER_ERROR"
)

func (d DeliveryService) String() string {
	return string(d)
}

type OrderTypes string

func (o OrderTypes) String() string {
	return string(o)
}

type Venus struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Price struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type Coordinates struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type Location struct {
	StreetAddress    string      `json:"street_address"`
	Apartment        string      `json:"apartment"`
	City             string      `json:"city"`
	Country          string      `json:"country"`
	Coordinates      Coordinates `json:"coordinates"`
	FormattedAddress string      `json:"formatted_address"`
}

type Delivery struct {
	Status       string    `json:"status"`
	Type         string    `json:"type"`
	Time         time.Time `json:"time"`
	Fee          Price     `json:"fee"`
	Location     Location  `json:"location"`
	SelfDelivery bool      `json:"self_delivery"`
}
type CategoryOrder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type OptionOrder struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	Price      Price  `json:"price"`
	PosID      string `json:"pos_id"`
	Count      int    `json:"count"`
	ValuePosID string `json:"value_pos_id"`
}

type OrderItem struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Count      int           `json:"count"`
	Category   CategoryOrder `json:"category"`
	PosID      string        `json:"pos_id"`
	Options    []OptionOrder `json:"options"`
	TotalPrice Price         `json:"total_price"`
	BasePrice  Price         `json:"base_price"`
	UnitPrice  Price         `json:"unit_price"`
}
type PreOrder struct {
	Time   time.Time `json:"preorder_time"`
	Status string    `json:"pre_order_status"`
}

type CashPayment struct {
	CashAmount   CashAmount   `json:"cash_amount"`
	CashToExpect CashToExpect `json:"cash_to_expect"`
}
type CashAmount struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type CashToExpect struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}
type Order struct {
	ID                  string      `json:"id"`
	Venue               Venus       `json:"venue"`
	Price               Price       `json:"price"`
	Delivery            Delivery    `json:"delivery"`
	Items               []OrderItem `json:"items"`
	CreatedAt           time.Time   `json:"created_at"`
	ConsumerComment     string      `json:"consumer_comment"`
	PickupEta           time.Time   `json:"pickup_eta"`
	AttributionID       string      `json:"attribution_id"`
	Type                string      `json:"type"`
	PreOrder            PreOrder    `json:"pre_order"`
	ConsumerName        string      `json:"consumer_name"`
	ConsumerPhoneNumber string      `json:"consumer_phone_number"`
	OrderNumber         string      `json:"order_number"`
	OrderStatus         string      `json:"order_status"`
	ModifiedAt          time.Time   `json:"modified_at"`
	Timezone            string      `json:"-"`
	CashPayment         CashPayment `json:"cash_payment"`
}

type OrderNotification struct {
	Id        string                `json:"id"`
	Type      string                `json:"type"`
	CreatedAt time.Time             `json:"created_at"`
	Body      OrderNotificationBody `json:"order"`
}

type OrderNotificationBody struct {
	Id          string `json:"id"`
	ResourceUrl string `json:"resource_url"`
	Status      string `json:"status"`
	VenueId     string `json:"venue_id"`
}

type AcceptOrderRequest struct {
	ID         string    `json:"-"`
	PickupTime time.Time `json:"adjusted_pickup_time"`
}

type RejectOrderRequest struct {
	ID     string `json:"-"`
	Reason string `json:"reason"`
}
