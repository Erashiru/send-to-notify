package pos

import (
	"context"
	"fmt"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	menuUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/kwaaka-team/orders-core/core/models"
)

func (bs *BasePosService) constructOrderComments(order models.Order, store coreStoreModels.Store) (string, string, string) {

	var (
		addressName        = "Адрес"
		orderCodeName      = "Код заказа"
		paymentCashName    = "Наличный"
		paymentDelayedName = "Безналичный"
		commentName        = "Комментарий"
		//deliveryName       = "Доставка"
		cutleryName      = "Столовые приборы"
		allergyName      = "Аллергия"
		courierPhoneName = "Номер курьера"
		paymentTypeName  = "Тип оплаты"
		pickUpToName     = "Приготовить к"
		quantityPerson   = "Количество персон"
	)

	commentSettings := store.Settings.CommentSetting

	if commentSettings.HasCommentSetting {
		addressName = commentSettings.AddressName
		orderCodeName = commentSettings.OrderCodeName
		paymentCashName = commentSettings.CashPaymentName
		paymentDelayedName = commentSettings.DelayedPaymentName
		commentName = commentSettings.CommentName
		//deliveryName = commentSettings.DeliveryName
		cutleryName = commentSettings.CutleryName
		allergyName = commentSettings.Allergy
		courierPhoneName = commentSettings.CourierPhoneName
		paymentTypeName = commentSettings.PaymentTypeName
		pickUpToName = commentSettings.PickUpToName
		quantityPerson = commentSettings.QuantityPerson
	}

	addressLabel := order.DeliveryAddress.Label
	if order.DeliveryAddress.Flat != "" {
		addressLabel = order.DeliveryAddress.Label + ", кв." + order.DeliveryAddress.Flat
	}

	if addressLabel == "" {
		if commentSettings.HasCommentSetting {
			addressLabel = commentSettings.CommentDynamicName.HasNotAddress
		}
	}
	customerPhone := order.Customer.PhoneNumber
	if !strings.Contains(customerPhone, "+") {
		customerPhone = "+77771111111"
	}

	orderCodes := make([]string, 0)

	if order.IsParentOrder {
		orderCodes = append(orderCodes, "VS")
	}

	if bs.isAnotherBill(order.RestaurantID, order.DeliveryService) {
		orderCodes = append(orderCodes, order.OrderCode)
	}

	if order.PickUpCode != "" {
		orderCodes = append(orderCodes, order.OrderCodePrefix+order.PickUpCode)
	}

	var orderCodeComment string
	if order.RestaurantID == "642020886372c553a507d208" {
		orderCodeComment += "\n"
	}

	var (
		addressComment     string
		paymentType        string
		paymentTypeComment string
		allergyComment     string
		//deliveryComment         string
		customerReqsComment     string
		cutleryRequestedComment string
	)

	deliveryServiceName := order.DeliveryService
	if deliveryServiceName == models.QRMENU.String() {
		deliveryServiceName = "Kwaaka Direct"
	} else {
		deliveryServiceName = cases.Title(language.Und, cases.NoLower).String(order.DeliveryService)
	}

	switch order.PaymentMethod {
	case models.PAYMENT_METHOD_CASH:
		paymentType = fmt.Sprintf("%s %s", paymentCashName, deliveryServiceName)
	case models.PAYMENT_METHOD_DELAYED:
		paymentType = fmt.Sprintf("%s %s", paymentDelayedName, deliveryServiceName)
	}

	var orderCommentArr = make([]string, 0, 4)

	pickupTimeComment := fmt.Sprintf("%s: %s", pickUpToName, order.EstimatedPickupTime.Value.Time.
		Add(time.Duration(store.Settings.TimeZone.UTCOffset)*time.Hour).
		Format("15:04:05"))

	specialRequirements := order.SpecialRequirements
	if specialRequirements == "" {
		specialRequirements = models.NO
		if commentSettings.HasCommentSetting {
			specialRequirements = commentSettings.CommentDynamicName.HasNotSpecialRequirements
		}
	}
	// Setteling order comment poles by using store core model bill parameteres
	if store.BillParameter.IsActive {
		if store.BillParameter.BillParameters.AddOrderCode {
			orderCodeComment += fmt.Sprintf("%s: %s", orderCodeName, strings.Join(orderCodes, ","))
		}

		if store.BillParameter.BillParameters.AddAddress {
			addressComment = fmt.Sprintf("%s: %s", addressName, addressLabel)
		}

		if store.BillParameter.BillParameters.AddPaymentType {
			paymentTypeComment = fmt.Sprintf("%s: %s", paymentTypeName, paymentType)
		}

		if store.BillParameter.BillParameters.AddComments {
			customerReqsComment = fmt.Sprintf("%s: %s", commentName, specialRequirements)
		}
	} else {
		orderCodeComment += fmt.Sprintf("%s: %s", orderCodeName, strings.Join(orderCodes, ","))
		addressComment = fmt.Sprintf("%s: %s", addressName, addressLabel)
		paymentTypeComment = fmt.Sprintf("%s: %s", paymentTypeName, paymentType)
		customerReqsComment = fmt.Sprintf("%s: %s", commentName, specialRequirements)
		//deliveryComment = fmt.Sprintf("%s: %s", deliveryName, cases.Title(language.Und, cases.NoLower).String(order.DeliveryService))
	}

	var cutleryRequested string
	switch order.CutleryRequested {
	case true:
		cutleryRequested = "Нужны"
		if commentSettings.HasCommentSetting {
			cutleryRequested = commentSettings.CommentDynamicName.HasCutlery
		}
	case false:
		cutleryRequested = "Не нужны"
		if commentSettings.HasCommentSetting {
			cutleryRequested = commentSettings.CommentDynamicName.HasNotCutlery
		}
	}
	cutleryRequestedComment = fmt.Sprintf("%s: %s", cutleryName, cutleryRequested)

	allergy := order.AllergyInfo
	if allergy == "" {
		allergy = models.NO

		if commentSettings.HasCommentSetting {
			allergy = commentSettings.CommentDynamicName.HasNotAllergy
		}
	}
	allergyComment = fmt.Sprintf("%s: %s", allergyName, allergy)

	var kitchenComment string

	if store.Settings.RetrySetting.IsActive && order.IsRetry {
		log.Info().Msgf("retry settings is active true for %s, order_id %s, delivery service %s", store.Name, order.OrderID, order.DeliveryService)
		orderCommentArr = append(orderCommentArr, store.Settings.RetrySetting.Message)
		kitchenComment = store.Settings.RetrySetting.Message + "\n"
	} else if order.IsRetry {
		log.Info().Msgf("retry notification without settings for %s, order_id %s, delivery service %s", store.Name, order.OrderID, order.DeliveryService)
		orderCommentArr = append(orderCommentArr, "!!Возможно дублированный заказ!! Проверьте на наличие дубликата!!")
		kitchenComment = store.Settings.RetrySetting.Message + "\n"
	}

	// Construct comment
	orderCommentArr = append(orderCommentArr, orderCodeComment)
	isAnotherBill := bs.isAnotherBill(order.RestaurantID, order.DeliveryService)
	if isAnotherBill {
		orderCommentArr = append(orderCommentArr, paymentType, allergyComment)
	} else {
		if order.Type == "PREORDER" {
			if addressComment != "" {
				orderCommentArr = append(orderCommentArr, "ПРЕДЗАКАЗ", allergyComment, paymentTypeComment, customerReqsComment, addressComment)
			} else {
				orderCommentArr = append(orderCommentArr, "ПРЕДЗАКАЗ", allergyComment, paymentTypeComment, customerReqsComment)
			}
		} else {
			if addressComment != "" {
				orderCommentArr = append(orderCommentArr, allergyComment, addressComment, customerReqsComment, paymentTypeComment)
			} else {
				orderCommentArr = append(orderCommentArr, allergyComment, customerReqsComment, paymentTypeComment)
			}
		}
	}

	courierPhone := order.Courier.PhoneNumber
	if courierPhone == "" {
		courierPhone = models.NO

		if commentSettings.HasCommentSetting {
			courierPhone = commentSettings.CommentDynamicName.HasNotCourierPhone
		}
	}

	courierPhoneComment := fmt.Sprintf("%s: %s", courierPhoneName, courierPhone)
	var orderComment string
	if order.DeliveryService == models.YANDEX.String() {
		orderComment = fmt.Sprintf("%s: %v\n", quantityPerson, order.Persons)
	}
	orderComment += strings.Join(orderCommentArr, "\n")
	if order.IsChildOrder && order.VirtualStoreComment != "" {
		orderComment += "\n" + order.VirtualStoreComment
	}

	customerCommentArr := []string{
		courierPhoneComment,
		allergyComment,
		customerReqsComment,
		cutleryRequestedComment,
	}
	customerComment := strings.Join(customerCommentArr, "\n")

	// Kitchen comment
	if order.AllergyInfo != "" {
		kitchenComment += fmt.Sprintf("Код заказа: %s\nАллергия: %s\nДоставка: %s", order.OrderCodePrefix+order.PickUpCode, order.AllergyInfo, order.DeliveryService)
	} else {
		kitchenComment += fmt.Sprintf("Код заказа: %s\nДоставка: %s", order.OrderCodePrefix+order.PickUpCode, order.DeliveryService)
	}
	if order.Type == "PREORDER" {
		kitchenComment = fmt.Sprintf("ПРЕДЗАКАЗ\n%s\n%s", pickupTimeComment, kitchenComment)
	}

	if len(customerComment) > 500 {
		customerComment = customerComment[:500]
	}

	if len(orderComment) > 500 {
		orderComment = orderComment[:500]
	}

	if len(kitchenComment) > 500 {
		kitchenComment = kitchenComment[:500]
	}

	return orderComment, customerComment, kitchenComment
}

