package models

type StoreKwaakaAdminConfig struct {
	IsIntegrated bool `bson:"is_integrated" json:"is_integrated"`
	IsActive     bool `bson:"is_active" json:"is_active"`
	//PaymentTypes DeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	CookingTime   int32    `bson:"cooking_time" json:"cooking_time"`
	StoreID       []string `bson:"store_id" json:"store_id"`
	SendToPos     bool     `bson:"send_to_pos"`
	IsMarketPlace bool     `bson:"is_market_place" json:"is_market_place"`
}
