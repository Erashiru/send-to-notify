package models

type StoreExpress24Config struct {
	StoreID               []string                   `bson:"store_id" json:"store_id"`
	MenuUrl               string                     `bson:"menu_url" json:"menu_url"`
	Username              string                     `bson:"menu_username" json:"username"`
	Password              string                     `bson:"menu_password" json:"password"`
	AdjustedPickupMinutes int                        `bson:"adjusted_pickup_minutes" json:"adjusted_pickup_minutes"`
	SendToPos             bool                       `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketplace         bool                       `bson:"is_marketplace" json:"is_marketplace"`
	PaymentTypes          DeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	IsOpen                bool                       `bson:"is_open" json:"is_open"`
	OrderCodePrefix       string                     `bson:"order_code_prefix" json:"order_code_prefix"`
	IgnoreStatusUpdate    bool                       `bson:"ignore_status_update" json:"ignore_status_update"`
	Token                 string                     `bson:"token" json:"token"`
	Vat                   int                        `bson:"vat" json:"vat"`
}
