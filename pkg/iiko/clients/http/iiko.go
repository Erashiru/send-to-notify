package http

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/iiko/clients"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	cli    *resty.Client
	apiKey string
	quit   chan struct{}
}

func New(conf *clients.Config) (clients.IIKO, error) {

	cli := resty.New().
		SetBaseURL(conf.BaseURL).
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		})

	c := &Client{
		cli:    cli,
		quit:   make(chan struct{}),
		apiKey: conf.ApiLogin,
	}

	if err := c.Auth(context.Background()); err != nil {
		return nil, ErrAuth
	}

	ticker := time.NewTicker(59 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := c.Auth(context.Background()); err != nil {
					return
				}
			case <-c.quit:
				ticker.Stop()
				return
			}
		}
	}()

	return c, nil
}
