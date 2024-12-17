package order

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/config"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	selector2 "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	MenuClient "github.com/kwaaka-team/orders-core/pkg/menu"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	models2 "github.com/kwaaka-team/orders-core/service/error_solutions/models"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/order_rules"
	paymentRepo "github.com/kwaaka-team/orders-core/service/payment/repository"
	"github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/store"
	storeGroupServicePkg "github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	errConstructor      = errors.New("constructor error")
	errWithNotification = errors.New("error with notification")
	ErrInvalidInterval  = errors.New("invalid interval")
)

type OrderCronService interface {
	ActualizeOrdersStatusByPosType(ctx context.Context, posType string) error
	GetOrderStat(ctx context.Context, interval string) (models.OrderStat, error)
	Get3plOrdersWithoutDriver(ctx context.Context, indriveCallTime int64) ([]models.Order, error)
	SendDeferStatusSubmission(ctx context.Context) error
}

type StatusUpdateService interface {
	UpdateOrderStatus(ctx context.Context, orderID, status, statusDescription string) error
}

type CreationService interface {
	CreateOrder(ctx context.Context, externalStoreID, deliveryService string, aggReq interface{}, storeSecret string) (models.Order, error)
}

type ReviewService interface {
}

type InfoSharingService interface {
	GetShaurmaFoodOrders(ctx context.Context, restaurantID, deliveryService, orderStatus, dateFrom, dateTo string, page, limit int64) ([]models.ShaurmaFoodOrdersInfo, int, error)
	GetOrder(ctx context.Context, orderID string) (models.Order, error)
	GetOrderByID(ctx context.Context, id string) (models.Order, error)
}

type CancellationService interface {
	CancelOrderByAggregator(ctx context.Context, orderID string, delivery string) error
}

type ServiceImpl struct {
	storeService      store.Service
	aggregatorFactory aggregator.Factory
	posFactory        pos.Factory
	repository        Repository

	globalConfig      config.Configuration
	menuClient        MenuClient.Client
	menuService       menu.Service
	storeClient       storeClient.Client
	storeGroupService storeGroupServicePkg.Service
	publisher         *Publisher

	orderRuleService order_rules.Service
	posSender        PosSender

	paymentRepo paymentRepo.PaymentsRepository

	cartService CartService

	errSolution error_solutions.Service
}

func NewServiceImpl(repository Repository) (*ServiceImpl, error) {
	if repository == nil {
		return nil, errors.New("order repository is nil")
	}
	return &ServiceImpl{
		repository: repository,
	}, nil
}

type ServiceFactory struct {
	StoreService store.Service

	AggregatorFactory aggregator.Factory
	PosFactory        pos.Factory
	Repository        Repository

	GlobalConfig      *config.Configuration
	MenuClient        MenuClient.Client
	StoreClient       storeClient.Client
	MenuService       *menu.Service
	StoreGroupService storeGroupServicePkg.Service
	Publisher         *Publisher
	PosSender         PosSender
	OrderRuleService  order_rules.Service

	PaymentRepo paymentRepo.PaymentsRepository

	CartService CartService

	ErrSolution error_solutions.Service
}

func (f ServiceFactory) Create() (*ServiceImpl, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	return &ServiceImpl{
		storeService:      f.StoreService,
		aggregatorFactory: f.AggregatorFactory,
		posFactory:        f.PosFactory,
		repository:        f.Repository,
		globalConfig:      *f.GlobalConfig,
		menuClient:        f.MenuClient,
		storeClient:       f.StoreClient,
		menuService:       *f.MenuService,
		storeGroupService: f.StoreGroupService,
		publisher:         f.Publisher,
		posSender:         f.PosSender,
		orderRuleService:  f.OrderRuleService,
		paymentRepo:       f.PaymentRepo,
		cartService:       f.CartService,
		errSolution:       f.ErrSolution,
	}, nil
}

func (f ServiceFactory) validate() error {
	if f.StoreService == nil {
		return errors.Wrap(errConstructor, "store factory is nil")
	}
	if f.AggregatorFactory == nil {
		return errors.Wrap(errConstructor, "aggregator factory is nil")
	}
	if f.PosFactory == nil {
		return errors.Wrap(errConstructor, "pos factory is nil")
	}
	if f.Repository == nil {
		return errors.Wrap(errConstructor, "repository is nil")
	}
	if f.GlobalConfig == nil {
		return errors.Wrap(errConstructor, "global config is nil")
	}
	if f.MenuClient == nil {
		return errors.Wrap(errConstructor, "menu client is nil")
	}
	if f.MenuService == nil {
		return errors.Wrap(errConstructor, "menu service is nil")
	}
	if f.Publisher == nil {
		return errors.Wrap(errConstructor, "observer3pl service is nil")
	}
	if f.ErrSolution == nil {
		return errors.Wrap(errConstructor, "error solution service is nil")
	}
	return nil
}

func (s *ServiceImpl) ActualizeOrdersStatusByPosType(ctx context.Context, posType string) error {
	orders, err := s.repository.FindActiveOrdersByPosType(ctx, posType, time.Now().UTC().Add(-2*time.Hour))
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, order := range orders {
		wg.Add(1)
		go func(localOrder models.Order) {
			defer wg.Done()

			var localErr error
			store, localErr := s.storeService.GetByID(ctx, localOrder.RestaurantID)
			if localErr != nil {
				log.Err(localErr).Msgf("get store by id, restaurant_id=%s", localOrder.RestaurantID)
				return
			}

			posService, localErr := s.posFactory.GetPosService(models.Pos(localOrder.PosType), store)
			if localErr != nil {
				log.Err(localErr).Msgf("get pos service, store_id=%s, pos_type=%s", store.ID, localOrder.PosType)
				return
			}

			posStatus, localErr := posService.GetOrderStatus(ctx, localOrder)
			if err != nil {
				log.Err(localErr).Msgf("get order status, order_id=%s, pos_type=%s", localOrder.ID, localOrder.PosType)
				return
			}

			if localErr = s.updateOrderStatus(ctx, store, localOrder, posStatus, ""); localErr != nil {
				log.Err(localErr).Msgf("update order status, order_id=%s, pos_type=%s", localOrder.ID, localOrder.PosType)
			}
		}(order)
	}

	wg.Wait()

	return nil
}

