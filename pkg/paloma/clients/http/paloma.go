package http

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/paloma/clients"
	"time"
)

type paloma struct {
	cli    *resty.Client
	apiKey string
	class  string
	quit   chan struct{}
}

func New(conf *clients.Config) (clients.Paloma, error) {

	cli := resty.New().
		SetBaseURL(conf.BaseURL).
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
		})

	c := &paloma{
		cli:    cli,
		quit:   make(chan struct{}),
		apiKey: conf.ApiKey,
		class:  conf.Class,
	}

	if err := c.Auth(context.Background()); err != nil {
		return nil, ErrAuth
	}

	ticker := time.NewTicker(tokenTimeout)
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
