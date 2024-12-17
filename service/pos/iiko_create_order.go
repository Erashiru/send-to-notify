package pos

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	menuUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

func isManga(storeId string) bool {
	mangaMap := map[string]struct{}{
		"6336ad1e1ccf421f062ee6a3": {},
		"6353d604531150d3081a88a8": {},
		"6353d9a5531150d3081a88bb": {},
		"6353d7b7531150d3081a88af": {},
		"6353d8f5531150d3081a88b4": {},
		"66067c9c11ade97820e2fb66": {},
	}

	if _, ok := mangaMap[storeId]; !ok {
		return false
	}

	return true
}

func isFarsh(storeId string) bool {
	farshMap := map[string]struct{}{
		"6274e50ae64402c9194b69ad": {},
		"627f7e5ccd4cf5567dc3d74e": {},
		"627f7d3f0066bf575db139c9": {},
		"626e3a88f5c32de682df3a9f": {},
		"63f32cff7e94dda798e6aa7a": {},
		"62711d3f4132a9fbdd256d4f": {},
	}

	if _, ok := farshMap[storeId]; !ok {
		return false
	}

	return true
}

func isRumi(storeId string) bool {
	rumiMap := map[string]struct{}{
		"63e33f6a3dbec7f34af0baaa": {},
		"63e360f45ae0673b2030378d": {},
		"62d6a528e67cb9c9f27ae50a": {},
		"6422d8e46aff02fcd948ae5a": {},
		"6422d8309d19b00fe6f9f41e": {},
		"6422d92b6aff02fcd948ae5c": {},
		"6423cbf22d57a44f7f937bad": {},
	}

	if _, ok := rumiMap[storeId]; !ok {
		return false
	}

	return true
}

func isBahandi(storeId string) bool {
	bahandiMap := map[string]struct{}{
		"65bb7b3bf32bd66dc387592b": {},
		"634917094b40f9033b215be8": {},
		"6344072d5d3b3efb5379118e": {},
		"634e529f117c451a38e238d4": {},
		"634664622c7064d66300bf5c": {},
		"6351096f498696d04553d07e": {},
		"634fd968b2c8cded05ebca09": {},
		"633d71ff4c9853e1aa90f6e5": {},
		"635139749d23554aeb2d6e4d": {},
		"63491452d130a5bab8cdeb84": {},
		"63513ac79d23554aeb2d6e64": {},
		"63513e5521aacb8df8655b7a": {},
		"634661c7102046f2cfdff9de": {},
		"634911f422a68f3fca43d2b5": {},
		"6351381b9d23554aeb2d6e3f": {},
		"634d31b505dd136a8b1217bb": {},
		"634e58c806ebe8b223f0f6f9": {},
		"6347d1e07df6425f5dad3f1b": {},
		"6349135cdfaa6e9e24cb4d00": {},
		"63513d9921aacb8df8655b73": {},
		"633d71ff4c9853e1aa90f6d5": {},
		"635139f49d23554aeb2d6e5c": {},
		"6356ca8e3e515feaf9d4d99c": {},
		"63d21bee4ea81a11af61fca7": {},
		"634fc02d8fc6b3e13660857d": {},
		"633d71ff4c9853e1aa90f9d5": {},
		"635256b1569642609eb51cab": {},
		"63513b6c9d23554aeb2d6e6b": {},
		"63513aa99d23554aeb2d6e61": {},
		"635139e99d23554aeb2d6e59": {},
		"634915ba49e96c7b3daf8814": {},
		"634fd730b2c8cded05ebca04": {},
		"634409d22b784bed6890a44b": {},
		"634d35b63940af93beb93f34": {},
		"635139929d23554aeb2d6e50": {},
		"634e59cb06ebe8b223f0f6fe": {},
		"639ac3ef113f588823bdd141": {},
		"636a453c847d2ccb5366f5b4": {},
		"635138e99d23554aeb2d6e45": {},
		"635135a519159c982fe20b26": {},
		"633d71ff4c9853e1aa90f7ee": {},
		"634e54d3cf15c6d40ed3ff0d": {},
		"635273904ca3858e923816a9": {},
		"633d71ff4c9853e1aa90f8d5": {},
		"634e7d33798aac535f9a0d26": {},
		"63525672569642609eb51ca8": {},
		"634e7c43798aac535f9a0d21": {},
		"6347d409258d4c093c0ed997": {},
		"63513d4e9d23554aeb2d6e83": {},
		"633d71ff4c9853e1aa90f555": {},
		"634d3495d46d90727f35d5cb": {},
		"634d2e968f26022d35c6b864": {},
		"6347cfe67df6425f5dad3f16": {},
		"63513b739d23554aeb2d6e6e": {},
		"63440249451afd98f952be5e": {},
		"63971aeffac1865089eb54d5": {},
		"634fd51fb2c8cded05ebc9fe": {},
		"635135bd19159c982fe20b29": {},
		"634d2fd55160c5fda85fc30e": {},
		"63511085039f207caa317324": {},
		"635282138027a4002844eb4f": {},
		"6405e705d0893bac7c19d94f": {},
		"6437a3adbfb06d111b0eccbd": {},
		"64a3e387e9c9b5a4e9b9d9f6": {},
		"64e70cfbb4b89e0925416997": {},
		"64e83f1aff0e40b395df66a7": {},
		"65335f6245bc16b02dbff32d": {},
		"653367c037f99653d6ba044c": {},
		"653bb1cf6224c6791c4a906e": {},
		"654b09e787d35e55b5252f3f": {},
		"6576c666304eb43eac3d8260": {},
		"65b390a754f52d75b09a3f70": {},
		"65b531c94ebf8ce908647f6a": {},
		"6613a2b333cf5fd14750c612": {},
	}

	if _, ok := bahandiMap[storeId]; !ok {
		return false
	}

	return true
}

