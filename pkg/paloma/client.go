package paloma

import (
	"github.com/kwaaka-team/orders-core/pkg/paloma/clients"
	"github.com/kwaaka-team/orders-core/pkg/paloma/clients/http"
	"github.com/pkg/errors"
)

var ErrInvalidProtocol = errors.New("invalid protocol")

func New(config *clients.Config) (clients.Paloma, error) {

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
