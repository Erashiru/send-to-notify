package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/dto"
)

func (c *Client) UploadMenu(ctx context.Context, menu dto.Menu, storeId string) error {
	var (
		response dto.UploadMenuResponse
		errResp  dto.CustomError
	)

	path := fmt.Sprintf("/api/v1/brands/%s/menus/%s", storeId, "site_id")

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetBody(&menu).
		SetResult(&response).
		SetError(errResp).
		Put(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return errResp.Error()
	}
	return nil
}
