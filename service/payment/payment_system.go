package payment

import (
	"context"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/order"
	"github.com/kwaaka-team/orders-core/service/payment/callcenter"
	"github.com/kwaaka-team/orders-core/service/payment/cash"
	"github.com/kwaaka-team/orders-core/service/payment/ioka"
	"github.com/kwaaka-team/orders-core/service/payment/kaspi"
	"github.com/kwaaka-team/orders-core/service/payment/kaspi_manual"
	"github.com/kwaaka-team/orders-core/service/payment/kaspi_salescout"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"github.com/kwaaka-team/orders-core/service/payment/multicard"
	"github.com/kwaaka-team/orders-core/service/payment/payme"
	"github.com/kwaaka-team/orders-core/service/payment/whatsapp"
	"github.com/kwaaka-team/orders-core/service/payment/wooppay"
	wppService "github.com/kwaaka-team/orders-core/service/whatsapp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
)

var paymentSystemNotFoundError = errors.New("payment system not found")

type PaymentSystem interface {
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

type PaymentSystemFactory struct {
	ioka                  *ioka.IokaService
	kaspi                 *kaspi.KaspiService
	payme                 *payme.PaymeService
	woopPayBaseUrl        string
	woopPayResultUrl      string
	callCenter            *callcenter.CallCenterService
	whatsappService       wppService.Service
	kaspiSaleScoutService *kaspi_salescout.KaspiSaleScoutService
	cashService           *cash.CashService
	multicardService      *multicard.MultiCardService
}

func NewFactory(iokaBaseUrl, iokaApiKey, paymeBaseUrl, paymeApiKey, woopPayBaseUrl, woopPayResultUrl string,
	wppService wppService.Service, kaspiSaleScoutBaseUrl, kaspiSaleScoutToken, kaspiSaleScoutMerchantID string,
	logger *zap.SugaredLogger, ocBaseUrl string, cartService *order.CartServiceImpl) (*PaymentSystemFactory, error) {
	iokaService, err := ioka.NewIokaService(iokaBaseUrl, iokaApiKey)
	if err != nil {
		return nil, err
	}

	kaspiService, err := kaspi.NewKaspiService()
	if err != nil {
		return nil, err
	}

	paymeService, err := payme.NewPaymeService(paymeBaseUrl, paymeApiKey)
	if err != nil {
		return nil, err
	}

	callcenterService, err := callcenter.NewCallCenterService()
	if err != nil {
		return nil, err
	}

	kaspiSaleScoutService, err := kaspi_salescout.NewKaspiSaleScoutService(kaspiSaleScoutBaseUrl, kaspiSaleScoutToken, kaspiSaleScoutMerchantID)
	if err != nil {
		return nil, err
	}

	cashService, err := cash.NewCashService()
	if err != nil {
		return nil, err
	}

	return &PaymentSystemFactory{
		ioka:                  iokaService,
		kaspi:                 kaspiService,
		payme:                 paymeService,
		woopPayBaseUrl:        woopPayBaseUrl,
		woopPayResultUrl:      woopPayResultUrl,
		callCenter:            callcenterService,
		whatsappService:       wppService,
		kaspiSaleScoutService: kaspiSaleScoutService,
		cashService:           cashService,
	}, nil
}

func (f *PaymentSystemFactory) GetPaymentSystem(paymentOrder models.PaymentOrder, store coreStoreModels.Store) (PaymentSystem, error) {
	switch paymentOrder.PaymentSystem {
	case models.IOKA:
		return f.ioka, nil
	case models.KASPI:
		return f.kaspi, nil
	case models.PAYME:
		return f.payme, nil
	case models.WHATSAPP:
		return f.getWhatsappService(store, f.whatsappService)
	case models.WOOPPAY:
		return f.getWoopPayService(store)
	case models.CallCenter:
		return f.callCenter, nil
	case models.KaspiManual:
		return f.getKaspiManualService(store)
	case models.KaspiSaleScout:
		return f.kaspiSaleScoutService, nil
	case models.Cash:
		return f.cashService, nil
	case models.MultiCard:
		return f.multicardService, nil
	}
	return nil, errors.Wrapf(paymentSystemNotFoundError, "delivery service %s not found", paymentOrder.PaymentSystem)
}

func (f *PaymentSystemFactory) getWoopPayService(store coreStoreModels.Store) (*wooppay.WoopPayService, error) {
	storePaymentService, err := store.GetStorePaymentService(models.WOOPPAY)
	if err != nil {
		return nil, err
	}

	return wooppay.NewWoopPayService(f.woopPayBaseUrl, f.woopPayResultUrl, storePaymentService.Username, storePaymentService.Password, storePaymentService.MerchantService)
}

func (f *PaymentSystemFactory) getKaspiManualService(store coreStoreModels.Store) (PaymentSystem, error) {
	if len(store.PaymentSystems) == 0 {
		return nil, errors.New("kaspi_manual is not exist in store")
	}

	for _, ps := range store.PaymentSystems {
		if strings.ToLower(ps.Name) != "kaspi" {
			continue
		}
		if !ps.IsActive {
			return nil, errors.New("kaspi_manual is inactive in store")
		}

		return kaspi_manual.NewKaspiManualService(ps.PaymentURL)
	}

	return nil, errors.New("kaspi_manual is not exist in store")
}

func (f *PaymentSystemFactory) getWhatsappService(store coreStoreModels.Store, wpp wppService.Service) (PaymentSystem, error) {
	if len(store.PaymentSystems) == 0 {
		return nil, errors.New("whatsapp is not exist in store")
	}

	for _, ps := range store.PaymentSystems {
		if strings.ToLower(ps.Name) != "whatsapp" {
			continue
		}
		if !ps.IsActive {
			return nil, errors.New("whatsapp is inactive in store")
		}

		return whatsapp.NewWhatsappService(wpp, ps.PaymentURL)
	}

	return whatsapp.NewWhatsappService(wpp, "")
}
