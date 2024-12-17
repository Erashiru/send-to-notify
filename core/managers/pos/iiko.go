package pos

import (
	"context"
	"fmt"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"

	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	iikoConf "github.com/kwaaka-team/orders-core/pkg/iiko/clients"
	IIKOClient "github.com/kwaaka-team/orders-core/pkg/iiko/clients/http"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	menuCoreModel "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	MenuUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	"github.com/rs/zerolog/log"
)

type IIKOManager struct {
	iikoClient         iikoConf.IIKO
	menuClient         menuCore.Client
	menu               coreMenuModels.Menu
	globalConfig       config.Configuration
	productsMap        map[string]coreMenuModels.Product
	attributesMap      map[string]coreMenuModels.Attribute
	attributeGroupsMap map[string]coreMenuModels.AttributeGroup
	promoMap           map[string]string
	aggregatorMenu     coreMenuModels.Menu
}

func NewIIKOManager(globalConfig config.Configuration, store coreStoreModels.Store, menu coreMenuModels.Menu, promo coreMenuModels.Promo, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu) (IIKOManager, error) {
	baseUrl := globalConfig.IIKOConfiguration.BaseURL
	if store.PosType == models.Syrve.String() {
		log.Info().Msgf("syrve client, api_key %s", store.IikoCloud.Key)
		baseUrl = globalConfig.SyrveConfiguration.BaseURL
	}
	if store.IikoCloud.CustomDomain != "" {
		baseUrl = store.IikoCloud.CustomDomain
	}

	iikoClient, err := IIKOClient.New(&iikoConf.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		ApiLogin: store.IikoCloud.Key,
	})

	if err != nil {
		log.Trace().Err(err).Msgf("Cant initialize %s Client. iiko configs: %+v", store.PosType, iikoClient)
		return IIKOManager{}, err
	}

	return IIKOManager{
		iikoClient:         iikoClient,
		menu:               menu,
		globalConfig:       globalConfig,
		menuClient:         menuClient,
		productsMap:        MenuUtils.ProductsMap(menu),
		attributesMap:      MenuUtils.AtributesMap(menu),
		attributeGroupsMap: MenuUtils.AtributeGroupsMap(menu),
		promoMap:           map[string]string{},
		aggregatorMenu:     aggregatorMenu,
	}, nil
}

func (manager IIKOManager) IsAliveStatus(ctx context.Context, store coreStoreModels.Store) (bool, error) {
	result, err := manager.iikoClient.IsAlive(ctx, iikoModels.IsAliveRequest{
		OrganizationIds:  []string{store.IikoCloud.OrganizationID},
		TerminalGroupIds: []string{store.IikoCloud.TerminalID},
	})
	if err != nil {
		return false, err
	}

	if len(result.IsAliveStatus) > 0 {
		return result.IsAliveStatus[0].IsAlive, nil
	}

	return false, errors.New("iiko is alive status is empty")
}

func (manager IIKOManager) CancelOrder(ctx context.Context, order models.Order, cancelReason, paymentStrategy string, store coreStoreModels.Store) error {

	removalTypeId := manager.detectCancelStrategy(order, store)

	_, err := manager.iikoClient.CancelDeliveryOrder(ctx, store.IikoCloud.OrganizationID, order.PosOrderID, removalTypeId)
	if err != nil {
		log.Err(err).Msg("cancel order in IIKOManger error")
		return err
	}

	return nil
}

func (manager IIKOManager) GetOrderStatus(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, error) {
	return "", nil
}

func (manager IIKOManager) detectCancelStrategy(order models.Order, store coreStoreModels.Store) string {

	timePassed := time.Now().Sub(order.OrderTime.Value.Time)

	if timePassed.Seconds() > 30 {
		if store.IikoCloud.RemovalTypeIdWithCharge != "" {
			return store.IikoCloud.RemovalTypeIdWithCharge
		}
		return ""
	}

	if store.IikoCloud.RemovalTypeIdWithoutCharge != "" {
		return store.IikoCloud.RemovalTypeIdWithoutCharge
	}

	return ""
}

