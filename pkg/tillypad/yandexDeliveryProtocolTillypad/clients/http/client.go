package http

import (
	"context"
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg"
	"github.com/kwaaka-team/orders-core/pkg/tillypad/yandexDeliveryProtocolTillypad/clients"
)

type clientImpl struct {
	restyCli     *resty.Client
	clientId     string
	clientSecret string
	pathPrefix   string
}

func NewClient(cfg clients.Config) (*clientImpl, error) {
	restyClient := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeaders(map[string]string{
			pkg.ContentTypeHeader: pkg.JsonType,
		}).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	client := &clientImpl{
		restyCli:     restyClient,
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
		pathPrefix:   cfg.PathPrefix,
	}

	if err := client.auth(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *clientImpl) auth() error {
	var err error

	bearerToken, err := c.GetAccessToken(context.Background(), c.clientId, c.clientSecret)
	if err != nil {
		return err
	}

	c.restyCli.SetAuthToken(bearerToken)

	return nil
}
