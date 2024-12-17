package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/yandex/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (c *Client) MenuImportInitiation(ctx context.Context, req models.MenuInitiationRequest) error {
	path := "/menu/import/initiation"

	var errResponse models.ErrorResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetError(&errResponse).
		SetBody(req).
		Post(path)

	if err != nil {
		return err
	}

	log.Info().Msgf("menu import initiation path: %s, request: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return errors.New(fmt.Sprintf("httpStatus: %s, message: %s, code: %v", resp.Status(), errResponse.Message, errResponse.Code))
	}

	log.Info().Msgf("menu import initiation code: %d, response: %s", resp.RawResponse.StatusCode, string(resp.Body()))

	return nil
}