func (s *ServiceImpl) SendDeferStatusSubmission(ctx context.Context) error {
	storeDeferSubmission := true
	stores, err := s.storeService.GetStoresBySelectorFilter(ctx, selector2.NewEmptyStoreSearch().SetDeferSubmission(&storeDeferSubmission))
	if err != nil {
		return err
	}

	for _, st := range stores {
		orders, _, err := s.repository.GetAllOrders(ctx, selector.EmptyOrderSearch().
			SetIsDeferSubmission(true).
			SetStoreID(st.ID).
			SetOrderTimeFrom(time.Now().UTC().Add(-1*time.Hour)))
		if err != nil {
			return err
		}

		for _, order := range orders {
			if !s.checkOrder(order) {
				continue
			}

			log.Info().Msgf("starting to send order is ready status for order %s", order.ID)

			if err := s.sendOrderReadyStatusToAgg(ctx, order, st); err != nil {
				log.Err(err).Msgf("error while changing status in restaurant %s , order id  %s", st.ID, order.ID)
				return err
			}
		}
	}

	return nil
}

func (s *ServiceImpl) checkOrder(order models.Order) bool {
	if len(order.Errors) != 0 ||
		order.FailReason.Code != "" ||
		!order.IsDeferSubmission ||
		order.Status == "FAILED" ||
		order.Status == "CANCELLED_BY_DELIVERY_SERVICE" ||
		order.Status == "CANCELLED_BY_POS_SYSTEM" ||
		order.Status == "CANCELLED" {
		return false
	}

	return true
}

func (s *ServiceImpl) sendOrderReadyStatusToAgg(ctx context.Context, order models.Order, store storeModels.Store) error {
	aggService, err := s.aggregatorFactory.GetAggregator(order.DeliveryService, store)
	if err != nil {
		return err
	}

	if !s.validateOrderTime(order, store) {
		return nil
	}

	switch order.DeliveryService {
	case models.GLOVO.String(), models.YANDEX.String():
		if err := aggService.UpdateOrderInAggregator(ctx, order, store, models.READY_FOR_PICKUP.String()); err != nil {
			log.Err(err).Msgf("error while updating status in %s, error: %s , orderID %s", order.DeliveryService, err, order.ID)
		}
	case models.WOLT.String():
		if err := aggService.UpdateOrderInAggregator(ctx, order, store, models.Ready.String()); err != nil {
			log.Err(err).Msgf("error while updating status in %s, error: %s , orderID %s", order.DeliveryService, err, order.ID)
		}
	}

	if err := s.repository.UpdateOrderDeferStatus(ctx, false, order.ID); err != nil {
		log.Err(err).Msgf("error while updating order in database %s , orderID  %s", err, order.ID)
	}

	return nil
}

func (s *ServiceImpl) validateOrderTime(order models.Order, store storeModels.Store) bool {
	currentTime := time.Now().UTC()

	lunchStart := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 12, 0, 0, 0, currentTime.Location())
	lunchEnd := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 14, 0, 0, 0, currentTime.Location())

	dinnerStart := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 19, 0, 0, 0, currentTime.Location())
	dinnerEnd := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 21, 0, 0, 0, currentTime.Location())

	var sendStatusTime time.Time
	if currentTime.After(lunchStart) && currentTime.Before(lunchEnd) || currentTime.After(dinnerStart) && currentTime.Before(dinnerEnd) {
		sendStatusTime = order.OrderTime.Value.Add(time.Duration(store.DeferSubmission.BusyTime) * time.Minute)
		if sendStatusTime.After(currentTime) {
			return true
		}

		return false
	} else {
		sendStatusTime = order.OrderTime.Value.Add(time.Duration(store.DeferSubmission.DefaultTime) * time.Minute)

		if currentTime.After(sendStatusTime) {
			return true
		}

		return false
	}
}

func (s *ServiceImpl) GetOrderStat(ctx context.Context, interval string) (models.OrderStat, error) {
	start, end, err := s.getInterval(interval)
	if err != nil {
		return models.OrderStat{}, err
	}

	log.Info().Msgf("interval for order stat: %v - %v", start, end)

	timeoutErrSolutions, err := s.errSolution.GetTimeoutErrorSolutions(ctx)
	if err != nil {
		return models.OrderStat{}, err
	}

	totalOrderNum, err := s.repository.GetOrderNumber(ctx,
		selector.EmptyOrderSearch().
			SetOrderTimeFrom(start).
			SetOrderTimeTo(end),
	)
	if err != nil {
		return models.OrderStat{}, errors.Wrap(err, "GetOrderNumber")
	}

	failedOrders, _, err := s.repository.GetAllOrders(ctx,
		selector.EmptyOrderSearch().
			SetOrderTimeFrom(start).
			SetOrderTimeTo(end).
			SetStatusFailedOrFailReasonNotEmpty(),
	)
	if err != nil {
		return models.OrderStat{}, errors.Wrap(err, "GetAllOrders")
	}

	timeoutNum, err := s.repository.GetOrderNumber(ctx,
		selector.EmptyOrderSearch().
			SetOrderTimeFrom(start).
			SetOrderTimeTo(end).
			SetFailedReasonTimeoutCodes(timeoutErrSolutions),
	)

	if err != nil {
		return models.OrderStat{}, errors.Wrap(err, "GetOrderNumber (timeout one)")
	}

	orderStat := models.OrderStat{
		TotalOrderNumber: float64(totalOrderNum),
		TotalFailed:      float64(len(failedOrders)) - float64(timeoutNum),
		TimeoutErrs:      float64(timeoutNum),
		Errors:           make(map[string]float64),
	}

	for _, failedOrder := range failedOrders {
		if failedOrder.FailReason.Code == "" {
			orderStat.Errors[pos.OTHER_FAIL_REASON_CODE]++
			continue
		}
		orderStat.Errors[failedOrder.FailReason.Code]++
	}
	orderStat.ConstructedErrMsg, err = s.constructMsgForFailedOrders(ctx, failedOrders, orderStat)
	if err != nil {
		return models.OrderStat{}, err
	}

	return orderStat, nil
}

