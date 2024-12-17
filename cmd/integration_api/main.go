package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/getsentry/sentry-go"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/kwaaka-team/orders-core/cmd"
	"github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/core/config"
	mongo2 "github.com/kwaaka-team/orders-core/core/database/drivers/mongo"
	deliverooManagers "github.com/kwaaka-team/orders-core/core/deliveroo/managers"
	"github.com/kwaaka-team/orders-core/core/externalapi/database"
	"github.com/kwaaka-team/orders-core/core/externalapi/database/drivers"
	externalManagers "github.com/kwaaka-team/orders-core/core/externalapi/managers"
	foodBandManagers "github.com/kwaaka-team/orders-core/core/foodband/managers"
	glovoManagers "github.com/kwaaka-team/orders-core/core/glovo/managers"
	externalPosIntegrationManagers "github.com/kwaaka-team/orders-core/core/integration_api/managers"
	"github.com/kwaaka-team/orders-core/core/integration_api/repository"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1"
	jowiManagers "github.com/kwaaka-team/orders-core/core/jowi/managers"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
	iikoManagers "github.com/kwaaka-team/orders-core/core/service/iiko/managers"
	starterAppManagers "github.com/kwaaka-team/orders-core/core/starter_app/managers"
	talabatManagers "github.com/kwaaka-team/orders-core/core/talabat/manager"
	woltManagers "github.com/kwaaka-team/orders-core/core/wolt/managers"
	firebase_client "github.com/kwaaka-team/orders-core/pkg/firebase"
	"github.com/kwaaka-team/orders-core/pkg/menu"
	menuModels "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	"github.com/kwaaka-team/orders-core/pkg/notify"
	"github.com/kwaaka-team/orders-core/pkg/order"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp/clients"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	"github.com/kwaaka-team/orders-core/service/aws_s3"
	"github.com/kwaaka-team/orders-core/service/bitrix"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	errorSolutionsRepo "github.com/kwaaka-team/orders-core/service/error_solutions/repository"
	"github.com/kwaaka-team/orders-core/service/gourmet"
	"github.com/kwaaka-team/orders-core/service/kwaaka_3pl"
	"github.com/kwaaka-team/orders-core/service/legal_entity_payment"
	legalEntitySrv "github.com/kwaaka-team/orders-core/service/legalentity"
	"github.com/kwaaka-team/orders-core/service/legalentity/models"
	legalEntityRepo "github.com/kwaaka-team/orders-core/service/legalentity/repository"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	orderServicePkg "github.com/kwaaka-team/orders-core/service/order"
	"github.com/kwaaka-team/orders-core/service/order/delivery"
	"github.com/kwaaka-team/orders-core/service/order_report"
	"github.com/kwaaka-team/orders-core/service/order_rules"
	paymentServicePkg "github.com/kwaaka-team/orders-core/service/payment"
	paymentRepository "github.com/kwaaka-team/orders-core/service/payment/repository"
	"github.com/kwaaka-team/orders-core/service/pos"
	posRepository "github.com/kwaaka-team/orders-core/service/pos/repository"
	"github.com/kwaaka-team/orders-core/service/promo_code"
	promoCodeRepo "github.com/kwaaka-team/orders-core/service/promo_code/repository"
	userPromoCodeRepo "github.com/kwaaka-team/orders-core/service/promo_code/user_repository"
	"github.com/kwaaka-team/orders-core/service/refund"
	"github.com/kwaaka-team/orders-core/service/restaurant_set"
	"github.com/kwaaka-team/orders-core/service/shaurma_food"
	"github.com/kwaaka-team/orders-core/service/sms"
	"github.com/kwaaka-team/orders-core/service/stoplist"
	storeServicePkg "github.com/kwaaka-team/orders-core/service/store"
	storeGroupServicePkg "github.com/kwaaka-team/orders-core/service/storegroup"
	wppService "github.com/kwaaka-team/orders-core/service/whatsapp"
	repository2 "github.com/kwaaka-team/orders-core/service/whatsapp/repository"
	"github.com/kwaaka-team/orders-core/service/whatsapp_business"
	_ "github.com/lib/pq"
	lumigotracer "github.com/lumigo-io/lumigo-go-tracer"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

