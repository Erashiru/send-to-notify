package aggregator

import (
	"context"
	"fmt"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models2 "github.com/kwaaka-team/orders-core/core/wolt/models"
	"github.com/kwaaka-team/orders-core/core/wolt/models_v2"
	"github.com/kwaaka-team/orders-core/pkg/wolt/clients"
	woltModels "github.com/kwaaka-team/orders-core/pkg/wolt/clients/dto"
	httpclient "github.com/kwaaka-team/orders-core/pkg/wolt/clients/http"
	"github.com/kwaaka-team/orders-core/service/menu"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math"
	"strconv"
	"strings"
	"time"
)

type woltService struct {
	cli         clients.Wolt
	menuService *menu.Service
}

func newWoltService(baseUrl string, store storeModels.Store, menuService *menu.Service) (*woltService, error) {
	woltCli, err := httpclient.NewClient(&clients.Config{
		ApiKey:   store.Wolt.ApiKey,
		Protocol: "http",
		StoreID:  store.ID,
		BaseURL:  baseUrl,
		Username: store.Wolt.MenuUsername,
		Password: store.Wolt.MenuPassword,
	})
	if err != nil {
		return nil, err
	}

	return &woltService{
		woltCli,
		menuService,
	}, nil
}

func (s *woltService) OpenStore(ctx context.Context, aggregatorStoreId string) error {
	err := s.cli.ManageStore(ctx, woltModels.IsStoreOpen{
		AvailableStore: "ONLINE",
		VenueId:        aggregatorStoreId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *woltService) GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error) {
	status, err := s.cli.GetStoreStatus(ctx, aggregatorStoreId)
	if err != nil {
		return false, err
	}
	//return status.Status.IsOpen, nil
	if !status.Status.IsOnline || !status.Status.IsOpen {
		return false, nil
	}
	return true, nil
}

func (s *woltService) getDayMapping() map[string]int {
	return map[string]int{
		"SUNDAY":    1,
		"MONDAY":    2,
		"TUESDAY":   3,
		"WEDNESDAY": 4,
		"THURSDAY":  5,
		"FRIDAY":    6,
		"SATURDAY":  7,
	}
}

func (s *woltService) toSystemSchedule(schedule []woltModels.OpeningTimes) []storeModels.Schedule {
	result := make([]storeModels.Schedule, 0, len(schedule))

	dayMapping := s.getDayMapping()

	for _, cur := range schedule {
		dayOfWeek, ok := dayMapping[cur.OpeningDay]
		if !ok {
			log.Err(errors.New("day mapping is wrong")).Msgf("opening day: %s", cur.OpeningDay)
			continue
		}

		result = append(result, storeModels.Schedule{
			DayOfWeek: dayOfWeek,
			TimeSlots: []storeModels.TimeSlot{
				{
					Opening: cur.OpeningTime,
					Closing: cur.ClosingTime,
				},
			},
		})
	}

	return result
}

func (s *woltService) GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error) {
	storeStatus, err := s.cli.GetStoreStatus(ctx, aggregatorStoreId)
	if err != nil {
		return storeModels.AggregatorSchedule{}, err
	}

	systemSchedule := s.toSystemSchedule(storeStatus.OpeningTimes)

	return storeModels.AggregatorSchedule{
		Schedule: systemSchedule,
	}, nil
}

func (s *woltService) attributesToWoltModelsBulk(attributes []menuModels.Attribute) woltModels.UpdateAttributes {
	res := woltModels.UpdateAttributes{
		Attribute: make([]woltModels.UpdateAttribute, 0, len(attributes)),
	}

	for i := 0; i < len(attributes); i++ {
		res.Attribute = append(res.Attribute, woltModels.UpdateAttribute{
			ExtID:       attributes[i].ExtID,
			IsAvailable: &attributes[i].IsAvailable,
		})
	}

	return res

}

