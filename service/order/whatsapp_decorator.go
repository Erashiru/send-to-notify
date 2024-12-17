package order

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/whatsapp/clients"
	"github.com/kwaaka-team/orders-core/service/pos"
	storeServicePkg "github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
	"net/url"
	"strings"
	"time"
)

type WhatsappServiceDecorator struct {
	service        CreationService
	whatsappClient clients.Whatsapp
	storeService   storeServicePkg.Service
	targetPosTypes map[string]struct{}
}

func NewWhatsappServiceDecorator(
	service CreationService,
	whatsappClient clients.Whatsapp,
	storeService storeServicePkg.Service,
	targetPosTypes ...models.Pos,
) (*WhatsappServiceDecorator, error) {
	if service == nil {
		return nil, errors.New("order service is nil")
	}

	if whatsappClient == nil {
		return nil, errors.New("whatsapp client is nil")
	}

	if storeService == nil {
		return nil, errors.New("store service is nil")
	}

	targetPosTypesMap := make(map[string]struct{}, 0)
	for i := range targetPosTypes {
		posType := targetPosTypes[i]
		targetPosTypesMap[posType.String()] = struct{}{}
	}

	return &WhatsappServiceDecorator{
		service:        service,
		whatsappClient: whatsappClient,
		storeService:   storeService,
		targetPosTypes: targetPosTypesMap,
	}, nil
}

func (s *WhatsappServiceDecorator) CreateOrder(ctx context.Context, externalStoreID, deliveryService string, aggReq interface{}, storeSecret string) (models.Order, error) {
	order, err := s.service.CreateOrder(ctx, externalStoreID, deliveryService, aggReq, storeSecret)

	log.Info().Msgf("starting to send notification to whatsapp: order id: %s", order.OrderID)

	if wppErr := s.sendMessage(ctx, order, err); wppErr != nil {
		log.Err(wppErr).Msg(fmt.Sprintf("send whatsapp message error %s", wppErr.Error()))
	}

	return order, err
}

// TODO: uncomment functionality when wpp will be fully ready
func (s *WhatsappServiceDecorator) sendMessage(ctx context.Context, order models.Order, orderErr error) error {
	store, err := s.storeService.GetByID(ctx, order.RestaurantID)
	if err != nil {
		return err
	}

	if !s.isToSendNotification(store) {
		return nil
	}

	var msg string
	switch {
	case orderErr == nil:
		if store.SendWhatsappNotification {
			msg = s.constructSuccessMessage(order, store)
		} else if store.ID == "66cc50a1c64b42b9aefdd89d" {
			msg = s.constructSuccessMsgForNavat(order, store)
		}
	default:
		if store.SendWhatsappNotification {
			msg = s.constructFailMessage(order, store, orderErr.Error())
		}
	}

	msg = url.QueryEscape(msg)

	if msg == "" {
		return nil
	}

	if wppErr := s.whatsappClient.SendMessage(ctx, store.WhatsappChatId, msg); wppErr != nil {
		return wppErr
	}

	return nil
}

func (s *WhatsappServiceDecorator) isToSendNotification(store storeModels.Store) bool {
	return true
}

