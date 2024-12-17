package http

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
)

func (p *paloma) GetOrderStatus(ctx context.Context, orderID string) (models2.OrderResponse, error) {
	var (
		result      models2.OrderResponse
		errResponse models2.ErrorResponse
	)

	response, err := p.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetQueryParam("method", "status").
		SetQueryParam("order_id", orderID).
		SetResult(&result).
		SetError(&errResponse).
		Get("")
	if err != nil {
		return models2.OrderResponse{}, err
	}

	if response.IsError() || response.StatusCode() >= 400 {
		return models2.OrderResponse{}, errResponse
	}

	return result, nil
}