func (s *woltService) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error) {
	result, err := s.cli.BulkAttribute(ctx, aggregatorStoreID, s.attributesToWoltModelsBulk(attributes))
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *woltService) IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error) {
	return restaurantSelfDelivery != true, nil
}

func (s *woltService) SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	webhook, ok := req.(models2.OrderNotification)
	if !ok {
		return nil, errors.New("casting error")
	}

	if webhook.Body.Status != models2.Created.String() {
		log.Info().Msgf("successfully created order, id=%s", webhook.Body.Id)
		return nil, errors.New(fmt.Sprintf("webhook action is not '%s'", models2.Created))
	}

	order, err := s.cli.GetOrder(context.Background(), webhook.Body.Id)
	if err != nil {
		return nil, err
	}

	utils.Beautify("getting virtual store order from wolt body", order)

	childRestaurantOrders := make(map[string]models2.Order)

	for i := range order.Items {
		product := order.Items[i]
		productRestaurantID, productID, err := splitVirtualStoreItemID(product.PosID, "_")
		if err != nil {
			log.Err(errors.New("not valid signature, len will be 2 with _")).Msgf("orders core, splitVirtualStoreItemID error, id: %s", product.PosID)
			continue
		}

		childProduct := product
		childProduct.PosID = productID
		childOrder, ok := childRestaurantOrders[productRestaurantID]
		if !ok {
			child := order
			child.Items = make([]models2.OrderItem, 0, 1)
			child.Items = append(child.Items, childProduct)
			child.Venue.ID = child.Venue.ID + "_" + productRestaurantID
			child.ID = order.ID + "_" + productRestaurantID

			childRestaurantOrders[productRestaurantID] = child
			continue
		}
		childOrder.Items = append(childOrder.Items, childProduct)
		childRestaurantOrders[productRestaurantID] = childOrder
	}

	res := make([]interface{}, 0, len(childRestaurantOrders))
	for _, val := range childRestaurantOrders {
		res = append(res, s.calculateChildOrderTotalSum(val))
	}

	return res, nil
}

func (s *woltService) calculateChildOrderTotalSum(order models2.Order) models2.Order {
	totalSum := 0

	for _, p := range order.Items {
		totalSum = totalSum + p.TotalPrice.Amount
	}

	order.Price.Amount = totalSum
	return order
}

func (s *woltService) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	order, ok := req.(models2.Order)
	if !ok {
		return "", errors.New("casting error")
	}

	return order.Venue.ID, nil
}

func (s *woltService) GetSystemCreateOrderRequestByAggregatorRequest(r interface{}, store storeModels.Store) (models.Order, error) {
	virtualClildOrder, ok := r.(models2.Order)
	if ok {
		var (
			req models.Order
			err error
		)
		req = s.toOrderRequest(virtualClildOrder, store)

		req.LogMessages.FromDelivery, err = utils.GetJsonFormatFromModel(req)
		if err != nil {
			return req, err
		}

		utils.Beautify("creating order request body to orders-core", req)

		if virtualClildOrder.OrderStatus == models2.Rejected.String() {
			req = rejectedOrder(req, virtualClildOrder.OrderStatus)
		}

		return req, nil
	}

	webhook, ok := r.(models2.OrderNotification)
	if !ok {
		return models.Order{}, errors.New("casting error")
	}

	if webhook.Body.Status != models2.Created.String() {
		log.Info().Msgf("successfully created order, id=%s", webhook.Body.Id)
		return models.Order{}, errors.New(fmt.Sprintf("webhook action is not '%s'", models2.Created))
	}

	var (
		req models.Order
	)

	//if store.RestaurantGroupID == "670f9fbc7e9927e85ba1ca07" {
	//	order, err := s.cli.GetOrderByV2(context.Background(), webhook.Body.Id)
	//	if err != nil {
	//		return models.Order{}, err
	//	}
	//
	//	utils.Beautify("getting order from wolt v2 body", order)
	//
	//	req = s.toOrderV2Request(order, store)
	//	req.LogMessages.FromDelivery, err = utils.GetJsonFormatFromModel(order)
	//	if err != nil {
	//		return req, err
	//	}
	//
	//	if order.OrderStatus == models2.Rejected.String() {
	//		req = rejectedOrder(req, order.OrderStatus)
	//	}
	//
	//} else {
	order, err := s.cli.GetOrder(context.Background(), webhook.Body.Id)
	if err != nil {
		return models.Order{}, err
	}

	utils.Beautify("getting order from wolt body", order)

	req = s.toOrderRequest(order, store)

	req.LogMessages.FromDelivery, err = utils.GetJsonFormatFromModel(order)
	if err != nil {
		return req, err
	}

	if order.OrderStatus == models2.Rejected.String() {
		req = rejectedOrder(req, order.OrderStatus)
	}
	//}

	utils.Beautify("creating order request body to orders-core", req)

	return req, nil
}

