package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	models3 "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"log"
	"strings"
	"time"

	"github.com/kwaaka-team/orders-core/core/models"
)

type Service string

const (
	Telegram      Service = "telegram"
	TimeFormat            = "2006-01-02 15:04:05"
	StoreIsOpened         = "storeIsOpened"
)

var (
	paymentSystemMap = map[string]string{
		"ioka":            "Ioka",
		"kaspi_salescout": "Kaspi Salescout",
		"cash":            "Cash",
		"by_cashier":      "By Cashier",
	}
	deliveryDispatcherMap = map[string]string{
		"yandex":  "Yandex",
		"indrive": "Indrive",
		"wolt":    "Wolt",
	}
)

func (s Service) String() string {
	return string(s)
}

type NotificationType string

const (
	CheckIn                             NotificationType = "check_in"
	SuccessCreateOrder                  NotificationType = "success_create_order"
	CreateOrder                         NotificationType = "create_order"
	CancelOrder                         NotificationType = "cancel_order"
	UpdateOrder                         NotificationType = "update_order"
	StoreClosed                         NotificationType = "store_closed"
	StoreStatusReport                   NotificationType = "store_status_report"
	OrderStat                           NotificationType = "order_statistic"
	OrderStatusChange                   NotificationType = "order_status_change"
	ThirdPartyError                     NotificationType = "third_party_error "
	NoCourier                           NotificationType = "no_courier"
	Refund                              NotificationType = "refund"
	Compensation                        NotificationType = "compensation"
	CancelDeliveryFromDispatcherPage    NotificationType = "cancel_delivery_from_dispatcher_page"
	NoDeliveryDispatcher                NotificationType = "no_delivery_dispatcher"
	AutoUpdatePublicateMenu             NotificationType = "auto_update_publicate_menu"
	PutProductToStopListWithErrSolution NotificationType = "put_product_to_stoplist_with_err_solution"
)

func (s NotificationType) String() string {
	return string(s)
}

var (
	reviewKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("1 - плохо", "review:1")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("2", "review:2")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("3", "review:3")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("4", "review:4")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("5 - отлично", "review:5")),
	)
)

func ConstructOrderReport(service Service, stat models.OrderStat) string {
	var msg string

	switch service {
	case Telegram:
		msg = stat.ConstructedErrMsg

	}
	return msg
}

type Report struct {
	Dividend float64
	Divisor  float64
}

func (r *Report) GetErrorMessage(format string, dividend float64) (string, float64) {
	num, percent := r.GetNumberAndPercent(dividend)
	return fmt.Sprintf(" • <b> %.0f - %.2f%% </b>", num, percent) + format, num
}

func (r *Report) GetNumberAndPercent(dividend float64) (float64, float64) {
	r.Dividend = dividend
	return dividend, r.GetPercentage()
}

func (r *Report) GetPercentage() float64 {
	if r.Divisor == 0 {
		return 0
	}
	return (r.Dividend / r.Divisor) * 100
}

func ConstructRestaurantMessageToNotify(order models.Order, store coreStoreModels.Store, errorInfo string) string {
	var msg string

	localTime, err := order.OrderTime.GetLocalTime()
	if err != nil {
		localTime = order.OrderTime.Value.Time
	}
	orderTime := localTime.Format(TimeFormat)

	paymentMethod := order.PaymentMethod
	if paymentMethod == models.PAYMENT_METHOD_DELAYED {
		paymentMethod = models.PAYMENT_METHOD_CARD
	}

	msg = fmt.Sprintf(
		"<b>[❌] Заказ %s № %s не удался в %s</b>\n"+
			"<b>Ресторан: </b>%s\n"+
			"<b>Тип оплаты: </b>%s\n"+
			"<b>Покупатель: </b>%s\n"+
			"<b>Номер телефона покупателя: </b>%s\n"+
			"<b>Код заказа: </b>%s\n"+
			"<b>Код выдачи: </b>%s\n"+
			"<b>Комментарий к заказу: </b>%s\n"+
			"<b>Дата создания: </b>%s\n"+
			"<b>OrderID заказа: </b>%s\n\n"+ // TODO: For testing purposes, easier to find order
			strings.ToTitle(order.DeliveryService), order.OrderCode, strings.ToTitle(order.PosType), store.Name,
		paymentMethod, order.Customer.Name, order.Customer.PhoneNumber, order.OrderCode, order.PickUpCode, order.AllergyInfo,
		orderTime, order.OrderID,
	)

	for index, product := range order.Products {
		msg += fmt.Sprintf("<b>%d. %s, x%v</b>\n", index+1, product.Name, product.Quantity)

		var modifierBody string
		for position, attribute := range product.Attributes {
			modifierBody += fmt.Sprintf("    %d. %s, x%v\n", position+1, attribute.Name, attribute.Quantity)
		}

		msg += modifierBody + "\n\n"
	}

	format := "<b>Тип ошибки: </b> %s\n" +
		"<b>Причина:</b> %s\n" +
		"<b>Что делать:</b> %s\n"

	errType, reason, solution := utils.GetReasonAndSolution(errorInfo)
	msg += fmt.Sprintf(format, errType, reason, solution)
	return msg
}