func (manager IIKOManager) fulfillOrderProducts(ctx context.Context, order models.Order, store coreStoreModels.Store) ([]iikoModels.Item, coreStoreModels.StoreDelivery, float64, error) {
	var (
		items      []iikoModels.Item
		serviceFee float64
	)
	//var err error

	aggregatorProducts, aggregatorAttributes := ActiveMenuPositions(ctx, manager.aggregatorMenu)

	var isMarketplace bool
	switch order.DeliveryService {
	case "glovo":
		isMarketplace = store.Glovo.IsMarketplace
	case "wolt":
		isMarketplace = store.Wolt.IsMarketplace
	case "chocofood":
		isMarketplace = store.Chocofood.IsMarketplace
	case "qr_menu":
		isMarketplace = store.QRMenu.IsMarketplace
	default:
		for _, deliveryService := range store.ExternalConfig {
			if deliveryService.Type == order.DeliveryService {
				isMarketplace = deliveryService.IsMarketplace
			}
		}
	}

	// Add delivery to order if store has own delivery
	var delivery coreStoreModels.StoreDelivery
	if store.Delivery != nil && len(store.Delivery) > 0 && !isMarketplace {
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
			items = append(items, iikoModels.Item{
				ProductId: delivery.ID,
				Price:     utils.PointerOfFloat(float64(delivery.Price)),
				Type:      "Product",
				Amount:    1,
			})
			order.EstimatedTotalPrice.Value = order.EstimatedTotalPrice.Value + float64(delivery.Price)
		}
	}

	// Add products to order
	for _, product := range order.Products {
		if len(product.Promos) != 0 {
			if product.Promos[0].Type == "GIFT" {
				orderItem := iikoModels.Item{
					ProductId: product.ID,
					Price:     utils.PointerOfFloat(0),
					Type:      "Product",
					Amount:    float64(product.Quantity),
				}

				items = append(items, orderItem)
				continue
			}
		}

		var modifiers []iikoModels.ItemModifier
		var modifiersIds = make(map[string]string)
		var modifiersPrice float64

		productPosID, exist := aggregatorProducts[product.ID]
		if exist {
			product.ID = productPosID
		} else {
			log.Warn().Msgf("Product with ID %s %s not matched", product.ID, product.Name)
		}

		// Add product attributes to order product
		for _, attribute := range product.Attributes {

			for _, deliveryService := range store.ExternalConfig {
				if deliveryService.Type == order.DeliveryService {
					modifiersPrice += attribute.Price.Value * float64(attribute.Quantity)
				}
			}

			if attribute.ID == models.ServiceFee {
				serviceFee += attribute.Price.Value * float64(product.Quantity) * float64(attribute.Quantity)
				continue
			}

			attributePosID, exist_ := aggregatorAttributes[attribute.ID]
			if exist_ {
				attribute.ID = attributePosID
			} else {
				log.Warn().Msgf("Attribute with ID %s %s not matched", attribute.ID, attribute.Name)
			}

			menuAttribute, ok := manager.attributesMap[attribute.ID]

			if ok {
				modifiersIds[attribute.ID] = attribute.ID
				attributeGroupID := MenuUtils.FindAttributeGroupID(product.ID, attribute.ID, manager.productsMap, manager.attributeGroupsMap, manager.attributesMap, store.IikoCloud.IsExternalMenu, menuAttribute.ParentAttributeGroup)

				var priceAmount float64
				if store.Settings.PriceSource == models.POSPriceSource {
					priceAmount = menuAttribute.Price
				} else {
					priceAmount = attribute.Price.Value
				}

				itemModifier := iikoModels.ItemModifier{
					ProductId: attribute.ID,
					Amount:    float64(attribute.Quantity),
					Price:     priceAmount,
				}

				if attributeGroupID != "" {
					itemModifier.ProductGroupId = attributeGroupID
				}
				if store.IikoCloud.IsExternalMenu {
					itemModifier.ProductGroupId = menuAttribute.ParentAttributeGroup
				}
				modifiers = append(modifiers, itemModifier)
			} else {
				menuAttributeProduct, ok_ := manager.productsMap[attribute.ID]

				if ok_ {
					var priceAmount float64
					if store.Settings.PriceSource == models.POSPriceSource {
						priceAmount = menuAttributeProduct.Price[0].Value
					} else {
						priceAmount = attribute.Price.Value
					}

					items = append(items, iikoModels.Item{
						ProductId: menuAttributeProduct.ProductID,
						Price:     utils.PointerOfFloat(priceAmount),
						Type:      "Product",
						Amount:    float64(attribute.Quantity * product.Quantity),
					})
				} else {
					return nil, coreStoreModels.StoreDelivery{}, 0, fmt.Errorf("ATTRIBUTE NOT FOUND IN POS MENU, ID %s, NAME %s", attribute.ID, attribute.Name)
				}
			}

		}

		menuProduct, ok := manager.productsMap[product.ID]

		if !ok {
			log.Info().Msgf("PRODUCT NOT FOUND IN POS MENU, ID %s, NAME %s", product.ID, product.Name)

			return nil, coreStoreModels.StoreDelivery{}, 0, errors.Wrap(errs.ErrProductNotFound, fmt.Sprintf("PRODUCT NOT FOUND IN POS MENU, ID %s, NAME %s", product.ID, product.Name))
		}

		// Add default attributes to order
		if menuProduct.MenuDefaultAttributes != nil && len(menuProduct.MenuDefaultAttributes) > 0 {
			for _, defaultAttribute := range menuProduct.MenuDefaultAttributes {

				defaultAttributePosID, exist__ := aggregatorAttributes[defaultAttribute.ExtID]
				if exist__ {
					defaultAttribute.ExtID = defaultAttributePosID
				}

				_, ok_ := modifiersIds[defaultAttribute.ExtID]

				if ok_ {
					continue
				}

				price := 0
				amount := 1
				if defaultAttribute.DefaultAmount > 0 {
					amount = defaultAttribute.DefaultAmount
				}

				menuDefaultAttribute, ok := manager.attributesMap[defaultAttribute.ExtID]
				var attributeGroupID string
				if ok {
					attributeGroupID = MenuUtils.FindAttributeGroupID(product.ID, defaultAttribute.ExtID, manager.productsMap, manager.attributeGroupsMap, manager.attributesMap, store.IikoCloud.IsExternalMenu, menuDefaultAttribute.ParentAttributeGroup)
				}

				defaultModifier := iikoModels.ItemModifier{
					ProductId: defaultAttribute.ExtID,
					Amount:    float64(amount),
					Price:     float64(price),
				}

				if attributeGroupID != "" {
					defaultModifier.ProductGroupId = attributeGroupID
				}

				modifiersIds[defaultAttribute.ExtID] = defaultAttribute.ExtID
				modifiers = append(modifiers, defaultModifier)
			}
		}

		var itemPrice float64
		switch store.Settings.PriceSource {
		case models.POSPriceSource:
			itemPrice = menuProduct.Price[0].Value
		case models.DeliveryServicePriceSource:
			itemPrice = product.Price.Value - modifiersPrice
		default:
			itemPrice = product.Price.Value - modifiersPrice
		}

		for _, promo := range product.Promos {
			switch promo.Type {
			case "FIXED":
				itemPrice = itemPrice - float64(promo.Discount)
			case "PERCENTAGE":
				itemPrice = itemPrice - float64(promo.Discount)
			}
		}

		orderItem := iikoModels.Item{
			ProductId: menuProduct.ProductID,
			Price:     utils.PointerOfFloat(itemPrice),
			Type:      "Product",
			Amount:    float64(product.Quantity),
			Modifiers: modifiers,
		}

		if menuProduct.ProductID == "" {
			orderItem.ProductId = menuProduct.ExtID
		}

		if menuProduct.SizeID != "" {
			orderItem.ProductSizeID = menuProduct.SizeID
		}

		items = append(items, orderItem)
	}
	//if err != nil {
	//	return nil, storeModels.StoreDelivery{}, err
	//}

	return items, delivery, serviceFee, nil
}

