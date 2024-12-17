package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"github.com/rs/zerolog/log"
)

func (c *Client) SendNotification(ctx context.Context, notificationInfo models.SendNotificationRequest) error {

	path := "/api/1/notifications/send"

	var errResponse models.ErrorResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(notificationInfo).
		Post(path)
	if err != nil {
		log.Err(err).Msgf("iiko send notifications error for pos order id: %s", notificationInfo.OrderId)
		return err
	}

	if resp.IsError() {
		log.Info().Msgf("iiko send notifications response error: %+v", resp)
		return fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return nil
}
