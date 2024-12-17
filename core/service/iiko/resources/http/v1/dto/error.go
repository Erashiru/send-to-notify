package dto

type ErrorResponse struct {
	Code        int    `json:"code,omitempty" example:"400"`
	Error       error  `json:"error,omitempty"`
	Msg         string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}