func (manager IIKOManager) getPayment(ctx context.Context, order models.Order, store coreStoreModels.Store, delivery coreStoreModels.StoreDelivery, serviceFee float64) (iikoModels.Payment, string, string, error) {

	var paymentTypes *coreStoreModels.DeliveryServicePaymentType
	switch order.DeliveryService {
	case "glovo":
		paymentTypes = &store.Glovo.PaymentTypes
	case "wolt":
		paymentTypes = &store.Wolt.PaymentTypes
	case "chocofood":
		paymentTypes = &store.Chocofood.PaymentTypes
	case "qr_menu":
		paymentTypes = &store.QRMenu.PaymentTypes
	case "express24":
		paymentTypes = &store.Express24.PaymentTypes
	case "kwaaka_admin":
		paymentTypesKwaakaAdmin := coreStoreModels.DeliveryServicePaymentType{
			CASH: coreStoreModels.PaymentType{
				PaymentTypeID:   order.PosPaymentInfo.PaymentTypeID,
				PaymentTypeKind: order.PosPaymentInfo.PaymentTypeKind,
			},
			DELAYED: coreStoreModels.PaymentType{
				PaymentTypeID:   order.PosPaymentInfo.PaymentTypeID,
				PaymentTypeKind: order.PosPaymentInfo.PaymentTypeKind,
			},
		}
		paymentTypes = &paymentTypesKwaakaAdmin
	case "starter_app":
		paymentTypes = &store.StarterApp.PaymentTypes

	default:
		for _, deliveryService := range store.ExternalConfig {
			if deliveryService.Type == order.DeliveryService {
				currentPaymentTypes := deliveryService.PaymentTypes
				paymentTypes = &currentPaymentTypes
			}
		}
	}

	if paymentTypes == nil {
		return iikoModels.Payment{}, "", "", errors.New("no any payment types")
	}

	var paymentType coreStoreModels.PaymentType
	var isProcessedExternally bool
	var orderType string
	var orderTypeService string

	switch order.PaymentMethod {
	case models.PAYMENT_METHOD_DELAYED:
		paymentType = paymentTypes.DELAYED
		orderType = paymentTypes.DELAYED.OrderType
		orderTypeService = paymentTypes.DELAYED.OrderTypeService
		isProcessedExternally = true
		if order.IsChildOrder {
			orderType = paymentTypes.DELAYED.OrderTypeForVirtualStore
		}
	case models.PAYMENT_METHOD_CASH:
		paymentType = paymentTypes.CASH
		orderType = paymentTypes.CASH.OrderType
		orderTypeService = paymentTypes.CASH.OrderTypeService
		isProcessedExternally = false
		if order.IsChildOrder {
			orderType = paymentTypes.CASH.OrderTypeForVirtualStore
		}
	default:
		log.Info().Msgf("PayMent Method: %v", order.PaymentMethod)
	}

	// Choosed custom IsProcessedExternally if exists.
	switch order.PaymentMethod {
	case "CASH":
		if paymentTypes.CASH.IsProcessedExternally != nil {
			isProcessedExternally = *paymentTypes.CASH.IsProcessedExternally
		}
	case "DELAYED":
		if paymentTypes.DELAYED.IsProcessedExternally != nil {
			isProcessedExternally = *paymentTypes.DELAYED.IsProcessedExternally
		}
	}

	return iikoModels.Payment{
		PaymentTypeKind:       paymentType.PaymentTypeKind,
		Sum:                   int(order.EstimatedTotalPrice.Value) + delivery.Price - int(order.PartnerDiscountsProducts.Value) - int(serviceFee),
		PaymentTypeID:         paymentType.PaymentTypeID,
		IsProcessedExternally: isProcessedExternally,
	}, orderType, orderTypeService, nil
}

