package http

import (
	"context"
)

func (c *Client) auth(ctx context.Context) {
	c.cli.SetQueryParam("token", c.token)
}
