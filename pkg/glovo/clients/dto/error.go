package dto

import (
	"fmt"
)

type CustomError struct {
	customError `json:"error"`
}

type customError struct {
	RequestID  string      `json:"requestId"`
	Code       string      `json:"code"`
	Msg        string      `json:"message"`
	Domain     string      `json:"domain"`
	StaticCode interface{} `json:"staticCode"`
}

func (c CustomError) Error() error {
	return fmt.Errorf("%s", c.Msg)
}
