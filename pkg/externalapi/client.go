package externalapi

import (
	"github.com/kwaaka-team/orders-core/pkg/externalapi/clients"
	"github.com/kwaaka-team/orders-core/pkg/externalapi/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewWebhookClient(conf *clients.Config) (clients.Client, error) {
	switch conf.Protocol {
	case "http":
		cli, err := http.NewClient(conf)
		if err != nil {
			return nil, err
		}
		return cli, nil
	}

	return nil, ErrInvalidProtocol
}
