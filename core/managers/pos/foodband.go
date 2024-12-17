package pos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/config"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/domain/foodband"
	MenuUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	foodbandClient "github.com/kwaaka-team/orders-core/pkg/posintegration"
	foodbandConf "github.com/kwaaka-team/orders-core/pkg/posintegration/clients"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

type FoodBandManager struct {
	foodbandCli        foodbandConf.FOODBAND
	productsMap        map[string]coreMenuModels.Product
	attributesMap      map[string]coreMenuModels.Attribute
	attributeGroupsMap map[string]coreMenuModels.AttributeGroup
	promoMap           map[string]string
	aggregatorMenu     coreMenuModels.Menu
	storeID            string
}

func NewFoodBandManager(globalConfig config.Configuration, store coreStoreModels.Store, menu coreMenuModels.Menu, promo coreMenuModels.Promo, aggregatorMenu coreMenuModels.Menu) (FoodBandManager, error) {
	foodBandCfg, err := getFoodBandCfg(store.ExternalConfig)
	if err != nil {
		log.Trace().Err(err).Msgf("Cant initialize %s Client in %s restaurant", store.PosType, store.ID)
		return FoodBandManager{}, err
	}
	retryMaxCount, err := strconv.Atoi(globalConfig.RetryConfiguration.Count)
	if err != nil {
		retryMaxCount = 0
	}

	foodBandClient, err := foodbandClient.NewClient(&foodbandConf.Config{
		CancelOrderUrl: foodBandCfg.WebhookConfig.OrderCancel,
		CreateOrderUrl: foodBandCfg.WebhookConfig.OrderCreate,
		ApiToken:       foodBandCfg.ClientSecret,
		RetryMaxCount:  retryMaxCount,
	})

	if err != nil {
		log.Trace().Err(err).Msgf("Cant initialize %s Client.", store.PosType)
		return FoodBandManager{}, err
	}

	return FoodBandManager{
		foodbandCli:        foodBandClient,
		productsMap:        MenuUtils.ProductsMap(menu),
		attributesMap:      MenuUtils.AtributesMap(menu),
		attributeGroupsMap: MenuUtils.AtributeGroupsMap(menu),
		promoMap:           map[string]string{},
		aggregatorMenu:     aggregatorMenu,
		storeID:            foodBandCfg.StoreID[0],
	}, nil
}

func getFoodBandCfg(externalConfig []coreStoreModels.StoreExternalConfig) (coreStoreModels.StoreExternalConfig, error) {
	for _, ext := range externalConfig {
		if ext.Type == models.FoodBand.String() {
			return ext, nil
		}
	}
	return coreStoreModels.StoreExternalConfig{}, fmt.Errorf("not found foodband cfg in restaurant")
}

func (manager FoodBandManager) CancelOrder(ctx context.Context, order models.Order, cancelReason, paymentStrategy string, store coreStoreModels.Store) error {
	err := manager.foodbandCli.CancelOrder(ctx, foodband.CancelOrderRequest{
		StoreID:         manager.storeID,
		OrderID:         order.PosOrderID,
		DeliveryService: order.DeliveryService,
		CancelReason:    cancelReason,
		PaymentStrategy: paymentStrategy,
	})
	if err != nil {
		log.Err(err).Msg("FOODBAND cancel order error")
		return err
	}

	log.Err(err).Msg("FOODBAND cancel order success")

	return nil
}

func (manager FoodBandManager) sendOrder(ctx context.Context, order any, store coreStoreModels.Store) (any, error) {
	posOrder, ok := order.(foodband.CreateOrderRequest)
	if !ok {
		return 0, validator.ErrCastingPos
	}
	utils.Beautify("FOODBAND Request Body", posOrder)

	retryCount, err := manager.foodbandCli.CreateOrder(ctx, posOrder)
	if err != nil {

		return retryCount, errors.Wrap(err, "foodband create order")
	}

	return retryCount, nil
}

func (manager FoodBandManager) constructOrderComment(ctx context.Context, order models.Order, store coreStoreModels.Store) string {
	var (
		commentName  = "Комментарий"
		allergyName  = "Аллергия"
		quantityName = "Количество персон"
	)

	commentSettings := store.Settings.CommentSetting

	if commentSettings.HasCommentSetting {
		commentName = commentSettings.CommentName
		allergyName = commentSettings.Allergy
		quantityName = commentSettings.QuantityPerson
	}

	orderComment := ""

	if order.SpecialRequirements != "" {
		orderComment = fmt.Sprintf("%s: %s. ", commentName, order.SpecialRequirements)
	}
	if order.AllergyInfo != "" {
		orderComment = fmt.Sprintf("%s%s: %s.", orderComment, allergyName, order.AllergyInfo)
	}
	if order.Persons != 0 {
		orderComment = fmt.Sprintf("%s%s: %d", orderComment, quantityName, order.Persons)
	}

	return orderComment
}

