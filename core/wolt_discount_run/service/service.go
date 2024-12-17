package service

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/wolt_discount_run/models"
	"github.com/rs/zerolog/log"
)

const jsonType = "application/json"

type Service struct {
	baseUrl string
	cli     *resty.Client
}

func NewService(baseUrl, token string) *Service {
	cli := resty.New().
		SetBaseURL(baseUrl).
		SetHeaders(map[string]string{
			"Content-Type":  jsonType,
			"Accept":        jsonType,
			"Authorization": token,
		})

	return &Service{
		baseUrl: baseUrl,
		cli:     cli,
	}
}

func (s *Service) WoltDiscountRun(ctx context.Context, req models.DiscountRunRequest) error {
	path := "/v3/wolt-discount-run"

	var errResp errors.ErrorResponse

	resp, err := s.cli.R().
		SetContext(ctx).
		SetError(&errResp).
		SetBody(req).
		Post(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("wolt discount run path: %s, request: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return fmt.Errorf("wolt discount run response status: %d, error: %+v, %+v", resp.StatusCode(), errResp, resp.Error())
	}

	log.Info().Msgf("wolt discount run response status: %d, %s", resp.StatusCode(), string(resp.Body()))

	return nil
}
