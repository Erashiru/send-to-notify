package dto

type UpdateItemStoplistRequest struct {
	Price     float64 `json:"price" example:"375.00"`
	Available bool    `json:"available" example:"false"`
}

type UpdateItemStoplistRespone struct {
	ID        string  `json:"id" example:"test_id"`
	Price     float64 `json:"price" example:"123"`
	Available bool    `json:"available" example:"false"`
}
