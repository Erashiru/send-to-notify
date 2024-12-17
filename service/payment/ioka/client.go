package ioka

import (
	"context"
	"fmt"
	"github.com/aws/smithy-go/time"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/service/payment/ioka/dto"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"github.com/pkg/errors"
	"strconv"
)

type IokaService struct {
	IokaClient
}

func (i IokaService) CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	return models.PaymentOrder{}, models.ErrUnsupportedMethod
}

func (i IokaService) GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error) {
	return "", models.ErrUnsupportedMethod
}

type IokaClient interface {
	CreateCustomer(ctx context.Context, customer models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	GetCustomerByID(ctx context.Context, customerID string) (dto.GetCustomerByIdResponse, error)
	GetCustomers(ctx context.Context) ([]dto.GetCustomerByIdResponse, error)
	DeleteCustomerByID(ctx context.Context, customerID string) error
	CreateSubscription(ctx context.Context, subscription models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error)
	CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error)
	OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error)

	UpdateSubscription(ctx context.Context, req dto.CreateSubscriptionRequest) (dto.CreateSubscriptionResponse, error)
	ChangeSubscriptionStatus(ctx context.Context, subscriptionID string, req dto.ChangeSubscriptionStatusRequest) (dto.CreateSubscriptionResponse, error)
	GetSubscriptions(ctx context.Context) ([]dto.CreateSubscriptionResponse, error)
	GetSubscriptionByID(ctx context.Context, subscriptionID string) (dto.CreateSubscriptionResponse, error)
	GetSubscriptionPayments(ctx context.Context, subscriptionID string) ([]dto.GetSubscriptionPaymentsResponse, error)
	GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error)
	SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error
	GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error)
	RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error)
}

type IokaClientImpl struct {
	restyClient *resty.Client
}

func NewIokaService(baseUrl, apiKey string) (*IokaService, error) {
	if baseUrl == "" {
		return nil, errors.New("base URL could not be empty")
	}

	client := resty.New().
		SetBaseURL(baseUrl).
		SetHeaders(map[string]string{
			"Content-Type": "application/json; charset=utf-8",
			"Accept":       "application/json; charset=utf-8",
			"API-KEY":      apiKey,
		})

	return &IokaService{
		&IokaClientImpl{
			restyClient: client,
		},
	}, nil
}

func (cl *IokaClientImpl) CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	var (
		response dto.CreatePaymentOrderResponse
		errResp  dto.ErrorResponse
	)

	path := "/v2/orders"

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
		return models.PaymentOrder{}, fmt.Errorf("ioka cli create payment order: %s", errResp.Message)
	}

	return cl.toPaymentOrder(response, paymentOrder)
}

func (cl *IokaClientImpl) GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error) {
	req, ok := r.(dto.WebhookEvent)
	if !ok {
		return models.WebhookEvent{}, errors.New("casting error")
	}

	res, err := cl.toWebhookEventModel(req)
	if err != nil {
		return models.WebhookEvent{}, err
	}

	return res, nil
}

