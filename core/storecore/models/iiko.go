package models

type PaymentType struct {
	PaymentTypeID            string `bson:"iiko_payment_type_id" json:"iiko_payment_type_id"`
	PaymentTypeKind          string `bson:"iiko_payment_type_kind" json:"iiko_payment_type_kind"`
	PromotionPaymentTypeID   string `bson:"promotion_payment_type_id"`
	OrderType                string `bson:"order_type" json:"order_type"`
	OrderTypeService         string `bson:"order_type_service" json:"order_type_service"`
	OrderTypeForVirtualStore string `bson:"order_type_for_virtual_store" json:"order_type_for_virtual_store"`
	IsProcessedExternally    *bool  `bson:"is_processed_externally,omitempty" json:"is_processed_externally"`
}

type StoreIikoConfig struct {
	OrganizationID                              string `bson:"organization_id" json:"organization_id"`
	TerminalID                                  string `bson:"terminal_id" json:"terminal_id"`
	Key                                         string `bson:"key" json:"key"`
	DiscountBalanceId                           string `bson:"discount_balance_id" json:"discount_balance_id"`
	SendKitchenComments                         bool   `bson:"send_kitchen_comments" json:"send_kitchen_comments"`
	StopListByBalance                           bool   `bson:"stoplist_by_balance" json:"stoplist_by_balance"`
	StopListBalanceLimit                        int    `bson:"stoplist_balance_limit" json:"stoplist_balance_limit"`
	IsExternalMenu                              bool   `bson:"is_external_menu" json:"is_external_menu"`
	ExternalMenuID                              string `bson:"external_menu_id" json:"external_menu_id"`
	PriceCategory                               string `bson:"price_category" json:"price_category"`
	HasCombo                                    bool   `bson:"has_combo" json:"has_combo"`
	IgnoreExternalMenuProductsWithZeroNullPrice bool   `bson:"ignore_external_menu_products_with_zero_null_price" json:"ignore_external_menu_products_with_zero_null_price"`
	CustomDomain                                string `bson:"custom_domain" json:"custom_domain"`
	RemovalTypeIdWithCharge                     string `bson:"removal_type_id_with_charge" json:"removal_type_id_with_charge"`
	RemovalTypeIdWithoutCharge                  string `bson:"removal_type_id_without_charge" json:"removal_type_id_without_charge"`
}
