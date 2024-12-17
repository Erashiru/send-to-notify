package external

import "github.com/pkg/errors"

var (
	ErrNoWebhookSubscription = errors.New("there's no webhook subscription")
	ErrNotImplemented        = errors.New("not implemented")
)
