package models

type OrderResponse struct {
	Message string `json:"message"`
}

type ClosingResponse struct {
	Until string `json:"until"`
}
