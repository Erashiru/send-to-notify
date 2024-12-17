package order_rules

import (
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func errorSwitch(err error) error {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return drivers.ErrNotFound
	case mongo.IsDuplicateKeyError(err):
		return drivers.ErrAlreadyExist
	default:
		return err
	}
}
