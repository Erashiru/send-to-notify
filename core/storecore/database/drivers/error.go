package drivers

import "github.com/pkg/errors"

var (
	ErrInvalid             = errors.New("invalid id")
	ErrAlreadyExist        = errors.New("already exist")
	ErrNotFound            = errors.New("not found")
	ErrInvalidConfigStruct = errors.New("invalid configuration structure")
)
