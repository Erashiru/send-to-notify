package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrAttributeGroupExtIDNotFound = errors.New("attribute group external ID not found")
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

func closeCur(cur *mongo.Cursor) {
	if err := cur.Close(context.Background()); err != nil {
		log.Err(err).Msg("closing cursor:")
	}
}
