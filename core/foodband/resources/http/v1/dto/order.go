package dto

type UpdateOrderStatusRequest struct {
	Status string `json:"status" example:"COOKING_STARTED"`
}
