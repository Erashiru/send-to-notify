package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BKOffer struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Name       string             `bson:"name" json:"name"`
	FinalPrice int                `bson:"final_price" json:"final_price"`
	Discount   int                `bson:"discount" json:"discount"`
	OfferID    int                `bson:"offer_id" json:"offer_id"`
	GlovoPrice int                `bson:"glovo_price" json:"glovo_price"`
	ProductID  string             `bson:"product_id" json:"product_id"`
	StartAt    string             `bson:"start_at" json:"start_at"`
	EndAt      string             `bson:"end_at" json:"end_at"`
	OnlyPrime  bool               `bson:"only_prime" json:"only_prime"`
	IsActive   bool               `bson:"is_active" json:"is_active"`
}
