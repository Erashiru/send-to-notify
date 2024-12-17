package burgerking

import (
	"github.com/kwaaka-team/orders-core/pkg/burgerking/clients"
	"github.com/kwaaka-team/orders-core/pkg/burgerking/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewBKClient(conf *clients.Config) (clients.BK, error) {
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
