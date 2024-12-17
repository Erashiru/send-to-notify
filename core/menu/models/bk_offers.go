package models

import coreModels "github.com/kwaaka-team/orders-core/core/models"

type BkOffers struct {
	ID         string          `bson:"_id,omitempty" json:"id"`
	OfferID    int             `bson:"offer_id" json:"offer_id"`
	Name       string          `bson:"name" json:"name"`
	FinalPrice int             `bson:"final_price" json:"final_price"`
	Discount   int             `bson:"discount" json:"discount"`
	GlovoPrice int             `bson:"glovo_price" json:"glovo_price"`
	IsActive   bool            `bson:"is_active" json:"is_active"`
	ProductID  string          `bson:"product_id" json:"product_id"`
	OnlyPrime  bool            `bson:"only_prime" json:"only_prime"`
	EndAt      string          `bson:"end_at" json:"end_at"`
	StartAt    string          `bson:"start_at" json:"start_at"`
	Days       []string        `bson:"days,omitempty" json:"days"`
	TimesOfDay string          `bson:"times_of_day" json:"times_of_day"`
	CreatedAt  coreModels.Time `bson:"created_at" json:"created_at"`
	UpdatedAt  coreModels.Time `bson:"updated_at" json:"updated_at"`
}
