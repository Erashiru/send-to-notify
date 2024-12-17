package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (c *Client) CreateTableOrder(ctx context.Context, createDeliveryBody models.CreateDeliveryRequest) (models.CreateDeliveryResponse, error) {
	path := "/api/1/order/create"

	var (
		response    models.CreateDeliveryResponse
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(createDeliveryBody).
		Post(path)
	if err != nil {
		return models.CreateDeliveryResponse{}, errResponse
	}
	log.Info().Msgf("req CreateTableOrder %+v", resp.Request)
	log.Info().Msgf("resp CreateTableOrder %+v", resp.RawResponse)

	if resp.IsError() {
		return models.CreateDeliveryResponse{}, fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return response, nil
}

func (c *Client) GetTables(ctx context.Context, req models.TableRequest) (models.TableResponse, error) {

	path := "/api/1/reserve/available_restaurant_sections"

	var (
		result      models.TableResponse
		errResponse models.ErrorResponse2
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).
		SetBody(req).
		Post(path)
	if err != nil {
		return result, errResponse
	}
	log.Info().Msgf("req GetTables %+v", resp.Request)
	log.Info().Msgf("resp GetTables %+v", resp.RawResponse)

	if resp.IsError() {
		return result, fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return result, nil
}

func (c *Client) GetOrdersByTables(ctx context.Context, req models.OrdersByTablesRequest) (models.OrdersByTablesResponse, error) {

	path := "/api/1/order/by_table"

	var (
		result      models.OrdersByTablesResponse
		errResponse models.ErrorResponse2
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).
		SetBody(req).
		Post(path)
	if err != nil {
		return result, errResponse
	}
	log.Info().Msgf("req GetOrdersByTables %+v", resp.Request)
	log.Info().Msgf("resp GetOrdersByTables %+v", resp.RawResponse)

	if resp.IsError() {
		return result, errors.New(errResponse.Description)
	}

	return result, nil
}

func (c *Client) GetOrdersByIDs(ctx context.Context, req models.GetOrdersByIDsRequest) (models.OrdersByTablesResponse, error) {

	path := "/api/1/order/by_id"

	var (
		result      models.OrdersByTablesResponse
		errResponse models.ErrorResponse2
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).
		SetBody(req).
		Post(path)
	if err != nil {
		return result, errResponse
	}
	log.Info().Msgf("req GetOrdersByIDs %+v", resp.Request)
	log.Info().Msgf("resp GetOrdersByIDs %+v", resp.RawResponse)

	if resp.IsError() {
		return result, fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return result, nil
}

func (c *Client) CloseTableOrder(ctx context.Context, req models.CloseTableOrderReq) (string, error) {
	path := "/api/1/order/close"

	var (
		result      models.ErrorResponse2
		errResponse models.ErrorResponse2
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).
		SetBody(req).
		Post(path)
	if err != nil {
		return "", errResponse
	}
	log.Info().Msgf("req CloseTableOrder %+v", resp.Request)
	log.Info().Msgf("resp CloseTableOrder %+v", resp.RawResponse)

	if resp.IsError() {
		return "", fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return result.CorID, nil
}

func (c *Client) AddOrdersPayment(ctx context.Context, req models.ChangePaymentReq) (string, error) {
	path := "/api/1/order/add_payments"

	var (
		result      models.ErrorResponse2
		errResponse models.ErrorResponse2
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).
		SetBody(req).
		Post(path)
	if err != nil {
		return "", errResponse
	}
	log.Info().Msgf("req AddOrdersPayment %+v", resp.Request)
	log.Info().Msgf("resp AddOrdersPayment %+v", resp.RawResponse)

	if resp.IsError() {
		return "", fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return result.CorID, nil
}

func (c *Client) GetCommandStatus(ctx context.Context, req models.GetCommandStatusReq) error {
	path := "/api/1/commands/status"

	var (
		result      models.GetCommandStatusResp
		errResponse models.ErrorResponse2
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).
		SetBody(req).
		Post(path)
	if err != nil {
		return errResponse
	}
	log.Info().Msgf("req GetCommandStatus %+v", resp.Request)
	log.Info().Msgf("resp GetCommandStatus %+v", resp.RawResponse)

	if resp.IsError() {
		return fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	if result.State != "Success" {
		return fmt.Errorf("%s", result.Exception.Message)
	}

	return nil
}
