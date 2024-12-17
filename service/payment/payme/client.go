package payme

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"github.com/kwaaka-team/orders-core/service/payment/payme/dto"
	"github.com/pkg/errors"
	"time"
)

type PaymeService struct {
	PaymeClient
}

func (p PaymeService) RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error) {
	return models.PaymentOrder{}, models.RefundResponse{}, models.ErrUnsupportedMethod
}

func (p PaymeService) CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	return models.PaymentOrder{}, models.ErrUnsupportedMethod
}

func (p PaymeService) GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error) {
	return "", models.ErrUnsupportedMethod
}

type PaymeClient interface {
	CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error)
	GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error)
	CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error)
	OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error)
	SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error
	GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error)
}

type PaymeClientImpl struct {
	restyClient *resty.Client
}

func NewPaymeService(baseUrl, apiKey string) (*PaymeService, error) {
	if baseUrl == "" {
		return nil, errors.New("base URL could not be empty")
	}

	client := resty.New().
		SetBaseURL(baseUrl).
		SetHeaders(map[string]string{
			"Content-Type": "application/json; charset=utf-8",
			"Accept":       "application/json; charset=utf-8",
			"X-auth":       apiKey,
		})

	return &PaymeService{
		&PaymeClientImpl{
			restyClient: client,
		},
	}, nil
}

func (cl *PaymeClientImpl) GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error) {
	return nil, models.ErrUnsupportedMethod
}

func (cl *PaymeClientImpl) CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	var (
		response dto.CreateReceiptResponse
		errResp  dto.ErrorResponse
	)

	req := cl.fromPaymentOrder(paymentOrder)

	req.ID = cl.generateRequestID()

	req.Method = "receipts.create"

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Post("")

	if err != nil {
		return models.PaymentOrder{}, err
	}

	if resp.IsError() {
		return models.PaymentOrder{}, fmt.Errorf("ioka cli create payment order: %s", errResp.Error.Message)
	}

	return cl.toPaymentOrder(response, paymentOrder)
}

func (cl *PaymeClientImpl) SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error {
	var (
		response dto.SendPaymentOrderToCustomerResponse
		errResp  dto.ErrorResponse
	)

	req := cl.fromPaymentOrderToNotification(paymentOrder)

	req.Id = cl.generateRequestID()

	req.Method = "receipts.send"

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Post("")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("ioka cli create payment order: %s", errResp.Error.Message)
	}

	return nil
}

func (cl *PaymeClientImpl) GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error) {
	req, ok := r.(dto.WebhookEvent)
	if !ok {
		return models.WebhookEvent{}, errors.New("casting error")
	}

	return cl.toWebhookEventModel(req)
}

func (cl *PaymeClientImpl) CreateCustomer(ctx context.Context, customer models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	return models.PaymentSystemCustomer{}, models.ErrUnsupportedMethod
}

func (cl *PaymeClientImpl) CreateSubscription(ctx context.Context, subscription models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	return models.PaymentSystemSubscription{}, models.ErrUnsupportedMethod
}

func (cl *PaymeClientImpl) OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error) {
	return models.ApplePaySessionOpenResponse{}, models.ErrUnsupportedMethod
}

func (cl *PaymeClientImpl) CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error) {
	return models.PaymentEvent{}, models.ErrUnsupportedMethod
}

func (cl *PaymeClientImpl) fromPaymentOrder(req models.PaymentOrder) dto.CreateReceiptRequest {
	res := dto.CreateReceiptRequest{
		Params: dto.CreateReceiptRequestParams{
			Amount: req.Amount,
			Account: dto.Account{
				OrderID: req.CartID,
			},
		},
	}
	return res
}

func (cl *PaymeClientImpl) toPaymentOrder(resp dto.CreateReceiptResponse, req models.PaymentOrder) (models.PaymentOrder, error) {
	req.PaymentOrderID = resp.Result.Receipt.ID

	return req, nil
}

func (cl *PaymeClientImpl) generateRequestID() int64 {
	// Capture the current timestamp in a special format
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)

	// Obtain a unique code using UUID
	uniqueID := uuid.New().ID()

	// Combine the time-based component and the UUID to create a globally unique ID
	id := currentTimestamp + int64(uniqueID)

	return id
}

func (cl *PaymeClientImpl) fromPaymentOrderToNotification(req models.PaymentOrder) dto.SendPaymentOrderToCustomerRequest {
	res := dto.SendPaymentOrderToCustomerRequest{
		Params: dto.SendPaymentOrderToCustomerRequestParams{
			Id:    req.PaymentOrderID,
			Phone: req.CustomerPhoneNumber,
		},
	}
	return res
}

func (cl *PaymeClientImpl) toWebhookEventModel(r dto.WebhookEvent) (models.WebhookEvent, error) {
	switch r.Method {
	case "CreateTransaction":
		return models.WebhookEvent{
			Event: models.PAYME_CREATE_TRANSACTION_METHOD,
			PaymentEvent: models.PaymentEvent{
				OrderID: r.Params.ID,
			},
		}, nil

	case "CheckPerformTransaction":
		return models.WebhookEvent{
			Event: models.PAYME_CHECK_PERFORM_TRANSACTION_METHOD,
		}, nil

	case "PerformTransaction":
		return models.WebhookEvent{
			Event: models.PAYME_PERFORM_TRANSACTION_METHOD,
			PaymentEvent: models.PaymentEvent{
				OrderID: r.Params.ID,
			},
		}, nil
	}

	return models.WebhookEvent{}, models.ErrUnsupportedPaymeMethod
}
