package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) GetCustomerTransactions(ctx context.Context, req models.GetTransactionInfoReq) (models.GetTransactionInfoResp, error) {
	path := "/api/1/loyalty/iiko/customer/transactions/by_revision"

	var (
		response    models.GetTransactionInfoResp
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(&req).
		Post(path)
	if err != nil {
		return models.GetTransactionInfoResp{}, err
	}

	if resp.IsError() {
		return models.GetTransactionInfoResp{}, fmt.Errorf("get customer transaction response has error: %s", resp.Error())
	}

	return response, nil
}
