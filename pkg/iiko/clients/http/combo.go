package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) GetCombos(ctx context.Context, req models.GetCombosRequest) (models.GetCombosResponse, error) {
	path := "/api/1/combo"

	var (
		response    models.GetCombosResponse
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
			utils.Beautify("get combos full response", *resp)
		}

		return models.GetCombosResponse{}, fmt.Errorf("get combos request failed, error: %v", errResponse)
	}

	if resp.IsError() {
		utils.Beautify("struct error", errResponse)
		return models.GetCombosResponse{}, fmt.Errorf("get combos response has error: %s", resp.Error())
	}

	return response, nil
}
