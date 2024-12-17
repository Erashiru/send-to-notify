package v1

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	ginAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/config/general"
	deliverooManagers "github.com/kwaaka-team/orders-core/core/deliveroo/managers"
	externalManagers "github.com/kwaaka-team/orders-core/core/externalapi/managers"
	foodBandManagers "github.com/kwaaka-team/orders-core/core/foodband/managers"
	glovoManagers "github.com/kwaaka-team/orders-core/core/glovo/managers"
	externalPosIntegrationManagers "github.com/kwaaka-team/orders-core/core/integration_api/managers"
	jowiManagers "github.com/kwaaka-team/orders-core/core/jowi/managers"
	iikoManagers "github.com/kwaaka-team/orders-core/core/service/iiko/managers"
	starterAppManagers "github.com/kwaaka-team/orders-core/core/starter_app/managers"
	talabatManagers "github.com/kwaaka-team/orders-core/core/talabat/manager"
	woltManagers "github.com/kwaaka-team/orders-core/core/wolt/managers"
	"github.com/kwaaka-team/orders-core/service/bitrix"
	"github.com/kwaaka-team/orders-core/service/gourmet"
	"github.com/kwaaka-team/orders-core/service/kwaaka_3pl"
	"github.com/kwaaka-team/orders-core/service/legal_entity_payment"
	menuServicePkg "github.com/kwaaka-team/orders-core/service/menu"
	"github.com/kwaaka-team/orders-core/service/order"
	"github.com/kwaaka-team/orders-core/service/order_report"
	"github.com/kwaaka-team/orders-core/service/payment"
	posService "github.com/kwaaka-team/orders-core/service/pos"
	"github.com/kwaaka-team/orders-core/service/promo_code"
	"github.com/kwaaka-team/orders-core/service/restaurant_set"
	"github.com/kwaaka-team/orders-core/service/shaurma_food"
	"github.com/kwaaka-team/orders-core/service/sms"
	"github.com/kwaaka-team/orders-core/service/stoplist"
	storeServicePkg "github.com/kwaaka-team/orders-core/service/store"
	storeGroupServicePkg "github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/kwaaka-team/orders-core/service/whatsapp"
	"github.com/kwaaka-team/orders-core/service/whatsapp_business"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

const (
	storePath     = "restaurant_id"
	orderPath     = "order_id"
	remoteId      = "remoteId"
	remoteOrderId = "remoteOrderId"
)

var ginLambda *ginAdapter.GinLambda

type Server struct {
	Router                        *gin.Engine
	orderService                  order.CreationService
	orderReviewService            order.ReviewService
	menuService                   *menuServicePkg.Service
	posService                    posService.Factory
	statusUpdateService           order.StatusUpdateService
	orderCronService              order.OrderCronService
	orderKwaaka3plService         kwaaka_3pl.Service
	storeService                  storeServicePkg.Service
	stopListService               stoplist.Service
	storeGroupService             storeGroupServicePkg.Service
	glovoManager                  glovoManagers.Order
	woltManager                   woltManagers.Order
	deliverooManager              deliverooManagers.Event
	externalOrderManager          externalManagers.OrderClient
	externalMenuManager           externalManagers.MenuClient
	externalAuthManager           externalManagers.AuthClient
	talabatOrderManager           talabatManagers.Order
	talabatMenuManager            talabatManagers.Menu
	starterAppOrderManager        starterAppManagers.Order
	iikoManager                   iikoManagers.Event
	posterService                 *posService.PosterService
	foodBandMenuManager           foodBandManagers.Menu
	foodBandOrderManager          foodBandManagers.Order
	foodBandStoreManager          foodBandManagers.Store
	externalPosIntegrationManager externalPosIntegrationManagers.ExternalPosIntegrationClient
	paymentManager                payment.Service
	jowiManager                   jowiManagers.JowiManager
	Config                        general.Configuration
	Logger                        *zap.SugaredLogger
	LegalEntityPaymentService     legal_entity_payment.Service
	TelegramService               order.TelegramService
	orderInfoSharingService       order.InfoSharingService
	orderCancellationService      order.CancellationService
	shaurmaFoodService            *shaurma_food.Service
	WppBusinessService            whatsapp_business.Service
	WhatsappService               whatsapp.Service
	PromoCode                     promo_code.Service
	orderReport                   order_report.OrderReport
	cartService                   order.CartService
	SmsService                    sms.Service
	bitrixService                 bitrix.Service
	restaurantSetService          restaurant_set.Service
	gourmetService                *gourmet.ServiceImpl
}