func ConstructOrderMessageToNotify(service Service, order models.Order, store coreStoreModels.Store, errorInfo string) string {
	var message string

	switch service {
	case Telegram:
		if order.DeliveryAddress.Label == "" {
			order.DeliveryAddress.Label = "None"
		}
		if order.Customer.PhoneNumber == "" {
			order.Customer.PhoneNumber = "None"
		}
		if order.PosOrderID == "" {
			order.PosOrderID = "None"
		}
		if order.AllergyInfo == "" {
			order.AllergyInfo = "None"
		}

		var orderTime string
		localTime, err := order.OrderTime.GetLocalTime()
		if err != nil {
			orderTime = order.OrderTime.Value.Format("2006-01-02 15:04:05")
		}
		orderTime = localTime.Format("2006-01-02 15:04:05")

		var paymentSystem = order.PaymentSystem
		if val, ok := paymentSystemMap[order.PaymentSystem]; ok {
			paymentSystem = val
		}

		message += fmt.Sprintf("<b>[❌] Заказ %s № %v не удался в %v</b>\n"+
			"<b>Ресторан: </b>%s\n"+
			"<b>Тип Оплаты: </b>%s\n"+
			"<b>Система Оплаты: </b>%s\n"+
			"<b>Order ID: </b>%v\n"+
			"<b>Статус: </b>%v\n"+
			"<b>Покупатель: </b>%v\n"+
			"<b>Номер телефона: </b>%v\n"+
			"<b>Адрес: </b>%v\n"+
			"<b>Код заказа: </b>%v\n"+
			"<b>Код выдачи: </b>%v\n"+
			"<b>Комментарий к заказу: </b>%v\n"+ // ? sure that this field
			"<b>ID в %v: </b>%v\n"+
			"<b>%v Store ID: </b>%v\n"+
			"<b>ID в %v: </b>%v\n"+
			"<b>Дата создания: </b>%v\n\n"+
			"<b>Ошибка: </b>%v\n"+
			"<b>Что делать?: </b>%v\n\n",
			strings.ToTitle(order.DeliveryService),
			order.OrderCode,
			strings.ToUpper(store.PosType),
			store.Name,
			order.PaymentMethod,
			paymentSystem,
			order.OrderID,
			order.Status,
			order.Customer.Name,
			order.Customer.PhoneNumber,
			order.DeliveryAddress.Label,
			order.OrderCode,
			order.PickUpCode,
			order.AllergyInfo,
			order.DeliveryService,
			order.OrderID,
			strings.ToTitle(order.DeliveryService),
			order.StoreID,
			strings.ToUpper(store.PosType),
			order.PosOrderID,
			orderTime,
			errorInfo,
			utils.GetSuggestedSolution(errorInfo, order.LogLinks.LogStreamLinkByOrderId),
		)

		for index, product := range order.Products {
			message += fmt.Sprintf("<b>%d. %s, x%v</b>\n", index+1, product.Name, product.Quantity)

			var modifierBody string
			for position, attribute := range product.Attributes {
				modifierBody += fmt.Sprintf("    %d. %s, x%v\n", position+1, attribute.Name, attribute.Quantity)
			}

			message += modifierBody + "\n"
		}
	default:
		message = "invalid service type"
	}

	return message
}

