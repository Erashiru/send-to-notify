package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/burgerking/clients"
	models2 "github.com/kwaaka-team/orders-core/pkg/burgerking/clients/models"
	"github.com/kwaaka-team/orders-core/pkg/burgerking/clients/models/utils"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	BaseURL string
	Env     string

	cli *resty.Client
}

func NewClient(conf *clients.Config) (clients.BK, error) {

	if conf.Address == "" {
		conf.Address = baseURL
	}

	cli := resty.New().
		SetBaseURL(conf.Address).
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
		})

	c := &Client{
		cli: cli,
	}

	return c, nil
}

func (c *Client) SendOrder(ctx context.Context, order models2.Order) (models2.OrderResponse, error) {
	path := "/orders/glovo"

	var (
		res         models2.OrderResponse
		errResponse models2.ErrorResponse
	)

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetBody(order).
		SetResult(&res).
		SetError(&errResponse).
		Post(path)

	if err != nil {
		if rsp != nil {
			utils.Beautify("response body", rsp)
		}

		return models2.OrderResponse{}, fmt.Errorf("err -> %w %s", errResponse, err.Error())
	}

	utils.Beautify("success response body", rsp)

	if rsp.IsError() {
		return models2.OrderResponse{Message: errResponse.Description}, fmt.Errorf("%w %v", errResponse, rsp.Error())
	}

	return res, nil
}

func (c *Client) CancelOrder(ctx context.Context, order models2.CancelOrderRequest) error {
	path := "/orders/glovo/cancellation"

	var (
		errResponse models2.ErrorResponse
	)

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetBody(order).
		SetError(&errResponse).
		Post(path)

	if err != nil {
		if rsp != nil {
			utils.Beautify("response body", rsp)
		}
		return fmt.Errorf("%w %s", errResponse, err.Error())
	}

	if rsp.IsError() {
		return fmt.Errorf("%w %s", errResponse, err.Error())
	}

	return nil
}

func (c *Client) SendMenu(ctx context.Context, req models2.SendMenuRequest) error {
	path := "/prod/burger_king"

	var (
		errResponse models2.ErrorResponse
	)

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetBody(req).
		SetError(&errResponse).
		Post(path)

	if err != nil {
		if rsp != nil {
			utils.Beautify("response body", rsp)
		}
		return fmt.Errorf("%w %s", errResponse, err.Error())
	}

	if rsp.IsError() {
		return fmt.Errorf("%w %s", errResponse, err.Error())
	}

	return nil
}

func (c *Client) SetSchedule(ctx context.Context, storeID string, req models2.SetScheduleRequest) error {
	path := fmt.Sprintf("/webhook/stores/%s/schedule", storeID)

	var (
		errResponse models2.ErrorResponse
	)

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetBody(req).
		SetError(&errResponse).
		Post(path)

	if err != nil {
		if rsp != nil {
			utils.Beautify("response body", rsp)
		}
		return fmt.Errorf("%w %s", errResponse, err.Error())
	}

	if rsp.IsError() {
		return fmt.Errorf("%w %s", errResponse, err.Error())
	}

	return nil
}

func (c *Client) GetClosing(ctx context.Context, storeID string) (models2.ClosingResponse, error) {
	path := fmt.Sprintf("/webhook/stores/%s/closing", storeID)

	var (
		resp        models2.ClosingResponse
		errResponse models2.ErrorResponse
	)

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetResult(&resp).
		SetError(&errResponse).
		Get(path)

	if err != nil {
		if rsp != nil {
			utils.Beautify("response body", rsp)
		}
		return models2.ClosingResponse{}, fmt.Errorf("%w %s", errResponse, err.Error())
	}

	if rsp.IsError() {
		return models2.ClosingResponse{}, fmt.Errorf("%w %s", errResponse, err.Error())
	}

	return resp, nil
}

func (c *Client) UpdateClosing(ctx context.Context, storeID string) (models2.ClosingResponse, error) {
	path := fmt.Sprintf("/webhook/stores/%s/closing", storeID)

	var (
		resp        models2.ClosingResponse
		errResponse models2.ErrorResponse
	)

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetResult(&resp).
		SetError(&errResponse).
		Put(path)

	if err != nil {
		if rsp != nil {
			utils.Beautify("response body", rsp)
		}
		return models2.ClosingResponse{}, fmt.Errorf("%w %s", errResponse, err.Error())
	}

	if rsp.IsError() {
		return models2.ClosingResponse{}, fmt.Errorf("%w %s", errResponse, err.Error())
	}

	return resp, nil
}

func (c *Client) DeleteClosing(ctx context.Context, storeID string) error {
	path := fmt.Sprintf("/webhook/stores/%s/closing", storeID)

	var (
		errResponse models2.ErrorResponse
	)

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		Delete(path)

	if err != nil {
		if rsp != nil {
			utils.Beautify("response body", rsp)
		}
		return errResponse
	}

	if rsp.IsError() {
		return errResponse
	}

	return nil
}

func (c *Client) StopLists(ctx context.Context, req models2.StopListRequest) error {
	path := "/webhook/stores/stop-lists"

	var (
		errResponse models2.ErrorResponse
	)

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetBody(req).
		SetError(&errResponse).
		Patch(path)

	if err != nil {
		if rsp != nil {
			utils.Beautify("response body", rsp)
		}
		return errResponse
	}

	if rsp.IsError() {
		return errResponse
	}

	return nil
}
