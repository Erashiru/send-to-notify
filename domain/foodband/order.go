package foodband

type CreateOrderRequest struct {
	StoreID string `json:"store_id"`
	Order   Order  `json:"order"`
}

type CancelOrderRequest struct {
	StoreID         string `json:"store_id"`
	OrderID         string `json:"order_id"`
	DeliveryService string `json:"delivery_service"`
	CancelReason    string `json:"cancel_reason"`
	PaymentStrategy string `json:"payment_strategy"`
}

type Order struct {
	ID                   string        `json:"id"`
	Type                 string        `json:"type"`
	Code                 string        `json:"code"`
	PickUpCode           string        `json:"pick_up_code"`
	CompleteBefore       string        `json:"complete_before"`
	Phone                string        `json:"phone"`
	DeliveryService      string        `json:"delivery_service"`
	DeliveryPoint        DeliveryPoint `json:"delivery_point"`
	Comment              string        `json:"comment"`
	Customer             Customer      `json:"customer"`
	Courier              Courier       `json:"courier"`
	Products             []Product     `json:"products"`
	Payments             []Payment     `json:"payments"`
	DeliveryFee          float64       `json:"delivery_fee"`
	DeliveryProviderType string        `json:"delivery_provider_type"`
}

type Courier struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type DeliveryPoint struct {
	Coordinates  Coordinates `json:"coordinates"`
	AddressLabel string      `json:"address_label"`
	Comment      string      `json:"comment"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Customer struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type Product struct {
	ID         string      `json:"id"`
	Quantity   int         `json:"quantity"`
	Price      float64     `json:"price"`
	Comment    string      `json:"comment"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	ID       string  `json:"id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type Payment struct {
	PaymentTypeKind string  `json:"payment_type_kind"`
	Sum             float64 `json:"sum"`
}