func ConstructSuccessMessage(order models.Order, store coreStoreModels.Store) string {
	var paymentTypeRu string
	if order.PosPaymentInfo.PaymentTypeKind == models.POS_PAYMENT_CASH || order.PaymentMethod == models.PAYMENT_METHOD_CASH {
		paymentTypeRu = "Наличные"
	} else {
		paymentTypeRu = "Карта"
	}

	deliveryService := order.DeliveryService
	switch order.DeliveryService {
	case models.QRMENU.String():
		deliveryService = "Kwaaka Direct"
	case models.KWAAKA_ADMIN.String():
		deliveryService = "Kwaaka Admin"
	}

	var paymentSystem = order.PaymentSystem
	if val, ok := paymentSystemMap[order.PaymentSystem]; ok {
		paymentSystem = val
	}

	var deliveryDispatcher = order.DeliveryDispatcher
	if val, ok := deliveryDispatcherMap[order.DeliveryDispatcher]; ok {
		deliveryDispatcher = val
	}

	var message string

	switch {
	case !order.Preorder.Time.Value.IsZero() && order.SendCourier:
		timeZone := checkStoreTimezone(store.Settings.TimeZone.TZ, store.Settings.TimeZone.UTCOffset)
		preorderTime := order.Preorder.Time.Value.Add(time.Duration(timeZone) * time.Hour)
		hours := preorderTime.Hour()
		minutes := preorderTime.Minute()
		timeStr := fmt.Sprintf("%02d:%02d", hours, minutes)
		orderType := "Предзаказ. Приготовить к " + timeStr
		message = fmt.Sprintf("<b>[✅] Предзаказ успешно отправлен на кассу:</b>\n<b>Ресторан:</b> %s\n<b>Адрес ресторана:</b> %s\n<b>Номер телефона ресторана:</b> %s\n<b>Тип заказа:</b> %s\n<b>Cервис Доставки:</b> %s\n<b>ID заказа:</b> %s\n<b>ID 3pl заказа:</b> %s\n<b>Код курьера:</b> %s\n<b>Тип Оплаты:</b> %s\n<b>Система Оплаты:</b> %s\n<b>Агрегатор:</b> %s\n\n", store.Name, store.Address.City+", "+store.Address.Street, store.StorePhoneNumber, orderType, deliveryDispatcher, order.ID, order.DeliveryOrderID, order.PickUpCode, paymentTypeRu, paymentSystem, deliveryService)

	case !order.Preorder.Time.Value.IsZero() && !order.SendCourier:
		timeZone := checkStoreTimezone(store.Settings.TimeZone.TZ, store.Settings.TimeZone.UTCOffset)
		preorderTime := order.Preorder.Time.Value.Add(time.Duration(timeZone) * time.Hour)
		hours := preorderTime.Hour()
		minutes := preorderTime.Minute()
		timeStr := fmt.Sprintf("%02d:%02d", hours, minutes)
		orderType := "Предзаказ & Самовывоз. Приготовить к " + timeStr
		message = fmt.Sprintf("<b>[✅] Ресторан:</b> %s\n<b>Адрес ресторана:</b> %s\n<b>Номер телефона ресторана:</b> %s\n<b>Тип заказа:</b> %s\n<b>ID заказа:</b> %s\n<b>Тип Оплаты:</b> %s\n<b>Система Оплаты:</b> %s\n<b>Агрегатор:</b> %s\n\n", store.Name, store.Address.City+", "+store.Address.Street, store.StorePhoneNumber, orderType, order.ID, paymentTypeRu, paymentSystem, deliveryService)

	case order.Preorder.Time.Value.IsZero() && !order.SendCourier:
		orderType := "Самовывоз"
		message = fmt.Sprintf("<b>[✅] Ресторан:</b> %s\n<b>Адрес ресторана:</b> %s\n<b>Номер телефона ресторана:</b> %s\n<b>Тип заказа:</b> %s\n<b>ID заказа:</b> %s\n<b>Тип Оплаты:</b> %s\n<b>Система Оплаты:</b> %s\n<b>Агрегатор:</b> %s\n\n", store.Name, store.Address.City+", "+store.Address.Street, store.StorePhoneNumber, orderType, order.ID, paymentTypeRu, paymentSystem, deliveryService)

	case order.Preorder.Time.Value.IsZero() && order.SendCourier:
		message = fmt.Sprintf("<b>[✅] Ресторан:</b> %s\n<b>Адрес ресторана:</b> %s\n<b>Номер телефона ресторана:</b> %s\n<b>Cервис Доставки:</b> %s\n<b>ID заказа:</b> %s\n<b>ID 3pl заказа:</b> %s\n<b>Код курьера:</b> %s\n<b>Тип Оплаты:</b> %s\n<b>Система Оплаты:</b> %s\n<b>Агрегатор:</b> %s\n\n", store.Name, store.Address.City+", "+store.Address.Street, store.StorePhoneNumber, deliveryDispatcher, order.ID, order.DeliveryOrderID, order.PickUpCode, paymentTypeRu, paymentSystem, deliveryService)
	}

	customerInfo := fmt.Sprintf("<b>Данные о клиенте:</b>\n<b>Имя:</b> %s\n<b>Номер:</b> %s\n<b>Адрес:</b> %s\n\n", order.Customer.Name, order.Customer.PhoneNumber, order.DeliveryAddress.Label)

	var orderBody = "<b>Продукты:</b>\n"

	for index, product := range order.Products {
		orderBody += fmt.Sprintf("<b>%d. %s, x%v</b>\n", index+1, product.Name, product.Quantity)

		var modifierBody string
		for position, attribute := range product.Attributes {
			modifierBody += fmt.Sprintf("    %d. %s, x%v\n", position+1, attribute.Name, attribute.Quantity)
		}

		orderBody += modifierBody + "\n"
	}

	return message + customerInfo + orderBody
}

