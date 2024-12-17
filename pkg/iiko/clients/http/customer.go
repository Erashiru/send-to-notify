package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) GetCustomerInfo(ctx context.Context, req models.GetCustomerInfoRequest) (models.GetCustomerInfoResponse, error) {
	path := "/api/1/loyalty/iiko/customer/info"

	var (
		response    models.GetCustomerInfoResponse
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
			utils.Beautify("get customer info full response", *resp)
		}

		return models.GetCustomerInfoResponse{}, fmt.Errorf("get customer info request failed, error: %v", errResponse)
	}

	if resp.IsError() {
		utils.Beautify("struct error", errResponse)
		return models.GetCustomerInfoResponse{}, fmt.Errorf("get customer info response has error: %s", resp.Error())
	}

	return response, nil
}
