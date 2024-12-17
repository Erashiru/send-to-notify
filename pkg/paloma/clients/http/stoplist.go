package http

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	models2 "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
)

func (p *paloma) GetStopList(ctx context.Context, pointID string) (models2.StopList, error) {
	var (
		result      models2.StopList
		errResponse models2.ErrorResponse
	)

	response, err := p.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetQueryParam("method", "stoplist").
		SetQueryParam("point_id", pointID).
		SetResult(&result).
		SetError(&errResponse).
		Get("")
	if err != nil {
		return models2.StopList{}, err
	}

	if response.IsError() || response.StatusCode() >= 400 {
		return models2.StopList{}, errResponse
	}

	utils.Beautify("paloma stoplist response", result)

	return result, nil
}