func (s *WhatsappServiceDecorator) constructSuccessMessage(order models.Order, store storeModels.Store) string {

	log.Info().Msgf("starting constructing whatsapp message: order id: %s, whatsapp chat id: %s", order.OrderID, store.WhatsappChatId)

	var orderInfo, clientInfo, productInfo string

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
	case models.WOLT.String():
		deliveryService = "Wolt"
	case models.GLOVO.String():
		deliveryService = "Glovo"
	case models.YANDEX.String():
		deliveryService = "Yandex"
	}

	var orderType string

	switch {
	case !order.Preorder.Time.Value.IsZero() && order.SendCourier:
		timeZone := s.checkStoreTimezone(store.Settings.TimeZone.TZ, store.Settings.TimeZone.UTCOffset)
		preorderTime := order.Preorder.Time.Value.Add(time.Duration(timeZone) * time.Hour)
		hours := preorderTime.Hour()
		minutes := preorderTime.Minute()
		timeStr := fmt.Sprintf("%02d:%02d", hours, minutes)
		orderType = "Предзаказ. Приготовить к " + timeStr

	case !order.Preorder.Time.Value.IsZero() && !order.SendCourier:
		timeZone := s.checkStoreTimezone(store.Settings.TimeZone.TZ, store.Settings.TimeZone.UTCOffset)
		preorderTime := order.Preorder.Time.Value.Add(time.Duration(timeZone) * time.Hour)
		hours := preorderTime.Hour()
		minutes := preorderTime.Minute()
		timeStr := fmt.Sprintf("%02d:%02d", hours, minutes)
		orderType = "Предзаказ & Самовывоз. Приготовить к " + timeStr
	case order.Preorder.Time.Value.IsZero() && !order.SendCourier:
		orderType = "Самовывоз"
	case order.Preorder.Time.Value.IsZero() && order.SendCourier:
		orderType = "Доставка"
	}

	orderInfo = fmt.Sprintf(
		"*[✅] Ресторан:* %s\n"+
			"*Адрес ресторана:* %s\n"+
			"*Тип заказа:* %s\n"+
			"*Тип оплаты:* %s\n"+
			"*ID заказа:* %s\n"+
			"*Код курьера:* %s\n"+
			"*Код заказа:* %s\n"+
			"*Агрегатор:* %s\n\n",
		order.RestaurantName, store.Address.City+", "+store.Address.Street, orderType, paymentTypeRu, order.ID, order.DeliveryOrderID, order.OrderCode, deliveryService)

	if order.AllergyInfo != "" {
		clientInfo = fmt.Sprintf("*Данные о клиенте:*\n*Имя:* %s\n*Номер:* %s\n*Адрес:* %s\n*Комментарий к заказу:* %s\n\n", order.Customer.Name, order.Customer.PhoneNumber, order.DeliveryAddress.Label, order.AllergyInfo)
	} else {
		clientInfo = fmt.Sprintf("*Данные о клиенте:*\n*Имя:* %s\n*Номер:* %s\n*Адрес:* %s\n\n", order.Customer.Name, order.Customer.PhoneNumber, order.DeliveryAddress.Label)
	}

	productInfo += "*Состав заказа:*\n"

	for index, product := range order.Products {
		productInfo += fmt.Sprintf("%d. %s, x%v\n", index+1, product.Name, product.Quantity)

		var modifierBody string
		for position, attribute := range product.Attributes {
			modifierBody += fmt.Sprintf("    %d. %s, x%v\n", position+1, attribute.Name, attribute.Quantity)
		}
		if modifierBody != "" {
			productInfo += modifierBody + "\n"
		}
	}

	productInfo += "\n\n"

	orderSum := fmt.Sprintf("*Сумма заказа:* %.0f\n", math.Ceil(order.EstimatedTotalPrice.Value))

	log.Info().Msgf("finished constructing whatsapp message: %s", orderInfo+clientInfo+productInfo+orderSum+orderSum)

	return orderInfo + clientInfo + productInfo + orderSum
}

func (s *WhatsappServiceDecorator) checkStoreTimezone(tz string, offset float64) float64 {
	switch tz {
	case "Asia/Almaty":
		return 5
	default:
		return offset
	}
}

var suggestedSolutions = map[string]string{
	"Creation timeout expired":                                    "При создании заказа нет ответа на запрос (скорее всего проблема с интернетом на точке)",
	"is excluded from menu for order's table":                     "продукт убрали из меню, необходимо проверить в айко меню, есть ли данный продукт под другим айди, если нет, то удалить из меню, если есть, заменить айди продукта",
	"PRODUCT NOT FOUND IN POS MENU":                               "по данному айди нет продукта в айко меню, необходимо проверить в айко меню, есть ли данный продукт под другим айди, если нет, то удалить из меню, если есть, заменить айди продукта",
	"Cannot find fixed group modifiers item":                      "отправляемый нами атрибут не существует у продукта в айко, необходимо его удалить и добавить валидный атрибут",
	"has invalid group amount":                                    "мы не отправляем атрибут в нужном количестве, необходимо добавить нужный атрибут в продукт",
	"doesn't belong to your api login included organization list": "по данному апи-логину нет организации, по которой отправляем заказ, необходимо проверить валидность параметров key, organization_id, terminal_id в iiko_cloud",
	"is inactive. Only active products can be added to order":     "продукт неактивный. необходимо проверить, не удален ли продукт или не стоит ли на стопе",
	"Не найден элемент в коллекции":                               "Продукт не найден в коллекции в пос меню. Нужно проверить имеется ли этот продукт под другим айди",
}

