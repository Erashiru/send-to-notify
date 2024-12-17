package deliveroo

import (
	"github.com/kwaaka-team/orders-core/pkg/deliveroo/clients"
	"github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewDeliverooClient(conf *clients.Config) (*http.Client, error) {
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
