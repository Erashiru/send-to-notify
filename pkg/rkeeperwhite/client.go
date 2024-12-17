package rkeeperwhite

import (
	"github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients"
	"github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewRKeeperClient(conf *clients.Config) (clients.RKeeper, error) {
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