func setCustomerPhoneNumber(ctx context.Context, req models.Order, store coreStoreModels.Store) models.Order {
	if req.Customer.PhoneNumberWithPlus != "" {
		return req
	}

	if !strings.Contains(req.Customer.PhoneNumber, "+") {
		req.Customer.PhoneNumber = "+77771111111"
	}

	if store.Settings.CommentSetting.DefaultCourierPhone != "" && req.DeliveryService == models.YANDEX.String() {
		req.Customer.PhoneNumber = store.Settings.CommentSetting.DefaultCourierPhone
	}

	if isManga(req.RestaurantID) && req.DeliveryService == models.YANDEX.String() {
		req.Customer.PhoneNumber = "+79999999999"
	}

	if isManga(req.RestaurantID) && req.DeliveryService == models.GLOVO.String() {
		req.Customer.PhoneNumber = "+78888888888"
	}

	if isManga(req.RestaurantID) && req.DeliveryService == models.WOLT.String() {
		req.Customer.PhoneNumber = "+77777777777"
	}

	if isFarsh(req.RestaurantID) && strings.Contains(req.Customer.PhoneNumber, " доб. ") {
		i := strings.Index(req.Customer.PhoneNumber, " доб. ")
		req.Customer.PhoneNumber = req.Customer.PhoneNumber[:i]
	}

	return req
}

