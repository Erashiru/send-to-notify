package express24_v2

import (
	"github.com/kwaaka-team/orders-core/pkg/express24_v2/clients"
	"github.com/kwaaka-team/orders-core/pkg/express24_v2/clients/http"
	"github.com/pkg/errors"
)

var (
	ErrInvalidProtocol = errors.New("invalid protocol")
)

func NewExpress24Client(conf *clients.Config) (clients.Express24V2, error) {
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
