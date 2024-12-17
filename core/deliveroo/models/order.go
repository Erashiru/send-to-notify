package models

type Order struct {
	ID                   string        `json:"id"`
	OrderNumber          string        `json:"order_number"`
	LocationID           string        `json:"location_id"`
	DisplayID            string        `json:"display_id"`
	Status               string        `json:"status"`
	StatusLog            []StatusLog   `json:"status_log"`
	At                   string        `json:"at"`
	StatusObject         StatusObject  `json:"status_object"`
	FulfillmentType      string        `json:"fulfillment_type"`
	OrderNotes           string        `json:"order_notes"`
	CutleryNotes         string        `json:"cutlery_notes"`
	ASAP                 bool          `json:"asap"`
	PrepareFor           string        `json:"prepare_for"`
	TableNumber          string        `json:"table_number"`
	Subtotal             Price         `json:"subtotal"`
	Delivery             Delivery      `json:"delivery"`
	TotalPrice           Price         `json:"total_price"`
	PartnerOrderSubtotal Price         `json:"partner_order_subtotal"`
	PartnerOrderTotal    Price         `json:"partner_order_total"`
	OfferDiscount        Price         `json:"offer_discount"`
	CashDue              Price         `json:"cash_due"`
	BagFee               Price         `json:"bag_fee"`
	Surcharge            Price         `json:"surcharge"`
	Items                []Item        `json:"items"`
	StartPreparingAt     string        `json:"start_preparing_at"`
	ConfirmAt            string        `json:"confirm_at"`
	Promotions           []Promotion   `json:"promotions"`
	RemakeDetails        RemakeDetails `json:"remake_details"`
	OrderCost            float64       `json:"order_cost"`
	IsTabletless         bool          `json:"is_tabletless"`
	Customer             Customer      `json:"customer"`
}

type StatusLog struct {
	At     string `json:"at"`
	Status string `json:"status"`
}

type Object struct {
	At     string `json:"at"`
	Status string `json:"status"`
}

type StatusObject struct {
	Status string `json:"status"`
}

type Price struct {
	Fractional   int    `json:"fractional"`
	CurrencyCode string `json:"currency_code"`
}

type Delivery struct {
	DeliveryFee       Price    `json:"delivery_fee"`
	DeliveryNotes     string   `json:"delivery_notes"`
	Line1             string   `json:"line1"`
	Line2             string   `json:"line2"`
	City              string   `json:"city"`
	Postcode          string   `json:"postcode"`
	ContactNumber     string   `json:"contact_number"`
	ContactAccessCode string   `json:"contact_access_code"`
	DeliverBy         string   `json:"deliver_by"`
	CustomerName      string   `json:"customer_name"`
	Location          Location `json:"location"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Item struct {
	PosItemID       string     `json:"pos_item_id"`
	Quantity        int        `json:"quantity"`
	Name            string     `json:"name"`
	OperationalName string     `json:"operational_name"`
	UnitPrice       Price      `json:"unit_price"`
	TotalPrice      Price      `json:"total_price"`
	MenuUnitPrice   Price      `json:"menu_unit_price"`
	DiscountAmount  Price      `json:"discount_amount"`
	Modifiers       []Modifier `json:"modifiers"`
}

type Modifier struct {
	Name  string `json:"name"`
	Price Price  `json:"price"`
}

type Promotion struct {
	ID         string       `json:"id"`
	Type       string       `json:"type"`
	Value      int          `json:"value"`
	PosItemIDs []PosItemIds `json:"pos_item_ids"`
}

type RemakeDetails struct {
	ParentOrderID string  `json:"parent_order_id"`
	Fault         string  `json:"fault"`
	OrderCost     float64 `json:"order_cost"`
	IsTabletless  bool    `json:"is_tabletless"`
}

type Customer struct {
	FirstName            string `json:"first_name"`
	ContactNumber        string `json:"contact_number"`
	ContactAccessCode    string `json:"contact_access_code"`
	OrderFrequencyAtSite string `json:"order_frequency_at_site"`
}

type PosItemIds struct {
	Id string `json:"id"`
}
