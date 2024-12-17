package models

type StoreQRMenuConfig struct {
	StoreID               []string                   `bson:"store_id" json:"store_id"`
	URL                   string                     `bson:"url" json:"url"`
	IsIntegrated          bool                       `bson:"is_integrated" json:"is_integrated"`
	PaymentTypes          DeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	Hash                  string                     `bson:"hash" json:"hash"`
	CookingTime           int32                      `bson:"cooking_time" json:"cooking_time"`
	DeliveryTime          int                        `bson:"delivery_time" json:"delivery_time"`
	NoTable               bool                       `bson:"no_table" json:"no_table"`
	Theme                 string                     `bson:"theme" json:"theme"`
	IsMarketplace         bool                       `bson:"is_marketplace" json:"is_marketplace"`
	SendToPos             bool                       `bson:"send_to_pos"`
	IgnoreStatusUpdate    bool                       `bson:"ignore_status_update"`
	AdjustedPickupMinutes int                        `bson:"adjusted_pickup_minutes" json:"adjusted_pickup_minutes"`
	BusyMode              bool                       `bson:"busy_mode" json:"busy_mode"`
	IsMarketplaceForIIKO  bool                       `bson:"is_marketplace_for_iiko" json:"is_marketplace_for_iiko"` //если true, тогда  создается заказ на кассе и ресторан назначает курьера
}
