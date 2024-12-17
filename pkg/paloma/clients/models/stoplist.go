package models

type StopList struct {
	PointId       int            `json:"point_id"`
	StopListItems []StopListItem `json:"items"`
}

type StopListItem struct {
	ObjectId int `json:"object_id"`
}
