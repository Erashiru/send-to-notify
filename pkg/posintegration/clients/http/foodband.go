package http

import (
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/posintegration/clients"
	"time"
)

type Client struct {
	cli            *resty.Client
	createOrderUrl string
	cancelOrderUrl string
	retryMaxCount  int
}

func New(conf *clients.Config) (clients.FOODBAND, error) {

	cli := resty.New().
		SetRetryCount(5).
		SetRetryWaitTime(1 * time.Second).
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Accept":        "application/json",
			"Authorization": conf.ApiToken,
		})

	c := &Client{
		cli:            cli,
		createOrderUrl: conf.CreateOrderUrl,
		cancelOrderUrl: conf.CancelOrderUrl,
		retryMaxCount:  conf.RetryMaxCount,
	}

	return c, nil
}