var reasons = map[string]string{
	"Creation timeout expired":                                    "При создании заказа нет ответа на запрос",
	"is excluded from menu for order's table":                     "Продукт убрали из меню",
	"PRODUCT NOT FOUND IN POS MENU":                               "По данному айди нет продукта в айко меню",
	"Cannot find fixed group modifiers item":                      "Отправляемый нами атрибут не существует у продукта в айко",
	"has invalid group amount":                                    "мы не отправляем атрибут в нужном количестве",
	"doesn't belong to your api login included organization list": "по данному апи-логину нет организации, по которой отправляем заказ",
	"is inactive. Only active products can be added to order":     "Продукт неактивный либо удален",
	"Не найден элемент в коллекции":                               "Продукт не найден в коллекции в пос меню",
}

func (s *WhatsappServiceDecorator) constructFailMessage(order models.Order, store storeModels.Store, errStr string) string {

	var paymentTypeRu string
	if order.PosPaymentInfo.PaymentTypeKind == models.POS_PAYMENT_CASH || order.PaymentMethod == models.PAYMENT_METHOD_CASH {
		paymentTypeRu = "Наличные"
	} else {
		paymentTypeRu = "Карта"
	}

	var deliverySrv string

	switch order.DeliveryService {
	case models.YANDEX.String():
		deliverySrv = "Yandex"
	case models.GLOVO.String():
		deliverySrv = "Glovo"
	case models.WOLT.String():
		deliverySrv = "Wolt"
	case models.QRMENU.String():
		deliverySrv = "QR Menu"
	case models.KWAAKA_ADMIN.String():
		deliverySrv = "Kwaaka Admin"
	}

	timeZone := s.checkStoreTimezone(store.Settings.TimeZone.TZ, store.Settings.TimeZone.UTCOffset)

	orderInfo := fmt.Sprintf(
		"*Заказ %s № %s не удался* \n"+
			"*Ресторан:* %s\n"+
			"*Тип оплаты:* %s\n"+
			"*Покупатель:* %s С\n"+
			"*Номер телефона:* %s\n"+
			"*Код заказа:* %s\n"+
			"*Код выдачи:* %s\n"+
			"*Комментарий к заказу:* %s\n"+
			"*Дата создания:* %s\n\n",
		deliverySrv,
		order.OrderCode,
		order.RestaurantName,
		paymentTypeRu,
		order.Customer.Name,
		order.Customer.PhoneNumber,
		order.OrderCode,
		order.PickUpCode,
		order.AllergyInfo,
		order.CreatedAt.Add(time.Duration(timeZone)*time.Hour).Format("15:04:05 02-01-2006"),
	)

	productInfo := "*Состав заказа:*\n"

	for index, product := range order.Products {
		productInfo += fmt.Sprintf("%d. %s, x%v\n", index+1, product.Name, product.Quantity)

		var modifierBody string
		for position, attribute := range product.Attributes {
			modifierBody += fmt.Sprintf("    %d. %s, x%v\n", position+1, attribute.Name, attribute.Quantity)
		}
		if modifierBody != "" {
			productInfo += modifierBody + "\n"
		}
	}

	productInfo += "\n"
	productInfo += fmt.Sprintf("*Сумма заказа:* %.0f\n", math.Ceil(order.EstimatedTotalPrice.Value))
	productInfo += "\n\n"

	var errInfo string

	switch order.DeliveryService {
	case models.GLOVO.String():
		errInfo += "*Что делать:*\n" +
			"1) Пробить данный заказ вручную на кассу\n" +
			"2) Проверить и прожимать статусы заказа в планшете Glovo\n\n"
	case models.WOLT.String():
		errInfo += "*Что делать:*\n" +
			"1) Пробить данный заказ вручную на кассу\n" +
			"2) Проверить и прожимать статусы заказа в планшете Wolt\n\n"
	case models.YANDEX.String():
		errInfo += "*Что делать:*\n" +
			"1) Пробить данный заказ вручную на кассу\n" +
			"2) Проверить и прожимать статусы заказа в личном кабинете Яндекс.Еда\n\n"
	}

	if errStr != "" {
		errInfo += fmt.Sprintf(
			"*Ошибка:* \n" + errStr + "\n\n")
	}

	var reason string
	for key := range reasons {
		if strings.Contains(errStr, key) {
			if reasonInMap, ok := reasons[key]; ok {
				reason += fmt.Sprintf("*Причина:*\n" + reasonInMap + "\n\n")
			}
		}
	}

	if reason == "" {
		switch order.FailReason.Code {
		case pos.CREATION_TIMEOUT_CODE:
			reason = fmt.Sprintf(
				"*Причина:*\n" +
					"Неправильный/Нестабильная работа интернет, " +
					"обмена апи моноблока/кассы/фронта с сервером, " +
					" отсутствие электроснабжения или интернета на моноблоке/кассе/фронте\n\n")
		case pos.NEED_TO_HEAL_ORDER_CODE:
			reason = fmt.Sprintf(
				"*Причина:*\n " +
					"Ошибка в составе заказа\n\n")
		case pos.PAYMENT_TYPE_MISSED_CODE:
			reason = fmt.Sprintf(
				"*Причина:*\n" +
					"Тип оплаты отсутствует\n\n")
		case pos.INTEGRATION_OFF_CODE:
			reason += fmt.Sprintf(
				"*Причина:*\n" +
					"Интеграция отключена либо закончилась лицензия\n\n")
		}
	}

	errInfo += reason

	var solution string
	for key := range suggestedSolutions {
		if strings.Contains(errStr, key) {
			if solutionInMap, ok := suggestedSolutions[key]; ok {
				solution = fmt.Sprintf("*Решение:*\n" + solutionInMap + "\n")
			}
		}
	}

	if solution == "" {
		solution = fmt.Sprintf("*Решение:*\n" + "Обратиться в службу поддержки\n")
	}

	errInfo += solution

	return orderInfo + productInfo + errInfo
}

