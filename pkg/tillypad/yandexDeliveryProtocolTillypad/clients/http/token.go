package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/rs/zerolog/log"
)

func (c *clientImpl) GetAccessToken(ctx context.Context, clientId, clientSecret string) (string, error) {
	path := c.pathPrefix + "/security/oauth/token"

	var (
		result struct {
			AccessToken string `json:"access_token"`
		}
		errResponse models.ErrorResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetError(&errResponse).
		SetFormData(map[string]string{
			"client_id":     clientId,
			"client_secret": clientSecret,
			"grant_type":    "client_credentials",
		}).
		SetResult(&result).
		Post(path)

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		log.Trace().Msgf("(GetAccessToken) request path: %s; request itself: %+v", resp.Request.URL, resp.Request.Body)
		return "", fmt.Errorf("(GetAccessToken) response: status code: %d; error %v: %v", resp.StatusCode(), errResponse, resp.Error())
	}

	return result.AccessToken, nil
}
