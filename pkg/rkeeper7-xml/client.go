package rkeeper7_xml

import (
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewClient(conf *clients.Config) (clients.RKeeper7, error) {
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
