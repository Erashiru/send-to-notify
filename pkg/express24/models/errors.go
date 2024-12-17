package models

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Code   int    `json:"code"`
}

func (er ErrorResponse) ToString() string {
	return er.Error
}
