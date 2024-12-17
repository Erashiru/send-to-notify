package models

type ErrorResponse struct {
	ErrorPlane  string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
}

func (er ErrorResponse) Error() string {
	return er.ErrorPlane
}

type ErrorClickUpTask struct {
	Err   string `json:"err"`
	Ecode string `json:"ecode"`
}

func (err ErrorClickUpTask) Error() string {
	return err.Err
}