func (manager IIKOManager) applyOrderDiscount(ctx context.Context, posOrder iikoModels.Order, order models.Order) iikoModels.Order {
	productsIDs := make([]string, 0, len(order.Products))
	for _, product := range order.Products {
		productsIDs = append(productsIDs, product.ID)
	}

	promosMap := make(map[string]menuCoreModel.PromoDiscount)
	giftMap := make(map[string]string)

	promos, err := manager.menuClient.GetStorePromos(ctx, menuCoreModel.GetPromosSelector{
		StoreID:         order.RestaurantID,
		DeliveryService: order.DeliveryService,
		ProductIDs:      productsIDs,
		IsActive:        true,
	})

	if err != nil {
		log.Error().Err(err).Msgf("Failed to get promos")
		return posOrder
	}

	for _, promo := range promos {
		for _, productID := range promo.ProductIds {
			promosMap[productID] = promo
		}

		for _, gift := range promo.ProductGifts {
			giftMap[gift.PromoId] = gift.ProductId
		}
	}

	var discount float64
	for i := 0; i < len(posOrder.Items); i++ {

		if productID, ok := giftMap[posOrder.Items[i].ProductId]; ok {
			posOrder.Items[i].ProductId = productID
			discount += *posOrder.Items[0].Price
			posOrder.Items[i].Price = utils.PointerOfFloat(0)
			continue
		}

		if promo, ok := promosMap[posOrder.Items[i].ProductId]; ok {

			if promo.Type == models.Discount.String() {
				discountAmount := float64(promo.Percent) * *posOrder.Items[i].Price / 100.0
				*posOrder.Items[i].Price -= discountAmount
				discount += discountAmount * posOrder.Items[i].Amount
			}

		}

		for j := 0; j < len(posOrder.Items[i].Modifiers); j++ {
			if promo, ok := promosMap[posOrder.Items[i].Modifiers[j].ProductId]; ok {

				if promo.Type == models.Discount.String() {
					discountAmount := float64(promo.Percent) * posOrder.Items[i].Modifiers[j].Price / 100.0
					posOrder.Items[i].Modifiers[j].Price -= discountAmount
					discount += discountAmount * posOrder.Items[i].Modifiers[j].Amount
				}
			}
		}
	}

	return posOrder
}

