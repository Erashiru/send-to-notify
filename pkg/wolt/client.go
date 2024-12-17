package wolt

import (
	"github.com/kwaaka-team/orders-core/pkg/wolt/clients"
	"github.com/kwaaka-team/orders-core/pkg/wolt/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewWoltClient(conf *clients.Config) (clients.Wolt, error) {
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
