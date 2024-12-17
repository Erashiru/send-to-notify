package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
)

func (c *clientImpl) GetMenu(ctx context.Context, storeId string) (models.Menu, error) {
	path := c.pathPrefix + fmt.Sprintf("/menu/%s/composition", storeId)

	var (
		result  models.Menu
		errResp models.ErrorResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetError(&errResp).
		SetResult(&result).
		Get(path)

	if err != nil {
		return models.Menu{}, err
	}

	if resp.IsError() {
		return models.Menu{}, fmt.Errorf("get composition error: %v", resp.Error())
	}

	return result, nil
}
