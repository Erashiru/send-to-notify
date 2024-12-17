package glovo

import (
	"github.com/kwaaka-team/orders-core/pkg/glovo/clients"
	"github.com/kwaaka-team/orders-core/pkg/glovo/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewGlovoClient(conf *clients.Config) (clients.Glovo, error) {
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
