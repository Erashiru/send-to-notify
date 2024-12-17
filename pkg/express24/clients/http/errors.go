package http

import "github.com/pkg/errors"

var (
	ErrAuth = errors.New("could not authenticate apiLogin")
)
