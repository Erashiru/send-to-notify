package http

import (
	"context"
	"fmt"
	models2 "github.com/kwaaka-team/orders-core/pkg/poster/clients/models"
)

func (c *Client) GetStopList(ctx context.Context) (models2.GetStopListResponse, error) {
	path := "/api/storage.getStorageLeftovers"

	var (
		response    models2.GetStopListResponse
		errResponse models2.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		Get(path)
	if err != nil {
		return response, fmt.Errorf("%v + %v", err, response)
	}

	if resp.IsError() {
		return response, errResponse
	}
	return response, nil
}
