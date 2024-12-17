package payment

import (
	"context"
	"encoding/json"
	"fmt"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	paymeDto "github.com/kwaaka-team/orders-core/service/payment/payme/dto"
	"github.com/kwaaka-team/orders-core/service/payment/repository"
	"github.com/kwaaka-team/orders-core/service/payment/wooppay/dto"
	"github.com/kwaaka-team/orders-core/service/refund"
	refundModels "github.com/kwaaka-team/orders-core/service/refund/models"
	storeServicePkg "github.com/kwaaka-team/orders-core/service/store"
	storeGroupServicePkg "github.com/kwaaka-team/orders-core/service/storegroup"
	wppService "github.com/kwaaka-team/orders-core/service/whatsapp"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"regexp"
	"time"
)

type Service interface {
	CreateCustomerInDB(ctx context.Context, customer models.PaymentSystemCustomer) error
	CreateCustomerInPaymentSystem(ctx context.Context, paymentSystem string, email string, restaurantID string) (string, error)
	CreateSubscriptionInPaymentSystem(ctx context.Context, paymentSystem, subscriptionID, cardID string, restaurantID string) (string, error)
	CreateSubscriptionInDB(ctx context.Context, subscription models.PaymentSystemSubscription) error
	WebhookEvent(ctx context.Context, paymentSystem string, paymentSystemReq interface{}) (interface{}, error)
	CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	SaveCardDetailsEvent(ctx context.Context, event models.Card) error
	SavePaymentDetailsEvent(ctx context.Context, req models.PaymentOrder) (models.PaymentOrder, error)
	CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error)
	OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error)
	CreatePaymentByApplePayAndSavePaymentDetails(ctx context.Context, paymentRequest models.ApplePayPayment) error
	WoopPayWebhookEvent(ctx context.Context, paymentSystemReq dto.WebhookEventRequest) error
	UpdatePaymentOrderByOrderID(ctx context.Context, cartID string, status string) error
	GetUnpaidPayments(ctx context.Context, minutes int) ([]models.PaymentOrder, error)
	SetNotificationCount(ctx context.Context, id string, count int) error
	CreateCustomerForSavedCardsPayments(ctx context.Context, paymentSystem, restaurantID string, customer models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	GetCustomerCards(ctx context.Context, restaurantID, paymentSystem string, customer models.PaymentSystemCustomer) ([]models.CustomerCards, error)
	GetUnpaidPaymentsByPaymentSystem(ctx context.Context, minutes int, paymentSystem string) ([]models.PaymentOrder, error)
	GetKaspiSaleScoutPaymentStatus(ctx context.Context, paymentID string) (string, error)
	RefundToCustomer(ctx context.Context, orderID, reason string, amount int) (models.PaymentOrder, models.RefundResponse, error, error)
	GetRefund(ctx context.Context, orderID string) (refundModels.Refund, error)
	CreatePaymentLinkForCustomerToPay(ctx context.Context, orderId string) (models.PaymentOrder, error)
}

type ServiceImpl struct {
	paymentSystemFactory *PaymentSystemFactory
	customersRepo        repository.CustomersRepository
	paymentsRepo         repository.PaymentsRepository
	subscriptionRepo     repository.SubscriptionsRepository
	notifyQueue          notifyQueue.SQSInterface
	paymentsQueueUrl     string
	logger               *zap.SugaredLogger
	storeService         storeServicePkg.Service
	storeGroupService    storeGroupServicePkg.Service
	refundRepo           refund.Repository
}

