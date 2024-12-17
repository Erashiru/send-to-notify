package http

import (
	"context"
	"fmt"

	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) GetOrganizations(ctx context.Context) ([]models.Info, error) {

	path := "/api/1/organizations"

	var result struct {
		Organizations []models.Info `json:"organizations"`
	}

	var (
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetBody(models.AdditionalInfo{ReturnAddInfo: true}).
		SetResult(&result).
		SetError(&errResponse).
		Post(path)

	if err != nil {
		return nil, errResponse
	}

	if resp.IsError() {
		return nil, fmt.Errorf("%s %w", resp.Error(), ErrResponse)
	}

	return result.Organizations, nil
}
