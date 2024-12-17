package http

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/poster/clients"
	"time"
)

type Client struct {
	cli   *resty.Client
	token string
	quit  chan struct{}
}

func New(conf *clients.Config) (clients.Poster, error) {

	cli := resty.New().
		SetBaseURL(conf.BaseURL).
		SetRetryCount(5).
		SetRetryWaitTime(3 * time.Second).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		})

	c := &Client{
		cli:   cli,
		quit:  make(chan struct{}),
		token: conf.Token,
	}

	c.auth(context.Background())

	return c, nil
}
