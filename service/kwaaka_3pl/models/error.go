package models

import "github.com/pkg/errors"

var (
	ErrDeliveryOrderIdIsEmpty = errors.New("3pl delivery order id is empty")
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
