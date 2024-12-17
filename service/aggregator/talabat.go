package aggregator

import (
	"context"
	"github.com/google/uuid"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models2 "github.com/kwaaka-team/orders-core/core/talabat/models"
	models3 "github.com/kwaaka-team/orders-core/core/wolt/models"
	httpclient "github.com/kwaaka-team/orders-core/pkg/talabat"
	"github.com/kwaaka-team/orders-core/pkg/talabat/clients"
	talabatModels "github.com/kwaaka-team/orders-core/pkg/talabat/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type talabatService struct {
	middlewareCli clients.TalabatMW
	menuCli       clients.TalabatMenu
	restaurantID  string
}

func newTalabatService(middlewareBaseURL, menuBaseUrl string, store storeModels.Store) (*talabatService, error) {
	middlewareCli, err := httpclient.NewMiddlewareClient(&clients.Config{
		Protocol: "http",
		BaseURL:  middlewareBaseURL,
		Username: store.Talabat.Username,
		Password: store.Talabat.Password,
	})
	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Talabat middleware Manager")
		return nil, err
	}

	menuCli, err := httpclient.NewMenuClient(&clients.Config{
		Protocol: "http",
		BaseURL:  menuBaseUrl,
		Username: store.Talabat.Username,
		Password: store.Talabat.Password,
	})
	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Talabat menu manager")
		return nil, err
	}

	return &talabatService{
		middlewareCli, menuCli, store.Talabat.RestaurantID,
	}, nil
}

func (s *talabatService) OpenStore(ctx context.Context, aggregatorStoreId string) error {
	return errors.New("method not implemented")
}

func (s *talabatService) GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error) {
	return false, errors.New("method not implemented")
}

func (s *talabatService) GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error) {
	return storeModels.AggregatorSchedule{}, errors.New("method not implemented")
}

func (s *talabatService) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error) {
	requestID := uuid.New().String()
	scheduledOn, err := time.Now().Add(2 * time.Minute).UTC().MarshalText()
	if err != nil {
		return "", err
	}
	err = s.menuCli.UpdateItemsAvailability(ctx, talabatModels.UpdateItemsAvailabilityRequest{
		RequestID:    requestID,
		RestaurantID: s.restaurantID,
		ScheduledOn:  string(scheduledOn),
		Availability: s.toAvailabilities(aggregatorStoreID, products, nil, isAvailable),
	})
	if err != nil {
		return requestID, err
	}

	return requestID, nil
}

func (s *talabatService) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error) {
	requestID := uuid.New().String()
	scheduledOn, err := time.Now().Add(2 * time.Minute).UTC().MarshalText()
	if err != nil {
		return "", err
	}
	err = s.menuCli.UpdateItemsAvailability(ctx, talabatModels.UpdateItemsAvailabilityRequest{
		RequestID:    requestID,
		RestaurantID: s.restaurantID,
		ScheduledOn:  string(scheduledOn),
		Availability: s.toAvailabilitiesBulk(aggregatorStoreID, products, nil),
	})
	if err != nil {
		return requestID, err
	}

	return requestID, nil
}

func (s *talabatService) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error) {
	requestID := uuid.New().String()
	scheduledOn, err := time.Now().Add(2 * time.Minute).UTC().MarshalText()
	if err != nil {
		return "", err
	}
	err = s.menuCli.UpdateItemsAvailability(ctx, talabatModels.UpdateItemsAvailabilityRequest{
		RequestID:    requestID,
		RestaurantID: s.restaurantID,
		ScheduledOn:  string(scheduledOn),
		Availability: s.toAvailabilitiesBulk(aggregatorStoreID, nil, attributes),
	})
	if err != nil {
		return requestID, err
	}

	return requestID, nil
}

func (s *talabatService) IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error) {
	return store.Talabat.IsMarketplace, nil
}

func (s *talabatService) SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	return nil, nil
}

func (s *talabatService) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	return "", nil
}

