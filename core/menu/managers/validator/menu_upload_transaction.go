package validator

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/custom"
	"github.com/pkg/errors"
)

type MenuUploadTransaction interface {
	ValidateCreate(req models.MenuUploadTransaction) error
	ValidateUpdate(req models.UpdateMenuUploadTransaction) error
}

type menuUploadTransactionImpl struct{}

var _ MenuUploadTransaction = (*menuUploadTransactionImpl)(nil)

func NewMenuUploadTransactionValidator() *menuUploadTransactionImpl {
	return &menuUploadTransactionImpl{}
}

func (mut menuUploadTransactionImpl) ValidateCreate(req models.MenuUploadTransaction) error {
	var errs custom.Error

	if req.StoreID == "" {
		errs.Append(errors.New("restaurant_id could not be empty"))
	}

	if req.Status == "" {
		errs.Append(errors.New("status could not be empty"))
	}

	if req.Service == "" {
		errs.Append(errors.New("service could not be empty"))
	}

	// fixme details is required?
	if len(req.ExtTransactions) != 0 {
		errs.Append(mut.extTransactions(req.ExtTransactions))
	}

	return errs.ErrorOrNil()
}

func (mut menuUploadTransactionImpl) ValidateUpdate(req models.UpdateMenuUploadTransaction) error {
	var errs custom.Error

	if req.RestaurantID != nil && *req.RestaurantID == "" {
		errs.Append(errors.New("restaurant_id could not be empty"))
	}

	if req.Status != nil && *req.Status == "" {
		errs.Append(errors.New("status could not be empty"))
	}

	if req.Service != nil && *req.Service == "" {
		errs.Append(errors.New("service could not be empty"))
	}

	// fixme details is required?
	if len(req.ExtTransactions) != 0 {
		errs.Append(mut.extTransactions(req.ExtTransactions))
	}

	return errs.ErrorOrNil()
}

func (mut menuUploadTransactionImpl) extTransactions(req models.ExtTransactions) error {
	var errs custom.Error
	for i := range req {
		if req[i].ID == "" {
			errs.Append(errors.New("ext_transactions.id could not be empty"))
		}
		if req[i].Status == "" {
			errs.Append(errors.New("ext_transactions.status could not be empty"))
		}

	}
	return errs.ErrorOrNil()
}
