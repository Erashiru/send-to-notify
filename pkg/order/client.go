package order

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/config/general"
	mongo2 "github.com/kwaaka-team/orders-core/core/database/drivers/mongo"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/kwaaka-team/orders-core/service/kwaaka_3pl"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/order"
	orderServicePkg "github.com/kwaaka-team/orders-core/service/order"
	"github.com/kwaaka-team/orders-core/service/order/delivery"
	"github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/stoplist"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"

	"github.com/kwaaka-team/orders-core/core/models/selector"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	"github.com/kwaaka-team/orders-core/core/managers"
	menuClient "github.com/kwaaka-team/orders-core/pkg/menu"
	menuCoreModel "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	"github.com/kwaaka-team/orders-core/pkg/order/dto"
	notifyClient "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	storeCoreModel "github.com/kwaaka-team/orders-core/pkg/store/dto"
	errorSolutionsRepo "github.com/kwaaka-team/orders-core/service/error_solutions/repository"
	telegramSvc "github.com/kwaaka-team/orders-core/service/order"
)

type Client interface {
	Client() *mongo.Client
	Close(ctx context.Context) error
	GetOrder(ctx context.Context, order dto.OrderSelector) (models.Order, error)
	GetOrdersWithFilters(ctx context.Context, query dto.OrderSelector) ([]models.Order, int, error)

	GetActiveOrders(ctx context.Context, query dto.ActiveOrderSelector) ([]models.Order, error)
	GetActivePreorders(ctx context.Context, query dto.ActiveOrderSelector) ([]models.Order, error)
	UpdateOrder(ctx context.Context, order models.Order) error
	CancelOrder(ctx context.Context, order models.CancelOrder) error
	UpdateOrderStatus(ctx context.Context, posOrderID, pos, status, errorDescription string) error
	UpdateOrderStatusInDS(ctx context.Context, orderID string, posStatus dto.PosStatus) error
	UpdateOrderStatusByID(ctx context.Context, orderID, pos, status string) error

	UpdateOrderModel(ctx context.Context, query dto.OrderSelector, order dto.UpdateOrder) error
	CreateOrderInDB(ctx context.Context, req models.Order) (models.Order, error)
	CreateOrderInPOS(ctx context.Context, req models.Order) (models.Order, error)

	SetPaidStatus(ctx context.Context, orderID string) error
	ManualUpdateStatus(ctx context.Context, req models.Order) error
	GetAllOrders(ctx context.Context, query dto.OrderSelector) ([]models.Order, error)
	CancelOrderInPos(ctx context.Context, order models.CancelOrderInPos) error

	GetDirectCallCenterActivePreorders(ctx context.Context, query dto.ActiveOrderSelector) ([]models.Order, error)
	GetHaniKarimaActivePreOrders(ctx context.Context, query dto.ActiveOrderSelector) ([]models.Order, error)

	GetOrdersForAutoCloseCron(ctx context.Context, autoCloseTime int, storeID string) ([]models.Order, error)
}

type OrderCoreClient struct {
	orderManager    managers.Order
	sqsCli          notifyClient.SQSInterface
	storeCli        storeClient.Client
	ds              drivers.DataStore
	globalConfig    config.Configuration
	errSolutions    error_solutions.Service
	stopListService stoplist.Service
	telegramService telegramSvc.TelegramService
}

