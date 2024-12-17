package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) IsAlive(ctx context.Context, req models.IsAliveRequest) (models.IsAliveResponse, error) {
	path := "/api/1/terminal_groups/is_alive"

	var (
		response    models.IsAliveResponse
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(&req).
		Post(path)

	if err != nil {
		if resp != nil {
			utils.Beautify("is alive full response", *resp)
		}

		return models.IsAliveResponse{}, fmt.Errorf("is alive request failed, error: %v", errResponse)
	}

	if resp.IsError() {
		utils.Beautify("struct error", errResponse)
		return models.IsAliveResponse{}, fmt.Errorf("is alive response has error: %s", resp.Error())
	}

	return response, nil
}
