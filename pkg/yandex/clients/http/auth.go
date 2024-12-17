package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/yandex/models"
	"github.com/pkg/errors"
)

func (c *Client) Auth(ctx context.Context) error {
	path := "/oauth2/token"

	var resp models.AuthResponse
	var errResponse models.ErrorResponse

	rsp, err := c.cli.R().
		SetContext(ctx).
		SetHeaders(map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
			"Accept":       "application/json",
		}).
		SetFormData(map[string]string{
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
		}).
		SetResult(&resp).
		SetError(&errResponse).
		Post(path)
	if err != nil {
		return err
	}

	if rsp.IsError() {
		return errors.New(fmt.Sprintf("code: %s, message: %s", errResponse.Code, errResponse.Message))
	}

	c.cli.SetHeader("Authorization", fmt.Sprintf("Bearer %s", resp.AccessToken))

	return nil
}

func (c *Client) Close(ctx context.Context) {
	close(c.quit)
}
