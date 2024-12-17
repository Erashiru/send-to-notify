package express24

import (
	"github.com/kwaaka-team/orders-core/pkg/express24/clients"
	"github.com/kwaaka-team/orders-core/pkg/express24/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewExpress24Client(conf *clients.Config) (clients.Express24, error) {
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
