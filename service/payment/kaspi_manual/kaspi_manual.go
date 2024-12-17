package kaspi_manual

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type KaspiManualClient interface {
	CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error)
	GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error)
	CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error)
	OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error)
	SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error
	GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error)
}

type KaspiManualService struct {
	KaspiManualClient
}

func (k KaspiManualService) RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error) {
	return models.PaymentOrder{}, models.RefundResponse{}, models.ErrUnsupportedMethod
}

func (k KaspiManualService) CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	return models.PaymentOrder{}, models.ErrUnsupportedMethod
}

func (k KaspiManualService) GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error) {
	return "", models.ErrUnsupportedMethod
}

type KaspiManualClientImpl struct {
	paymentUrl string
}

func NewKaspiManualService(paymentUrl string) (*KaspiManualService, error) {
	if paymentUrl == "" {
		return nil, errors.New("store kaspi payment url is empty")
	}
	return &KaspiManualService{KaspiManualClient: &KaspiManualClientImpl{
		paymentUrl: paymentUrl,
	}}, nil
}

func (cl *KaspiManualClientImpl) GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error) {
	return nil, models.ErrUnsupportedMethod
}

func (cl *KaspiManualClientImpl) CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error) {
	return models.PaymentEvent{}, models.ErrUnsupportedMethod
}
func (cl *KaspiManualClientImpl) SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error {
	return models.ErrUnsupportedMethod
}
func (cl *KaspiManualClientImpl) CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	return models.PaymentSystemCustomer{}, models.ErrUnsupportedMethod
}

func (cl *KaspiManualClientImpl) CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	return models.PaymentSystemSubscription{}, models.ErrUnsupportedMethod
}

func (cl *KaspiManualClientImpl) GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error) {
	return models.WebhookEvent{}, models.ErrUnsupportedMethod
}

func (cl *KaspiManualClientImpl) CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	paymentOrder.PaymentOrderStatus = models.UNPAID
	paymentOrder.PaymentStatusHistory = append(paymentOrder.PaymentOrderStatusHistory, models.StatusHistory{Status: models.UNPAID, Time: time.Now()})
	paymentOrder.CheckoutURL = fmt.Sprintf("https://pay.kwaaka.direct?sum=%s&link=%s", strconv.Itoa(paymentOrder.Amount/100), cl.paymentUrl)
	return paymentOrder, nil
}

func (cl *KaspiManualClientImpl) OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error) {
	return models.ApplePaySessionOpenResponse{}, models.ErrUnsupportedMethod
}