func rejectedOrder(order models.Order, woltStatus string) models.Order {
	log.Info().Msgf("order canceled by delivery service, order id: %s, status: %s", order.OrderID, woltStatus)
	order.Status = models.STATUS_CANCELLED_BY_DELIVERY_SERVICE.ToString()
	order.StatusesHistory = append(order.StatusesHistory, models.OrderStatusUpdate{
		Name: models.STATUS_CANCELLED_BY_DELIVERY_SERVICE.ToString(),
		Time: time.Now().UTC(),
	})
	return order
}

func (s *woltService) MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	log.Info().Msg("delivery service: wolt")

	var statusConfig []storeModels.Status

	switch order.Type {
	case models.Instant:
		log.Info().Msgf("purchase type: INSTANT")
		statusConfig = store.Wolt.PurchaseTypes.Instant
	case models.Preorder:
		log.Info().Msgf("purchase type: PREORDER")
		statusConfig = store.Wolt.PurchaseTypes.Preorder
	}

	if order.IsPickedUpByCustomer {
		log.Info().Msgf("purchase type: TAKEAWAY")
		statusConfig = store.Wolt.PurchaseTypes.TakeAway
	}

	for _, status := range statusConfig {
		if status.PosStatus == posStatus.String() {
			log.Info().Msgf("[SPECIAL MATCHING] pos status: %v, wolt status: %v", status.PosStatus, status.Status)
			return status.Status
		}
	}

	log.Info().Msgf("[DEFAULT MATCHING], pos status: %v", posStatus)

	switch posStatus {
	case models.ACCEPTED, models.WAIT_COOKING, models.READY_FOR_COOKING, models.COOKING_STARTED, models.WAIT_SENDING:

		if order.Type == models.Preorder {
			return models.Confirm.String()
		}

		return models.Accept.String()
	case models.COOKING_COMPLETE, models.CLOSED, models.READY_FOR_PICKUP, models.ON_WAY, models.DELIVERED, models.OUT_FOR_DELIVERY:
		return models.Ready.String()
	case models.CANCELLED_BY_POS_SYSTEM:
		return models.Reject.String()
	case models.PICKED_UP_BY_CUSTOMER:
		return models.Delivered.String()
	}
	return ""
}

