package wooppay

import (
	"context"
	"github.com/aws/smithy-go/time"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"github.com/kwaaka-team/orders-core/service/payment/wooppay/dto"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type WoopPayService struct {
	WoopPayClient
}

func (w WoopPayService) RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error) {
	return models.PaymentOrder{}, models.RefundResponse{}, models.ErrUnsupportedMethod
}

func (w WoopPayService) CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	return models.PaymentOrder{}, models.ErrUnsupportedMethod
}

func (w WoopPayService) GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error) {
	return "", models.ErrUnsupportedMethod
}

type WoopPayClient interface {
	CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error)
	GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error)
	CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error)
	OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error)
	SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error
	GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error)
}

type WoopPayClientImpl struct {
	restyClient     *resty.Client
	resultUrl       string
	merchantService string
	login           string
}

func NewWoopPayService(baseUrl, resultUrl, login, password, merchantService string) (*WoopPayService, error) {
	if baseUrl == "" {
		return nil, errors.New("woopPay base URL could not be empty")
	}

	if login == "" {
		return nil, errors.New("woopPay login could not be empty")
	}

	if password == "" {
		return nil, errors.New("woopPay password could not be empty")
	}

	if resultUrl == "" {
		return nil, errors.New("woopPay result URL could not be empty")
	}

	if merchantService == "" {
		return nil, errors.New("woopPay merchant service could not be empty")
	}

	client := resty.New().
		SetBaseURL(baseUrl).
		SetHeaders(map[string]string{
			"Content-Type": "application/json; charset=utf-8",
			"Accept":       "application/json; charset=utf-8",
		})

	clientImpl := &WoopPayClientImpl{
		restyClient:     client,
		resultUrl:       resultUrl,
		merchantService: merchantService,
		login:           login,
	}
	if err := clientImpl.Auth(context.Background(), login, password); err != nil {
		return nil, err
	}

	return &WoopPayService{
		clientImpl,
	}, nil
}

func (cl *WoopPayClientImpl) Auth(ctx context.Context, login, password string) error {

	path := "/auth"

	var resp dto.AuthResponse
	var errResponse dto.AuthErrorResponse

	rsp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(dto.AuthRequest{
			Login:    login,
			Password: password,
		}).
		SetResult(&resp).
		SetError(&errResponse).
		Post(path)
	if err != nil {
		return err
	}

	if rsp.IsError() {
		return errors.New(errResponse.Error.Title)
	}

	cl.restyClient.SetHeader("Authorization", resp.Data.Token)

	return nil
}

func (cl *WoopPayClientImpl) CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	return models.PaymentSystemCustomer{}, models.ErrUnsupportedMethod
}

func (cl *WoopPayClientImpl) CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	return models.PaymentSystemSubscription{}, models.ErrUnsupportedMethod
}

func (cl *WoopPayClientImpl) GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error) {
	req, ok := r.(dto.WebhookEventRequest)
	if !ok {
		return models.WebhookEvent{}, errors.New("casting error")
	}

	res, err := cl.toWebhookEventModel(req)
	if err != nil {
		return models.WebhookEvent{}, err
	}

	return res, nil
}

func (cl *WoopPayClientImpl) GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error) {
	return nil, models.ErrUnsupportedMethod
}

func (cl *WoopPayClientImpl) CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	path := "/merchant/invoice"

	var resp dto.CreatePaymentInvoiceResponse
	var errResponse dto.CreatePaymentInvoiceErrorResponse

	req := cl.fromPaymentOrder(paymentOrder)

	utils.Beautify("wooppay create invoice request: ", req)

	rsp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(req).
		SetResult(&resp).
		SetError(&errResponse).
		Post(path)

	utils.Beautify("wooppay create invoice response: ", rsp)

	if err != nil {
		return models.PaymentOrder{}, err
	}

	if rsp.IsError() {
		return models.PaymentOrder{}, errors.New(errResponse.Error.Title)
	}

	paymentOrder = cl.toPaymentOrder(resp, paymentOrder)

	return paymentOrder, nil
}

func (cl *WoopPayClientImpl) CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error) {
	return models.PaymentEvent{}, models.ErrUnsupportedMethod
}

func (cl *WoopPayClientImpl) OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error) {
	return models.ApplePaySessionOpenResponse{}, models.ErrUnsupportedMethod
}

func (cl *WoopPayClientImpl) SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error {
	return models.ErrUnsupportedMethod
}

func (cl *WoopPayClientImpl) fromPaymentOrder(req models.PaymentOrder) dto.CreatePaymentInvoiceRequest {
	return dto.CreatePaymentInvoiceRequest{
		Amount:     int64(req.Amount / 100),
		Currency:   req.Currency,
		Merchant:   cl.login,
		Service:    cl.merchantService,
		ExternalID: req.CartID,
		Regulations: dto.Regulations{
			ClientIdentification: dto.ClientIdentification{
				Phone: cl.fromPhoneNumber(req.CustomerPhoneNumber),
			},
			Behavior: dto.Behavior{
				ResultURL: dto.ResultURL{
					URL:  cl.resultUrl,
					Type: "post",
				},
				SuccessURL:   req.SuccessURL,
				FailureURL:   req.FailureURL,
				AutoRedirect: true,
			},
		},
	}
}

func (cl *WoopPayClientImpl) fromPhoneNumber(phone string) string {
	if strings.Contains(phone, "+") {
		return strings.Trim(phone, "+")
	}
	return phone
}

func (cl *WoopPayClientImpl) toPaymentOrder(resp dto.CreatePaymentInvoiceResponse, req models.PaymentOrder) models.PaymentOrder {
	req.PaymentOrderID = strconv.Itoa(resp.Data.Attributes.OperationID)
	req.PaymentInvoiceID = resp.Data.ID
	req.CheckoutURL = resp.Data.URL
	return req
}

func (cl *WoopPayClientImpl) toWebhookEventModel(r dto.WebhookEventRequest) (models.WebhookEvent, error) {
	event := models.PaymentEvent{
		CapturedAmount: r.Amount,
		Payer: models.Payer{
			PanMasked: r.CardMask,
			CardID:    r.CardHash,
		},
		ProcessingFee: r.Commission,
		OrderID:       r.OperationID,
		Status:        "PAID",
	}
	if r.Date != "" {
		createdAt, err := time.ParseDateTime(strings.Replace(r.Date, " ", "T", 1))
		if err != nil {
			return models.WebhookEvent{}, err
		}
		event.CreatedAt = createdAt
	}

	return models.WebhookEvent{
		PaymentEvent: event,
	}, nil
}
