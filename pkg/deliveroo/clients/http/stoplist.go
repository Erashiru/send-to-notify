package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/dto"
)

func (c *Client) GetUnavailabilities(ctx context.Context, req dto.GetUnavailabilitiesRequest) (dto.GetUnavailabilitiesResponse, error) {
	path := fmt.Sprintf("api/v1/brands/%s/menus/%s/item_unavailabilities/%s", req.BrandID, req.MenuID, req.SiteID)

	var (
		response dto.GetUnavailabilitiesResponse
		errResp  dto.CustomError
	)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetError(errResp).
		SetResult(&response).
		Get(path)

	if err != nil {
		return dto.GetUnavailabilitiesResponse{}, err
	}

	if resp.IsError() {
		return dto.GetUnavailabilitiesResponse{}, errResp.Error()
	}
	return response, nil
}

func (c *Client) UpdateUnavailabileItems(ctx context.Context, req dto.UpdateUnavailabilitesRequest) error {
	path := fmt.Sprintf("api/v1/brands/%s/menus/%s/item_unavailabilities/%s", req.BrandID, req.MenuID, req.SiteID)

	var (
		errResp dto.CustomError
	)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		SetError(errResp).
		SetBody(req.UpdateUnavailabilitesRequestBody).
		Put(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return errResp.Error()
	}
	return nil
}