func (manager IIKOManager) toPayment(ctx context.Context, order models.Order) (iikoModels.Payment, error) {
	var paymentTypeKind string

	switch order.PaymentMethod {
	case "CASH":
		paymentTypeKind = "Cash"
	case "DELAYED":
		paymentTypeKind = "Card"
	default:
		return iikoModels.Payment{}, errors.New("payment type is missed")
	}

	isProcessedExternally := true
	if order.PosPaymentInfo.IsProcessedExternally != nil {
		isProcessedExternally = *order.PosPaymentInfo.IsProcessedExternally
	}

	return iikoModels.Payment{
		PaymentTypeKind:       paymentTypeKind,
		Sum:                   int(order.EstimatedTotalPrice.Value) - int(order.PartnerDiscountsProducts.Value), // TODO: delivery price?
		PaymentTypeID:         order.PosPaymentInfo.PaymentTypeID,
		IsProcessedExternally: isProcessedExternally,
	}, nil
}

func (manager IIKOManager) toItemsAndCombos(ctx context.Context, order models.Order) ([]iikoModels.Item, []iikoModels.Combo) {
	orderItems := make([]iikoModels.Item, 0, len(order.Products))
	combos := make([]iikoModels.Combo, 0)

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
			continue
		}

		orderItem := iikoModels.Item{
			ProductId: product.ID,
			Type:      "Product",
			Amount:    float64(product.Quantity),
			Price:     utils.PointerOfFloat(product.Price.Value),
		}

		for _, attribute := range product.Attributes {
			orderItem.Modifiers = append(orderItem.Modifiers, iikoModels.ItemModifier{
				ProductId:      attribute.ID,
				Amount:         float64(attribute.Quantity),
				ProductGroupId: attribute.GroupID,
				Price:          attribute.Price.Value,
			})
		}

		orderItems = append(orderItems, orderItem)
	}

	return orderItems, combos
}

