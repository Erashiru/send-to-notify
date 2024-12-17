package http

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) Auth(ctx context.Context) error {

	path := "/api/1/access_token"

	var resp models.AuthResponse
	var errResponse models.ErrorResponse

	rsp, err := c.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(models.AuthRequest{ApiLogin: c.apiKey}).
		SetResult(&resp).
		SetError(&errResponse).
		Post(path)
	if err != nil {
		return err
	}

	if rsp.IsError() {
		return errors.New(errResponse.Description)
	}

	c.cli.SetHeader("Authorization", fmt.Sprintf("Bearer %s", resp.AccessToken))
	return nil
}

func (c *Client) Close(ctx context.Context) {
	close(c.quit)
}
