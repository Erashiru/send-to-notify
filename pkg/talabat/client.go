package talabat

import (
	"github.com/kwaaka-team/orders-core/pkg/talabat/clients"
	"github.com/kwaaka-team/orders-core/pkg/talabat/clients/http"
	"github.com/pkg/errors"
)

func NewMenuClient(conf *clients.Config) (clients.TalabatMenu, error) {
	switch conf.Protocol {
	case "http":
		cli, err := http.NewMenu(conf)
		if err != nil {
			return nil, err
		}
		return cli, nil
	}

	return nil, errors.New("invalid protocol")
}

func NewMiddlewareClient(conf *clients.Config) (clients.TalabatMW, error) {
	switch conf.Protocol {
	case "http":
		cli, err := http.NewMW(conf)
		if err != nil {
			return nil, err
		}
		return cli, nil
	}

	return nil, errors.New("invalid protocol")
}
