package dto

type ErrorResponse struct {
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}