func (manager IIKOManager) ConstructPosOrderWithBase(ctx context.Context, order models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	var err error

	orderItems, combos := manager.toItemsAndCombos(ctx, order)

	payment, err := manager.toPayment(ctx, order)
	if err != nil {
		return iikoModels.CreateDeliveryRequest{}, order, err
	}

	orderComment, customerComment, kitchenComment := ConstructOrderComments(ctx, order, store)

	iikoOrder := iikoModels.Order{
		Phone:   order.Customer.PhoneNumber,
		Combos:  combos,
		Comment: orderComment,
		Customer: &iikoModels.Customer{
			Name:    order.Customer.Name,
			Comment: customerComment,
			Gender:  "NotSpecified",
		},
		Payments: []iikoModels.Payment{
			payment,
		},
		Items: orderItems,
	}

	if order.PosPaymentInfo.OrderType != "" || store.IikoCloud.SendKitchenComments {
		if order.AllergyInfo == "" {
			order.AllergyInfo = models.NO
		}

		if len(iikoOrder.Items) != 0 {
			iikoOrder.Items[len(iikoOrder.Items)-1].Comment = kitchenComment
		}

		iikoOrder.OrderTypeID = order.PosPaymentInfo.OrderType

		if order.Type == models.ORDER_TYPE_INSTANT {
			iikoOrder.CompleteBefore = ""
		}
	}

	orderServiceType := iikoModels.SERVICE_TYPE_DELIVERY_BY_COURIER
	if order.IsPickedUpByCustomer || order.IsMarketplace {
		orderServiceType = iikoModels.SERVICE_TYPE_DELIVERY_BY_CLIENT
	}

	if order.PosPaymentInfo.OrderTypeService != "" {
		orderServiceType = order.PosPaymentInfo.OrderTypeService
	}

	if order.PosPaymentInfo.OrderType == "" {
		iikoOrder.OrderServiceType = orderServiceType
	}

	// Collect complete before date
	completeBeforeDate := order.EstimatedPickupTime.Value.Time

	if completeBeforeDate.IsZero() {
		completeBeforeDate = time.Now().UTC().Add(time.Hour)
	}

	// Convert UTC time to stores local time
	completeBeforeDate = completeBeforeDate.Add(time.Duration(store.Settings.TimeZone.UTCOffset)*time.Hour - 1*time.Minute)

	completeBefore := completeBeforeDate.Format("2006-01-02 15:04:05.000")

	switch orderServiceType {
	case iikoModels.SERVICE_TYPE_DELIVERY_BY_CLIENT, iikoModels.SERVICE_TYPE_DELIVERY_PICKUP:
		switch order.Type {
		case models.ORDER_TYPE_PREORDER:
			iikoOrder.CompleteBefore = completeBefore
		default:
			completeBefore = ""
		}
	default:
		iikoOrder.CompleteBefore = completeBefore
	}

	if orderServiceType == iikoModels.SERVICE_TYPE_DELIVERY_BY_COURIER {
		var streetName string
		// TODO:
		switch order.DeliveryService {
		case "glovo":
			streetName = "Доставка Глово"
		case "wolt":
			streetName = "Wolt доставка"
		case "chocofood":
			streetName = "Chocofood доставка"
		case "qr_menu":
			streetName = "Qr доставка"
		default:
			for _, deliveryService := range store.ExternalConfig {
				if deliveryService.Type == order.DeliveryService {
					streetName = fmt.Sprintf("%v доставка", deliveryService.Type)
				}
			}
		}

		iikoOrder.DeliveryPoint = &iikoModels.DeliveryPoint{
			Coordinates: &iikoModels.Coordinates{
				Longitude: order.DeliveryAddress.Longitude,
				Latitude:  order.DeliveryAddress.Latitude,
			},
			Address: &iikoModels.Address{
				Street: &iikoModels.Street{
					Name: streetName,
					City: store.Address.City,
				},
				House: "1",
			},
			Comment: order.DeliveryAddress.Label,
		}
	}

	// Init IIKO order
	transportToFrontTimeoutStr := manager.globalConfig.IIKOConfiguration.TransportToFrontTimeout
	transportToFrontTimeout := 0
	if transportToFrontTimeoutStr != "" {
		transportToFrontTimeout, err = strconv.Atoi(transportToFrontTimeoutStr)

		if err != nil {
			transportToFrontTimeout = 180
		}
	}

	return iikoModels.CreateDeliveryRequest{
		OrganizationID:  store.IikoCloud.OrganizationID,
		TerminalGroupID: store.IikoCloud.TerminalID,
		Order:           &iikoOrder,
		CreateOrderSettings: &iikoModels.CreateOrderSettings{
			TransportToFrontTimeout: transportToFrontTimeout,
		},
	}, order, nil
}

