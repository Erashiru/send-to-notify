package models

type YTimesConfig struct {
	PointId   string `json:"point_id" bson:"point_id"`
	AuthToken string `json:"auth_token" bson:"auth_token"`
}
