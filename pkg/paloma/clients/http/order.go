package http

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
)

func (p *paloma) CreateOrder(ctx context.Context, pointID string, req models2.Order) (models2.OrderResponse, error) {
	var (
		result      models2.OrderResponse
		errResponse models2.ErrorResponse
	)

	response, err := p.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetQueryParam("method", "order").
		SetQueryParam("point_id", pointID).
		SetBody(&req).
		SetResult(&result).
		SetError(&errResponse).
		Post("/")
	if err != nil {
		return models2.OrderResponse{}, err
	}

	if response.IsError() || response.StatusCode() >= 400 {
		return models2.OrderResponse{}, errResponse
	}

	return result, nil
}