func (s *WhatsappServiceDecorator) constructSuccessMsgForNavat(order models.Order, store storeModels.Store) string {

	var paymentTypeRu string
	if order.PosPaymentInfo.PaymentTypeKind == models.POS_PAYMENT_CASH || order.PaymentMethod == models.PAYMENT_METHOD_CASH {
		paymentTypeRu = "Наличные"
	} else {
		paymentTypeRu = "Карта"
	}

	var orderInfo, clientInfo, productInfo, pickUpCode string
	deliveryService := order.DeliveryService
	switch order.DeliveryService {
	case models.WOLT.String():
		deliveryService = "Wolt"
	case models.GLOVO.String():
		deliveryService = "Glovo"
	case models.YANDEX.String():
		deliveryService = "Yandex"
	}

	orderInfo = fmt.Sprintf(
		"*[✅] Ресторан:* %s\n"+
			"*Тип оплаты:* %s\n"+
			"*Адрес ресторана:* %s\n"+
			"*Код заказа:* %s\n"+
			"*Агрегатор:* %s\n\n",
		order.RestaurantName, paymentTypeRu, store.Address.City+", "+store.Address.Street, order.OrderCode, deliveryService)

	if order.DeliveryService == models.GLOVO.String() {
		pickUpCode = fmt.Sprintf("*GLOVO PICK UP CODE:* %s\n", order.PickUpCode)
	}

	if order.AllergyInfo != "" {
		clientInfo = fmt.Sprintf("*Данные о клиенте:*\n*Имя:* %s\n*Комментарий к заказу:* %s\n\n", order.Customer.Name, order.AllergyInfo)
	} else {
		clientInfo = fmt.Sprintf("*Данные о клиенте:*\n*Имя:* %s\n\n", order.Customer.Name)
	}

	for _, product := range order.Products {
		productInfo += fmt.Sprintf("- %s, x%v\n", product.Name, product.Quantity)

		var modifierBody string
		for _, attribute := range product.Attributes {
			modifierBody += fmt.Sprintf("    - %s, x%v\n", attribute.Name, attribute.Quantity)
		}
		if modifierBody != "" {
			productInfo += modifierBody + "\n"
		}
	}
	orderSum := fmt.Sprintf("\n\n*Сумма заказа:* %.0f\n", math.Ceil(order.EstimatedTotalPrice.Value))

	return orderInfo + pickUpCode + clientInfo + "*Состав заказа:*\n" + productInfo + orderSum
}
