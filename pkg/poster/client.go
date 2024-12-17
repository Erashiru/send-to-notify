package poster

import (
	"github.com/kwaaka-team/orders-core/pkg/poster/clients"
	"github.com/kwaaka-team/orders-core/pkg/poster/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewClient(conf *clients.Config) (clients.Poster, error) {
	switch conf.Protocol {
	case "http":
		cli, err := http.New(conf)
		if err != nil {
			return nil, err
		}
		return cli, nil
	}

	return nil, ErrInvalidProtocol
}