func (s *ServiceImpl) constructMsgForFailedOrders(ctx context.Context, orders []models.Order, orderStat models.OrderStat) (string, error) {

	var (
		msg, avoidable, notAvoidable  string
		avoidableNum, notAvoidableNum float64
		r                             = &telegram.Report{}
		mapAvoidable                  = make(map[string]float64)
		mapNotAvoidable               = make(map[string]float64)
	)

	r.Divisor = orderStat.TotalFailed
	totalErrorNumber := orderStat.TimeoutErrs + orderStat.TotalFailed

	msg += "<b>Информация по заказам за предыдущий день</b>\n"
	msg += fmt.Sprintf("<b>[✅] Общее количество заказов:</b> %.0f\n", orderStat.TotalOrderNumber)
	msg += fmt.Sprintf("<b>[❌] Количество неуспешных в кассу заказов вызванных неполадкой интернета/обмена с POS сервером/электроснабжением кассы: </b> %.0f - %.2f%% (%.2f%%)\n",
		orderStat.TimeoutErrs, s.getPercent(orderStat.TotalOrderNumber, orderStat.TimeoutErrs), s.getPercent(totalErrorNumber, orderStat.TimeoutErrs))
	msg += fmt.Sprintf("<b>[❌] Общее количество ошибок:</b> %.0f - %.2f%%  (%.2f%%)\n",
		orderStat.TotalFailed, s.getPercent(orderStat.TotalOrderNumber, orderStat.TotalFailed), s.getPercent(totalErrorNumber, orderStat.TotalFailed))

	if orderStat.TotalFailed == 0 {
		msg += "<b>Ошибочных заказов на последние 24 часа не было </b>"
		return msg, nil
	}

	errSolutions, err := s.errSolution.GetAllErrorSolutions(ctx)
	if err != nil {
		return "", err
	}

	avoidableNum, notAvoidableNum = s.GetAvoidableOrderNum(errSolutions, orderStat)

	for _, errorSolution := range errSolutions {
		if errorSolution.IsTimeout {
			continue
		}
		switch errorSolution.Avoidable {
		case true:
			r.Divisor = avoidableNum
			format := fmt.Sprintf("%s (код: %s)\n", errorSolution.BusinessName, errorSolution.Code)
			if orderStat.Errors[errorSolution.Code] > 0 {
				avoidableMsg, num := r.GetErrorMessage(format, orderStat.Errors[errorSolution.Code])
				mapAvoidable[avoidableMsg] = num
			}
		case false:
			r.Divisor = notAvoidableNum
			format := fmt.Sprintf("%s (код: %s)\n", errorSolution.BusinessName, errorSolution.Code)
			if orderStat.Errors[errorSolution.Code] > 0 {
				notAvoidableMsg, num := r.GetErrorMessage(format, orderStat.Errors[errorSolution.Code])
				mapNotAvoidable[notAvoidableMsg] = num
			}
		}
	}

	msg += s.formatStateErrorMessage("Ошибки, которые мы могли избежать", avoidableNum, orderStat.TotalFailed, orderStat.TotalOrderNumber)
	msg += s.formatStateErrorMessage("Ошибки, которые мы не могли избежать", notAvoidableNum, orderStat.TotalFailed, orderStat.TotalOrderNumber)

	avoidable = fmt.Sprintf("\n<b>Ошибки, которые мы могли избежать: </b> %.0f из них\n", avoidableNum) + s.mapSortByCount(mapAvoidable)
	notAvoidable = fmt.Sprintf("\n<b>Ошибки, которые мы не могли избежать: </b> %.0f из них\n", notAvoidableNum) + s.mapSortByCount(mapNotAvoidable)

	return msg + avoidable + notAvoidable, nil
}

func (s *ServiceImpl) GetAvoidableOrderNum(errSolutions []models2.ErrorSolution, orderStat models.OrderStat) (float64, float64) {
	var avoidableNum, notAvoidableNum float64
	for _, errorSolution := range errSolutions {
		if errorSolution.IsTimeout {
			continue
		}
		switch errorSolution.Avoidable {
		case true:
			if orderStat.Errors[errorSolution.Code] > 0 {
				avoidableNum += orderStat.Errors[errorSolution.Code]
			}
		case false:
			if orderStat.Errors[errorSolution.Code] > 0 {
				notAvoidableNum += orderStat.Errors[errorSolution.Code]
			}
		}
	}
	return avoidableNum, notAvoidableNum
}

func (s *ServiceImpl) mapSortByCount(errorsMap map[string]float64) string {
	type kv struct {
		Key   string
		Value float64
	}

	var (
		sortedValues []kv
		result       string
	)

	for k, v := range errorsMap {
		sortedValues = append(sortedValues, kv{Key: k, Value: v})
	}

	sort.Slice(sortedValues, func(i, j int) bool {
		return sortedValues[i].Value > sortedValues[j].Value
	})

	for _, kv := range sortedValues {
		result += kv.Key
	}

	return result
}

