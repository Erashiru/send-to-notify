package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/talabat/models"
	"github.com/pkg/errors"
)

func (c *Client) AuthMenu(ctx context.Context) error {
	path := "/api/User/GetToken"

	var resp models.AuthMenuResponse
	var errResponse models.ErrorResponse

	rsp, err := c.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(models.AuthMenuRequest{
			Username: c.username,
			Password: c.password,
		}).
		SetResult(&resp).
		SetError(&errResponse).
		Post(path)
	if err != nil {
		return errors.Wrapf(err, rsp.String())
	}

	if rsp.IsError() {
		return errors.New(fmt.Sprintf("type: %s, title: %s, detail: %s, instance: %s", errResponse.Type, errResponse.Title, errResponse.Detail, errResponse.Instance))
	}

	c.cli.SetHeader("Authorization", fmt.Sprintf("Bearer %s", resp.AccessToken))

	return nil
}

func (c *Client) AuthMW(ctx context.Context) error {
	path := "/v2/login"

	var resp models.AuthMWResponse
	var errResponse models.AuthMWErrorResponse

	rsp, err := c.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetFormData(map[string]string{
			"grant_type": "client_credentials",
			"username":   c.username,
			"password":   c.password,
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
