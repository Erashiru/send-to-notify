package models

type UserStore struct {
	ID               string   `bson:"_id,omitempty" json:"id"`
	Username         string   `bson:"username" json:"username"`
	StoreId          string   `bson:"restaurant_id" json:"store_id"`
	StoreGroupId     string   `bson:"restaurant_group_id" json:"store_group_id"`
	FCMTokens        []string `bson:"fcm_tokens,omitempty" json:"fcm_tokens,omitempty"`
	SendNotification bool     `bson:"send_notification" json:"send_notification"`
}
