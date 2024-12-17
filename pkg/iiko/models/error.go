package models

type ErrorResponse struct {
	CorID       string `json:"correlation_id"`
	Description string `json:"error_description"`
	Err         string `json:"error"`
}

func (er ErrorResponse) Error() string {
	return er.Description
}

type ExternalErrorResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type ErrorResponse2 struct {
	CorID       string `json:"correlationId"`
	Description string `json:"errorDescription"`
	Err         string `json:"error"`
}

func (er ErrorResponse2) Error() string {
	return er.Description
}