func (s *woltService) UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error {
	switch aggregatorStatus {

	case models.Accept.String():
		var adjustedPickUpTime *time.Time = nil
		if !store.Wolt.IgnorePickupTime {
			adjustedPickUpTime = &order.EstimatedPickupTime.Value.Time
		}

		log.Info().Msgf("result estimated pick up time: %v", adjustedPickUpTime)

		switch order.RestaurantSelfDelivery {
		case true:
			deliveryTime := time.Now().UTC().Add(time.Minute * time.Duration(store.Wolt.CookingTime+store.Wolt.AdjustedPickupMinutes))
			if err := s.cli.AcceptSelfDeliveryOrder(ctx, woltModels.AcceptSelfDeliveryOrderOrderRequest{
				ID:           order.OrderID,
				DeliveryTime: &deliveryTime,
			}); err != nil {
				log.Trace().Err(err).Msgf("accept selfDelivery order, order_id=%v, deliveryTime=%v", order.OrderID, deliveryTime)
				return err
			}
			log.Info().Msgf("success accept selfDelivery order, order_id=%v, deliveryTime=%v", order.OrderID, deliveryTime)

		case false:
			if err := s.cli.AcceptOrder(ctx, woltModels.AcceptOrderRequest{
				ID:         order.OrderID,
				PickupTime: adjustedPickUpTime,
			}); err != nil {
				log.Trace().Err(err).Msgf("accept order, order_id=%v, pick_up_time=%v", order.OrderID, adjustedPickUpTime)
				return err
			}
			log.Info().Msgf("success accept order, order_id=%v, pick_up_time=%v", order.OrderID, adjustedPickUpTime)
		}

	case models.Reject.String():

		orderID := order.OrderID
		reason := "reject"
		if err := s.cli.RejectOrder(ctx, woltModels.RejectOrderRequest{
			ID:     orderID,
			Reason: reason,
		}); err != nil {
			log.Trace().Err(err).Msgf("reject order, order_id=%v, reason=%v", orderID, reason)
			return err
		}

		log.Info().Msgf("success reject order, order_id=%v, reason=%v", orderID, reason)

	case models.Ready.String():

		orderID := order.OrderID
		if err := s.cli.MarkOrder(ctx, orderID); err != nil {
			log.Trace().Err(err).Msgf("mark order, order_id=%v", orderID)
			return err
		}

		log.Info().Msgf("success mark order, order_id=%v", orderID)

	case models.Confirm.String():

		orderID := order.OrderID
		if err := s.cli.ConfirmPreOrder(ctx, orderID); err != nil {
			log.Trace().Err(err).Msgf("confirm pre-order, order_id=%v", orderID)
			return err
		}

		log.Info().Msgf("success confirm pre-order, order_id=%v", orderID)

	case models.Delivered.String():

		orderID := order.OrderID
		if err := s.cli.DeliveredOrder(ctx, orderID); err != nil {
			log.Trace().Err(err).Msgf("delivered order, order_id=%v", orderID)
			return err
		}
	}
	log.Info().Msgf("aggregator status is empty: %v", aggregatorStatus)
	return nil
}

