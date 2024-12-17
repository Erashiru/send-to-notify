package models

type ErrorResponse struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

func (e ErrorResponse) Error() string {
	return e.Info
}
