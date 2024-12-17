package models

type AuthInfo struct {
	Token       string   `bson:"token"`
	Restaurants []string `bson:"restaurant_ids"`
	PosType     string   `bson:"pos_type"`
	Active      bool     `bson:"is_active"`
}