func (s *woltService) toOrderRequest(req models2.Order, store storeModels.Store) models.Order {

	cookingTime := s.cookingTimeFromProducts(req, store)

	var pickUpTime = req.CreatedAt.UTC().Add(time.Duration(cookingTime) * time.Minute)

	if store.Wolt.BusyMode {
		pickUpTime = pickUpTime.UTC().Add(time.Duration(store.Wolt.AdjustedPickupMinutes) * time.Minute)

		if pickUpTime.Sub(req.PickupEta.UTC()) > time.Duration(25)*time.Minute {
			pickUpTime = req.PickupEta.UTC().Add(time.Duration(25) * time.Minute)
		}
	} else if store.Wolt.ScheduledBusyMode && len(store.Wolt.ScheduledBusyModeTime) > 0 {
		for _, scheduledTime := range store.Wolt.ScheduledBusyModeTime {
			if s.isNowWithinRange(scheduledTime.From, scheduledTime.To) {
				pickUpTime = pickUpTime.UTC().Add(time.Duration(store.Wolt.AdjustedPickupMinutes) * time.Minute)
				if pickUpTime.Sub(req.PickupEta.UTC()) > time.Duration(25)*time.Minute {
					pickUpTime = req.PickupEta.UTC().Add(time.Duration(25) * time.Minute)
				}
				log.Info().Msgf("wolt %+v: for order_id: %s and pickUpTime: %v", scheduledTime, req.ID, pickUpTime)
			}
		}
	}

	if pickUpTime.Before(req.PickupEta.UTC()) {
		pickUpTime = req.PickupEta.UTC()
	}

	res := models.Order{
		OrderID:         req.ID,
		StoreID:         req.Venue.ID,
		OrderCode:       req.OrderNumber,
		PickUpCode:      req.OrderNumber,
		RestaurantID:    store.ID,
		PosType:         store.PosType,
		Type:            strings.ToUpper(req.Type),
		DeliveryService: models2.WOLT.String(),
		Status:          models2.STATUS_NEW.String(),
		StatusesHistory: []models.OrderStatusUpdate{
			{
				Name: models2.STATUS_NEW.String(),
				Time: req.CreatedAt.UTC(),
			},
		},
		OrderTime: models.TransactionTime{
			Value:    models.Time{Time: req.CreatedAt.UTC()},
			TimeZone: store.Settings.TimeZone.TZ,
		},
		EstimatedPickupTime: models.TransactionTime{
			Value:    models.Time{Time: pickUpTime},
			TimeZone: store.Settings.TimeZone.TZ,
		},
		Preorder: models.PreOrder{
			Time: models.TransactionTime{
				Value:    models.Time{Time: req.PreOrder.Time},
				TimeZone: store.Settings.TimeZone.TZ,
			},
			Status: req.PreOrder.Status,
		},
		UtcOffsetMinutes: strconv.Itoa(int(store.Settings.TimeZone.UTCOffset)),
		Currency:         req.Price.Currency,
		AllergyInfo:      req.ConsumerComment,
		EstimatedTotalPrice: models.Price{
			Value:        float64(req.Price.Amount) / 100,
			CurrencyCode: req.Price.Currency,
		},
		TotalCustomerToPay: TotalCustomerToPay(req),
		DeliveryFee: models.Price{
			Value:        float64(req.Delivery.Fee.Amount) / 100,
			CurrencyCode: req.Delivery.Fee.Currency,
		},
		Customer: models.Customer{
			Name:        req.ConsumerName,
			PhoneNumber: req.ConsumerPhoneNumber,
		},
		DeliveryAddress: models.DeliveryAddress{
			Label:     req.Delivery.Location.FormattedAddress,
			Latitude:  req.Delivery.Location.Coordinates.Latitude,
			Longitude: req.Delivery.Location.Coordinates.Longitude,
			City:      req.Delivery.Location.City,
			Street:    req.Delivery.Location.StreetAddress,
		},
		PaymentMethod:          "DELAYED",
		RestaurantSelfDelivery: req.Delivery.SelfDelivery,
		IsCashPayment:          isCashPayment(req),
	}

	if req.Delivery.Type == "takeaway" && req.Type == "preorder" {
		res.EstimatedPickupTime = models.TransactionTime{
			Value:    models.Time{Time: req.PreOrder.Time},
			TimeZone: store.Settings.TimeZone.TZ,
		}
		res.IsPickedUpByCustomer = true
	}

	res.Products = make([]models.OrderProduct, 0, len(req.Items))
	for _, product := range req.Items {
		res.Products = append(res.Products, s.toOrderProduct(product))
	}

	return res
}

