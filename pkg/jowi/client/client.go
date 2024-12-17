package client

import (
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/jowi"
)

type Client struct {
	cli    *resty.Client
	apiKey string
	sig    string
	quit   chan struct{}
}

func New(conf jowi.Config) (jowi.Jowi, error) {
	if conf.BaseURL == "" {
		conf.BaseURL = BaseURL
	}

	cli := resty.New().
		SetBaseURL(conf.BaseURL).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
		})

	c := &Client{
		cli:    cli,
		quit:   make(chan struct{}),
		apiKey: conf.ApiKey,
		sig:    generateSig(conf.ApiKey, conf.ApiSecret),
	}

	return c, nil
}
