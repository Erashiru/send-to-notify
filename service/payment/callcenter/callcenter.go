package callcenter

import (
	"context"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"time"
)

type CallCenterClient interface {
	CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error)
	GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error)
	CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error)
	OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error)
	SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error
	GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error)
}

type CallCenterService struct {
	CallCenterClient
}

func (c CallCenterService) RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error) {
	return models.PaymentOrder{}, models.RefundResponse{}, models.ErrUnsupportedMethod
}

func (c CallCenterService) CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	return models.PaymentOrder{}, models.ErrUnsupportedMethod
}

func (c CallCenterService) GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error) {
	return "", models.ErrUnsupportedMethod
}

type CallCenterClientImpl struct {
}

func NewCallCenterService() (*CallCenterService, error) {
	return &CallCenterService{CallCenterClient: &CallCenterClientImpl{}}, nil
}

func (cl *CallCenterClientImpl) GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error) {
	return nil, models.ErrUnsupportedMethod
}

func (cl *CallCenterClientImpl) CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error) {
	return models.PaymentEvent{}, models.ErrUnsupportedMethod
}
func (cl *CallCenterClientImpl) SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error {
	return models.ErrUnsupportedMethod
}
func (cl *CallCenterClientImpl) CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	return models.PaymentSystemCustomer{}, models.ErrUnsupportedMethod
}

func (cl *CallCenterClientImpl) CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	return models.PaymentSystemSubscription{}, models.ErrUnsupportedMethod
}

func (cl *CallCenterClientImpl) GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error) {
	return models.WebhookEvent{}, models.ErrUnsupportedMethod
}

func (cl *CallCenterClientImpl) CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {

	paymentOrder.PaymentOrderStatus = models.UNPAID
	paymentOrder.PaymentStatusHistory = append(paymentOrder.PaymentOrderStatusHistory, models.StatusHistory{Status: models.UNPAID, Time: time.Now()})

	return paymentOrder, nil
}

func (cl *CallCenterClientImpl) OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error) {
	return models.ApplePaySessionOpenResponse{}, models.ErrUnsupportedMethod
}
