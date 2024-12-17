package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/talabat/models"
	"github.com/pkg/errors"
)

func (c *Client) MarkOrderPrepared(ctx context.Context, orderToken string) error {
	path := fmt.Sprintf("/v2/orders/%s/preparation-completed", orderToken)

	var errResponse models.AuthMWErrorResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		Post(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		if resp.IsError() {
			return errors.New(fmt.Sprintf("code: %s, message: %s", errResponse.Code, errResponse.Message))
		}
	}

	return nil
}

func (c *Client) AcceptOrder(ctx context.Context, req models.AcceptOrderRequest) error {
	path := fmt.Sprintf("/v2/order/status/%s", req.OrderToken)

	var errResponse models.AuthMWErrorResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(req).
		Post(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		if resp.IsError() {
			return errors.New(fmt.Sprintf("code: %s, message: %s", errResponse.Code, errResponse.Message))
		}
	}

	return nil
}

func (c *Client) RejectOrder(ctx context.Context, req models.RejectOrderRequest) error {
	path := fmt.Sprintf("/v2/order/status/%s", req.OrderToken)

	var errResponse models.AuthMWErrorResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(req).
		Post(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		if resp.IsError() {
			return errors.New(fmt.Sprintf("code: %s, message: %s", errResponse.Code, errResponse.Message))
		}
	}

	return nil
}

func (c *Client) OrderPickedUp(ctx context.Context, req models.OrderPickedUpRequest) error {
	path := fmt.Sprintf("/v2/order/status/%s", req.OrderToken)

	var errResponse models.AuthMWErrorResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(req).
		Post(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		if resp.IsError() {
			return errors.New(fmt.Sprintf("code: %s, message: %s", errResponse.Code, errResponse.Message))
		}
	}

	return nil
}
