package models

import "time"

type StoreWoltConfig struct {
	StoreID               []string                   `bson:"store_id"`
	MenuUsername          string                     `bson:"menu_username"`
	MenuPassword          string                     `bson:"menu_password"`
	ApiKey                string                     `bson:"api_key"`
	AdjustedPickupMinutes int                        `bson:"adjusted_pickup_minutes"`
	CookingTime           int                        `bson:"cooking_time"`
	BusyMode              bool                       `bson:"busy_mode"`
	MenuUrl               string                     `bson:"menu_url"`
	SendToPos             bool                       `bson:"send_to_pos"`
	IsMarketplace         bool                       `bson:"is_marketplace"`
	PaymentTypes          DeliveryServicePaymentType `bson:"payment_types"`
	PurchaseTypes         PurchaseTypes              `bson:"purchase_types"`
	IgnoreStatusUpdate    bool                       `bson:"ignore_status_update"`
	AutoAcceptOn          bool                       `bson:"auto_accept_on"`
	PostAutoAcceptOn      bool                       `bson:"post_auto_accept_on,omitempty"`
	IsOpen                bool                       `bson:"is_open"`
	IgnorePickupTime      bool                       `bson:"ignore_pickup_time"`
	OrderCodePrefix       string                     `bson:"order_code_prefix"`
	ScheduledBusyMode     bool                       `bson:"scheduled_busy_mode"`
	ScheduledBusyModeTime []ScheduledBusyModeTime    `bson:"scheduled_busy_mode_time"`
	StoreAutoOpen         bool                       `bson:"store_auto_open"`
}

type PurchaseTypes struct {
	Instant  []Status `bson:"instant"`
	Preorder []Status `bson:"preorder"`
	TakeAway []Status `bson:"takeaway"`
}

type Status struct {
	PosStatus string `bson:"pos_status" json:"pos_status"`
	Status    string `bson:"status" json:"status"`
}

type ScheduledBusyModeTime struct {
	From time.Time `bson:"from" json:"from"`
	To   time.Time `bson:"to" json:"to"`
}
