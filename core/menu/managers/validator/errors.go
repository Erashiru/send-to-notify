package validator

import "github.com/pkg/errors"

var (
	ErrProductNotModifier = errors.New("products not updated")
	ErrEmptyStores        = errors.New("empty external stores")
)
