package models

type ErrorResponse struct {
	CorID       string `json:"correlationId"`
	Description string `json:"errorDescription"`
	Err         string `json:"error"`
}

func (er ErrorResponse) Error() string {
	return er.Err
}
