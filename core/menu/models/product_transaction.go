package models

import coreModels "github.com/kwaaka-team/orders-core/core/models"

type ProductTransaction struct {
	TxID         string          `bson:"transaction_id"`
	MenuID       string          `bson:"menu_id"`
	RestaurantID string          `bson:"restaurant_id"`
	Delivery     string          `bson:"delivery"`
	CreatedAt    coreModels.Time `bson:"created_at"`
	Products     Products        `bson:"products"`
	Attributes   Attributes      `bson:"attributes"`
}
