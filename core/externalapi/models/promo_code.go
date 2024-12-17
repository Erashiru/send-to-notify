package models

import (
	"github.com/kwaaka-team/orders-core/core/qrmenu/models"
	"time"
)

type PromoCode struct {
	ID                string           `json:"id" bson:"_id,omitempty"`
	Code              string           `json:"code" bson:"code"`
	Name              string           `json:"name" bson:"name"`
	Description       string           `json:"description" bson:"description"`
	Link              string           `json:"link" bson:"link"`
	RestaurantIDs     []string         `json:"restaurant_ids" bson:"restaurant_ids"`
	UsageTime         int              `json:"usage_time" bson:"usage_time"`
	DeliveryType      []string         `json:"delivery_type" bson:"delivery_type"` // delivery or self-delivery
	CreatedAt         time.Time        `json:"created_at" bson:"created_at"`
	ValidFrom         time.Time        `json:"valid_from" bson:"valid_from"`
	ValidUntil        time.Time        `json:"valid_until" bson:"valid_until"`
	MinimumOrderPrice int              `json:"minimum_order_price" bson:"minimum_order_price"`
	PromoCodeCategory string           `json:"promo_code_category" bson:"promo_code_category"` // gift, sale
	SaleType          string           `json:"sale_type" bson:"sale_type"`                     // percentage, currency
	Sale              int              `json:"sale" bson:"sale"`                               // if percentage - 10 || if currency 1000
	Available         bool             `json:"available" bson:"available"`
	ForAllProduct     bool             `json:"for_all_product" bson:"for_all_product"`
	Product           []models.Product `json:"products" bson:"products"`
	IsDeleted         bool             `json:"is_deleted" bson:"is_deleted"`
}

type UserAndPromoCode struct {
	UserId        string   `json:"user_id" bson:"user_id"`
	PromoCode     string   `json:"promo_code" bson:"promo_code"`
	UsageTime     int      `json:"usage_time" bson:"usage_time"`
	RestaurantIds []string `json:"restaurant_ids" bson:"restaurant_ids,omitempty"`
}

type ValidateUserPromoCode struct {
	RestaurantID string           `json:"restaurant_id"`
	UserId       string           `json:"user_id"`
	TotalSum     int              `json:"total_sum"`
	PromoCode    string           `json:"promo_code"`
	CartProducts []models.Product `json:"cart_products"` // products
	DeliveryType string           `json:"delivery_type"` // delivery or self-delivery
}

type ValidateUserPromoCodeResponse struct {
	Exist      bool             `json:"exist"`
	Comment    string           `json:"comment"`
	TotalPrice float64          `json:"total_price"`
	SalePrice  int              `json:"sale_price"`
	Products   []models.Product `json:"products"`
}

type UpdatePromoCode struct {
	ID                string            `json:"id"`
	Code              *string           `json:"code"`
	Name              *string           `json:"name"`
	Description       *string           `json:"description"`
	Link              *string           `json:"link"`
	RestaurantIDs     *[]string         `json:"restaurant_ids"`
	UsageTime         *int              `json:"usage_time"`
	DeliveryType      *[]string         `json:"delivery_type"` // delivery or self-delivery
	CreatedAt         *time.Time        `json:"created_at"`
	ValidFrom         *time.Time        `json:"valid_from"`
	ValidUntil        *time.Time        `json:"valid_until"`
	MinimumOrderPrice *int              `json:"minimum_order_price"`
	PromoCodeCategory *string           `json:"promo_code_category"` // gift, sale
	SaleType          *string           `json:"sale_type"`           // percentage, currency
	Sale              *int              `json:"sale"`                //  if percentage - 10 || if currency 1000
	Available         *bool             `json:"available"`
	ForAllProduct     *bool             `json:"for_all_product"`
	Product           *[]models.Product `json:"products"`
	IsDeleted         *bool             `json:"is_deleted"`
}

type UserPromoCodeUsageTimeRequest struct {
	UserId         string `json:"user_id"`
	PromoCodeValue string `json:"promo_code"`
	RestaurantId   string `json:"restaurant_id"`
}
