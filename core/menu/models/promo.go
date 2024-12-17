package models

import "time"

type Promo struct {
	Id                       string               `bson:"_id" json:"id"`
	Name                     string               `bson:"name" json:"name"`
	Type                     string               `bson:"type" json:"type"`
	DeliveryPrice            int                  `bson:"delivery_price" json:"delivery_price"`
	StartDate                time.Time            `bson:"start_date" json:"start_date"`
	EndDate                  time.Time            `bson:"end_date" json:"end_date"`
	StartTime                string               `bson:"start_time" json:"start_time"`
	EndTime                  string               `bson:"end_time" json:"end_time"`
	RestaurantIds            []string             `bson:"restaurant_ids" json:"restaurant_ids"`
	ProductIds               []string             `bson:"product_ids" json:"product_ids"`
	Percent                  int                  `bson:"percent" json:"percent"`
	Budget                   int                  `bson:"budget" json:"budget"`
	IsPosIntegrated          bool                 `bson:"is_pos_integrated" json:"is_pos_integrated"`
	PosDiscount              PosDiscount          `bson:"pos_discount" json:"pos_discount"`
	PosType                  string               `bson:"pos_type" json:"pos_type"`
	DeliveryService          string               `bson:"delivery_service" json:"delivery_service"`
	IsActive                 bool                 `bson:"is_active" json:"is_active"`
	CreatedAt                time.Time            `bson:"created_at" json:"created_at"`
	ProductGifts             []ProductGift        `bson:"product_gifts" json:"product_gifts"`
	CategoryPercentage       []CategoryPercentage `bson:"category_percentage" json:"category_percentage"`
	CategoryFixed            []CategoryFixed      `bson:"category_fixed" json:"category_fixed"`
	OrderDiscount            []OrderDiscount      `bson:"order_discounts" json:"order_discount"`
	PercentageForEachProduct bool                 `bson:"percentage_for_each_product" json:"percentage_for_each_product"`
	ProductsPercentage       []ProductPercentage  `bson:"products_percentage" json:"products_percentage"`
}

type OrderDiscount struct {
	Type       string `bson:"type" json:"type"`
	DiscountID string `bson:"discount_id" json:"discount_id"`
	Amount     int    `bson:"amount" json:"amount"`
}

type CategoryPercentage struct {
	CategoryId string `bson:"category_id" json:"category_id"`
	Percent    int    `bson:"percent" json:"percent"`
	DiscountId string `bson:"discount_id" json:"discount_id"`
}

type CategoryFixed struct {
	CategoryId string `bson:"category_id" json:"category_id"`
	Sum        int    `bson:"sum" json:"sum"`
	DiscountId string `bson:"discount_id" json:"discount_id"`
}

type ProductGift struct {
	ProductId string `bson:"product_id" json:"product_id"`
	PromoId   string `bson:"promo_id" json:"promo_id"`
}

type PosDiscount struct {
	DiscountId string `bson:"discount_id" json:"discount_id"`
	Type       string `bson:"type" json:"type"`
}

type ProductPercentage struct {
	ProductID string `bson:"product_id" json:"product_id"`
	Percent   int    `bson:"percent" json:"percent"`
}
