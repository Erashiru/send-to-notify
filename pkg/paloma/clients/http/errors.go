package http

import "github.com/pkg/errors"

var (
	ErrAuth   = errors.New("auth failed")
	ErrClass  = errors.New("auth class is empty")
	ErrApiKey = errors.New("auth api key is empty")
)
