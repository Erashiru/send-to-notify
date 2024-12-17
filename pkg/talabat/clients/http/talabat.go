package http

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/talabat/clients"
	"time"
)

type Client struct {
	cli      *resty.Client
	quit     chan struct{}
	username string
	password string
}

func NewMenu(conf *clients.Config) (clients.TalabatMenu, error) {
	cli := resty.New().
		SetBaseURL(conf.BaseURL).
		SetHeaders(map[string]string{
			"Content-Type": "application/json-patch+json",
			"Accept":       "application/json",
		})

	c := &Client{
		cli:      cli,
		quit:     make(chan struct{}),
		username: conf.Username,
		password: conf.Password,
	}

	if err := c.AuthMenu(context.Background()); err != nil {
		return nil, err
	}

	ticker := time.NewTicker(59 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := c.AuthMenu(context.Background()); err != nil {
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

func NewMW(conf *clients.Config) (clients.TalabatMW, error) {
	cli := resty.New().
		SetBaseURL(conf.BaseURL).
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Accept":       "application/json",
		})

	c := &Client{
		cli:      cli,
		quit:     make(chan struct{}),
		username: conf.Username,
		password: conf.Password,
	}

	if err := c.AuthMW(context.Background()); err != nil {
		return nil, err
	}

	ticker := time.NewTicker(59 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := c.AuthMW(context.Background()); err != nil {
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
