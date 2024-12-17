package yandex

import (
	"github.com/kwaaka-team/orders-core/pkg/yandex/clients"
	"github.com/kwaaka-team/orders-core/pkg/yandex/clients/http"
	"github.com/pkg/errors"
)

func NewClient(conf *clients.Config) (clients.Yandex, error) {
	switch conf.Protocol {
	case "http":
		cli, err := http.New(conf)
		if err != nil {
			return nil, err
		}
		return cli, nil
	}

	return nil, errors.New("invalid protocol")
}
