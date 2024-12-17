package multicard

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	ocErrs "github.com/kwaaka-team/orders-core/core/errors"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/pkg/multicard"
	"github.com/kwaaka-team/orders-core/pkg/multicard/dto"
	"github.com/kwaaka-team/orders-core/service/order"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	multicardDto "github.com/kwaaka-team/orders-core/service/payment/multicard/dto"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type MultiCardService struct {
	MulticardClient
}

type MulticardClient interface {
	CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error)
	GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error)
	CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error)
	OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error)
	SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error
	GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error)
	GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error)
	RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error)
	CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
}

type MulticardClientImpl struct {
	multicardCli multicard.Client
	log          *zap.SugaredLogger
	cartService  *order.CartServiceImpl
	callbackUrl  string
	storeId      string
	secret       string
}

func NewMultricardService(client multicard.Client, log *zap.SugaredLogger, ocBaseUrl string, cartService *order.CartServiceImpl, storeId, secret string) (*MultiCardService, error) {
	callback := fmt.Sprintf(ocBaseUrl + "/v1/multicard/webhooks")

	return &MultiCardService{
		&MulticardClientImpl{multicardCli: client, log: log, callbackUrl: callback, cartService: cartService, storeId: storeId, secret: secret},
	}, nil
}

func (m *MulticardClientImpl) CreateCustomer(ctx context.Context, req models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	return models.PaymentSystemCustomer{}, nil
}

func (m *MulticardClientImpl) CreateSubscription(ctx context.Context, req models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	return models.PaymentSystemSubscription{}, nil
}

func (m *MulticardClientImpl) GetSystemWebhookEventRequestByPaymentSystemRequest(r interface{}) (models.WebhookEvent, error) {
	const op = "multicard.GetSystemWebhookEventRequestByPaymentSystemRequest"

	logger := m.log.With(zap.Fields(
		zap.String("op", op),
		zap.Any("request", r),
	))

	req, ok := r.(multicardDto.Webhook)
	if !ok {
		logger.Errorf("casting error")
		return models.WebhookEvent{}, errors.New("casting error")
	}

	resp, err := m.toWebhookEvent(req)
	if err != nil {
		logger.Errorf("failed to get webhook event")
		return models.WebhookEvent{}, err
	}

	return resp, nil
}

func (m *MulticardClientImpl) CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	const op = "multicardService.CreatePaymentOrder"

	logger := m.log.WithOptions(zap.Fields(
		zap.String("op", op),
		zap.Any("paymentOrder", paymentOrder),
	))

	resp, err := m.multicardCli.CreatePaymentInvoice(ctx, dto.CreatePaymentInvoiceRequest{
		Amount:      paymentOrder.Amount,
		InvoiceId:   paymentOrder.CartID,
		ReturnUrl:   paymentOrder.BackURL,
		CallbackUrl: m.callbackUrl,
		Ofds:        make([]dto.Ofd, 0),
	})
	if err != nil {
		logger.Errorf("failed to create payment invoice %s", err.Error())
		return models.PaymentOrder{}, fmt.Errorf("%s: %w", op, err)
	}

	if !resp.Success {
		logger.Errorf("failed to create payment invoice: response status is fail")
		return models.PaymentOrder{}, fmt.Errorf("%s: %w", op, errors.New("response of payment invoice is fail"))
	}

	paymentOrder.PaymentOrderID = resp.Data.Uuid
	paymentOrder.CheckoutURL = resp.Data.CheckoutUrl

	return paymentOrder, nil
}

func (m *MulticardClientImpl) CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error) {
	return models.PaymentEvent{}, nil
}

func (m *MulticardClientImpl) OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error) {
	return models.ApplePaySessionOpenResponse{}, nil
}

func (m *MulticardClientImpl) GetCustomerCards(ctx context.Context, customerID string) ([]models.CustomerCards, error) {
	return nil, nil
}

func (m *MulticardClientImpl) GetPaymentStatusByID(ctx context.Context, paymentID string) (string, error) {
	return "", nil
}

func (m *MulticardClientImpl) RefundPayment(ctx context.Context, paymentOrder models.PaymentOrder, amount int) (models.PaymentOrder, models.RefundResponse, error) {
	const op = "multicardService.RefundPayment"

	log := m.log.With(zap.Fields(
		zap.String("op", op),
		zap.String("payment_order_id", paymentOrder.ExternalID),
	))

	log.Infof("starting refund for payment %s", paymentOrder.ExternalID)

	payment, _, err := m.multicardCli.ReturnFunds(ctx, paymentOrder.MulticardRefundUuid)
	if err != nil {
		return models.PaymentOrder{}, models.RefundResponse{}, err
	}

	refundResp := models.RefundResponse{
		PaymentID: payment.Uuid,
		Status:    payment.Status,
	}

	paymentOrder.PaymentOrderStatus = payment.Status
	paymentOrder.PaymentOrderStatusHistory = append(paymentOrder.PaymentOrderStatusHistory, models.StatusHistory{
		Status: payment.Status,
		Time:   time.Now().UTC(),
	})

	return paymentOrder, refundResp, nil
}

func (m *MulticardClientImpl) SendPaymentOrderToCustomer(ctx context.Context, paymentOrder models.PaymentOrder) error {
	return nil
}

func (m *MulticardClientImpl) CreatePaymentLink(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	return models.PaymentOrder{}, nil
}

func (m *MulticardClientImpl) toWebhookEvent(req multicardDto.Webhook) (models.WebhookEvent, error) {
	_, ok := m.checkSign(req)
	if !ok {
		return models.WebhookEvent{}, fmt.Errorf("signs does not match")
	}

	return models.WebhookEvent{
		Event: models.PAYMENT_APPROVED,
		PaymentEvent: models.PaymentEvent{
			OrderID:             req.InvoiceUuid,
			Status:              models.PAID,
			CapturedAmount:      strconv.Itoa(req.Amount),
			CreatedAt:           req.PaymentTime.Time,
			MulticardRefundUuid: req.Uuid,
			Payer: models.Payer{
				Phone:     req.Phone,
				PanMasked: req.CardPan,
			},
		},
		OrderEvent: models.OrderEvent{
			Status: models.PAID,
		},
	}, nil
}

func (m *MulticardClientImpl) checkSign(req multicardDto.Webhook) (coreModels.Cart, bool) {
	cart, err := m.getCart(req.InvoiceId)
	if err != nil {
		return coreModels.Cart{}, false
	}

	signStr := strings.TrimSpace(m.storeId + cart.ID + strconv.Itoa(req.Amount) + m.secret)

	hasher := md5.New()
	hasher.Write([]byte(signStr))

	predictedHash := hex.EncodeToString(hasher.Sum(nil))

	if predictedHash != req.Sign {
		return coreModels.Cart{}, false
	}

	return cart, true
}

func (m *MulticardClientImpl) getCart(cartId string) (coreModels.Cart, error) {
	ctx := context.Background()

	cart, err := m.cartService.GetKwaakaAdminCartByOrderID(ctx, cartId)
	if err != nil {
		if errors.Is(err, ocErrs.ErrNotFound) {
			cart, err := m.cartService.GetQRMenuCartByID(ctx, cartId)
			if err != nil {
				return coreModels.Cart{}, err
			}
			return cart, nil
		}
		return coreModels.Cart{}, err
	}

	return cart, nil
}
