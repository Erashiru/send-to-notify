package models_v2

type Venue struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	ExternalVenueId string `json:"external_venue_id"`
}

type BasketPrice struct {
	Total          Total          `json:"total"`
	PriceBreakdown PriceBreakdown `json:"price_breakdown"`
}

type PriceBreakdown struct {
	TotalBeforeDiscounts Total `json:"total_before_discounts"`
	TotalDiscounts       Total `json:"total_discounts"`
}

type Total struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

type Delivery struct {
	Status       string   `json:"status"`
	Type         string   `json:"type"`     // homedelivery, takeaway, eatin
	Time         string   `json:"time"`     // ISO 8601
	Location     Location `json:"location"` //  only for self-delivery venues
	SelfDelivery bool     `json:"self_delivery"`
}

type Coordinates struct {
	Lon float64 `json:"log"`
	Lat float64 `json:"lat"`
}

type Location struct {
	StreetAddress    string      `json:"street_address"`
	Apartment        string      `json:"apartment"`
	City             string      `json:"city"`
	Country          string      `json:"country"`
	Coordinates      Coordinates `json:"coordinates"`
	FormattedAddress string      `json:"formatted_address"`
}

type Fees struct {
	Total          Total         `json:"total"`
	PriceBreakdown FeesBreakDown `json:"price_breakdown"`
	Parts          []Part        `json:"parts"`
}

type FeesBreakDown struct {
	TotalBeforeDiscounts Total `json:"total_before_discounts"`
	TotalDiscounts       Total `json:"total_discounts"`
	Liability            Total `json:"liability"`
}

type Part struct {
	Type  string `json:"type"`
	Total Total  `json:"total"`
}

type Option struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Value      string      `json:"value"`
	Price      Total       `json:"price"`
	PosId      string      `json:"pos_id"`
	Count      int         `json:"count"`
	ValuePosId string      `json:"value_pos_id"`
	Deposit    interface{} `json:"deposit"`
}

type ItemPrice struct {
	UnitPrice      Total          `json:"unit_price"`
	Total          Total          `json:"total"`
	PriceBreakdown PriceBreakdown `json:"price_breakdown"`
}

type Category struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Item struct {
	ItemType      string      `json:"item_type"`
	Id            string      `json:"id"`
	Count         int         `json:"count"`
	PosId         string      `json:"pos_id"`
	Sku           string      `json:"sku"`  // only for retail venues
	Gtin          string      `json:"gtin"` // only for retail venues
	Options       []Option    `json:"options"`
	ItemPrice     ItemPrice   `json:"item_price"`
	Name          string      `json:"name"`
	Category      Category    `json:"category"`
	Deposit       interface{} `json:"deposit"`
	IsBundleOffer bool        `json:"is_bundle_offer"`
}

type Preorder struct {
	Time   string `json:"time"`   // ISO 8601
	Status string `json:"status"` // waiting, confirmed (waiting - waiting to be confirmed by the venue staff; confirmed - has been confirmed by venue staff)
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
	Id                  string      `json:"id"`
	Venue               Venue       `json:"venue"`
	BasketPrice         BasketPrice `json:"basket_price"`
	Delivery            Delivery    `json:"delivery"`
	Fees                Fees        `json:"fees"`
	Items               []Item      `json:"items"`
	CreatedAt           string      `json:"created_at"`
	ConsumerComment     string      `json:"consumer_comment"` // comment by client
	PickupEta           string      `json:"pickup_eta"`
	AttributionId       string      `json:"attribution_id"` // source where customer placed an order
	Type                string      `json:"type"`           // preorder, instant
	PreOrder            Preorder    `json:"pre_order"`      // only for type (preorder)
	ConsumerName        string      `json:"consumer_name"`
	ConsumerPhoneNumber string      `json:"consumer_phone_number"` // by settings Wolt
	OrderNumber         string      `json:"order_number"`
	OrderStatus         string      `json:"order_status"`
	ModifiedAt          string      `json:"modified_at"`
	CompanyTaxId        string      `json:"company_tax_id"`      // in some markets, consumers are allowed to bill their company for the purchase
	LoyaltyCardNumber   string      `json:"loyalty_card_number"` // only for retail venues
	CashPayment         CashPayment `json:"cash_payment"`
}
