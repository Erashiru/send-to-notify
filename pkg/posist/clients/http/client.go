package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg"
	"github.com/kwaaka-team/orders-core/pkg/posist/clients"
)

type posist struct {
	restyCli    *resty.Client
	customerKey string
}

func New(conf *clients.Config) (clients.Posist, error) {
	restyCli := resty.New().
		SetBaseURL(conf.BaseURL).
		SetHeaders(map[string]string{
			pkg.ContentTypeHeader: pkg.JsonType,
			pkg.AcceptHeader:      pkg.JsonType,
			pkg.AuthHeader:        conf.AuthBasic,
		})

	c := &posist{
		restyCli:    restyCli,
		customerKey: conf.CustomerKey,
	}

	return c, nil
}
