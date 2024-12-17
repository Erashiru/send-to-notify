package mongodb

import (
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func errorSwitch(err error) error {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return errors.Wrap(drivers.ErrNotFound, "not found error")
	case mongo.IsDuplicateKeyError(err):
		return errors.Wrap(drivers.ErrAlreadyExist, "duplicate key error")
	default:
		return err
	}
}
