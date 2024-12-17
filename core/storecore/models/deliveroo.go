package models

type StoreDeliverooConfig struct {
	BaseURL            string `bson:"base_url" json:"base_url"`
	Username           string `bson:"username" json:"username"`
	Password           string `bson:"password" json:"password"`
	StoreID            string `bson:"store_id" json:"store_id"`
	SendToPos          bool   `bson:"send_to_pos" json:"send_to_pos"`
	OrderCodePrefix    string `bson:"order_code_prefix" json:"order_code_prefix"`
	IgnoreStatusUpdate bool   `bson:"ignore_status_update" json:"ignore_status_update"`
}
