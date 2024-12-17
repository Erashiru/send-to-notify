package whatsapp

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"github.com/kwaaka-team/orders-core/service/whatsapp"
	"math"
	"net/url"
	"strconv"
)

type WhatsappService struct {
	WhatsappClient
}

func (w WhatsappService) RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error) {
	return models.PaymentOrder{}, models.RefundResponse{}, models.ErrUnsupportedMethod
}

func (w WhatsappService) CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	return models.PaymentOrder{}, models.ErrUnsupportedMethod
}

func (w WhatsappService) GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error) {
	return "", models.ErrUnsupportedMethod
}

type WhatsappClient interface {
	CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error)
	GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error)
	CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error)
	OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error)
	SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error
	GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error)
}

type WhatsappClientImpl struct {
	wppService whatsapp.Service
	paymentUrl string
}

func NewWhatsappService(wpp whatsapp.Service, paymentUrl string) (*WhatsappService, error) {
	if wpp == nil {
		return nil, errors.New("whatsapp service is null")
	}

	return &WhatsappService{
		&WhatsappClientImpl{wppService: wpp, paymentUrl: paymentUrl},
	}, nil
}

func (w *WhatsappClientImpl) CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	return models.PaymentSystemCustomer{}, models.ErrUnsupportedMethod
}

func (w *WhatsappClientImpl) CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	return models.PaymentSystemSubscription{}, models.ErrUnsupportedMethod
}
func (cl *WhatsappClientImpl) GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error) {
	return nil, models.ErrUnsupportedMethod
}

func (w *WhatsappClientImpl) GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error) {
	return w.wppService.GetSystemWebhookEventRequestByPaymentSystemRequest(r)
}

func (w *WhatsappClientImpl) CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	paymentOrder.PaymentOrderStatus = models.UNPAID
	paymentOrder.PaymentOrderID = uuid.New().String()
	paymentOrder.CheckoutURL = fmt.Sprintf("https://pay.kwaaka.direct?sum=%s&link=%s", strconv.Itoa(paymentOrder.Amount/100), w.paymentUrl)

	msg, err := w.formRestaurantMessage(ctx, paymentOrder)
	if err != nil {
		return models.PaymentOrder{}, err
	}

	if err := w.wppService.SendMessage(ctx, paymentOrder.WhatsappPaymentChatId, msg, ""); err != nil {
		return models.PaymentOrder{}, err
	}

	return paymentOrder, nil
}

func (w *WhatsappClientImpl) CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error) {
	return models.PaymentEvent{}, models.ErrUnsupportedMethod
}

func (w *WhatsappClientImpl) OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error) {
	return models.ApplePaySessionOpenResponse{}, models.ErrUnsupportedMethod
}

func (w *WhatsappClientImpl) SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error {
	return models.ErrUnsupportedMethod
}

func (w *WhatsappClientImpl) formRestaurantMessage(ctx context.Context, paymentOrder models.PaymentOrder) (string, error) {
	amount := math.Ceil(float64(paymentOrder.Amount) / 100)
	customerInfo := fmt.Sprintf("[✅] Данные о клиенте:\nИмя: %s\nНомер телефона:%s\nCумма оплаты: %.2f\nID Оплаты: %s\n", paymentOrder.CustomerName, paymentOrder.CustomerPhoneNumber, amount, paymentOrder.PaymentOrderID)

	msg := "Ответьте на это сообщение “Да”, чтобы подтвердить оплату клиентом.\n"

	return url.QueryEscape(customerInfo + msg), nil
}
