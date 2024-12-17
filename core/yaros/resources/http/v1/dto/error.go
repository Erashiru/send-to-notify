package dto

import "github.com/pkg/errors"

var (
	ErrTokenIsNotValid = errors.New("token is not valid")
	ErrTokenExpired    = errors.New("access token has been expired. you should request a new one")
)

type ErrorResponse struct {
	Code        int    `json:"code,omitempty" example:"400"`
	Error       error  `json:"error,omitempty"`
	Msg         string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`
}