func NewService(paymentSystemFactory *PaymentSystemFactory,
	customersRepo repository.CustomersRepository,
	paymentsRepo repository.PaymentsRepository,
	subscriptionRepo repository.SubscriptionsRepository,
	notifyQueue notifyQueue.SQSInterface, paymentsQueueUrl string, logger *zap.SugaredLogger,
	storeService storeServicePkg.Service,
	storeGroupService storeGroupServicePkg.Service,
	refundRepo refund.Repository,
) (Service, error) {
	if paymentSystemFactory == nil {
		return nil, errors.New("payment system factory is nil")
	}

	return &ServiceImpl{
		paymentSystemFactory: paymentSystemFactory,
		customersRepo:        customersRepo,
		paymentsRepo:         paymentsRepo,
		subscriptionRepo:     subscriptionRepo,
		notifyQueue:          notifyQueue,
		paymentsQueueUrl:     paymentsQueueUrl,
		logger:               logger,
		storeService:         storeService,
		storeGroupService:    storeGroupService,
		refundRepo:           refundRepo,
	}, nil
}

func (s *ServiceImpl) GetKaspiSaleScoutPaymentStatus(ctx context.Context, paymentID string) (string, error) {
	paymentSystem, err := s.paymentSystemFactory.GetPaymentSystem(models.PaymentOrder{
		PaymentSystem: models.KaspiSaleScout,
	}, coreStoreModels.Store{})

	if err != nil {
		return "", err
	}

	return paymentSystem.GetPaymentStatusByID(ctx, paymentID)
}

