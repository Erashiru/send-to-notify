package models

type StorePalomaConfig struct {
	PointID               string `bson:"point_id" json:"point_id"`
	ApiKey                string `bson:"api_key" json:"api_key"`
	StopListByBalance     bool   `bson:"stoplist_by_balance" json:"stoplist_by_balance"`
	StopListBalanceLimit  int    `bson:"stoplist_balance_limit" json:"stoplist_balance_limit"`
	AggregatorPriceTypeId string `bson:"aggregator_price_type_id" json:"aggregator_price_type_id"`
}
