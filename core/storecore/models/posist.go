package models

type StorePosistConfig struct {
	CustomerKey string `json:"customer_key" bson:"customer_key"`
	TabId       string `json:"tab_id" bson:"tab_id"`
	AuthBasic   string `json:"auth_basic" bson:"auth_basic"`
	Username    string `json:"username" bson:"username"`
	Password    string `json:"password" bson:"password"`
}
