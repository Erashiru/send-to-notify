package kwaaka_3pl

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	models3 "github.com/kwaaka-team/orders-core/core/menu/models"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/menu"
	"github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/service/kwaaka_3pl/models"
	"github.com/kwaaka-team/orders-core/service/order"
	"github.com/kwaaka-team/orders-core/service/order/delivery"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"time"
)

type Service interface {
	Create3plOrder(ctx context.Context, req models.CreateDeliveryRequest) error
	Cancel3plOrder(ctx context.Context, orderID string, orderInfo models2.OrderInfoForTelegramMsg) error
	GetDeliveryInfoByOrderId(ctx context.Context, orderId string) (models2.GetDeliveryInfoResp, error)
	SetKwaaka3plDispatcher(ctx context.Context, req models2.SetKwaaka3plDispatcherRequest) error
	GetOrdersByCustomerPhone(ctx context.Context, query models2.GetOrdersByCustomerPhoneRequest) (models2.GetOrdersByCustomerPhoneResponse, error)
	GetCustomerByDeliveryId(ctx context.Context, deliveryId string) (models2.Customer, string, error)
	BulkCreate3plOrder(ctx context.Context, orders []models2.Order, byAdmin bool) error
	ListPotentialProviders(ctx context.Context, req models.ListProvidersRequest) ([]models.ProviderResponse, error)
	GetDeliveryPrice(ctx context.Context, deliveryOrderID string) (float64, error)
	GetOrderForTelegramByDeliveryOrderId(ctx context.Context, deliveryId string) (models2.OrderInfoForTelegramMsg, error)
	GetDeliveryDispatcherPrices(ctx context.Context, deliveryIDs []string) (models.GetDeliveryDispatcherPricesResponse, error)
	Instant3plOrder(ctx context.Context, req models2.Order) error
	CancelCourierSearch(ctx context.Context, deliveryOrderID string) error
	MapIikoStatusTo3plStatus(ctx context.Context, iikoStatus, customerPhoneNumber, storeID string) error
	Save3plHistory(ctx context.Context, deliveryOrderId string, newDeliveryAddress models2.DeliveryAddress, newCustomer models2.Customer) error
	GetOrderByOrderID(ctx context.Context, orderID string) (models2.Order, error)
	Get3plDeliveryInfo(ctx context.Context, deliveries []string) ([]models2.GetDeliveryInfoResp, error)
	NoDispatcherMessage(ctx context.Context) error
	ActualizeDeliveryInfoByDeliveryIDs(ctx context.Context, deliveryIDs []string) error
	PerformerLookupMoreThan15Minute(ctx context.Context) error
	InsertChangeDeliveryHistory(ctx context.Context, deliveryID string, history models2.ChangesHistory) error
	// Todo temporary methods
	GetDeliveryInfoForReport(ctx context.Context, orderId string) (models2.GetDeliveryInfoResp, error)
	GetDeliveryPriceFromPQ(ctx context.Context, deliveryOrderID string) (float64, error)
	GetDeliveryStatus(ctx context.Context, deliveryID string) (string, error)
	GetDeliveryFromPostgres(ctx context.Context, deliveryID string) (models2.DeliveryOrderPQ, error)
	GetDeliveryProdiver(ctx context.Context, deliveryID string) (string, error)
}

type ServiceImpl struct {
	restyClient        *resty.Client
	repository         order.Repository
	storeService       store.Service
	sqsCli             que.SQSInterface
	queueUrl           string
	logger             *zap.SugaredLogger
	telegramService    order.TelegramService
	menuClient         menu.Client
	deliveryRepository delivery.Repository
}

func NewKwaaka3plService(sqsCli que.SQSInterface, queueUrl string, repository order.Repository, storeService store.Service, baseUrl, authToken string, logger *zap.SugaredLogger, telegram order.TelegramService, menuCli menu.Client, deliveryRepository delivery.Repository) (*ServiceImpl, error) {
	if baseUrl == "" {
		return nil, errors.New("base URL could not be empty")
	}
	if repository == nil {
		return nil, errors.New("repository could not be nil")
	}
	if sqsCli == nil {
		return nil, errors.New("sqsCli is nil")
	}
	if queueUrl == "" {
		return nil, errors.New("queueName is nil")
	}
	if deliveryRepository == nil {
		return nil, errors.New("delivery repository could not be nil")
	}

	client := resty.New().
		SetBaseURL(baseUrl).SetHeader("authorization", authToken)

	return &ServiceImpl{
		restyClient:        client,
		repository:         repository,
		storeService:       storeService,
		sqsCli:             sqsCli,
		queueUrl:           queueUrl,
		logger:             logger,
		telegramService:    telegram,
		menuClient:         menuCli,
		deliveryRepository: deliveryRepository,
	}, nil
}

func (cl *ServiceImpl) Create3plOrder(ctx context.Context, req models.CreateDeliveryRequest) error {

	message, err := json.Marshal(req)
	if err != nil {
		return err
	}

	cl.logger.Infof("queue message body: %s", string(message))

	if err = cl.sqsCli.SendSQSMessageToFIFO(ctx, cl.queueUrl, string(message), req.ID); err != nil {
		cl.logger.Errorf("couldn't create 3pl Order: %s", err.Error())
		return err
	}

	return nil
}

func (cl *ServiceImpl) cancel3plOrder(ctx context.Context, deliveryOrderID string, orderInfo models2.OrderInfoForTelegramMsg) error {
	path := "/v1/delivery/cancel"

	var errResponse models.ErrorResponse

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetQueryParam("delivery_id", deliveryOrderID).
		SetBody(orderInfo).
		Put(path)
	if err != nil {
		return fmt.Errorf("cancel order error: %v", err)
	}

	log.Info().Msgf("request 3pl service CancelCourierSearch %+v \n", resp.Request.Body)
	log.Info().Msgf("response 3pl service CancelCourierSearch %+v", string(resp.Body()))

	if resp.IsError() {
		return fmt.Errorf("cancel order response error: %v", errResponse.Message)
	}
	return nil
}

func (s *ServiceImpl) Cancel3plOrder(ctx context.Context, orderID string, orderInfo models2.OrderInfoForTelegramMsg) error {

	order, err := s.repository.FindOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	delivery, err := s.deliveryRepository.GetDeliveryByDeliveryID(ctx, order.DeliveryOrderID)
	if err != nil {
		return err
	}

	store, err := s.storeService.GetByID(ctx, order.RestaurantID)
	if err != nil {
		return err
	}

	if delivery.CancelState == models.CancelUnavailable {
		err := s.telegramService.SendMessageToQueue(telegram.CancelDeliveryFromDispatcherPage, order, store, "", "", "", models3.Product{})
		if err != nil {
			log.Err(err).Msgf("error: construct and send cancel delivery from dispatcher page for delivery id: %s", delivery.Id)
		}
		return nil
	}

	return s.cancel3plOrder(ctx, order.DeliveryOrderID, orderInfo)
}