func (manager IIKOManager) constructPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	// for new integrate restaurants with combo
	if store.IikoCloud.HasCombo {
		return manager.ConstructPosOrderWithBase(ctx, order, store)
	}

	var isMarketplace bool
	switch order.DeliveryService {
	case "glovo":
		isMarketplace = store.Glovo.IsMarketplace
	case "wolt":
		isMarketplace = store.Wolt.IsMarketplace
	case "chocofood":
		isMarketplace = store.Chocofood.IsMarketplace
	case "qr_menu":
		isMarketplace = store.QRMenu.IsMarketplace
	case "express24":
		isMarketplace = store.Express24.IsMarketplace
	case "kwaaka_admin":
		isMarketplace = store.KwaakaAdmin.IsMarketPlace
	case "starter_app":
		isMarketplace = store.StarterApp.IsMarketPlace
	default:
		for _, deliveryService := range store.ExternalConfig {
			if deliveryService.Type == order.DeliveryService {
				isMarketplace = deliveryService.IsMarketplace
				if order.IsMarketplace && !isMarketplace {
					isMarketplace = true
				}
			}
		}
	}

	customerPhone := order.Customer.PhoneNumber
	if !strings.Contains(customerPhone, "+") {
		customerPhone = "+77771111111"
	}

	if store.Settings.CommentSetting.DefaultCourierPhone != "" && order.DeliveryService == models.YANDEX.String() {
		customerPhone = store.Settings.CommentSetting.DefaultCourierPhone
	}

	// Collect order products
	orderItems, delivery, serviceFee, err := manager.fulfillOrderProducts(ctx, order, store)
	if err != nil {
		return nil, order, err
	}

	if serviceFee != 0 {
		order.HasServiceFee = true
		order.ServiceFeeSum = serviceFee
	}

	payment, orderType, orderTypeService, err := manager.getPayment(ctx, order, store, delivery, serviceFee)

	if err != nil {
		log.Trace().Err(err).Msg("cant find payment for order")
		return nil, order, err
	}

	orderComment, customerComment, kitchenComment := ConstructOrderComments(ctx, order, store)

	// Init IIKO order
	iikoOrder := iikoModels.Order{
		Phone:   customerPhone,
		Comment: orderComment,
		// temp solve
		Customer: &iikoModels.Customer{
			Name:    order.Customer.Name,
			Comment: customerComment,
			Gender:  "NotSpecified",
		},
		Payments:      []iikoModels.Payment{payment},
		Items:         orderItems,
		DiscountsInfo: &iikoModels.DiscountsInfo{},
	}

	if orderType != "" || store.IikoCloud.SendKitchenComments {
		if order.AllergyInfo == "" {
			order.AllergyInfo = models.NO
		}

		if len(iikoOrder.Items) != 0 {
			iikoOrder.Items[len(iikoOrder.Items)-1].Comment = kitchenComment
		}

		iikoOrder.OrderTypeID = orderType

		if order.Type == models.ORDER_TYPE_INSTANT {
			iikoOrder.CompleteBefore = ""
		}
	}

	orderServiceType := iikoModels.SERVICE_TYPE_DELIVERY_BY_COURIER
	if order.IsPickedUpByCustomer || isMarketplace {
		orderServiceType = iikoModels.SERVICE_TYPE_DELIVERY_BY_CLIENT
	}

	if orderTypeService != "" {
		orderServiceType = orderTypeService
	}

	if orderTypeService == "" && orderType == "" {
		iikoOrder.OrderServiceType = orderServiceType
	}

	// Collect complete before date
	completeBeforeDate := order.EstimatedPickupTime.Value.Time.UTC()

	if completeBeforeDate.IsZero() {
		completeBeforeDate = time.Now().UTC().Add(time.Hour)
	}

	completeBeforeDate = completeBeforeDate.Add(time.Duration(store.Settings.TimeZone.UTCOffset) * time.Hour)

	completeBefore := completeBeforeDate.Format("2006-01-02 15:04:05.000")

	switch orderServiceType {
	case iikoModels.SERVICE_TYPE_DELIVERY_BY_CLIENT, iikoModels.SERVICE_TYPE_DELIVERY_PICKUP:
		switch order.Type {
		case models.ORDER_TYPE_PREORDER:
			iikoOrder.CompleteBefore = completeBefore
		default:
			completeBefore = ""
		}
	default:
		iikoOrder.CompleteBefore = completeBefore
	}

	var deliveryPoint iikoModels.DeliveryPoint
	if orderServiceType == iikoModels.SERVICE_TYPE_DELIVERY_BY_COURIER {
		var streetName string
		switch order.DeliveryService {
		case "glovo":
			streetName = "Доставка Глово"
		case "wolt":
			streetName = "Wolt доставка"
		case "chocofood":
			streetName = "Chocofood доставка"
		case "qr_menu":
			streetName = "Qr доставка"
		default:
			for _, deliveryService := range store.ExternalConfig {
				if deliveryService.Type == order.DeliveryService {
					streetName = fmt.Sprintf("%v доставка", deliveryService.Type)
				}
			}
		}

		deliveryPoint = iikoModels.DeliveryPoint{
			Coordinates: &iikoModels.Coordinates{
				Longitude: order.DeliveryAddress.Longitude,
				Latitude:  order.DeliveryAddress.Latitude,
			},
			Address: &iikoModels.Address{
				Street: &iikoModels.Street{
					Name: streetName,
					City: store.Address.City,
				},
				House: "1",
			},
			Comment: order.DeliveryAddress.Label,
		}

		iikoOrder.DeliveryPoint = &deliveryPoint
	}

	if order.PosPaymentInfo.OrderTypeService != "" && order.PosPaymentInfo.OrderTypeService != "DeliveryByCourier" {
		log.Info().Msgf("order type service=%s", order.PosPaymentInfo.OrderTypeService)
		iikoOrder.DeliveryPoint = nil
	}

	// Init IIKO order
	transportToFrontTimeoutStr := manager.globalConfig.IIKOConfiguration.TransportToFrontTimeout
	transportToFrontTimeout := 0
	if transportToFrontTimeoutStr != "" {
		transportToFrontTimeout, err = strconv.Atoi(transportToFrontTimeoutStr)

		if err != nil {
			transportToFrontTimeout = 180
		}
	}

	// Apply discounts on the order
	iikoOrder = manager.applyOrderDiscount(ctx, iikoOrder, order)

	return iikoModels.CreateDeliveryRequest{
		OrganizationID:  store.IikoCloud.OrganizationID,
		TerminalGroupID: store.IikoCloud.TerminalID,
		Order:           &iikoOrder,
		CreateOrderSettings: &iikoModels.CreateOrderSettings{
			TransportToFrontTimeout: transportToFrontTimeout,
		},
	}, order, nil
}

