package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/yaros/models"
)

func (c *Client) CreateOrder(ctx context.Context, restID string, order models.OrderRequest) (models.OrderResponse, error) {
	path := fmt.Sprintf("/orders/%s", restID)

	var (
		response    models.OrderResponse
		errResponse models.ErrorResponse
	)
	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(order).
		Post(path)
	if err != nil {
		utils.Beautify("yaros create delivery order error", err)
		utils.Beautify("yaros create delivery order errorResponse body", errResponse)
		return models.OrderResponse{}, errResponse
	}
	if resp.IsError() {
		utils.Beautify("yaros create delivery order response", resp)
		utils.Beautify("yaros create delivery order errorResponse body", errResponse)
		return models.OrderResponse{}, fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return response, nil
}

func (c *Client) UpdateOrder(ctx context.Context, restID string, update models.OrderRequest) (models.OrderResponse, error) {
	path := fmt.Sprintf("/orders/%s", restID)

	var (
		response    models.OrderResponse
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(update).
		Patch(path)
	if err != nil {
		utils.Beautify("yaros update delivery order error", err)
		utils.Beautify("yaros update delivery order errorResponse body", errResponse)
		return models.OrderResponse{}, errResponse
	}

	if resp.IsError() {
		utils.Beautify("yaros update delivery order response", resp)
		utils.Beautify("yaros update delivery order errorResponse body", errResponse)
		return models.OrderResponse{}, fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return response, nil
}

func (c *Client) GetOrders(ctx context.Context, restID, infoSystem, department string) (models.OrderResponse, error) {

	path := fmt.Sprintf("/orders/%s?infosystem=%s&department=%s", restID, infoSystem, department)

	var (
		response    models.OrderResponse
		errResponse models.ErrorResponse
	)
	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		Get(path)

	if err != nil {
		utils.Beautify("yaros update delivery order error", err)
		utils.Beautify("yaros update delivery order errorResponse body", errResponse)
		return models.OrderResponse{}, errResponse
	}

	if resp.IsError() {
		utils.Beautify("yaros update delivery order response", resp)
		utils.Beautify("yaros update delivery order errorResponse body", errResponse)
		return models.OrderResponse{}, fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return response, nil
}