func (s *talabatService) GetSystemCreateOrderRequestByAggregatorRequest(r interface{}, store storeModels.Store) (models.Order, error) {
	req, ok := r.(models2.CreateOrderRequest)
	if !ok {
		return models.Order{}, errors.New("casting error")
	}

	estimatedTotalPriceValue, err := strconv.ParseFloat(req.Price.GrandTotal, 64)
	if err != nil {
		return models.Order{}, err
	}

	var customerCashPaymentAmount float64
	if req.Price.PayRestaurant != "" {
		amount, err := strconv.ParseFloat(req.Price.PayRestaurant, 64)
		if err != nil {
			return models.Order{}, err
		}
		customerCashPaymentAmount = amount
	}

	orderProducts, err := s.toProducts(req.Products)
	if err != nil {
		return models.Order{}, err
	}

	order := models.Order{
		OrderID:             req.Token,
		OrderCode:           req.Code,
		DeliveryService:     "talabat",
		SpecialRequirements: req.Comments.CustomerComment,
		AllergyInfo:         req.Comments.CustomerComment,
		Customer: models.Customer{
			Name:        req.Customer.FirstName + " " + req.Customer.LastName,
			PhoneNumber: req.Customer.MobilePhone,
		},
		DeliveryAddress: models.DeliveryAddress{
			Longitude: req.Delivery.Address.Longitude,
			Latitude:  req.Delivery.Address.Latitude,
			Label:     s.toLabel(req.Delivery.Address),
		},
		EstimatedPickupTime: models.TransactionTime{
			Value:    models.Time{Time: s.toEstimatedPickUpTime(req.Delivery, store.Talabat.AdjustedPickupMinutes)},
			TimeZone: store.Settings.TimeZone.TZ,
		},
		PickUpCode:       req.ShortCode,
		UtcOffsetMinutes: strconv.Itoa(int(store.Settings.TimeZone.UTCOffset)),
		EstimatedTotalPrice: models.Price{
			Value:        estimatedTotalPriceValue,
			CurrencyCode: store.Settings.Currency,
		},
		DeliveryFee: models.Price{
			Value:        s.toDeliveryFeeValue(req.Price.DeliveryFees),
			CurrencyCode: store.Settings.Currency,
		},
		CustomerCashPaymentAmount: models.Price{
			Value:        customerCashPaymentAmount,
			CurrencyCode: store.Settings.Currency,
		},
		Products: orderProducts,
		OrderTime: models.TransactionTime{
			Value:    models.Time{Time: req.CreatedAt.UTC()},
			TimeZone: store.Settings.TimeZone.TZ,
		},
		StatusesHistory: []models.OrderStatusUpdate{
			{
				Name: "NEW",
				Time: time.Now().UTC(),
			},
		},
	}

	if req.ExpeditionType == "pickup" {
		order.IsPickedUpByCustomer = true
		order.PickUpCode = req.Pickup.PickupCode
		order.EstimatedPickupTime = models.TransactionTime{
			Value:    models.Time{Time: s.toEstimatedPickUpTimeCustomerPickUp(req.Pickup.PickupTime, store.Talabat.AdjustedPickupMinutes)},
			TimeZone: store.Settings.TimeZone.TZ,
		}
	}

	switch req.Payment.Status {
	case "paid":
		order.PaymentMethod = "DELAYED"
	case "pending":
		order.PaymentMethod = "CASH"
	}

	switch req.PreOrder {
	case true:
		order.Type = "PREORDER"
		order.Preorder = models.PreOrder{
			Time: models.TransactionTime{
				Value:    order.EstimatedPickupTime.Value,
				TimeZone: store.Settings.TimeZone.TZ,
			},
		}
	case false:
		order.Type = "INSTANT"
	}

	return order, nil
}

func (s *talabatService) MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	log.Info().Msg("delivery service: chocofood")

	switch posStatus {
	case models.ACCEPTED, models.WAIT_COOKING, models.READY_FOR_COOKING, models.COOKING_STARTED, models.WAIT_SENDING:
		return models.OrderAccepted.String()
	case models.COOKING_COMPLETE, models.CLOSED, models.READY_FOR_PICKUP, models.ON_WAY, models.DELIVERED, models.OUT_FOR_DELIVERY:
		return models.OrderPrepared.String()
	case models.CANCELLED_BY_POS_SYSTEM:
		return models.OrderRejected.String()
	case models.PICKED_UP_BY_CUSTOMER:
		return models.OrderPickedUp.String()
	default:
		return ""
	}
}

