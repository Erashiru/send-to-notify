package http

import "context"

func (c *Client) Auth(ctx context.Context) error {
	return nil
}

func (c *Client) Close(ctx context.Context) {
	close(c.quit)
}