func (s *ServiceImpl) CreateCustomerInDB(ctx context.Context, customer models.PaymentSystemCustomer) error {
	if err := s.validateCustomer(customer); err != nil {
		return err
	}

	_, err := s.customersRepo.GetCustomerByEmail(ctx, customer.Email)
	if err == nil {
		return errors.Errorf("customer with email %s already exists in db", customer.Email)
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	_, err = s.customersRepo.InsertCustomer(ctx, customer)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) WoopPayWebhookEvent(ctx context.Context, paymentSystemReq dto.WebhookEventRequest) error {
	paymentOrderInDB, err := s.paymentsRepo.GetPaymentOrderByPaymentOrderID(ctx, paymentSystemReq.OperationID)
	if err != nil {
		return err
	}

	store, err := s.storeService.GetByID(ctx, paymentOrderInDB.RestaurantID)
	if err != nil {
		return err
	}

	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(models.PaymentOrder{
		PaymentSystem: models.WOOPPAY,
	}, store)
	if err != nil {
		return err
	}

	webhook, err := paymentService.GetSystemWebhookEventRequestByPaymentSystemRequest(paymentSystemReq)
	if err != nil {
		return err
	}

	_, err = s.SavePaymentDetailsEvent(ctx, models.PaymentOrder{
		PaymentOrderID:       webhook.PaymentEvent.OrderID,
		Payer:                webhook.PaymentEvent.Payer,
		PaymentOrderStatus:   "PAID",
		PaymentStatusHistory: []models.StatusHistory{{Status: webhook.PaymentEvent.Status}},
		CartID:               paymentOrderInDB.CartID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) WebhookEvent(ctx context.Context, paymentSystem string, paymentSystemReq interface{}) (interface{}, error) {
	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(models.PaymentOrder{PaymentSystem: paymentSystem}, coreStoreModels.Store{
		PaymentSystems: []coreStoreModels.PaymentSystem{
			{
				Name:     paymentSystem,
				IsActive: true,
			},
		},
	})
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	s.logger.Infof("Webhook: %s", paymentSystemReq)
	webhook, err := paymentService.GetSystemWebhookEventRequestByPaymentSystemRequest(paymentSystemReq)
	if err != nil {
		s.logger.Error(err)
		if errors.Is(err, wppService.ErrWpp) {
			return nil, err
		}
		return nil, err
	}

	switch paymentSystem {
	case models.IOKA, models.WHATSAPP, models.MultiCard:
		return s.iokaWebhookEvent(ctx, webhook)
	case models.PAYME:
		return s.paymeWebhookEvent(ctx, webhook)

	default:
		return nil, errors.New("unsupported payment system")
	}

}

func (s *ServiceImpl) CreateCustomerForSavedCardsPayments(ctx context.Context, paymentSystem, restaurantID string, customer models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	//if err := s.validateEmail(customer.Email); err != nil {
	//	return models.PaymentSystemCustomer{}, err
	//}

	store, err := s.storeService.GetByID(ctx, restaurantID)
	if err != nil {
		return models.PaymentSystemCustomer{}, err
	}

	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(models.PaymentOrder{
		PaymentSystem: paymentSystem,
	}, store)
	if err != nil {
		return models.PaymentSystemCustomer{}, err
	}

	customer, err = paymentService.CreateCustomer(ctx, customer)
	if err != nil {
		return models.PaymentSystemCustomer{}, err
	}

	return customer, nil
}

func (s *ServiceImpl) CreateCustomerInPaymentSystem(ctx context.Context, paymentSystem string, email string, restaurantID string) (string, error) {
	if err := s.validateEmail(email); err != nil {
		return "", err
	}

	store, err := s.storeService.GetByID(ctx, restaurantID)
	if err != nil {
		return "", err
	}

	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(models.PaymentOrder{
		PaymentSystem: paymentSystem,
	}, store)
	if err != nil {
		return "", err
	}

	customer, err := s.customersRepo.GetCustomerByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("email does not exist")
		}
		return "", err
	}

	customer, err = paymentService.CreateCustomer(ctx, customer)
	if err != nil {
		return "", err
	}

	if err := s.customersRepo.UpdateCustomer(ctx, customer); err != nil {
		return "", err
	}

	return customer.CheckoutURL, nil
}

func (s *ServiceImpl) CreateSubscriptionInPaymentSystem(ctx context.Context, paymentSystem, subscriptionID, cardID, restaurantID string) (string, error) {
	subscription, err := s.subscriptionRepo.GetSubscriptionByID(ctx, subscriptionID)
	if err != nil {
		return "", err
	}
	subscription.Payer.CardID = cardID

	store, err := s.storeService.GetByID(ctx, restaurantID)
	if err != nil {
		return "", err
	}

	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(models.PaymentOrder{
		PaymentSystem: paymentSystem,
	}, store)
	if err != nil {
		return "", err
	}

	res, err := paymentService.CreateSubscription(ctx, subscription)
	if err != nil {
		return "", err
	}

	if err := s.subscriptionRepo.UpdateSubscription(ctx, res); err != nil {
		return "", err
	}

	return res.Schedule.Status, nil
}

func (s *ServiceImpl) CreatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	paymentOrder.PaymentOrderStatus = "NEW"

	switch paymentOrder.PaymentSystem {
	case models.WHATSAPP:
		paymentOrder.PaymentOrderStatusHistory = []models.StatusHistory{
			{
				Status: "NEW",
				Time:   time.Now().UTC(),
			},
			{
				Status: "UNPAID",
				Time:   time.Now().UTC(),
			},
		}
	default:
		paymentOrder.PaymentOrderStatusHistory = []models.StatusHistory{
			{
				Status: "NEW",
				Time:   time.Now().UTC(),
			},
		}
	}

	store, err := s.storeService.GetByID(ctx, paymentOrder.RestaurantID)
	if err != nil {
		return models.PaymentOrder{}, err
	}

	storeGroup, err := s.storeGroupService.GetStoreGroupByID(ctx, store.RestaurantGroupID)
	if err != nil {
		return models.PaymentOrder{}, err
	}

	paymentOrder.RestaurantName = store.Name
	paymentOrder.RestaurantGroupName = storeGroup.Name
	paymentOrder.WhatsappPaymentChatId = store.WhatsappPaymentChatId

	paymentOrder, err = s.paymentsRepo.InsertPaymentOrder(ctx, paymentOrder)
	if err != nil {
		return models.PaymentOrder{}, err
	}

	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(paymentOrder, store)
	if err != nil {
		return models.PaymentOrder{}, err
	}

	res, err := paymentService.CreatePaymentOrder(ctx, paymentOrder)
	if err != nil {
		return models.PaymentOrder{}, err
	}

	if err = s.paymentsRepo.UpdatePaymentOrder(ctx, res); err != nil {
		return models.PaymentOrder{}, err
	}

	return res, nil
}

