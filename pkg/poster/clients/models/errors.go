package models

type ErrorResponse struct {
	Code    int    `json:"error"`
	Message string `json:"message"`
	Field   string `json:"field"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}
