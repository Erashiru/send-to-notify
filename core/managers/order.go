package managers

import (
	"context"
	"fmt"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	notifierpkg "github.com/kwaaka-team/orders-core/core/managers/notifier"
	telegram2 "github.com/kwaaka-team/orders-core/core/managers/notifier/telegram"
	"github.com/kwaaka-team/orders-core/core/managers/notifier/whatsapp"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
	storeDto "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/kwaaka-team/orders-core/pkg/wolt/clients/http"
	"github.com/kwaaka-team/orders-core/service/kwaaka_3pl"
	models3 "github.com/kwaaka-team/orders-core/service/kwaaka_3pl/models"
	"strconv"
	"strings"

	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"

	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"time"

	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	"github.com/kwaaka-team/orders-core/core/managers/aggregator"
	"github.com/kwaaka-team/orders-core/core/managers/pos"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/custom"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	storeSelector "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	MenuClient "github.com/kwaaka-team/orders-core/pkg/menu"
	MenuModels "github.com/kwaaka-team/orders-core/pkg/menu/dto"

	"github.com/kwaaka-team/orders-core/pkg/que"
	StoreClient "github.com/kwaaka-team/orders-core/pkg/store"

	"github.com/rs/zerolog/log"

	aggregatorService "github.com/kwaaka-team/orders-core/service/aggregator"
	orderService "github.com/kwaaka-team/orders-core/service/order"
	posService "github.com/kwaaka-team/orders-core/service/pos"
)

type Order interface {
	GetOrder(ctx context.Context, query selector.Order) (models.Order, error)
	GetActiveOrders(ctx context.Context, query selector.Order) ([]models.Order, error)
	GetActivePreorders(ctx context.Context, query selector.Order) ([]models.Order, error)
	UpdateOrder(ctx context.Context, order models.Order) error
	CancelOrder(ctx context.Context, order models.CancelOrder) error
	UpdateOrderStatus(ctx context.Context, orderID, pos, status, errorDescription string) error
	UpdateOrderStatusByID(ctx context.Context, id, pos, status string) error
	UpdateOrderStatusInDS(ctx context.Context, id string, status models.PosStatus) error
	GetOrdersWithFilters(ctx context.Context, query selector.Order) ([]models.Order, int, error)
	CreateOrderDB(ctx context.Context, req models.Order) (models.Order, error)
	CreateOrderInPOS(ctx context.Context, req models.Order) (models.Order, error)
	SetPaidStatus(ctx context.Context, orderID string) error
	ManualUpdateStatus(ctx context.Context, order models.Order) error
	GetAllOrders(ctx context.Context, query selector.Order) ([]models.Order, error)
	CancelOrderInPos(ctx context.Context, req models.CancelOrderInPos) error
	GetDirectCallCenterActivePreorders(ctx context.Context, query selector.Order) ([]models.Order, error)
	GetHaniKarimaActivePreOrders(ctx context.Context, query selector.Order) ([]models.Order, error)
	GetOrdersForAutoCloseCron(ctx context.Context, query selector.Order) ([]models.Order, error)
}

type OrderManager struct {
	ds                drivers.DataStore
	orderRepo         drivers.OrderRepository
	globalConfig      config.Configuration
	sqsCli            que.SQSInterface
	menuClient        MenuClient.Client
	storeClient       StoreClient.Client
	storeService      store.Service
	aggregatorFactory aggregatorService.Factory
	posFactory        posService.Factory
	repository        orderService.Repository
	kwaaka3plService  kwaaka_3pl.Service
	telegramService   orderService.TelegramServiceImpl
}

var nonCancelableStatuses = map[string]struct{}{
	models3.ComingToPickup: {},
	models3.PickedUp:       {},
	models3.Delivered:      {},
	models3.Returning:      {},
	models3.Returned:       {},
	models3.Failed:         {},
	models3.Cancelled:      {},
}

var errConstructor error = errors.New("OrderManager constructor error ")

func NewOrderManager(
	ds drivers.DataStore,
	globalConfig config.Configuration,
	menuClient MenuClient.Client,
	storeClient StoreClient.Client,
	sqsCli que.SQSInterface,
	storeService store.Service,
	aggregatorFactory aggregatorService.Factory,
	posFactory posService.Factory,
	orderRepo orderService.Repository,
	kwaaka3plService kwaaka_3pl.Service,
	telegramService orderService.TelegramServiceImpl,
) (Order, error) {
	if storeService == nil {
		return nil, errors.Wrap(errConstructor, "storeFactory is nil")
	}
	if aggregatorFactory == nil {
		return nil, errors.Wrap(errConstructor, "aggregatorFactory is nil")
	}
	if posFactory == nil {
		return nil, errors.Wrap(errConstructor, "posFactory is nil")
	}
	if orderRepo == nil {
		return nil, errors.Wrap(errConstructor, "orderRepo is nil")
	}

	return &OrderManager{
		ds:           ds,
		orderRepo:    ds.OrderRepository(),
		globalConfig: globalConfig,
		sqsCli:       sqsCli,
		menuClient:   menuClient,
		storeClient:  storeClient,

		storeService:      storeService,
		aggregatorFactory: aggregatorFactory,
		posFactory:        posFactory,
		repository:        orderRepo,
		kwaaka3plService:  kwaaka3plService,
		telegramService:   telegramService,
	}, nil
}