func ConstructCancelMessageToNotify(order models.Order, storeName string, comment string) string {
	message := fmt.Sprintf(
		"<b>[❌] Ресторан:</b> %s\n"+
			"<b>Сервис доставки:</b> %s\n"+
			"<b>Имя клиента:</b> %s\n"+
			"<b>Номер клиента для связи:</b> %s\n"+
			"<b>Причина отмены:</b> %s\n\n",
		storeName, order.DeliveryService, order.Customer.Name, order.Customer.PhoneNumber, comment)

	var orderBody = "<b>Продукты:</b>\n"

	for index, product := range order.Products {
		orderBody += fmt.Sprintf("<b>%d. %s, x%v</b>\n", index+1, product.Name, product.Quantity)

		var modifierBody string
		for position, attribute := range product.Attributes {
			modifierBody += fmt.Sprintf("    %d. %s, x%v\n", position+1, attribute.Name, attribute.Quantity)
		}

		orderBody += modifierBody + "\n"
	}

	return message + orderBody
}

func ConstructStoreClosedToNotify(store coreStoreModels.Store, msg, deliveryService string) string {
	location, err := time.LoadLocation(store.Settings.TimeZone.TZ)
	if err != nil {
		log.Printf("location parsing error occured so it will use asia/almaty timezone to send message to telegram: %s", err)
		location, _ = time.LoadLocation("Asia/Almaty")
	}
	closedTime := time.Now().In(location).Format(TimeFormat)

	switch msg {
	case StoreIsOpened:
		return fmt.Sprintf("Ресторан закрыт и был открыт по автооткрытию\nТочка: %s\nАгрегатор: %s\nГород: %s\nВремя закрытия: %s\n", store.Name, deliveryService, store.City, closedTime)
	default:
		return fmt.Sprintf("Ресторан закрыт\nТочка: %s\nАгрегатор: %s\nГород: %s\nВремя закрытия: %s\n", store.Name, deliveryService, store.City, closedTime)
	}
}