func (s *talabatService) UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error {
	switch aggregatorStatus {
	case models.OrderAccepted.String():
		if err := s.middlewareCli.AcceptOrder(ctx, talabatModels.AcceptOrderRequest{
			OrderToken:     order.OrderID,
			RemoteOrderId:  order.OrderID,
			Status:         models.OrderAccepted.String(),
			AcceptanceTime: time.Now().UTC().Format(time.RFC3339),
		}); err != nil {
			log.Trace().Err(err).Msgf("error accept talabat order, order_token=%v", order.OrderID)
			return err
		}
		log.Info().Msgf("success accept talabat order, order_token=%v", order.OrderID)

	case models.OrderRejected.String():
		if err := s.middlewareCli.RejectOrder(ctx, talabatModels.RejectOrderRequest{
			OrderToken: order.OrderID,
			Reason:     "TECHNICAL_PROBLEM",
			Status:     models.OrderRejected.String(),
			Message:    "TECHNICAL_PROBLEM",
		}); err != nil {
			log.Trace().Err(err).Msgf("error reject talabat order, order_token=%v", order.OrderID)
			return err
		}
		log.Info().Msgf("success reject talabat order, order_token=%v", order.OrderID)

	case models.OrderPickedUp.String():
		if err := s.middlewareCli.OrderPickedUp(ctx, talabatModels.OrderPickedUpRequest{
			OrderToken: order.OrderID,
			Status:     models.OrderPickedUp.String(),
		}); err != nil {
			log.Trace().Err(err).Msgf("error reject talabat order, order_token=%v", order.OrderID)
			return err
		}
		log.Info().Msgf("success reject talabat order, order_token=%v", order.OrderID)

	case models.OrderPrepared.String():
		if err := s.middlewareCli.MarkOrderPrepared(ctx, order.OrderID); err != nil {
			log.Trace().Err(err).Msgf("error mark order prepared talabat order, order_token=%v", order.OrderID)
			return err
		}
		log.Info().Msgf("success mark order prepared talabat order, order_token=%v", order.OrderID)
	}

	log.Info().Msgf("talabat aggregator status is invalid: %v", aggregatorStatus)
	return nil
}

func (s *talabatService) toProducts(req []models2.Product) ([]models.OrderProduct, error) {
	res := make([]models.OrderProduct, 0, len(req))
	for _, p := range req {
		product, err := s.toProduct(p)
		if err != nil {
			return nil, err
		}
		res = append(res, product)
	}
	return res, nil
}

func (s *talabatService) toProduct(req models2.Product) (models.OrderProduct, error) {
	productUnitPrice, err := strconv.ParseFloat(req.UnitPrice, 64)
	if err != nil {
		return models.OrderProduct{}, err
	}

	productQuantity, err := strconv.Atoi(req.Quantity)
	if err != nil {
		return models.OrderProduct{}, err
	}

	res := models.OrderProduct{
		ID:   req.RemoteCode,
		Name: req.Name,
		Price: models.Price{
			Value: productUnitPrice,
		},
		Quantity: productQuantity,
	}
	res.Attributes = make([]models.ProductAttribute, 0, len(req.SelectedToppings))
	for _, t := range req.SelectedToppings {
		price, err := strconv.ParseFloat(t.Price, 64)
		if err != nil {
			return models.OrderProduct{}, err
		}
		res.Attributes = append(res.Attributes, models.ProductAttribute{
			ID:       t.RemoteCode,
			Quantity: t.Quantity,
			Price: models.Price{
				Value: price,
			},
			Name: t.Name,
		})
	}
	return res, nil
}

