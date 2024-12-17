package http

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/yandex/clients"
	"time"
)

type Client struct {
	cli          *resty.Client
	quit         chan struct{}
	clientID     string
	clientSecret string
}

func New(conf *clients.Config) (clients.Yandex, error) {
	cli := resty.New().
		SetBaseURL(conf.BaseURL)

	c := &Client{
		cli:          cli,
		quit:         make(chan struct{}),
		clientID:     conf.ClientID,
		clientSecret: conf.ClientSecret,
	}

	if err := c.Auth(context.Background()); err != nil {
		return nil, err
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
