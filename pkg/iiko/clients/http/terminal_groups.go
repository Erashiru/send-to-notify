package http

import (
	"context"
	"github.com/pkg/errors"

	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) GetTerminalGroups(ctx context.Context, organizationID string) (models.TerminalGroupsResponse, error) {

	type TerminalGroupsRequest struct {
		Organizations []string `json:"organizationIds"`
	}

	var req = TerminalGroupsRequest{
		Organizations: []string{organizationID},
	}

	path := "/api/1/terminal_groups"

	var (
		terminalGroups models.TerminalGroupsResponse
		errR           models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&terminalGroups).
		SetError(&errR).
		Post(path)

	if err != nil {
		return models.TerminalGroupsResponse{}, err
	}

	if resp.IsError() {
		return models.TerminalGroupsResponse{}, errors.New("iiko get terminal groups - stop list error")
	}

	return terminalGroups, nil
}
