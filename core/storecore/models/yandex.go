package models

type StoreExternalConfig struct {
	StoreID                  []string                   `bson:"store_id" json:"store_id"`
	Type                     string                     `bson:"type" json:"type"`
	ServiceType              string                     `bson:"service_type" json:"service_type"`
	WebhookConfig            WebhookConfig              `bson:"webhook_config" json:"webhook_config"`
	MenuUrl                  string                     `bson:"menu_url" json:"menu_url"`
	SendToPos                bool                       `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketplace            bool                       `bson:"is_marketplace" json:"is_marketplace"`
	PaymentTypes             DeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	PurchaseTypes            PurchaseTypes              `bson:"purchase_types" json:"purchase_types"`
	ClientSecret             string                     `bson:"client_secret" json:"client_secret"`
	IgnoreStatusUpdate       bool                       `bson:"ignore_status_update" json:"ignore_status_update"`
	WebhookURL               string                     `bson:"webhook_url" json:"webhook_url"`
	AuthToken                string                     `bson:"auth_token" json:"auth_token"`
	AutoAcceptOn             bool                       `bson:"auto_accept_on" json:"auto_accept_on"`
	PostAutoAcceptOn         bool                       `bson:"post_auto_accept_on" json:"post_auto_accept_on"`
	WebhookProductStoplist   string                     `bson:"webhook_product_stoplist" json:"webhook_product_stoplist"`
	WebhookAttributeStoplist string                     `bson:"webhook_attribute_stoplist" json:"webhook_attribute_stoplist"`
	OrderCodePrefix          string                     `bson:"order_code_prefix" json:"order_code_prefix"`
	AdjustedPickupMinutes    int                        `bson:"adjusted_pickup_minutes"`
}

type WebhookConfig struct {
	OrderCreate   string `bson:"order_create" json:"order_create"`
	OrderCancel   string `bson:"order_cancel" json:"order_cancel"`
	RetryMaxCount int    `bson:"retry_max_count" json:"retry_max_count"`
}
