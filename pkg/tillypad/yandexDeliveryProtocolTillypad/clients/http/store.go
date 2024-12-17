package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
)

func (c *clientImpl) GetStores(ctx context.Context) (models.GetStoreResponse, error) {
	path := c.pathPrefix + "/restaurants"

	var (
		result  models.GetStoreResponse
		errResp models.ErrorResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetError(&errResp).
		SetResult(&result).
		Get(path)

	if err != nil {
		return models.GetStoreResponse{}, err
	}

	if resp.IsError() {
		return models.GetStoreResponse{}, fmt.Errorf("get pos stores error: %v", resp.Error())
	}

	return result, nil
}
