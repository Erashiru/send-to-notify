package http

import "github.com/pkg/errors"

var (
	ErrAuth     = errors.New("could not authenticate apiLogin")
	ErrResponse = errors.New("response error")
	ErrStopList = errors.New("no stop list for given organization")
	ErrNotFound = errors.New("not found for given organization")
)
