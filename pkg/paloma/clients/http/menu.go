package http

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
)

func (p *paloma) GetMenu(ctx context.Context, pointID string) (models2.Menu, error) {
	var (
		result      models2.Menu
		errResponse models2.ErrorResponse
	)

	response, err := p.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetQueryParam("method", "menu").
		SetQueryParam("point_id", pointID).
		SetResult(&result).
		SetError(&errResponse).
		Get("")
	if err != nil {
		return models2.Menu{}, err
	}

	if response.IsError() || response.StatusCode() >= 400 {
		return models2.Menu{}, errResponse
	}

	return result, nil
}
