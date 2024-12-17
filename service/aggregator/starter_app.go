package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	starterAppModels "github.com/kwaaka-team/orders-core/core/starter_app/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models3 "github.com/kwaaka-team/orders-core/core/wolt/models"
	"github.com/kwaaka-team/orders-core/pkg/starterapp"
	starterappCli "github.com/kwaaka-team/orders-core/pkg/starterapp/clients"
	"github.com/kwaaka-team/orders-core/pkg/starterapp/clients/dto"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type starterAppService struct {
	cli         starterappCli.StarterApp
	menuService *menu.Service
}

func newStarterAppService(baseUrl string, store storeModels.Store, menuService *menu.Service) (*starterAppService, error) {
	cli, err := starterapp.NewStarterAppClient(&starterappCli.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		ApiKey:   store.StarterApp.ApiKey,
	})
	if err != nil {
		log.Trace().Err(err).Msg("can't initialize StarterApp menu client ")
		return nil, errors.Wrap(constructorError, err.Error())
	}

	return &starterAppService{
		cli,
		menuService,
	}, nil
}

func (s starterAppService) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error) {
	shopID, err := strconv.Atoi(aggregatorStoreID)
	if err != nil {
		return "", err
	}
	req := make([]dto.MealOfferRequest, 0)
	for i := range products {
		quantity := 0
		if isAvailable {
			quantity = 1000
		}
		id, err := strconv.Atoi(products[i].StarterAppOfferID)
		if err != nil {
			log.Err(err).Msgf("can't parse product starter app offer id UpdateStopListByProducts, product ext id: %s", products[i].ExtID)
			continue
		}
		mealId, err := strconv.Atoi(products[i].StarterAppID)
		if err != nil {
			log.Err(err).Msgf("can't parse product starter app id UpdateStopListByProducts, product ext id: %s", products[i].ExtID)
			continue
		}
		if len(products[i].Price) == 0 {
			log.Err(err).Msgf("price len 0 UpdateStopListByProducts, product ext id: %s", products[i].ExtID)
			continue
		}
		req = append(req, dto.MealOfferRequest{
			ID:       int64(id),
			PosId:    products[i].ExtID,
			Quantity: quantity,
			Price:    products[i].Price[0].Value,
			InMenu:   isAvailable,
			MealId:   int64(mealId),
		})
	}
	return "", s.cli.UpdateMealOffers(ctx, req, shopID)
}

func (s starterAppService) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error) {
	shopID, err := strconv.Atoi(aggregatorStoreID)
	if err != nil {
		return "", err
	}
	req := make([]dto.MealOfferRequest, 0)
	for i := range products {
		quantity := 0
		if products[i].IsAvailable {
			quantity = 1000
		}
		id, err := strconv.Atoi(products[i].StarterAppOfferID)
		if err != nil {
			log.Err(err).Msgf("can't parse product starter app offer id UpdateStopListByProducts, product ext id: %s", products[i].ExtID)
			continue
		}
		mealId, err := strconv.Atoi(products[i].StarterAppID)
		if err != nil {
			log.Err(err).Msgf("can't parse product starter app id UpdateStopListByProducts, product ext id: %s", products[i].ExtID)
			continue
		}
		if len(products[i].Price) == 0 {
			log.Err(err).Msgf("price len 0 UpdateStopListByProducts, product ext id: %s", products[i].ExtID)
			continue
		}
		req = append(req, dto.MealOfferRequest{
			ID:       int64(id),
			PosId:    products[i].ExtID,
			Quantity: quantity,
			Price:    products[i].Price[0].Value,
			InMenu:   products[i].IsAvailable,
			MealId:   int64(mealId),
		})
	}
	return "", s.cli.UpdateMealOffers(ctx, req, shopID)
}

