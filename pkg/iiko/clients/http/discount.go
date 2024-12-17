package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) GetDiscounts(ctx context.Context, organizationID string) (models.StoreDiscountsResponse, error) {

	type TerminalGroupsRequest struct {
		Organizations []string `json:"organizationIds"`
	}

	var (
		path = "/api/1/discounts"
		req  = TerminalGroupsRequest{
			Organizations: []string{organizationID},
		}
		response    models.StoreDiscountsResponse
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(&req).
		Post(path)
	if err != nil {
		return models.StoreDiscountsResponse{}, fmt.Errorf("http/discount.go - fn GetDiscounts - fn c.cli.R(): counln't send request while getting discounts: %w", err)
	}
	if resp.IsError() {
		return models.StoreDiscountsResponse{}, fmt.Errorf("http/discount.go - fn GetDiscounts - if resp.IsError(): counln't get response while getting discounts: %w", err)
	}

	return response, nil
}
