package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) GetOrderTypes(ctx context.Context, req models.OrderTypesRequest) (models.OrderTypesResponse, error) {

	path := "/api/1/deliveries/order_types"

	var (
		orderTypes models.OrderTypesResponse
		errR       models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&orderTypes).
		SetError(&errR).
		Post(path)

	if err != nil {
		return models.OrderTypesResponse{}, err
	}

	if resp.IsError() {
		return models.OrderTypesResponse{}, fmt.Errorf("order types %w - %s", ErrNotFound, err)
	}

	return orderTypes, nil
}
