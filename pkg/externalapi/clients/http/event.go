package http

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/chocofood/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/externalapi/clients"
	dto2 "github.com/kwaaka-team/orders-core/pkg/externalapi/clients/dto"
)

type Client struct {
	AuthToken   string
	restyClient *resty.Client
}

func NewClient(cfg *clients.Config) (clients.Client, error) {

	client := resty.New().
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
			authHeader:        cfg.AuthToken,
		}).
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime)

	cl := &Client{
		restyClient: client,
		AuthToken:   cfg.AuthToken,
	}

	return cl, nil
}

func (cli *Client) UpdateOrderWebhook(ctx context.Context, order dto2.Order, path string) error {
	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&order).
		Post(path)
	if err != nil {
		if resp != nil {
			utils.Beautify("response body", resp)
		}
		return err
	}

	utils.Beautify("base url", cli.restyClient.BaseURL)
	utils.Beautify("request body", resp.Request.Body)
	utils.Beautify("request headers", resp.Request.Header)

	if resp.IsError() {
		utils.Beautify("update order webhook error body", resp)
		return fmt.Errorf("update order webhook error body %s", resp)
	}

	return nil
}

func (cli *Client) UpdateProductStopList(ctx context.Context, product dto2.Product, path string) error {
	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&product).
		Post(path)
	if err != nil {
		if resp != nil {
			utils.Beautify("response body", resp)
		}
		return err
	}

	if resp.IsError() {
		utils.Beautify("update product stoplist webhook error body", resp)
		return fmt.Errorf("update product stoplist webhook error body %s", resp)
	}

	return nil
}

func (cli *Client) UpdateModifierStopList(ctx context.Context, modifier dto2.Modifier, path string) error {
	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&modifier).
		Post(path)
	if err != nil {
		if resp != nil {
			utils.Beautify("response body", resp)
		}
		return err
	}

	if resp.IsError() {
		utils.Beautify("update modifier stoplist webhook error body", resp)
		return fmt.Errorf("update modifier stoplist webhook error body %s", resp)
	}

	return nil
}