func (manager OrderManager) GetAllOrders(ctx context.Context, query selector.Order) ([]models.Order, error) {
	orders, err := manager.orderRepo.GetAllOrders(ctx, query)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (manager OrderManager) ManualUpdateStatus(ctx context.Context, order models.Order) error {
	store, err := manager.storeClient.FindStore(ctx, storeDto.StoreSelector{
		ID:      order.RestaurantID,
		PosType: order.PosType,
	})
	if err != nil {
		return err
	}

	posManager, err := pos.NewPosManager(order.PosType, manager.globalConfig, manager.ds, store, coreMenuModels.Menu{}, coreMenuModels.Promo{}, manager.menuClient, coreMenuModels.Menu{}, manager.storeClient)
	if err != nil {
		log.Trace().Err(validator.ErrPosInitialize).Msgf("posName: %s, restaurantID: %s", order.PosType, store.ID)
		return validator.ErrPosInitialize
	}
	if posManager == nil {
		return validator.ErrPosInitialize
	}

	externalStatus, err := posManager.GetOrderStatus(ctx, order, store)
	if err != nil {
		return err
	}

	posStatus, err := manager.compareStatuses(order.PosType, externalStatus, order.Status)
	if err != nil {
		return err
	}

	if err := manager.ds.OrderRepository().UpdateOrderStatus(ctx, selector.EmptyOrderSearch().SetID(order.ID), posStatus.String(), ""); err != nil {
		return err
	}

	if posStatus == models.FAILED {
		log.Info().Msg("order was failed, skipped aggregator update")
		if queErr := manager.sendMessageToQue(ctx, telegram.UpdateOrder, order, store, externalStatus); queErr != nil {
			log.Err(validator.ErrSendingToQue).Msg("")
		}

		return nil
	}

	aggregatorStatus, _ := manager.compareAggregatorStatuses(order, posStatus, store)

	if err := manager.updateOrderInAggregator(ctx, order, store, aggregatorStatus); err != nil {
		return err
	}

	return nil
}

func (manager OrderManager) sendOrder(ctx context.Context, order models.Order, store coreStoreModels.Store, menu coreMenuModels.Menu, promo coreMenuModels.Promo, aggregatorMenu coreMenuModels.Menu) (models.Order, error) {

	posManager, err := pos.NewPosManager(store.PosType, manager.globalConfig, manager.ds, store, menu, promo, manager.menuClient, aggregatorMenu, manager.storeClient)
	if err != nil {
		log.Trace().Err(validator.ErrPosInitialize).Msgf("posName: %s, restaurantID: %s", store.PosType, store.ID)
		return order, validator.ErrPosInitialize
	}

	if posManager == nil {
		log.Trace().Err(validator.ErrPosInitialize).Msg("")
		return order, validator.ErrPosInitialize
	}

	utils.Beautify("successfully creating pos manager", posManager)

	order, err = posManager.CreateOrder(ctx, order, store)
	if err != nil {
		log.Trace().Err(err).Msg("cant send order to POS")
		return order, err
	}

	return order, nil
}

func (manager OrderManager) CreateOrderDB(ctx context.Context, req models.Order) (models.Order, error) {

	store, err := manager.storeClient.FindStore(ctx, storeDto.StoreSelector{
		ID:              req.RestaurantID,
		DeliveryService: req.DeliveryService,
		ExternalStoreID: req.StoreID,
		PosType:         req.PosType,
	})
	if err != nil {
		log.Err(err).Msgf("Find store error")
		return req, err
	}

	order, err := manager.orderRepo.InsertOrder(ctx, req)

	if err != nil {

		if errors.Is(err, errs.ErrAlreadyExist) {
			log.Err(err).Msgf("Insert order error")
			return manager.passOrder(order)
		}

		if queErr := manager.sendMessageToQue(ctx, telegram.CreateOrder, req, store, err.Error()); queErr != nil {
			log.Err(validator.ErrSendingToQue).Msg("")
		}

		log.Err(err).Msgf("Insert order error")

		return manager.failOrder(ctx, order, err)
	}

	return order, nil
}

func (manager OrderManager) CreateOrderInPOS(ctx context.Context, req models.Order) (models.Order, error) {
	store, err := manager.storeClient.FindStore(ctx, storeDto.StoreSelector{
		ID:              req.RestaurantID,
		ExternalStoreID: req.StoreID,
		PosType:         req.PosType,
	})

	if err != nil {
		msg := err.Error() + " store"
		if queErr := manager.sendMessageToQue(ctx, telegram.CreateOrder, req, store, msg); queErr != nil {
			log.Err(validator.ErrSendingToQue).Msg("")
		}
		log.Err(err).Msgf("Find store error")
		return req, err
	}

	order, err := manager.GetOrder(ctx, selector.EmptyOrderSearch().SetID(req.ID))

	if err != nil {
		msg := err.Error() + " GetOrderDB"
		if queErr := manager.sendMessageToQue(ctx, telegram.CreateOrder, req, store, msg); queErr != nil {
			log.Err(validator.ErrSendingToQue).Msg("")
		}
		log.Err(err).Msgf("GetOrder error")
		return req, err
	}
	order.IsRetry = req.IsRetry

	var (
		promo coreMenuModels.Promo
	)

	var (
		sendToPos    bool
		autoAcceptOn bool
	)

	switch order.DeliveryService {
	case "qr_menu":
		sendToPos = store.QRMenu.SendToPos
	case "glovo":
		sendToPos = store.Glovo.SendToPos
	case "wolt":
		sendToPos = store.Wolt.SendToPos
	case "express24":
		sendToPos = store.Express24.SendToPos
	case "deliveroo":
		sendToPos = store.Deliveroo.SendToPos
	case "kwaaka_admin":
		sendToPos = store.KwaakaAdmin.SendToPos
	case "starter_app":
		sendToPos = store.StarterApp.SendToPos
	default:
		for _, deliveryService := range store.ExternalConfig {
			if deliveryService.Type == order.DeliveryService {
				sendToPos = deliveryService.SendToPos
				promo, _ = manager.menuClient.GetPromos(ctx, req.StoreID, MenuModels.DeliveryService(order.DeliveryService))
				autoAcceptOn = deliveryService.AutoAcceptOn
			}
		}
	}

	if !sendToPos {
		return manager.skipOrder(ctx, order, nil)
	}

	if order.Status == models.STATUS_CANCELLED_BY_DELIVERY_SERVICE.ToString() {
		log.Info().Msgf("order canceled by delivery service, order id: %s, status: %s", order.OrderID, order.Status)
		return models.Order{}, fmt.Errorf("send preorder in POS error, order canceled by delivery service. order id: %s, status: %s", order.OrderID, order.Status)
	}

	menu, err := manager.menuClient.GetMenuByID(ctx, store.MenuID)
	if err != nil {
		msg := err.Error() + " GetMenuByID"
		if queErr := manager.sendMessageToQue(ctx, telegram.CreateOrder, req, store, msg); queErr != nil {
			log.Err(validator.ErrSendingToQue).Msg("")
		}
		return manager.failOrder(ctx, order, err)
	}
	aggregatorMenu, err := manager.getActiveMenu(ctx, store, order)
	if err != nil {
		log.Err(err).Msgf("restaurant name %s, active menu for %s is wrong", store.Name, order.DeliveryService)
		if queErr := manager.sendMessageToQue(ctx, telegram.CreateOrder, req, store, err.Error()); queErr != nil {
			log.Err(validator.ErrSendingToQue).Msg("")
		}

		return manager.failOrder(ctx, order, err)
	}

	if isHaniRestDelivery(order.StoreID, order.DeliveryService) {
		order = haniRestAddDeliveryProduct(order, menu)
	}

	order, err = manager.sendOrder(ctx, order, store, menu, promo, aggregatorMenu)
	if err != nil {
		if errors.Is(err, posService.ErrRetry) {
			return order, posService.ErrRetry
		}
		msg := err.Error() + " sendOrder"
		if queErr := manager.sendMessageToQue(ctx, telegram.CreateOrder, req, store, msg); queErr != nil {
			log.Err(validator.ErrSendingToQue).Msg("")
		}

		order.Status = "FAILED"
		return manager.failOrder(ctx, order, err)

	}

	if order.Type == "PREORDER" {
		order.Status = models.NEW.String()
	}

	if autoAcceptOn {
		order.Status = models.ACCEPTED.String()
	}

	return manager.successOrder(ctx, order)
}

func (manager OrderManager) getActiveMenu(ctx context.Context, store coreStoreModels.Store, order models.Order) (coreMenuModels.Menu, error) {
	var aggregatorMenuID string

	for _, menu := range store.Menus {
		if menu.Delivery == order.DeliveryService && menu.IsActive {
			aggregatorMenuID = menu.ID
			break
		}
	}

	if aggregatorMenuID != "" {
		aggregatorMenu, err := manager.menuClient.GetMenuByID(ctx, aggregatorMenuID)
		if err != nil {
			return coreMenuModels.Menu{}, err
		}

		return aggregatorMenu, nil
	}

	return coreMenuModels.Menu{}, nil
}

func (manager OrderManager) sendMessageToQue(ctx context.Context, notificationType telegram.NotificationType, order models.Order, store coreStoreModels.Store, err string) error {

	var (
		message   string
		chatIDs   []string
		service   = telegram.Telegram
		queueName = manager.globalConfig.QueConfiguration.Telegram
	)

	if order.DeliveryService == models.QRMENU.String() {
		chatIDs = append(chatIDs, manager.globalConfig.NotificationConfiguration.KwaakaDirectTelegramChatId)
	}

	switch notificationType {
	case telegram.CancelOrder:
		message = telegram.ConstructCancelMessageToNotify(order, store.Name, order.CancelReason.Reason)
		chatIDs = append(chatIDs, store.Telegram.CancelChatID)
	case telegram.CreateOrder, telegram.UpdateOrder:
		message = telegram.ConstructOrderMessageToNotify(service, order, store, err)
		chatIDs = append(chatIDs, manager.globalConfig.NotificationConfiguration.TelegramChatID)
		if manager.notDirect(store.Telegram.CreateOrderChatID) {
			chatIDs = append(chatIDs, store.Telegram.CreateOrderChatID)
		}
		if order.DeliveryService == models.YANDEX.String() && strings.Contains(err, "Creation timeout expired, order automatically transited to error creation status") {
			chatIDs = append(chatIDs, manager.globalConfig.NotificationConfiguration.YandexErrorChatID)
		}
	case telegram.StoreClosed:
		message = telegram.ConstructStoreClosedToNotify(store, "", "")
		chatIDs = append(chatIDs, store.Telegram.StoreStatusChatId)
		if store.Telegram.StoreStatusChatId == "" {
			chatIDs = append(chatIDs, "-1002038506041")
		}
	}

	chatIDs = manager.deleteEmptyChatIDs(chatIDs)

	log.Info().Msgf("telegram creds: queue%v  chat_id %v ", queueName, chatIDs)

	for _, chatID := range chatIDs {
		if err := manager.sqsCli.SendMessage(queueName, message, chatID, store.Telegram.TelegramBotToken); err != nil {
			log.Err(err).Msgf("queue name: %v, order_id %v, notification type: %v, delivery service %v", queueName, order.OrderID, notificationType, order.DeliveryService)
			return err
		}
	}

	return nil
}

func (manager OrderManager) notDirect(chatID string) bool {
	// Aula, Kazbek Saraishyk,Kazbek Bokeikhana
	directChatIDs := []string{"-4265935199", "-4270681463", "-4220435967"}
	for _, id := range directChatIDs {
		if chatID == id {
			return false
		}
	}
	return true
}

func (manager OrderManager) deleteEmptyChatIDs(chatIDs []string) []string {
	var newChatIDs []string
	for _, chatID := range chatIDs {
		if chatID != "" {
			newChatIDs = append(newChatIDs, chatID)
		}
	}
	return newChatIDs
}

func (manager OrderManager) failOrder(ctx context.Context, req models.Order, err error) (models.Order, error) {

	var errs custom.Error

	log.Trace().Err(err).Msgf("%v", validator.ErrFailed)

	req.Status = string(models.STATUS_FAILED)
	req.StatusesHistory = append(req.StatusesHistory, models.OrderStatusUpdate{
		Name: string(models.STATUS_FAILED),
		Time: models.TimeNow().Time,
	})
	if updateErr := manager.orderRepo.UpdateOrder(ctx, req); updateErr != nil {
		return req, validator.ErrFailed
	}

	if err != nil {
		log.Trace().Err(err).Msg("Error while saving order")
		return req, err
	}

	errs.Append(err, validator.ErrFailed)

	return req, errs
}

func (manager OrderManager) skipOrder(ctx context.Context, req models.Order, err error) (models.Order, error) {
	switch err {
	case nil:
		log.Info().Msg("Integration is off or order has discount, skipping...")
	default:
		log.Trace().Err(err).Msg("Skipping order...")
	}

	req.Status = string(models.STATUS_SKIPPED)
	req.StatusesHistory = append(req.StatusesHistory, models.OrderStatusUpdate{
		Name: string(models.STATUS_SKIPPED),
		Time: models.TimeNow().Time,
	})

	err = manager.orderRepo.UpdateOrder(ctx, req)
	if err != nil {
		log.Trace().Err(err).Msg("Error while saving order")
		return req, err
	}

	return req, nil
}

func (manager OrderManager) passOrder(req models.Order) (models.Order, error) {
	log.Info().Msg("Order already exist, skipping...")
	return req, errors.Wrap(validator.ErrPassed, fmt.Sprintf("order %s passed", req.OrderID))
}

func (manager OrderManager) successOrder(ctx context.Context, req models.Order) (models.Order, error) {

	log.Info().Msgf("Order created successfully, status: %v", req.Status)

	req.StatusesHistory = append(req.StatusesHistory, models.OrderStatusUpdate{
		Name: req.Status,
		Time: models.TimeNow().Time,
	})

	if err := manager.orderRepo.UpdateOrder(ctx, req); err != nil {
		return models.Order{}, err
	}

	return req, nil
}

func (manager OrderManager) GetOrder(ctx context.Context, query selector.Order) (models.Order, error) {
	return manager.ds.OrderRepository().GetOrder(ctx, query)
}

func (manager OrderManager) GetActiveOrders(ctx context.Context, query selector.Order) ([]models.Order, error) {

	query = query.SetOrderTimeFrom(time.Now().UTC().Add(-2 * time.Hour)).SetIgnoreStatus("CLOSED")

	orders, _, err := manager.ds.OrderRepository().GetOrders(ctx, query)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (manager OrderManager) GetActivePreorders(ctx context.Context, query selector.Order) ([]models.Order, error) {
	query = query.SetPickupTimeTo(time.Now().UTC().Add(time.Minute * 25)).
		SetType("PREORDER").
		SetStatus(models.WAIT_SENDING.String()).
		SetPickupTimeFrom(time.Now().UTC())

	orders, _, err := manager.ds.OrderRepository().GetOrders(ctx, query)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (manager OrderManager) GetDirectCallCenterActivePreorders(ctx context.Context, query selector.Order) ([]models.Order, error) {

	deliveryServices := []string{coreStoreModels.QRMENU.String(), models.KWAAKA_ADMIN.String()}

	query = query.SetPickupTimeTo(time.Now().UTC().Add(time.Hour * 2)).
		SetType("PREORDER").
		SetStatus(models.WAIT_SENDING.String()).
		SetPickupTimeFrom(time.Now().UTC()).SetDeliveryServices(deliveryServices)

	orders, _, err := manager.ds.OrderRepository().GetOrders(ctx, query)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (manager OrderManager) GetHaniKarimaActivePreOrders(ctx context.Context, query selector.Order) ([]models.Order, error) {
	//для хани и каримы отправлять предзаказы сразу на кассу
	haniRestaurantGroup, err := manager.storeClient.FindStoreGroup(ctx, storeSelector.StoreGroup{ID: "6683fd660077b538b9497c26"})
	if err != nil {
		log.Err(err).Msgf("get hani restaurant group error")
		return nil, err
	}

	karimaRestaurantGroup, err := manager.storeClient.FindStoreGroup(ctx, storeSelector.StoreGroup{ID: "664c37dd3458dbb02d09d191"})
	if err != nil {
		log.Err(err).Msgf("get karima restaurant group error")
		return nil, err
	}

	ekiRestaurantGroup, err := manager.storeClient.FindStoreGroup(ctx, storeSelector.StoreGroup{ID: "664c58922f60d21e6940f94a"})
	if err != nil {
		log.Err(err).Msgf("get eki restaurant group error")
		return nil, err
	}

	storeIDs := append(haniRestaurantGroup.StoreIds, karimaRestaurantGroup.StoreIds...)
	storeIDs = append(storeIDs, ekiRestaurantGroup.StoreIds...)

	deliveryServices := []string{coreStoreModels.QRMENU.String(), models.KWAAKA_ADMIN.String()}

	query = query.SetType("PREORDER").
		SetStatus(models.WAIT_SENDING.String()).
		SetPickupTimeFrom(time.Now().UTC()).SetRestaurants(storeIDs).SetDeliveryServices(deliveryServices)

	orders, _, err := manager.ds.OrderRepository().GetOrders(ctx, query)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (manager OrderManager) UpdateOrder(ctx context.Context, order models.Order) error {
	return manager.orderRepo.UpdateOrder(ctx, order)
}

func (manager OrderManager) CancelOrder(ctx context.Context, req models.CancelOrder) error {
	order, err := manager.orderRepo.CancelOrder(ctx, req)
	if err != nil {
		return err
	}

	if err := manager.CancelOrderInPos(ctx, models.CancelOrderInPos{
		OrderID:         order.OrderID,
		DeliveryService: order.DeliveryService,
		CancelReason:    models.CancelReason{Reason: "CANCEL FROM DELIVERY SERVICE"},
	}); err != nil {
		log.Err(err).Msgf("error while canceling order in pos %s, order id %s", err.Error(), order.OrderID)
		return err
	}

	return nil
}

func (manager OrderManager) CancelOrderInPos(ctx context.Context, req models.CancelOrderInPos) error {
	order, err := manager.orderRepo.GetOrder(ctx, selector.Order{
		OrderID:         req.OrderID,
		DeliveryService: req.DeliveryService,
	})
	if err != nil {
		log.Err(err).Msgf("OrderCore -> OrderManager -> Cancel Order In POS -> GetOrder, order_id %v", req.OrderID)
		return err
	}

	order.Status = string(models.STATUS_CANCELLED_BY_DELIVERY_SERVICE)
	order.StatusesHistory = append(order.StatusesHistory, models.OrderStatusUpdate{
		Name: string(models.STATUS_CANCELLED_BY_DELIVERY_SERVICE),
		Time: time.Now().UTC(),
	})
	order.CancelReason = req.CancelReason
	order.PaymentStrategy = req.PaymentStrategy

	err = manager.orderRepo.UpdateOrder(ctx, order)
	if err != nil {
		log.Err(err).Msgf("OrderCore -> OrderManager -> Cancel Order -> UpdateOrder, order_id %s, id %s", req.OrderID, order.OrderID)
		return err
	}

	store, err := manager.storeClient.FindStore(ctx, storeDto.StoreSelector{
		ID: order.RestaurantID,
	})
	if err != nil {
		log.Err(err).Msgf("OrderCore -> OrderManager -> Cancel Order -> FindStore for %s, order_id %s", order.RestaurantID, req.OrderID)
		return err
	}

	if err := manager.sendMessageToQue(ctx, telegram.CancelOrder, order, store, ""); err != nil {
		log.Err(validator.ErrSendingToQue).Msg("")
	}

	if store.PosType == "" {
		return nil
	}

	posManager, err := pos.NewPosManager(store.PosType, manager.globalConfig, manager.ds, store, coreMenuModels.Menu{}, coreMenuModels.Promo{}, nil, coreMenuModels.Menu{}, manager.storeClient)
	if err != nil {
		log.Trace().Err(validator.ErrPosInitialize).Msgf("posName: %s, restaurantID: %s", store.PosType, store.ID)
		return validator.ErrPosInitialize
	}

	if order.DeliveryService == models.WOLT.String() && order.PosType == iikoModels.IIKO.String() {
		agg, err := manager.aggregatorFactory.GetAggregator(req.DeliveryService, store)
		if err != nil {
			log.Err(err).Msgf("OrderCore -> OrderManager -> Cancel Order -> GetAggregator for %s, order_id %s", order.RestaurantID, req.OrderID)
			return err
		}
		woltOrder, err := agg.GetAggregatorOrder(ctx, req.OrderID)
		if err != nil {
			log.Err(err).Msgf("OrderCore -> OrderManager -> Cancel Order -> GetAggregatorOrder for %s, order_id %s", order.RestaurantID, req.OrderID)
			return err
		}
		if woltOrder.Delivery.Status == http.AcceptedWoltOrderStatus ||
			woltOrder.Delivery.Status == http.ReadyWoltOrderStatus ||
			woltOrder.Delivery.Status == http.DeliveredWoltOrderStatus {
			return errors.New("order is already in PRODUCTION state, can not cancel order in POS")
		}
	}

	storeGroup, err := manager.storeClient.FindStoreGroup(ctx, storeSelector.NewEmptyStoreGroupSearch().SetID(store.RestaurantGroupID))
	if err != nil {
		log.Err(err).Msgf("get store group for store: %s", store.RestaurantGroupID)
		return err
	}
	if storeGroup.CancelOrderAllowed {
		log.Info().Msgf("cancel order allowed for store group: %s", storeGroup.ID)
		err = posManager.CancelOrder(ctx, order, req.CancelReason.Reason, req.PaymentStrategy, store)
		if err != nil {
			log.Err(err).Msgf("OrderCore -> OrderManager -> Cancel Order -> CancelOrder, pos_order_id %s, order_id %s", order.PosOrderID, req.OrderID)
			if errors.Is(err, errs.ErrUnsupportedMethod) {
				return nil
			}
			return err
		}
		return nil
	}

	switch order.DeliveryService {
	case coreStoreModels.QRMENU.String(), models.KWAAKA_ADMIN.String():
		err = posManager.CancelOrder(ctx, order, req.CancelReason.Reason, req.PaymentStrategy, store)
		if err != nil {
			log.Err(err).Msgf("OrderCore -> OrderManager -> Cancel Order -> CancelOrder, pos_order_id %s, order_id %s", order.PosOrderID, req.OrderID)
			if errors.Is(err, errs.ErrUnsupportedMethod) {
				return nil
			}
			return err
		}
	default:
		err = posManager.UpdateOrderProblem(ctx, store.IikoCloud.OrganizationID, order.PosOrderID)
		if err != nil {
			log.Err(err).Msgf("OrderCore -> OrderManager -> Cancel Order -> CancelOrder, pos_order_id %s, order_id %s", order.PosOrderID, req.OrderID)
			if errors.Is(err, errs.ErrUnsupportedMethod) {
				return nil
			}
			return err
		}
	}

	return nil
}

func (manager OrderManager) UpdateOrderStatus(ctx context.Context, posOrderID, pos, externalStatus, errorDescription string) error {
	query := selector.EmptyOrderSearch()
	switch pos {
	case models.RKeeper.String():
		query = query.SetOrderID(posOrderID)
	case models.WaitSending.String():
		query = query.SetID(posOrderID)
	case models.BurgerKing.String():
		query = query.SetOrderID(posOrderID)
	case models.Poster.String():
		query = query.SetOrderID(posOrderID)
	case models.Kwaaka.String():
		query = query.SetID(posOrderID)
	default:
		query = query.SetPosOrderID(posOrderID)
	}

	order, err := manager.orderRepo.GetOrder(ctx, query)
	if err != nil {
		return fmt.Errorf("order %w", err)
	}

	store, err := manager.storeClient.FindStore(ctx, storeDto.StoreSelector{
		ID:      order.RestaurantID,
		PosType: order.PosType,
	})
	if err != nil {
		return err
	}

	posStatus, err := manager.compareStatuses(order.PosType, externalStatus, order.Status)
	if err != nil {
		return err
	}

	log.Info().Msgf("Pos status = %s", posStatus)

	if !manager.validateOrderStatusQueue(order.Status, store.PosType, posStatus) {
		log.Err(errs.ErrValidateOrderStatusQueue)
		return errs.ErrValidateOrderStatusQueue
	}

	if err = manager.ds.OrderRepository().UpdateOrderStatus(ctx, query, posStatus.String(), errorDescription); err != nil {
		return err
	}

	if posStatus == models.FAILED {
		log.Info().Msg("order was failed, skipped aggregator update")
		if queErr := manager.sendMessageToQue(ctx, telegram.UpdateOrder, order, store, errorDescription); queErr != nil {
			log.Err(validator.ErrSendingToQue).Msg("")
		}
		return nil
	}

	if store.Kwaaka3PL.Is3pl {
		err := manager.sendKwaaka3plOrder(ctx, order, store, posStatus)
		if err != nil {
			log.Err(err).Msg("send 3 pl order error")
			return err
		}
	}

	aggregatorStatus, newPosStatus := manager.compareAggregatorStatuses(order, posStatus, store)

	shouldSubmit := manager.shouldSubmitStatus(order, store, aggregatorStatus)

	if !order.IsChildOrder || !shouldSubmit {
		if err := manager.updateOrderInAggregator(ctx, order, store, aggregatorStatus); err != nil {
			return err
		}
	}

	if order.DeliveryService == models.QRMENU.String() || order.DeliveryService == models.KWAAKA_ADMIN.String() {
		var (
			notifier notifierpkg.Notifier
			err      error
			chatID   *int64
		)

		storeGroup, err := manager.storeClient.FindStoreGroup(ctx, storeSelector.StoreGroup{
			StoreIDs: []string{order.RestaurantID},
		})
		if err != nil {
			return err
		}

		// TODO: applying strategy based on user's chat id is not correct, implement logic of strategy choice based different parameter
		// (probably by checking type of order in 'order' struct)
		chatID, err = manager.telegramService.GetUserChatID(ctx)
		var numErr *strconv.NumError
		if err != nil && !errors.As(err, &numErr) {
			return err
		}
		if chatID != nil {
			notifier, err = telegram2.NewTelegramNotifier(manager.globalConfig.OrderBotToken, manager.telegramService)
			if err != nil {
				log.Err(err).Msg("(UpdateOrderStatus): error initializing telegram bot api")
				return err
			}
		} else {
			notifier = whatsapp.NewWhatsappNotifier(manager.sqsCli, manager.globalConfig.QueueUrls.WhatsappMessagesQueueUrl)
		}

		if err := notifier.Notify(ctx, newPosStatus, order, storeGroup, store); err != nil {
			return err
		}
	}

	return nil
}

func (manager OrderManager) shouldSubmitStatus(order models.Order, store coreStoreModels.Store, aggStatus string) bool {
	if !store.DeferSubmission.IsDeferSumbission || !order.IsDeferSubmission {
		return false
	}

	switch aggStatus {
	case models.READY_FOR_PICKUP.String(), models.Ready.String():
		return true
	default:
		return false
	}
}

func (manager OrderManager) sendKwaaka3plOrder(ctx context.Context, order models.Order, store coreStoreModels.Store, posStatus models.PosStatus) error {

	if posStatus == models.CANCELLED_BY_POS_SYSTEM && order.DeliveryOrderID != "" {

		if err := manager.kwaaka3plService.Save3plHistory(ctx, order.DeliveryOrderID, models.DeliveryAddress{}, models.Customer{}); err != nil {
			log.Err(err).Msgf("error: save 3pl history for delivery order id: %s", order.DeliveryOrderID)
			return err
		}

		err := manager.kwaaka3plService.CancelCourierSearch(ctx, order.DeliveryOrderID)
		if err != nil {
			log.Error().Msgf("couldn't cancel courier search in 3pl order %s", err.Error())
			return fmt.Errorf("core/managers/order - fn sendKwaaka3plOrder - fn CancelCourierSearch: %w", err)
		}
		return nil
	}

	if (posStatus == models.OUT_FOR_DELIVERY || posStatus == models.PICKED_UP_BY_CUSTOMER) && order.DeliveryDispatcher == coreStoreModels.SELFDELIVERY.String() {
		var posStatusStr string
		switch posStatus {
		case models.OUT_FOR_DELIVERY:
			posStatusStr = models.OUT_FOR_DELIVERY_str
		case models.PICKED_UP_BY_CUSTOMER:
			posStatusStr = models.PICKED_UP_BY_CUSTOMER_str
		default:
			return nil
		}

		err := manager.kwaaka3plService.MapIikoStatusTo3plStatus(ctx, posStatusStr, order.Customer.PhoneNumber, order.RestaurantID)
		if err != nil {
			log.Error().Msgf("couldn't map iiko status to 3pl")
			return fmt.Errorf("core/managers/order - fn sendKwaaka3plOrder - fn MapIikoStatuseTo3plStatus: %w", err)
		}
		return nil
	}

	if order.IsMarketplace || !order.SendCourier {
		return nil
	}
	if order.DeliveryOrderID != "" || order.DeliveryDispatcher == "" || order.DeliveryDispatcher == coreStoreModels.SELFDELIVERY.String() {
		return nil
	}

	switch posStatus {
	case models.NEW, models.ACCEPTED, models.COOKING_STARTED, models.ON_WAY, models.WAIT_SENDING, models.CANCELLED_BY_POS_SYSTEM, models.FAILED:
		return nil
	default:
		log.Info().Msgf("going to create 3pl order, posStatus: %s, orderID: %s", posStatus.String(), order.ID)
	}

	items := make([]models3.Item, 0, len(order.Products))
	for _, product := range order.Products {
		items = append(items, models3.Item{
			Name:     product.Name,
			ID:       product.ID,
			Quantity: product.Quantity,
			Price:    product.Price.Value,
		})
	}

	if err := manager.kwaaka3plService.Create3plOrder(ctx, models3.CreateDeliveryRequest{
		ID:                order.ID,
		FullDeliveryPrice: order.FullDeliveryPrice,
		Provider:          order.DeliveryDispatcher,
		PickUpTime:        time.Now().Add(time.Duration(order.DispatcherDeliveryTime) * time.Minute),
		DeliveryAddress: models3.Address{
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
		StoreAddress: models3.Address{
			Label:   store.Address.City + ", " + store.Address.Street,
			Lon:     store.Address.Coordinates.Longitude,
			Lat:     store.Address.Coordinates.Latitude,
			Comment: store.Address.Entrance,
		},
		CustomerInfo: models3.CustomerInfo{
			Name:  manager.setCustomerName(order.Customer.Name),
			Phone: manager.setCustomerPhoneNumber(order.Customer.PhoneNumber),
			Email: manager.setCustomerEmail(order.Customer.Email),
		},
		StoreInfo: models3.StoreInfo{
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

	proposals, err := manager.kwaaka3plService.ListPotentialProviders(ctx, models3.ListProvidersRequest{
		Address: models3.OrderAddress{
			City:   order.DeliveryAddress.City,
			Street: order.DeliveryAddress.Street,
			Coordinates: models3.Coordinates{
				Lat: order.DeliveryAddress.Latitude,
				Lon: order.DeliveryAddress.Longitude,
			},
			Language: "ru",
		},
		RestaurantCoordinates: models3.Coordinates{
			Lat: store.Address.Coordinates.Latitude,
			Lon: store.Address.Coordinates.Longitude,
		},
		MinPreparationTimeMinutes: int(order.EstimatedPickupTime.Value.Sub(time.Now().UTC()).Minutes()),
		ItemsSettings: models3.ItemsSettings{
			Quantity: len(items),
			Size: models3.Size{
				Height: 0.3,
				Width:  0.3,
				Length: 0.3,
			},
			Weight: 1,
		},
		IndriveAvailable: store.Kwaaka3PL.IndriveAvailable,
		WoltAvailable:    store.Kwaaka3PL.WoltDriveAvailable,
		YandexAvailable:  store.Kwaaka3PL.YandexAvailable,
	})
	if err != nil {
		log.Error().Msgf("couldn't get proposals to save in order %s", err.Error())
	}

	var proposalsToSave []models.Proposal

	for _, proposal := range proposals {
		if proposal.Provider == nil {
			continue
		}
		proposalsToSave = append(proposalsToSave, models.Proposal{
			Price:               proposal.Provider.Price,
			TimeEstimateMinutes: proposal.Provider.TimeEstimateMinutes,
			ProviderService:     proposal.Provider.ProviderService,
			Priority:            proposal.Provider.Priority,
		})
	}

	if err = manager.repository.SetProposals(ctx, order.OrderID, proposalsToSave); err != nil {
		log.Error().Msgf("couldn't save proposals in order %s", err.Error())
	}

	return nil
}

func (manager OrderManager) cancelAndFindAnotherProvider(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	deliveryInfo, err := manager.kwaaka3plService.GetDeliveryInfoByOrderId(ctx, order.OrderID)
	if err != nil {
		return err
	}

	for _, status := range deliveryInfo.Statuses {
		if _, exist := nonCancelableStatuses[status.Status]; exist {
			return nil
		}
	}

	items := make([]models3.Item, 0, len(order.Products))
	for _, product := range order.Products {
		items = append(items, models3.Item{
			Name:     product.Name,
			ID:       product.ID,
			Quantity: product.Quantity,
			Price:    product.Price.Value,
		})
	}

	providerProposals, err := manager.kwaaka3plService.ListPotentialProviders(ctx, models3.ListProvidersRequest{
		Address: models3.OrderAddress{
			City:   order.DeliveryAddress.City,
			Street: order.DeliveryAddress.Street,
			Coordinates: models3.Coordinates{
				Lat: order.DeliveryAddress.Latitude,
				Lon: order.DeliveryAddress.Longitude,
			},
			Language: "ru",
		},
		RestaurantCoordinates: models3.Coordinates{
			Lat: store.Address.Coordinates.Latitude,
			Lon: store.Address.Coordinates.Longitude,
		},
		MinPreparationTimeMinutes: int(order.EstimatedPickupTime.Value.Sub(time.Now().UTC()).Minutes()),
		ItemsSettings: models3.ItemsSettings{
			Quantity: len(items),
			Size: models3.Size{
				Height: 0.3,
				Width:  0.3,
				Length: 0.3,
			},
			Weight: 1,
		},
	})
	if err != nil {
		return err
	}

	var newDispatcher models3.Delivery
	min := 0
	for _, promise := range providerProposals {
		if promise.Provider == nil {
			continue
		}
		if order.DeliveryDispatcher == promise.Provider.ProviderService {
			continue
		}

		if min == 0 {
			min = promise.Provider.Priority
			newDispatcher.FullDeliveryPrice = float64(promise.Provider.Price.Amount)
			newDispatcher.Dispatcher = promise.Provider.ProviderService
			newDispatcher.DeliveryTime = int32(promise.Provider.TimeEstimateMinutes)
			newDispatcher.KwaakaChargedDeliveryPrice = promise.Provider.Price.KwaakaChargeSum
			newDispatcher.Priority = promise.Provider.Priority
		}
		if promise.Provider.Priority < min {
			newDispatcher.FullDeliveryPrice = float64(promise.Provider.Price.Amount)
			newDispatcher.Dispatcher = promise.Provider.ProviderService
			newDispatcher.DeliveryTime = int32(promise.Provider.TimeEstimateMinutes)
			newDispatcher.KwaakaChargedDeliveryPrice = promise.Provider.Price.KwaakaChargeSum
			newDispatcher.Priority = promise.Provider.Priority
			min = promise.Provider.Priority
		}
	}

	if newDispatcher.Dispatcher == "" {
		return nil
	}

	err = manager.kwaaka3plService.Cancel3plOrder(ctx, order.ID, models.OrderInfoForTelegramMsg{
		RestaurantName:        store.Name,
		RestaurantAddress:     store.Address.City + " " + store.Address.Street,
		RestaurantPhoneNumber: store.StorePhoneNumber,
		OrderId:               order.ID,
		Id3plOrder:            order.DeliveryOrderID,
		CustomerName:          order.Customer.Name,
		CustomerPhoneNumber:   order.Customer.PhoneNumber,
		CustomerAddress:       order.DeliveryAddress.Label,
		DeliveryService:       order.DeliveryDispatcher,
	})
	if err != nil {
		return err
	}

	if err = manager.kwaaka3plService.Create3plOrder(ctx, models3.CreateDeliveryRequest{
		ID:                order.ID,
		FullDeliveryPrice: newDispatcher.FullDeliveryPrice,
		Provider:          newDispatcher.Dispatcher,
		PickUpTime:        time.Now().Add(time.Duration(order.DispatcherDeliveryTime) * time.Minute),
		DeliveryAddress: models3.Address{
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
		StoreAddress: models3.Address{
			Label:   store.Address.City + ", " + store.Address.Street,
			Lon:     store.Address.Coordinates.Longitude,
			Lat:     store.Address.Coordinates.Latitude,
			Comment: store.Address.Entrance,
		},
		CustomerInfo: models3.CustomerInfo{
			Name:  manager.setCustomerName(order.Customer.Name),
			Phone: manager.setCustomerPhoneNumber(order.Customer.PhoneNumber),
			Email: manager.setCustomerEmail(order.Customer.Email),
		},
		StoreInfo: models3.StoreInfo{
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

func (manager OrderManager) updateOrderInAggregator(ctx context.Context, order models.Order, store coreStoreModels.Store, aggregatorStatus string) error {
	if order.DeliveryService == models.YANDEX.String() {
		return nil
	}

	aggregatorManager, err := aggregator.NewAggregatorManager(order.DeliveryService, manager.globalConfig, store)
	if err != nil {
		return err
	}

	switch order.DeliveryService {
	case models.EMENU.String():
		if err := aggregatorManager.UpdateOrderStatus(ctx, selector.EmptyOrderStatusUpdate().SetOrderID(order.OrderID).SetOrderStatus(aggregatorStatus)); err != nil {
			return err
		}

	case models.GLOVO.String(), models.CHOCOFOOD.String():
		if aggregatorStatus == "" {
			log.Info().Msgf("aggregator status is empty, order id %s", order.OrderID)
			return nil
		}

		if err := aggregatorManager.UpdateOrderStatus(ctx, selector.EmptyOrderStatusUpdate().
			SetOrderID(order.OrderID).
			SetOrderStatus(aggregatorStatus).
			SetStoreID(order.StoreID)); err != nil {
			return err
		}

	case models.WOLT.String():
		switch aggregatorStatus {

		case models.Accept.String():
			var adjustedPickUpTime *time.Time = nil

			if !store.Wolt.IgnorePickupTime {
				adjustedPickUpTime = &order.EstimatedPickupTime.Value.Time
			}

			log.Info().Msgf("result estimated pick up time: %v", adjustedPickUpTime)

			switch order.RestaurantSelfDelivery {
			case false:
				if err := aggregatorManager.AcceptOrder(ctx, order.OrderID, adjustedPickUpTime); err != nil {
					log.Trace().Err(err).Msgf("accept order, order_id=%v, pick_up_time=%v", order.OrderID, adjustedPickUpTime)
					return err
				}
				log.Info().Msgf("success accept order, order_id=%v, pick_up_time=%v", order.OrderID, adjustedPickUpTime)

			case true:
				deliveryTime := time.Now().UTC().Add(time.Minute * time.Duration(store.Wolt.CookingTime+store.Wolt.AdjustedPickupMinutes))
				if err := aggregatorManager.AcceptSelfDeliveryOrder(ctx, order.OrderID, &deliveryTime); err != nil {
					log.Trace().Err(err).Msgf("accept selfDelivery order, order_id=%v, deliveryTime=%v", order.OrderID, deliveryTime)
					return err
				}
				log.Info().Msgf("success accept selfDelivery order, order_id=%v, deliveryTime=%v", order.OrderID, deliveryTime)
			}

		case models.Reject.String():
			if store.Wolt.IgnoreStatusUpdate {
				log.Info().Msgf("Ignore status update: reject")
				return nil
			}

			if err := aggregatorManager.RejectOrder(ctx, order.OrderID, ""); err != nil {
				return err
			}

		case models.Ready.String():
			if store.Wolt.IgnoreStatusUpdate {
				log.Info().Msgf("Ignore status update: ready")
				return nil
			}

			if err := aggregatorManager.MarkOrder(ctx, order.OrderID); err != nil {
				return err
			}

		case models.Confirm.String():
			if store.Wolt.IgnoreStatusUpdate {
				log.Info().Msgf("Ignore status update: confirm")
				return nil
			}

			if err := aggregatorManager.ConfirmPreOrder(ctx, order.OrderID); err != nil {
				return err
			}

		case models.Delivered.String():
			if store.Wolt.IgnoreStatusUpdate {
				log.Info().Msgf("Ignore status update: delivered")
				return nil
			}

			if err := aggregatorManager.DeliveredOrder(ctx, order.OrderID); err != nil {
				return err
			}

		default:
			log.Info().Msgf("aggregator status is empty: %v", aggregatorStatus)
			return nil
		}
	case models.TALABAT.String():
		switch aggregatorStatus {
		case models.OrderAccepted.String():
			pickupTime := time.Now().UTC()
			if err := aggregatorManager.AcceptOrder(ctx, order.OrderID, &pickupTime); err != nil {
				return err
			}
		case models.OrderRejected.String():
			if err := aggregatorManager.RejectOrder(ctx, order.OrderID, "TECHNICAL_PROBLEM"); err != nil {
				return err
			}
		case models.OrderPickedUp.String():
			if err := aggregatorManager.DeliveredOrder(ctx, order.OrderID); err != nil {
				return err
			}
		case models.OrderPrepared.String():
			if err := aggregatorManager.MarkOrder(ctx, order.OrderID); err != nil {
				return err
			}
		default:
			log.Info().Msgf("aggregator status is invalid: %v", aggregatorStatus)
			return nil
		}

	case models.DELIVEROO.String():
		switch aggregatorStatus {
		case models.Rejected.String():
			if err := aggregatorManager.RejectOrder(ctx, order.OrderID, ""); err != nil {
				return err
			}
		case "":
			log.Info().Msgf("need to define what to do with other statuses")
			return nil
		default:
			if err := aggregatorManager.UpdateOrderStatus(ctx, selector.EmptyOrderStatusUpdate().SetOrderID(order.OrderID).SetOrderStatus(aggregatorStatus)); err != nil {
				return err
			}
		}

	case models.KWAAKA_ADMIN.String(), models.QRMENU.String():
		return nil

	case models.STARTERAPP.String():
		if err := aggregatorManager.UpdateOrderStatus(ctx, selector.EmptyOrderStatusUpdate().SetOrderID(order.OrderID).SetOrderStatus(aggregatorStatus)); err != nil {
			return err
		}

	default:
		log.Info().Msgf("unsupported delivery service: %v", order.DeliveryService)
		return nil
	}

	return nil
}

func (manager OrderManager) GetOrdersWithFilters(ctx context.Context, query selector.Order) ([]models.Order, int, error) {
	return manager.ds.OrderRepository().GetOrders(ctx, query)
}

func (manager OrderManager) UpdateOrderStatusByID(ctx context.Context, orderID, pos, externalStatus string) error {
	order, err := manager.orderRepo.GetOrder(ctx,
		selector.EmptyOrderSearch().
			SetID(orderID))
	if err != nil {
		return err
	}

	posStatus, err := manager.compareStatuses(pos, externalStatus, order.Status)
	if err != nil {
		return err
	}

	return manager.ds.OrderRepository().UpdateOrderStatusByID(ctx, orderID, pos, posStatus.String())
}

func (manager OrderManager) UpdateOrderStatusInDS(ctx context.Context, orderID string, posStatus models.PosStatus) error {
	order, err := manager.orderRepo.GetOrder(ctx,
		selector.EmptyOrderSearch().
			SetID(orderID))
	if err != nil {
		return err
	}

	store, err := manager.storeClient.FindStore(ctx, storeDto.StoreSelector{
		ID: order.RestaurantID,
	})
	if err != nil {
		return err
	}

	aggregatorStatus, newPosStatus := manager.compareAggregatorStatuses(order, posStatus, store)

	if err := manager.updateOrderInAggregator(ctx, order, store, aggregatorStatus); err != nil {
		return err
	}

	if err := manager.ds.OrderRepository().UpdateOrderStatus(ctx, selector.OrderSearch().SetID(order.ID), newPosStatus, ""); err != nil {
		return err
	}

	return nil
}

func (manager OrderManager) compareAggregatorStatuses(order models.Order, posStatus models.PosStatus, store coreStoreModels.Store) (string, string) {
	switch order.DeliveryService {
	case models.WOLT.String():
		log.Info().Msg("delivery service: wolt")

		var statusConfig []coreStoreModels.Status

		switch order.Type {
		case models.Instant:
			log.Info().Msgf("purchase type: INSTANT")
			statusConfig = store.Wolt.PurchaseTypes.Instant
		case models.Preorder:
			log.Info().Msgf("purchase type: PREORDER")
			statusConfig = store.Wolt.PurchaseTypes.Preorder
		}

		if order.IsPickedUpByCustomer {
			log.Info().Msgf("purchase type: TAKEAWAY")
			statusConfig = store.Wolt.PurchaseTypes.TakeAway
		}

		for _, status := range statusConfig {
			if status.PosStatus == posStatus.String() {
				log.Info().Msgf("[SPECIAL MATCHING] pos status: %v, wolt status: %v", status.PosStatus, status.Status)
				return status.Status, posStatus.String()
			}
		}

		log.Info().Msgf("[DEFAULT MATCHING], pos status: %v", posStatus)

		switch posStatus {
		case models.ACCEPTED, models.WAIT_COOKING, models.READY_FOR_COOKING, models.COOKING_STARTED, models.WAIT_SENDING:

			if order.Type == models.Preorder {
				return models.Confirm.String(), posStatus.String()
			}

			return models.Accept.String(), posStatus.String()
		case models.COOKING_COMPLETE, models.CLOSED, models.READY_FOR_PICKUP, models.ON_WAY, models.DELIVERED, models.OUT_FOR_DELIVERY:
			return models.Ready.String(), posStatus.String()
		case models.CANCELLED_BY_POS_SYSTEM:
			return models.Reject.String(), posStatus.String()
		case models.PICKED_UP_BY_CUSTOMER:
			return models.Delivered.String(), posStatus.String()
		default:
			return "", posStatus.String()
		}
	case models.GLOVO.String():
		log.Info().Msg("delivery service: glovo")

		for _, status := range store.Glovo.PurchaseTypes.Instant {
			if status.PosStatus == posStatus.String() {
				log.Info().Msgf("[SPECIAL MATCHING] pos status: %v, glovo status: %v", status.PosStatus, status.Status)
				return status.Status, posStatus.String()
			}
		}

		log.Info().Msgf("[DEFAULT MATCHING], pos status: %v", posStatus)

		switch posStatus {
		case models.ACCEPTED, models.COOKING_STARTED, models.WAIT_SENDING:
			return models.ACCEPTED.String(), posStatus.String()
		case models.READY_FOR_PICKUP, models.COOKING_COMPLETE, models.CLOSED:
			return models.READY_FOR_PICKUP.String(), posStatus.String()
		case models.OUT_FOR_DELIVERY:
			return models.OUT_FOR_DELIVERY.String(), posStatus.String()
		case models.PICKED_UP_BY_CUSTOMER:
			return models.PICKED_UP_BY_CUSTOMER.String(), posStatus.String()
		default:
			return "", posStatus.String()
		}

	case models.CHOCOFOOD.String():
		log.Info().Msg("delivery service: chocofood")

		switch posStatus {
		case models.ACCEPTED, models.COOKING_STARTED:
			return models.Accepted.String(), posStatus.String()
		case models.READY_FOR_PICKUP, models.COOKING_COMPLETE:
			return models.ReadyForPickup.String(), posStatus.String()
		case models.OUT_FOR_DELIVERY:
			return models.OutForDelivery.String(), posStatus.String()
		case models.PICKED_UP_BY_CUSTOMER:
			return models.PickedUpByCustomer.String(), posStatus.String()
		case models.CANCELLED_BY_POS_SYSTEM:
			return models.Rejected.String(), posStatus.String()

		default:
			return "", posStatus.String()
		}
	case models.EMENU.String():
		for _, external := range store.ExternalConfig {
			if external.Type == models.EMENU.String() {
				for _, status := range external.PurchaseTypes.Instant {
					if status.PosStatus == posStatus.String() {
						log.Info().Msgf("[SPECIAL MATCHING] pos status: %v, emenu status: %v", status.PosStatus, status.Status)
						return status.Status, posStatus.String()
					}
				}
			}
		}

		log.Info().Msgf("[DEFAULT MATCHING], pos status: %v", posStatus)

		switch posStatus {
		case models.ACCEPTED:
			return models.ACCEPTED.String(), posStatus.String()
		case models.COOKING_STARTED:
			return models.COOKING_STARTED.String(), posStatus.String()
		case models.READY_FOR_PICKUP, models.COOKING_COMPLETE, models.PICKED_UP_BY_CUSTOMER, models.OUT_FOR_DELIVERY:
			return models.COOKING_COMPLETE.String(), posStatus.String()
		case models.CLOSED:
			return models.CLOSED.String(), posStatus.String()
		case models.CANCELLED_BY_POS_SYSTEM:
			return "", posStatus.String()
		default:
			return "", posStatus.String()
		}

	case models.YANDEX.String():
		for _, external := range store.ExternalConfig {
			if external.Type == models.YANDEX.String() {
				for _, status := range external.PurchaseTypes.Instant {
					if status.PosStatus == posStatus.String() {
						log.Info().Msgf("[SPECIAL MATCHING] pos status: %v, yandex status: %v", status.PosStatus, status.Status)
						return status.Status, posStatus.String()
					}
				}
			}
		}
		log.Info().Msgf("[DEFAULT MATCHING], pos status: %v", posStatus)
		return "", posStatus.String()
	case models.TALABAT.String():
		log.Info().Msg("delivery service: chocofood")

		switch posStatus {
		case models.ACCEPTED, models.WAIT_COOKING, models.READY_FOR_COOKING, models.COOKING_STARTED, models.WAIT_SENDING:
			return models.OrderAccepted.String(), posStatus.String()
		case models.COOKING_COMPLETE, models.CLOSED, models.READY_FOR_PICKUP, models.ON_WAY, models.DELIVERED, models.OUT_FOR_DELIVERY:
			return models.OrderPrepared.String(), posStatus.String()
		case models.CANCELLED_BY_POS_SYSTEM:
			return models.OrderRejected.String(), posStatus.String()
		case models.PICKED_UP_BY_CUSTOMER:
			return models.OrderPickedUp.String(), posStatus.String()
		default:
			return "", posStatus.String()
		}
	case models.DELIVEROO.String():
		switch posStatus {
		case models.ACCEPTED, models.COOKING_STARTED:
			return models.AcceptedDeliveroo.String(), posStatus.String()
		case models.CANCELLED_BY_POS_SYSTEM:
			return models.RejectedDeliveroo.String(), posStatus.String()
		default:
			return "", posStatus.String()
		}
	case models.KWAAKA_ADMIN.String():
		return posStatus.String(), posStatus.String()
	case models.QRMENU.String():
		return posStatus.String(), posStatus.String()

	case models.STARTERAPP.String():
		switch posStatus {
		case models.ACCEPTED, models.READY_FOR_COOKING, models.WAIT_SENDING:
			return models.Created.String(), posStatus.String()
		case models.COOKING_STARTED:
			return models.InProgress.String(), posStatus.String()
		case models.COOKING_COMPLETE, models.CLOSED, models.READY_FOR_PICKUP, models.ON_WAY, models.OUT_FOR_DELIVERY:
			return models.Cooked.String(), posStatus.String()
		case models.CANCELLED_BY_POS_SYSTEM:
			return models.Canceled.String(), posStatus.String()
		default:
			return "", posStatus.String()
		}
	default:
		return "", posStatus.String()
	}
}

func compareFoodBandStatuses(externalStatus, previousStatus string) (models.PosStatus, error) {
	previousStatusPriority, _ := getFoodBandStatusPriority(previousStatus)
	externalStatusPriority, status := getFoodBandStatusPriority(externalStatus)

	if externalStatusPriority == 0 {
		return 0, models.StatusIsNotExist
	}

	if previousStatusPriority > externalStatusPriority {
		return status, models.InvalidStatusPriority
	}
	return status, nil
}

func getFoodBandStatusPriority(status string) (int, models.PosStatus) {
	switch status {
	case "ACCEPTED":
		return 1, models.ACCEPTED
	case "COOKING_STARTED":
		return 2, models.COOKING_STARTED
	case "COOKING_COMPLETE":
		return 3, models.COOKING_COMPLETE
	case "READY_FOR_PICKUP":
		return 4, models.READY_FOR_PICKUP
	case "OUT_FOR_DELIVERY":
		return 5, models.OUT_FOR_DELIVERY
	case "PICKED_UP_BY_CUSTOMER":
		return 6, models.PICKED_UP_BY_CUSTOMER
	case "DELIVERED":
		return 6, models.DELIVERED
	case "CLOSED":
		return 7, models.CLOSED
	default:
		return 0, 0
	}
}

func (manager OrderManager) compareStatuses(pos, externalStatus, previousStatus string) (models.PosStatus, error) {
	switch pos {
	case models.RKeeper7XML.String():
		switch externalStatus {
		case "0":
			return models.ACCEPTED, nil
		case "1":
			return models.CLOSED, nil
		default:
			return 0, models.StatusIsNotExist
		}
	case models.IIKO.String(), models.Syrve.String():
		switch externalStatus {
		case "WAIT_SENDING":
			return models.WAIT_SENDING, nil
		case "PAYMENT_NEW":
			return models.PAYMENT_NEW, nil
		case "PAYMENT_IN_PROGRESS":
			return models.PAYMENT_IN_PROGRESS, nil
		case "PAYMENT_SUCCESS":
			return models.PAYMENT_SUCCESS, nil
		case "PAYMENT_CANCELLED":
			return models.PAYMENT_CANCELLED, nil
		case "PAYMENT_WAITING":
			return models.PAYMENT_DELETED, nil
		case "PAYMENT_DELETED":
			return models.PAYMENT_WAITING, nil
		case "CookingStarted":
			return models.COOKING_STARTED, nil
		case "WaitCooking", "ReadyForCooking", "Unconfirmed":
			return models.ACCEPTED, nil
		case "Waiting":
			return models.READY_FOR_PICKUP, nil
		case "CookingCompleted":
			return models.COOKING_COMPLETE, nil
		case "Delivered":
			return models.PICKED_UP_BY_CUSTOMER, nil
		case "OnWay":
			return models.OUT_FOR_DELIVERY, nil
		case "Closed":
			return models.CLOSED, nil
		case "NEW":
			return models.NEW, nil
		case "Cancelled":
			return models.CANCELLED_BY_POS_SYSTEM, nil
		case "Error":
			return models.FAILED, nil
		default:
			return 0, models.StatusIsNotExist
		}
	case models.JOWI.String():
		switch externalStatus {
		case "WAIT_SENDING":
			return models.WAIT_SENDING, nil
		case "0":
			return models.NEW, nil
		case "1":
			return models.ACCEPTED, nil
		case "2":
			return models.CANCELLED_BY_POS_SYSTEM, nil
		case "3":
			return models.OUT_FOR_DELIVERY, nil
		case "4":
			return models.CLOSED, nil
		default:
			return 0, models.StatusIsNotExist
		}
	case models.RKeeper.String():
		switch externalStatus {
		case "WAIT_SENDING":
			return models.WAIT_SENDING, nil
		case "Canceled":
			return models.CANCELLED_BY_POS_SYSTEM, nil
		case "Created":
			return models.ACCEPTED, nil
		case "Cooking":
			return models.COOKING_STARTED, nil
		case "Ready":
			return models.COOKING_COMPLETE, nil
		case "Complited", "IssuedOut":
			return models.CLOSED, nil
		default:
			return 0, models.StatusIsNotExist
		}
	case models.Paloma.String():
		switch externalStatus {
		case "new":
			return models.ACCEPTED, nil
		case "cooking":
			return models.COOKING_STARTED, nil
		case "on_way":
			return models.ON_WAY, nil
		case "completed":
			return models.CLOSED, nil
		case "canceled":
			return models.CANCELLED_BY_POS_SYSTEM, nil
		default:
			return 0, models.StatusIsNotExist
		}
	case models.BurgerKing.String():
		switch externalStatus {
		case models.ACCEPTED.String(), models.COOKING_STARTED.String(), models.WAIT_SENDING.String():
			return models.ACCEPTED, nil
		case models.READY_FOR_PICKUP.String(), models.COOKING_COMPLETE.String(), models.CLOSED.String():
			return models.READY_FOR_PICKUP, nil
		case models.OUT_FOR_DELIVERY.String():
			return models.OUT_FOR_DELIVERY, nil
		case models.PICKED_UP_BY_CUSTOMER.String():
			return models.PICKED_UP_BY_CUSTOMER, nil
		default:
			return 0, models.StatusIsNotExist
		}
	case models.FoodBand.String():
		return compareFoodBandStatuses(externalStatus, previousStatus)
	case models.Poster.String():
		switch externalStatus {
		case "closed":
			return models.ACCEPTED, nil
		default:
			return 0, models.StatusIsNotExist
		}
	case models.Yaros.String():
		switch externalStatus {
		case "accepted":
			return models.ACCEPTED, nil
		case "rejected":
			return models.CANCELLED_BY_POS_SYSTEM, nil
		default:
			return 0, models.StatusIsNotExist
		}
	case models.CTMax.String():
		switch externalStatus {
		case "ACCEPTED":
			return models.ACCEPTED, nil
		case "COOKING_STARTED":
			return models.COOKING_STARTED, nil
		case "COOKING_COMPLETE":
			return models.COOKING_COMPLETE, nil
		case "CLOSED":
			return models.CLOSED, nil
		default:
			return 0, models.StatusIsNotExist
		}
	case models.Kwaaka.String():
		switch externalStatus {
		case "ACCEPTED":
			return models.ACCEPTED, nil
		case "COOKING_STARTED":
			return models.COOKING_STARTED, nil
		case "COOKING_COMPLETE":
			return models.COOKING_COMPLETE, nil
		case "CLOSED":
			return models.CLOSED, nil

		default:
			return 0, models.StatusIsNotExist
		}
	default:
		return 0, models.PosSystemIsIncorrect
	}
}

func (manager OrderManager) validateOrderStatusQueue(exStatus, posType string, posStatus models.PosStatus) bool {

	if posType != models.Kwaaka.String() {
		return true
	}

	queue := make(map[string]int)
	statuses := []string{models.NEW.String(), models.ACCEPTED.String(), models.COOKING_STARTED.String(), models.COOKING_COMPLETE.String(), models.CLOSED.String()}

	for i, status := range statuses {
		queue[status] = i
	}

	return queue[exStatus] < queue[posStatus.String()]
}

func (manager OrderManager) SetPaidStatus(ctx context.Context, orderID string) error {
	return manager.orderRepo.SetPaidStatus(ctx, orderID)
}

func isNecessaryUpdateOrderStatus(order models.Order, store coreStoreModels.Store) bool {
	if order.DeliveryService != models.WOLT.String() {
		return true
	}
	for _, v := range order.StatusesHistory {
		histStatus := ConvertPosStatusToAggregator(v.Name, order, store)
		currStatus := ConvertPosStatusToAggregator(order.Status, order, store)
		if histStatus == currStatus || currStatus == "" {
			return false
		}
	}
	return true
}

func ConvertPosStatusToAggregator(posStatus string, order models.Order, store coreStoreModels.Store) string {

	var statusConfig []coreStoreModels.Status
	switch order.Type {
	case models.Instant:
		statusConfig = store.Wolt.PurchaseTypes.Instant
	case models.Preorder:
		statusConfig = store.Wolt.PurchaseTypes.Preorder
	}
	if order.IsPickedUpByCustomer {
		statusConfig = store.Wolt.PurchaseTypes.TakeAway
	}
	for _, status := range statusConfig {
		if status.PosStatus == posStatus {
			return status.Status
		}
	}

	switch posStatus {
	case models.ACCEPTED.String(), models.WAIT_COOKING.String(), models.READY_FOR_COOKING.String(), models.COOKING_STARTED.String(), models.WAIT_SENDING.String():
		if order.Type == models.Preorder {
			return models.Confirm.String()
		}
		return models.Accept.String()
	case models.COOKING_COMPLETE.String(), models.CLOSED.String(), models.READY_FOR_PICKUP.String(), models.ON_WAY.String(), models.DELIVERED.String(), models.OUT_FOR_DELIVERY.String():
		return models.Ready.String()
	case models.CANCELLED_BY_POS_SYSTEM.String():
		return models.Reject.String()
	case models.PICKED_UP_BY_CUSTOMER.String():
		return models.Delivered.String()
	default:
		return ""
	}

}

func (s *OrderManager) setCustomerEmail(email string) string {
	customerEmail := email
	if customerEmail == "" {
		customerEmail = models.Default3plCustomerEmail
	}

	return customerEmail
}

func (s *OrderManager) setCustomerPhoneNumber(phoneNumber string) string {
	customerPhone := phoneNumber
	if customerPhone == "" {
		customerPhone = models.Default3plCustomerPhone
	}

	return customerPhone
}

func (s *OrderManager) setCustomerName(name string) string {
	customerName := name
	if customerName == "" {
		customerName = models.Default3plCustomerName
	}

	return customerName
}

func isHaniRestDelivery(storeID, deliveryService string) bool {
	if deliveryService == models.QRMENU.String() || deliveryService == models.KWAAKA_ADMIN.String() {
		haniRestIDs := []string{"6683fc9339a3222785df695f", "6683fd3f0077b538b9497c24", "6691122b5aafae72c5da35cb", "669112e0076d30366ac63add", "669113bdcb3c0bfc666a8335"}

		for _, haniRestID := range haniRestIDs {
			if haniRestID == storeID {
				return true
			}
		}
	}

	return false
}

func haniRestAddDeliveryProduct(order models.Order, menu coreMenuModels.Menu) models.Order {
	productMap := make(map[string]coreMenuModels.Product)

	deliveryProductID := "c668b9c1-2f16-433d-98fb-4391fbd1414e"

	if order.EstimatedTotalPrice.Value < 10000 {
		deliveryProductID = "7efa7198-257b-40a7-aa02-9cdddeeb2a09"
	}

	for _, product := range menu.Products {
		productMap[product.ExtID] = product
	}

	if _, ok := productMap[deliveryProductID]; !ok {
		return order
	}

	order.Products = append(order.Products, models.OrderProduct{
		ID:   productMap[deliveryProductID].ExtID,
		Name: productMap[deliveryProductID].Name[0].Value,
		Price: models.Price{
			Value: productMap[deliveryProductID].Price[0].Value,
		},
		Quantity: 1,
	})

	order.EstimatedTotalPrice.Value += productMap[deliveryProductID].Price[0].Value

	return order
}

func (manager OrderManager) GetOrdersForAutoCloseCron(ctx context.Context, query selector.Order) ([]models.Order, error) {

	orders, err := manager.orderRepo.GetAllOrders(ctx, query)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