func ConstructStoreStatusReportToNotify(store coreStoreModels.Store, durations []coreStoreModels.OpenTimeDuration) string {
	if len(durations) == 0 {
		return ""
	}

	var message string
	var uptime int

	sum := 0
	for _, duration := range durations {
		service := duration.DeliveryService
		actualTime := duration.ActualOpenTimeDuration
		totalTime := duration.TotalOpenTimeDuration

		if totalTime != 0 {
			uptime = int(actualTime.Minutes() / totalTime.Minutes() * 100)
		} else {
			uptime = 0
		}

		sum += uptime

		message += fmt.Sprintf("-%s = %d%%\n", service, uptime)
	}

	message = fmt.Sprintf("%s Uptime:\n%s = %d%% (avg)\n%s", store.Name, store.Name, sum/len(durations), message)

	return message
}

func ConstructOrderStatusChangeMessage(msg *tgbotapi.MessageConfig, status string) {
	var text string

	switch status {
	case models.COOKING_STARTED.String():
		text = "Ваш заказ начал готовится. Среднее время готовки - 20 минут"
	case models.COOKING_COMPLETE.String():
		text = "Ваш заказ готов!"
	case models.CLOSED.String():
		text = "Ваш заказ доставлен. Пожалуйста, оцените его"
		msg.ReplyMarkup = reviewKeyboard
	}

	msg.Text = text
}

func ConstructNoCourierMessage(order models.Order, store coreStoreModels.Store) string {

	log.Printf("construct no courier message for order id: %s", order.ID)

	var paymentTypeRu string

	if order.PosPaymentInfo.PaymentTypeKind == models.PAYMENT_METHOD_CASH {
		paymentTypeRu = "Наличные"
	} else {
		paymentTypeRu = "Карта"
	}

	deliveryService := order.DeliveryService
	switch deliveryService {
	case models.QRMENU.String():
		deliveryService = "Kwaaka Direct"
	case models.KWAAKA_ADMIN.String():
		deliveryService = "Kwaaka Admin"
	}

	var paymentSystem = order.PaymentSystem
	if val, ok := paymentSystemMap[order.PaymentSystem]; ok {
		paymentSystem = val
	}

	var deliveryDispatcher = order.DeliveryDispatcher
	if val, ok := deliveryDispatcherMap[order.DeliveryDispatcher]; ok {
		deliveryDispatcher = val
	}

	message := fmt.Sprintf(
		"<b>[❌] Ресторан:</b> %s\n"+
			"<b>Адрес ресторана:</b> %s\n"+
			"<b>Номер телефона ресторана:</b> %s\n"+
			"<b>Cервис Доставки:</b> %s\n"+
			"<b>ID заказа:</b> %s\n"+
			"<b>ID 3pl заказа:</b> %s\n"+
			"<b>Код курьера:</b> %s\n"+
			"<b>Тип Оплаты:</b> %s\n"+
			"<b>Система Оплаты:</b> %s\n"+
			"<b>Агрегатор:</b> %s\n\n",
		store.Name, store.Address.City+", "+store.Address.Street, store.StorePhoneNumber, deliveryDispatcher,
		order.OrderID, order.DeliveryOrderID, order.PickUpCode, paymentTypeRu, paymentSystem, deliveryService)

	customerInfo := fmt.Sprintf("<b>Данные о клиенте:</b>\n<b>Имя:</b> %s\n<b>Номер:</b> %s\n<b>Адрес:</b> %s\n\n", order.Customer.Name, order.Customer.PhoneNumber, order.DeliveryAddress.Label)

	noCourier := "<b>Заказ готов, но курьер не найден. Необходимо связаться с аккаунт менеджером сервиса доставки</b>\n"

	return message + customerInfo + noCourier
}

