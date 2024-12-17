package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/yaros/models"
)

func (c *Client) GetItems(ctx context.Context, restID string) (models.GetItemsResponse, error) {
	path := fmt.Sprintf("/goods/%s", restID)

	var (
		resp        models.GetItemsResponse
		errResponse models.ErrorResponse
	)

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetResult(&resp).
		SetError(&errResponse).
		Get(path)
	if err != nil {
		return models.GetItemsResponse{}, err
	}
	if rsp.IsError() {
		return models.GetItemsResponse{}, fmt.Errorf("%s %w", rsp.Error(), ErrResponse)
	}
	return resp, nil
}

func (c *Client) GetCategories(ctx context.Context, restID string) (models.GetCategoriesResponse, error) {
	path := fmt.Sprintf("/categories/%s", restID)

	var (
		resp        models.GetCategoriesResponse
		errResponse models.ErrorResponse
	)

	rsp, err := c.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetResult(&resp).
		SetError(&errResponse).
		Get(path)
	if err != nil {
		return models.GetCategoriesResponse{}, err
	}

	if rsp.IsError() {
		return models.GetCategoriesResponse{}, fmt.Errorf("%s %w", rsp.Error(), ErrResponse)
	}
	return resp, nil
}