func NewServer(
	orderService order.CreationService,
	orderReviewService order.ReviewService,
	menuService *menuServicePkg.Service,
	posService posService.Factory,
	statusUpdateService order.StatusUpdateService,
	orderCronService order.OrderCronService,
	orderKwaaka3plService kwaaka_3pl.Service,
	storeService storeServicePkg.Service,
	stopListService stoplist.Service,
	storeGroupService storeGroupServicePkg.Service,
	glovoManager glovoManagers.Order,
	woltManager woltManagers.Order,
	deliverooManager deliverooManagers.Event,
	externalOrderManager externalManagers.OrderClient,
	externalMenuManager externalManagers.MenuClient,
	externalAuthManager externalManagers.AuthClient,
	talabatOrderManager talabatManagers.Order,
	talabatMenuManager talabatManagers.Menu,
	starterAppOrderManager starterAppManagers.Order,
	iikoManager iikoManagers.Event,
	posterService *posService.PosterService,
	foodBandMenuManager foodBandManagers.Menu,
	foodBandOrderManager foodBandManagers.Order,
	foodBandStoreManager foodBandManagers.Store,
	externalPosIntegrationManager externalPosIntegrationManagers.ExternalPosIntegrationClient,
	paymentManager payment.Service,
	jowiManager jowiManagers.JowiManager,
	config general.Configuration,
	logger *zap.SugaredLogger,
	isProduction bool,
	legalEntityPaymentService legal_entity_payment.Service,
	telegramService order.TelegramService,
	orderInfoSharingService order.InfoSharingService,
	orderCancellationService order.CancellationService,
	shaurmaFoodService *shaurma_food.Service,
	wppBusinessService whatsapp_business.Service,
	whatsappService whatsapp.Service,
	promoCodeService promo_code.Service,
	orderReport order_report.OrderReport,
	cartService order.CartService,
	SmsService sms.Service,
	bitrixService bitrix.Service,
	restaurantSetService restaurant_set.Service,
	gourmetService *gourmet.ServiceImpl,
) *Server {

	server := &Server{
		Router:                        gin.Default(),
		orderService:                  orderService,
		orderReviewService:            orderReviewService,
		menuService:                   menuService,
		posService:                    posService,
		statusUpdateService:           statusUpdateService,
		orderCronService:              orderCronService,
		orderKwaaka3plService:         orderKwaaka3plService,
		storeService:                  storeService,
		storeGroupService:             storeGroupService,
		stopListService:               stopListService,
		glovoManager:                  glovoManager,
		woltManager:                   woltManager,
		deliverooManager:              deliverooManager,
		externalOrderManager:          externalOrderManager,
		externalMenuManager:           externalMenuManager,
		externalAuthManager:           externalAuthManager,
		talabatMenuManager:            talabatMenuManager,
		talabatOrderManager:           talabatOrderManager,
		starterAppOrderManager:        starterAppOrderManager,
		iikoManager:                   iikoManager,
		posterService:                 posterService,
		foodBandMenuManager:           foodBandMenuManager,
		foodBandOrderManager:          foodBandOrderManager,
		foodBandStoreManager:          foodBandStoreManager,
		externalPosIntegrationManager: externalPosIntegrationManager,
		jowiManager:                   jowiManager,
		Config:                        config,
		Logger:                        logger,
		paymentManager:                paymentManager,
		LegalEntityPaymentService:     legalEntityPaymentService,
		TelegramService:               telegramService,
		orderInfoSharingService:       orderInfoSharingService,
		orderCancellationService:      orderCancellationService,
		shaurmaFoodService:            shaurmaFoodService,
		WppBusinessService:            wppBusinessService,
		WhatsappService:               whatsappService,
		PromoCode:                     promoCodeService,
		orderReport:                   orderReport,
		cartService:                   cartService,
		SmsService:                    SmsService,
		bitrixService:                 bitrixService,
		restaurantSetService:          restaurantSetService,
		gourmetService:                gourmetService,
	}

	ginLambda = ginAdapter.New(server.Router)
	server.register(server.Router, isProduction)

	server.Router.RedirectTrailingSlash = true
	server.Router.RedirectFixedPath = true
	server.Router.HandleMethodNotAllowed = true

	return server
}

