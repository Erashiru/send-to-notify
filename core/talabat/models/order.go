package models

import "time"

type CancelOrderRequest struct {
	RemoteID      string `json:"-"`
	RemoteOrderID string `json:"-"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type CreateOrderResponse struct {
	RemoteResponse RemoteResponse `json:"remoteResponse"`
}

type RemoteResponse struct {
	RemoteOrderId string `json:"remoteOrderId"`
}
type CreateOrderRequest struct {
	Token              string             `json:"token"`
	Code               string             `json:"code"`
	Comments           Comments           `json:"comments"`
	CreatedAt          time.Time          `json:"createdAt"`
	Customer           Customer           `json:"customer"`
	Delivery           Delivery           `json:"delivery"`
	Discounts          []Discounts        `json:"discounts"`
	ExpeditionType     string             `json:"expeditionType"`
	ExpiryDate         string             `json:"expiryDate"`
	ExtraParameters    ExtraParameters    `json:"extraParameters"`
	LocalInfo          LocalInfo          `json:"localInfo"`
	Payment            Payment            `json:"payment"`
	Test               bool               `json:"test"`
	ShortCode          string             `json:"shortCode"`
	PreOrder           bool               `json:"preOrder"`
	Pickup             PickUp             `json:"pickup"`
	PlatformRestaurant PlatformRestaurant `json:"platformRestaurant"`
	Price              Price              `json:"price"`
	Products           []Product          `json:"products"`
	CorporateTaxID     string             `json:"corporateTaxId"`
	CallbackUrls       CallbackUrls       `json:"callbackUrls"`
}

type PickUp struct {
	PickupTime string `json:"pickupTime"`
	PickupCode string `json:"pickupCode"`
}

type Comments struct {
	CustomerComment string `json:"customerComment"`
}

type Customer struct {
	Email       string   `json:"email"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	MobilePhone string   `json:"mobilePhone"`
	Flags       []string `json:"flags"`
}

type Address struct {
	Postcode             int     `json:"postcode"`
	City                 string  `json:"city"`
	Street               string  `json:"street"`
	Number               string  `json:"number"`
	Longitude            float64 `json:"longitude"`
	Latitude             float64 `json:"latitude"`
	Intercom             string  `json:"intercom"`
	Floor                string  `json:"floor"`
	FlatNumber           string  `json:"flatNumber"`
	Entrance             string  `json:"entrance"`
	DeliveryMainArea     string  `json:"deliveryMainArea"`
	DeliveryInstructions string  `json:"deliveryInstructions"`
	DeliveryArea         string  `json:"deliveryArea"`
	Company              string  `json:"company"`
	Building             string  `json:"building"`
}

type Delivery struct {
	Address              Address `json:"address"`
	ExpectedDeliveryTime string  `json:"expectedDeliveryTime"`
	ExpressDelivery      bool    `json:"expressDelivery"`
	RiderPickupTime      string  `json:"riderPickupTime"`
}

type Discounts struct {
	Name   string `json:"name"`
	Amount string `json:"amount"`
	Type   string `json:"type"`
}

type ExtraParameters struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}

type LocalInfo struct {
	CountryCode    string `json:"countryCode"`
	CurrencySymbol string `json:"currencySymbol"`
	Platform       string `json:"platform"`
	PlatformKey    string `json:"platformKey"`
}

type Payment struct {
	Status string `json:"status"`
	Type   string `json:"type"`
}

type PlatformRestaurant struct {
	ID string `json:"id"`
}

type DeliveryFees struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type Price struct {
	DeliveryFees        []DeliveryFees `json:"deliveryFees"`
	GrandTotal          string         `json:"grandTotal"`
	PayRestaurant       string         `json:"payRestaurant"`
	RiderTip            string         `json:"riderTip"`
	TotalNet            string         `json:"totalNet"`
	VatTotal            string         `json:"vatTotal"`
	CollectFromCustomer string         `json:"collectFromCustomer"`
}

type SelectedToppings struct {
	Children   []Children `json:"children"`
	Name       string     `json:"name"`
	Price      string     `json:"price"`
	Quantity   int        `json:"quantity"`
	ID         string     `json:"id"`
	RemoteCode string     `json:"remoteCode"`
	Type       string     `json:"type"`
}

type Children struct {
	Children   []Children `json:"children"`
	Name       string     `json:"name"`
	Price      string     `json:"price"`
	Quantity   int        `json:"quantity"`
	ID         string     `json:"id"`
	RemoteCode any        `json:"remoteCode"`
	Type       string     `json:"type"`
}

type Variation struct {
	Name string `json:"name"`
}

type Product struct {
	CategoryName     string             `json:"categoryName"`
	Name             string             `json:"name"`
	PaidPrice        string             `json:"paidPrice"`
	Quantity         string             `json:"quantity"`
	RemoteCode       string             `json:"remoteCode"`
	SelectedToppings []SelectedToppings `json:"selectedToppings"`
	UnitPrice        string             `json:"unitPrice"`
	Comment          string             `json:"comment"`
	ID               string             `json:"id"`
	Variation        Variation          `json:"variation"`
}

type CallbackUrls struct {
	OrderAcceptedURL string `json:"orderAcceptedUrl"`
	OrderRejectedURL string `json:"orderRejectedUrl"`
	OrderPickedUpURL string `json:"orderPickedUpUrl"`
	OrderPreparedURL string `json:"orderPreparedUrl"`
}