func (manager FoodBandManager) constructPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	orderComment := manager.constructOrderComment(ctx, order, store)
	products, err := manager.fulfillOrderProducts(ctx, order, store)
	if err != nil {
		return nil, order, err
	}

	payment, err := getPayment(order, store)
	if err != nil {
		return nil, order, err
	}
	deliveryProviderType := manager.getDeliveryProviderType(order)

	orderFoodBand := foodband.Order{
		ID:              uuid.New().String(),
		Type:            order.Type,
		Code:            order.OrderCode,
		PickUpCode:      order.PickUpCode,
		CompleteBefore:  collectCompleteBeforeDate(store, order, deliveryProviderType),
		Phone:           getPhone(order.Customer.PhoneNumber),
		DeliveryService: order.DeliveryService,
		DeliveryPoint:   getDeliveryPoint(order.DeliveryAddress, deliveryProviderType),
		Comment:         orderComment,
		Customer: foodband.Customer{
			Name: order.Customer.Name,
		},
		Courier: foodband.Courier{
			Name:        order.Courier.Name,
			PhoneNumber: order.Courier.PhoneNumber,
		},
		Products:             products,
		Payments:             []foodband.Payment{payment},
		DeliveryFee:          order.DeliveryFee.Value,
		DeliveryProviderType: deliveryProviderType,
	}

	return foodband.CreateOrderRequest{
		StoreID: manager.storeID,
		Order:   orderFoodBand,
	}, order, nil
}

func (manager FoodBandManager) CreateOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, error) {
	posOrder, _, err := manager.constructPosOrder(ctx, order, store)
	if err != nil {
		log.Trace().Err(validator.ErrCastingPos).Msg("")
		return order, err
	}
	utils.Beautify("successfully construct FOODBAND order(any)", posOrder)

	posOrderRequest, ok := posOrder.(foodband.CreateOrderRequest)
	if !ok {
		return order, validator.ErrCastingPos
	}
	utils.Beautify("successfully casting FOODBAND order(struct)", posOrderRequest)

	resp, err := manager.sendOrder(ctx, posOrderRequest, coreStoreModels.Store{})
	retryCount := resp.(int)
	if err != nil {
		return models.Order{
			PosOrderID: posOrderRequest.Order.ID,
			CreationResult: models.CreationResult{
				OrderInfo: models.OrderInfo{
					ID:             posOrderRequest.Order.ID,
					OrganizationID: posOrderRequest.StoreID,
					CreationStatus: "ERROR",
				},
				ErrorDescription: err.Error(),
			},
			RetryCount: retryCount,
			IsRetry:    retryCount != 0,
		}, err
	}

	if ok {
		order.PosOrderID = posOrderRequest.Order.ID
		order.CreationResult = models.CreationResult{
			OrderInfo: models.OrderInfo{
				ID:             posOrderRequest.Order.ID,
				OrganizationID: posOrderRequest.StoreID,
			},
		}
		order.RetryCount = retryCount
		order.IsRetry = retryCount != 0
	}

	utils.Beautify("finished order model result", order)

	return order, nil
}

func (manager FoodBandManager) GetOrderStatus(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, error) {
	return "", errs.ErrUnsupportedMethod
}

func (manager FoodBandManager) getDeliveryProviderType(order models.Order) string {
	if order.IsPickedUpByCustomer {
		return models.FOODBAND_CUSTOMER_PICKUP
	}
	if order.RestaurantSelfDelivery {
		return models.FOODBAND_DELIVERY_RESTAURANT
	}
	return models.FOODBAND_DELIVERY_AGGREGATOR
}

func collectCompleteBeforeDate(store coreStoreModels.Store, order models.Order, deliveryProviderType string) string {
	completeBeforeDate := order.EstimatedPickupTime.Value.Time
	if completeBeforeDate.IsZero() {
		completeBeforeDate = time.Now().UTC().Add(time.Hour)
	}
	completeBeforeDate = completeBeforeDate.Add(time.Duration(store.Settings.TimeZone.UTCOffset)*time.Hour - 1*time.Minute)

	completeBefore := completeBeforeDate.Format("2006-01-02 15:04:05.000")

	if deliveryProviderType == models.FOODBAND_DELIVERY_RESTAURANT || deliveryProviderType == models.FOODBAND_CUSTOMER_PICKUP {
		if order.Type == models.ORDER_TYPE_PREORDER {
			return completeBefore
		}
		return ""
	}

	return completeBefore
}

func getPhone(customerPhone string) string {
	if !strings.Contains(customerPhone, "+") {
		customerPhone = "+77771111111"
	}
	return customerPhone
}

