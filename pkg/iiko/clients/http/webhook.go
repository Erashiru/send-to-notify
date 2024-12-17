package http

import (
	"context"
	"fmt"

	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) GetWebhookSetting(ctx context.Context, organizationID string) (models.GetWebhookSettingResponse, error) {

	path := "/api/1/webhooks/settings"

	var (
		response    models.GetWebhookSettingResponse
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetBody(models.GetWebhookSettingRequest{OrganizationID: organizationID}).
		SetResult(&response).
		SetError(&errResponse).
		Post(path)

	if err != nil {
		return models.GetWebhookSettingResponse{}, errResponse
	}

	if resp.IsError() {
		return models.GetWebhookSettingResponse{}, fmt.Errorf("%s %w", resp.Error(), ErrResponse)
	}

	return response, nil
}

func (c *Client) UpdateWebhookSetting(ctx context.Context, request models.UpdateWebhookRequest) (models.CorID, error) {

	path := "/api/1/webhooks/update_settings"

	var (
		response    models.CorID
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(request).
		Post(path)

	if err != nil {
		return models.CorID{}, err
	}

	if resp.IsError() {
		return models.CorID{}, fmt.Errorf("%s %w", resp.Error(), ErrResponse)
	}

	return response, nil
}
