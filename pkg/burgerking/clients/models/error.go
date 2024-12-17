package models

type ErrorResponse struct {
	Code        int    `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Reason      string `json:"reason,omitempty"`
	ErrorPlane  string `json:"error"`
}

func (er ErrorResponse) Error() string {
	return er.ErrorPlane
}