func (s starterAppService) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error) {
	shopID, err := strconv.Atoi(aggregatorStoreID)
	if err != nil {
		return "", err
	}
	req := make([]dto.MealOfferRequest, 0)
	for i := range attributes {
		quantity := 0
		if attributes[i].IsAvailable {
			quantity = 1000
		}
		id, err := strconv.Atoi(attributes[i].StarterAppOfferID)
		if err != nil {
			log.Err(err).Msgf("can't parse attributes starter app offer id UpdateStopListByProducts, attributes ext id: %s", attributes[i].ExtID)
			continue
		}
		mealId, err := strconv.Atoi(attributes[i].StarterAppID)
		if err != nil {
			log.Err(err).Msgf("can't parse attributes starter app id UpdateStopListByProducts, attributes ext id: %s", attributes[i].ExtID)
			continue
		}

		req = append(req, dto.MealOfferRequest{
			ID:       int64(id),
			PosId:    attributes[i].ExtID,
			Quantity: quantity,
			Price:    attributes[i].Price,
			InMenu:   attributes[i].IsAvailable,
			MealId:   int64(mealId),
		})
	}
	return "", s.cli.UpdateMealOffers(ctx, req, shopID)
}

func (s starterAppService) MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	log.Info().Msg("delivery service: starter_app")

	switch posStatus {
	case models.ACCEPTED, models.READY_FOR_COOKING, models.WAIT_SENDING:
		return models.Created.String()
	case models.COOKING_STARTED:
		return models.InProgress.String()
	case models.COOKING_COMPLETE, models.CLOSED, models.READY_FOR_PICKUP, models.ON_WAY, models.OUT_FOR_DELIVERY:
		return models.Cooked.String()
	case models.CANCELLED_BY_POS_SYSTEM:
		return models.Canceled.String()
	}

	return ""
}

func (s starterAppService) UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error {
	if aggregatorStatus == "" {
		log.Info().Msgf("aggregator status is empty, order id %s", order.OrderID)
		return nil
	}

	switch aggregatorStatus {
	case models.Canceled.String():
		if err := s.cli.SendOrderErrorNotification(ctx, dto.SendOrderErrorNotificationRequest{
			IsOrderSent: false,
			IsPosError:  true,
			Error: dto.SendOrderErrorNotificationError{
				Message:  order.CancelReason.Reason,
				Request:  "",
				Response: "",
			},
		}, order.OrderID); err != nil {
			return err
		}

	default:
		if err := s.cli.ChangeOrderStatus(ctx, dto.ChangeOrderStatusRequest{
			Status: aggregatorStatus,
		}, order.OrderID); err != nil {
			return err
		}
	}

	return nil
}