func (iikoSvc *iikoService) toItemsAndCombos(order models.Order) ([]iikoModels.Item, []iikoModels.Combo, int) {
	orderItems := make([]iikoModels.Item, 0, len(order.Products))
	combos := make([]iikoModels.Combo, 0)

	amount := 0

	for _, product := range order.Products {
		if product.IsCombo {
			for _, attribute := range product.Attributes {
				orderItem := iikoModels.Item{
					ProductId: attribute.ID,
					Type:      "Product",
					Amount:    float64(attribute.Quantity),
					ComboInformation: &iikoModels.ComboInformation{
						ComboID:       product.ID,
						ComboSourceID: product.SourceActionID,
						ComboGroupID:  attribute.GroupID,
					},
				}

				orderItems = append(orderItems, orderItem)
			}

			combos = append(combos, iikoModels.Combo{
				ID:        product.ID,
				Name:      product.Name,
				Amount:    product.Quantity,
				Price:     int(product.Price.Value),
				SourceID:  product.SourceActionID,
				ProgramID: product.ProgramID,
			})

			amount += int(product.Price.Value)
			continue
		}

		orderItem := iikoModels.Item{
			ProductId:     product.ID,
			Type:          "Product",
			Amount:        float64(product.Quantity),
			Price:         utils.PointerOfFloat(product.Price.Value),
			ProductSizeID: product.SizeId,
		}

		amount += int(product.Price.Value)

		unique := make(map[string]iikoModels.ItemModifier)

		for _, attribute := range product.Attributes {
			if val, ok := unique[attribute.ID]; !ok {
				unique[attribute.ID] = iikoModels.ItemModifier{
					ProductId:      attribute.ID,
					Amount:         float64(attribute.Quantity),
					ProductGroupId: attribute.GroupID,
					Price:          attribute.Price.Value,
				}
			} else {
				val.Amount += float64(attribute.Quantity)
				unique[attribute.ID] = val
			}

		}

		for _, val := range unique {
			orderItem.Modifiers = append(orderItem.Modifiers, val)
			amount += int(val.Price * val.Amount)
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, combos, amount
}

func (iikoSvc *iikoService) addCommonPaymentType(paymentTypeKind string, payments []iikoModels.Payment, order models.Order, isProcessedExternally bool, priceSource string, amount int) []iikoModels.Payment {

	payment := iikoModels.Payment{
		PaymentTypeKind:       paymentTypeKind,
		Sum:                   int(order.EstimatedTotalPrice.Value) - int(order.PartnerDiscountsProducts.Value),
		PaymentTypeID:         order.PosPaymentInfo.PaymentTypeID,
		IsProcessedExternally: isProcessedExternally,
	}

	if priceSource == models.POSPriceSource {
		payment.Sum = amount
	}

	// TODO: primer plohogo koda
	if order.DeliveryService == models.KWAAKA_ADMIN.String() && len(order.Promos) != 0 {
		payment.Sum -= payment.Sum * order.Promos[0].Discount / 100
	}

	// Todo: понять как отправить на кассу 0 сумму, и возможно ли это. Судя по ответам IIKO пока нет
	if order.DeliveryService == models.QRMENU.String() && order.PaymentMethod == models.PAYMENT_METHOD_CASH {
		payment.IsProcessedExternally = false
	}

	payments = append(payments, payment)

	return payments
}

func (iikoSvc *iikoService) addPromotionPaymentType(payments []iikoModels.Payment, order models.Order, isProcessedExternally bool) []iikoModels.Payment {
	if order.PosPaymentInfo.PromotionPaymentTypeID != "" && order.PartnerDiscountsProducts.Value != 0 {
		payments = append(payments, iikoModels.Payment{
			PaymentTypeKind:       "Card",
			Sum:                   int(order.PartnerDiscountsProducts.Value),
			PaymentTypeID:         order.PosPaymentInfo.PromotionPaymentTypeID,
			IsProcessedExternally: isProcessedExternally,
		})
	}

	return payments
}

func (iikoSvc *iikoService) getPaymentTypeKind(order models.Order) (models.Order, string, error) {
	var paymentTypeKind string

	switch order.PosPaymentInfo.PaymentTypeKind {
	case "Cash":
		paymentTypeKind = "Cash"
	case "Card":
		paymentTypeKind = "Card"
	default:
		order.FailReason.Code = PAYMENT_TYPE_MISSED_CODE
		order.FailReason.Message = PAYMENT_TYPE_MISSED
		return order, "", errors.New("payment type is missed")
	}

	return order, paymentTypeKind, nil
}

func (iikoSvc *iikoService) toPayments(order models.Order, priceSource string, amount int) ([]iikoModels.Payment, error) {

	order, paymentTypeKind, err := iikoSvc.getPaymentTypeKind(order)
	if err != nil {
		return nil, err
	}

	payments := make([]iikoModels.Payment, 0, 2)

	isProcessedExternally := true
	if order.PosPaymentInfo.IsProcessedExternally != nil {
		isProcessedExternally = *order.PosPaymentInfo.IsProcessedExternally
	}

	payments = iikoSvc.addCommonPaymentType(paymentTypeKind, payments, order, isProcessedExternally, priceSource, amount)
	payments = iikoSvc.addPromotionPaymentType(payments, order, isProcessedExternally)

	return payments, nil
}

func belongsToAttributeGroup(reqProduct coreMenuModels.Product,
	reqAttributeId string,
	reqAttributeGroupsMap map[string]coreMenuModels.AttributeGroup) bool {

	for _, groupId := range reqProduct.AttributesGroups {
		group, ok := reqAttributeGroupsMap[groupId]
		if !ok {
			continue
		}

		for _, attributeId := range group.Attributes {
			if attributeId == reqAttributeId {
				return true
			}
		}
	}

	return false
}

func bindModifiersToLastProducts(items []iikoModels.Item, productsMap map[string]coreMenuModels.Product,
	attributeGroupsMap map[string]coreMenuModels.AttributeGroup) []iikoModels.Item {

	lastParentItem := make(map[string]string)

	allModifiers := make([]iikoModels.ItemModifier, 0, 4)

	// get all modifiers
	for _, item := range items {
		allModifiers = append(allModifiers, item.Modifiers...)
	}

	for _, item := range items {
		product, ok := productsMap[item.ProductId]
		if !ok {
			continue
		}

		// check modifier parent last product
		for _, modifier := range allModifiers {
			if belongsToAttributeGroup(product, modifier.ProductId, attributeGroupsMap) {
				lastParentItem[modifier.ProductId] = product.ProductID
			}
		}
	}

	anotherModifiersToProduct := map[string][]iikoModels.ItemModifier{}

	for index, item := range items {
		productModifiers := make([]iikoModels.ItemModifier, 0, len(items))

		var anotherModifiers []iikoModels.ItemModifier

		// check if assigned modifiers to last product
		if anotherModifier, ok := anotherModifiersToProduct[item.ProductId]; ok {
			anotherModifiers = anotherModifier
		}

		// map unique modifiers with increased amount
		unique := make(map[string]iikoModels.ItemModifier)
		for _, anotherModifier := range anotherModifiers {
			if val, ok := unique[anotherModifier.ProductId]; ok {
				val.Amount += anotherModifier.Amount
				unique[anotherModifier.ProductId] = val
			} else {
				unique[anotherModifier.ProductId] = anotherModifier
			}
		}

		for _, currentModifier := range item.Modifiers {
			// if not exist parent
			if val, exist := lastParentItem[currentModifier.ProductId]; !exist {
				productModifiers = append(productModifiers, currentModifier)
			} else {
				// if parent current product
				if val == item.ProductId {
					// if id is equal another modifier increasing amount
					if mod, _ok := unique[currentModifier.ProductId]; _ok {
						currentModifier.Amount += mod.Amount
						// delete modifier from map
						delete(unique, currentModifier.ProductId)
					}

					productModifiers = append(productModifiers, currentModifier)
				} else {
					// adding modifier to another product map
					anotherModifiersToProduct[val] = append(anotherModifiersToProduct[val], currentModifier)
				}
			}
		}

		if len(unique) != 0 {
			for _, val := range unique {
				productModifiers = append(productModifiers, val)
			}
		}

		items[index].Modifiers = productModifiers
	}

	return items
}

func (iikoSvc *iikoService) sendOrder(ctx context.Context, posOrder iikoModels.CreateDeliveryRequest) (iikoModels.CreateDeliveryResponse, error) {

	//var errs custom.Error
	utils.Beautify("IIKO Request Body", posOrder)

	createResponse, err := iikoSvc.iikoClient.CreateDeliveryOrder(ctx, posOrder)
	if err != nil {
		log.Err(err).Msg("iikoService, sendOrder -IIKO error")
		//errs.Append(err, validator.ErrIgnoringPos)
		return iikoModels.CreateDeliveryResponse{}, err
	}

	if createResponse.OrderInfo != nil && createResponse.OrderInfo.ErrorInfo != nil {
		log.Err(err).Msgf("iikoService, sendOrder - IIKO Response ErrorInfo error: %+v", *createResponse.OrderInfo.ErrorInfo)
		//errs.Append(err, errors.New(createResponse.OrderInfo.ErrorInfo.Message))
		return iikoModels.CreateDeliveryResponse{}, err
	}

	return createResponse, nil
}

func (iikoSvc *iikoService) fillRequestData(order models.Order, orderItems []iikoModels.Item, combos []iikoModels.Combo, payments []iikoModels.Payment, orderComment, customerComment string) iikoModels.Order {
	phoneNumber := order.Customer.PhoneNumber

	if order.Customer.PhoneNumberWithPlus != "" {
		phoneNumber = order.Customer.PhoneNumberWithPlus
	}

	iikoOrder := iikoModels.Order{
		Phone:   phoneNumber,
		Combos:  combos,
		Comment: orderComment,
		Customer: &iikoModels.Customer{
			Name:    order.Customer.Name,
			Comment: customerComment,
			Gender:  "NotSpecified",
		},
		Payments:   payments,
		Items:      orderItems,
		OperatorID: order.OperatorID,
	}

	//TODO: REFACTOR ME PLS
	if order.DeliveryService == "kwaaka_admin" && len(order.Promos) != 0 {

		for _, promo := range order.Promos {
			if promo.IikoDiscountId != "" {
				discountInfo := iikoModels.DiscountsInfo{
					Discounts: []iikoModels.Discount{
						{
							Type:           "RMS",
							DiscountTypeId: promo.IikoDiscountId,
						},
					},
				}

				iikoOrder.DiscountsInfo = &discountInfo
			}
		}

	}

	return iikoOrder
}

func (iikoSvc *iikoService) addKitchenComments(req models.Order, iikoOrder iikoModels.Order, kitchenComment string, store coreStoreModels.Store) (models.Order, iikoModels.Order) {
	if !store.IikoCloud.SendKitchenComments {
		return req, iikoOrder
	}

	if req.AllergyInfo == "" {
		req.AllergyInfo = models.NO
	}

	if len(iikoOrder.Items) != 0 {
		iikoOrder.Items[len(iikoOrder.Items)-1].Comment = kitchenComment
	}

	return req, iikoOrder
}

func (iikoSvc *iikoService) setOrderType(iikoOrder iikoModels.Order, orderType string) iikoModels.Order {
	iikoOrder.OrderTypeID = orderType
	return iikoOrder
}

func (iikoSvc *iikoService) setOrderServiceType(req models.Order, iikoOrder iikoModels.Order, store coreStoreModels.Store) iikoModels.Order {
	orderServiceType := iikoModels.SERVICE_TYPE_DELIVERY_BY_COURIER

	if req.IsPickedUpByCustomer || req.IsMarketplace && req.DeliveryService != models.YANDEX.String() {
		orderServiceType = iikoModels.SERVICE_TYPE_DELIVERY_BY_CLIENT
	}

	if req.SendCourier && req.DeliveryService == coreStoreModels.QRMENU.String() && req.DeliveryDispatcher == coreStoreModels.SELFDELIVERY.String() {
		orderServiceType = iikoModels.SERVICE_TYPE_DELIVERY_BY_COURIER
	}

	iikoOrder.OrderServiceType = orderServiceType

	return iikoOrder
}

func (iikoSvc *iikoService) getCompleteBefore(store coreStoreModels.Store, req models.Order) string {
	completeBeforeDate := req.EstimatedPickupTime.Value.Time

	if completeBeforeDate.IsZero() {
		completeBeforeDate = time.Now().UTC().Add(time.Hour)
	}

	// Convert UTC time to stores local time
	completeBeforeDate = completeBeforeDate.Add(time.Duration(store.Settings.TimeZone.UTCOffset) * time.Hour)

	completeBefore := completeBeforeDate.Format("2006-01-02 15:04:05.000")

	return completeBefore
}

func (iikoSvc *iikoService) setDeliveryAddress(store coreStoreModels.Store, req models.Order, iikoOrder iikoModels.Order) iikoModels.Order {
	if iikoOrder.OrderServiceType != "" && iikoOrder.OrderServiceType != iikoModels.SERVICE_TYPE_DELIVERY_BY_COURIER {
		return iikoOrder
	}

	if iikoOrder.OrderTypeID != "" && req.PosPaymentInfo.OrderTypeService != iikoModels.SERVICE_TYPE_DELIVERY_BY_COURIER {
		return iikoOrder
	}

	switch req.DeliveryService {
	case models.YANDEX.String():
		if !req.IsMarketplace {
			return iikoOrder
		}

		iikoOrder.DeliveryPoint = iikoSvc.addressByFields(req.DeliveryAddress)

	default:
		var streetName string

		switch req.DeliveryService {
		case "glovo":
			streetName = "Доставка Глово"
		case "wolt":
			streetName = "Wolt доставка"
		case "express24":
			streetName = "Express24 доставка"
		case "chocofood":
			streetName = "Chocofood доставка"
		case "qr_menu":
			streetName = "Qr доставка"
		default:
			for _, deliveryService := range store.ExternalConfig {
				if deliveryService.Type == req.DeliveryService {
					streetName = fmt.Sprintf("%v доставка", deliveryService.Type)
				}
			}
		}

		if req.DeliveryAddress.Street != "" {
			streetName += fmt.Sprintf(" %s", req.DeliveryAddress.Street)
		}

		iikoOrder.DeliveryPoint = &iikoModels.DeliveryPoint{
			Coordinates: &iikoModels.Coordinates{
				Latitude:  req.DeliveryAddress.Latitude,
				Longitude: req.DeliveryAddress.Longitude,
			},
			Address: &iikoModels.Address{
				Street: &iikoModels.Street{
					Name: streetName,
					City: store.Address.City,
				},
				House: "1",
			},
			Comment: req.DeliveryAddress.Label,
		}
	}

	return iikoOrder
}

func (iikoSvc *iikoService) setCreationResult(req models.Order, response iikoModels.CreateDeliveryResponse) models.Order {
	req.CreationResult = models.CreationResult{
		CorrelationId: response.CorrelationID,
		OrderInfo: models.OrderInfo{
			ID:             response.OrderInfo.ID,
			OrganizationID: response.OrderInfo.OrganizationID,
			Timestamp:      int64(response.OrderInfo.Timestamp),
			CreationStatus: response.OrderInfo.CreationStatus,
		},
	}
	return req
}

func (iikoSvc *iikoService) setError(req models.Order, response iikoModels.CreateDeliveryResponse) models.Order {
	req.CreationResult.ErrorDescription = response.OrderInfo.ErrorInfo.Description + response.OrderInfo.ErrorInfo.Message
	return req

}

func (iikoSvc *iikoService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration,
	store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	var err error

	order, err = prepareAnOrder(ctx, order, store, menu, aggregatorMenu, menuClient)
	if err != nil {
		return order, err
	}

	utils.Beautify("prepared order", order)

	posOrder, order, err := iikoSvc.constructPosOrder(ctx, order, store)
	if err != nil {
		log.Trace().Err(validator.ErrCastingPos).Msg("")
		return order, err
	}

	posOrder.Order.ExternalNumber = order.OrderCode

	if isManga(order.RestaurantID) {
		posOrder.Order.Items = bindModifiersToLastProducts(posOrder.Order.Items, menuUtils.ProductsMap(menu), menuUtils.AtributeGroupsMap(menu))
	}

	utils.Beautify("constructed order", posOrder)

	order, err = iikoSvc.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	createResponse, err := iikoSvc.sendOrder(ctx, posOrder)
	if err != nil {
		if errors.Is(err, validator.ErrCastingPos) {
			order.FailReason.Code = NEED_TO_HEAL_ORDER_CODE
			order.FailReason.Message = err.Error()
		}
		return order, err
	}

	utils.Beautify("create order response", createResponse)

	order = setPosOrderId(order, createResponse.OrderInfo.ID)

	order = iikoSvc.setCreationResult(order, createResponse)

	if createResponse.OrderInfo.ErrorInfo != nil {
		order = iikoSvc.setError(order, createResponse)
	}

	utils.Beautify("result order body", order)

	return order, nil

}

func (iikoSvc *iikoService) setCompleteBefore(ctx context.Context, req models.Order, iikoOrder iikoModels.Order, store coreStoreModels.Store) iikoModels.Order {
	completeBefore := iikoSvc.getCompleteBefore(store, req)

	if req.Type == models.ORDER_TYPE_PREORDER || iikoOrder.OrderServiceType == iikoModels.SERVICE_TYPE_DELIVERY_BY_COURIER || (iikoOrder.OrderTypeID != "" && req.PosPaymentInfo.OrderTypeService == iikoModels.SERVICE_TYPE_DELIVERY_BY_COURIER) {
		iikoOrder.CompleteBefore = completeBefore
	}

	return iikoOrder
}

func (iikoSvc *iikoService) constructPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (iikoModels.CreateDeliveryRequest, models.Order, error) {
	var err error

	orderItems, combos, amount := iikoSvc.toItemsAndCombos(order)

	// TODO: refactor
	var delivery coreStoreModels.StoreDelivery
	if store.Delivery != nil && len(store.Delivery) > 0 && !order.IsMarketplace {
		if len(store.Delivery) == 1 && store.Delivery[0].IsActive {
			delivery = store.Delivery[0]
		} else if len(store.Delivery) > 1 {
			for _, storeDelivery := range store.Delivery {
				if storeDelivery.IsActive && storeDelivery.Price == int(order.DeliveryFee.Value) {
					delivery = storeDelivery
					break
				}
			}
		}

		if delivery.ID != "" {
			orderItems = append(orderItems, iikoModels.Item{
				ProductId: delivery.ID,
				Price:     utils.PointerOfFloat(float64(delivery.Price)),
				Type:      "Product",
				Amount:    1,
			})
			order.EstimatedTotalPrice.Value = order.EstimatedTotalPrice.Value + float64(delivery.Price)
			amount += delivery.Price
		}
	}

	payments, err := iikoSvc.toPayments(order, store.Settings.PriceSource, amount)
	if err != nil {
		return iikoModels.CreateDeliveryRequest{}, order, err
	}

	orderComment, customerComment, kitchenComment := iikoSvc.constructOrderComments(order, store)

	iikoOrder := iikoSvc.fillRequestData(order, orderItems, combos, payments, orderComment, customerComment)

	order, iikoOrder = iikoSvc.addKitchenComments(order, iikoOrder, kitchenComment, store)

	if order.PosPaymentInfo.OrderType != "" {
		iikoOrder = iikoSvc.setOrderType(iikoOrder, order.PosPaymentInfo.OrderType)
	} else {
		iikoOrder = iikoSvc.setOrderServiceType(order, iikoOrder, store)
	}

	iikoOrder = iikoSvc.setDeliveryAddress(store, order, iikoOrder)

	iikoOrder = iikoSvc.setCompleteBefore(ctx, order, iikoOrder, store)

	return iikoModels.CreateDeliveryRequest{
		OrganizationID:  store.IikoCloud.OrganizationID,
		TerminalGroupID: store.IikoCloud.TerminalID,
		Order:           &iikoOrder,
		CreateOrderSettings: &iikoModels.CreateOrderSettings{
			TransportToFrontTimeout: iikoSvc.transportToFrontTimeout,
		},
	}, order, nil
}

func (iikoSvc *iikoService) addressByFields(address models.DeliveryAddress) *iikoModels.DeliveryPoint {
	r := &iikoModels.DeliveryPoint{
		Coordinates: &iikoModels.Coordinates{
			Latitude:  address.Latitude,
			Longitude: address.Longitude,
		},
		Address: &iikoModels.Address{
			Street: &iikoModels.Street{
				Name: address.Street,
				City: address.City,
			},
			House:     address.HouseNumber,
			Entrance:  address.Entrance,
			Doorphone: address.Intercom,
			Floor:     address.Floor,
			Type:      iikoModels.LegacyAddressType,
		},
		Comment: address.Label,
	}

	if r.Address.Street.City == "" {
		r.Address.Street.City = address.Label
	}
	if r.Address.House == "" {
		r.Address.House = address.Label
	}

	return r
}
