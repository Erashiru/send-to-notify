package order

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"github.com/rs/zerolog/log"
)

const (
	retriesNumber   = 5
	retriesWaitTime = 1 * time.Second
)

const (
	acceptHeader      = "Accept"
	contentTypeHeader = "Content-Type"

	jsonType = "application/json"
)

type Order interface {
	OrderStatus(ctx context.Context, event models.EventInfo) error
	OrderTableStatus(ctx context.Context, event models.EventInfo) error
}

type order struct {
	cli *resty.Client
}

func NewOrderClient(integrationBaseURL string) Order {

	restyCli := resty.New().
		SetBaseURL(integrationBaseURL).
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
		})
	return &order{
		cli: restyCli,
	}
}

func (o order) OrderStatus(ctx context.Context, event models.EventInfo) error {

	path := "/v1/orders/iiko/updates"

	var errRsp errResponse

	resp, err := o.cli.R().
		SetContext(ctx).
		SetBody(event).
		SetError(&errRsp).
		Post(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("%s", errRsp.Detail)
	}

	log.Info().Msgf(fmt.Sprintf("%s - %s", resp.Status(), string(resp.Body())))

	return nil
}

func (o order) OrderTableStatus(ctx context.Context, event models.EventInfo) error {

	path := "/v1/orders/iiko/table/updates"

	var errRsp errResponse

	resp, err := o.cli.R().
		SetContext(ctx).
		SetBody(event).
		SetError(&errRsp).
		Post(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("%s", errRsp.Detail)
	}

	log.Info().Msgf(fmt.Sprintf("%s - %s", resp.Status(), string(resp.Body())))
	return nil
}