func (s *ServiceImpl) getPercent(dividend, divisor float64) float64 {
	if divisor == 0 {
		return 0
	}
	return (divisor * 100) / dividend
}

func (s *ServiceImpl) sortByPosAggErrorOrders(orders []models.Order, totalErrNum float64) string {

	var (
		iikoNum, rkeeper7XMLNum, rkeeperNum, otherNum, woltNum, glovoNum, yandexNum float64
		msg                                                                         string
	)

	for _, order := range orders {
		switch order.PosType {
		case models.IIKO.String():
			iikoNum++
		case models.RKeeper.String():
			rkeeperNum++
		case models.RKeeper7XML.String():
			rkeeper7XMLNum++
		default:
			otherNum++
		}

		switch order.DeliveryService {
		case models.WOLT.String():
			woltNum++
		case models.GLOVO.String():
			glovoNum++
		case models.YANDEX.String():
			yandexNum++
		}
	}

	msg = fmt.Sprintf("<b>Информация по агрегаторам:</b>\n")
	msg += fmt.Sprintf("<b> • wolt:</b> %.0f - %.2f%% <b> - </b>\n", woltNum, totalErrNum/woltNum)
	msg += fmt.Sprintf("<b> • glovo:</b> %.0f - %.2f%% <b> - </b>\n", glovoNum, totalErrNum/glovoNum)
	msg += fmt.Sprintf("<b> • yandex:</b> %.0f - %.2f%% <b> - </b>\n\n\n", yandexNum, totalErrNum/yandexNum)

	msg = fmt.Sprintf("<b>Информация по POS :</b>\n")
	msg += fmt.Sprintf("<b> • iiko:</b> %.0f - %.2f%% <b> - </b>\n", iikoNum, totalErrNum/iikoNum)
	msg += fmt.Sprintf("<b> • rkeeper:</b> %.0f - %.2f%% <b> - </b>\n", rkeeperNum, totalErrNum/rkeeperNum)
	msg += fmt.Sprintf("<b> • rkeeper_wsa:</b> %.0f - %.2f%% <b> - </b>\n", rkeeper7XMLNum, totalErrNum/rkeeper7XMLNum)
	msg += fmt.Sprintf("<b> • другие:</b> %.0f - %.2f%% <b> - </b>\n", otherNum, totalErrNum/otherNum)

	return msg
}

func (s *ServiceImpl) formatStateErrorMessage(name string, num, totalFailed, total float64) string {
	percentage := 0.0
	if num != 0 {
		percentage = (num * 100) / totalFailed
	}
	return fmt.Sprintf("      • %s: %.0f - %.2f%% (%.2f%%)\n", name, num, percentage, (num*100)/total)
}

func (s *ServiceImpl) updateOrderStatus(ctx context.Context, store storeModels.Store, order models.Order, status, statusDescription string) error {
	posService, err := s.posFactory.GetPosService(models.Pos(store.PosType), store)
	if err != nil {
		return err
	}

	newSystemStatus, err := posService.MapPosStatusToSystemStatus(status, order.Status)
	if err != nil {
		return err
	}

	if s.isIgnoreRepeatedOrderStatus(order.Status, newSystemStatus.String()) {
		log.Info().Msgf("ignoring repeated order status, order_id=%s, current system status=%s, new system status=%s", order.ID, order.Status, newSystemStatus.String())
		return nil
	}

	if err = s.repository.UpdateOrderStatusByID(ctx, order.ID, newSystemStatus.String()); err != nil {
		return err
	}

	err = s.publisher.NotifySubscribers(ctx, order, store, newSystemStatus)
	if err != nil {
		log.Err(err).Msgf("createOrder subscribers error")
	}

	ignoreStatusUpdate, err := s.ignoreStatusUpdateInAggregator(store, order, newSystemStatus)
	if err != nil {
		return err
	}

	if ignoreStatusUpdate {
		return nil
	}

	aggregatorService, err := s.aggregatorFactory.GetAggregator(order.DeliveryService, store)
	if err != nil {
		return err
	}

	aggStatus := aggregatorService.MapSystemStatusToAggregatorStatus(order, newSystemStatus, store)

	if err = aggregatorService.UpdateOrderInAggregator(ctx, order, store, aggStatus); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) isIgnoreRepeatedOrderStatus(currentSystemStatus, newSystemStatus string) bool {
	return currentSystemStatus == newSystemStatus
}

func (s *ServiceImpl) ignoreStatusUpdateInAggregator(store storeModels.Store, order models.Order, newSystemStatus models.PosStatus) (bool, error) {
	ignoreStatusUpdate, err := s.storeService.IgnoreStatusUpdate(store, order.DeliveryService)
	if err != nil {
		return false, err
	}

	if ignoreStatusUpdate && newSystemStatus != models.ACCEPTED {
		log.Info().Msgf("ignoring status update, order_id=%s, current system status=%s, new system status=%s", order.ID, order.Status, newSystemStatus.String())
		return true, nil
	}

	return false, nil
}