func (s *ServiceImpl) cancelCourierSearch(ctx context.Context, deliveryOrderID string, orderInfo models2.OrderInfoForTelegramMsg) error {

	path := "/v1/delivery/cancel"

	var errResponse models.ErrorResponse

	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetQueryParam("delivery_id", deliveryOrderID).
		SetBody(orderInfo).
		Put(path)
	if err != nil {
		log.Error().Msgf("kwaaka_3pl/service - fn cancelCourierSearch - fn s.restyClient.R(): request error: delivery order id: %s", deliveryOrderID)
		return err
	}

	log.Info().Msgf("request 3pl service CancelCourierSearch: resp.Request.Body %+v \n", resp.Request.Body)
	log.Info().Msgf("response 3pl service CancelCourierSearch: string(resp.Body()) %+v", string(resp.Body()))

	if resp.IsError() {
		log.Error().Msgf("kwaaka_3pl/service - fn cancelCourierSearch - fn resp.IsError(): response error: %v", resp.IsError())
		return fmt.Errorf("cancel 3pl order response error: %s, %s", resp.Error(), errResponse.Message)
	}

	log.Info().Msgf("success while cancelling 3pl order - cancel_3pl/service - fn cancelCourierSearch: cancel 3pl order: %s", deliveryOrderID)

	return nil
}

func (s *ServiceImpl) CancelCourierSearch(ctx context.Context, deliveryOrderID string) error {

	log.Info().Msgf("kwaaka_3pl/service - fn CancelCourierSearch: start CancelCourierSearch: %s", deliveryOrderID)

	ord, err := s.repository.GetOrderBy3plDeliveryID(ctx, deliveryOrderID)
	if err != nil {
		s.logger.Errorf("error while get order by 3pl delivery id: %s error: %v", deliveryOrderID, err)
		return err
	}

	log.Info().Msgf("kwaaka_3pl/service - fn CancelCourierSearch - fn GetOrderBy3plDeliveryID: success receiving order: delivery order id: %s", deliveryOrderID)

	st, err := s.storeService.GetByID(ctx, ord.RestaurantID)
	if err != nil {
		s.logger.Errorf("error while getting store by restaurant id: %s error: %v", ord.StoreID, err)
		return err
	}

	log.Info().Msgf("kwaaka_3pl/service - fn CancelCourierSearch - fn GetByID: success receiving store: restaurant id: %s", st.ID)

	orderInfo := models2.OrderInfoForTelegramMsg{
		RestaurantName:        ord.RestaurantName,
		RestaurantAddress:     st.Address.City + " " + st.Address.Street,
		RestaurantPhoneNumber: st.StorePhoneNumber,
		OrderId:               ord.ID,
		Id3plOrder:            deliveryOrderID,
		CustomerName:          ord.Customer.Name,
		CustomerPhoneNumber:   ord.Customer.PhoneNumber,
		CustomerAddress:       ord.DeliveryAddress.Label,
		DeliveryService:       ord.DeliveryDispatcher,
	}
	err = s.cancelCourierSearch(ctx, deliveryOrderID, orderInfo)
	if err != nil {
		s.logger.Errorf("fn cancelCourierSearch - error while cancel 3pl order: %s error:  %v", deliveryOrderID, err)
		return err
	}

	log.Info().Msgf("kwaaka_3pl/service - fn CancelCourierSearch: success func CancelCourierSearch")

	return nil
}

func (s *ServiceImpl) Save3plHistory(ctx context.Context, deliveryOrderId string, newDeliveryAddress models2.DeliveryAddress, newCustomer models2.Customer) error {
	ord, err := s.repository.GetOrderBy3plDeliveryID(ctx, deliveryOrderId)
	if err != nil {
		s.logger.Errorf("error while get order by 3pl delivery id: %s error: %v", deliveryOrderId, err)
		return err
	}

	err = s.repository.Save3plDeliveryHistoryAndSetEmptyDispatcherService(ctx, ord.ID, models2.History3plDelivery{
		DeliveryOrderID:            ord.DeliveryOrderID,
		DeliveryDispatcher:         ord.DeliveryDispatcher,
		FullDeliveryPrice:          ord.FullDeliveryPrice,
		RestaurantPayDeliveryPrice: ord.RestaurantPayDeliveryPrice,
		KwaakaChargedDeliveryPrice: ord.KwaakaChargedDeliveryPrice,
		DeliveryAddress:            ord.DeliveryAddress,
		Customer:                   ord.Customer,
	}, newDeliveryAddress, newCustomer)
	if err != nil {
		s.logger.Errorf("error while set empty dispatcher service for order id:%s error:%v", ord.ID, err)
		return err
	}

	log.Info().Msgf("kwaaka_3pl/service - fn Save3plHistory: success func Save3plHistory")

	return nil
}

