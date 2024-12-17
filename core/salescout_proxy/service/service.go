package salescout_proxy

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/service/payment/kaspi_salescout/dto"
	"github.com/rs/zerolog/log"
)

type Service struct {
	saleScoutCli *resty.Client
	merchantID   string
}

func NewKaspiSaleScoutService(baseUrl, token, merchantID string) (*Service, error) {
	if baseUrl == "" {
		return nil, fmt.Errorf("base URL could not be empty")
	}

	client := resty.New().
		SetBaseURL(baseUrl).
		SetHeaders(map[string]string{
			"Content-Type":  "application/json; charset=utf-8",
			"Accept":        "application/json; charset=utf-8",
			"authorization": "Bearer " + token,
		})

	return &Service{
		saleScoutCli: client,
		merchantID:   merchantID,
	}, nil
}

func (cl *Service) CreatePaymentLink(ctx context.Context, req dto.CreatePaymentOrderRequest) (dto.CreatePaymentOrderResponse, error) {
	var (
		response dto.CreatePaymentOrderResponse
		errResp  dto.ErrorResponse
	)

	path := "/api/kaspi-api/create-link"

	resp, err := cl.saleScoutCli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto.CreatePaymentOrderResponse{}, err
	}

	if resp.IsError() {
		return dto.CreatePaymentOrderResponse{}, fmt.Errorf("kaspi saleScout cli create payment order: %s", errResp.Message)
	}

	return response, nil
}

func (cl *Service) CreatePaymentToken(ctx context.Context, req dto.CreatePaymentOrderRequest) (dto.CreatePaymentTokenResponse, error) {
	var (
		response dto.CreatePaymentTokenResponse
		errResp  dto.ErrorResponse
	)

	path := "/api/kaspi-api/create-token"

	resp, err := cl.saleScoutCli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto.CreatePaymentTokenResponse{}, err
	}

	if resp.IsError() {
		return dto.CreatePaymentTokenResponse{}, fmt.Errorf("kaspi saleScout cli create payment order: %s", errResp.Message)
	}

	return response, nil
}

func (cl *Service) GetPaymentStatusByID(ctx context.Context, paymentID string) (dto.PaymentStatusResponse, error) {
	var (
		response dto.PaymentStatusResponse
		errResp  dto.ErrorResponse
	)

	path := fmt.Sprintf("/api/kaspi-api/status/%s", paymentID)

	resp, err := cl.saleScoutCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetResult(&response).
		Get(path)

	if err != nil {
		return dto.PaymentStatusResponse{}, err
	}

	if resp.IsError() {
		return dto.PaymentStatusResponse{}, fmt.Errorf("kaspi saleScout cli: %s", errResp.Message)
	}

	return response, nil
}

func (cl *Service) RefundPayment(ctx context.Context, req dto.RefundRequest) (dto.RefundResponse, error) {
	var (
		response dto.RefundResponse
		errResp  dto.ErrorResponse
		path     = "/api/kaspi-api/refund"
	)

	resp, err := cl.saleScoutCli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Post(path)
	if err != nil {
		return dto.RefundResponse{}, fmt.Errorf("starterApp salescout request error: %w", err)
	}
	if resp.IsError() {
		log.Error().Msgf("starterApp salescout refund error: code: %d, message: %s", errResp.StatusCode, errResp.Message)
		return dto.RefundResponse{}, fmt.Errorf("starterApp salescout refund error: code: %d, message: %s", errResp.StatusCode, errResp.Message)
	}

	return response, nil
}
