package validator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/auth/models/selector"
	"github.com/kwaaka-team/orders-core/core/errors"
)

type User interface {
	ValidateUID(ctx context.Context, query selector.User) error
}

type userImpl struct{}

var _ User = (*userImpl)(nil)

func NewUserValidator() User {
	return &userImpl{}
}

func (s *userImpl) ValidateUID(ctx context.Context, query selector.User) error {
	var errs errors.Error

	if !query.HasUID() {
		errs.Append(errors.ErrInvalid)
	}

	return errs.ErrorOrNil()
}