func (s *ServiceImpl) MapIikoStatusTo3plStatus(ctx context.Context, iikoStatus, customerPhoneNumber, storeID string) error {

	s.logger.Errorf("kwaaka_3pl/service - fn MapIikoStatusTo3plStatus: map iiko status to 3pl status: %s", iikoStatus)

	path := "/v1/delivery/self-delivery-status"

	var (
		errResponse models.ErrorResponse
		req         = models2.ActualizeSelfDeliveryStatusRequest{
			IikoStatus:          iikoStatus,
			CustomerPhoneNumber: customerPhoneNumber,
			StoreID:             storeID,
		}
	)

	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(req).
		Post(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("request actualize self delivery iiko status %+v \n", resp.Request.Body)
	log.Info().Msgf("response actualize self delivery iiko status %+v", string(resp.Body()))

	if resp.IsError() {
		return fmt.Errorf("actualize iiko self delivery status response error: %s, %s", resp.Error(), errResponse.Message)
	}

	s.logger.Errorf("success while actualizing self delivery status - cancel_3pl/service - fn cancelCourierSearch: cancel 3pl order")

	return nil
}

func (s *ServiceImpl) GetDeliveryInfoByOrderId(ctx context.Context, orderId string) (models2.GetDeliveryInfoResp, error) {
	order, err := s.repository.FindOrderByOrderID(ctx, orderId)
	if err != nil {
		return models2.GetDeliveryInfoResp{}, err
	}

	if order.DeliveryOrderID == "" {
		s.logger.Infof("3pl delivery order is empty for order_id: %s", orderId)
		return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
	}

	store, err := s.storeService.GetByID(ctx, order.RestaurantID)
	if err != nil {
		log.Err(err).Msgf("error: get store for store id: %s", order.RestaurantID)
		return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
	}

	delivery, err := s.deliveryRepository.GetDeliveryByDeliveryID(ctx, order.DeliveryOrderID)
	if err != nil {
		log.Err(err).Msgf("error: get delivery by delivery id: %s", delivery.Id)
		return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
	}

	noCourierTime := time.Now().UTC().Add(time.Duration(-10) * time.Minute)
	if delivery.Courier.TrackingUrl == "" && order.OrderTime.Value.After(noCourierTime) {
		err = s.telegramService.SendMessageToQueue(telegram.NoCourier, order, store, "", "", "", models3.Product{})
		if err != nil {
			log.Err(err).Msgf("cant send no courier message to queue, order_id: %s", order.OrderID)
		}
	}

	maxCookingTime, err := s.getMaxCookingTimeForStore(ctx, store.MenuID, order.Products)
	if err != nil {
		log.Err(err).Msgf("error: get max cooking time for store menu id: %s", store.MenuID)
		return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
	}

	if maxCookingTime == 0 {
		maxCookingTime = int(store.QRMenu.CookingTime)
	}

	return models2.GetDeliveryInfoResp{
		DeliveryID: delivery.Id,
		Statuses:   delivery.StatusHistory,
		DeliveryOrder: models2.GetDeliveryOrderTrackingUrl{
			TrackingUrl:  delivery.Courier.TrackingUrl,
			Latitude:     delivery.Courier.Latitude,
			Longitude:    delivery.Courier.Longitude,
			CourierPhone: delivery.Courier.CourierPhone,
		},
		CookingTime: maxCookingTime,
		BusyMode:    store.QRMenu.BusyMode,
		CancelState: delivery.CancelState,
	}, nil
}

// Todo temporary
func (s *ServiceImpl) GetDeliveryInfoForReport(ctx context.Context, orderId string) (models2.GetDeliveryInfoResp, error) {
	order, err := s.repository.FindOrderByOrderID(ctx, orderId)
	if err != nil {
		return models2.GetDeliveryInfoResp{}, err
	}
	if order.DeliveryOrderID == "" {
		s.logger.Infof("3pl delivery order is empty for order_id: %s", orderId)
		return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
	}
	store, err := s.storeService.GetByID(ctx, order.RestaurantID)
	if err != nil {
		log.Err(err).Msgf("error: get store for store id: %s", order.RestaurantID)
		return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
	}

	switch len([]rune(order.DeliveryOrderID)) {
	case 1, 2, 3, 4, 5:
		deliveryFromPQ, err := s.GetDeliveryFromPostgres(ctx, order.DeliveryOrderID)
		if err != nil {
			log.Err(err).Msgf("error: get delivery from postgres: %s", order.DeliveryOrderID)
			return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
		}
		deliveryStatusesFromPQ, err := s.GetDeliveryStatusFromPostgres(ctx, order.DeliveryOrderID)
		if err != nil {
			log.Err(err).Msgf("error: get delivery status from postgres: %s", order.DeliveryOrderID)
			return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
		}
		maxCookingTime, err := s.getMaxCookingTimeForStore(ctx, store.MenuID, order.Products)
		if err != nil {
			log.Err(err).Msgf("error: get max cooking time for store menu id: %s", store.MenuID)
			return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
		}
		if maxCookingTime == 0 {
			maxCookingTime = int(store.QRMenu.CookingTime)
		}
		return models2.GetDeliveryInfoResp{
			DeliveryID: deliveryFromPQ.Id,
			Statuses:   models2.ToModel(deliveryStatusesFromPQ),
			DeliveryOrder: models2.GetDeliveryOrderTrackingUrl{
				TrackingUrl:  deliveryFromPQ.Courier.TrackingUrl,
				CourierPhone: deliveryFromPQ.Courier.Phone,
			},
			CookingTime: maxCookingTime,
			BusyMode:    store.QRMenu.BusyMode,
		}, nil
	default:
		deliveryFromMongo, err := s.deliveryRepository.GetDeliveryByDeliveryID(ctx, order.DeliveryOrderID)
		if err != nil {
			log.Err(err).Msgf("error: get delivery from mongo: %s", order.DeliveryOrderID)
			return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
		}
		maxCookingTime, err := s.getMaxCookingTimeForStore(ctx, store.MenuID, order.Products)
		if err != nil {
			log.Err(err).Msgf("error: get max cooking time for store menu id: %s", store.MenuID)
			return models2.GetDeliveryInfoResp{}, models.ErrDeliveryOrderIdIsEmpty
		}
		if maxCookingTime == 0 {
			maxCookingTime = int(store.QRMenu.CookingTime)
		}
		return models2.GetDeliveryInfoResp{
			DeliveryID: deliveryFromMongo.Id,
			Statuses:   deliveryFromMongo.StatusHistory,
			DeliveryOrder: models2.GetDeliveryOrderTrackingUrl{
				TrackingUrl:  deliveryFromMongo.Courier.TrackingUrl,
				Latitude:     deliveryFromMongo.Courier.Latitude,
				Longitude:    deliveryFromMongo.Courier.Longitude,
				CourierPhone: deliveryFromMongo.Courier.CourierPhone,
			},
			CookingTime: maxCookingTime,
			BusyMode:    store.QRMenu.BusyMode,
			CancelState: deliveryFromMongo.CancelState,
		}, nil
	}
}

func (s *ServiceImpl) GetDeliveryProdiver(ctx context.Context, deliveryID string) (string, error) {
	if deliveryID == "" {
		return "", models.ErrDeliveryOrderIdIsEmpty
	}

	switch len([]rune(deliveryID)) {
	case 1, 2, 3, 4, 5:
		deliveryFromPQ, err := s.GetDeliveryFromPostgres(ctx, deliveryID)
		if err != nil {
			log.Err(err).Msgf("error: get delivery from postgres: %s", deliveryID)
			return "", models.ErrDeliveryOrderIdIsEmpty
		}
		return deliveryFromPQ.DeliveryService.Name, nil
	default:
		deliveryFromMongo, err := s.deliveryRepository.GetDeliveryByDeliveryID(ctx, deliveryID)
		if err != nil {
			log.Err(err).Msgf("error: get delivery from mongo: %s", deliveryID)
			return "", models.ErrDeliveryOrderIdIsEmpty
		}
		return deliveryFromMongo.DeliveryService, nil
	}
}

// Todo temporary
func (s *ServiceImpl) GetDeliveryStatus(ctx context.Context, deliveryID string) (string, error) {
	if deliveryID == "" {
		return "", models.ErrDeliveryOrderIdIsEmpty
	}
	switch len([]rune(deliveryID)) {
	case 1, 2, 3, 4, 5:
		deliveryStatusesFromPQ, err := s.GetDeliveryStatusFromPostgres(ctx, deliveryID)
		if err != nil {
			log.Err(err).Msgf("error: get delivery status from postgres: %s", deliveryID)
			return "", fmt.Errorf("error: get delivery status from postgres: %s", deliveryID)
		}
		if deliveryStatusesFromPQ != nil && len(deliveryStatusesFromPQ) != 0 {
			return deliveryStatusesFromPQ[len(deliveryStatusesFromPQ)-1].Status, nil
		}
	default:
		deliveryFromMongo, err := s.deliveryRepository.GetDeliveryByDeliveryID(ctx, deliveryID)
		if err != nil {
			log.Err(err).Msgf("error: get delivery from mongo: %s", deliveryID)
			return "", fmt.Errorf("error: get delivery from mongo: %s", deliveryID)
		}
		return deliveryFromMongo.Status, nil
	}
	return "", fmt.Errorf("error: get delivery status: %s", deliveryID)
}

// Todo temporary
func (s *ServiceImpl) GetDeliveryFromPostgres(ctx context.Context, deliveryID string) (models2.DeliveryOrderPQ, error) {
	path := fmt.Sprintf("/v1/delivery/get-delivery?delivery_id=%s", deliveryID)

	var (
		result      models2.DeliveryOrderPQ
		errResponse models.ErrorResponse
	)

	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).
		Get(path)
	if err != nil {
		return models2.DeliveryOrderPQ{}, err
	}

	log.Info().Msgf("request 3pl service get delivery price %+v\n", resp.Request.URL)

	if resp.IsError() {
		return models2.DeliveryOrderPQ{}, fmt.Errorf("error get 3pl delivery price: %v", resp.Error())
	}

	log.Info().Msgf("response 3pl service get delivery price %+v\n", resp.String())

	return result, nil
}