func (s *talabatService) toLabel(address models2.Address) string {
	res := ""
	if address.City != "" {
		res = res + address.City + ","
	}

	if address.Street != "" {
		res = res + address.Street + ","
	}

	if address.Number != "" {
		res = res + address.Number + ","
	}
	if address.DeliveryArea != "" {
		res = res + address.DeliveryArea + ","
	}

	if address.Entrance != "" {
		res = res + address.Entrance + ","
	}

	if address.Intercom != "" {
		res = res + "intercom: " + address.Intercom + ","
	}

	if address.FlatNumber != "" {
		res = res + address.FlatNumber + ","
	}

	if address.Building != "" {
		res = res + address.Building + ","
	}

	if address.DeliveryInstructions != "" {
		res = res + address.DeliveryInstructions + ","
	}

	if len(res) > 0 && res[len(res)-1:] == "," {
		res = res[:len(res)-1]
	}

	return res
}

func (s *talabatService) toEstimatedPickUpTime(delivery models2.Delivery, adjuctedPickUpMinutes int) time.Time {
	if delivery.ExpectedDeliveryTime != "" {
		estimatedPickUpTime, err := time.Parse(time.RFC3339, delivery.ExpectedDeliveryTime)
		if err != nil {
			estimatedPickUpTime = time.Now().UTC().Add(time.Duration(adjuctedPickUpMinutes) * time.Minute)
		}
		return estimatedPickUpTime
	}

	if delivery.RiderPickupTime != "" {
		estimatedPickUpTime, err := time.Parse(time.RFC3339, delivery.RiderPickupTime)
		if err != nil {
			estimatedPickUpTime = time.Now().Add(time.Duration(adjuctedPickUpMinutes) * time.Minute)
		}
		return estimatedPickUpTime
	}
	return time.Now().Add(time.Duration(adjuctedPickUpMinutes) * time.Minute)
}

func (s *talabatService) toDeliveryFeeValue(req []models2.DeliveryFees) float64 {
	var res float64
	for _, f := range req {
		res = res + f.Value
	}
	return res
}

func (s *talabatService) toEstimatedPickUpTimeCustomerPickUp(reqTime string, adjuctedPickUpMinutes int) time.Time {
	if reqTime != "" {
		estimatedPickUpTime, err := time.Parse(time.RFC3339, reqTime)
		if err != nil {
			estimatedPickUpTime = time.Now().UTC().Add(time.Duration(adjuctedPickUpMinutes) * time.Minute)
		}
		return estimatedPickUpTime
	}
	return time.Now().UTC().Add(time.Duration(adjuctedPickUpMinutes) * time.Minute)
}

func (s *talabatService) toAvailabilities(storeID string, products menuModels.Products, attributes menuModels.Attributes, isAvailable bool) []talabatModels.Availability {
	items := make([]talabatModels.ItemStoplist, 0, len(products)+len(attributes))
	for _, product := range products {
		items = append(items, s.toItem(product.ExtID, isAvailable))
	}
	for _, attribute := range attributes {
		items = append(items, s.toItem(attribute.ExtID, isAvailable))
	}

	return []talabatModels.Availability{
		{
			BranchId: storeID,
			Items:    items,
		},
	}
}

func (s *talabatService) toAvailabilitiesBulk(storeID string, products menuModels.Products, attributes menuModels.Attributes) []talabatModels.Availability {
	items := make([]talabatModels.ItemStoplist, 0, len(products)+len(attributes))
	for _, product := range products {
		items = append(items, s.toItem(product.ExtID, product.IsAvailable))
	}
	for _, attribute := range attributes {
		items = append(items, s.toItem(attribute.ExtID, attribute.IsAvailable))
	}

	return []talabatModels.Availability{
		{
			BranchId: storeID,
			Items:    items,
		},
	}
}

func (s *talabatService) toItem(itemID string, isAvailable bool) talabatModels.ItemStoplist {
	return talabatModels.ItemStoplist{
		ItemId: itemID,
		Status: s.toStatus(isAvailable),
	}
}

func (s *talabatService) toStatus(isAvailable bool) int {
	if isAvailable {
		return 0
	}
	return 1
}

func (s *talabatService) GetAggregatorOrder(ctx context.Context, orderID string) (models3.Order, error) {
	return models3.Order{}, nil
}

func (s *talabatService) SendOrderErrorNotification(ctx context.Context, req interface{}) error {
	return nil
}

func (s *talabatService) SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error {
	return nil
}
