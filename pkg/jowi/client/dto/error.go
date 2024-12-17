package dto

type ErrorResponse struct {
	Status  int    `json:"status"`
	Error   int    `json:"error"`
	Message string `json:"message"`
}
