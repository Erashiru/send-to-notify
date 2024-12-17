package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/foodband/models/utils"
	"github.com/kwaaka-team/orders-core/domain/foodband"
	"time"
)

func (c *Client) CreateOrder(ctx context.Context, createOrderBody foodband.CreateOrderRequest) (int, error) {
	var err error
	i := 0
	for i <= c.retryMaxCount {
		err = c.SendOrder(ctx, createOrderBody)
		if err == nil {
			return i, nil
		}
		time.Sleep(3 * time.Second)
		i++
	}
	return i, err
}

func (c *Client) SendOrder(ctx context.Context, createOrderBody foodband.CreateOrderRequest) error {
	var (
		errResponse any
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(createOrderBody).
		Post(c.createOrderUrl)
	if err != nil {
		utils.Beautify("foodband create order error", err)
		utils.Beautify("foodband create order resp", resp)
		return err
	}

	if resp.IsError() {
		utils.Beautify("foodband create order response", resp)
		utils.Beautify("foodband create order errorResponse body", errResponse.(string))
		return fmt.Errorf(errResponse.(string))
	}
	if resp.StatusCode() > 299 {
		return fmt.Errorf("foodband create order error: %v", resp.StatusCode())
	}

	return nil
}

func (c *Client) CancelOrder(ctx context.Context, cancelOrderBody foodband.CancelOrderRequest) error {
	var (
		errResponse any
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(cancelOrderBody).
		Post(c.cancelOrderUrl)
	if err != nil {
		utils.Beautify("foodband cancel order error", err)
		return err
	}

	if resp.IsError() {
		utils.Beautify("foodband cancel order response", resp)
		utils.Beautify("foodband cancel order errorResponse body", errResponse)
		return fmt.Errorf(errResponse.(string))
	}
	if resp.StatusCode() > 299 {
		return fmt.Errorf("foodband cancel order error: %v", resp.StatusCode())
	}

	return nil
}