func (server *Server) register(engine *gin.Engine, isProduction bool) {
	engine.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	engine.Use(server.sentryMiddleware())

	api := engine.Group("/api")
	{
		api.POST("/update-order-status", server.UpdateOrderStatusByPosTypes)
		api.POST("/update-stoplist", server.UpdateStopListByPosTypes)
		api.POST("/update-stoplist-by-section", server.UpdateStopListBySection)

		api.POST("/generate-new-aggregator-menu", server.GenerateNewAggregatorMenu)
		api.POST("/auto-update-aggregator-menu", server.AutoUpdateAggregatorMenu)

		api.GET("/get-stat", server.GetOrderStat)

		api.GET("/get-wolt-csv", server.GetWoltMenuInCsv)

		api.POST("/set-markup-to-aggregator-menu", server.SetMarkUpToAggregatorMenu)
		api.POST("/generate-aggregator-menu-from-pos-menu", server.GenerateAggregatorMenuFromPosMenu)

		api.POST("/3pl/create-delivery", server.CreateDelivery3plCron)
		api.POST("/payment/send-in-whatsapp", server.SendPayment)

		api.POST("/kaspi-salescout/create-order", server.KaspiSaleScoutCronEvent)
		api.POST("/report/delivery-dispatcher-price", server.DeliveryDispatcherPrice)
		api.POST("/send-defer-orders", server.SendDeferOrders)
		api.POST("/bitrix/event", server.SendMessageToCustomerFromBitrixLead)
		api.POST("/delivery/no-dispatcher-message", server.NoDispatcherMessage)
		api.POST("/delivery/performer-lookup-time", server.PerformerLookupMoreThan15Minute)

		api.POST("/send-telegram-message", server.SendTelegramMessage)
	}

	posIIKO := engine.Group("/iiko")
	posIIKO.Use(server.iikoProductsSecretMiddleware())

	{
		posIIKO.POST("/events", server.EventIIKO)
	}

	posSyrve := engine.Group("/syrve")
	posSyrve.Use(server.iikoProductsSecretMiddleware())

	{
		posSyrve.POST("/events", server.EventSyrve)
	}

	ytimes := engine.Group("/ytimes")
	{
		ytimes.Use(server.ytimesSecretMiddleware())
		ytimes.POST("/remote-order/status", server.YTimesUpdateOrderStatus)
		ytimes.POST("/menu/changed", server.YTimesMenuUpdates)
	}

	jowi := engine.Group("/jowi")
	{
		jowi.POST("/events", server.JowiEvents)
	}

	poster := engine.Group("/poster")

	{
		poster.GET("/code", server.CodeReceiverHandlerPoster)
		poster.POST("/events", server.WebHookEventsHandlerPoster)
	}

	// v1
	v1 := engine.Group("/v1")

	{

		// endpoint for external integration
		integration := v1.Group("/integration")
		{
			// endpoint for external integration POS
			posIntegration := integration.Group("/pos")

			swagger := posIntegration.Group("/swagger")

			posIntegration.POST("/stoplist", server.ExternalPosIntegrationStopList)
			posIntegration.POST("/order", server.ExternalPosIntegrationUpdateOrder)
			posIntegration.GET("/:restaurant_id/orders", server.ExternalPosIntegrationGetOrders)

			swagger.GET("/*any", ginSwagger.WrapHandler(
				swaggerFiles.Handler,
				ginSwagger.URL("doc.json"),
			))
		}
	}

	{
		stores := v1.Group("/stores")

		stores.GET("", server.GetStoresFoodBand)
		authStores := stores.Group(fmt.Sprintf("/:%s", "store_id"))
		authStores.Use(server.Auth)
		{
			authStores.POST(fmt.Sprintf("/manage/:%s", "delivery_service"), server.ManageAggregatorStoreFoodBand)
			authStores.POST(fmt.Sprintf("/menus/:%s", "delivery_service"), server.UploadMenuFoodBand)
			authStores.GET(fmt.Sprintf("/menus/:%s/:%s", "delivery_service", "transaction_id"), server.GetMenuUploadStatusFoodBand)
			authStores.POST(fmt.Sprintf("/orders/:%s/status", "order_id"), server.UpdateOrderStatusFoodBand)
			authStores.POST(fmt.Sprintf("/products/:%s", "product_id"), server.StopListProductFoodBand)
			authStores.POST(fmt.Sprintf("/attributes/:%s", "attribute_id"), server.StopListAttributeFoodBand)
		}
	}

	{
		posYaros := v1.Group("/yaros")
		posYaros.Use(YarosSecretMiddleware(server.Config.YarosConfiguration.Token))
		{
			posYaros.PATCH("/order-update", server.OrderUpdateYaros)
			posYaros.PATCH("/stoplist-update", server.StoplistUpdateYaros)
		}
	}

	{
		glovo := v1.Group("/glovo")
		glovo.Use(server.secretMiddleware(server.Config.GlovoConfiguration.Token))
		{
			glovo.POST("/placeOrder", server.CreateOrderGlovo)
			glovo.POST("/cancelOrder", server.CancelOrderGlovo)
		}

		wolt := v1.Group("/wolt")

		{
			wolt.POST("/placeOrder", server.CreateOrderWolt)
		}

		deliveroo := v1.Group("/deliveroo")
		//deliveroo.Use(SecretMiddleware(server.opts.Token))
		{
			deliveroo.POST("/order-events", server.OrderEventDeliveroo)
			deliveroo.POST("/menu-events", server.MenuEventDeliveroo)

		}

		express24 := v1.Group("/express24")
		express24.Use(server.express24SecretMiddleware(server.Config.Express24Configuration.Token))
		{
			express24.POST("/order-receive", server.ReceiveOrder)
		}

		{
			oauth := v1.Group("/security/oauth")
			{
				oauth.POST("/token", server.CreateToken)
			}

			external := v1.Group("")

			external.Use(server.authorizeJWT(server.Config.AppSecret))

			{
				orders := external.Group("/orders")
				{
					orders.POST("", server.GetOrders)
				}

				nomenclature := external.Group("/nomenclature")
				{
					nomenclature.GET(fmt.Sprintf("/:%s/composition", storePath), server.GetRetailMenu)
					nomenclature.GET(fmt.Sprintf("/:%s/availability", storePath), server.GetRemains)
				}

				menu := external.Group("/menu")
				{
					menu.GET(fmt.Sprintf("/:%s/composition", storePath), server.GetMenu)
					menu.GET(fmt.Sprintf("/:%s/availability", storePath), server.GetAvailability)
					menu.GET(fmt.Sprintf("/:%s/promos", storePath), server.GetPromos)
				}

				restaurant := external.Group("/restaurants")
				{
					restaurant.GET("", server.GetRestaurants)
				}

				order := external.Group("/order")
				{
					order.POST("", server.CreateOrder)
					order.GET(fmt.Sprintf("/:%s", orderPath), server.GetOrder)
					order.DELETE(fmt.Sprintf("/:%s", orderPath), server.CancelOrder)
					order.PUT(fmt.Sprintf("/:%s", orderPath), server.UpdateOrder)
					order.GET(fmt.Sprintf("/:%s/status", orderPath), server.GetOrderStatus)
				}
			}

		}

		talabat := v1.Group("/talabat")
		{
			talabat.POST("requestResult", server.MenuUploadCallbackTalabat)
			talabat.POST(fmt.Sprintf("/remoteId/:%s/remoteOrder/:%s/posOrderStatus", remoteId, remoteOrderId), server.CancelOrderTalabat)
			talabat.POST(fmt.Sprintf("/order/:%s", remoteId), server.CreateOrderTalabat)
		}

		starterApp := v1.Group("/starterapp")
		starterApp.Use(server.StarterAppAuth)
		{
			starterApp.POST("/order", server.CreateOrderStarterApp)
		}

	}
	{
		whatsappPayment := v1.Group("/whatsapp")
		{
			whatsappPayment.POST("/webhooks", server.WhatsappWebhooks)
		}
	}
	{
		ioka := v1.Group("/ioka")
		{
			ioka.POST("/webhooks", server.IokaWebhookEvent)
		}
	}
	{
		multicard := v1.Group("/multicard")
		{
			multicard.POST("/webhooks", server.MulticardWebhook)
		}
	}
	{
		payme := v1.Group("/payme")
		{
			payme.POST("/webhooks", server.PaymeWebhookEvent)
		}
	}
	{
		woopPay := v1.Group("/wooppay")
		{
			woopPay.POST("/webhooks", server.WoopPayWebhookEvent)
		}
	}
	{
		qrMenu := v1.Group("/qr-menu")
		qrMenu.GET("/notify-unpaid-customers", server.NotifyAllUnpaidCustomers)
		qrMenu.Use(server.secretMiddleware(server.Config.KwaakaQrMenuToken))
		{
			qrMenu.POST("/placeOrder", server.CreateOrderQRmenu)
			qrMenu.POST("/applePay/session", server.OpenApplePaySessionByQRmenu)
			qrMenu.POST("/applePay/createPayment/:payment_order_id", server.CreateApplePayPayment)
			qrMenu.GET("/twogis-review-link/:restaurant_id", server.GetTwoGisReviewLink)
			wppBusiness := qrMenu.Group("/wpp-business")
			{
				wppBusiness.POST("/send-verification-code", server.SendVerificationCode)
			}
			sms := qrMenu.Group("/sms")
			{
				sms.POST("/send-verification-code", server.SendVerificationCodeBySms)
			}
			promoCode := qrMenu.Group("/promo-code")
			{
				promoCode.POST("/validate-promo-code", server.ValidatePromoCodeForUser)
				promoCode.POST("/add-usage-time", server.AddUserPromoCodeUsageTimeToDB)
			}
			restaurantSet := qrMenu.Group("/restaurant-set")
			{
				restaurantSet.GET("/get/:restaurant_set_id", server.GetRestaurantSetById)
				restaurantSet.POST("/create", server.CreateRestaurantSet)
				restaurantSet.GET("/get-with-rest-group/:restaurant_set_id", server.GetRestaurantSetWithRestGroup)
				restaurantSet.GET("/get-by-domain-name", server.GetRestaurantSetByDomainName)
			}
		}
	}
	{
		kwaakaAdmin := v1.Group("/kwaaka-admin")
		kwaakaAdmin.Use(server.secretMiddleware(server.Config.KwaakaAdminToken))
		kwaakaAdmin.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowCredentials: true,
			AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Phone", "Usergroup"},
		}))
		{
			kwaakaAdmin.POST("/placeOrder", server.CreateOrderKwaakaAdmin)
			kwaakaAdmin.DELETE("/cancelOrder/:order_id", server.CancelOrderKwaakaAdmin)
			kwaakaAdmin.POST("/setOrdersDispatcher", server.SetOrdersDispatcher)
			kwaakaAdmin.PUT("/cancelOrder", server.CancelOrderDispatcher)
			kwaakaAdmin.POST("/courier-search-cancel/:delivery_order_id", server.CancelCourierSearch)
			kwaakaAdmin.POST("/save-3pl-history/:delivery_order_id", server.Save3plHistory)
			kwaakaAdmin.POST("/delivery-3pl-info", server.Get3plDeliveryInfo)
			kwaakaAdmin.POST("/map-iiko-statuses-to-3pl-statuses", server.MapIIKOstatusTo3plStatus)

			kwaakaAdmin.POST("/stoplist/product", server.KwaakaAdminStopListByProductID)
			kwaakaAdmin.POST("/stoplist/attribute", server.KwaakaAdminStopListByAttributeID)
			kwaakaAdmin.GET("/get_all_stores/:restaurant_group_id", server.GetRestaurantsByGroupId)
			kwaakaAdmin.GET("/:restaurant_group_id", server.GetStoresInRestaurantGroupByQuery)
			kwaakaAdmin.GET("/get-order-by-delivery-id/:delivery_id", server.GetCustomerByDeliveryId)
			kwaakaAdmin.GET("/get-order-for-telegram-message/:delivery_id", server.GetOrderForTelegramByDeliveryOrderId)
			kwaakaAdmin.POST("/bulk-create-order/:order_id", server.BulkCreate3plOrder)
			kwaakaAdmin.POST("/busy-mode", server.UpdateKwaakaAdminBusyMode)
			kwaakaAdmin.POST("/delivery-change-history", server.InsertChangeDeliveryHistory)

			kwaakaAdmin.GET("/get-discounts/:store_id", server.GetCustomerDiscount)
			kwaakaAdmin.POST("/get-discount-history/:store_id", server.GetDiscountHistory)
			kwaakaAdmin.GET("/get-restaurant-discounts/:store_id", server.GetDiscountsForStore)

			kwaakaAdmin.POST("/upsert-menu", server.UpsertMenu)

			kwaakaAdmin.POST("/:restaurant-id/customer-data", server.CreateStorePhoneEmail)

			kwaakaAdmin.POST("/polygon", server.CreatePolygon)
			kwaakaAdmin.PUT("/polygon", server.UpdatePolygon)
			kwaakaAdmin.GET("/polygon/:restaurant_id", server.GetPolygonByRestaurantID)

			kwaakaAdmin.POST("/payment-order", server.CreatePaymentOrder)
			kwaakaAdmin.POST("/payment-order/create-payment-link/:order_id", server.CreatePaymentLink)
			kwaakaAdmin.PUT("/payment-order", server.UpdatePaymentOrderStatus)
			kwaakaAdmin.GET("/set-seq-number/:restaurant_id", server.SetActualSeqNumber)
			kwaakaAdmin.POST("/order-report/restaurant", server.OrdersReportForRestaurant)
			kwaakaAdmin.POST("/order-report/restaurant/totals", server.OrderReportForRestaurantTotals)
			kwaakaAdmin.POST("/order-report/kwaaka/totals", server.OrderReportForKwaakaTotals)
			kwaakaAdmin.POST("/order-report/xlsx", server.OrderReportToXlsx)
			kwaakaAdmin.POST("/refund/:order_id", server.RefundPayment)
			kwaakaAdmin.GET("/get-refund/:order_id", server.GetRefund)

			kwaakaAdmin.POST("/menu/:menu_id/product/:product_id/name", server.AddNameInProduct)
			kwaakaAdmin.POST("/menu/:menu_id/product/:product_id/description", server.AddDescriptionInProduct)
			kwaakaAdmin.POST("/menu/:menu_id/section/:section_id/name", server.AddNameInSection)
			kwaakaAdmin.POST("/menu/:menu_id/section/:section_id/description", server.AddDescriptionInSection)
			kwaakaAdmin.POST("/menu/:menu_id/attribute-group/:attribute_group_id/name", server.AddNameInAttributeGroup)
			kwaakaAdmin.POST("/menu/:menu_id/attribute/:attribute_id/name", server.AddNameInAttribute)
			kwaakaAdmin.PATCH("/menu/:menu_id/product/:product_id/name", server.ChangeNameInProduct)
			kwaakaAdmin.PATCH("/menu/:menu_id/product/:product_id/description", server.ChangeDescriptionInProduct)
			kwaakaAdmin.PATCH("/menu/:menu_id/section/:section_id/name", server.ChangeNameInSection)
			kwaakaAdmin.PATCH("/menu/:menu_id/section/:section_id/description", server.ChangeDescriptionInSection)
			kwaakaAdmin.PATCH("/menu/:menu_id/attribute-group/:attribute_group_id/name", server.ChangeNameInAttributeGroup)
			kwaakaAdmin.PATCH("/menu/:menu_id/attribute/:attribute_id/name", server.ChangeNameInAttribute)
			kwaakaAdmin.POST("/menu/:menu_id/product/:product_id/regulatory-information", server.AddRegulatoryInformation)
			kwaakaAdmin.PATCH("/menu/:menu_id/product/:product_id/regulatory-information", server.ChangeRegulatoryInformation)
		}

		promoCode := kwaakaAdmin.Group("/promo-code")
		{
			promoCode.POST("", server.CreatePromoCode)
			promoCode.PUT("", server.UpdatePromoCode)
			promoCode.GET("/id/:promo-code-id", server.GetPromoCodeByID)
			promoCode.GET("/code/:promo-code", server.GetAvailablePromoCodeByCode)
			promoCode.GET("/:promo-code/restaurant/:restaurant-id", server.GetPromoCodeByCodeAndRestaurantId)
			promoCode.GET("/restaurant/:restaurant-id", server.GetPromoCodesByRestaurantID)
		}

		legalEntityPayment := kwaakaAdmin.Group("/legal-entity-payment")
		{
			legalEntityPayment.POST("/create", server.CreateLegalEntityPayment)
			legalEntityPayment.GET("/:legal_entity_payment_id", server.GetLegalEntityPaymentByID)
			legalEntityPayment.POST("/list", server.GetListLegalEntityPayment)
			legalEntityPayment.PUT("/update", server.UpdateLegalEntityPayment)
			legalEntityPayment.DELETE("/:legal_entity_payment_id", server.DeleteLegalEntityPayment)
			legalEntityPayment.POST("/payment-analytics", server.GetLegalEntityPaymentAnalytics)
			legalEntityPayment.POST("/upload-pdf", server.UploadPDF)
			legalEntityPayment.PATCH("/create-bill", server.CreateBill)
			legalEntityPayment.PATCH("/confirm-payment", server.ConfirmPayment)
		}

		dispatcher := kwaakaAdmin.Group("/dispatcher")
		{
			dispatcher.GET("/customer/phone/:phone/orders", server.GetOrdersByCustomerPhone)
		}

		telegram := kwaakaAdmin.Group("/telegram")
		{
			telegram.POST("/send-message", server.SendTelegramMessage)
			telegram.POST("/send-3pl-message", server.Send3plErrorMsg)
			telegram.POST("/send-compensation-message", server.SendCompensationMessage)
		}
		whatsapp := kwaakaAdmin.Group("/whatsapp")
		{
			whatsapp.POST("/send-newsletter", server.SendNewsletter)
			whatsapp.POST("/send-message", server.SendMessage)
		}
	}
	{
		shaurmaFood := v1.Group("/shaurma-food")
		shaurmaFood.Use(server.secretMiddleware(server.Config.ShaurmaFoodToken))
		shaurmaFood.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowCredentials: true,
			AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Phone", "Usergroup"},
		}))
		{
			shaurmaFood.GET("/orders", server.GetOrdersShaurmaFood)
		}
		{
			shaurmaFood.POST("/set-external-menu-prices-to-aggregator-menus", server.SetExternalMenuPricesToAggregatorMenus)
			shaurmaFood.POST("/restaurant/:restaurant_id/:aggregator/update-images-and-descriptions", server.UpdateImagesAndDescriptionsInAggregatorMenus)
		}
	}

	{
		gourmetApi := v1.Group("/gourmet")
		gourmetApi.Use(server.secretMiddleware(server.Config.GourmetToken))
		gourmetApi.Use(cors.New(cors.Config{
			AllowAllOrigins:  true,
			AllowCredentials: true,
			AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Phone", "Usergroup"},
		}))

		gourmetApi.GET(fmt.Sprintf("/restaurants/:%s/tables", "restaurantId"), server.GourmetGetTables)
		gourmetApi.GET(fmt.Sprintf("/restaurants/:%s/orders", "restaurantId"), server.GourmetGetOrders)
		gourmetApi.POST(fmt.Sprintf("/restaurants/:%s/orders/:%s/pay", "restaurantId", "orderId"), server.GourmetPay)
	}

}

func (server *Server) GinProxy(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request ginproxy: %v\n", req)
	return ginLambda.Proxy(req)
}