func (manager IIKOManager) sendOrder(ctx context.Context, order any, store coreStoreModels.Store) (any, error) {

	//var errs custom.Error
	posOrder, ok := order.(iikoModels.CreateDeliveryRequest)

	if !ok {
		return "", validator.ErrCastingPos
	}

	utils.Beautify("IIKO Request Body", posOrder)

	createResponse, err := manager.iikoClient.CreateDeliveryOrder(ctx, posOrder)
	if err != nil {
		log.Err(err).Msg("IIKOManager, sendOrder - IIKO error")
		//errs.Append(err, validator.ErrIgnoringPos)
		return "", err
	}

	if createResponse.OrderInfo != nil && createResponse.OrderInfo.ErrorInfo != nil {
		log.Err(err).Msgf("IIKOManager, sendOrder - IIKO Response ErrorInfo error: %+v", *createResponse.OrderInfo.ErrorInfo)
		//errs.Append(err, errors.New(createResponse.OrderInfo.ErrorInfo.Message))
		return "", err
	}

	return createResponse, nil
}

func (manager IIKOManager) CreateOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, error) {
	posOrder, order, err := manager.constructPosOrder(ctx, order, store)
	if err != nil {
		log.Trace().Err(validator.ErrCastingPos).Msg("")
		return order, err
	}

	utils.Beautify("successfully construct IIKO order(any)", posOrder)

	posOrderRequest, ok := posOrder.(iikoModels.CreateDeliveryRequest)
	if !ok {
		return order, validator.ErrCastingPos
	}

	utils.Beautify("successfully casting IIKO order(struct)", posOrderRequest)

	createResponse, err := manager.sendOrder(ctx, posOrderRequest, coreStoreModels.Store{})
	if err != nil {
		return order, err
	}

	utils.Beautify("successfully send order to IIKO, response (any)", createResponse)

	response, ok := createResponse.(iikoModels.CreateDeliveryResponse)
	if ok {
		order.PosOrderID = response.OrderInfo.ID
		order.CreationResult = models.CreationResult{
			CorrelationId: response.CorrelationID,
			OrderInfo: models.OrderInfo{
				ID:             response.OrderInfo.ID,
				OrganizationID: response.OrderInfo.OrganizationID,
				Timestamp:      int64(response.OrderInfo.Timestamp),
				CreationStatus: response.OrderInfo.CreationStatus,
			},
		}

		if response.OrderInfo.ErrorInfo != nil {
			order.CreationResult.ErrorDescription = response.OrderInfo.ErrorInfo.Description + response.OrderInfo.ErrorInfo.Message
		}

		utils.Beautify("successfully send order to IIKO, response (struct)", response)
	}

	// order, err = manager.RetrieveOrder(ctx, order, store.IikoCloud.OrganizationID)
	// if err != nil {
	//	return order, err
	// }

	utils.Beautify("finished order model result", order)

	return order, nil
}

func (manager IIKOManager) RetrieveOrder(ctx context.Context, order models.Order, organizationID string) (models.Order, error) {
	for i := 0; i < 4; i++ {
		time.Sleep(1 * time.Second)
		retrieveResponse, err := manager.iikoClient.RetrieveDeliveryOrder(ctx, organizationID, order.PosOrderID)
		if err != nil {
			log.Err(err).Msg("IIKO retrieve order failed error")
			return order, err
		}

		if retrieveResponse.CreationStatus == "Error" {
			log.Err(err).Msg("IIKO retrieve order, creation status is error")
			return order, errors.New(retrieveResponse.ErrorInfo.Description)
		}
	}

	return order, nil
}

func (manager IIKOManager) UpdateOrderProblem(ctx context.Context, organizationID, posOrderID string) error {

	problem := iikoModels.UpdateOrderProblem{
		OrderId:        posOrderID,
		OrganizationId: organizationID,
		HasProblem:     true,
		Problem:        "ОТМЕНА ЗАКАЗА",
	}

	err := manager.iikoClient.UpdateOrderProblem(ctx, problem)
	if err != nil {
		log.Err(err).Msgf("error: UpdateOrderProblem for pos_order_id: %s", posOrderID)
		return err
	}

	return nil
}
