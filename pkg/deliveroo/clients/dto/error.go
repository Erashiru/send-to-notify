package dto

import "fmt"

type CustomError struct {
	customError `json:"error"`
}

type customError struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

func (c CustomError) Error() error {
	return fmt.Errorf("%s", c.Msg)
}

type StoplistUpdateResponse struct {
	TransactionID string `json:"transaction_id"`
}