// игнорируем дату
func (s woltService) isNowWithinRange(from time.Time, to time.Time) bool {
	now := time.Now()

	normalize := func(t time.Time) time.Time {
		return time.Date(1970, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	}

	fromNormalized := normalize(from)
	toNormalized := normalize(to)
	nowNormalized := normalize(now)

	return nowNormalized.After(fromNormalized) && nowNormalized.Before(toNormalized)
}

func (s *woltService) cookingTimeFromProducts(req models2.Order, store storeModels.Store) int32 {
	var productIds []string
	for _, product := range req.Items {
		productIds = append(productIds, product.PosID)
	}

	menu, err := s.menuService.GetAggregatorMenuIfExists(context.Background(), store, storeModels.WOLT.String())
	if err != nil {
		return int32(store.Wolt.CookingTime)
	}

	cookingTime, err := s.menuService.GetLongestCookingTimeByProductIds(context.Background(), menu.ID, productIds)

	if err != nil || cookingTime == 0 {
		return int32(store.Wolt.CookingTime)
	}

	return cookingTime
}

func (s *woltService) cookingTimeFromProducts2(req models_v2.Order, store storeModels.Store) int32 {
	var productIds []string
	for _, product := range req.Items {
		productIds = append(productIds, product.PosId)
	}

	menu, err := s.menuService.GetAggregatorMenuIfExists(context.Background(), store, storeModels.WOLT.String())
	if err != nil {
		return int32(store.Wolt.CookingTime)
	}

	cookingTime, err := s.menuService.GetLongestCookingTimeByProductIds(context.Background(), menu.ID, productIds)

	if err != nil || cookingTime == 0 {
		return int32(store.Wolt.CookingTime)
	}

	return cookingTime
}

func (s *woltService) toOrderProduct(req models2.OrderItem) models.OrderProduct {
	res := models.OrderProduct{
		ID:   req.PosID,
		Name: req.Name,
		Price: models.Price{
			Value: float64(req.BasePrice.Amount) / 100,
		},
		Quantity: req.Count,
	}

	res.Attributes = make([]models.ProductAttribute, 0, len(req.Options))
	for _, attribute := range req.Options {
		res.Attributes = append(res.Attributes, models.ProductAttribute{
			ID:      attribute.ValuePosID,
			GroupID: attribute.PosID,
			Name:    attribute.Value,
			Price: models.Price{
				Value: float64(attribute.Price.Amount) / 100,
			},
			Quantity: attribute.Count,
		})
	}

	return res
}

func (s *woltService) toOrderV2Request(req models_v2.Order, store storeModels.Store) models.Order {
	createdAt, err := time.Parse(time.RFC3339, req.CreatedAt)
	if err != nil {
		log.Err(err).Msgf("created at time parse error")
	}

	pickUpEta, err := time.Parse(time.RFC3339, req.PickupEta)
	if err != nil {
		log.Err(err).Msgf("pickup eta time parse error")
	}

	preOrderTime, err := time.Parse(time.RFC3339, req.PreOrder.Time)
	if err != nil {
		log.Err(err).Msgf("preorder time parse error")
	}

	cookingTime := s.cookingTimeFromProducts2(req, store)

	var pickUpTime = createdAt.UTC().Add(time.Duration(cookingTime) * time.Minute)

	if store.Wolt.BusyMode {
		pickUpTime = pickUpTime.UTC().Add(time.Duration(store.Wolt.AdjustedPickupMinutes) * time.Minute)

		if pickUpTime.Sub(pickUpEta.UTC()) > time.Duration(25)*time.Minute {
			pickUpTime = pickUpEta.UTC().Add(time.Duration(25) * time.Minute)
		}
	} else if store.Wolt.ScheduledBusyMode && len(store.Wolt.ScheduledBusyModeTime) > 0 {
		for _, scheduledTime := range store.Wolt.ScheduledBusyModeTime {
			if s.isNowWithinRange(scheduledTime.From, scheduledTime.To) {
				pickUpTime = pickUpTime.UTC().Add(time.Duration(store.Wolt.AdjustedPickupMinutes) * time.Minute)
				if pickUpTime.Sub(pickUpEta.UTC()) > time.Duration(25)*time.Minute {
					pickUpTime = pickUpEta.UTC().Add(time.Duration(25) * time.Minute)
				}
				log.Info().Msgf("wolt %+v: for order_id: %s and pickUpTime: %v", scheduledTime, req.Id, pickUpTime)
			}
		}
	}

	if pickUpTime.Before(pickUpEta.UTC()) {
		pickUpTime = pickUpEta.UTC()
	}

	order := models.Order{
		OrderID:         req.Id,
		StoreID:         req.Venue.Id,
		OrderCode:       req.OrderNumber,
		PickUpCode:      req.OrderNumber,
		RestaurantID:    store.ID,
		PosType:         store.PosType,
		Type:            strings.ToUpper(req.Type),
		DeliveryService: models2.WOLT.String(),
		Status:          models2.STATUS_NEW.String(),
		StatusesHistory: []models.OrderStatusUpdate{
			{
				Name: models2.STATUS_NEW.String(),
				Time: createdAt.UTC(),
			},
		},
		Preorder: models.PreOrder{
			Time: models.TransactionTime{
				Value:    models.Time{Time: preOrderTime.UTC()},
				TimeZone: store.Settings.TimeZone.TZ,
			},
			Status: req.PreOrder.Status,
		},
		OrderTime: models.TransactionTime{
			Value:    models.Time{Time: createdAt.UTC()},
			TimeZone: store.Settings.TimeZone.TZ,
		},
		EstimatedPickupTime: models.TransactionTime{
			Value:    models.Time{Time: pickUpTime},
			TimeZone: store.Settings.TimeZone.TZ,
		},
		UtcOffsetMinutes: strconv.Itoa(int(store.Settings.TimeZone.UTCOffset)),
		Currency:         req.BasketPrice.Total.Currency,
		AllergyInfo:      req.ConsumerComment,
		EstimatedTotalPrice: models.Price{
			Value:        float64(req.BasketPrice.PriceBreakdown.TotalBeforeDiscounts.Amount) / 100,
			CurrencyCode: req.BasketPrice.Total.Currency,
		},
		TotalCustomerToPay: models.Price{
			Value:        float64(req.BasketPrice.Total.Amount) / 100,
			CurrencyCode: req.BasketPrice.Total.Currency,
		},
		DeliveryFee: models.Price{
			Value:        float64(req.Fees.PriceBreakdown.TotalDiscounts.Amount) / 100,
			CurrencyCode: req.Fees.PriceBreakdown.TotalDiscounts.Currency,
		},
		PartnerDiscountsProducts: models.Price{
			Value:        math.Abs(float64(req.BasketPrice.PriceBreakdown.TotalDiscounts.Amount)) / 100,
			CurrencyCode: req.BasketPrice.PriceBreakdown.TotalDiscounts.Currency,
		},
		Customer: models.Customer{
			Name:        req.ConsumerName,
			PhoneNumber: req.ConsumerPhoneNumber,
		},
		DeliveryAddress: models.DeliveryAddress{
			Label:     req.Delivery.Location.StreetAddress,
			Latitude:  req.Delivery.Location.Coordinates.Lat,
			Longitude: req.Delivery.Location.Coordinates.Lon,
		},
		PaymentMethod:          "DELAYED",
		RestaurantSelfDelivery: req.Delivery.SelfDelivery,
	}

	if req.Delivery.Type == "takeaway" {
		order.IsPickedUpByCustomer = true
	}

	if req.Type == "preorder" && !preOrderTime.IsZero() {
		order.EstimatedPickupTime = models.TransactionTime{
			Value:    models.Time{Time: preOrderTime},
			TimeZone: store.Settings.TimeZone.TZ,
		}
	}

	order.Products = make([]models.OrderProduct, 0, len(req.Items))
	for _, product := range req.Items {
		order.Products = append(order.Products, s.toOrderProductV2(product))
	}

	return order
}

func (s *woltService) toOrderProductV2(req models_v2.Item) models.OrderProduct {
	res := models.OrderProduct{
		ID:   req.PosId,
		Name: req.Name,
		Price: models.Price{
			Value: float64(req.ItemPrice.UnitPrice.Amount) / 100 / float64(req.Count),
		},
		Quantity: req.Count,
	}

	res.Attributes = make([]models.ProductAttribute, 0, len(req.Options))
	for _, attribute := range req.Options {
		res.Attributes = append(res.Attributes, models.ProductAttribute{
			ID:      attribute.ValuePosId,
			GroupID: attribute.PosId,
			Name:    attribute.Value,
			Price: models.Price{
				Value: float64(attribute.Price.Amount) / 100,
			},
			Quantity: attribute.Count,
		})
	}

	return res
}

func (s *woltService) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error) {
	if err := s.cli.UpdateStopListByProducts(ctx, aggregatorStoreID, products, isAvailable); err != nil {
		return "", err
	}

	return "", nil
}