// Todo temporary
func (s *ServiceImpl) GetDeliveryStatusFromPostgres(ctx context.Context, deliveryID string) ([]models2.DeliveryStatusHistory, error) {
	path := fmt.Sprintf("/v1/delivery/get-delivery-status?delivery_id=%s", deliveryID)

	var (
		result      []models2.DeliveryStatusHistory
		errResponse models.ErrorResponse
	)

	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).
		Get(path)
	if err != nil {
		return []models2.DeliveryStatusHistory{}, err
	}

	log.Info().Msgf("request 3pl service get delivery price %+v\n", resp.Request.URL)

	if resp.IsError() {
		return []models2.DeliveryStatusHistory{}, fmt.Errorf("error get 3pl delivery price: %v", resp.Error())
	}

	log.Info().Msgf("response 3pl service get delivery price %+v\n", resp.String())

	return result, nil
}

func (s *ServiceImpl) Get3plDeliveryInfo(ctx context.Context, deliveries []string) ([]models2.GetDeliveryInfoResp, error) {

	var deliveriesInfoResponse []models2.GetDeliveryInfoResp

	if len(deliveries) == 0 {
		return nil, nil
	}

	deliveriesInfo, err := s.deliveryRepository.GetDeliveries(ctx, deliveries)
	if err != nil {
		log.Err(err).Msgf("error: get delivery info for delivery ids: %v", deliveries)
		return nil, err
	}

	for _, deliveryInfo := range deliveriesInfo {

		order, err := s.repository.GetOrderBy3plDeliveryID(ctx, deliveryInfo.Id)
		if err != nil {
			log.Err(err).Msgf("error: get order for 3pl delivery order id: %s", deliveryInfo.Id)
			continue
		}

		store, err := s.storeService.GetByID(ctx, order.RestaurantID)
		if err != nil {
			log.Err(err).Msgf("error: get store for store id: %s", order.RestaurantID)
			continue
		}

		noCourierTime := time.Now().UTC().Add(time.Duration(-10) * time.Minute)
		if deliveryInfo.Courier.TrackingUrl == "" && order.OrderTime.Value.After(noCourierTime) {
			err = s.telegramService.SendMessageToQueue(telegram.NoCourier, order, store, "", "", "", models3.Product{})
			if err != nil {
				log.Err(err).Msgf("cant send no courier message to queue, order_id: %s", order.OrderID)
			}
		}

		maxCookingTime, err := s.getMaxCookingTimeForStore(ctx, store.MenuID, order.Products)
		if err != nil {
			log.Err(err).Msgf("error: get max cooking time for store menu id: %s", store.MenuID)
			continue
		}

		if maxCookingTime == 0 {
			maxCookingTime = int(store.QRMenu.CookingTime)
		}

		deliveriesInfoResponse = append(deliveriesInfoResponse, models2.GetDeliveryInfoResp{
			DeliveryID: deliveryInfo.Id,
			Statuses:   deliveryInfo.StatusHistory,
			DeliveryOrder: models2.GetDeliveryOrderTrackingUrl{
				TrackingUrl:  deliveryInfo.Courier.TrackingUrl,
				Latitude:     deliveryInfo.Courier.Latitude,
				Longitude:    deliveryInfo.Courier.Longitude,
				CourierPhone: deliveryInfo.Courier.CourierPhone,
			},
			CookingTime: maxCookingTime,
			BusyMode:    store.QRMenu.BusyMode,
			CancelState: deliveryInfo.CancelState,
		})
	}

	return deliveriesInfoResponse, nil
}

func (s *ServiceImpl) getMaxCookingTimeForStore(ctx context.Context, posMenuID string, products []models2.OrderProduct) (int, error) {

	posMenu, err := s.menuClient.GetMenuByID(ctx, posMenuID)
	if err != nil {
		return 0, err
	}

	mapItems := make(map[string]bool)
	for _, cartProduct := range products {
		mapItems[cartProduct.ID] = true
	}

	cookingTime := 0

	for _, posProduct := range posMenu.Products {
		if mapItems[posProduct.ExtID] && cookingTime < int(posProduct.CookingTime) {
			cookingTime = int(posProduct.CookingTime)
		}
	}

	return cookingTime, nil
}

