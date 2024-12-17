package models

type StoreRKeeperConfig struct {
	ObjectId                      int    `bson:"object_id" json:"object_id"`
	ApiKey                        string `bson:"api_key" json:"api_key"`
	PriceTypeID                   int    `bson:"price_type_id,omitempty" json:"price_type_id,omitempty"`
	IgnoreUpsertProductWithPrice0 bool   `bson:"ignore_upsert_product_with_price_0" json:"ignore_upsert_product_with_price_0"`
}

type StoreRKeeper7XMLConfig struct {
	Domain                 string `bson:"domain" json:"domain"`
	Username               string `bson:"username" json:"username"`
	Password               string `bson:"password" json:"password"`
	UCSUsername            string `bson:"ucs_username" json:"ucs_username"`
	UCSPassword            string `bson:"ucs_password" json:"ucs_password"`
	Token                  string `bson:"token" json:"token"`
	ObjectID               string `bson:"object_id" json:"object_id"`
	Anchor                 string `bson:"anchor" json:"anchor"`
	LicenseInstanceGUID    string `bson:"license_instance_guid" json:"license_instance_guid"`
	DefaultTable           string `bson:"default_table" json:"default_table"`
	StationID              string `bson:"station_id" json:"station_id"`
	StationCode            string `bson:"station_code"`
	SeqNumber              int    `bson:"seq_number" json:"seq_number"`
	ChildItems             int    `bson:"child_items"`
	ClassificatorItemIdent int    `bson:"classificator_item_ident"`
	ClassificatorPropMask  string `bson:"classificator_prop_mask"`
	MenuItemsPropMask      string `bson:"menu_items_prop_mask"`
	PropFilter             string `bson:"prop_filter"`
	Cashier                string `bson:"cashier"`
	OrderType              string `bson:"order_type"`
	OrderTypeCode          string `bson:"order_type_code"`
	PrepayReasonId         string `bson:"prepay_reason_id"`
	TradeGroupId           string `bson:"trade_group_id"`
	PriceTypeId            string `bson:"price_type_id"`
	IsLifeTimeLicence      bool   `bson:"is_life_time_licence" json:"is_life_time_licence"`
}
