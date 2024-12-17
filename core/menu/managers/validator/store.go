package validator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models/custom"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/pkg/errors"
)

type Store interface {
	ValidateExternalID(ctx context.Context, query selector.Store) error
	ValidateDelivery(ctx context.Context, query selector.Store) error
	ValidateListStores(ctx context.Context, query selector.Menu) error
}

type storeImpl struct{}

var _ Store = (*storeImpl)(nil)

func NewStoreValidator() *storeImpl {
	return &storeImpl{}
}

func (s *storeImpl) ValidateExternalID(ctx context.Context, query selector.Store) error {
	var errs custom.Error

	if !query.HasExternalStoreID() && !query.HasID() {
		errs.Append(errors.New("store id is missing"))
	}

	if !query.HasDeliveryService() && !query.HasID() {
		errs.Append(errors.New("delivery service is missing"))
	}

	return errs.ErrorOrNil()
}

func (s *storeImpl) ValidateDelivery(ctx context.Context, query selector.Store) error {
	var errs custom.Error

	if !query.HasDeliveryService() {
		errs.Append(errors.New("delivery service is missing"))
	}

	return errs.ErrorOrNil()
}

func (s *storeImpl) ValidateListStores(ctx context.Context, query selector.Menu) error {
	var errs custom.Error

	if !query.HasProductExtID() {
		errs.Append(errors.New("product_id couldn't be empty"))
	}

	return errs.ErrorOrNil()
}
