package models

type StoreTalabatConfig struct {
	Username              string                     `bson:"username" json:"username"`
	Password              string                     `bson:"password" json:"password"`
	RestaurantID          string                     `bson:"restaurant_id" json:"restaurant_id"`
	BranchID              []string                   `bson:"store_id" json:"store_id"`
	RemoteBranchID        string                     `bson:"remote_branch_id" json:"remote_branch_id"`
	ChainID               string                     `bson:"chain_id" json:"chain_id"`
	VendorID              string                     `bson:"vendor_id" json:"vendor_id"`
	AdjustedPickupMinutes int                        `bson:"adjusted_pickup_minutes" json:"adjusted_pickup_minutes"`
	SendToPos             bool                       `bson:"send_to_pos" json:"send_to_pos"`
	IsMarketplace         bool                       `bson:"is_marketplace" json:"is_marketplace"`
	PaymentTypes          DeliveryServicePaymentType `bson:"payment_types" json:"payment_types"`
	IsNewMenu             bool                       `bson:"is_new_menu" json:"is_new_menu"`
	OrderCodePrefix       string                     `bson:"order_code_prefix"`
	IgnoreStatusUpdate    bool                       `bson:"ignore_status_update"`
}
