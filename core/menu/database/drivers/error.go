package drivers

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalid             = errors.New("invalid id")
	ErrInvalidConfigStruct = errors.New("invalid configuration structure")
	ErrAlreadyExist        = errors.New("already exist")
	ErrNotFound            = errors.New("not found")
	ErrEmptySequenceID     = errors.New("empty sequence name")
)