func (s *woltService) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error) {
	if isSendRemains {
		updateProducts := s.toUpdateInventoryProducts(products)
		if len(updateProducts.Data) == 0 {
			return "", nil
		}
		if err := s.cli.UpdateMenuItemInventory(ctx, aggregatorStoreID, updateProducts); err != nil {
			log.Err(err).Msgf("UpdateInventoryProducts Wolt error")
			if strings.Contains(err.Error(), "429") {
				s.retryUpdateMenuItemInventory(ctx, updateProducts, aggregatorStoreID, 4)
			}
		}
		time.Sleep(10 * time.Second)
	}

	if err := s.cli.UpdateStopListByProductsBulk(ctx, aggregatorStoreID, products); err != nil {
		log.Err(err).Msgf("UpdateStopListByProducts Wolt error")
		if strings.Contains(err.Error(), "429") {
			s.retryUpdateStopListByProductsBulk(ctx, aggregatorStoreID, products, 4)
		}
		return "", err
	}

	return "", nil
}

func (s *woltService) GetAggregatorOrder(ctx context.Context, orderID string) (models2.Order, error) {
	log.Info().Msgf("start to get aggregator order of order_id: %s", orderID)
	order, err := s.cli.GetOrder(ctx, orderID)
	if err != nil {
		return models2.Order{}, err
	}
	return order, nil
}

