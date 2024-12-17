package models

import (
	"time"
)

type Cart struct {
	ID                      string          `bson:"_id,omitempty" json:"id,omitempty"`
	Table                   string          `json:"table" bson:"table"`
	RestaurantID            string          `json:"restaurant_id" bson:"restaurant_id"`
	Customer                Customer        `json:"customer" bson:"customer"`
	Currency                string          `json:"currency" bson:"currency"`
	MenuID                  string          `json:"menu_id" bson:"menu_id"`
	CreatedAt               time.Time       `json:"created_at" bson:"created_at"`
	ModifiedAt              time.Time       `json:"modified_at" bson:"modified_at"`
	Items                   []CartProduct   `json:"items" bson:"items"`
	TotalSum                float64         `json:"total_sum" bson:"total_sum"`
	DeliveryAddress         DeliveryAddress `json:"delivery_address" bson:"delivery_address"`
	OrderComment            string          `json:"order_comment" bson:"order_comment"`
	IsPickedUpByCustomer    bool            `json:"is_picked_up_by_customer" bson:"is_picked_up_by_customer"`
	PreOrderTime            time.Time       `json:"pre_order_time" bson:"pre_order_time"`
	PaymentType             string          `json:"payment_type" bson:"payment_type"`
	PaymentTypeID           string          `json:"payment_type_id" bson:"payment_type_id"`     // kwaaka admin case
	PaymentTypeKind         string          `json:"payment_type_kind" bson:"payment_type_kind"` // kwaaka admin case
	Delivery                Delivery        `json:"delivery" bson:"delivery"`
	Status                  string          `json:"status" bson:"status"`
	BackUrl                 string          `json:"back_url" bson:"back_url"`
	SuccessUrl              string          `json:"success_url" bson:"success_url"`
	FailureUrl              string          `json:"failure_url" bson:"failure_url"`
	PaymentSystem           string          `json:"payment_system" bson:"payment_system"`
	CardID                  string          `json:"card_id" bson:"card_id"`
	PaymentSystemCustomerID string          `json:"payment_system_customer_id" bson:"payment_system_customer_id"`
}

type CartProduct struct {
	ProductID   string      `json:"product_id,omitempty" bson:"product_id"`
	Quantity    int         `json:"quantity" bson:"quantity"`
	Price       float64     `json:"price" bson:"price"`
	Name        string      `json:"name,omitempty" bson:"name"`
	Attributes  []Attribute `json:"attributes,omitempty" bson:"attributes,omitempty"`
	CookingTime int32       `json:"cooking_time" bson:"cooking_time"`
}

type Attribute struct {
	AttributeID string  `json:"attribute_id" bson:"attribute_id"` // attr id in pos
	Name        string  `json:"name" bson:"name"`
	Price       float64 `json:"price" bson:"price"`
	Quantity    int     `json:"quantity,omitempty" bson:"quantity"`
}

type Delivery struct {
	Dispatcher                 string  `json:"dispatcher" bson:"dispatcher"`
	DeliveryTime               int32   `json:"delivery_time" bson:"delivery_time"`
	ClientDeliveryPrice        float64 `json:"client_delivery_price" bson:"client_delivery_price"`
	FullDeliveryPrice          float64 `json:"full_delivery_price" bson:"full_delivery_price"`
	DropOffScheduleTime        string  `json:"drop_off_schedule_time" bson:"drop_off_schedule_time"`
	KwaakaChargedDeliveryPrice float64 `json:"delivery_service_fee" bson:"kwaaka_charged_delivery_price"`
}

// ToDo, Delete after a couple of months. It's used for excel report handler
type OldCart struct {
	ID                      string          `bson:"_id,omitempty" json:"id,omitempty"`
	Table                   string          `json:"table" bson:"table"`
	RestaurantID            string          `json:"restaurant_id" bson:"restaurant_id"`
	Customer                Customer        `json:"customer" bson:"customer"`
	Currency                string          `json:"currency" bson:"currency"`
	MenuID                  string          `json:"menu_id" bson:"menu_id"`
	CreatedAt               time.Time       `json:"created_at" bson:"created_at"`
	ModifiedAt              time.Time       `json:"modified_at" bson:"modified_at"`
	Items                   []CartProduct   `json:"items" bson:"items"`
	TotalSum                float64         `json:"total_sum" bson:"total_sum"`
	DeliveryAddress         DeliveryAddress `json:"delivery_address" bson:"delivery_address"`
	OrderComment            string          `json:"order_comment" bson:"order_comment"`
	IsPickedUpByCustomer    bool            `json:"is_picked_up_by_customer" bson:"is_picked_up_by_customer"`
	PreOrderTime            time.Time       `json:"pre_order_time" bson:"pre_order_time"`
	PaymentType             string          `json:"payment_type" bson:"payment_type"`
	Delivery                Delivery        `json:"delivery" bson:"delivery"`
	Status                  string          `json:"status" bson:"status"`
	BackUrl                 string          `json:"back_url" bson:"back_url"`
	SuccessUrl              string          `json:"success_url" bson:"success_url"`
	FailureUrl              string          `json:"failure_url" bson:"failure_url"`
	PaymentSystem           string          `json:"payment_system" bson:"paymentsystem"`
	CardID                  string          `json:"card_id" bson:"card_id"`
	PaymentSystemCustomerID string          `json:"payment_system_customer_id" bson:"payment_system_customer_id"`
}
