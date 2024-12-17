package http

import (
	"context"
	"fmt"
	models2 "github.com/kwaaka-team/orders-core/pkg/poster/clients/models"
)

func (c *Client) GetSpots(ctx context.Context) (models2.GetSpotsResponse, error) {
	path := "/api/access.getSpots"

	var (
		response    models2.GetSpotsResponse
		errResponse models2.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		Get(path)
	if err != nil {
		return models2.GetSpotsResponse{}, fmt.Errorf("%v + %v", err, response)
	}

	if resp.IsError() {
		return models2.GetSpotsResponse{}, errResponse
	}

	if response.Message != "" {
		return models2.GetSpotsResponse{}, fmt.Errorf("get spots error: %s, status: %d", response.Message, response.Code)
	}

	return response, nil
}
