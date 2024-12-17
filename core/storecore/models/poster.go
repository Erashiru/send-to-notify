package models

type StorePosterConfig struct {
	AccountNumber           int    `bson:"account_number" json:"account_number"`
	AccountNumberString     string `bson:"account_number_string" json:"account_number_string"`
	ApplicationNumber       int    `json:"application_number" bson:"application_number"`
	ApplicationSecret       string `json:"application_secret" bson:"application_secret"`
	ApplicationID           int    `json:"application_id" bson:"application_id"`
	Token                   string `json:"token" bson:"token"`
	Password                string `json:"password" bson:"password"`
	Username                string `json:"username" bson:"username"`
	CompanyName             string `json:"company_name" bson:"company_name"`
	SpotId                  string `json:"spot_id" bson:"spot_id"`
	CookingTime             int32  `json:"cooking_time" bson:"cooking_time"`
	StopListByBalance       bool   `json:"stop_list_by_balance" bson:"stop_list_by_balance"`
	IgnoreStopListFromAdmin bool   `json:"ignore_stop_list_from_admin" bson:"ignore_stop_list_from_admin"`
	IgnorePaymentType       bool   `json:"ignore_payment_type" bson:"ignore_payment_type"`
	PaymentType             int32  `json:"payment_type" bson:"payment_type"`
}