func ConstructRefundMessage(order models.Order, store coreStoreModels.Store, amount, reason string) string {
	log.Printf("construct refund message")

	var paymentMethod string
	if order.PosPaymentInfo.PaymentTypeKind == models.PAYMENT_METHOD_CASH {
		paymentMethod = "Наличные"
	} else {
		paymentMethod = "Карта"
	}

	var deliveryService string
	switch order.DeliveryService {
	case models.QRMENU.String():
		deliveryService = "Kwaaka Direct"
	case models.KWAAKA_ADMIN.String():
		deliveryService = "Kwaaka Admin"
	default:
		return ""
	}

	var paymentSystem = order.PaymentSystem
	if val, ok := paymentSystemMap[order.PaymentSystem]; ok {
		paymentSystem = val
	}

	var deliveryDispatcher = order.DeliveryDispatcher
	if val, ok := deliveryDispatcherMap[order.DeliveryDispatcher]; ok {
		deliveryDispatcher = val
	}

	message := fmt.Sprintf(
		"<b>[🔄] Возврат клиенту успешно совершен</b> \n"+
			"<b>Ресторан:</b> %s\n"+
			"<b>Номер телефона ресторана:</b> %s\n"+
			"<b>Cервис Доставки:</b> %s\n"+
			"<b>ID заказа:</b> %s\n"+
			"<b>ID 3pl заказа:</b> %s\n"+
			"<b>Код курьера:</b> %s\n"+
			"<b>Тип Оплаты:</b> %s\n"+
			"<b>Система Оплаты:</b> %s\n"+
			"<b>Сумма Возврата:</b> %s\n"+
			"<b>Причина:</b> %s\n"+
			"<b>Агрегатор:</b> %s\n\n\n",
		store.Name, store.StorePhoneNumber, deliveryDispatcher,
		order.ID, order.DeliveryOrderID, order.PickUpCode, paymentMethod, paymentSystem, amount, reason, deliveryService)

	return message
}

func ConstructCompensationMessage(order models.Order, store coreStoreModels.Store, compensationID, compensationNum, compensationText string) string {
	log.Printf("construct compensation message")

	var paymentMethod string
	if order.PosPaymentInfo.PaymentTypeKind == models.PAYMENT_METHOD_CASH {
		paymentMethod = "Наличные"
	} else {
		paymentMethod = "Карта"
	}

	var deliveryService string
	switch order.DeliveryService {
	case models.QRMENU.String():
		deliveryService = "Kwaaka Direct"
	case models.KWAAKA_ADMIN.String():
		deliveryService = "Kwaaka Admin"
	default:
		return ""
	}

	var paymentSystem = order.PaymentSystem
	if val, ok := paymentSystemMap[order.PaymentSystem]; ok {
		paymentSystem = val
	}

	var deliveryDispatcher = order.DeliveryDispatcher
	if val, ok := deliveryDispatcherMap[order.DeliveryDispatcher]; ok {
		deliveryDispatcher = val
	}

	compensationMsg := fmt.Sprintf(
		"<b>[✅] Заявка на компенсацию</b> \n"+
			"<b>ID Заявки:</b> %s\n"+
			"<b>Номер заявки у ресторана:</b> %s\n"+
			"<b>Текст компенсации:</b>\n"+
			"%s\n\n", compensationID, compensationNum, compensationText)

	commonMsg := fmt.Sprintf(
		"<b>Ресторан:</b> %s\n"+
			"<b>Адрес ресторана:</b> %s\n"+
			"<b>Номер телефона ресторана:</b> %s\n"+
			"<b>Cервис Доставки:</b> %s\n"+
			"<b>ID заказа:</b> %s\n"+
			"<b>ID 3pl заказа:</b> %s\n"+
			"<b>Код заказа:</b> %s\n"+
			"<b>Тип Оплаты:</b> %s\n"+
			"<b>Система Оплаты:</b> %s\n"+
			"<b>Агрегатор:</b> %s\n\n",
		store.Name, store.Address.City+", "+store.Address.Street, store.StorePhoneNumber, deliveryDispatcher,
		order.ID, order.DeliveryOrderID, order.PickUpCode, paymentMethod, paymentSystem, deliveryService)

	customerInfo := fmt.Sprintf(
		"<b>Данные о клиенте:</b>\n"+
			"<b>Имя:</b> %s\n"+
			"<b>Номер:</b> %s\n"+
			"<b>Адрес:</b> %s\n\n",
		order.Customer.Name, order.Customer.PhoneNumber, order.DeliveryAddress.Label)

	var products = "<b>Продукты:</b>\n"
	for index, product := range order.Products {
		products += fmt.Sprintf("<b>%d. %s, x%v</b>\n", index+1, product.Name, product.Quantity)

		var modifierBody string
		for position, attribute := range product.Attributes {
			modifierBody += fmt.Sprintf("    %d. %s, x%v\n", position+1, attribute.Name, attribute.Quantity)
		}
		products += modifierBody + "\n"
	}

	return compensationMsg + commonMsg + customerInfo + products
}

