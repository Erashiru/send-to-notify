package models

import "time"

type Order struct {
	ID                   string          `json:"id"`
	RestaurantID         string          `json:"restaurant_id"`
	Customer             Customer        `json:"customer"`
	DeliveryAddress      DeliveryAddress `json:"delivery_address"`
	Currency             string          `json:"currency"`
	Items                []Item          `json:"items"`
	TotalSum             float64         `json:"total_sum"`
	OrderComment         string          `json:"order_comment"`
	IsPickedUpByCustomer bool            `json:"is_picked_up_by_customer"`
	PaymentType          PaymentType     `json:"payment_type"`
	PreOrderTime         time.Time       `json:"pre_order_time"`
	Delivery             Delivery        `json:"delivery"`
	OperatorID           string          `json:"operator_id"`
	OperatorName         string          `json:"operator_name"`
	Discount             Discount        `json:"discount"`
}

type Delivery struct {
	Dispatcher                 string  `json:"dispatcher" bson:"dispatcher"`
	DeliveryTime               int32   `json:"delivery_time" bson:"delivery_time"`
	ClientDeliveryPrice        float64 `json:"client_delivery_price" bson:"client_delivery_price"`
	FullDeliveryPrice          float64 `json:"full_delivery_price" bson:"full_delivery_price"`
	KwaakaChargedDeliveryPrice float64 `json:"delivery_service_fee" bson:"kwaaka_charged_delivery_price"`
	DropOffScheduleTime        string  `json:"drop_off_schedule_time" bson:"drop_off_schedule_time"`
}

type PaymentType struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}
type Customer struct {
	Name                string `json:"name"`
	PhoneNumber         string `json:"phone_number"`
	Email               string `json:"email"`
	WhatsAppPhoneNumber string `json:"whats_app_phone_number"`
}
type Coordinates struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
type DeliveryAddress struct {
	Street       string      `json:"street"`
	Apartment    string      `json:"apartment"`
	Floor        string      `json:"floor"`
	City         string      `json:"city"`
	LocationType string      `json:"location_type"`
	BuildingName string      `json:"building_name"`
	Entrance     string      `json:"entrance"`
	DoorBellInfo string      `json:"door_bell_info"`
	Coordinates  Coordinates `json:"coordinates"`
	Comment      string      `json:"comment"`
}
type Attribute struct {
	AttributeID string  `json:"attribute_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}
type Item struct {
	ProductID   string      `json:"product_id"`
	Quantity    int         `json:"quantity"`
	Price       float64     `json:"price"`
	Name        string      `json:"name"`
	Attributes  []Attribute `json:"attributes,omitempty"`
	CookingTime int32       `json:"cooking_time"`
}

type Discount struct {
	Type           string `json:"type"`
	IikoDiscountId string `json:"iiko_discount_id"`
	Percent        int    `json:"percent"`
}
