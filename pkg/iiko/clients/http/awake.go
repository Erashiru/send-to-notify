package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) AwakeTerminal(ctx context.Context, req models.IsAliveRequest) (models.AwakeResponse, error) {
	path := "/api/1/terminal_groups/awake"

	var (
		response    models.AwakeResponse
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

		return models.AwakeResponse{}, fmt.Errorf("awake request failed, error: %v", errResponse)
	}

	if resp.IsError() {
		utils.Beautify("struct error", errResponse)
		return models.AwakeResponse{}, fmt.Errorf("awake response has error: %s", resp.Error())
	}

	return response, nil
}
