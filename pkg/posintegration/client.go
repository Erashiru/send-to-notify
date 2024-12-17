package posintegration

import (
	"github.com/kwaaka-team/orders-core/pkg/posintegration/clients"
	"github.com/kwaaka-team/orders-core/pkg/posintegration/clients/http"
)

func NewClient(conf *clients.Config) (clients.FOODBAND, error) {
	cli, err := http.New(conf)
	if err != nil {
		return nil, err
	}

	return cli, nil
}
