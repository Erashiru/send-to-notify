package http

import (
	"context"
	"fmt"
	models2 "github.com/kwaaka-team/orders-core/pkg/poster/clients/models"
	"github.com/rs/zerolog/log"
	"net/url"
)

func (c *Client) CreateOrder(ctx context.Context, req models2.CreateOrderRequest) (models2.CreateOrderResponse, error) {
	path := "/api/incomingOrders.createIncomingOrder"

	var (
		response    models2.CreateOrderResponse
		errResponse models2.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(&req).
		Post(path)
	if err != nil {
		return models2.CreateOrderResponse{}, fmt.Errorf("%v + %v", err, response)
	}

	if resp.IsError() {
		return models2.CreateOrderResponse{}, errResponse
	}

	if response.Message != "" {
		return models2.CreateOrderResponse{}, fmt.Errorf("get products error: %s, status: %d", response.Message, response.Code)
	}

	log.Info().Msgf("poster CreateOrder response: %+v", response)

	return response, nil
}

func (c *Client) GetOrder(ctx context.Context, id string) (models2.CreateOrderResponse, error) {
	path := "/api/incomingOrders.getIncomingOrder"

	var (
		response    models2.CreateOrderResponse
		errResponse models2.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetQueryParam("incoming_order_id", id).
		Get(path)
	if err != nil {
		return models2.CreateOrderResponse{}, fmt.Errorf("%v + %v", err, response)
	}

	if resp.IsError() {
		return models2.CreateOrderResponse{}, errResponse
	}

	if response.Message != "" {
		return models2.CreateOrderResponse{}, fmt.Errorf("get products error: %s, status: %d", response.Message, response.Code)
	}

	return response, nil
}

func (c *Client) GetOrders(ctx context.Context, req models2.GetOrdersRequest) (models2.GetOrdersResponse, error) {
	path := "/api/incomingOrders.getIncomingOrders"
	var (
		response    models2.GetOrdersResponse
		errResponse models2.ErrorResponse
	)

	queryParams := url.Values{
		"status":    {req.Status},
		"date_from": {req.DateFrom},
		"date_to":   {req.DateTo},
	}
	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetQueryParamsFromValues(queryParams).
		Get(path)
	if err != nil {
		return models2.GetOrdersResponse{}, fmt.Errorf("%v + %v", err, response)
	}

	if resp.IsError() {
		return models2.GetOrdersResponse{}, errResponse
	}

	return response, nil
}