func (cl *IokaClientImpl) CreateCustomer(ctx context.Context, customer models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	var (
		response dto.CreateCustomerResponse
		errResp  dto.ErrorResponse
	)

	path := "/v2/customers"

	req := cl.fromCustomer(customer)
	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Post(path)

	if err != nil {
		return models.PaymentSystemCustomer{}, err
	}

	if resp.IsError() {
		return models.PaymentSystemCustomer{}, fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	res, err := cl.toCustomer(response, customer)
	if err != nil {
		return models.PaymentSystemCustomer{}, err
	}

	return res, nil
}

func (cl *IokaClientImpl) CreateSubscription(ctx context.Context, subscription models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	var (
		response dto.CreateSubscriptionResponse
		errResp  dto.ErrorResponse
	)

	path := "/v2/subscriptions"

	req := cl.fromSubscription(subscription)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Post(path)

	if err != nil {
		return models.PaymentSystemSubscription{}, err
	}

	if resp.IsError() {
		return models.PaymentSystemSubscription{}, fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	res, err := cl.toSubscription(response, subscription)
	if err != nil {
		return models.PaymentSystemSubscription{}, err
	}

	return res, nil
}

func (cl *IokaClientImpl) GetCustomerByID(ctx context.Context, customerID string) (dto.GetCustomerByIdResponse, error) {
	var (
		response dto.GetCustomerByIdResponse
		errResp  dto.ErrorResponse
	)

	path := fmt.Sprintf("/v2/customers/%s", customerID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetResult(&response).
		Get(path)

	if err != nil {
		return dto.GetCustomerByIdResponse{}, err
	}

	if resp.IsError() {
		return dto.GetCustomerByIdResponse{}, fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	return response, nil
}

func (cl *IokaClientImpl) DeleteCustomerByID(ctx context.Context, customerID string) error {
	var (
		errResp dto.ErrorResponse
	)

	path := fmt.Sprintf("/v2/customers/%s", customerID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		Delete(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	return nil
}

func (cl *IokaClientImpl) UpdateSubscription(ctx context.Context, req dto.CreateSubscriptionRequest) (dto.CreateSubscriptionResponse, error) {
	var (
		response dto.CreateSubscriptionResponse
		errResp  dto.ErrorResponse
	)

	path := "/v2/subscriptions"

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Patch(path)

	if err != nil {
		return dto.CreateSubscriptionResponse{}, err
	}

	if resp.IsError() {
		return dto.CreateSubscriptionResponse{}, fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	return response, nil
}

func (cl *IokaClientImpl) GetSubscriptions(ctx context.Context) ([]dto.CreateSubscriptionResponse, error) {
	var (
		response []dto.CreateSubscriptionResponse
		errResp  dto.ErrorResponse
	)

	path := "/v2/subscriptions"

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	return response, nil
}

func (cl *IokaClientImpl) GetSubscriptionByID(ctx context.Context, subscriptionID string) (dto.CreateSubscriptionResponse, error) {
	var (
		response dto.CreateSubscriptionResponse
		errResp  dto.ErrorResponse
	)

	path := fmt.Sprintf("/v2/subscriptions/%s", subscriptionID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetResult(&response).
		Get(path)

	if err != nil {
		return dto.CreateSubscriptionResponse{}, err
	}

	if resp.IsError() {
		return dto.CreateSubscriptionResponse{}, fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	return response, nil
}

func (cl *IokaClientImpl) GetSubscriptionPayments(ctx context.Context, subscriptionID string) ([]dto.GetSubscriptionPaymentsResponse, error) {
	var (
		response []dto.GetSubscriptionPaymentsResponse
		errResp  dto.ErrorResponse
	)

	path := fmt.Sprintf("/v2/subscriptions/%s/payments", subscriptionID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	return response, nil
}

func (cl *IokaClientImpl) GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error) {
	var (
		response []dto.GetCustomerCardsResponse
		errResp  dto.ErrorResponse
	)

	path := fmt.Sprintf("/v2/customers/%s/cards", customerID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	return cl.toCards(response), nil
}

func (cl *IokaClientImpl) ChangeSubscriptionStatus(ctx context.Context, subscriptionID string, req dto.ChangeSubscriptionStatusRequest) (dto.CreateSubscriptionResponse, error) {
	var (
		response dto.CreateSubscriptionResponse
		errResp  dto.ErrorResponse
	)

	path := fmt.Sprintf("/v2/subscriptions/%s", subscriptionID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Patch(path)

	if err != nil {
		return dto.CreateSubscriptionResponse{}, err
	}

	if resp.IsError() {
		return dto.CreateSubscriptionResponse{}, fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	return response, nil
}

func (cl *IokaClientImpl) GetCustomers(ctx context.Context) ([]dto.GetCustomerByIdResponse, error) {
	var (
		response []dto.GetCustomerByIdResponse
		errResp  dto.ErrorResponse
	)

	path := "/v2/customers"

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("ioka cli: %s", errResp.Message)
	}

	return response, nil
}

func (cl *IokaClientImpl) CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error) {
	var (
		response dto.PaymentEvent
		errResp  dto.ErrorResponse
	)

	path := fmt.Sprintf("/v2/orders/%s/payments/tool", paymentRequest.OrderID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&paymentRequest).
		SetError(&errResp).
		SetResult(&response).
		Post(path)

	if err != nil {
		return models.PaymentEvent{}, err
	}

	if resp.IsError() {
		return models.PaymentEvent{}, fmt.Errorf("ioka cli create applePay payment: %s", errResp.Message)
	}

	return cl.toPaymentEventModel(response)
}

func (cl *IokaClientImpl) OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error) {
	var (
		response models.ApplePaySessionOpenResponse
		errResp  dto.ErrorResponse
	)

	path := fmt.Sprintf("/v2/payment-methods/%s/apple-pay-session", request.OrderID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&request).
		SetError(&errResp).
		SetResult(&response).
		Post(path)

	if err != nil {
		return models.ApplePaySessionOpenResponse{}, err
	}

	if resp.IsError() {
		return models.ApplePaySessionOpenResponse{}, fmt.Errorf("ioka cli create applePay payment: %s", errResp.Message)
	}

	return response, nil
}

func (cl *IokaClientImpl) RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error) {

	var (
		req      = cl.fromPaymentOrderToRefundRequest(paymentOrder, amount)
		response dto.RefundResponse
		errResp  dto.ErrorResponse
		path     = fmt.Sprintf("/v2/orders/%s/refunds", paymentOrder.PaymentOrderID)
	)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetResult(&response).
		SetError(&errResp).
		Post(path)

	if err != nil {
		traceInfo := resp.Request.TraceInfo()
		return models.PaymentOrder{}, models.RefundResponse{}, fmt.Errorf("request refund error: conn time: %s, server time: %s, response time: %s, error: %w", traceInfo.ConnTime, traceInfo.ServerTime, traceInfo.ResponseTime, err)
	}
	if resp.IsError() {
		return models.PaymentOrder{}, models.RefundResponse{}, fmt.Errorf("ioka/client - fn RefundPayment - resp.IsError(): %s", errResp.Message)
	}

	paymentOrder, err = cl.fromRefundResponseToPaymentOrder(response, paymentOrder)

	refundResponse := cl.fromRefundResponseToModel(response)

	return paymentOrder, refundResponse, err
}

func (cl *IokaClientImpl) SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error {
	return models.ErrUnsupportedMethod
}

func (cl *IokaClientImpl) fromSubscription(req models.PaymentSystemSubscription) dto.CreateSubscriptionRequest {
	return dto.CreateSubscriptionRequest{
		CustomerId:  req.Payer.CustomerID,
		CardId:      req.Payer.CardID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
		ExtraInfo:   req.ExtraInfo,
		NextPay:     req.Schedule.NextPay,
		Step:        req.Schedule.Step,
		Unit:        req.Schedule.Unit,
	}
}

func (cl *IokaClientImpl) toSubscription(iokaResp dto.CreateSubscriptionResponse, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	req.SubscriptionID = iokaResp.ID
	createdAt, err := time.ParseDateTime(iokaResp.CreatedAt)
	if err != nil {
		return models.PaymentSystemSubscription{}, err
	}
	req.CreatedAtPS = createdAt
	req.Schedule.Status = iokaResp.Schedule.Status

	req.Payer = models.Payer{
		Type:          iokaResp.Payer.Type,
		PanMasked:     iokaResp.Payer.PanMasked,
		ExpiryDate:    iokaResp.Payer.ExpiryDate,
		Holder:        iokaResp.Payer.Holder,
		PaymentSystem: iokaResp.Payer.PaymentSystem,
		Emitter:       iokaResp.Payer.Emitter,
		CustomerID:    iokaResp.Payer.CustomerID,
		CardID:        iokaResp.Payer.CardID,
		Email:         iokaResp.Payer.Email,
		Phone:         iokaResp.Payer.Phone,
	}
	return req, nil
}

func (cl *IokaClientImpl) toCards(req []dto.GetCustomerCardsResponse) []models.CustomerCards {
	res := make([]models.CustomerCards, 0, len(req))

	for _, card := range req {
		res = append(res, cl.toCard(card))
	}

	return res
}

func (cl *IokaClientImpl) toCard(req dto.GetCustomerCardsResponse) models.CustomerCards {
	return models.CustomerCards{
		ID:            req.ID,
		CustomerID:    req.CustomerID,
		CreatedAt:     req.CreatedAt,
		PanMasked:     req.PanMasked,
		ExpiryDate:    req.ExpiryDate,
		Holder:        req.Holder,
		PaymentSystem: req.PaymentSystem,
		Emitter:       req.Emitter,
		CvcRequired:   req.CvcRequired,
		Error: models.Error{
			Code:    req.Error.Code,
			Message: req.Error.Message,
		},
		Action: models.Action{
			URL: req.Action.URL,
		},
	}
}

func (cl *IokaClientImpl) fromCustomer(req models.PaymentSystemCustomer) dto.CreateCustomerRequest {
	return dto.CreateCustomerRequest{
		ExternalId:     req.ExternalID,
		Email:          req.Email,
		Phone:          req.Phone,
		PhoneCheckDate: req.PhoneCheckDate,
		Channel:        req.Channel,
	}
}

func (cl *IokaClientImpl) toCustomer(iokaResp dto.CreateCustomerResponse, customer models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	createdAt, err := time.ParseDateTime(iokaResp.Customer.CreatedAt)
	if err != nil {
		return models.PaymentSystemCustomer{}, err
	}

	accounts, err := cl.toAccounts(iokaResp.Customer.Accounts)
	if err != nil {
		return models.PaymentSystemCustomer{}, err
	}

	customer.PaymentSystemCustomerID = iokaResp.Customer.ID
	customer.CreatedAt = createdAt
	customer.Status = iokaResp.Customer.Status
	customer.Email = iokaResp.Customer.Email
	customer.Phone = iokaResp.Customer.Phone
	customer.CheckoutURL = iokaResp.Customer.CheckoutURL
	customer.AccessToken = iokaResp.Customer.AccessToken
	customer.CustomerAccessToken = iokaResp.Customer.AccessToken
	customer.Accounts = accounts

	return customer, nil
}

func (cl *IokaClientImpl) toAccounts(req []dto.Accounts) ([]models.Accounts, error) {
	res := make([]models.Accounts, 0, len(req))

	for _, v := range req {
		acc, err := cl.toAccount(v)
		if err != nil {
			return nil, err
		}
		res = append(res, acc)
	}

	return res, nil
}

func (cl *IokaClientImpl) toAccount(req dto.Accounts) (models.Accounts, error) {
	createdAt, err := time.ParseDateTime(req.CreatedAt)
	if err != nil {
		return models.Accounts{}, err
	}

	return models.Accounts{
		ID:         req.ID,
		ShopID:     req.ShopID,
		CustomerID: req.CustomerID,
		Status:     req.Status,
		Name:       req.Name,
		Amount:     req.Amount,
		Currency:   req.Currency,
		CreatedAt:  createdAt,
		ExternalID: req.ExternalID,
		Resources:  cl.toResources(req.Resources),
	}, nil
}

func (cl *IokaClientImpl) toResources(req []dto.Resources) []models.Resources {
	res := make([]models.Resources, 0, len(req))

	for _, v := range req {
		res = append(res, cl.toResource(v))
	}

	return res
}

func (cl *IokaClientImpl) toResource(req dto.Resources) models.Resources {
	return models.Resources{
		ID:        req.ID,
		Iban:      req.Iban,
		IsDefault: req.IsDefault,
	}
}

func (cl *IokaClientImpl) toWebhookEventModel(r dto.WebhookEvent) (models.WebhookEvent, error) {
	orderEvent, err := cl.toOrderEventModel(r.OrderEvent)
	if err != nil {
		return models.WebhookEvent{}, err
	}

	paymentEvent, err := cl.toPaymentEventModel(r.PaymentEvent)
	if err != nil {
		return models.WebhookEvent{}, err
	}

	refundEvent, err := cl.toRefundEventModel(r.RefundEvent)
	if err != nil {
		return models.WebhookEvent{}, err
	}

	installmentEvent, err := cl.toInstallmentEvent(r.InstallmentEvent)
	if err != nil {
		return models.WebhookEvent{}, err
	}

	customerEvent, err := cl.toCustomerEventModel(r.CustomerEvent)
	if err != nil {
		return models.WebhookEvent{}, err
	}

	cardEvent, err := cl.toCardEventModel(r.CardEvent)
	if err != nil {
		return models.WebhookEvent{}, err
	}

	return models.WebhookEvent{
		Event:            r.Event,
		OrderEvent:       orderEvent,
		PaymentEvent:     paymentEvent,
		RefundEvent:      refundEvent,
		InstallmentEvent: installmentEvent,
		CardEvent:        cardEvent,
		CustomerEvent:    customerEvent,
	}, nil
}

func (cl *IokaClientImpl) toOrderEventModel(r dto.OrderEvent) (models.OrderEvent, error) {
	event := models.OrderEvent{
		ID:             r.ID,
		ShopID:         r.ShopID,
		Status:         r.Status,
		Amount:         r.Amount,
		Currency:       r.Currency,
		CaptureMethod:  r.CaptureMethod,
		ExternalID:     r.ExternalID,
		Description:    r.Description,
		SubscriptionID: r.SubscriptionID,
	}
	if r.CreatedAt != "" {
		createdAt, err := time.ParseDateTime(r.CreatedAt)
		if err != nil {
			return models.OrderEvent{}, err
		}
		event.CreatedAt = createdAt
	}

	if r.DueDate != "" {
		dueDate, err := time.ParseDateTime(r.DueDate)
		if err != nil {
			return models.OrderEvent{}, err
		}
		event.DueDate = dueDate
	}

	return event, nil
}

func (cl *IokaClientImpl) toPaymentEventModel(r dto.PaymentEvent) (models.PaymentEvent, error) {
	event := models.PaymentEvent{
		ID:             r.ID,
		OrderID:        r.OrderID,
		Status:         r.Status,
		ApprovedAmount: strconv.Itoa(r.ApprovedAmount),
		CapturedAmount: strconv.Itoa(r.CapturedAmount),
		RefundedAmount: strconv.Itoa(r.RefundedAmount),
		ProcessingFee:  strconv.Itoa(r.ProcessingFee),
		Payer:          cl.toPayerModel(r.Payer),
		Error:          cl.toErrorModel(r.Error),
		Acquirer:       cl.toAcquirerModel(r.Acquirer),
		Action:         cl.toActionModel(r.Action),
	}
	if r.CreatedAt != "" {
		createdAt, err := time.ParseDateTime(r.CreatedAt)
		if err != nil {
			return models.PaymentEvent{}, err
		}
		event.CreatedAt = createdAt
	}
	return event, nil
}

func (cl *IokaClientImpl) toRefundEventModel(r dto.RefundEvent) (models.RefundEvent, error) {
	event := models.RefundEvent{
		ID:        r.ID,
		PaymentID: r.PaymentID,
		OrderID:   r.OrderID,
		Status:    r.Status,
		Error:     cl.toErrorModel(r.Error),
		Acquirer:  cl.toAcquirerModel(r.Acquirer),
		Amount:    r.Amount,
		Reason:    r.Reason,
		Author:    r.Author,
	}

	if r.CreatedAt != "" {
		createdAt, err := time.ParseDateTime(r.CreatedAt)
		if err != nil {
			return models.RefundEvent{}, err
		}
		event.CreatedAt = createdAt
	}

	return event, nil
}

func (cl *IokaClientImpl) toInstallmentEvent(r dto.InstallmentEvent) (models.InstallmentEvent, error) {
	event := models.InstallmentEvent{
		ID:             r.ID,
		OrderID:        r.OrderID,
		RedirectURL:    r.RedirectURL,
		Status:         r.Status,
		ProcessingFee:  r.ProcessingFee,
		MonthlyPayment: r.MonthlyPayment,
		InterestRate:   r.InterestRate,
		EffectiveRate:  r.EffectiveRate,
		Iin:            r.Iin,
		Phone:          r.Phone,
		Period:         r.Period,
		ContractNumber: r.ContractNumber,
		Error:          cl.toErrorModel(r.Error),
	}

	if r.CreatedAt != "" {
		createdAt, err := time.ParseDateTime(r.CreatedAt)
		if err != nil {
			return models.InstallmentEvent{}, err
		}
		event.CreatedAt = createdAt
	}

	if r.ContractSignedAt != "" {
		contractSignedAt, err := time.ParseDateTime(r.ContractSignedAt)
		if err != nil {
			return models.InstallmentEvent{}, err
		}
		event.ContractSignedAt = contractSignedAt
	}

	return event, nil
}

func (cl *IokaClientImpl) toCustomerEventModel(r dto.CustomerEvent) (models.CustomerEvent, error) {
	event := models.CustomerEvent{
		ID:         r.ID,
		ShopID:     r.ShopID,
		Status:     r.Status,
		ExternalID: r.ExternalID,
		Email:      r.Email,
		Phone:      r.Phone,
	}

	if r.CreatedAt != "" {
		createdAt, err := time.ParseDateTime(r.CreatedAt)
		if err != nil {
			return models.CustomerEvent{}, err
		}
		event.CreatedAt = createdAt
	}

	return event, nil
}

func (cl *IokaClientImpl) toCardEventModel(r dto.CardEvent) (models.CardEvent, error) {
	event := models.CardEvent{
		ID:            r.ID,
		CustomerID:    r.CustomerID,
		Status:        r.Status,
		PanMasked:     r.PanMasked,
		ExpiryDate:    r.ExpiryDate,
		Holder:        r.Holder,
		PaymentSystem: r.PaymentSystem,
		Emitter:       r.Emitter,
		CvcRequired:   r.CvcRequired,
		Error:         cl.toErrorModel(r.Error),
		Action:        cl.toActionModel(r.Action),
	}

	if r.CreatedAt != "" {
		createdAt, err := time.ParseDateTime(r.CreatedAt)
		if err != nil {
			return models.CardEvent{}, err
		}
		event.CreatedAt = createdAt
	}

	return event, nil
}

func (cl *IokaClientImpl) toPayerModel(r dto.Payer) models.Payer {
	return models.Payer{
		Type:          r.Type,
		PanMasked:     r.PanMasked,
		ExpiryDate:    r.ExpiryDate,
		Holder:        r.Holder,
		PaymentSystem: r.PaymentSystem,
		Emitter:       r.Emitter,
		Email:         r.Email,
		Phone:         r.Phone,
		CustomerID:    r.CustomerID,
		CardID:        r.CardID,
	}
}

func (cl *IokaClientImpl) toErrorModel(r dto.Error) models.Error {
	return models.Error{
		Code:    r.Code,
		Message: r.Message,
	}
}

func (cl *IokaClientImpl) toAcquirerModel(r dto.Acquirer) models.Acquirer {
	return models.Acquirer{
		Name:      r.Name,
		Reference: r.Reference,
	}
}

func (cl *IokaClientImpl) toActionModel(r dto.Action) models.Action {
	return models.Action{
		URL: r.URL,
	}
}

func (cl *IokaClientImpl) toPaymentOrder(iokaResp dto.CreatePaymentOrderResponse, req models.PaymentOrder) (models.PaymentOrder, error) {
	req.PaymentOrderID = iokaResp.Order.ID
	createdAt, err := time.ParseDateTime(iokaResp.Order.CreatedAt)
	if err != nil {
		return models.PaymentOrder{}, err
	}
	req.CreatedAtPaymentSystem = createdAt
	req.PaymentOrderStatus = iokaResp.Order.Status
	req.PaymentOrderStatusHistory = append(req.PaymentOrderStatusHistory, models.StatusHistory{
		Status: iokaResp.Order.Status,
		Time:   createdAt,
	})
	req.CustomerID = iokaResp.Order.CustomerID
	req.CardID = iokaResp.Order.CardID
	req.CheckoutURL = iokaResp.Order.CheckoutURL

	return req, nil
}

func (cl *IokaClientImpl) fromRefundResponseToPaymentOrder(iokaRefundResp dto.RefundResponse, req models.PaymentOrder) (models.PaymentOrder, error) {
	createdAt, err := time.ParseDateTime(iokaRefundResp.CreatedAt)
	if err != nil {
		return models.PaymentOrder{}, fmt.Errorf("fn fromRefundResponseToPaymentOrder: couldn't parse created at time in ioka response: %w", err)
	}
	req.PaymentOrderStatus = "REFUNDED"
	req.UpdatedAt = createdAt
	req.PaymentOrderStatusHistory = append(req.PaymentOrderStatusHistory, models.StatusHistory{
		Status: "REFUNDED",
		Time:   createdAt,
	})
	return req, nil
}

func (cl *IokaClientImpl) fromRefundResponseToModel(iokaRefundResp dto.RefundResponse) models.RefundResponse {
	return models.RefundResponse{
		ID:        iokaRefundResp.ID,
		PaymentID: iokaRefundResp.PaymentID,
		OrderID:   iokaRefundResp.OrderID,
		Status:    iokaRefundResp.Status,
		CreatedAt: iokaRefundResp.CreatedAt,
		Error: models.ErrorRefund{
			Code:    iokaRefundResp.Error.Code,
			Message: iokaRefundResp.Error.Message,
		},
		Acquirer: models.AcquirerRefund{
			Name:      iokaRefundResp.Acquirer.Name,
			Reference: iokaRefundResp.Acquirer.Reference,
		},
	}
}

func (cl *IokaClientImpl) fromPaymentOrder(req models.PaymentOrder) dto.CreatePaymentOrderRequest {
	return dto.CreatePaymentOrderRequest{
		Amount:        req.Amount,
		Currency:      req.Currency,
		CaptureMethod: req.CaptureMethod,
		ExternalID:    req.ExternalID,
		Description:   req.Description,
		Mcc:           req.Mcc,
		ExtraInfo: dto.ExtraInfo{
			RestaurantGroupName: req.RestaurantGroupName,
			RestaurantName:      req.RestaurantName,
			CustomerName:        req.CustomerName,
			CustomerPhoneNumber: req.CustomerPhoneNumber,
		},
		Attempts:   req.Attempts,
		DueDate:    req.DueDate,
		CustomerID: req.CustomerID,
		CardID:     req.CardID,
		BackURL:    req.BackURL,
		SuccessURL: req.SuccessURL,
		FailureURL: req.FailureURL,
		Template:   req.Template,
	}
}

func (cl *IokaClientImpl) fromPaymentOrderToRefundRequest(req models.PaymentOrder, amount int) dto.RefundRequest {
	return dto.RefundRequest{
		Amount: amount * 100,
		Reason: req.RefundReason,
	}
}