func (cl *ServiceImpl) get3plDeliveryInfo(ctx context.Context, deliveryOrderId string) (models2.GetDeliveryInfoResp, error) {
	path := fmt.Sprintf("/v1/delivery/receive?delivery_id=%s", deliveryOrderId)

	var (
		errResponse   models.ErrorResponse
		deliveryOrder models2.Delivery3plOrder
	)
	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&deliveryOrder).
		Get(path)
	if err != nil {
		return models2.GetDeliveryInfoResp{}, err
	}

	log.Info().Msgf("request get 3pl delivery info %+v \n", resp.Request.Body)
	log.Info().Msgf("response get 3pl delivery info %+v", string(resp.Body()))

	if resp.IsError() {
		return models2.GetDeliveryInfoResp{}, fmt.Errorf("get 3pl delivery info error: %s, %s with status: %s", resp.Error(), errResponse.Message, resp.Status())
	}

	return models2.GetDeliveryInfoResp{
		Statuses:      deliveryOrder.StatusHistory,
		DeliveryOrder: deliveryOrder.Courier,
	}, nil
}

func (s *ServiceImpl) SetKwaaka3plDispatcher(ctx context.Context, req models2.SetKwaaka3plDispatcherRequest) error {
	order, err := s.repository.FindOrderByID(ctx, req.OrderID)
	if err != nil {
		return err
	}

	if order.IsMarketplace || !order.SendCourier {
		return errors.Errorf("can't set dispatcher, isMarketPlace: %v, isPickupByCustomer: %v for order id: %s", order.IsMarketplace, order.IsPickedUpByCustomer, order.ID)
	}

	if order.DeliveryOrderID != "" {
		delivery, err := s.deliveryRepository.GetDeliveryByDeliveryID(ctx, order.DeliveryOrderID)
		if err != nil {
			log.Err(err).Msgf("error: get delivery by delivery order id: %s", order.DeliveryOrderID)
		}

		if delivery.Status == models.Delivered ||
			delivery.Status == models.Cancelled ||
			delivery.Status == models.Returning ||
			delivery.Status == models.Returned ||
			delivery.Status == models.Failed {

			log.Info().Msgf("dispatcher is already setted. Save delivery info to history_3pl_delivery_info and create new delivery for order id: %s", order.ID)
			err = s.repository.Save3plDeliveryHistoryAndSetEmptyDispatcherService(ctx, order.ID, models2.History3plDelivery{
				DeliveryOrderID:            order.DeliveryOrderID,
				DeliveryDispatcher:         order.DeliveryDispatcher,
				FullDeliveryPrice:          order.FullDeliveryPrice,
				RestaurantPayDeliveryPrice: order.RestaurantPayDeliveryPrice,
				KwaakaChargedDeliveryPrice: order.KwaakaChargedDeliveryPrice,
				DeliveryAddress:            order.DeliveryAddress,
				Customer:                   order.Customer,
			}, req.DeliveryAddress, req.Customer)
		}
		order, err = s.repository.FindOrderByID(ctx, req.OrderID)
		if err != nil {
			return err
		}
	}

	order.DeliveryDispatcher = req.Dispatcher
	order.DeliveryOrderPromiseID = req.DeliveryOrderPromiseID
	order.FullDeliveryPrice = req.FullDeliveryPrice

	if err = s.repository.UpdateOrder(ctx, order); err != nil {
		return err
	}
	s.logger.Infof("set new 3pl dispatcher for order: %s", order.OrderID)

	store, err := s.storeService.GetByID(ctx, order.RestaurantID)
	if err != nil {
		return err
	}
	if !store.Kwaaka3PL.Is3pl {
		return errors.New("store is not integrated with kwaaka 3pl")
	}

	items := make([]models.Item, 0, len(order.Products))
	for _, product := range order.Products {
		items = append(items, models.Item{
			Name:     product.Name,
			ID:       product.ID,
			Quantity: product.Quantity,
			Price:    product.Price.Value,
		})
	}

	if err = s.Create3plOrder(ctx, models.CreateDeliveryRequest{
		ID:                order.ID,
		FullDeliveryPrice: order.FullDeliveryPrice,
		Provider:          order.DeliveryDispatcher,
		PickUpTime:        time.Now().Add(time.Duration(order.DispatcherDeliveryTime) * time.Minute),
		DeliveryAddress: models.Address{
			Label:        order.DeliveryAddress.Label,
			Lat:          order.DeliveryAddress.Latitude,
			Lon:          order.DeliveryAddress.Longitude,
			Comment:      order.DeliveryAddress.Comment,
			BuildingName: order.DeliveryAddress.BuildingName,
			Street:       order.DeliveryAddress.Street,
			Flat:         order.DeliveryAddress.Flat,
			Porch:        order.DeliveryAddress.Porch,
			Floor:        order.DeliveryAddress.Floor,
		},
		StoreAddress: models.Address{
			Label:   store.Address.City + ", " + store.Address.Street,
			Lon:     store.Address.Coordinates.Longitude,
			Lat:     store.Address.Coordinates.Latitude,
			Comment: store.Address.Entrance,
		},
		CustomerInfo: models.CustomerInfo{
			Name:  s.setCustomerName(order.Customer.Name),
			Phone: s.setCustomerPhoneNumber(order.Customer.PhoneNumber),
			Email: s.setCustomerEmail(order.Customer.Email),
		},
		StoreInfo: models.StoreInfo{
			Name:  store.Name,
			Phone: store.StorePhoneNumber,
			Email: store.Settings.Email,
		},
		PickUpCode:      order.PickUpCode,
		Currency:        order.Currency,
		Comment:         order.SpecialRequirements,
		Items:           items,
		ExternalStoreID: store.Kwaaka3PL.IndriveStoreID,
		TaxiClass:       store.Kwaaka3PL.TaxiClass,
	}); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) setCustomerEmail(email string) string {
	customerEmail := email
	if customerEmail == "" {
		customerEmail = models2.Default3plCustomerEmail
	}

	return customerEmail
}

func (s *ServiceImpl) setCustomerPhoneNumber(phoneNumber string) string {
	customerPhone := phoneNumber
	if customerPhone == "" {
		customerPhone = models2.Default3plCustomerPhone
	}

	return customerPhone
}

func (s *ServiceImpl) setCustomerName(name string) string {
	customerName := name
	if customerName == "" {
		customerName = models2.Default3plCustomerName
	}

	return customerName
}