func (s *woltService) toUpdateInventoryProducts(products []menuModels.Product) woltModels.WoltInventory {
	woltInventory := woltModels.WoltInventory{}

	for _, p := range products {
		if p.Balance <= 0 {
			p.Balance = 0
		}
		woltInventory.Data = append(woltInventory.Data, woltModels.Item{
			ExtID:     p.ExtID,
			Inventory: int(p.Balance)},
		)
	}

	return woltInventory
}

func (s *woltService) retryUpdateMenuItemInventory(ctx context.Context, updateProducts woltModels.WoltInventory, venueID string, retryCount int) {
	log.Info().Msgf("retrying UpdateMenuItemInventory Wolt; retryCount=%d", retryCount)
	for i := 0; i < retryCount; i++ {
		time.Sleep(3 * time.Second)
		if err := s.cli.UpdateMenuItemInventory(ctx, venueID, updateProducts); err != nil {
			log.Err(err).Msgf("UpdateInventoryProducts Wolt error; retryCount: %d", i+1)
		} else {
			log.Info().Msgf("UpdateInventoryProducts Wolt success; retryCount: %d", i+1)
			return
		}
	}
}

func (s *woltService) retryUpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, retryCount int) {
	log.Info().Msgf("retrying UpdateStopListByProducts Wolt error; retryCount=%d", len(products))
	for i := 0; i < retryCount; i++ {
		time.Sleep(3 * time.Second)
		if err := s.cli.UpdateStopListByProductsBulk(ctx, aggregatorStoreID, products); err != nil {
			log.Err(err).Msgf("UpdateStopListByProducts Wolt error; retryCount: %d", i+1)
		} else {
			log.Info().Msgf("UpdateStopListByProducts Wolt success; retryCount: %d", i+1)
			return
		}
	}
}

func TotalCustomerToPay(req models2.Order) models.Price {
	if req.CashPayment.CashAmount.Amount > 0 {
		return models.Price{
			Value:        float64(req.CashPayment.CashAmount.Amount) / 100,
			CurrencyCode: req.CashPayment.CashAmount.Currency,
		}

	}
	return models.Price{
		Value:        float64(req.Price.Amount) / 100,
		CurrencyCode: req.Price.Currency,
	}
}
func isCashPayment(req models2.Order) bool {
	return req.CashPayment.CashAmount.Amount > 0
}

func (s *woltService) SendOrderErrorNotification(ctx context.Context, req interface{}) error {
	return nil
}

func (s *woltService) SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error {
	return nil
}