func (s *ServiceImpl) CreatePaymentByApplePay(ctx context.Context, paymentRequest models.ApplePayPayment) (models.PaymentEvent, error) {
	paymentOrder, err := s.paymentsRepo.GetPaymentOrderByPaymentOrderID(ctx, paymentRequest.OrderID)
	if err != nil {
		return models.PaymentEvent{}, err
	}

	store, err := s.storeService.GetByID(ctx, paymentOrder.RestaurantID)
	if err != nil {
		return models.PaymentEvent{}, err
	}
	paymentOrder.PaymentSystem = paymentRequest.PaymentSystem

	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(paymentOrder, store)
	if err != nil {
		return models.PaymentEvent{}, err
	}

	res, err := paymentService.CreatePaymentByApplePay(ctx, paymentRequest)
	if err != nil {
		return models.PaymentEvent{}, err
	}

	return res, nil
}

func (s *ServiceImpl) OpenApplePaySession(ctx context.Context, request models.ApplePaySessionOpenRequest) (models.ApplePaySessionOpenResponse, error) {
	paymentOrder, err := s.paymentsRepo.GetPaymentOrderByPaymentOrderID(ctx, request.OrderID)
	if err != nil {
		return models.ApplePaySessionOpenResponse{}, err
	}

	store, err := s.storeService.GetByID(ctx, paymentOrder.RestaurantID)
	if err != nil {
		return models.ApplePaySessionOpenResponse{}, err
	}

	paymentOrder.PaymentSystem = request.PaymentSystem

	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(paymentOrder, store)
	if err != nil {
		return models.ApplePaySessionOpenResponse{}, err
	}

	res, err := paymentService.OpenApplePaySession(ctx, request)
	if err != nil {
		return models.ApplePaySessionOpenResponse{}, err
	}

	return res, nil
}

func (s *ServiceImpl) CreateSubscriptionInDB(ctx context.Context, subscription models.PaymentSystemSubscription) error {

	customer, err := s.customersRepo.GetCustomerByEmail(ctx, subscription.Payer.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("email does not exist")
		}
		return err
	}

	subscription.Payer.CustomerID = customer.PaymentSystemCustomerID

	_, err = s.subscriptionRepo.InsertSubscription(ctx, subscription)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) SaveCardDetailsEvent(ctx context.Context, event models.Card) error {

	customer, err := s.customersRepo.GetCustomerByPaymentSystemCustomerID(ctx, event.CustomerID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New("customer does not exist")
		}
		return err
	}

	customer = s.saveCardChangesInCustomer(customer, event)

	if err := s.customersRepo.UpdateCustomer(ctx, customer); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) SavePaymentDetailsEvent(ctx context.Context, req models.PaymentOrder) (models.PaymentOrder, error) {
	paymentOrder, err := s.paymentsRepo.GetPaymentOrderByPaymentOrderID(ctx, req.PaymentOrderID)
	if err != nil {
		return paymentOrder, err
	}

	oldPaymentOrderStatus := paymentOrder.PaymentOrderStatus

	if oldPaymentOrderStatus == models.PAID {
		return models.PaymentOrder{}, nil
	}

	paymentOrder = s.savePaymentOrderChanges(req, paymentOrder)

	if err := s.paymentsRepo.UpdatePaymentOrder(ctx, paymentOrder); err != nil {
		return paymentOrder, err
	}

	//add create order logic using query
	paymentOrderJson, err := json.Marshal(paymentOrder)
	if err != nil {
		return paymentOrder, err
	}
	s.logger.Infof("queue url: %s", s.paymentsQueueUrl)
	s.logger.Infof("queue message body: %s", string(paymentOrderJson))

	if paymentOrder.PaymentOrderStatus != models.PAID {
		return paymentOrder, nil
	}

	if err := s.notifyQueue.SendSQSMessageToFIFO(ctx, s.paymentsQueueUrl, string(paymentOrderJson), req.PaymentOrderID); err != nil {
		return paymentOrder, err
	}

	return paymentOrder, nil
}