func (s *ServiceImpl) GetOrdersByCustomerPhone(ctx context.Context, query models2.GetOrdersByCustomerPhoneRequest) (models2.GetOrdersByCustomerPhoneResponse, error) {
	orders, totalCount, err := s.repository.GetAllOrders(ctx, selector.EmptyOrderSearch().
		SetRestaurants([]string{query.RestaurantID}).
		SetCustomerNumber(query.CustomerPhone).
		SetPage(query.Pagination.Page).
		SetLimit(query.Pagination.Limit).SetSorting("order_time.value", -1))
	if err != nil {
		return models2.GetOrdersByCustomerPhoneResponse{}, err
	}

	averageBill, err := s.repository.GetAverageBill(ctx,
		selector.EmptyOrderSearch().
			SetRestaurants([]string{query.RestaurantID}).
			SetCustomerNumber(query.CustomerPhone))
	if err != nil {
		return models2.GetOrdersByCustomerPhoneResponse{}, err
	}

	resOrders := make([]models2.OrderByCustomerPhone, 0, len(orders))
	for i := range orders {
		ord := models2.OrderByCustomerPhone{
			ID:                  orders[i].ID,
			DeliveryService:     orders[i].DeliveryService,
			RestaurantID:        orders[i].RestaurantID,
			RestaurantName:      orders[i].RestaurantName,
			OrderID:             orders[i].OrderID,
			OrderCode:           orders[i].OrderCode,
			Status:              orders[i].Status,
			StatusesHistory:     orders[i].StatusesHistory,
			OrderTime:           orders[i].OrderTime,
			Customer:            orders[i].Customer,
			Products:            orders[i].Products,
			EstimatedTotalPrice: orders[i].EstimatedTotalPrice,
		}
		resOrders = append(resOrders, ord)
	}

	var result = models2.GetOrdersByCustomerPhoneResponse{
		Orders: resOrders,
		CustomerOrderHistory: models2.CustomerOrderHistory{
			Phone:       query.CustomerPhone,
			AverageBill: averageBill,
			Amount:      totalCount,
		},
	}

	if len(orders) != 0 {
		result.CustomerOrderHistory.Name = orders[0].Customer.Name
		result.CustomerOrderHistory.LastOrder = orders[0].OrderTime.Value.Time
	}

	return result, nil
}

func (s *ServiceImpl) GetCustomerByDeliveryId(ctx context.Context, deliveryId string) (models2.Customer, string, error) {
	order, err := s.repository.FindOrderByDeliveryOrderId(ctx, deliveryId)
	if err != nil {
		return models2.Customer{}, "", err
	}

	return order.Customer, order.RestaurantID, nil
}

func (s *ServiceImpl) GetOrderForTelegramByDeliveryOrderId(ctx context.Context, deliveryId string) (models2.OrderInfoForTelegramMsg, error) {

	order, err := s.repository.FindOrderByDeliveryOrderId(ctx, deliveryId)
	if err != nil {
		return models2.OrderInfoForTelegramMsg{}, err
	}

	store, err := s.storeService.GetByID(ctx, order.RestaurantID)
	if err != nil {
		return models2.OrderInfoForTelegramMsg{}, nil
	}

	return models2.OrderInfoForTelegramMsg{
		RestaurantName:        store.Name,
		RestaurantAddress:     store.Address.City + store.Address.Street,
		RestaurantPhoneNumber: store.StorePhoneNumber,
		OrderId:               order.ID,
		Id3plOrder:            deliveryId,
		CustomerPhoneNumber:   order.Customer.PhoneNumber,
		CustomerName:          order.Customer.Name,
		CustomerAddress:       order.DeliveryAddress.Label,
		DeliveryService:       order.DeliveryDispatcher,
	}, nil
}

func (s *ServiceImpl) BulkCreate3plOrder(ctx context.Context, orders []models2.Order, byAdmin bool) error {

	s.logger.Infof("start to bul create 3pl orders: %v", orders)

	for _, order := range orders {
		if !byAdmin {
			if order.Status == models2.FAILED.String() || order.Status == models2.CANCELLED_BY_POS_SYSTEM.String() || order.Status == string(models2.STATUS_CANCELLED) || order.Status == string(models2.STATUS_CANCELLED_BY_DELIVERY_SERVICE) || order.Status == string(models2.STATUS_WAIT_SENDING) || order.Status == string(models2.STATUS_SKIPPED) {
				continue
			}
		}

		s.logger.Infof("order that need 3pl: %v", order)

		store, err := s.storeService.GetByID(ctx, order.StoreID)
		if err != nil {
			s.telegramService.SendMessageToQueue(telegram.ThirdPartyError, order, store, err.Error(),
				fmt.Sprintf("Ошибка при создании доставки для заказа %s, ресторан %s не найден, error: %s",
					order.OrderID, order.StoreID, err.Error()), "", models3.Product{})
			return err
		}

		items := make([]models.Item, 0, len(order.Products))
		for _, product := range order.Products {
			items = append(items, models.Item{
				Name:     product.Name,
				ID:       product.ID,
				Quantity: product.Quantity,
				Price:    product.Price.Value,
			})
		}

		err = s.Create3plOrder(ctx, models.CreateDeliveryRequest{
			ID:                order.ID,
			FullDeliveryPrice: order.FullDeliveryPrice,
			Provider:          order.DeliveryDispatcher,
			PickUpTime:        time.Now().Add(time.Duration(order.DispatcherDeliveryTime) * time.Minute),
			DeliveryAddress: models.Address{
				Label:        order.DeliveryAddress.Label,
				Lat:          order.DeliveryAddress.Latitude,
				Lon:          order.DeliveryAddress.Longitude,
				Comment:      order.DeliveryAddress.Comment,
				City:         order.DeliveryAddress.City,
				BuildingName: order.DeliveryAddress.BuildingName,
				Street:       order.DeliveryAddress.Street,
				Flat:         order.DeliveryAddress.Flat,
				Porch:        order.DeliveryAddress.Porch,
				Floor:        order.DeliveryAddress.Floor,
			},
			StoreAddress: models.Address{
				Label:   store.Address.City + ", " + store.Address.Street,
				Lon:     store.Address.Coordinates.Longitude,
				Lat:     store.Address.Coordinates.Latitude,
				Comment: store.Address.Entrance,
			},
			CustomerInfo: models.CustomerInfo{
				Name:  s.setCustomerName(order.Customer.Name),
				Phone: s.setCustomerPhoneNumber(order.Customer.PhoneNumber),
				Email: s.setCustomerEmail(order.Customer.Email),
			},
			StoreInfo: models.StoreInfo{
				Name:  store.Name,
				Phone: store.StorePhoneNumber,
				Email: store.Settings.Email,
			},
			PickUpCode:      order.PickUpCode,
			Currency:        order.Currency,
			Comment:         order.SpecialRequirements,
			Items:           items,
			ExternalStoreID: store.Kwaaka3PL.IndriveStoreID,
			TaxiClass:       store.Kwaaka3PL.TaxiClass,
		})
		if err != nil {
			log.Err(err).Msg("error while bulk create 3pl order")
			return err
		}
	}
	return nil
}

func (s *ServiceImpl) ListPotentialProviders(ctx context.Context, req models.ListProvidersRequest) ([]models.ProviderResponse, error) {
	path := "/v1/delivery/providers"

	var result []models.ProviderResponse

	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetBody(&req).
		SetResult(&result).
		Post(path)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error list potential providers: %v", resp.Error())
	}

	return result, nil
}

