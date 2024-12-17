package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"

	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) GetStopList(ctx context.Context, req models.StopListRequest) (models.StopListResponse, error) {

	path := "/api/1/stop_lists"

	var (
		stopList models.StopListResponse
		errR     models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&stopList).
		SetError(&errR).
		Post(path)

	if err != nil {
		return models.StopListResponse{}, err
	}

	if resp.IsError() {
		return models.StopListResponse{}, fmt.Errorf("%w - %s", ErrStopList, err)
	}

	utils.Beautify("iiko stoplist response", stopList)

	return stopList, nil
}
