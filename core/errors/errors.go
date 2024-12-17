package errors

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrExpiredToken    = errors.New("token has expired")
	ErrTokenIsNotValid = errors.New("token is not valid")

	ErrInvalid             = errors.New("Invalid id")
	ErrInvalidConfigStruct = errors.New("Invalid configuration structure")
	ErrAlreadyExist        = errors.New("Already exist")
	ErrNotFound            = errors.New("Not found")
	ErrEmptySequenceID     = errors.New("Empty sequence name")

	ErrEmpty     = errors.New("empty")
	ErrEventType = errors.New("event type not supported")
	ErrTimeout   = errors.New("timeout")

	ErrUnsupportedMethod = errors.New("unsupported method")

	ErrNoWebhookSubscription = errors.New("there's no webhook subscription")

	ErrProductNotFound = errors.New("PRODUCT NOT FOUND IN POS MENU")
	ErrStoreNotFound   = errors.New("store not found")

	ErrValidateOrderStatusQueue = errors.New("error: validate order status queue")
)

var (
	ErrMissingReview = errors.New("missing review")
)

func ErrorSwitch(err error) error {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return ErrNotFound
	case mongo.IsDuplicateKeyError(err):
		return ErrAlreadyExist
	default:
		return err
	}
}

func CloseCur(cur *mongo.Cursor) {
	if err := cur.Close(context.Background()); err != nil {
		log.Err(err).Msg("closing cursor:")
	}
}