func (s starterAppService) GetSystemCreateOrderRequestByAggregatorRequest(req interface{}, store storeModels.Store) (models.Order, error) {
	r, ok := req.(starterAppModels.Order)
	if !ok {
		return models.Order{}, errors.New("casting error")
	}

	orderCreatedAt := time.Now().UTC()
	cookingTime := r.CookingTime
	estimatedPickupTime := orderCreatedAt.Add(time.Duration(cookingTime) * time.Minute)

	isPickUpByCustomer, err := s.getIsPickedUpByCustomer(r.DeliveryType)
	if err != nil {
		return models.Order{}, err
	}

	orderType := "INSTANT"
	var preorder models.PreOrder
	if r.IsPreorder {
		// r.DeliveryDatetime in UTC
		if r.DeliveryDatetime.Before(orderCreatedAt) {
			return models.Order{}, errors.New("preorder time is before than now")
		}

		orderType = "PREORDER"
		estimatedPickupTime = r.DeliveryDatetime.Add(-time.Minute * time.Duration(r.DeliveryDuration))

		if isPickUpByCustomer {
			estimatedPickupTime = r.DeliveryDatetime
			customerPickupMinTime := orderCreatedAt.Add(time.Duration(cookingTime) * time.Minute)
			if estimatedPickupTime.Before(customerPickupMinTime) {
				log.Info().Msgf("estimated pickup time before order cooking complete time, estimated pickeup time: %v, cooking complete: %v\n", estimatedPickupTime, orderCreatedAt)
				return models.Order{}, fmt.Errorf("estimated pickup time before cooking complete time, estimated pickeup time: %v, cooking complete: %v\n", estimatedPickupTime, orderCreatedAt)
			}
		}

		if estimatedPickupTime.Before(orderCreatedAt) {
			log.Info().Msgf("estimated pickup time before order created time, estimated pickeup time: %v, created: %v\n", estimatedPickupTime, orderCreatedAt)
			return models.Order{}, fmt.Errorf("estimated pickup time before order created time, estimated pickeup time: %v, created: %v\n", estimatedPickupTime, orderCreatedAt)
		}

		preorder = models.PreOrder{
			Time: models.TransactionTime{
				Value:    models.Time{Time: r.DeliveryDatetime},
				TimeZone: store.Settings.TimeZone.TZ,
			},
			Status: "waiting",
		}

	}

	products, err := s.toOrderProducts(r)
	if err != nil {
		return models.Order{}, err
	}

	order := models.Order{
		OrderID:         r.GlobalId,
		StoreID:         r.ShopId,
		OrderCode:       r.StarterId,
		PickUpCode:      r.StarterId,
		RestaurantID:    store.ID,
		PosType:         store.PosType,
		Type:            orderType,
		DeliveryService: models.STARTERAPP.String(),
		Status:          models.STATUS_NEW.ToString(),
		StatusesHistory: []models.OrderStatusUpdate{
			{
				Name: models.STATUS_NEW.ToString(),
				Time: orderCreatedAt,
			},
		},
		OrderTime: models.TransactionTime{
			Value:    models.Time{Time: orderCreatedAt},
			TimeZone: store.Settings.TimeZone.TZ,
		},
		EstimatedPickupTime: models.TransactionTime{
			Value:    models.Time{Time: estimatedPickupTime},
			TimeZone: store.Settings.TimeZone.TZ,
		},
		Preorder:         preorder,
		UtcOffsetMinutes: strconv.Itoa(int(store.Settings.TimeZone.UTCOffset)),
		Currency:         store.Settings.Currency,
		AllergyInfo:      r.Comment,
		EstimatedTotalPrice: models.Price{
			Value:        r.TotalPrice,
			CurrencyCode: store.Settings.Currency,
		},
		TotalCustomerToPay: models.Price{
			Value:        r.TotalPrice,
			CurrencyCode: store.Settings.Currency,
		},
		DeliveryFee: models.Price{
			Value:        r.DeliveryPrice,
			CurrencyCode: store.Settings.Currency,
		},
		Customer: models.Customer{
			Name:        r.Username,
			PhoneNumber: r.UserPhone,
		},
		DeliveryAddress: models.DeliveryAddress{
			Label:     r.Address.Street,
			Longitude: r.Address.Longitude,
			Latitude:  r.Address.Latitude,
			City:      r.Address.City,
			Comment:   r.Address.Comment,
			Flat:      r.Address.Flat,
			Floor:     r.Address.Floor,
			Porch:     r.Address.Entrance,
		},
		PaymentMethod:        s.getPaymentMethod(r.PaymentType),
		IsPickedUpByCustomer: isPickUpByCustomer,
		Products:             products,
	}

	return order, nil
}

func (s starterAppService) getPaymentMethod(paymentMethod string) string {
	if paymentMethod == "cash" {
		return models.PAYMENT_METHOD_CASH
	}
	return models.PAYMENT_METHOD_DELAYED
}

func (s starterAppService) getIsPickedUpByCustomer(deliveryType string) (bool, error) {
	switch deliveryType {
	case "pickup":
		return true, nil
	case "courier":
		return false, nil
	}

	return false, errors.New("unsupported delivery type")
}

func (s starterAppService) toOrderProducts(req starterAppModels.Order) ([]models.OrderProduct, error) {
	products := make([]models.OrderProduct, 0, len(req.OrderItems))
	for _, item := range req.OrderItems {

		product := models.OrderProduct{
			ID:   item.ExtID,
			Name: item.Name,
			Price: models.Price{
				Value: item.TotalPrice,
			},
			Quantity: item.Quantity,
		}

		product.Attributes = make([]models.ProductAttribute, 0, len(item.Modifiers))
		for _, mod := range item.Modifiers {
			product.Attributes = append(product.Attributes, models.ProductAttribute{
				ID: mod.ExtId,
				Price: models.Price{
					Value: float64(mod.Price),
				},
				Quantity: mod.Amount,
				Name:     mod.Name,
			})
		}

		products = append(products, product)
	}

	return products, nil
}

func (s starterAppService) GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error) {
	return storeModels.AggregatorSchedule{}, errors.New("method not implemented")
}

func (s starterAppService) GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error) {
	return false, errors.New("method not implemented")
}

func (s starterAppService) OpenStore(ctx context.Context, aggregatorStoreId string) error {
	return errors.New("method not implemented")
}

