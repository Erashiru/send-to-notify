package models

type StoreGlovoConfig struct {
	StoreID                            []string                   `bson:"store_id" json:"store_id"`
	MenuUrl                            string                     `bson:"menu_url" json:"menu_url"`
	SendToPos                          bool                       `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketplace                      bool                       `bson:"is_marketplace" json:"is_marketplace"`
	PaymentTypes                       DeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	PurchaseTypes                      PurchaseTypes              `bson:"purchase_types" json:"purchase_types"`
	IsOpen                             bool                       `bson:"is_open" json:"is_open"`
	PartnersUsername                   string                     `bson:"partners_username" json:"partners_username"`
	PartnersPassword                   string                     `bson:"partners_password" json:"partners_password"`
	AdditionalPreparationTimeInMinutes int                        `bson:"additional_preparation_time_in_minutes"  json:"additional_preparation_time_in_minutes"`
	OrderCodePrefix                    string                     `bson:"order_code_prefix"`
	ScheduledOpenStore                 bool                       `bson:"scheduled_open_store"`
	AdjustedPickupMinutes              int                        `bson:"adjusted_pickup_minutes"`
	IgnoreStatusUpdate                 bool                       `bson:"ignore_status_update"`
	ScheduledBusyMode                  bool                       `bson:"scheduled_busy_mode"`
	ScheduledBusyModeTime              []ScheduledBusyModeTime    `bson:"scheduled_busy_mode_time"`
	StoreAutoOpen                      bool                       `bson:"store_auto_open"`
	AutoAcceptOn                       bool                       `bson:"auto_accept_on"`
	PostAutoAcceptOn                   bool                       `bson:"post_auto_accept_on,omitempty"`
}
