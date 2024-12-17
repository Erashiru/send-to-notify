package pos

import (
	"context"
	"errors"
	"fmt"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	"github.com/kwaaka-team/orders-core/core/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
)

type BasePosManager interface {
	CreateOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, error) // TODO: Dont duplicate this method in each POS manager
	GetOrderStatus(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, error)
	CancelOrder(ctx context.Context, order models.Order, cancelReason, paymentStrategy string, store coreStoreModels.Store) error
	UpdateOrderProblem(ctx context.Context, organizationID, posOrderID string) error
}

func NewPosManager(posName string, globalConfig config.Configuration, ds drivers.DataStore, store coreStoreModels.Store, menu coreMenuModels.Menu, promo coreMenuModels.Promo, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu, storeCli storeClient.Client) (BasePosManager, error) {

	switch posName {
	case models.BurgerKing.String():
		manager, err := NewBKManager(globalConfig, ds.BKOfferRepository())
		if err != nil {
			return nil, err
		}
		return manager, nil
	case models.IIKO.String(), models.Syrve.String():
		manager, err := NewIIKOManager(globalConfig, store, menu, promo, menuClient, aggregatorMenu)
		if err != nil {
			return nil, err
		}
		return manager, nil
	case models.JOWI.String():
		manager, err := NewJOWIManager(globalConfig)
		if err != nil {
			return nil, err
		}
		return manager, nil
	case models.RKeeper.String():
		manager, err := NewRKeeperManager(globalConfig, menu, aggregatorMenu, store)
		if err != nil {
			return nil, err
		}
		return manager, nil
	case models.Paloma.String():
		manager, err := NewPalomaManager(globalConfig, menu, aggregatorMenu, store)
		if err != nil {
			return nil, err
		}
		return manager, nil
	case models.Poster.String():
		manager, err := NewPosterManager(globalConfig, menu, aggregatorMenu, store)
		if err != nil {
			return nil, err
		}
		return manager, nil
	case models.FoodBand.String():
		manager, err := NewFoodBandManager(globalConfig, store, menu, promo, aggregatorMenu)
		if err != nil {
			return nil, err
		}
		return manager, nil
	case models.Yaros.String():
		manager, err := NewYarosManager(globalConfig, menu, aggregatorMenu, store)
		if err != nil {
			return nil, err
		}
		return manager, nil
	case models.RKeeper7XML.String():
		manager, err := NewRKeeper7XMLManager(globalConfig, store, storeCli)
		if err != nil {
			return nil, err
		}
		return manager, nil
	}

	return nil, errors.New("invalid pos name")
}