func NewClient() (*OrderCoreClient, error) {
	opts, err := config.LoadConfig(context.Background())
	if err != nil {
		return nil, err
	}

	ds, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSDB,
	})
	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionConfig()
	encoderCfg.EncoderConfig.TimeKey = "timestamp"
	encoderCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncoderConfig.StacktraceKey = ""
	l, err := encoderCfg.Build()
	if err != nil {
		return nil, err
	}
	logger := l.Sugar()
	defer logger.Sync()

	if err != nil {
		return nil, fmt.Errorf("cannot create datastore %s: %v", opts.DSName, err)
	}

	if err = ds.Connect(ds.Client()); err != nil {
		return nil, fmt.Errorf("cannot connect to datastore: %s", err)
	}

	menuCli, err := menuClient.New(menuCoreModel.Config{
		Region:    opts.Region,
		SecretEnv: os.Getenv(models.SECRET_ENV),
		MongoCli:  ds.Client(),
	})
	if err != nil {
		return nil, fmt.Errorf("menu-core cli err %s", err)
	}

	storeCli, err := storeClient.NewClient(storeCoreModel.Config{
		Region:    opts.Region,
		SecretEnv: os.Getenv(models.SECRET_ENV),
		MongoCli:  ds.Client(),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot initialize Store Client")

	}

	sqsCli := notifyClient.NewSQS(sqs.NewFromConfig(opts.AwsConfig))

	mongoClient := ds.Client()
	db := mongoClient.Database(opts.DSDB)

	storeFactory, aggFactory, posFactory, orderRepo, deliveryRepo, err := CreateServices(db, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize Store Client")
	}

	notificationConfig := general.NotificationConfiguration{
		TelegramChatID:    opts.TelegramChatID,
		TelegramChatToken: opts.TelegramChatToken,
		YandexErrorChatID: opts.YandexErrorChatID,
		OrderBotToken:     opts.OrderBotToken,
	}
	telegramRepo := telegram.NewTelegramRepo(db.Client().Database(opts.DSDB))
	telegramService, err := orderServicePkg.NewTelegramService(sqsCli, opts.QueConfiguration.Telegram, notificationConfig, telegramRepo)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize telegram service %v", err)
	}

	kwaaka3pl, err := kwaaka_3pl.NewKwaaka3plService(sqsCli, opts.Kwaaka3plQueue, orderRepo, storeFactory, opts.Kwaaka3pl.Kwaaka3plBaseUrl, opts.Kwaaka3pl.Kwaaka3plAuthToken, logger, telegramService, menuCli, deliveryRepo)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize kwaaka 3pl client: %v", err)
	}

	orderManager, err := managers.NewOrderManager(ds, opts, menuCli, storeCli, sqsCli, storeFactory, aggFactory, posFactory, orderRepo, kwaaka3pl, telegramService)
	if err != nil {
		return nil, err
	}

	errSolutionRepo, err := errorSolutionsRepo.NewMongoRepository(db)
	if err != nil {
		return nil, err
	}
	errSolutionService, err := error_solutions.NewErrorSolutionService(errSolutionRepo)
	if err != nil {
		return nil, err
	}

	stopListService, err := stoplist.CreateStopListServiceByWebhook(db, opts, 1)
	if err != nil {
		return nil, err
	}

	return &OrderCoreClient{
		orderManager:    orderManager,
		sqsCli:          sqsCli,
		storeCli:        storeCli,
		ds:              ds,
		globalConfig:    opts,
		errSolutions:    errSolutionService,
		stopListService: stopListService,
		telegramService: telegramService,
	}, nil
}

