package dto

import storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"

type StoreQRMenuConfig struct {
	URL                   string                                 `json:"url"`
	IsIntegrated          bool                                   `json:"is_integrated"`
	PaymentTypes          storeModels.DeliveryServicePaymentType `json:"payment_types"`
	Hash                  string                                 `json:"hash"`
	CookingTime           int                                    `json:"cooking_time"`
	DeliveryTime          int                                    `json:"delivery_time"`
	NoTable               bool                                   `json:"no_table"`
	Theme                 string                                 `json:"theme"`
	IsMarketplace         bool                                   `json:"is_marketplace"`
	AdjustedPickupMinutes int                                    `json:"adjusted_pickup_minutes"`
	BusyMode              bool                                   `json:"busy_mode"`
}