func ConstructCancelDeliveryFromDispatcherPage(order models.Order, store coreStoreModels.Store) string {

	log.Printf("construct cancel delivery from dispatcher page for order id: %s", order.ID)

	deliveryService := order.DeliveryService
	switch deliveryService {
	case models.QRMENU.String():
		deliveryService = "Kwaaka Direct"
	case models.KWAAKA_ADMIN.String():
		deliveryService = "Kwaaka Admin"
	}

	var deliveryDispatcher = order.DeliveryDispatcher
	if val, ok := deliveryDispatcherMap[order.DeliveryDispatcher]; ok {
		deliveryDispatcher = val
	}

	message := fmt.Sprintf(
		"<b>[❌🚚] Ресторан:</b> %s\n"+
			"<b>Адрес ресторана:</b> %s\n"+
			"<b>Номер телефона ресторана:</b> %s\n"+
			"<b>Cервис Доставки:</b> %s\n"+
			"<b>ID заказа:</b> %s\n"+
			"<b>ID 3pl заказа:</b> %s\n"+
			"<b>Код курьера:</b> %s\n"+
			"<b>Агрегатор:</b> %s\n\n",
		store.Name, store.Address.City+" "+store.Address.Street, store.StorePhoneNumber, deliveryDispatcher,
		order.OrderID, order.DeliveryOrderID, order.PickUpCode, deliveryService)

	customerInfo := fmt.Sprintf("<b>Данные о клиенте:</b>\n<b>Имя:</b> %s\n<b>Номер:</b> %s\n<b>Адрес:</b> %s\n\n", order.Customer.Name, order.Customer.PhoneNumber, order.DeliveryAddress.Label)

	cancelDeliveryFromDispatcherPageMsg := "<b>Необходимо отменить заказ в Личном Кабинете провайдера</b>\n"

	return message + customerInfo + cancelDeliveryFromDispatcherPageMsg
}

func ConstructNoDeliveryDispatcherMessage(order models.Order, store coreStoreModels.Store) string {
	message := fmt.Sprintf(
		"<b>[❌] Ресторан:</b> %s\n"+
			"<b>Адрес ресторана:</b> %s\n"+
			"<b>Номер телефона ресторана:</b> %s\n"+
			"<b>ID заказа:</b> %s\n"+
			"<b>Cервис:</b> %s\n\n",
		store.Name, store.Address.City+", "+store.Address.Street, store.StorePhoneNumber,
		order.OrderID, order.DeliveryService)

	customerInfo := fmt.Sprintf("<b>Данные о клиенте:</b>\n<b>Имя:</b> %s\n<b>Номер:</b> %s\n<b>Адрес:</b> %s\n\n", order.Customer.Name, order.Customer.PhoneNumber, order.DeliveryAddress.Label)

	noDispatcher := "<b>Заказ готов, но доставка не создана. Необходимо связаться с аккаунт менеджером сервиса доставки</b>\n"

	return message + customerInfo + noDispatcher
}

func ConstructPutProductToStopListWithErrSolutionMessage(store coreStoreModels.Store, product models3.Product, err string) string {
	if len(product.ExtID) == 0 {
		return ""
	}

	message := fmt.Sprintf(
		"<b>[❌] Restaurant: </b> %s\n"+
			"<b>Product/Attribute id: </b> %s\n"+
			"<b>Product/Attribute name: </b> %s\n"+
			"<b>Error business name: </b> %s\n\n",

		store.Name, product.ExtID, product.Name[0].Value, err)

	errSolution := "<b>Поставили продукт/атрибут на стоп в доступных агрегаторах</b>\n"

	return message + errSolution
}

func checkStoreTimezone(tz string, offset float64) float64 {
	switch tz {
	case "Asia/Almaty":
		return 5
	default:
		return offset
	}
}