func (s *ServiceImpl) GetDeliveryPrice(ctx context.Context, deliveryOrderID string) (float64, error) {
	path := fmt.Sprintf("/v1/delivery/delivery-price?delivery_id=%s", deliveryOrderID)

	var (
		result      float64
		errResponse models.ErrorResponse
	)

	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).
		Post(path)
	if err != nil {
		return 0, err
	}

	log.Info().Msgf("request 3pl service get delivery price %+v\n", resp.Request.URL)

	if resp.IsError() {
		return 0, fmt.Errorf("error get 3pl delivery price: %v", resp.Error())
	}

	log.Info().Msgf("response 3pl service get delivery price %+v\n", resp.String())

	return result, nil
}

func (s *ServiceImpl) GetDeliveryPriceFromPQ(ctx context.Context, deliveryOrderID string) (float64, error) {
	path := fmt.Sprintf("/v1/delivery/get-delivery-price?delivery_id=%s", deliveryOrderID)

	var (
		result      float64
		errResponse models.ErrorResponse
	)

	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&result).
		Get(path)
	if err != nil {
		return 0, err
	}

	log.Info().Msgf("request 3pl service get delivery price %+v\n", resp.Request.URL)

	if resp.IsError() {
		return 0, fmt.Errorf("error get 3pl delivery price: %v", resp.Error())
	}

	log.Info().Msgf("response 3pl service get delivery price %+v\n", resp.String())

	return result, nil
}

func (s *ServiceImpl) GetDeliveryDispatcherPrices(ctx context.Context, deliveryIDs []string) (models.GetDeliveryDispatcherPricesResponse, error) {
	path := fmt.Sprintf("/v1/delivery/delivery-prices")

	var (
		result      models.GetDeliveryDispatcherPricesResponse
		errResponse models.ErrorResponse
	)

	req := models.GetDeliveryDispatcherPricesRequest{
		DeliveryIDs: deliveryIDs,
	}

	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(req).
		SetResult(&result).
		Post(path)
	if err != nil {
		return models.GetDeliveryDispatcherPricesResponse{}, err
	}

	log.Info().Msgf("request 3pl service get delivery prices %+v\n", resp.Request.Body)

	if resp.IsError() {
		return models.GetDeliveryDispatcherPricesResponse{}, fmt.Errorf("error get 3pl delivery prices: %v", resp.Error())
	}

	log.Info().Msgf("response 3pl service get delivery prices %+v\n", string(resp.Body()))

	return result, nil
}

