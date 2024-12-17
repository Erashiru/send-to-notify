package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
)

func (c *clientImpl) GetPromos(ctx context.Context, storeId string) (models.Promo, error) {
	path := c.pathPrefix + fmt.Sprintf("/menu/%s/promos", storeId)

	var (
		result  models.Promo
		errResp models.ErrorResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetError(&errResp).
		SetResult(&result).
		Get(path)

	if err != nil {
		return models.Promo{}, err
	}

	if resp.IsError() {
		return models.Promo{}, fmt.Errorf("get promos error: %v", resp.Error())
	}

	return result, nil
}
