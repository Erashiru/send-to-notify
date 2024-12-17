package http

import (
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/yaros/models"

	"context"
)

func (c *Client) ServerCheck(ctx context.Context) error {
	path := "/check"

	var errResponse models.ErrorResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		Get(path)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("%w - %s", ErrStopList, err)
	}
	return nil
}
