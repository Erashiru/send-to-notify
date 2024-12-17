package moysklad

import (
	"github.com/kwaaka-team/orders-core/pkg/moysklad/clients"
	"github.com/kwaaka-team/orders-core/pkg/moysklad/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewMoySkladClient(conf *clients.Config) (clients.MoySklad, error) {
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
