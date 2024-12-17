package cash

import (
	"context"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"time"
)

type CashService struct {
	CashClient
}

func (c CashService) RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error) {
	return models.PaymentOrder{}, models.RefundResponse{}, models.ErrUnsupportedMethod
}

func (c CashService) CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	return models.PaymentOrder{}, models.ErrUnsupportedMethod
}

func (c CashService) GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error) {
	return "", models.ErrUnsupportedMethod
}

type CashClient interface {
	CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error)
	GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error)
	CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error)
	OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error)
	SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error
	GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error)
}

type CashClientImpl struct{}

func NewCashService() (*CashService, error) {
	return &CashService{}, nil
}

func (cl *CashClientImpl) GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error) {
	return nil, models.ErrUnsupportedMethod
}

func (cl *CashClientImpl) CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error) {
	return models.PaymentEvent{}, models.ErrUnsupportedMethod
}
func (cl *CashClientImpl) SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error {
	return models.ErrUnsupportedMethod
}
func (cl *CashClientImpl) CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	return models.PaymentSystemCustomer{}, models.ErrUnsupportedMethod
}

func (cl *CashClientImpl) CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	return models.PaymentSystemSubscription{}, models.ErrUnsupportedMethod
}

func (cl *CashClientImpl) GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error) {
	return models.WebhookEvent{}, models.ErrUnsupportedMethod
}

func (cl *CashClientImpl) CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {

	paymentOrder.PaymentOrderStatus = models.UNPAID
	paymentOrder.PaymentStatusHistory = append(paymentOrder.PaymentOrderStatusHistory, models.StatusHistory{Status: models.UNPAID, Time: time.Now()})

	return paymentOrder, nil
}

func (cl *CashClientImpl) OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error) {
	return models.ApplePaySessionOpenResponse{}, models.ErrUnsupportedMethod
}