func (s *ServiceImpl) CreatePaymentByApplePayAndSavePaymentDetails(ctx context.Context, paymentRequest models.ApplePayPayment) error {

	res, err := s.CreatePaymentByApplePay(ctx, paymentRequest)
	if err != nil {
		return err
	}

	_, err = s.SavePaymentDetailsEvent(ctx, models.PaymentOrder{
		PaymentOrderID:       res.OrderID,
		PaymentID:            res.ID,
		Payer:                res.Payer,
		Acquirer:             res.Acquirer,
		Action:               res.Action,
		Error:                res.Error,
		PaymentOrderStatus:   res.Status,
		PaymentStatusHistory: []models.StatusHistory{{Status: res.Status}},
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) GetUnpaidPayments(ctx context.Context, minutes int) ([]models.PaymentOrder, error) {
	return s.paymentsRepo.GetUnpaidPayments(ctx, minutes)
}

func (s *ServiceImpl) GetUnpaidPaymentsByPaymentSystem(ctx context.Context, minutes int, paymentSystem string) ([]models.PaymentOrder, error) {
	return s.paymentsRepo.GetUnpaidPaymentsByPaymentSystem(ctx, minutes, paymentSystem)
}

func (s *ServiceImpl) SetNotificationCount(ctx context.Context, id string, count int) error {
	return s.paymentsRepo.SetNotificationCount(ctx, id, count)
}

func (s *ServiceImpl) validateCustomer(customer models.PaymentSystemCustomer) error {
	if err := s.validateEmail(customer.Email); err != nil {
		return err
	}

	if err := s.validatePhone(customer.Phone); err != nil {
		return err
	}

	if err := s.validateChannel(customer.Channel); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) validateEmail(email string) error {
	emailReg, err := regexp.Compile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if err != nil {
		return err
	}

	if !emailReg.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func (s *ServiceImpl) validatePhone(phone string) error {
	phoneReg, err := regexp.Compile(`^(\+77)\d{9}?$`)
	if err != nil {
		return err
	}
	if !phoneReg.MatchString(phone) {
		return errors.New("invalid phone number format")
	}

	return nil
}

func (s *ServiceImpl) validateChannel(channel string) error {
	if channel == "" {
		return errors.New("invalid channel")
	}

	return nil
}

func (s *ServiceImpl) saveCardChangesInCustomer(customer models.PaymentSystemCustomer, event models.Card) models.PaymentSystemCustomer {
	for i, card := range customer.Cards {
		if card.ID != event.ID {
			continue
		}

		customer.Cards[i] = models.Card{
			ID:            event.ID,
			Status:        event.Status,
			CreatedAt:     event.CreatedAt,
			PanMasked:     event.PanMasked,
			ExpiryDate:    event.ExpiryDate,
			PaymentSystem: event.PaymentSystem,
			Emitter:       event.Emitter,
			CvcRequired:   event.CvcRequired,
			CustomerID:    event.CustomerID,
		}
		customer.Status = event.CustomerStatus

		return customer
	}

	customer.Cards = append(customer.Cards, models.Card{
		ID:            event.ID,
		Status:        event.Status,
		CreatedAt:     event.CreatedAt,
		PanMasked:     event.PanMasked,
		ExpiryDate:    event.ExpiryDate,
		PaymentSystem: event.PaymentSystem,
		Emitter:       event.Emitter,
		CvcRequired:   event.CvcRequired,
		CustomerID:    event.CustomerID,
	})
	customer.Status = event.CustomerStatus

	return customer
}

func (s *ServiceImpl) savePaymentOrderChanges(req models.PaymentOrder, paymentOrder models.PaymentOrder) models.PaymentOrder {
	paymentOrder.PaymentID = req.PaymentID
	paymentOrder.Payer = req.Payer
	paymentOrder.Acquirer = req.Acquirer
	paymentOrder.Action = req.Action
	paymentOrder.Error = req.Error
	paymentOrder.PaymentStatusHistory = append(paymentOrder.PaymentStatusHistory, models.StatusHistory{
		Status: req.PaymentStatusHistory[0].Status,
		Time:   time.Now().UTC(),
	})
	paymentOrder.PaymentOrderStatus = req.PaymentOrderStatus
	paymentOrder.PaymentOrderStatusHistory = append(paymentOrder.PaymentOrderStatusHistory, models.StatusHistory{
		Status: req.PaymentOrderStatus,
		Time:   time.Now().UTC(),
	})
	paymentOrder.RefundReason = req.RefundReason
	paymentOrder.RefundAuthor = req.RefundAuthor
	paymentOrder.RefundAmount = req.RefundAmount
	paymentOrder.MulticardRefundUuid = req.MulticardRefundUuid

	return paymentOrder
}

func (s *ServiceImpl) paymeCreateTransaction(ctx context.Context, webhook models.WebhookEvent) (interface{}, error) {
	paymentOrder, err := s.paymentsRepo.GetPaymentOrderByPaymentOrderID(ctx, webhook.PaymentEvent.OrderID)
	if err != nil {
		s.logger.Errorf("payme create transaction error: %s", err)
		return nil, err
	}

	return paymeDto.WebhookResultResponse{
		Result: paymeDto.WebhookResult{
			Transaction: paymentOrder.ExternalID,
			CreateTime:  time.Now().Unix(),
			State:       1,
		},
	}, nil
}

func (s *ServiceImpl) paymeCheckPerformTransaction() (interface{}, error) {
	return paymeDto.WebhookResultResponse{
		Result: paymeDto.WebhookResult{
			Allow: true,
		},
	}, nil
}

func (s *ServiceImpl) paymePerformTransaction(ctx context.Context, webhook models.WebhookEvent) (interface{}, error) {
	paymentOrder, err := s.SavePaymentDetailsEvent(ctx, models.PaymentOrder{
		PaymentOrderID:       webhook.PaymentEvent.OrderID,
		PaymentOrderStatus:   models.PAID,
		PaymentStatusHistory: []models.StatusHistory{{Status: models.PAID}},
	})
	if err != nil {
		s.logger.Errorf("payme perform transaction error: %s", err)
		return nil, err
	}
	return paymeDto.WebhookResultResponse{
		Result: paymeDto.WebhookResult{
			Transaction: paymentOrder.ExternalID,
			PerformTime: time.Now().Unix(),
			State:       2,
		},
	}, nil
}

func (s *ServiceImpl) paymeWebhookEvent(ctx context.Context, webhook models.WebhookEvent) (interface{}, error) {
	switch webhook.Event {
	case models.PAYME_CREATE_TRANSACTION_METHOD:
		return s.paymeCreateTransaction(ctx, webhook)

	case models.PAYME_CHECK_PERFORM_TRANSACTION_METHOD:
		return s.paymeCheckPerformTransaction()

	case models.PAYME_PERFORM_TRANSACTION_METHOD:
		return s.paymePerformTransaction(ctx, webhook)

	default:
		return nil, errors.New("unsupported webhook event type")
	}
}

func (s *ServiceImpl) iokaWebhookEvent(ctx context.Context, webhook models.WebhookEvent) (interface{}, error) {
	switch webhook.Event {
	case models.CARD_APPROVED, models.CARD_DECLINED:
		err := s.SaveCardDetailsEvent(ctx, models.Card{
			ID:             webhook.CardEvent.ID,
			Status:         webhook.CardEvent.Status,
			CreatedAt:      webhook.CardEvent.CreatedAt,
			PanMasked:      webhook.CardEvent.PanMasked,
			ExpiryDate:     webhook.CardEvent.ExpiryDate,
			PaymentSystem:  webhook.CardEvent.PaymentSystem,
			Emitter:        webhook.CardEvent.Emitter,
			CvcRequired:    webhook.CardEvent.CvcRequired,
			CustomerID:     webhook.CardEvent.CustomerID,
			CustomerStatus: webhook.CustomerEvent.Status,
		})
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		return nil, nil

	case models.PAYMENT_DECLINED, models.PAYMENT_APPROVED, models.PAYMENT_CAPTURED, models.PAYMENT_CANCELED, models.PAYMENT_ACTION_REQUIRED,
		models.ORDER_EXPIRED, models.CAPTURE_DECLINED, models.CANCEL_DECLINED:
		_, err := s.SavePaymentDetailsEvent(ctx, models.PaymentOrder{
			PaymentOrderID:       webhook.PaymentEvent.OrderID,
			PaymentID:            webhook.PaymentEvent.ID,
			Payer:                webhook.PaymentEvent.Payer,
			Acquirer:             webhook.PaymentEvent.Acquirer,
			Action:               webhook.PaymentEvent.Action,
			Error:                webhook.PaymentEvent.Error,
			PaymentOrderStatus:   webhook.OrderEvent.Status,
			MulticardRefundUuid:  webhook.PaymentEvent.MulticardRefundUuid,
			PaymentStatusHistory: []models.StatusHistory{{Status: webhook.PaymentEvent.Status}},
		})
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		return nil, nil

	case models.REFUND_APPROVED, models.REFUND_DECLINED:
		_, err := s.SavePaymentDetailsEvent(ctx, models.PaymentOrder{
			PaymentOrderID:       webhook.PaymentEvent.OrderID,
			PaymentID:            webhook.PaymentEvent.ID,
			Payer:                webhook.PaymentEvent.Payer,
			Acquirer:             webhook.PaymentEvent.Acquirer,
			Action:               webhook.PaymentEvent.Action,
			Error:                webhook.PaymentEvent.Error,
			PaymentOrderStatus:   "REFUND_" + webhook.RefundEvent.Status,
			PaymentStatusHistory: []models.StatusHistory{{Status: "REFUND_" + webhook.RefundEvent.Status}},
			RefundAmount:         webhook.RefundEvent.Amount,
			RefundAuthor:         webhook.RefundEvent.Author,
			RefundReason:         webhook.RefundEvent.Reason,
		})
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		return nil, nil
	default:
		return nil, errors.New("unsupported webhook event type")
	}
}

func (s *ServiceImpl) UpdatePaymentOrderByOrderID(ctx context.Context, cartID string, status string) error {

	order, err := s.paymentsRepo.GetPaymentOrderByOrderID(ctx, cartID)
	if err != nil {
		return err
	}

	if order.PaymentOrderStatus == status {
		s.logger.Infof("status to be updated is the same as current status in payment order")
		return nil
	}

	order.PaymentOrderStatus = status
	order.PaymentOrderStatusHistory = append(order.PaymentStatusHistory, models.StatusHistory{
		Status: status,
		Time:   time.Now().UTC(),
	})

	if err := s.paymentsRepo.UpdatePaymentOrderStatusHistory(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) GetCustomerCards(ctx context.Context, restaurantID, paymentSystem string, customer models.PaymentSystemCustomer) ([]models.CustomerCards, error) {
	store, err := s.storeService.GetByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(models.PaymentOrder{
		PaymentSystem: paymentSystem,
	}, store)
	if err != nil {
		return nil, err
	}

	cards, err := paymentService.GetCustomerCards(ctx, customer.PaymentSystemCustomerID)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (s *ServiceImpl) RefundToCustomer(ctx context.Context, orderID, reason string, amount int) (models.PaymentOrder, models.RefundResponse, error, error) {

	paymentOrder, err := s.paymentsRepo.GetPaymentOrderByOrderID(ctx, orderID)
	if err != nil {
		s.logger.Errorf("get payment order by order id: %s", err)
		return models.PaymentOrder{}, models.RefundResponse{}, fmt.Errorf("payment/service - fn RefundCustomerToCustomer - fn GetPaymentOrderByOrderID: get payment order by order id error: order id %v: %w", orderID, err), nil
	}

	switch {
	case amount < 100 && paymentOrder.PaymentSystem == models.IOKA:
		return models.PaymentOrder{}, models.RefundResponse{}, nil, fmt.Errorf("payment/service - fn RefundCustomerToCustomer: refund amount is less than 100: amount: %d", paymentOrder.Amount)
	case amount < 1 && paymentOrder.PaymentSystem == models.KaspiSaleScout:
		return models.PaymentOrder{}, models.RefundResponse{}, nil, fmt.Errorf("payment/service - fn RefundCustomerToCustomer: refund amount is less than 1: amount: %d", paymentOrder.Amount)
	case amount > paymentOrder.Amount/100:
		return models.PaymentOrder{}, models.RefundResponse{}, nil, fmt.Errorf("payment/service - fn RefundCustomerToCustomer: refund amount exceeds the purchase amount: %d", paymentOrder.Amount)
	}

	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(paymentOrder, coreStoreModels.Store{})
	if err != nil {
		s.logger.Errorf("get payment system error: %s", err)
		return models.PaymentOrder{}, models.RefundResponse{}, fmt.Errorf("payment/service - fn RefundCustomerToCustomer - fn GetPaymentSystem: get payment system error: payment system: %s, error: %w", paymentOrder.PaymentSystem, err), nil
	}

	paymentOrder.RefundReason = reason

	paymentOrder, refundResponse, err := paymentService.RefundPayment(ctx, paymentOrder, amount)
	if err != nil {
		s.logger.Errorf("refund payment to customer error: %s", err)
		return models.PaymentOrder{}, models.RefundResponse{}, fmt.Errorf("payment/service - fn RefundCustomerToCustomer - fn RefundPayment: error: order id: %s, amount: %d, api error: %w",
			paymentOrder.OrderID, paymentOrder.Amount, err), nil
	}

	if err = s.paymentsRepo.UpdatePaymentOrder(ctx, paymentOrder); err != nil {
		return models.PaymentOrder{}, models.RefundResponse{}, fmt.Errorf("payment/service - fn RefundCustomerToCustomer - fn UpdatePaymentOrder: update payment order error: %w", err), nil
	}

	if err = s.refundRepo.InsertRefundInfo(ctx, refundModels.Refund{
		ID:            paymentOrder.ExternalID,
		Amount:        amount,
		Reason:        paymentOrder.RefundReason,
		OrderID:       paymentOrder.OrderID,
		PaymentID:     refundResponse.PaymentID,
		PaymentSystem: paymentOrder.PaymentSystem,
	}); err != nil {
		return models.PaymentOrder{}, models.RefundResponse{}, fmt.Errorf("payment/service - fn RefundCustomerToCustomer - fn InsertRefundInfo: inserting refund to db error: %w", err), nil
	}

	return paymentOrder, refundResponse, nil, nil
}

func (s *ServiceImpl) GetRefund(ctx context.Context, orderID string) (refundModels.Refund, error) {
	res, err := s.refundRepo.GetRefund(ctx, orderID)
	if err != nil {
		return refundModels.Refund{}, err
	}
	return res, nil
}

func (s *ServiceImpl) CreatePaymentLinkForCustomerToPay(ctx context.Context, orderId string) (models.PaymentOrder, error) {
	paymentOrder, err := s.paymentsRepo.GetPaymentOrderByOrderID(ctx, orderId)
	if err != nil {
		return models.PaymentOrder{}, fmt.Errorf("payment/service - fn CreatePaymentLinkForCustomerToPay - fn GetPaymentOrderByOrderID: %w", err)
	}

	paymentService, err := s.paymentSystemFactory.GetPaymentSystem(paymentOrder, coreStoreModels.Store{})
	if err != nil {
		s.logger.Errorf("get payment system error: %s", err)
		return models.PaymentOrder{}, fmt.Errorf("payment/service - fn CreatePaymentLinkForCustomerToPay - fn GetPaymentSystem: error: %w", err)
	}

	paymentOrderRes, err := paymentService.CreatePaymentLink(ctx, paymentOrder)
	if err != nil {
		return models.PaymentOrder{}, err
	}

	if err = s.paymentsRepo.UpdatePaymentOrder(ctx, paymentOrderRes); err != nil {
		return models.PaymentOrder{}, err
	}

	return paymentOrderRes, nil
}
