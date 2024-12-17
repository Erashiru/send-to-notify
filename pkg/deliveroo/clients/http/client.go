package http

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/deliveroo/clients"
	"github.com/pkg/errors"
	"time"
)

type Client struct {
	restyClient        *resty.Client
	BaseUrl            string
	Username, Password string
	quit               chan struct{}
}

func NewClient(cfg *clients.Config) (*Client, error) {

	if cfg.Username == "" || cfg.Password == "" {
		return nil, errors.New("username or password is not provided")
	}
	if cfg.BaseURL == "" {
		return nil, errors.New("base URL could not be empty")
	}

	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
		})

	cl := &Client{
		restyClient: client,
		Username:    cfg.Username,
		Password:    cfg.Password,
		BaseUrl:     cfg.BaseURL,
		quit:        make(chan struct{}),
	}

	if err := cl.Auth(context.Background()); err != nil {
		return nil, err
	}

	ticker := time.NewTicker(tokenTimeout)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := cl.Auth(context.Background()); err != nil {
					return
				}
			case <-cl.quit:
				ticker.Stop()
				return
			}
		}
	}()

	return cl, nil
}