func prepareAnOrder(ctx context.Context, order models.Order, store coreStoreModels.Store, menu coreMenuModels.Menu, aggregatorMenu coreMenuModels.Menu, menuCli menuCore.Client) (models.Order, error) {

	var (
		serviceFee float64
		err        error
	)

	order = setCustomerPhoneNumber(ctx, order, store)

	if isHaniRestDelivery(order.StoreID, order.DeliveryService) {
		order = haniRestAddDeliveryProduct(order, menu)
	}

	order, promosMap, giftMap, promoWithPercentMap, err := getPromosMap(ctx, order, menuCli)
	if err != nil {
		return order, err
	}

	order, serviceFee, err = fullFillProducts(
		order, store,
		menuUtils.ProductsMap(menu),
		menuUtils.AtributesMap(menu),
		menuUtils.AtributeGroupsMap(menu),
		menuUtils.ComboMap(menu),
		aggregatorMenu,
		promosMap,
		promoWithPercentMap,
	)
	if err != nil {
		return order, err
	}

	if serviceFee != 0 {
		order.HasServiceFee = true
		order.ServiceFeeSum = serviceFee
		order.EstimatedTotalPrice.Value = order.EstimatedTotalPrice.Value - serviceFee
	}

	order = applyOrderDiscount(ctx, order, promosMap, giftMap, promoWithPercentMap)

	return order, nil
}

func getExternalCfgByName(externalConfig []coreStoreModels.StoreExternalConfig, configType string) (coreStoreModels.StoreExternalConfig, error) {
	for _, ext := range externalConfig {
		if ext.Type == configType {
			return ext, nil
		}
	}
	return coreStoreModels.StoreExternalConfig{}, fmt.Errorf("not found %s cfg in restaurant", configType)
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
