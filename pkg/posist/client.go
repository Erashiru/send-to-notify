package posist

import (
	"errors"
	"github.com/kwaaka-team/orders-core/pkg/posist/clients"
	"github.com/kwaaka-team/orders-core/pkg/posist/clients/http"
)

var ErrInvalidProtocol = errors.New("invalid protocol")

func New(config *clients.Config) (clients.Posist, error) {
	switch config.Protocol {
	case "http":
		cli, err := http.New(config)
		if err != nil {
			return nil, err
		}

		return cli, nil
	}

	return nil, ErrInvalidProtocol
}
