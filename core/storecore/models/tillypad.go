package models

type TillyPadConfig struct {
	ClientId     string `json:"client_id" bson:"client_id"`
	ClientSecret string `json:"client_secret" bson:"client_secret"`
	PointId      string `json:"point_id" bson:"point_id"`
	PathPrefix   string `json:"path_prefix" bson:"path_prefix"`
}
