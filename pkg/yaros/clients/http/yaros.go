package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/yaros/clients"
)

type Client struct {
	cli *resty.Client
}

func New(conf *clients.Config) (*Client, error) {

	cli := resty.New().
		SetBaseURL(conf.BaseURL).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
		}).
		SetBasicAuth(conf.Username, conf.Password)

	return &Client{
		cli: cli,
	}, nil
}
