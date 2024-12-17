package models

type StoreChocofoodConfig struct {
	StoreID       []string                   `bson:"store_id" json:"store_id"`
	MenuUrl       string                     `bson:"menu_url" json:"menu_url"`
	SendToPos     bool                       `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketplace bool                       `bson:"is_marketplace" json:"is_marketplace"`
	PaymentTypes  DeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
}
