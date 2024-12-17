package http

import (
	"context"
	"fmt"

	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) GetMenu(ctx context.Context, organizationID string) (models.GetMenuResponse, error) {
	path := "/api/1/nomenclature"

	var (
		response    models.GetMenuResponse
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(models.OrganizationRequest{OrganizationID: organizationID}).
		Post(path)

	if err != nil {
		return models.GetMenuResponse{}, errResponse
	}

	if resp.IsError() {
		return models.GetMenuResponse{}, fmt.Errorf("%s %w", resp.Error(), ErrResponse)
	}

	return response, nil
}

func (c *Client) GetExternalMenu(ctx context.Context, organizationID, externalMenuID, priceCategoryId string) (models.GetExternalMenuResponse, error) {
	path := "/api/2/menu/by_id"

	var (
		response    models.GetExternalMenuResponse
		errResponse models.ExternalErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(models.GetExternalMenuRequest{
			ExternalMenuID:  externalMenuID,
			OrganizationIDS: []string{organizationID},
			PriceCategoryId: priceCategoryId,
		}).
		Post(path)

	if err != nil {
		return models.GetExternalMenuResponse{}, err
	}

	if resp.IsError() {
		return models.GetExternalMenuResponse{}, fmt.Errorf("%s %w", resp.Error(), ErrResponse)
	}

	return response, nil
}
