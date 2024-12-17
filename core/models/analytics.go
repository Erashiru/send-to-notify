package models

import (
	"fmt"
	"time"
)

type GetPopularProduct struct {
	RestaurantIDs []string
	StartDate     time.Time
	EndDate       time.Time
	Aggregator    []string
	SortBy        string
	Direction     int
	Limit         int
}

type OrderTotalAmount struct {
	RestaurantIDs []string
	StartDate     time.Time
	EndDate       time.Time
	Aggregator    []string
	Success       bool
}

type GetTotalOrdersByDate struct {
	RestaurantIDs []string
	StartDate     time.Time
	EndDate       time.Time
}

type TotalOrdersByDate struct {
	Date             `bson:"_id" json:"id"`
	DeliveryServices []ByDeliveryService `bson:"delivery_services" json:"delivery_services"`
}

type Date struct {
	Year  int `bson:"year" json:"year"`
	Month int `bson:"month" json:"month"`
	Day   int `bson:"day" json:"day"`
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

type ByDeliveryService struct {
	Service      string `bson:"service" json:"service"`
	TotalAmount  int    `bson:"total_amount" json:"total_amount"`
	TotalOrders  int    `bson:"total_orders" json:"total_orders"`
	CurrencyCode string `bson:"currency_code" json:"currency_code"`
}

func (ot *OrderTotalAmount) SetPreviousPeriod() {
	daysDifference := int(ot.EndDate.Sub(ot.StartDate).Hours() / 24)
	ot.EndDate = ot.StartDate
	ot.StartDate = ot.StartDate.AddDate(0, 0, -daysDifference)
}

type OrderTotalAmountQuantity struct {
	TotalAmount  int    `bson:"total_amount" json:"total_amount"`
	TotalOrders  int    `bson:"total_orders" json:"total_orders"`
	CurrencyCode string `bson:"currency_code" json:"currency_code"`
}

type AnalyticsTotalAmountQuantity struct {
	Current   OrderTotalAmountQuantity
	Previous  OrderTotalAmountQuantity
	Cancelled OrderTotalAmountQuantity
}
