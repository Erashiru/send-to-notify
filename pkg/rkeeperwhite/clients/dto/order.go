package dto

import (
	"time"
)

type CreateOrderParam struct {
	Async Sync  `json:"async"`
	Order Order `json:"order"`
}

type Order struct {
	OriginalOrderId   string               `json:"originalOrderId,omitempty"`
	Customer          *PersonInfo          `json:"customer,omitempty"`
	Payment           Payment              `json:"payment,omitempty"`
	ExpeditionType    string               `json:"expeditionType,omitempty"`
	Pickup            PickUp               `json:"pickup,omitempty"`
	Delivery          Delivery             `json:"delivery,omitempty"`
	Products          []CreateOrderProduct `json:"products"`
	Comment           string               `json:"comment,omitempty"`
	Price             *Price               `json:"price,omitempty"`
	PersonsQuantity   int                  `json:"personsQuantity,omitempty"`
	TableCode         int                  `json:"tableCode,omitempty"`
	OrderCategoryCode int                  `json:"orderCategoryCode,omitempty"`
	OrderTypeCode     int                  `json:"orderTypeCode,omitempty"`
	PrePayments       []PrePayment         `json:"prePayments,omitempty"`
}

type PrePayment struct {
	Amount   int    `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type Payment struct {
	Type string `json:"type,omitempty"`
}

type PickUp struct {
	Courier      *PersonInfo `json:"courier,omitempty"`
	ExpectedTime time.Time   `json:"expectedTime,omitempty"`
	Taker        string      `json:"taker,omitempty"`
}

type PersonInfo struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}

type Price struct {
	Total int `json:"total"`
}

type Delivery struct {
	ExpectedTime time.Time        `json:"expectedTime,omitempty"`
	Address      *DeliveryAddress `json:"address,omitempty"`
}

type Name struct {
	Name string `json:"name,omitempty"`
}

type Coordinate struct {
	Latitude  int `json:"latitude,omitempty"`
	Longitude int `json:"longitude,omitempty"`
}

type DeliveryAddress struct {
	FullAddress string     `json:"fullAddress,omitempty"`
	Subway      string     `json:"subway,omitempty"`
	Region      string     `json:"region,omitempty"`
	City        Name       `json:"city,omitempty"`
	Street      Name       `json:"street,omitempty"`
	HouseNumber string     `json:"houseNumber,omitempty"`
	FlatNumber  string     `json:"flatNumber,omitempty"`
	Entrance    string     `json:"entrance,omitempty"`
	Intercom    string     `json:"intercom,omitempty"`
	Floor       string     `json:"floor,omitempty"`
	Coordinates Coordinate `json:"coordinates,omitempty"`
}

type CreateOrderProduct struct {
	Id          string                  `json:"id"`
	Name        string                  `json:"name,omitempty"`
	Price       string                  `json:"price,omitempty"`
	Quantity    int                     `json:"quantity,omitempty"`
	Ingredients []CreateOrderIngredient `json:"ingredients,omitempty"`
}

type CreateOrderIngredient struct {
	Id       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Quantity int    `json:"quantity"`
	Price    string `json:"price,omitempty"`
}

type GetOrderRequest struct {
	TaskType string `json:"taskType"`
	Params   Params `json:"params"`
}

type GetOrderTaskResponse struct {
	TaskResponse   OrderResponse  `json:"taskResponse"`
	ResponseCommon ResponseCommon `json:"responseCommon"`
}

type OrderResponse struct {
	OrderResponseBody OrderResponseBody `json:"order"`
}

type OrderResponseBodyStatus struct {
	Value         string `json:"value"`
	IsBillPrinted bool   `json:"isBillPrinted"`
}

type OrderResponseBodyProducts struct {
	Ingredients []interface{} `json:"ingredients"`
	Id          int           `json:"id"`
	Name        string        `json:"name"`
	Quantity    int           `json:"quantity"`
	Price       int           `json:"price"`
}

type OrderResponseBody struct {
	OriginalOrderId string                      `json:"originalOrderId"`
	CreatedAt       string                      `json:"createdAt"`
	Status          OrderResponseBodyStatus     `json:"status"`
	Products        []OrderResponseBodyProducts `json:"products"`
	Comment         string                      `json:"comment"`
	Price           Price                       `json:"price"`
	AppliedPayments []interface{}               `json:"appliedPayments"`
	PersonsQuantity int                         `json:"personsQuantity"`
	TableCode       int                         `json:"tableCode"`
	WaiterId        int                         `json:"waiterId"`
	SubState        string                      `json:"substate"`
	DiscountIds     []interface{}               `json:"discountIds"`
}