func ConstructOrderComments(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, string, string) {

	var (
		addressName        = "Адрес"
		orderCodeName      = "Код заказа"
		paymentCashName    = "Наличный"
		paymentDelayedName = "Безналичный"
		commentName        = "Комментарий"
		deliveryName       = "Доставка"
		cutleryName        = "Столовые приборы"
		allergyName        = "Аллергия"
		courierPhoneName   = "Номер курьера"
		paymentTypeName    = "Тип оплаты"
		pickUpToName       = "Приготовить к"
		quantityPerson     = "Количество персон"
	)

	commentSettings := store.Settings.CommentSetting

	if commentSettings.HasCommentSetting {
		addressName = commentSettings.AddressName
		orderCodeName = commentSettings.OrderCodeName
		paymentCashName = commentSettings.CashPaymentName
		paymentDelayedName = commentSettings.DelayedPaymentName
		commentName = commentSettings.CommentName
		deliveryName = commentSettings.DeliveryName
		cutleryName = commentSettings.CutleryName
		allergyName = commentSettings.Allergy
		courierPhoneName = commentSettings.CourierPhoneName
		paymentTypeName = commentSettings.PaymentTypeName
		pickUpToName = commentSettings.PickUpToName
		quantityPerson = commentSettings.QuantityPerson
	}

	addressLabel := order.DeliveryAddress.Label
	if addressLabel == "" {
		addressLabel = models.NO
		if commentSettings.HasCommentSetting {
			addressLabel = commentSettings.CommentDynamicName.HasNotAddress
		}
	}
	customerPhone := order.Customer.PhoneNumber
	if !strings.Contains(customerPhone, "+") {
		customerPhone = "+77771111111"
	}

	var isFromVS string
	if order.IsParentOrder {
		isFromVS += " VS"
	}

	orderCodes := []string{order.OrderCode + isFromVS}
	if order.PickUpCode != "" {
		orderCodes = append(orderCodes, order.PickUpCode)
	}

	anotherCheque := map[string]struct{}{
		"633d71ff4c9853e1aa90f7ee": struct{}{},
		"633d71ff4c9853e1aa90f8d5": struct{}{},
		"6349135cdfaa6e9e24cb4d00": struct{}{},
		"634917094b40f9033b215be8": struct{}{},
		"634d35b63940af93beb93f34": struct{}{},
		"634664622c7064d66300bf5c": struct{}{},
		"634d3495d46d90727f35d5cb": struct{}{},
		"634661c7102046f2cfdff9de": struct{}{},
		"6347d409258d4c093c0ed997": struct{}{},
		"635139e99d23554aeb2d6e59": struct{}{},
		"63513b739d23554aeb2d6e6e": struct{}{},
		"63440249451afd98f952be5e": struct{}{},
		"63513e5521aacb8df8655b7a": struct{}{},
		"6347cfe67df6425f5dad3f16": struct{}{},
		"634fd968b2c8cded05ebca09": struct{}{},
		"635139749d23554aeb2d6e4d": struct{}{},
		"63511085039f207caa317324": struct{}{},
		"63513b6c9d23554aeb2d6e6b": struct{}{},
		"634e529f117c451a38e238d4": struct{}{},
		"63513ac79d23554aeb2d6e64": struct{}{},
		"63513aa99d23554aeb2d6e61": struct{}{},
		"635135bd19159c982fe20b29": struct{}{},
		"635256b1569642609eb51cab": struct{}{},
		"635139929d23554aeb2d6e50": struct{}{},
		"634e7c43798aac535f9a0d21": struct{}{},
		"634e58c806ebe8b223f0f6f9": struct{}{},
		"636a453c847d2ccb5366f5b4": struct{}{},
		"634911f422a68f3fca43d2b5": struct{}{},
		"633d71ff4c9853e1aa90f9d5": struct{}{},
		"633d71ff4c9853e1aa90f6d5": struct{}{},
		"6344072d5d3b3efb5379118e": struct{}{},
		"634fd730b2c8cded05ebca04": struct{}{},
		"635138e99d23554aeb2d6e45": struct{}{},
		"634fd51fb2c8cded05ebc9fe": struct{}{},
		"6351381b9d23554aeb2d6e3f": struct{}{},
		"635135a519159c982fe20b26": struct{}{},
		"634915ba49e96c7b3daf8814": struct{}{},
		"634fc02d8fc6b3e13660857d": struct{}{},
		"633d71ff4c9853e1aa90f555": struct{}{},
		"63513d9921aacb8df8655b73": struct{}{},
		"633d71ff4c9853e1aa90f6e5": struct{}{},
		"634e54d3cf15c6d40ed3ff0d": struct{}{},
		"634d31b505dd136a8b1217bb": struct{}{},
		"634d2e968f26022d35c6b864": struct{}{},
		"635273904ca3858e923816a9": struct{}{},
		"6347d1e07df6425f5dad3f1b": struct{}{},
		"635139f49d23554aeb2d6e5c": struct{}{},
		"6356ca8e3e515feaf9d4d99c": struct{}{},
		"635282138027a4002844eb4f": struct{}{},
		"634e7d33798aac535f9a0d26": struct{}{},
		"63525672569642609eb51ca8": struct{}{},
		"634e59cb06ebe8b223f0f6fe": struct{}{},
		"634409d22b784bed6890a44b": struct{}{},
		"6351096f498696d04553d07e": struct{}{},
		"63491452d130a5bab8cdeb84": struct{}{},
		"63513d4e9d23554aeb2d6e83": struct{}{},
		"634d2fd55160c5fda85fc30e": struct{}{},
		"63971aeffac1865089eb54d5": struct{}{},
		"639ac3ef113f588823bdd141": struct{}{},
		"63d21bee4ea81a11af61fca7": struct{}{},
		"64a3e387e9c9b5a4e9b9d9f6": struct{}{},
		"64e83f1aff0e40b395df66a7": struct{}{},
		"64e70cfbb4b89e0925416997": struct{}{},
		"6405e705d0893bac7c19d94f": struct{}{},
		"654b09e787d35e55b5252f3f": struct{}{},
		"65335f6245bc16b02dbff32d": struct{}{},
		"653367c037f99653d6ba044c": struct{}{},
		"653bb1cf6224c6791c4a906e": struct{}{},
		"6576c666304eb43eac3d8260": struct{}{},
		"66a3a709cbe0ec36e7e7fdc1": struct{}{},
		"66cc75a11035ee7bbbb640cc": struct{}{},
		// Aroma
		"642020886372c553a507d208": struct{}{},
	}

	var orderCodeComment string
	if order.RestaurantID == "642020886372c553a507d208" {
		orderCodeComment += "\n"
	}

	_, ok := anotherCheque[order.RestaurantID]

	var (
		addressComment          string
		paymentType             string
		paymentTypeComment      string
		allergyComment          string
		deliveryComment         string
		quantityPersonComment   string
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

		if store.BillParameter.BillParameters.AddDelivery {
			deliveryComment = fmt.Sprintf("%s: %s", deliveryName, cases.Title(language.Und, cases.NoLower).String(order.DeliveryService))
		}

	} else {
		orderCodeComment += fmt.Sprintf("%s: %s", orderCodeName, strings.Join(orderCodes, ","))
		addressComment = fmt.Sprintf("%s: %s", addressName, addressLabel)
		paymentTypeComment = fmt.Sprintf("%s: %s", paymentTypeName, paymentType)
		customerReqsComment = fmt.Sprintf("%s: %s", commentName, specialRequirements)
		deliveryComment = fmt.Sprintf("%s: %s", deliveryName, cases.Title(language.Und, cases.NoLower).String(order.DeliveryService))
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
	switch ok {
	case true:
		orderCommentArr = append(orderCommentArr, paymentType, allergyComment)
	default:
		if order.Type == "PREORDER" {
			orderCommentArr = append(orderCommentArr, "ПРЕДЗАКАЗ", allergyComment, pickupTimeComment, paymentTypeComment, customerReqsComment, deliveryComment, addressComment)
		} else {
			orderCommentArr = append(orderCommentArr, allergyComment, paymentTypeComment, pickupTimeComment, customerReqsComment, deliveryComment, addressComment)
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
		if store.BillParameter.BillParameters.AddQuantityPersons {
			orderComment = fmt.Sprintf("%s: %v\n", quantityPerson, order.Persons)
		}
	}
	orderComment += strings.Join(orderCommentArr, "\n")
	if order.IsChildOrder && order.VirtualStoreComment != "" {
		orderComment += "\n" + order.VirtualStoreComment
	}

	customerCommentArr := []string{
		paymentTypeComment,
		courierPhoneComment,
		orderCodeComment,
		allergyComment,
		deliveryComment,
		quantityPersonComment,
		customerReqsComment,
		cutleryRequestedComment,
	}
	customerComment := strings.Join(customerCommentArr, "\n")

	// Kitchen comment
	kitchenComment += fmt.Sprintf("Код заказа: %s\nДоставка: %s\nАллергия: %s", order.OrderCode, order.DeliveryService, order.AllergyInfo)
	if order.Type == "PREORDER" {
		kitchenComment = fmt.Sprintf("ПРЕДЗАКАЗ\n%s\n%s", pickupTimeComment, kitchenComment)
	}

	return orderComment, customerComment, kitchenComment
}

func ActiveMenuPositions(ctx context.Context, aggregatorMenu coreMenuModels.Menu) (map[string]string, map[string]string) {
	productsMap := make(map[string]string, len(aggregatorMenu.Products))

	for _, product := range aggregatorMenu.Products {
		if product.PosID == "" {
			productsMap[product.ExtID] = product.ExtID
			continue
		}

		productsMap[product.ExtID] = product.PosID
	}

	attributesMap := make(map[string]string, len(aggregatorMenu.Attributes))

	for _, attribute := range aggregatorMenu.Attributes {
		if attribute.PosID == "" || attribute.PosID != attribute.ExtID {
			attributesMap[attribute.ExtID] = attribute.ExtID
			continue
		}

		attributesMap[attribute.ExtID] = attribute.PosID
	}

	return productsMap, attributesMap
}
