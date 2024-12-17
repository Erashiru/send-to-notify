package models

type RuleType string

const (
	OrderAddition RuleType = "order_addition"
)

type SupplementType string

const (
	ExceedOrderAmount = "exceed_order_amount"
)

type OrderRule struct {
	Type                  RuleType               `bson:"type"`
	SupplementType        SupplementType         `bson:"supplement_type"`
	RestaurantIDs         []string               `bson:"restaurant_ids"`
	OrderAmount           int                    `bson:"order_amount"`
	SupplementaryProducts []SupplementaryProduct `bson:"supplementary_products"`
}

type SupplementaryProduct struct {
	ProductId string  `bson:"product_id"`
	Name      string  `bson:"name"`
	Price     float64 `bson:"price"`
	Quantity  int     `bson:"quantity"`
}
