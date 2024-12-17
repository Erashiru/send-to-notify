package client

import (
	"context"
	"fmt"
	dto2 "github.com/kwaaka-team/orders-core/pkg/jowi/client/dto"
	"github.com/pkg/errors"
)

func (c *Client) CreateOrder(ctx context.Context, order dto2.RequestCreateOrder) (dto2.ResponseOrder, error) {
	path := "/v3/orders"

	var response dto2.ResponseOrder
	var errResponse dto2.ErrorResponse

	order.ApiKey = c.apiKey
	order.Sig = c.sig

	result, err := c.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(order).
		SetResult(&response).
		SetError(&errResponse).
		Post(path)
	if err != nil {
		return dto2.ResponseOrder{}, err
	}

	if result.IsError() {
		return dto2.ResponseOrder{}, errors.New(errResponse.Message)
	}

	// FIXME: JOWI returns 200 even if response has errors, so decided to add ErrorResponse to ResponseModel
	if response.Status != 1 {
		return dto2.ResponseOrder{}, errors.New(response.Message)
	}

	return response, nil
}

func (c *Client) GetOrder(ctx context.Context, restaurantID, orderID string) (dto2.ResponseOrder, error) {
	path := fmt.Sprintf("/v3/orders/%s", orderID)

	var response dto2.ResponseOrder
	var errResponse dto2.ErrorResponse

	// FIXME resty doesn't send request with body (Content-Type=application/json) if method is equal GET
	result, err := c.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetQueryParams(map[string]string{
			"api_key":       c.apiKey,
			"sig":           c.sig,
			"restaurant_id": restaurantID,
			"id":            orderID,
			// TODO: restaurant_id, order_id as param?
		}).
		SetResult(&response).
		SetError(&errResponse).
		Get(path)
	if err != nil {
		return dto2.ResponseOrder{}, err
	}

	if result.IsError() {
		return dto2.ResponseOrder{}, errors.New(errResponse.Message)
	}

	// FIXME: JOWI returns 200 even if response has errors, so decided to add ErrorResponse to ResponseModel
	if response.Status != 1 {
		return dto2.ResponseOrder{}, errors.New(response.Message)
	}

	return response, nil
}

func (c *Client) CancelOrder(ctx context.Context, orderID string, cancelOrder dto2.RequestCancelOrder) (dto2.ResponseOrder, error) {
	path := fmt.Sprintf("/v3/orders/%s/cancel", orderID)

	var response dto2.ResponseOrder
	var errResponse dto2.ErrorResponse

	cancelOrder.ApiKey = c.apiKey
	cancelOrder.Sig = c.sig

	// TODO: cancel order body is not valid in documentation (restaurant_id required too)
	result, err := c.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(cancelOrder).
		SetResult(&response).
		SetError(&errResponse).
		Post(path)
	if err != nil {
		return dto2.ResponseOrder{}, err
	}

	if result.IsError() {
		return dto2.ResponseOrder{}, errors.New(errResponse.Message)
	}

	// FIXME: JOWI returns 200 even if response has errors, so decided to add ErrorResponse to ResponseModel
	if response.Status != 1 {
		return dto2.ResponseOrder{}, errors.New(response.Message)
	}

	return response, nil
}
