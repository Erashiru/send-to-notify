package models

import "time"

type Order struct {
	ID                   string          `json:"id"`
	Table                string          `json:"table"`
	RestaurantID         string          `json:"restaurant_id"`
	Customer             Customer        `json:"customer"`
	DeliveryAddress      DeliveryAddress `json:"delivery_address"`
	Currency             string          `json:"currency"`
	Items                []Product       `json:"items"`
	TotalSum             float64         `json:"total_sum"`
	OrderComment         string          `json:"order_comment"`
	IsPickedUpByCustomer bool            `json:"is_picked_up_by_customer"`
	PreOrderTime         time.Time       `json:"pre_order_time"`
	PaymentType          string          `json:"payment_type"`
	Delivery             Delivery        `json:"delivery"`
	PromoCode            string          `json:"promo_code"`
}

type Delivery struct {
	Dispatcher                 string  `json:"dispatcher"`
	DeliveryTime               int32   `json:"delivery_time"`
	ClientDeliveryPrice        float64 `json:"client_delivery_price"`
	FullDeliveryPrice          float64 `json:"full_delivery_price"`
	KwaakaChargedDeliveryPrice float64 `json:"delivery_service_fee"`
	DropOffScheduleTime        string  `json:"drop_off_schedule_time"`
}

type Product struct {
	ProductID   string      `json:"product_id,omitempty"`
	Quantity    int         `json:"quantity"`
	Price       float64     `json:"price"`
	Name        string      `json:"name,omitempty"`
	Attributes  []Attribute `json:"attributes,omitempty"`
	CookingTime int32       `json:"cooking_time"`
}

type Attribute struct {
	AttributeID string  `json:"attribute_id"` // attr id in pos
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type Customer struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
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
type Coordinates struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}
