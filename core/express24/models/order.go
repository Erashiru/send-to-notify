package models

import (
	"time"
)

type DeliveryService string

const (
	EXPRESS24 DeliveryService = "express24"
)

func (d DeliveryService) String() string {
	return string(d)
}

type Event struct {
	EventId      string `json:"event_id"`
	NewOrder     *Order `json:"new_order"`
	OrderChanged *Order `json:"order_changed"`
}

type Order struct {
	Id           int       `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
	Payment      Payment   `json:"payment"`
	Store        Store     `json:"store"`
	Delivery     Delivery  `json:"delivery"`
	User         User      `json:"user"`
	Products     []Product `json:"products"`
	OrderComment string    `json:"order_comment"`
	TotalPrice   float64   `json:"total_price"`
}

type Product struct {
	Id         int     `json:"id"`
	ExternalId string  `json:"external_id"`
	Name       string  `json:"name"`
	Qty        int     `json:"qty"`
	Params     []Param `json:"params"`
	Price      float64 `json:"price"`
}

type User struct {
	Name string `json:"name"`
}

type Delivery struct {
	Type    int     `json:"type"`
	Address Address `json:"address"`
	Comment string  `json:"comment"`
	Price   float64 `json:"price"`
}

type Address struct {
	Text   string  `json:"text"`
	Region Region  `json:"region"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
}

type Region struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type Store struct {
	ExternalId string `json:"external_id"`
	Branch     Branch `json:"branch"`
}

type Branch struct {
	ExternalId string `json:"external_id"`
}
type Payment struct {
	Id     string `json:"id"`
	Status int    `json:"status"`
}
type Param struct {
	Id         int      `json:"id"`
	ExternalId string   `json:"external_id"`
	Name       string   `json:"name"`
	Options    []Option `json:"options"`
}

type Option struct {
	Id         int     `json:"id"`
	ExternalId string  `json:"external_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
}
