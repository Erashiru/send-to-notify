package http

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
)

func (p *paloma) GetPoints(ctx context.Context, authKey string) ([]models2.Point, error) {
	var (
		result      []models2.Point
		errResponse models2.ErrorResponse
	)

	response, err := p.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetQueryParam("method", "points").
		SetQueryParam("authkey", authKey).
		SetResult(&result).
		SetError(&errResponse).
		Get("")
	if err != nil {
		return nil, err
	}

	if response.IsError() || response.StatusCode() >= 400 {
		return nil, errResponse
	}

	return result, nil
}
