package validator

import (
	"context"
	menuValidator "github.com/kwaaka-team/orders-core/core/menu/managers/validator/menu"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/custom"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/pkg/errors"
)

type Menu interface {
	ValidateID(ctx context.Context, query selector.Menu) error
	ValidateMenu(ctx context.Context, menu models.Menu) error
}

type MenuImpl struct{}

var _ Menu = (*MenuImpl)(nil)

func NewMenuValidator() *MenuImpl {
	return &MenuImpl{}
}

var (
	ErrAddTransaction = errors.New("could not add transactions to change")
)

func (s *MenuImpl) ValidateID(ctx context.Context, query selector.Menu) error {
	var errs custom.Error

	if !query.HasMenuID() {
		errs.Append(errors.New("menu id is missing"))
	}

	return errs.ErrorOrNil()
}

func (s *MenuImpl) ValidateMenu(ctx context.Context, menu models.Menu) error {
	var errs custom.Error

	factory, err := menuValidator.NewValidatorMenu(menu.Delivery)
	if err != nil {
		return err
	}
	errs.Append(factory.Validate(ctx, menu))

	return errs.ErrorOrNil()
}