func CreateServices(db *mongo.Database, globalConfig config.Configuration) (*store.ServiceImpl, aggregator.Factory, pos.Factory, order.Repository, delivery.Repository, error) {
	storeRepository, err := store.NewStoreMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	menuRepository, err := menuServicePkg.NewMenuMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	menuService, err := menuServicePkg.NewMenuService(menuRepository, nil, nil, nil)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	storeFactory, err := store.NewService(storeRepository)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	aggFactory, err := aggregator.NewFactory(
		globalConfig.WoltConfiguration.BaseURL,
		globalConfig.GlovoConfiguration.BaseURL, globalConfig.GlovoConfiguration.Token,
		globalConfig.TalabatConfiguration.MiddlewareBaseURL, globalConfig.TalabatConfiguration.MenuBaseUrl,
		globalConfig.Express24Configuration.BaseURL, globalConfig.StarterAppConfiguration.BaseUrl, menuService, nil, config.Configuration{},
	)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	var anotherBillRepository pos.AnotherBillRepository
	anotherBillRepository, err = pos.NewMongoAnotherBillRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	bkOfferRepository, err := mongo2.NewBKOfferRepository2(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	sqsCli := notifyClient.NewSQS(sqs.NewFromConfig(globalConfig.AwsConfig))

	posFactory, err := pos.NewFactory(
		anotherBillRepository, sqsCli, globalConfig.RetryConfiguration.QueueName,
		globalConfig.IIKOConfiguration.BaseURL, globalConfig.IIKOConfiguration.TransportToFrontTimeout,
		globalConfig.PosterConfiguration.BaseURL, globalConfig.PalomaConfiguration.BaseURL, globalConfig.PalomaConfiguration.Class,
		globalConfig.JowiConfiguration.BaseURL, globalConfig.JowiConfiguration.ApiKey, globalConfig.JowiConfiguration.ApiSecret,
		globalConfig.RKeeperBaseURL, globalConfig.RKeeperApiKey, globalConfig.BurgerKingConfiguration.BaseURL, bkOfferRepository,
		globalConfig.RKeeper7XMLConfiguration.LicenseBaseURL, globalConfig.SyrveConfiguration.BaseURL, globalConfig.YarosConfiguration.BaseURL, globalConfig.YarosConfiguration.InfoSystem, globalConfig.TillypadConfiguration.BaseUrl, globalConfig.Ytimes.BaseUrl, globalConfig.Ytimes.Token, globalConfig.PosistConfiguration.BaseUrl,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	var orderRepo order.Repository
	orderRepo, err = order.NewMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	var deliveryRepo delivery.Repository
	deliveryRepo, err = delivery.NewDeliveryMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return storeFactory, aggFactory, posFactory, orderRepo, deliveryRepo, nil

}

func (client *OrderCoreClient) GetConfig() config.Configuration {
	return client.globalConfig
}

func (client *OrderCoreClient) GetDataStore() drivers.DataStore {
	return client.ds
}

func (client OrderCoreClient) Client() *mongo.Client {
	return client.ds.Client()
}

func (client OrderCoreClient) Close(ctx context.Context) error {
	return client.ds.Close(ctx)
}

func (client OrderCoreClient) ManualUpdateStatus(ctx context.Context, req models.Order) error {
	return client.orderManager.ManualUpdateStatus(ctx, req)
}

func (client OrderCoreClient) UpdateOrderModel(ctx context.Context, query dto.OrderSelector, req dto.UpdateOrder) error {
	order, err := client.orderManager.GetOrder(ctx, selector.EmptyOrderSearch().
		SetID(query.ID).
		SetOrderID(query.OrderID).
		SetPosOrderID(query.PosOrderID).
		SetDeliveryService(query.DeliveryService))
	if err != nil {
		return err
	}

	if req.PosID != "" {
		order.PosOrderID = req.PosID
	}

	if !req.ReadingTime.Value.IsZero() {
		order.ReadingTime = append(order.ReadingTime, models.TransactionTime{
			Value:        req.ReadingTime.Value,
			UTCOffset:    req.ReadingTime.UTCOffset,
			TimeZone:     req.ReadingTime.TimeZone,
			RestaurantID: req.RestaurantID,
		})
	}

	if err := client.orderManager.UpdateOrder(ctx, order); err != nil {
		return err
	}

	return nil
}

func (client OrderCoreClient) CreateOrderInDB(ctx context.Context, req models.Order) (models.Order, error) {
	res, err := client.orderManager.CreateOrderDB(ctx, req)
	if err != nil {
		return models.Order{}, err
	}

	return res, nil
}

func (client OrderCoreClient) CreateOrderInPOS(ctx context.Context, req models.Order) (models.Order, error) {

	var (
		product models2.Product
	)

	res, err := client.orderManager.CreateOrderInPOS(ctx, req)
	if err == nil {
		return res, nil
	}
	errorMessage := err.Error()

	st, err := client.storeCli.FindStore(ctx, storeCoreModel.StoreSelector{
		ID: req.RestaurantID,
	})
	if err != nil {
		return models.Order{}, err
	}

	errorSolutions, err := client.errSolutions.GetAllErrorSolutions(ctx)
	if err != nil {
		return models.Order{}, err
	}

	errSolutionByCode, addToStopListStatus, err := client.errSolutions.GetErrorSolutionByCode(ctx, st, pos.MatchingCodes(errorMessage, errorSolutions))
	if err != nil {
		return models.Order{}, err
	}

	req.FailReason = models.FailReason{
		Code:         errSolutionByCode.Code,
		Message:      errorMessage,
		Reason:       errSolutionByCode.Reason,
		BusinessName: errSolutionByCode.BusinessName,
		Solution:     errSolutionByCode.Solution,
	}
	req.Status = models.FAILED.String()

	var productID string
	if addToStopListStatus {
		productID = pos.GetProductIDFromRegexp(errorMessage, errSolutionByCode)
		productErrorCodes := map[string]bool{
			"21": true, "4": true, "1": true,
		}
		attributeErrorCodes := map[string]bool{
			"25": true, "5": true, "7": true, "21": true, "27": true, "28": true,
		}
		if len(productID) > 0 {
			var err error
			switch {
			case productErrorCodes[errSolutionByCode.Code]:
				err = client.stopListService.UpdateStopListByPosProductID(ctx, false, st.ID, productID)
				if err == nil {
					log.Info().Msgf("successfully put product with id: %s to stop with error solution code: %s for store_id : %s", productID, errSolutionByCode.Code, st.ID)
				}
				for _, orderProduct := range req.Products {
					if orderProduct.ID == productID {
						product.ExtID = orderProduct.ID
						product.Name = append(product.Name, models2.LanguageDescription{Value: orderProduct.Name})
					}
				}
			case attributeErrorCodes[errSolutionByCode.Code]:
				err = client.stopListService.UpdateStopListByAttributeID(ctx, false, st.ID, productID)
				if err == nil {
					log.Info().Msgf("successfully put attribute with id: %s to stop with error solution code: %s for store_id : %s", productID, errSolutionByCode.Code, st.ID)
				}
				for _, orderProduct := range req.Products {
					if len(orderProduct.Attributes) > 0 {
						for _, orderAttribute := range orderProduct.Attributes {
							if orderAttribute.ID == productID {
								product.ExtID = orderAttribute.ID
								product.Name = append(product.Name, models2.LanguageDescription{Value: orderAttribute.Name})
							}
						}
					}
				}
			default:
				return models.Order{}, fmt.Errorf("unsupported error code to update stoplist bu pos product/attribute id : %s", errSolutionByCode.Code)
			}
			if err != nil {
				return models.Order{}, err
			}
			log.Info().Msgf("send stoplist status true, for product/attribute id: %s, store_id:%s, error solution code:%s", productID, st.ID, errSolutionByCode.Code)

			if errSolutionByCode.SendToTelegram {
				if err := client.telegramService.SendMessageToQueue(telegram.PutProductToStopListWithErrSolution, req, st, req.FailReason.BusinessName, "", "", product); err != nil {
					return models.Order{}, err
				}
			}
		}
	}

	if err := client.UpdateOrder(ctx, req); err != nil {
		return models.Order{}, err
	}
	log.Info().Msgf("OrderCoreClient: update order_id: %s with code: %s and message:%s", req.OrderID, req.FailReason.Code, req.FailReason.Message)

	return models.Order{}, err
}

func (client OrderCoreClient) GetOrder(ctx context.Context, query dto.OrderSelector) (models.Order, error) {

	order, err := client.orderManager.GetOrder(ctx, selector.EmptyOrderSearch().
		SetID(query.ID).
		SetPosOrderID(query.PosOrderID).
		SetDeliveryService(query.DeliveryService).
		SetOrderID(query.OrderID).
		SetIsParentOrder(query.IsParentOrder).
		SetStoreID(query.StoreID))

	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (client OrderCoreClient) GetActiveOrders(ctx context.Context, query dto.ActiveOrderSelector) ([]models.Order, error) {

	orders, err := client.orderManager.GetActiveOrders(ctx, selector.EmptyOrderSearch().
		SetPosType(query.PosType).
		SetStoreID(query.StoreID).
		SetDeliveryService(query.DeliveryService).
		SetOrderCode(query.OrderCode))

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (client OrderCoreClient) GetDirectCallCenterActivePreorders(ctx context.Context, query dto.ActiveOrderSelector) ([]models.Order, error) {
	orders, err := client.orderManager.GetDirectCallCenterActivePreorders(ctx, selector.EmptyOrderSearch().
		SetPosType(query.PosType).
		SetStoreID(query.StoreID).
		SetDeliveryService(query.DeliveryService).
		SetOrderCode(query.OrderCode))

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (client OrderCoreClient) GetHaniKarimaActivePreOrders(ctx context.Context, query dto.ActiveOrderSelector) ([]models.Order, error) {
	orders, err := client.orderManager.GetHaniKarimaActivePreOrders(ctx, selector.EmptyOrderSearch().
		SetPosType(query.PosType).
		SetStoreID(query.StoreID).
		SetDeliveryService(query.DeliveryService).
		SetOrderCode(query.OrderCode))

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (client OrderCoreClient) UpdateOrder(ctx context.Context, order models.Order) error {
	return client.orderManager.UpdateOrder(ctx, order)
}

func (client OrderCoreClient) CancelOrder(ctx context.Context, order models.CancelOrder) error {
	return client.orderManager.CancelOrder(ctx, order)
}

func (client OrderCoreClient) UpdateOrderStatus(ctx context.Context, posOrderID, pos, status, errorDescription string) error {
	return client.orderManager.UpdateOrderStatus(ctx, posOrderID, pos, status, errorDescription)
}

func (client OrderCoreClient) UpdateOrderStatusByID(ctx context.Context, orderID, pos, status string) error {
	return client.orderManager.UpdateOrderStatusByID(ctx, orderID, pos, status)
}

func (client OrderCoreClient) UpdateOrderStatusInDS(ctx context.Context, orderID string, status dto.PosStatus) error {
	return client.orderManager.UpdateOrderStatusInDS(ctx, orderID, models.PosStatus(status))
}

func (client OrderCoreClient) GetOrdersWithFilters(ctx context.Context, query dto.OrderSelector) ([]models.Order, int, error) {
	orders, total, err := client.orderManager.GetOrdersWithFilters(ctx, selector.Order{
		Pagination: selector.Pagination{
			Limit: query.Limit,
			Page:  query.Page,
		},
		Sorting: selector.Sorting{
			Param:     query.Param,
			Direction: query.Dir,
		},
		ID:                   query.ID,
		OrderCode:            query.OrderCode,
		DeliveryService:      query.DeliveryService,
		Restaurants:          query.Restaurants,
		OrderTimeTo:          query.OrderTimeTo,
		OrderTimeFrom:        query.OrderTimeFrom,
		OnlyActive:           query.OnlyActive,
		Status:               query.Status,
		CustomerNumber:       query.Customer.PhoneNumber,
		ExternalStoreID:      query.ExternalStoreID,
		StoreID:              query.StoreID,
		PosType:              query.PosType,
		IsPickedUpByCustomer: query.IsPickedUpByCustomer,
		DeliveryArray:        query.DeliveryArray,
	})

	if err != nil {
		return nil, 0, err
	}

	return orders, total, err
}

func (client OrderCoreClient) GetActivePreorders(ctx context.Context, query dto.ActiveOrderSelector) ([]models.Order, error) {

	orders, err := client.orderManager.GetActivePreorders(ctx, selector.EmptyOrderSearch().
		SetPosType(query.PosType).
		SetStoreID(query.StoreID).
		SetDeliveryService(query.DeliveryService))

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (client OrderCoreClient) GetAllOrders(ctx context.Context, query dto.OrderSelector) ([]models.Order, error) {
	orders, err := client.orderManager.GetAllOrders(ctx, selector.EmptyOrderSearch().
		SetOrderCode(query.OrderCode))
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (client OrderCoreClient) SetPaidStatus(ctx context.Context, orderID string) error {
	return client.orderManager.SetPaidStatus(ctx, orderID)
}

func (client OrderCoreClient) CancelOrderInPos(ctx context.Context, order models.CancelOrderInPos) error {
	return client.orderManager.CancelOrderInPos(ctx, order)
}

func (client OrderCoreClient) GetOrdersForAutoCloseCron(ctx context.Context, autoCloseTime int, storeID string) ([]models.Order, error) {

	if autoCloseTime == 0 {
		autoCloseTime = 60
	}

	orders, err := client.orderManager.GetAllOrders(ctx, selector.EmptyOrderSearch().
		SetRestaurants([]string{storeID}).SetEstimatedPickupTimeTo(time.Now().UTC().Add(-(time.Duration(autoCloseTime) * time.Minute))).
		SetCreatedAtTimeFrom(time.Now().UTC().Add(-(time.Duration(autoCloseTime)*time.Minute)-30*time.Minute)).
		SetCookingCompleteClosedStatus(&[]bool{true}[0]))
	if err != nil {
		return nil, err
	}

	return orders, nil
}
