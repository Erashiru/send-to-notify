package http

import "github.com/pkg/errors"

var (
	ErrResponse = errors.New("response error")
	ErrStopList = errors.New("no stop list for given organization")
	ErrNotFound = errors.New("not found for given organization")
)
