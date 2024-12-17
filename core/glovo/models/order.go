package models

type Order struct {
	OrderID             string `json:"order_id"`
	StoreID             string `json:"store_id"`
	OrderTime           string `json:"order_time"`
	EstimatedPickupTime string `json:"estimated_pickup_time"`
	UtcOffsetMinutes    string `json:"utc_offset_minutes"`
	PaymentMethod       string `json:"payment_method"`
	Currency            string `json:"currency"`
	OrderCode           string `json:"order_code"`
	AllergyInfo         string `json:"allergy_info"`
	SpecialRequirements string `json:"special_requirements"`

	EstimatedTotalPrice            int  `json:"estimated_total_price"`
	DeliveryFee                    *int `json:"delivery_fee"`
	MinimumBasketSurcharge         int  `json:"minimum_basket_surcharge"`
	CustomerCashPaymentAmount      int  `json:"customer_cash_payment_amount"`
	PartnerDiscountsProducts       int  `json:"partner_discounts_products"`
	PartnerDiscountedProductsTotal int  `json:"partner_discounted_products_total"`

	GlovoDiscountsProducts  int `json:"glovo_discounts_products"`
	DiscountedProductsTotal int `json:"discounted_products_total"`

	TotalCustomerToPay int `json:"total_customer_to_pay"`

	Courier              Courier         `json:"courier"`
	Customer             Customer        `json:"customer"`
	Products             []ProductOrder  `json:"products"`
	DeliveryAddress      DeliveryAddress `json:"delivery_address"`
	BundledOrders        []string        `json:"bundled_orders"`
	PickUpCode           string          `json:"pick_up_code"`
	IsPickedUpByCustomer bool            `json:"is_picked_up_by_customer"`
	CutleryRequested     bool            `json:"cutlery_requested"`

	LoyaltyCard string `json:"loyalty_card"`

	Timezone string `json:"-"`
}

type Courier struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type Customer struct {
	Name             string           `json:"name"`
	PhoneNumber      string           `json:"phone_number"`
	Hash             string           `json:"hash"`
	InvoicingDetails InvoicingDetails `json:"invoicing_details"`
}
type InvoicingDetails struct {
	CompanyName    string `json:"company_name"`
	CompanyAddress string `json:"company_address"`
	TaxID          string `json:"tax_id"`
}

type ProductOrder struct {
	ID                 string           `json:"id"`
	PurchasedProductID string           `json:"purchased_product_id"`
	Name               string           `json:"name"`
	Price              int              `json:"price"`
	Quantity           int              `json:"quantity"`
	Attributes         []AttributeOrder `json:"attributes"`
}

type AttributeOrder struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type DeliveryAddress struct {
	Label     string  `json:"label"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Replacement struct {
	PurchasedProductID string       `json:"purchased_product_id"`
	Product            ProductOrder `json:"product"`
}

type AddedProduct struct {
	ID         string           `json:"id"`
	Quantity   int              `json:"quantity"`
	Attributes []AttributeOrder `json:"attributes"`
}

type ModifyOrderProductRequest struct {
	ID               int64          `json:"-"`
	Replacements     []Replacement  `json:"replacements" binding:"required"`
	RemovedPurchases []string       `json:"removed_purchases" binding:"required"`
	AddedProducts    []AddedProduct `json:"added_products"`
}

type StoreStatus struct {
	Until string `json:"until"`
}

type OrderTypes string

const (
	INSTANT  OrderTypes = "INSTANT"
	PREORDER OrderTypes = "PREORDER"
)

func (o OrderTypes) String() string {
	return string(o)
}

type DeliveryService string

const (
	GLOVO DeliveryService = "glovo"
)

func (d DeliveryService) String() string {
	return string(d)
}
