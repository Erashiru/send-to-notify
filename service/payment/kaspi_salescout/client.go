package kaspi_salescout

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/service/payment/kaspi_salescout/dto"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type KaspiSaleScoutService struct {
	KaspiSaleScoutClient
}

func (k *KaspiSaleScoutImpl) CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	return models.PaymentSystemCustomer{}, models.ErrUnsupportedMethod
}

func (k *KaspiSaleScoutImpl) CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	return models.PaymentSystemSubscription{}, models.ErrUnsupportedMethod
}

func (k *KaspiSaleScoutImpl) GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error) {
	return models.WebhookEvent{}, models.ErrUnsupportedMethod
}

func (k *KaspiSaleScoutImpl) CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error) {
	return models.PaymentEvent{}, models.ErrUnsupportedMethod
}

func (k *KaspiSaleScoutImpl) OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error) {
	return models.ApplePaySessionOpenResponse{}, models.ErrUnsupportedMethod
}

func (k *KaspiSaleScoutImpl) SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error {
	return models.ErrUnsupportedMethod
}

func (k *KaspiSaleScoutImpl) GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error) {
	return nil, models.ErrUnsupportedMethod
}

type KaspiSaleScoutClient interface {
	GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error)
	CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error)
	GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error)
	CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error)
	OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error)
	SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error
	GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error)
	CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error)
}

type KaspiSaleScoutImpl struct {
	restyClient *resty.Client
	merchantID  string
}

func NewKaspiSaleScoutService(baseUrl, token, merchantID string) (*KaspiSaleScoutService, error) {
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

	return &KaspiSaleScoutService{
		&KaspiSaleScoutImpl{
			restyClient: client,
			merchantID:  merchantID,
		},
	}, nil
}

func (cl *KaspiSaleScoutImpl) GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error) {
	var (
		response dto.PaymentStatusResponse
		errResp  dto.ErrorResponse
	)

	path := fmt.Sprintf("/api/kaspi-api/status/%s", paymentID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetResult(&response).
		Get(path)

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("kaspi saleScout cli: %s", errResp.Message)
	}

	return response.Status, nil
}

func (cl *KaspiSaleScoutImpl) CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {

	switch paymentOrder.OrderSource {
	case models2.KWAAKA_ADMIN.String():
		paymentOrder.PaymentOrderStatus = models.REDIRECTED_TO_PAYMENT
		paymentOrder.PaymentOrderStatusHistory = append(paymentOrder.PaymentOrderStatusHistory, models.StatusHistory{
			Status: models.REDIRECTED_TO_PAYMENT,
			Time:   time.Now().UTC(),
		})

		paymentOrder.CheckoutURL = fmt.Sprintf("https://pay.kwaaka.direct?sum=%s&order_id=%s", strconv.Itoa(paymentOrder.Amount/100), paymentOrder.OrderID)
		return paymentOrder, nil

	case models2.QRMENU.String():
		var err error

		paymentOrder, err = cl.CreatePaymentLink(ctx, paymentOrder)
		if err != nil {
			return models.PaymentOrder{}, err
		}
		return paymentOrder, nil
	}

	return paymentOrder, fmt.Errorf("invalid payment order source: %s", paymentOrder.OrderSource)
}

func (cl *KaspiSaleScoutImpl) CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	var (
		response dto.CreatePaymentOrderResponse
		errResp  dto.ErrorResponse
	)

	path := "/api/kaspi-api/create-link"

	req := cl.fromPaymentOrder(paymentOrder)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Post(path)

	if err != nil {
		return models.PaymentOrder{}, err
	}

	if resp.IsError() {
		return models.PaymentOrder{}, fmt.Errorf("kaspi saleScout cli create payment order: %s", errResp.Message)
	}

	return cl.toPaymentOrder(response, paymentOrder)
}

func (cl *KaspiSaleScoutImpl) RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error) {
	log.Info().Msgf("starting refund request for kaspi_salescout: amount: %d", amount)

	var (
		response dto.RefundResponse
		errResp  dto.ErrorResponse
		path     = "/api/kaspi-api/refund"
	)

	req, err := cl.toRefundRequest(paymentOrder, float64(amount))
	if err != nil {
		return models.PaymentOrder{}, models.RefundResponse{}, err
	}

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Post(path)
	if err != nil {
		return models.PaymentOrder{}, models.RefundResponse{}, fmt.Errorf("kaspi salescout request error: %w", err)
	}
	if resp.IsError() {
		log.Error().Msgf("kaspi salescout refund error: code: %d, message: %s", errResp.StatusCode, errResp.Message)
		return models.PaymentOrder{}, models.RefundResponse{}, fmt.Errorf("kaspi salescout refund error: code: %d, message: %s", errResp.StatusCode, errResp.Message)
	}

	paymentOrder = cl.fromRefundResponseToPaymentOrder(response, paymentOrder)

	refundResponse := cl.fromRefundResponseToModel(response, paymentOrder)

	return paymentOrder, refundResponse, nil
}

func (cl *KaspiSaleScoutImpl) fromRefundResponseToModel(response dto.RefundResponse, paymentOrder models.PaymentOrder) models.RefundResponse {
	var currentTime = time.Now().UTC()

	return models.RefundResponse{
		ID:        strconv.Itoa(response.ReturnOperationID),
		PaymentID: paymentOrder.PaymentOrderID,
		Status:    response.Status,
		CreatedAt: fmt.Sprintf(currentTime.Format("2006-01-02 15:04:05")),
	}
}

func (cl *KaspiSaleScoutImpl) fromRefundResponseToPaymentOrder(response dto.RefundResponse, paymentOrder models.PaymentOrder) models.PaymentOrder {
	paymentOrder.PaymentOrderStatus = "REFUNDED"
	paymentOrder.UpdatedAt = time.Now().UTC()
	paymentOrder.PaymentInvoiceID = strconv.Itoa(response.ReturnOperationID)
	paymentOrder.PaymentOrderStatusHistory = append(paymentOrder.PaymentOrderStatusHistory, models.StatusHistory{
		Status: "REFUNDED",
		Time:   time.Now().UTC(),
	})
	return paymentOrder
}

func (cl *KaspiSaleScoutImpl) toRefundRequest(paymentOrder models.PaymentOrder, amount float64) (dto.RefundRequest, error) {
	paymentId, err := strconv.Atoi(paymentOrder.PaymentOrderID)
	if err != nil {
		return dto.RefundRequest{}, err
	}
	return dto.RefundRequest{
		Amount:     amount,
		PaymentID:  paymentId,
		MerchantID: cl.merchantID,
	}, nil
}

func (cl *KaspiSaleScoutImpl) fromPaymentOrder(req models.PaymentOrder) dto.CreatePaymentOrderRequest {
	return dto.CreatePaymentOrderRequest{
		Amount:        float64(req.Amount / 100),
		TransactionID: req.OrderID,
		MerchantID:    cl.merchantID,
	}
}

func (cl *KaspiSaleScoutImpl) toPaymentOrder(resp dto.CreatePaymentOrderResponse, req models.PaymentOrder) (models.PaymentOrder, error) {
	status := models.UNPAID
	if resp.Status == "Error" {
		status = "ERROR"
	}
	req.PaymentOrderID = strconv.Itoa(resp.PaymentId)
	req.PaymentOrderStatus = status
	req.PaymentOrderStatusHistory = append(req.PaymentOrderStatusHistory, models.StatusHistory{
		Status: status,
		Time:   time.Now().UTC(),
	})
	req.CheckoutURL = resp.PaymentLink

	return req, nil
}
