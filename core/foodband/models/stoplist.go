package models

type StopListReq struct {
	ID           string
	Price        float64
	IsAvailable  bool
	RestaurantID string
}
