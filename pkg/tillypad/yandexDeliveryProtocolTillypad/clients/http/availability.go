package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
)

func (c *clientImpl) GetAvailability(ctx context.Context, storeId string) (models.StopListResponse, error) {
	path := c.pathPrefix + fmt.Sprintf("/menu/%s/availability", storeId)

	var (
		result  models.StopListResponse
		errResp models.ErrorResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetError(&errResp).
		SetResult(&result).
		Get(path)

	if err != nil {
		return models.StopListResponse{}, err
	}

	if resp.IsError() {
		return models.StopListResponse{}, fmt.Errorf("get availability error: %v", resp.Error())
	}

	utils.Beautify("tillypad stoplist response", result)

	return result, nil
}
