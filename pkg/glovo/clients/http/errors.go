package http

import "github.com/pkg/errors"

var (
	ErrNotExist   = errors.New("not exist")
	ErrInvalid    = errors.New("invalid")
	ErrBadRequest = errors.New("bad request")
)