func (s *ServiceImpl) Instant3plOrder(ctx context.Context, order models2.Order) error {
	st, err := s.storeService.GetByID(ctx, order.RestaurantID)
	if err != nil {
		return err
	}

	if order.Status == models2.FAILED.String() || order.Status == models2.CANCELLED_BY_POS_SYSTEM.String() || order.Status == string(models2.STATUS_CANCELLED) || order.Status == string(models2.STATUS_CANCELLED_BY_DELIVERY_SERVICE) || order.Status == string(models2.STATUS_WAIT_SENDING) || order.Status == string(models2.STATUS_SKIPPED) {
		return fmt.Errorf("cannot create 3pl: order's status is %s", order.Status)
	}

	s.logger.Infof("order that need instant 3pl: %v", order)

	items := make([]models.Item, 0, len(order.Products))
	for _, product := range order.Products {
		items = append(items, models.Item{
			Name:     product.Name,
			ID:       product.ID,
			Quantity: product.Quantity,
			Price:    product.Price.Value,
		})
	}

	err = s.Create3plOrder(ctx, models.CreateDeliveryRequest{
		ID:                order.ID,
		FullDeliveryPrice: order.FullDeliveryPrice,
		Provider:          order.DeliveryDispatcher,
		PickUpTime:        time.Now().Add(time.Duration(order.DispatcherDeliveryTime) * time.Minute),
		DeliveryAddress: models.Address{
			Label:        order.DeliveryAddress.Label,
			Lat:          order.DeliveryAddress.Latitude,
			Lon:          order.DeliveryAddress.Longitude,
			Comment:      order.DeliveryAddress.Comment,
			City:         order.DeliveryAddress.City,
			BuildingName: order.DeliveryAddress.BuildingName,
			Street:       order.DeliveryAddress.Street,
			Flat:         order.DeliveryAddress.Flat,
			Porch:        order.DeliveryAddress.Porch,
			Floor:        order.DeliveryAddress.Floor,
		},
		StoreAddress: models.Address{
			Label:   st.Address.City + ", " + st.Address.Street,
			Lon:     st.Address.Coordinates.Longitude,
			Lat:     st.Address.Coordinates.Latitude,
			Comment: st.Address.Entrance,
		},
		CustomerInfo: models.CustomerInfo{
			Name:  s.setCustomerName(order.Customer.Name),
			Phone: s.setCustomerPhoneNumber(order.Customer.PhoneNumber),
			Email: s.setCustomerEmail(order.Customer.Email),
		},
		StoreInfo: models.StoreInfo{
			Name:  st.Name,
			Phone: st.StorePhoneNumber,
			Email: st.Settings.Email,
		},
		PickUpCode:      order.PickUpCode,
		Currency:        order.Currency,
		Comment:         order.SpecialRequirements,
		Items:           items,
		ExternalStoreID: st.Kwaaka3PL.IndriveStoreID,
		TaxiClass:       st.Kwaaka3PL.TaxiClass,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) GetOrderByOrderID(ctx context.Context, orderID string) (models2.Order, error) {
	return s.repository.FindOrderByOrderID(ctx, orderID)
}

func (s *ServiceImpl) NoDispatcherMessage(ctx context.Context) error {
	log.Info().Msgf("no dispatcher message start")

	timeNow := time.Now().UTC()

	orders, _, err := s.repository.GetAllOrders(ctx, selector.EmptyOrderSearch().
		SetDeliveryServices([]string{models2.KWAAKA_ADMIN.String(), models2.QRMENU.String()}).
		SetOrderTimeFrom(timeNow.Add(-time.Hour*1)).
		SetOrderTimeTo(timeNow))
	if err != nil {
		return err
	}

	log.Info().Msgf("kwaaka_admin, qr_menu order quantity: %d", len(orders))

	for _, order := range orders {

		switch order.Status {
		case models2.STATUS_DELIVERED.ToString(), models2.STATUS_CLOSED.ToString(), models2.STATUS_CANCELLED_BY_DELIVERY_SERVICE.ToString(), models2.STATUS_CANCELLED_BY_POS_SYSTEM.ToString(), models2.STATUS_FAILED.ToString(), models2.STATUS_CANCELLED.ToString():
			continue
		}

		for _, status := range order.StatusesHistory {
			if status.Name == string(models2.STATUS_COOKING_COMPLETE) && status.Time.Before(timeNow.Add(-time.Minute*5)) && order.DeliveryDispatcher == "" && order.SendCourier == true {
				log.Info().Msgf("no dispatcher order, order_id: %s", order.OrderID)

				store, err := s.storeService.GetByID(ctx, order.RestaurantID)
				if err != nil {
					log.Err(fmt.Errorf("get store by id error, store id: %s", order.StoreID))
					store = storeModels.Store{}
				}

				if err := s.telegramService.SendMessageToQueue(telegram.NoDeliveryDispatcher, order, store, "", "", "", models3.Product{}); err != nil {
					log.Err(err).Msgf("send message to queue telegram error")
					break
				}

				log.Info().Msgf("send message success, order_id: %s", order.OrderID)
				break
			}
		}
	}

	return nil
}

func (s *ServiceImpl) ActualizeDeliveryInfoByDeliveryIDs(ctx context.Context, deliveryIDs []string) error {

	s.logger.Errorf("start to actualize delivery info by delivery ids: %v", deliveryIDs)

	var (
		errResponse models.ErrorResponse
		req         struct {
			DeliveryIDs []string `json:"delivery_ids"`
		}
		InProgressDeliveryIDs []string
	)

	for _, deliveryID := range deliveryIDs {

		delivery, err := s.deliveryRepository.GetDeliveryByDeliveryID(ctx, deliveryID)
		if err != nil {
			return err
		}
		if delivery.Status == models.OrderCreated ||
			delivery.Status == models.PerformerLookup ||
			delivery.Status == models.ComingToPickup ||
			delivery.Status == models.PickedUp ||
			delivery.Status == models.Returning {

			InProgressDeliveryIDs = append(InProgressDeliveryIDs, deliveryID)
		}

	}

	path := "/v1/delivery/actualize-delivery"

	req.DeliveryIDs = InProgressDeliveryIDs

	resp, err := s.restyClient.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(req).
		Post(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("request actualize delivery info by deliveryIDs %+v \n", resp.Request.Body)
	log.Info().Msgf("response actualize delivery info by deliveryIDs %+v", string(resp.Body()))

	if resp.IsError() {
		return fmt.Errorf("actualize delivery info by deliveryIDS response error: %s, %s", resp.Error(), errResponse.Message)
	}

	return nil
}

func (s *ServiceImpl) PerformerLookupMoreThan15Minute(ctx context.Context) error {
	log.Info().Msgf("performer lookup more than 15 minute cron start")

	timeNow := time.Now().UTC()

	deliveries, err := s.deliveryRepository.GetAllDeliveries(ctx, selector.EmptyDelivery3plSearch().
		SetStatus(models.PerformerLookup).
		SetUpdatedTimeTo(timeNow.Add(-time.Minute*15)).
		SetCreatedTimeFrom(timeNow.Add(-time.Hour*1)).
		SetCreatedTimeTo(timeNow))
	if err != nil {
		log.Err(err).Msgf("GetAllDeliveries error: get deliveries with query")
		return err
	}

	for _, delivery := range deliveries {
		log.Info().Msgf("performer lookup more than 15 minute for delivery order id: %s", delivery.Id)
		order, err := s.repository.GetOrderBy3plDeliveryID(ctx, delivery.Id)
		if err != nil {
			log.Err(err).Msgf("error: GetOrderBy3plDeliveryID for delivery order id: %s", delivery.Id)
			continue
		}
		store, err := s.storeService.GetByID(ctx, order.RestaurantID)
		if err != nil {
			log.Err(err).Msgf("error: GetByID for id: %s", order.RestaurantID)
			continue
		}
		orderInfo, err := s.GetOrderForTelegramByDeliveryOrderId(ctx, delivery.Id)
		if err != nil {
			log.Err(err).Msgf("error: GetOrderForTelegramByDeliveryOrderId for delivery order id: %s", delivery.Id)
			continue
		}

		msg := "<b>Долгий поиск курьера (15 минут). Необходимо обратиться в службу поддержки провайдера</b>\n"
		msgOrder := s.convertOrderToMessage(orderInfo)
		if err := s.telegramService.SendMessageToQueue(telegram.ThirdPartyError, order, store, "", msg+msgOrder, "", models3.Product{}); err != nil {
			log.Err(err).Msgf("error: SendMessageToQueue for delivery order id: %s", delivery.Id)
		}

		if err := s.Cancel3plOrder(ctx, order.ID, orderInfo); err != nil {
			log.Err(err).Msgf("error: Cancel3plOrder for delivery order id: %s", delivery.Id)
			continue
		}

		if err := s.Save3plHistory(ctx, delivery.Id, models2.DeliveryAddress{}, models2.Customer{}); err != nil {
			log.Err(err).Msgf("error: Save3plHistory for delivery order id: %s", delivery.Id)
			continue
		}

		if err := s.BulkCreate3plOrder(ctx, []models2.Order{order}, false); err != nil {
			log.Err(err).Msgf("error: BulkCreate3plOrder for delivery order id: %s", delivery.Id)
			continue
		}
		log.Info().Msgf("successfully create delivery with cron: performer lookup more than 15 minute for delivery order id: %s", delivery.Id)
	}

	return nil
}

func (s *ServiceImpl) convertOrderToMessage(req models2.OrderInfoForTelegramMsg) string {

	return fmt.Sprintf("<b>Ресторан: </b> %s\n"+
		"<b>Адрес ресторана: </b> %s\n"+
		"<b>Номер телефона ресторана:  </b> %s\n"+
		"<b>Сервис доставки:</b> %s\n"+
		"<b>ID заказа:</b> %s\n"+
		"<b>ID 3pl:</b> %s\n"+
		"<b>Данные о клиенте:</b>\n"+
		"<b>Имя:</b> %s\n"+
		"<b>Номер:</b> %s\n"+
		"<b>Адрес:</b> %s\n\n\n", req.RestaurantName, req.RestaurantAddress, req.RestaurantPhoneNumber,
		req.DeliveryService, req.OrderId, req.Id3plOrder, req.CustomerName, req.CustomerPhoneNumber, req.CustomerAddress)

}

func (s *ServiceImpl) InsertChangeDeliveryHistory(ctx context.Context, deliveryID string, history models2.ChangesHistory) error {
	log.Info().Msgf("insert change delivery history for delivery ID: %s", deliveryID)
	return s.deliveryRepository.InsertChangeHistory(ctx, deliveryID, history.Username, history.Action)
}
