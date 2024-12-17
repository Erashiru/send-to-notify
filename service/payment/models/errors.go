package models

import "github.com/pkg/errors"

var ErrUnsupportedMethod = errors.New("unsupported method for payment system")
var ErrUnsupportedPaymeMethod = errors.New("unsupported payme method")