func run() error {
	log.Println("RUNNING INTEGRATION API")

	encoderCfg := zap.NewProductionConfig()
	encoderCfg.EncoderConfig.TimeKey = "timestamp"
	encoderCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncoderConfig.StacktraceKey = ""
	l, err := encoderCfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	logger := l.Sugar()
	defer logger.Sync()

	ctx := context.Background()
	opts, err := general.LoadConfig(ctx)
	if err != nil {
		return err
	}

	if err = sentry.Init(sentry.ClientOptions{
		Dsn:              opts.IntegrationDSN,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		SendDefaultPII:   true,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					logger.Info("request", req)
				}
			}
			return event
		},
		Debug:            false,
		AttachStacktrace: true,
	}); err != nil {
		return err
	}

	session := cmd.GetSession()

	s3Service := aws_s3.NewS3Service(session)

	cognitoSvc := cognitoidentityprovider.New(session)

	orderCli, err := order.NewClient()
	if err != nil {
		return err
	}
	defer orderCli.Close(ctx)

	storeCli, err := store.NewClient(storeModels.Config{
		MongoCli: orderCli.Client(),
	})
	if err != nil {
		return err
	}

	menuCli, err := menu.New(menuModels.Config{
		MongoCli: orderCli.Client(),
	})
	if err != nil {
		return err
	}

	externalDataStore, err := database.New(drivers.DataStoreConfig{
		URL:           opts.DSURL,
		DataStoreName: opts.DSName,
		DataBaseName:  opts.DSDB,
	})

	if err != nil {
		return err
	}

	if err = externalDataStore.Connect(); err != nil {
		return err
	}

	sqsCli := notifyQueue.NewSQS(sqs.NewFromConfig(opts.AwsConfig))
	notifyCli, err := notify.New()
	if err != nil {
		return err
	}

	glovoManager := glovoManagers.NewOrder(orderCli, storeCli)
	woltManager := woltManagers.NewOrder(orderCli, storeCli, opts.WoltConfiguration.BaseURL)
	deliverooManager := deliverooManagers.NewEvent(storeCli, orderCli, menuCli)
	externalOrderManager := externalManagers.NewOrderClientManager(orderCli, storeCli)
	externalMenuManager := externalManagers.NewMenuClientManager(menuCli, storeCli)
	externalAuthManager := externalManagers.NewAuthClientManager(externalDataStore, storeCli, opts.AppSecret, opts.EmenuConfiguration)
	talabatOrderManager := talabatManagers.NewOrder(orderCli, storeCli, logger)
	talabatMenuManager := talabatManagers.NewMenu(menuCli, logger)
	starterAppOrderManager := starterAppManagers.NewOrder(orderCli)

	foodBandMenuManager := foodBandManagers.NewMenuManager(menuCli, storeCli, opts.AwsSession, logger)
	foodBandStoreManager := foodBandManagers.NewStoreManager(storeCli, logger)
	foodBandOrderManager := foodBandManagers.NewOrderManager(orderCli, logger)

	jowiManager := jowiManagers.NewJowiManager(storeCli, orderCli)

	externalPosIntegrationRepository, err := repository.NewExternalPosIntegrationAuthRepository(repository.DBInfo{
		Driver:      opts.DSName,
		MongoClient: orderCli.Client(),
		DBName:      opts.DSDB,
	})
	if err != nil {
		return err
	}

	externalPosIntegrationManager := externalPosIntegrationManagers.NewExternalPosIntegrationManager(orderCli, menuCli, storeCli, externalPosIntegrationRepository)

	ds := orderCli.GetDataStore().Client().Database(opts.DSDB)

	cfg := orderCli.GetConfig()

	telegramService, err := orderServicePkg.NewTelegramService(sqsCli, cfg.QueConfiguration.Telegram, opts.NotificationConfiguration, &telegram.Repository{})
	if err != nil {
		return err
	}

	whatsappService, err := whatsapp.NewWhatsappClient(&clients.Config{
		Instance:  cfg.WhatsAppConfiguration.Instance,
		AuthToken: cfg.WhatsAppConfiguration.AuthToken,
		BaseURL:   cfg.WhatsAppConfiguration.BaseUrl,
		Insecure:  true,
		Protocol:  "http",
	})
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Addr,
		Username: cfg.RedisConfig.Username,
		Password: cfg.RedisConfig.Password,
	})

	storeService, stopListService, aggFactory, posFactory, orderRepo, menuService, storeGroupService, paymentFactory, customerRepo, subscriptionRepo, paymentRepo, kwaaka3plService, orderRuleService, orderReport, cartService, restaurantSetService, refundRepo, errorSolutionService, err := createServices(
		ds,
		cfg,
		s3Service,
		sqsCli,
		telegramService,
		logger,
		whatsappService,
		orderCli,
		menuCli,
		cognitoSvc,
		opts.IntegrationBaseURL,
	)
	if err != nil {
		return err
	}

	posterStoreAuthRepo, err := posRepository.NewPosterStoreAuthsMongoRepository(ds)
	if err != nil {
		return err
	}

	iikoManager := iikoManagers.NewEvent(storeCli, menuCli, notifyCli, orderCli, sqsCli,
		opts.IIKOConfiguration.BaseURL,
		opts.RetryConfiguration.QueueName,
		opts.RetryConfiguration.Count,
		opts.IntegrationBaseURL,
		opts.QueConfiguration.OfflineOrders,
		stopListService,
		storeService,
		telegramService,
		errorSolutionService,
	)

	var paymentService paymentServicePkg.Service

	publisher := &orderServicePkg.Publisher{}
	subscriber3PL, err := kwaaka_3pl.NewSubscriber3PL(kwaaka3plService)
	if err != nil {
		return err
	}
	publisher.AddSubscriber(subscriber3PL)

	posSender, err := orderServicePkg.NewPosSender(cfg, menuCli, storeCli, menuService, storeService, orderRepo, orderCli.GetDataStore(), posFactory, errorSolutionService, stopListService, telegramService, sqsCli)
	if err != nil {
		return err
	}

	orderServiceImpl, err := createOrderService(cfg, menuCli, storeCli, storeService, aggFactory, posFactory, orderRepo, menuService, storeGroupService, publisher, posSender, orderRuleService, paymentRepo, cartService, errorSolutionService)
	if err != nil {
		return err
	}

	emptyPosServiceImpl, err := orderServicePkg.NewEmptyPosService(storeService, aggFactory, orderRepo)
	if err != nil {
		return err
	}

	var orderService orderServicePkg.CreationService = orderServiceImpl

	orderService, err = orderServicePkg.NewErrorNotificationDecorator(orderServiceImpl, storeService, telegramService, whatsappService)
	if err != nil {
		return err
	}

	var emptyPosService orderServicePkg.CreationService = emptyPosServiceImpl

	emptyPosService, err = orderServicePkg.NewTelegramServiceDecorator(emptyPosService, &telegramService, storeService, nil, coreOrderModels.Kwaaka)
	if err != nil {
		return err
	}

	emptyPosService, err = orderServicePkg.NewWhatsappServiceDecorator(emptyPosService, whatsappService, storeService, coreOrderModels.Kwaaka)
	if err != nil {
		return err
	}

	firebaseMsgService, err := createFirebaseMsgService(ctx, s3Service, opts.FirebaseConfiguration)
	if err != nil {
		return err
	}
	emptyPosService, err = orderServicePkg.NewFirebaseServiceDecorator(emptyPosService, storeCli, firebaseMsgService)
	if err != nil {
		return err
	}

	orderService, err = orderServicePkg.NewCreateOrderRouter(storeService, orderService, emptyPosService)
	if err != nil {
		return err
	}

	orderService, err = orderServicePkg.NewTelegramServiceDecorator(orderService, &telegramService, storeService, []coreOrderModels.Aggregator{coreOrderModels.QRMENU, coreOrderModels.KWAAKA_ADMIN}, coreOrderModels.RKeeper7XML)
	if err != nil {
		return err
	}

	orderService, err = orderServicePkg.NewWhatsappServiceDecorator(orderService, whatsappService, storeService, coreOrderModels.RKeeper)
	if err != nil {
		return err
	}

	paymentService, err = createPaymentService(paymentFactory, customerRepo, subscriptionRepo, paymentRepo, sqsCli, cfg.QueueUrls.PaymentsQueueUrl, logger, storeService, storeGroupService, refundRepo)
	if err != nil {
		return err
	}

	var orderCronService orderServicePkg.OrderCronService = orderServiceImpl
	var statusUpdateService orderServicePkg.StatusUpdateService = orderServiceImpl
	var orderReviewService orderServicePkg.ReviewService = orderServiceImpl
	var orderInfoSharingService orderServicePkg.InfoSharingService = orderServiceImpl
	var orderCancellationService orderServicePkg.CancellationService = orderServiceImpl

	posterService, err := pos.NewPosterService(nil, cfg.PosterConfiguration.BaseURL, "", posterStoreAuthRepo, storeCli, menuCli, cfg.ApplicationID, cfg.ApplicationSecret, cfg.RedirectURI, logger)
	if err != nil {
		return err
	}

	legalEntityRepo := legalEntityRepo.NewLegalEntityRepo(ds.Client().Database(opts.DSDB))
	legalEntityService := legalEntitySrv.NewLegalEntityService(models.S3Info{
		KwaakaFilesBucket: opts.KwaakaFilesBucket, KwaakaFilesBaseUrl: opts.KwaakaFilesBaseUrl,
	}, legalEntityRepo, s3Service)

	shaurmaFoodService, err := shaurma_food.NewService(storeService, menuService, posFactory, storeGroupService)
	if err != nil {
		return err
	}

	wppBusinessService, err := whatsapp_business.NewWppBusinessService(cfg.WhatsappBusinessConfiguration, redisClient)
	if err != nil {
		return err
	}

	newsletterRepo := repository2.NewNewsletterRepository(ds)

	wppService, err := wppService.NewWhatsappService(whatsappService, cfg.WhatsAppConfiguration.Instance, cfg.WhatsAppConfiguration.AuthToken, cfg.WhatsAppConfiguration.BaseUrl, newsletterRepo, storeService, orderCli, storeGroupService, redisClient)
	if err != nil {
		return err
	}

	db, err := newPostgresDB(opts.LegalEntityPaymentCfg.DBCfg)
	if err != nil {
		return err
	}
	legalEntityPaymentRepo, err := legal_entity_payment.NewRepository(db)
	if err != nil {
		return err
	}
	legalEntityPaymentService, err := legal_entity_payment.NewService(opts, legalEntityPaymentRepo, s3Service, legalEntityService, wppService)
	if err != nil {
		return err
	}

	promoCodeRepo, err := promoCodeRepo.NewMongoRepository(ds.Client().Database(opts.DSDB))
	if err != nil {
		return err
	}

	userPromoCodes, err := userPromoCodeRepo.NewMongoRepository(ds.Client().Database(opts.DSDB))
	if err != nil {
		return err
	}

	promoCodeService, err := promo_code.NewPromoCodeService(logger, promoCodeRepo, userPromoCodes, storeService, menuService)
	if err != nil {
		return err
	}

	smsService, err := sms.NewSmsService(opts.SmsLogin, opts.SmsPassword, redisClient, storeGroupService)
	if err != nil {
		return err
	}

	bitrixService := bitrix.NewBitrixService(wppService)

	gourmetService, err := gourmet.NewServiceImpl(storeService, opts.IIKOConfiguration.BaseURL)
	if err != nil {
		return err
	}

	server := v1.NewServer(orderService, orderReviewService, menuService, posFactory, statusUpdateService, orderCronService, kwaaka3plService, storeService, stopListService, storeGroupService, glovoManager, woltManager, deliverooManager,
		externalOrderManager, externalMenuManager, externalAuthManager, talabatOrderManager, talabatMenuManager, starterAppOrderManager, iikoManager, posterService, foodBandMenuManager, foodBandOrderManager, foodBandStoreManager, externalPosIntegrationManager,
		paymentService, jowiManager, opts, logger, cmd.IsLambda(), legalEntityPaymentService, telegramService, orderInfoSharingService, orderCancellationService, shaurmaFoodService, wppBusinessService, wppService, promoCodeService, orderReport,
		cartService, smsService, bitrixService, restaurantSetService, gourmetService)

	if cmd.IsLambda() {
		wrappedHandler := lumigotracer.WrapHandler(server.GinProxy, &lumigotracer.Config{})
		lambda.Start(wrappedHandler)
		logger.Info("that log for checking order of executing of functions after start")
	} else {
		if err = http.ListenAndServe(":8080", server.Router); err != nil {
			logger.Fatalf("listen and serve error")
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
}

func createPaymentService(paymentFactory *paymentServicePkg.PaymentSystemFactory,
	customerRepo paymentRepository.CustomersRepository,
	subscriptionRepo paymentRepository.SubscriptionsRepository,
	paymentRepo paymentRepository.PaymentsRepository,
	notifyQueue notifyQueue.SQSInterface, paymentsQueueUrl string,
	logger *zap.SugaredLogger,
	storeService storeServicePkg.Service,
	storeGroupService storeGroupServicePkg.Service,
	refundRepo refund.Repository) (paymentServicePkg.Service, error) {
	paymentService, err := paymentServicePkg.NewService(paymentFactory, customerRepo, paymentRepo, subscriptionRepo, notifyQueue, paymentsQueueUrl, logger, storeService, storeGroupService, refundRepo)
	if err != nil {
		return nil, err
	}

	return paymentService, nil
}

func createOrderService(
	globalConfig config.Configuration,
	menuCli menu.Client, storeCli store.Client,
	storeService storeServicePkg.Service,
	aggFactory aggregator.Factory,
	posFactory pos.Factory,
	orderRepo orderServicePkg.Repository,
	menuService *menuServicePkg.Service,
	storeGroupService storeGroupServicePkg.Service,
	publisher *orderServicePkg.Publisher,
	posSender orderServicePkg.PosSender,
	orderRuleService order_rules.Service,
	paymentRepo paymentRepository.PaymentsRepository,
	cartService orderServicePkg.CartService,
	errSolutionService error_solutions.Service,
) (*orderServicePkg.ServiceImpl, error) {

	sf := orderServicePkg.ServiceFactory{
		StoreService:      storeService,
		AggregatorFactory: aggFactory,
		PosFactory:        posFactory,
		Repository:        orderRepo,
		GlobalConfig:      &globalConfig,
		MenuClient:        menuCli,
		StoreClient:       storeCli,
		MenuService:       menuService,
		StoreGroupService: storeGroupService,
		Publisher:         publisher,
		PosSender:         posSender,
		OrderRuleService:  orderRuleService,
		PaymentRepo:       paymentRepo,
		CartService:       cartService,
		ErrSolution:       errSolutionService,
	}
	orderService, err := sf.Create()
	if err != nil {
		return nil, err
	}
	return orderService, nil
}

func createServices(db *mongo.Database, cfg config.Configuration, s3Service aws_s3.Service, sqsCli notifyQueue.SQSInterface, telegramService orderServicePkg.TelegramService, logger *zap.SugaredLogger, whatsapp clients.Whatsapp, orderCli order.Client, menuCli menu.Client, cognito *cognitoidentityprovider.CognitoIdentityProvider, ocBaseUrl string) (storeServicePkg.Service,
	stoplist.Service,
	aggregator.Factory,
	pos.Factory,
	orderServicePkg.Repository,
	*menuServicePkg.Service,
	storeGroupServicePkg.Service,
	*paymentServicePkg.PaymentSystemFactory,
	paymentRepository.CustomersRepository,
	paymentRepository.SubscriptionsRepository,
	paymentRepository.PaymentsRepository,
	kwaaka_3pl.Service,
	order_rules.Service,
	order_report.OrderReport,
	orderServicePkg.CartService,
	restaurant_set.Service,
	refund.Repository,
	error_solutions.Service,
	error,
) {

	var anotherBillRepository pos.AnotherBillRepository
	anotherBillRepository, err := pos.NewMongoAnotherBillRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	storeRepository, err := storeServicePkg.NewStoreMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	storeFactory, err := storeServicePkg.NewService(storeRepository)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	var menuRepo menuServicePkg.Repository
	menuRepo, err = menuServicePkg.NewMenuMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	menuService, err := menuServicePkg.NewMenuService(menuRepo, storeFactory, s3Service, sqsCli)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	storeGroupRepository, err := storeGroupServicePkg.NewMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	stopListService, err := stoplist.CreateStopListServiceByWebhook(db, cfg, 1)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	storeGroupService, err := storeGroupServicePkg.NewService(storeGroupRepository)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	aggFactory, err := aggregator.NewFactory(
		cfg.WoltConfiguration.BaseURL,
		cfg.GlovoConfiguration.BaseURL, cfg.GlovoConfiguration.Token,
		cfg.TalabatConfiguration.MiddlewareBaseURL, cfg.TalabatConfiguration.MenuBaseUrl,
		cfg.Express24Configuration.BaseURL, cfg.StarterAppConfiguration.BaseUrl, menuService, cognito, cfg,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	bkOfferRepository, err := mongo2.NewBKOfferRepository2(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	posFactory, err := pos.NewFactory(
		anotherBillRepository, sqsCli, cfg.RetryConfiguration.QueueName,
		cfg.IIKOConfiguration.BaseURL, cfg.IIKOConfiguration.TransportToFrontTimeout, cfg.PosterConfiguration.BaseURL,
		cfg.PalomaConfiguration.BaseURL, cfg.PalomaConfiguration.Class, cfg.JowiConfiguration.BaseURL,
		cfg.JowiConfiguration.ApiKey, cfg.JowiConfiguration.ApiSecret, cfg.RKeeperConfiguration.RKeeperBaseURL, cfg.RKeeperConfiguration.RKeeperApiKey,
		cfg.BurgerKingConfiguration.BaseURL, bkOfferRepository, cfg.RKeeper7XMLConfiguration.LicenseBaseURL,
		cfg.SyrveConfiguration.BaseURL, cfg.YarosConfiguration.BaseURL, cfg.YarosConfiguration.InfoSystem, cfg.TillypadConfiguration.BaseUrl, cfg.Ytimes.BaseUrl, cfg.Ytimes.Token, cfg.PosistConfiguration.BaseUrl,
	)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	var orderRepo orderServicePkg.Repository
	orderRepo, err = orderServicePkg.NewMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	var deliveryRepo delivery.Repository
	deliveryRepo, err = delivery.NewDeliveryMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	newsletterRepo := repository2.NewNewsletterRepository(db)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Addr,
		Username: cfg.RedisConfig.Username,
		Password: cfg.RedisConfig.Password,
	})

	wppService, err := wppService.NewWhatsappService(whatsapp, cfg.WhatsAppConfiguration.Instance, cfg.WhatsAppConfiguration.AuthToken, cfg.WhatsAppConfiguration.BaseUrl, newsletterRepo, storeFactory, orderCli, storeGroupService, redisClient)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	cartRepo := orderServicePkg.NewCartRepository(db)
	cartService := orderServicePkg.NewCartService(cartRepo)

	paymentFactory, err := paymentServicePkg.NewFactory(cfg.IokaConfiguration.BaseUrl, cfg.IokaConfiguration.ApiKey, cfg.PaymeConfiguration.BaseUrl, cfg.PaymeConfiguration.ApiKey, cfg.WoopPayConfiguration.BaseUrl, cfg.WoopPayConfiguration.ResultUrl, wppService, cfg.KaspiSaleScoutConfiguration.BaseUrl, cfg.KaspiSaleScoutConfiguration.Token, cfg.KaspiSaleScoutConfiguration.MerchantID, logger, ocBaseUrl, cartService)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	customerRepo, err := paymentRepository.NewCustomersMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	subscriptionRepo, err := paymentRepository.NewSubscriptionsMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	paymentsRepo, err := paymentRepository.NewPaymentsMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	kwaaka3pl, err := kwaaka_3pl.NewKwaaka3plService(sqsCli, cfg.Kwaaka3pl.Kwaaka3plQueue, orderRepo, storeFactory, cfg.Kwaaka3pl.Kwaaka3plBaseUrl, cfg.Kwaaka3pl.Kwaaka3plAuthToken, logger, telegramService, menuCli, deliveryRepo)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	orderRuleRepo, err := order_rules.NewOrderRulesMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	orderRuleService, err := order_rules.NewOrderRuleService(orderRuleRepo)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	orderReport := order_report.NewOrderReportService(storeFactory, orderRepo, kwaaka3pl, cartService, storeGroupService)

	restaurantSetRepo, err := restaurant_set.NewMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	restaurantSetService, err := restaurant_set.NewService(restaurantSetRepo, storeGroupService)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	refundRepo, err := refund.NewMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	errorSolutionRepo, err := errorSolutionsRepo.NewMongoRepository(db)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	errorSolutionsService, err := error_solutions.NewErrorSolutionService(errorSolutionRepo)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
	}

	return storeFactory, stopListService, aggFactory, posFactory, orderRepo, menuService, storeGroupService, paymentFactory, customerRepo, subscriptionRepo, paymentsRepo, kwaaka3pl, orderRuleService, orderReport, cartService, restaurantSetService, refundRepo, errorSolutionsService, nil
}

func createFirebaseMsgService(ctx context.Context, s3Service aws_s3.Service, opts general.FirebaseConfiguration) (*firebase_client.MessageService, error) {
	configBytes, err := getFirebaseConfigsFromS3(s3Service, opts)
	if err != nil {
		return nil, err
	}

	configs := option.WithCredentialsJSON(configBytes)
	firebaseMsgService, err := firebase_client.NewFirebaseMessageService(ctx, configs)
	if err != nil {
		return nil, err
	}
	return firebaseMsgService, nil
}

func getFirebaseConfigsFromS3(s3Service aws_s3.Service, opts general.FirebaseConfiguration) ([]byte, error) {
	return s3Service.GetObject(opts.S3BucketName, opts.S3FileKey)
}

func newPostgresDB(cfg general.LegalEntityPaymentDB) (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Println("can't connect to legal entity payment DB")
	}

	if err = db.Ping(); err != nil {
		log.Println("can't connect to legal entity payment DB")
	}

	return db, nil
}
