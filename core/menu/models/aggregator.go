package models

const (
	GLOVO        AggregatorName = "glovo"
	WOLT         AggregatorName = "wolt"
	YANDEX       AggregatorName = "yandex"
	EMENU        AggregatorName = "emenu"
	CHOCOFOOD    AggregatorName = "chocofood"
	MOYSKLAD     AggregatorName = "moysklad"
	EXPRESS24    AggregatorName = "express24"
	EXPRESS24_V2 AggregatorName = "express24_v2"
	TALABAT      AggregatorName = "talabat"
	DELIVEROO    AggregatorName = "deliveroo"
	MOY_SKLAD    AggregatorName = "moysklad"
	SINGLECHOICE TypeSelection  = "SingleChoice"
	MULTICHOICE  TypeSelection  = "MultiChoice"
	HOMEDELIVERY DeliveryMethod = "homedelivery"
	TAKEAWAY     DeliveryMethod = "takeaway"
	STARTERAPP   AggregatorName = "starter_app"
	BASEQUANTITY int            = 5 //Minimum quantity of product in express24 for using stoplist in update product
)

type AggregatorName string
type TypeSelection string
type DeliveryMethod string

func (a AggregatorName) String() string {
	return string(a)
}

func (a TypeSelection) String() string {
	return string(a)
}

func (a DeliveryMethod) String() string {
	return string(a)
}

type AggregatorsConfig struct {
	BaseURL  string
	Username string
	Password string
	Token    string
	StoreID  string
	ApiKey   string
}
