package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/yaros/models"
)

func (c *Client) GetStopList(ctx context.Context, restID string) (models.StopListResponse, error) {
	path := fmt.Sprintf("/stopList/%s", restID)

	var (
		stopList models.StopListResponse
		errR     models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetResult(&stopList).
		SetError(&errR).
		Get(path)

	if err != nil {
		return models.StopListResponse{}, err
	}

	if resp.IsError() {
		return models.StopListResponse{}, fmt.Errorf("%w - %s", ErrStopList, err)
	}

	return stopList, nil
}
