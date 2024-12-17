package http

import (
	"context"
	"fmt"
	models2 "github.com/kwaaka-team/orders-core/pkg/express24/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (c *Client) Auth(ctx context.Context) error {

	path := "/api/external/auth"

	var resp models2.AuthResponse
	var errResponse models2.ErrorResponse

	rsp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(models2.AuthRequest{Login: c.Username, Password: c.Password}).
		SetResult(&resp).
		SetError(&errResponse).
		Post(path)
	if err != nil {
		return err
	}

	if rsp.IsError() {
		return errors.New(errResponse.ToString())
	}

	c.restyClient.SetHeader(authHeader, fmt.Sprintf("express24:%s", resp.Token))
	log.Info().Msgf("express24:%s", resp.Token)

	return nil
}

func (c *Client) Close(ctx context.Context) {
	close(c.quit)
}