func getPayment(order models.Order, store coreStoreModels.Store) (foodband.Payment, error) {
	var paymentTypes *coreStoreModels.DeliveryServicePaymentType
	switch order.DeliveryService {
	case models.GLOVO.String():
		paymentTypes = &store.Glovo.PaymentTypes
	case models.WOLT.String():
		paymentTypes = &store.Wolt.PaymentTypes
	case models.CHOCOFOOD.String():
		paymentTypes = &store.Chocofood.PaymentTypes
	case models.QRMENU.String():
		paymentTypes = &store.QRMenu.PaymentTypes
	default:
		for _, deliveryService := range store.ExternalConfig {
			if deliveryService.Type == order.DeliveryService {
				currentPaymentTypes := deliveryService.PaymentTypes
				paymentTypes = &currentPaymentTypes
			}
		}
	}

	if paymentTypes == nil {
		return foodband.Payment{}, fmt.Errorf("no any payment types")
	}

	var paymentType string

	switch order.PaymentMethod {
	case models.PAYMENT_METHOD_DELAYED:
		paymentType = models.PAYMENT_METHOD_CARD
	case models.PAYMENT_METHOD_CASH:
		paymentType = models.PAYMENT_METHOD_CASH
	default:
		log.Info().Msgf("Payment Method: %v", order.PaymentMethod)
	}

	return foodband.Payment{
		PaymentTypeKind: paymentType,
		Sum:             order.EstimatedTotalPrice.Value - order.PartnerDiscountsProducts.Value,
	}, nil
}

func (manager FoodBandManager) fulfillOrderProducts(ctx context.Context, order models.Order, store coreStoreModels.Store) ([]foodband.Product, error) {
	var items []foodband.Product

	aggregatorProducts, aggregatorAttributes := ActiveMenuPositions(ctx, manager.aggregatorMenu)

	// Add products to order

	for _, product := range order.Products {

		var modifiers []foodband.Attribute
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

			attributePosID, exist_ := aggregatorAttributes[attribute.ID]
			if exist_ {
				attribute.ID = attributePosID
			} else {
				log.Warn().Msgf("Attribute with ID %s %s not matched", attribute.ID, attribute.Name)
			}

			menuAttribute, ok := manager.attributesMap[attribute.ID]

			if ok {
				modifiersIds[attribute.ID] = attribute.ID

				var priceAmount float64
				if store.Settings.PriceSource == models.POSPriceSource {
					priceAmount = menuAttribute.Price
				} else {
					priceAmount = attribute.Price.Value
				}

				itemModifier := foodband.Attribute{
					ID:       attribute.ID,
					Quantity: attribute.Quantity,
					Price:    priceAmount,
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

					items = append(items, foodband.Product{
						ID:       menuAttributeProduct.ProductID,
						Price:    priceAmount,
						Quantity: attribute.Quantity * product.Quantity,
					})
				} else {
					return nil, fmt.Errorf("ATTRIBUTE NOT FOUND IN POS MENU, ID %s, NAME %s", attribute.ID, attribute.Name)
				}
			}

		}

		menuProduct, ok := manager.productsMap[product.ID]

		if !ok {
			log.Info().Msgf("PRODUCT NOT FOUND IN POS MENU, ID %s, NAME %s", product.ID, product.Name)

			id, exist_ := manager.promoMap[product.ID]

			var isGift string
			if len(product.Promos) != 0 {
				isGift = product.Promos[0].Type
			}

			if !exist_ || isGift != "GIFT" {
				return nil, errors.Wrap(errs.ErrProductNotFound, fmt.Sprintf("PRODUCT NOT FOUND IN POS MENU, ID %s, NAME %s", product.ID, product.Name))
			}

			orderItem := foodband.Product{
				ID:         id,
				Price:      0, // ?
				Quantity:   product.Quantity,
				Attributes: modifiers,
			}
			items = append(items, orderItem)
			continue
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

				defaultModifier := foodband.Attribute{
					ID:       defaultAttribute.ExtID,
					Quantity: amount,
					Price:    float64(price),
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
			case models.PROMO_TYPE_FIXED:
				itemPrice = itemPrice - float64(promo.Discount)
			case models.PROMO_TYPE_PERCENTAGE:
				itemPrice = itemPrice - float64(promo.Discount)
			}
		}

		orderItem := foodband.Product{
			ID:         menuProduct.ProductID,
			Price:      itemPrice,
			Quantity:   product.Quantity,
			Attributes: modifiers,
		}

		if menuProduct.ProductID == "" {
			orderItem.ID = menuProduct.ExtID
		}

		items = append(items, orderItem)
	}

	return items, nil
}

func getDeliveryPoint(deliveryAddress models.DeliveryAddress, orderServiceType string) foodband.DeliveryPoint {
	if orderServiceType == models.FOODBAND_DELIVERY_RESTAURANT {
		return foodband.DeliveryPoint{
			Coordinates: foodband.Coordinates{
				Latitude:  deliveryAddress.Latitude,
				Longitude: deliveryAddress.Longitude,
			},
			AddressLabel: deliveryAddress.Label,
		}
	}
	return foodband.DeliveryPoint{}
}

func (manager FoodBandManager) UpdateOrderProblem(ctx context.Context, organizationID, posOrderID string) error {
	return ErrUnsupportedMethod
}
