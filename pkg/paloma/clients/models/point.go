package models

type Point struct {
	PointId int    `json:"point_id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}