func (s starterAppService) IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error) {
	return restaurantSelfDelivery != true, nil
}

func (s starterAppService) SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	order, ok := req.(starterAppModels.Order)
	if !ok {
		return nil, errors.New("casting error")
	}

	childRestaurantOrders := make(map[string]starterAppModels.Order)

	for i := range order.OrderItems {
		product := order.OrderItems[i]

		productRestaurantID, productID, err := splitVirtualStoreItemID(product.ExtID, "_")
		if err != nil {
			log.Err(errors.New("not valid signature, len will be 2 with _")).Msgf("orders core, splitVirtualStoreItemID error, id: %s", product.ExtID)
			continue
		}
		childProduct := product
		childProduct.ExtID = productID
		childOrder, ok := childRestaurantOrders[productRestaurantID]
		if !ok {
			child := order
			child.OrderItems = make([]starterAppModels.OrderItem, 0, 1)
			child.OrderItems = append(child.OrderItems, childProduct)
			child.ShopId = child.ShopId + "_" + productRestaurantID
			child.GlobalId = order.GlobalId + "_" + productRestaurantID

			childRestaurantOrders[productRestaurantID] = child
			continue
		}
		childOrder.OrderItems = append(childOrder.OrderItems, childProduct)
		childRestaurantOrders[productRestaurantID] = childOrder
	}

	res := make([]interface{}, 0, len(childRestaurantOrders))
	for _, val := range childRestaurantOrders {
		res = append(res, s.calculateChildOrderTotalSum(val))
	}

	return res, nil
}

func (s starterAppService) calculateChildOrderTotalSum(order starterAppModels.Order) starterAppModels.Order {
	totalSum := float64(0)

	for _, p := range order.OrderItems {
		totalSum = totalSum + p.TotalPrice
	}

	order.Price = totalSum
	order.TotalPrice = totalSum
	return order
}

func (s starterAppService) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	order, ok := req.(starterAppModels.Order)
	if !ok {
		return "", errors.New("casting error")
	}

	return order.ShopId, nil
}

func (s starterAppService) GetAggregatorOrder(ctx context.Context, orderID string) (models3.Order, error) {
	return models3.Order{}, errors.New("method not implemented")
}

func (s starterAppService) SendOrderErrorNotification(ctx context.Context, req interface{}) error {
	var starterReq starterAppModels.OrderDto
	order, ok := req.(models.Order)
	if !ok {
		starterReq, ok = req.(starterAppModels.OrderDto)
		if !ok {
			return fmt.Errorf("send order notification for starter app casting error")
		}
	}
	type message struct {
		Code              string `json:"code"`
		HumanReadableText string `json:"humanReadableText"`
		Message           string `json:"message"`
	}
	type response struct {
		Error string `json:"error"`
	}
	m := message{
		Code:              "Common",
		HumanReadableText: "Ошибка при отправке заказа",
		Message:           order.FailReason.Message,
	}
	messageJsonData, err := json.Marshal(m)
	if err != nil {
		return errors.Errorf("error while marshalling starter app error message")
	}
	starterReqJsonData, err := json.Marshal(starterReq)
	if err != nil {
		return errors.Errorf("error while marshalling starter app error starterReq")
	}
	r := response{
		Error: order.FailReason.Message,
	}
	responseJsonData, err := json.Marshal(r)
	if err != nil {
		return errors.Errorf("error while marshalling starter app error response")
	}

	var errMsg dto.SendOrderErrorNotificationError
	if order.OrderID != "" {

		errMsg = dto.SendOrderErrorNotificationError{
			Message:  string(messageJsonData),
			Request:  order.LogMessages.FromDelivery,
			Response: string(responseJsonData),
		}
	} else {
		errMsg = dto.SendOrderErrorNotificationError{
			Message:  string(messageJsonData),
			Request:  string(starterReqJsonData),
			Response: string(responseJsonData),
		}
	}

	if err := s.cli.SendOrderErrorNotification(ctx, dto.SendOrderErrorNotificationRequest{
		IsOrderSent: false,
		IsPosError:  true,
		Error:       errMsg,
	}, order.OrderID); err != nil {
		return fmt.Errorf("send order notification for starter app error: %+v", err)
	}

	return nil
}

func (s starterAppService) SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error {
	return nil
}
