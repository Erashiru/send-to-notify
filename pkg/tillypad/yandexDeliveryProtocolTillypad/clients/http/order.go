package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/rs/zerolog/log"
)

func (c *clientImpl) GetOrder(ctx context.Context, orderId string) (models.Order, error) {
	path := c.pathPrefix + fmt.Sprintf("/order/%s", orderId)

	var (
		result  models.Order
		errResp models.ErrorResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetError(&errResp).
		SetResult(&result).
		Get(path)

	if err != nil {
		return models.Order{}, err
	}

	if resp.IsError() {
		return models.Order{}, fmt.Errorf("get order error: %v", resp.Error())
	}

	return result, nil
}

func (c *clientImpl) GetOrderStatus(ctx context.Context, orderId string) (models.OrderStatusResponse, error) {
	path := c.pathPrefix + fmt.Sprintf("/order/%s/status", orderId)

	var (
		result  models.OrderStatusResponse
		errResp models.ErrorResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetError(&errResp).
		SetResult(&result).
		Get(path)

	if err != nil {
		return models.OrderStatusResponse{}, err
	}

	if resp.IsError() {
		return models.OrderStatusResponse{}, fmt.Errorf("get order status error: %v", resp.Error())
	}

	return result, nil
}

func (c *clientImpl) CreateOrder(ctx context.Context, order models.Order) (models.CreationResult, error) {
	path := c.pathPrefix + "/order"

	log.Info().Msgf("TILLYPAD REQUEST FULL URL: %s", c.restyCli.BaseURL+path)

	var (
		result  models.CreationResult
		errResp models.ErrorResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetError(&errResp).
		SetBody(&order).
		SetResult(&result).
		Post(path)

	if err != nil {
		return models.CreationResult{}, err
	}

	log.Info().Msgf("TILLYPAD CREATE ORDER RESPONSE URL: %s", resp.Request.URL)

	if resp.IsError() {
		log.Info().Msgf("TILLYPAD GOT INTO RESPONSE ERROR: %s, CODE: %d", errResp.Description, errResp.Code)
		return models.CreationResult{}, fmt.Errorf("create order error: %v", resp.Error())
	}

	return result, nil
}
