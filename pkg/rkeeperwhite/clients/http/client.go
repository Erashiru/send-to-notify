package http

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	tokenTimeout    = 90
	retriesNumber   = 3
	retriesWaitTime = 1 * time.Second
)

const (
	acceptHeader      = "Accept"
	contentTypeHeader = "Content-Type"
	jsonType          = "application/json"
	AggregatorAuth    = "AggregatorAuthentication"
)

type Client struct {
	restyClient *resty.Client
}

func NewClient(cfg *clients.Config) (clients.RKeeper, error) {

	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base url could not be empty")
	}

	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
			AggregatorAuth:    cfg.ApiKey,
		}).
		SetTimeout(40 * time.Second).
		SetRetryCount(1).
		SetRetryWaitTime(retriesWaitTime).
		SetProxy("http://max0Q8ga:mAxsWX6F2Y@185.120.79.119:50100")

	log.Info().Msgf("client: %+v", client)

	cl := Client{
		restyClient: client,
	}

	return cl, nil
}