func (s *ServiceImpl) UpdateOrderStatus(ctx context.Context, orderID, status, statusDescription string) error {
	order, err := s.repository.FindOrderByPosOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	store, err := s.storeService.GetByID(ctx, order.RestaurantID)
	if err != nil {
		return err
	}

	if err = s.updateOrderStatus(ctx, store, order, status, statusDescription); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) GetShaurmaFoodOrders(ctx context.Context, restaurantID, deliveryService, orderStatus, dateFrom, dateTo string, page, limit int64) ([]models.ShaurmaFoodOrdersInfo, int, error) {

	shaurmaFoodRestiIDs, err := s.getShaurmaFoodRestaurantIDs(ctx)
	if err != nil {
		return nil, 0, err
	}
	restiIDs := make([]string, 0, 1)

	if restaurantID == "" {
		restiIDs = shaurmaFoodRestiIDs
	} else {
		if !s.isShaurmaFood(restaurantID, shaurmaFoodRestiIDs) {
			return nil, 0, errors.New("restaurant id is not for ShaurmaFood")
		}
		restiIDs = append(restiIDs, restaurantID)
	}

	query := selector.EmptyOrderSearch().
		SetRestaurants(restiIDs).
		SetDeliveryService(deliveryService).
		SetStatus(orderStatus).
		SetLimit(limit).
		SetPage(page)

	if dateFrom != "" {
		date, err := time.Parse("2006-01-02T15:04:05-0700", dateFrom)
		if err != nil {
			return nil, 0, err
		}
		query = query.SetOrderTimeFrom(date)
	}
	if dateTo != "" {
		date, err := time.Parse("2006-01-02T15:04:05-0700", dateTo)
		if err != nil {
			return nil, 0, err
		}
		query = query.SetOrderTimeTo(date)
	}

	orders, total, err := s.repository.GetAllOrders(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	ordersInfo := getCertainDataFromOrder(orders)

	pagesCount := math.Ceil(float64(total) / float64(limit))
	return ordersInfo, int(pagesCount), nil
}

func (s *ServiceImpl) GetOrder(ctx context.Context, orderID string) (models.Order, error) {
	order, err := s.repository.FindOrderByOrderID(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (s *ServiceImpl) GetOrderByID(ctx context.Context, id string) (models.Order, error) {
	order, err := s.repository.FindOrderByID(ctx, id)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func getCertainDataFromOrder(orders []models.Order) []models.ShaurmaFoodOrdersInfo {

	ordersInfo := make([]models.ShaurmaFoodOrdersInfo, 0, len(orders))
	for _, order := range orders {
		ordersInfo = append(ordersInfo, models.ShaurmaFoodOrdersInfo{
			RestaurantID:    order.RestaurantID,
			PosOrderID:      order.PosOrderID,
			OrderCode:       order.OrderCode,
			DeliveryService: order.DeliveryService,
			Products:        order.Products,
			PosPaymentInfo:  order.PosPaymentInfo,
			StatusesHistory: order.StatusesHistory,
			CreationResult:  order.CreationResult,
			Errors:          order.Errors,
		})
	}
	return ordersInfo
}

func (s *ServiceImpl) CreateOrder(ctx context.Context, externalStoreID, deliveryService string, aggReq interface{}, storeSecret string) (models.Order, error) {
	st, err := s.storeService.GetByExternalIdAndDeliveryService(ctx, externalStoreID, deliveryService)
	if err != nil {
		return models.Order{}, err
	}

	utils.Beautify("getting store body", st)

	if deliveryService == models.STARTERAPP.String() {
		aggReq, err = s.addExtIDAndNamesAndCookingTime(aggReq, st)
		if err != nil {
			return models.Order{}, err
		}
	}

	isSecretValid, err := s.storeService.IsSecretValid(st, deliveryService, storeSecret)
	if err != nil {
		return models.Order{}, err
	}

	if !isSecretValid {
		return models.Order{}, errors.New("store secret is not valid")
	}

	agg, err := s.aggregatorFactory.GetAggregator(deliveryService, st)
	if err != nil {
		return models.Order{}, err
	}

	if st.Settings.HasVirtualStore {
		return s.splitVirtualStoreOrder(ctx, st, aggReq, agg, deliveryService)
	}

	aggregatorRequestBody, err := utils.GetJsonFormatFromModel(aggReq)
	if err != nil {
		return models.Order{}, err
	}

	req, err := agg.GetSystemCreateOrderRequestByAggregatorRequest(aggReq, st)
	if err != nil {
		return models.Order{}, err
	}

	utils.Beautify("aggregator request to system request", req)

	req = fillRequestData(req, st, aggregatorRequestBody)

	req, err = s.saveProductImages(ctx, req, st, deliveryService)
	if err != nil {
		return models.Order{}, err
	}

	req = s.setRestaurantCharge(ctx, req, st)

	req = s.setDeferSubmission(req, st)

	switch req.DeliveryService {
	case models.QRMENU.String():
		req.CookingTime = st.QRMenu.CookingTime

		cart, err := s.cartService.GetQRMenuCartByID(ctx, strings.Split(req.OrderID, "_")[0])
		if err != nil {
			return req, err
		}
		req.PaymentSystem = cart.PaymentSystem
	case models.KWAAKA_ADMIN.String():
		req.CookingTime = st.QRMenu.CookingTime

		cart, err := s.cartService.GetKwaakaAdminCartByOrderID(ctx, strings.Split(req.OrderID, "_")[0])
		if err != nil {
			return req, err
		}
		req.PaymentSystem = cart.PaymentType
	}

	errSolutions, err := s.errSolution.GetAllErrorSolutions(ctx)
	if err != nil {
		return models.Order{}, err
	}

	order, err := s.saveOrderToDb(ctx, req)

	//TODO delete
	skipErrors := false
	if deliveryService == models.QRMENU.String() {
		skipErrors = true
	}
	if err != nil && !skipErrors {
		order.FailReason, _, err = s.errSolution.SetFailReason(ctx, st, err.Error(), pos.MatchingCodes(err.Error(), errSolutions), "")
		if err != nil {
			return order, err
		}
		updateOrderErr := s.repository.UpdateOrder(ctx, order)
		if updateOrderErr != nil {
			return order, updateOrderErr
		}

		return order, err
	}

	positionsOnStopError := s.getItemsAvailableStatus(ctx, order, st.MenuID)
	if positionsOnStopError != nil {
		log.Err(err).Msgf("update order: %s with positions on stop object error", order.OrderID)
	}

	isSendToPos, err := s.storeService.IsSendToPos(st, deliveryService)
	if err != nil {
		return order, err
	}

	if !isSendToPos {
		log.Info().Msgf("order_id = %s, order_code = %s, send_to_pos = %t, order = %+v", order.ID, order.OrderCode, isSendToPos, order)

		order.FailReason, _, err = s.errSolution.SetFailReason(ctx, st, "", "", pos.INTEGRATION_OFF_CODE)
		if err != nil {
			return order, err
		}

		err := s.repository.UpdateOrder(ctx, order)
		if err != nil {
			return order, err
		}

		if err = s.repository.UpdateOrderStatusByID(ctx, order.ID, string(models.STATUS_SKIPPED)); err != nil {
			return order, err
		}

		return order, nil
	}

	order.OrderCodePrefix, err = s.storeService.GetOrderCodePrefix(ctx, st, deliveryService)
	if err != nil {
		return order, err
	}

	if order.Type == "PREORDER" {
		return order, s.waitSendingOrder(ctx, order, st)
	}

	order.IsMarketplace, err = agg.IsMarketPlace(order.RestaurantSelfDelivery, st)
	if err != nil {
		return order, err
	}

	payments, err := s.storeService.GetPaymentTypes(st, deliveryService, order.PosPaymentInfo)
	if err != nil {
		return models.Order{}, err
	}

	order.PosPaymentInfo = s.getPayments(order.PaymentMethod, payments)

	isAutoAcceptOn, err := s.storeService.IsAutoAccept(st, deliveryService)
	isPostAutoAcceptOn, err := s.storeService.IsPostAutoAccept(st, deliveryService)
	if err != nil {
		return order, err
	}

	if isAutoAcceptOn {
		order.Status = models.ACCEPTED.String()

		aggStatus := agg.MapSystemStatusToAggregatorStatus(order, models.ACCEPTED, st)

		if err = agg.UpdateOrderInAggregator(ctx, order, st, aggStatus); err != nil {
			return models.Order{}, err
		}

		if err = s.repository.UpdateOrderStatusByID(ctx, order.ID, models.ACCEPTED.String()); err != nil {
			return models.Order{}, err
		}
		if deliveryService == models.YANDEX.String() {
			time.Sleep(7 * time.Second)
		}
	}

	order, err = s.orderRuleService.UseTheOrderRules(ctx, order)
	if err != nil {
		return models.Order{}, err
	}

	if order.Status == models.STATUS_CANCELLED_BY_DELIVERY_SERVICE.ToString() {
		return order, fmt.Errorf("order canceled by delivery service, order id: %s, status: %s", order.OrderID, order.Status)
	}

	order, err = s.posSender.SendPosRequest(ctx, order, st, deliveryService)
	if err != nil {

		if aggSendNotifErr := agg.SendOrderErrorNotification(ctx, order); aggSendNotifErr != nil {
			log.Err(err).Msgf("%+v", aggSendNotifErr)
		}

		if err2 := s.repository.UpdateOrder(ctx, order); err2 != nil {
			return order, err2
		}

		if errors.Is(err, pos.ErrRetry) {
			return order, nil
		}

		return order, err
	}
	if isPostAutoAcceptOn {
		order.Status = models.ACCEPTED.String()

		aggStatus := agg.MapSystemStatusToAggregatorStatus(order, models.ACCEPTED, st)

		if err = agg.UpdateOrderInAggregator(ctx, order, st, aggStatus); err != nil {
			return models.Order{}, err
		}

		if err = s.repository.UpdateOrderStatusByID(ctx, order.ID, models.ACCEPTED.String()); err != nil {
			return models.Order{}, err
		}
		if deliveryService == models.YANDEX.String() {
			time.Sleep(7 * time.Second)
		}
	}
	order, err = s.successOrder(ctx, order)
	if err != nil {
		return order, nil
	}

	return order, nil
}

func fillRequestData(req models.Order, st storeModels.Store, aggregatorRequestBody string) models.Order {
	req.RestaurantID = st.ID
	req.PosType = st.PosType
	req.RestaurantName = st.Name

	var logStream models.LogStream
	req.LogLinks = models.LogLinks{
		LogStreamLink:          logStream.GetLink(),
		LogStreamLinkByOrderId: logStream.GetLinkWithPattern(req.OrderID),
	}

	if req.LogMessages.FromDelivery == "" {
		req.LogMessages.FromDelivery = aggregatorRequestBody
	}

	return req
}

func (s *ServiceImpl) saveOrderToDb(ctx context.Context, req models.Order) (models.Order, error) {
	order, err := s.repository.InsertOrder(ctx, req)
	if err != nil {
		log.Err(err).Msgf("orders core, insert order error")

		if errors.Is(err, errs.ErrAlreadyExist) {
			log.Info().Msg("Order already exist, skipping...")
			req.FailReason = models.FailReason{
				Code:    pos.ORDER_ALREADY_EXIST_CODE,
				Message: pos.ORDER_ALREADY_EXIST,
			}
			if err1 := s.repository.UpdateOrder(ctx, req); err1 != nil {
				return models.Order{}, err
			}
			return req, errors.Wrap(validator.ErrPassed, fmt.Sprintf("order %s passed", req.OrderID))
		}

		return s.failOrder(ctx, order, err.Error())
	}

	return order, nil
}

func (s *ServiceImpl) saveProductImages(ctx context.Context, order models.Order, store storeModels.Store, deliveryService string) (models.Order, error) {
	aggrMenu, err := s.menuService.GetAggregatorMenuIfExists(ctx, store, deliveryService)
	if err != nil {
		s.failOrder(ctx, order, err.Error())
	}

	aggrProductMap := make(map[string]coreMenuModels.Product)
	for i := range aggrMenu.Products {
		aggrProductMap[aggrMenu.Products[i].ExtID] = aggrMenu.Products[i]
	}

	for i := range order.Products {
		if _, ok := aggrProductMap[order.Products[i].ID]; ok {
			order.Products[i].ImageURLs = aggrProductMap[order.Products[i].ID].ImageURLs
		}
	}

	return order, nil
}

func (s *ServiceImpl) failOrder(ctx context.Context, req models.Order, errMessage string) (models.Order, error) {

	log.Trace().Msgf("%s: %v", errMessage, validator.ErrFailed)

	req.Status = string(models.STATUS_FAILED)

	req.StatusesHistory = append(req.StatusesHistory, models.OrderStatusUpdate{
		Name: string(models.STATUS_FAILED),
		Time: models.TimeNow().Time,
	})

	errSolutions, err := s.errSolution.GetAllErrorSolutions(ctx)
	if err != nil {
		return models.Order{}, err
	}

	req.FailReason, _, err = s.errSolution.SetFailReason(ctx, storeModels.Store{}, errMessage, pos.MatchingCodes(errMessage, errSolutions), "")

	if updateErr := s.repository.UpdateOrder(ctx, req); updateErr != nil {
		return req, errors.Wrap(errWithNotification, validator.ErrFailed.Error())
	}

	log.Info().Msgf("updated fail order_id:%s with fail_reason: %+v", req.OrderID, req.FailReason.Message)

	if err != nil {
		log.Trace().Err(err).Msg("Error while saving order")
		return req, errors.Wrap(errWithNotification, err.Error())
	}

	return req, errors.Wrap(errWithNotification, validator.ErrFailed.Error())
}

func (s *ServiceImpl) successOrder(ctx context.Context, req models.Order) (models.Order, error) {

	log.Info().Msgf("Order created successfully, status: %s, order id: %s", req.Status, req.ID)

	req.StatusesHistory = append(req.StatusesHistory, models.OrderStatusUpdate{
		Name: req.Status,
		Time: models.TimeNow().Time,
	})

	// ignoring a possible error during the order update
	if err := s.repository.UpdateOrder(ctx, req); err != nil {
		return req, err
	}

	return req, nil
}

func (s *ServiceImpl) waitSendingOrder(ctx context.Context, order models.Order, store storeModels.Store) error {

	log.Info().Msgf("Order waiting sending, status: %v", string(models.STATUS_WAIT_SENDING))

	systemStatus := models.WAIT_SENDING

	if err := s.repository.UpdateOrderStatusByID(ctx, order.ID, systemStatus.String()); err != nil {
		return err
	}

	if !order.IsChildOrder {
		agg, err := s.aggregatorFactory.GetAggregator(order.DeliveryService, store)
		if err != nil {
			return err
		}

		aggStatus := agg.MapSystemStatusToAggregatorStatus(order, systemStatus, store)

		if err = agg.UpdateOrderInAggregator(ctx, order, store, aggStatus); err != nil {
			return err
		}
	}

	return nil
}

func (s *ServiceImpl) getPayments(paymentMethod string, paymentTypes storeModels.DeliveryServicePaymentType) models.PosPaymentInfo {

	switch paymentMethod {
	case models.PAYMENT_METHOD_CASH:
		orderType := paymentTypes.CASH.OrderType
		if paymentTypes.CASH.OrderTypeForVirtualStore != "" {
			orderType = paymentTypes.CASH.OrderTypeForVirtualStore
		}

		return models.PosPaymentInfo{
			PaymentTypeID:          paymentTypes.CASH.PaymentTypeID,
			OrderType:              orderType,
			OrderTypeService:       paymentTypes.CASH.OrderTypeService,
			IsProcessedExternally:  paymentTypes.CASH.IsProcessedExternally,
			PaymentTypeKind:        paymentTypes.CASH.PaymentTypeKind,
			PromotionPaymentTypeID: paymentTypes.CASH.PromotionPaymentTypeID,
		}
	case models.PAYMENT_METHOD_DELAYED:
		orderType := paymentTypes.DELAYED.OrderType
		if paymentTypes.DELAYED.OrderTypeForVirtualStore != "" {
			orderType = paymentTypes.DELAYED.OrderTypeForVirtualStore
		}

		return models.PosPaymentInfo{
			PaymentTypeID:          paymentTypes.DELAYED.PaymentTypeID,
			OrderType:              orderType,
			OrderTypeService:       paymentTypes.DELAYED.OrderTypeService,
			IsProcessedExternally:  paymentTypes.DELAYED.IsProcessedExternally,
			PaymentTypeKind:        paymentTypes.DELAYED.PaymentTypeKind,
			PromotionPaymentTypeID: paymentTypes.DELAYED.PromotionPaymentTypeID,
		}
	}

	return models.PosPaymentInfo{}
}

func (s *ServiceImpl) getInterval(interval string) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	switch interval {
	case "day":
		start := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC)
		end := start.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		return start, end, nil
	default:
		return time.Time{}, time.Time{}, ErrInvalidInterval
	}
}

func (s *ServiceImpl) getShaurmaFoodRestaurantIDs(ctx context.Context) ([]string, error) {
	shaurmaFoodGroupIDs := []string{"642ab068d5ad369ab4647d44", "6604030be5cb51f5698fafd5"}

	shaurmaFoodRestaurantIDs := make([]string, 0, 1)

	for _, id := range shaurmaFoodGroupIDs {
		group, err := s.storeGroupService.GetStoreGroupByID(ctx, id)
		if err != nil {
			return nil, err
		}
		shaurmaFoodRestaurantIDs = append(shaurmaFoodRestaurantIDs, group.StoreIds...)
	}

	return shaurmaFoodRestaurantIDs, nil
}

func (s *ServiceImpl) isShaurmaFood(restiID string, shaurmaFoodRestiIDs []string) bool {
	for _, shaurmaFoodRestiID := range shaurmaFoodRestiIDs {
		if restiID == shaurmaFoodRestiID {
			return true
		}
	}
	return false
}

func (s *ServiceImpl) CancelOrderByAggregator(ctx context.Context, orderID string, delivery string) error {
	order, err := s.repository.FindOrderByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	if delivery != order.DeliveryService {
		return errors.New("You cannot cancel order of this aggregator")
	}

	st, err := s.storeService.GetByID(ctx, order.StoreID)
	if err != nil {
		return err
	}

	posService, err := s.posFactory.GetPosService(models.Pos(order.PosType), st)
	if err != nil {
		return err
	}

	if err = posService.CancelOrder(ctx, order, st); err != nil {
		return err
	}

	if err = s.updateOrderStatus(ctx, st, order, string(models.STATUS_CANCELLED_BY_DELIVERY_SERVICE), ""); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) Get3plOrdersWithoutDriver(ctx context.Context, indriveCallTime int64) ([]models.Order, error) {

	iikoOrders, err := s.repository.GetIIKO3plOrdersForCron(ctx, indriveCallTime)
	if err != nil {
		log.Info().Err(err)
		return nil, err
	}

	orderForCron, err := s.repository.Get3plOrdersForCron(ctx, indriveCallTime)
	if err != nil {
		log.Info().Err(err)
		return nil, err
	}

	iikoOrders = append(iikoOrders, orderForCron...)

	uniqueOrders := make(map[string]models.Order)
	for _, order := range iikoOrders {
		uniqueOrders[order.ID] = order
	}

	var unique []models.Order
	for _, order := range uniqueOrders {
		log.Info().Msgf("preorder with id %s was selected as unique", order.ID)
		unique = append(unique, order)
	}

	return unique, nil
}

func (s *ServiceImpl) setRestaurantCharge(ctx context.Context, order models.Order, store storeModels.Store) models.Order {
	if store.RestaurantCharge.IsRestaurantChargeOn == false {
		return order
	}

	restaurantCharge := order.TotalCustomerToPay.Value * store.RestaurantCharge.RestaurantChargePercent / 100

	if restaurantCharge < store.RestaurantCharge.MinRestaurantCharge {
		restaurantCharge = store.RestaurantCharge.MinRestaurantCharge
	}

	if restaurantCharge > store.RestaurantCharge.MaxRestaurantCharge {
		restaurantCharge = store.RestaurantCharge.MaxRestaurantCharge
	}

	order.RestaurantCharge.Value = restaurantCharge

	return order
}

func (s *ServiceImpl) setDeferSubmission(order models.Order, store storeModels.Store) models.Order {
	if store.DeferSubmission.IsDeferSumbission {
		switch order.DeliveryService {
		case models.WOLT.String(), models.YANDEX.String(), models.GLOVO.String():
			order.IsDeferSubmission = true
		}
	}

	return order
}

func (s *ServiceImpl) getItemsAvailableStatus(ctx context.Context, order models.Order, posMenuID string) error {

	menu, err := s.menuService.GetMenuById(ctx, posMenuID)
	if err != nil {
		return err
	}

	var positionsOnStop []models.PositionsOnStop

	switch order.PosType {
	case models.Paloma.String():
		positionsOnStop = s.getPalomaProductsStatus(menu.Products, order.Products)
	default:
		positionsOnStop = s.getPalomaProductsStatus(menu.Products, order.Products)
	}

	for _, item := range order.Products {
		positionsOnStop = append(positionsOnStop, s.getAttributeStatus(menu.Attributes, item.Attributes)...)
	}

	order.PositionsOnStop = positionsOnStop

	return s.repository.UpdateOrder(ctx, order)
}

func (s *ServiceImpl) getProductStatus(menuProducts []coreMenuModels.Product, orderProducts []models.OrderProduct) []models.PositionsOnStop {
	menuProductMap := make(map[string]bool, len(menuProducts))
	result := make([]models.PositionsOnStop, 0)
	for _, menuProduct := range menuProducts {
		menuProductMap[menuProduct.ExtID] = menuProduct.IsAvailable
	}
	for _, orderProduct := range orderProducts {
		available, exist := menuProductMap[orderProduct.ID]
		if !exist {
			continue
		}
		if !available {
			result = append(result, models.PositionsOnStop{
				ID:   orderProduct.ID,
				Name: orderProduct.Name,
			})
		}
	}
	return result
}

func (s *ServiceImpl) getPalomaProductsStatus(menuProducts []coreMenuModels.Product, orderProducts []models.OrderProduct) []models.PositionsOnStop {
	menuProductMap := make(map[string]bool, len(menuProducts))
	result := make([]models.PositionsOnStop, 0)
	for _, menuProduct := range menuProducts {
		menuProductMap[menuProduct.ProductID] = menuProduct.IsAvailable
	}
	for _, orderProduct := range orderProducts {
		available, exist := menuProductMap[orderProduct.ID]
		if !exist {
			continue
		}
		if !available {
			result = append(result, models.PositionsOnStop{
				ID:   orderProduct.ID,
				Name: orderProduct.Name,
			})
		}
	}
	return result
}

func (s *ServiceImpl) getAttributeStatus(menuAttributes []coreMenuModels.Attribute, orderAttributes []models.ProductAttribute) []models.PositionsOnStop {
	menuAttributeMap := make(map[string]bool, len(menuAttributes))
	result := make([]models.PositionsOnStop, 0)
	for _, menuAttribute := range menuAttributes {
		menuAttributeMap[menuAttribute.ExtID] = menuAttribute.IsAvailable
	}
	for _, orderAttribute := range orderAttributes {
		available, exist := menuAttributeMap[orderAttribute.ID]
		if !exist {
			continue
		}
		if !available {
			result = append(result, models.PositionsOnStop{
				ID:   orderAttribute.ID,
				Name: orderAttribute.Name,
			})
		}
	}
	return result
}
