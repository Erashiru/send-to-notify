package dto

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"time"
)

type GetPromosSelector struct {
	StoreID         string   `json:"store_id"`
	ExternalStoreID string   `json:"external_store_id"`
	DeliveryService string   `json:"delivery_service"`
	ProductIDs      []string `json:"product_ids"`
	IsActive        bool     `json:"is_active"`
}

type PromoDiscount struct {
	Id                              string               `bson:"_id"`
	Name                            string               `bson:"name"`
	Type                            string               `bson:"type"`
	DeliveryPrice                   int                  `bson:"delivery_price"`
	HasAggregatorAndPartnerDiscount bool                 `bson:"has_aggregator_and_partner_discount"`
	StartDate                       time.Time            `bson:"start_date"`
	EndDate                         time.Time            `bson:"end_date"`
	StartTime                       string               `bson:"start_time"`
	EndTime                         string               `bson:"end_time"`
	RestaurantIds                   []string             `bson:"restaurant_ids"`
	ProductIds                      []string             `bson:"product_ids"`
	Percent                         int                  `bson:"percent"`
	Budget                          int                  `bson:"budget"`
	IsPosIntegrated                 bool                 `bson:"is_pos_integrated"`
	PosDiscount                     PosDiscountModel     `bson:"pos_discount"`
	PosType                         string               `bson:"pos_type"`
	DeliveryService                 string               `bson:"delivery_service"`
	IsActive                        bool                 `bson:"is_active"`
	CreatedAt                       time.Time            `bson:"created_at"`
	ProductGifts                    []ProductGift        `bson:"product_gifts"`
	CategoryPercentage              []CategoryPercentage `bson:"category_percentage"`
	CategoryFixed                   []CategoryFixed      `bson:"category_fixed"`
	OrderDiscount                   []OrderDiscount      `bson:"order_discounts"`
	PercentageForEachProduct        bool                 `bson:"percentage_for_each_product"`
	ProductsPercentage              []ProductPercentage  `bson:"products_percentage"`
}

type PosDiscountModel struct {
	DiscountId string `bson:"discount_id"`
	Type       string `bson:"type"`
}

type PosDiscount struct {
	CategoryPercentage []CategoryPercentage `json:"category_percentage"`
	CategoryFixed      []CategoryFixed      `json:"category_fixed"`
	OrderDiscounts     []OrderDiscount      `json:"order_discount"`
	Type               string               `json:"type"`
	Amount             int                  `json:"amount"`
	ProductIDs         []string             `json:"product_ids"`
}

type OrderDiscount struct {
	Type       string `json:"type"`
	DiscountID string `json:"discount_id"`
	Amount     int    `json:"amount"`
}

type CategoryPercentage struct {
	CategoryId string `json:"category_id"`
	Percent    int    `json:"percent"`
	DiscountId string `json:"discount_id"`
}

type CategoryFixed struct {
	CategoryId string `json:"category_id"`
	Sum        int    `json:"sum"`
	DiscountId string `json:"discount_id"`
}

type ProductGift struct {
	ProductId string `bson:"product_id"`
	PromoId   string `bson:"promo_id"`
}

type ProductPercentage struct {
	ProductID string `bson:"product_id"`
	Percent   int    `bson:"percent"`
}

func FromDiscounts(req models.Promo) PosDiscount {
	var response PosDiscount

	response.ProductIDs = req.ProductIds
	response.Type = req.Type
	response.Amount = req.Percent

	for _, fixed := range req.CategoryFixed {
		response.CategoryFixed = append(response.CategoryFixed, CategoryFixed{
			CategoryId: fixed.CategoryId,
			Sum:        fixed.Sum,
			DiscountId: fixed.DiscountId,
		})
	}

	for _, percentage := range req.CategoryPercentage {
		response.CategoryPercentage = append(response.CategoryPercentage, CategoryPercentage{
			CategoryId: percentage.CategoryId,
			Percent:    percentage.Percent,
			DiscountId: percentage.DiscountId,
		})
	}

	for _, discount := range req.OrderDiscount {
		response.OrderDiscounts = append(response.OrderDiscounts, OrderDiscount{
			Type:       discount.Type,
			DiscountID: discount.DiscountID,
			Amount:     discount.Amount,
		})
	}

	return response
}

func FromPromoDiscounts(req []models.Promo) []PromoDiscount {
	// Convert DB promo to dto promo

	response := make([]PromoDiscount, 0, len(req))

	for _, promo := range req {
		promoDiscount := PromoDiscount{
			Id:              promo.Id,
			Name:            promo.Name,
			Type:            promo.Type,
			DeliveryPrice:   promo.DeliveryPrice,
			StartDate:       promo.StartDate,
			EndDate:         promo.EndDate,
			StartTime:       promo.StartTime,
			EndTime:         promo.EndTime,
			RestaurantIds:   promo.RestaurantIds,
			ProductIds:      promo.ProductIds,
			PosType:         promo.PosType,
			Percent:         promo.Percent,
			Budget:          promo.Budget,
			IsPosIntegrated: promo.IsPosIntegrated,
			PosDiscount: PosDiscountModel{
				DiscountId: promo.PosDiscount.DiscountId,
				Type:       promo.PosDiscount.Type,
			},
			IsActive:                 promo.IsActive,
			CreatedAt:                promo.CreatedAt,
			PercentageForEachProduct: promo.PercentageForEachProduct,
		}

		for _, discount := range promo.OrderDiscount {
			promoDiscount.OrderDiscount = append(promoDiscount.OrderDiscount, OrderDiscount{
				Type:       discount.Type,
				DiscountID: discount.DiscountID,
				Amount:     discount.Amount,
			})
		}

		for _, fixed := range promo.CategoryFixed {
			promoDiscount.CategoryFixed = append(promoDiscount.CategoryFixed, CategoryFixed{
				CategoryId: fixed.CategoryId,
				Sum:        fixed.Sum,
				DiscountId: fixed.DiscountId,
			})
		}

		for _, percentage := range promo.CategoryPercentage {
			promoDiscount.CategoryPercentage = append(promoDiscount.CategoryPercentage, CategoryPercentage{
				CategoryId: percentage.CategoryId,
				Percent:    percentage.Percent,
				DiscountId: percentage.DiscountId,
			})
		}

		for _, productGift := range promo.ProductGifts {
			promoDiscount.ProductGifts = append(promoDiscount.ProductGifts, ProductGift{
				ProductId: productGift.ProductId,
				PromoId:   productGift.PromoId,
			})
		}

		for _, product := range promo.ProductsPercentage {
			promoDiscount.ProductsPercentage = append(promoDiscount.ProductsPercentage, ProductPercentage{
				ProductID: product.ProductID,
				Percent:   product.Percent,
			})
		}

		response = append(response, promoDiscount)
	}

	return response
}
