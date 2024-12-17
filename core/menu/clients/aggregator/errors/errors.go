package errors

import (
	"github.com/pkg/errors"
)

var (
	ErrAggregatorNotFound             = errors.New("aggregator not found")
	ErrNotImplemented                 = errors.New("method not implemented")
	ErrNoPermissionForPublishWoltMenu = errors.New("no permission for publish Wolt menu, menu has promo")
)
