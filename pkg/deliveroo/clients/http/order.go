package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/dto"
)

func (c *Client) UpdateOrderStatus(ctx context.Context, req dto.UpdateOrderStatusRequest, orderID string) error {
	path := fmt.Sprintf("/api/v1/orders/%s", orderID)

	var (
		errResp dto.CustomError
	)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetBody(&req).
		SetError(errResp).
		Patch(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return errResp.Error()
	}
	return nil
}

func (c *Client) CreateSyncStatus(ctx context.Context, req dto.CreateSyncStatusRequest, orderID string) error {
	path := fmt.Sprintf("/api/v1/orders/%s/sync_status", orderID)

	var (
		errResp dto.CustomError
	)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetBody(&req).
		SetError(errResp).
		Post(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return errResp.Error()
	}
	return nil
}
