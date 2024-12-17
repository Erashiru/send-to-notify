package aggregator

import (
	"context"
	models2 "github.com/kwaaka-team/orders-core/core/glovo/models"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models3 "github.com/kwaaka-team/orders-core/core/wolt/models"
	glovoClient "github.com/kwaaka-team/orders-core/pkg/glovo"
	"github.com/kwaaka-team/orders-core/pkg/glovo/clients"
	glovoModels "github.com/kwaaka-team/orders-core/pkg/glovo/clients/dto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type glovoService struct {
	cli clients.Glovo
}

func newGlovoService(baseURL, token string) (*glovoService, error) {
	cli, err := glovoClient.NewGlovoClient(&clients.Config{
		Protocol: "http",
		BaseURL:  baseURL,
		ApiKey:   token,
	})
	if err != nil {
		return nil, err
	}

	return &glovoService{cli}, nil
}

func (s *glovoService) OpenStore(ctx context.Context, aggregatorStoreId string) error {
	if err := s.cli.OpenStore(ctx, glovoModels.StoreManageRequest{
		StoreID: aggregatorStoreId,
	}); err != nil {
		return err
	}

	return nil
}

func (s *glovoService) GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error) {
	status, err := s.cli.GetStoreStatus(ctx, aggregatorStoreId)
	if err != nil {
		return false, err
	}

	return status.Until == "", nil
}

func (s *glovoService) toSystemSchedule(ctx context.Context, schedule glovoModels.StoreScheduleResponse) storeModels.AggregatorSchedule {
	result := storeModels.AggregatorSchedule{
		Timezone: schedule.Timezone,
		Schedule: make([]storeModels.Schedule, 0, len(schedule.Schedule)),
	}

	for _, cur := range schedule.Schedule {
		timeSlots := make([]storeModels.TimeSlot, 0, len(cur.TimeSlots))

		for _, timeSlot := range cur.TimeSlots {
			timeSlots = append(timeSlots, storeModels.TimeSlot{
				Opening: timeSlot.Opening,
				Closing: timeSlot.Closing,
			})
		}

		result.Schedule = append(result.Schedule, storeModels.Schedule{
			DayOfWeek: cur.DayOfWeek,
			TimeSlots: timeSlots,
		})
	}

	return result
}

func (s *glovoService) GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error) {
	storeSchedule, err := s.cli.GetStoreSchedule(ctx, aggregatorStoreId)
	if err != nil {
		return storeModels.AggregatorSchedule{}, err
	}

	return s.toSystemSchedule(ctx, storeSchedule), nil
}

func (s *glovoService) IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error) {
	return restaurantSelfDelivery != true, nil
}

func (s *glovoService) getRestaurantSelfDelivery(deliveryFee *int) bool {
	if deliveryFee != nil {
		return true
	}

	return false
}

func (s *glovoService) getDeliveryFee(deliveryFee *int) float64 {
	if deliveryFee != nil {
		return float64(*deliveryFee) / 100
	}

	return 0
}

func (s *glovoService) GetAggregatorCreateOrderRequestBySystemRequest(req models.Order) (interface{}, error) {
	layout := "2006-01-02 15:04:05"
	offsetMinutes := req.OrderTime.UTCOffset
	if offsetMinutes == "" {
		offsetMinutes = "300"
	}
	offsetDuration, err := time.ParseDuration(offsetMinutes + "m")
	if err != nil {
		offsetDuration = 300 * time.Minute
	}
	res := models2.Order{
		OrderID:             req.OrderID,
		StoreID:             req.StoreID,
		OrderTime:           req.OrderTime.Value.Add(offsetDuration).Format(layout),
		EstimatedPickupTime: req.EstimatedPickupTime.Value.Add(offsetDuration).Format(layout),
		PaymentMethod:       req.PaymentMethod,
		Currency:            req.Currency,
		Courier: models2.Courier{
			Name:        req.Courier.Name,
			PhoneNumber: req.Courier.PhoneNumber,
		},
		Customer: models2.Customer{
			Name:        req.Customer.Name,
			PhoneNumber: req.Customer.PhoneNumber,
			Hash:        req.Customer.Hash,
			InvoicingDetails: models2.InvoicingDetails{
				CompanyName:    req.Customer.InvoicingDetails.CompanyName,
				CompanyAddress: req.Customer.InvoicingDetails.CompanyAddress,
				TaxID:          req.Customer.InvoicingDetails.TaxID,
			},
		},
		OrderCode:            req.OrderCode,
		AllergyInfo:          req.AllergyInfo,
		UtcOffsetMinutes:     req.UtcOffsetMinutes,
		PickUpCode:           req.PickUpCode,
		IsPickedUpByCustomer: req.IsPickedUpByCustomer,
		CutleryRequested:     req.CutleryRequested,
	}

	return res, nil
}

func (s *glovoService) GetSystemCreateOrderRequestByAggregatorRequest(r interface{}, store storeModels.Store) (models.Order, error) {
	req, ok := r.(models2.Order)
	if !ok {
		return models.Order{}, errors.New("casting error")
	}

	layout := "2006-01-02 15:04:05" // The layout for the string representation of time

	orderTime, err := s.parseTimeWithOffset(req.OrderTime, layout, req.UtcOffsetMinutes)
	if err != nil {
		log.Err(err).Msgf("order time parse error")
		orderTime = time.Now()
	}
	estimatedPickUpTime, err := s.parseTimeWithOffset(req.EstimatedPickupTime, layout, req.UtcOffsetMinutes)
	if err != nil {
		log.Err(err).Msgf("estimated parse error")
		estimatedPickUpTime = orderTime.Add(30 * time.Minute)
	}
	tz := req.UtcOffsetMinutes

	res := models.Order{
		OrderID:         req.OrderID,
		StoreID:         req.StoreID,
		OrderCode:       req.OrderCode,
		Type:            models2.INSTANT.String(),
		DeliveryService: models2.GLOVO.String(),
		PickUpCode:      req.PickUpCode,
		Status:          models2.STATUS_NEW.String(),
		StatusesHistory: []models.OrderStatusUpdate{
			{
				Name: models2.STATUS_NEW.String(),
				Time: time.Now().UTC(),
			},
		},
		OrderTime: models.TransactionTime{
			Value:    models.Time{Time: orderTime.UTC()},
			TimeZone: tz,
		},
		EstimatedPickupTime: models.TransactionTime{
			Value:    models.Time{Time: estimatedPickUpTime.UTC().Add(20 * time.Minute)},
			TimeZone: tz,
		},
		UtcOffsetMinutes:    req.UtcOffsetMinutes,
		PaymentMethod:       req.PaymentMethod,
		Currency:            req.Currency,
		AllergyInfo:         req.AllergyInfo,
		SpecialRequirements: req.SpecialRequirements,
		EstimatedTotalPrice: models.Price{
			Value:        float64(req.EstimatedTotalPrice) / 100,
			CurrencyCode: req.Currency,
		},
		DeliveryFee: models.Price{
			Value:        s.getDeliveryFee(req.DeliveryFee),
			CurrencyCode: req.Currency,
		},
		MinimumBasketSurcharge: models.Price{
			Value:        float64(req.MinimumBasketSurcharge) / 100,
			CurrencyCode: req.Currency,
		},
		CustomerCashPaymentAmount: models.Price{
			Value:        float64(req.CustomerCashPaymentAmount) / 100,
			CurrencyCode: req.Currency,
		},
		PartnerDiscountsProducts: models.Price{
			Value:        float64(req.PartnerDiscountsProducts) / 100,
			CurrencyCode: req.Currency,
		},
		PartnerDiscountedProductsTotal: models.Price{
			Value:        float64(req.PartnerDiscountedProductsTotal) / 100,
			CurrencyCode: req.Currency,
		},
		TotalCustomerToPay: models.Price{
			Value:        float64(req.TotalCustomerToPay) / 100,
			CurrencyCode: req.Currency,
		},
		Courier: models.Courier{
			Name:        req.Courier.Name,
			PhoneNumber: req.Courier.PhoneNumber,
		},
		Customer: models.Customer{
			Name:        req.Customer.Name,
			PhoneNumber: req.Customer.PhoneNumber,
			Hash:        req.Customer.Hash,
			InvoicingDetails: models.CustomerInvoicingDetails{
				CompanyName:    req.Customer.InvoicingDetails.CompanyName,
				CompanyAddress: req.Customer.InvoicingDetails.CompanyAddress,
				TaxID:          req.Customer.InvoicingDetails.TaxID,
			},
		},
		DeliveryAddress: models.DeliveryAddress{
			Label:     req.DeliveryAddress.Label,
			Longitude: req.DeliveryAddress.Longitude,
			Latitude:  req.DeliveryAddress.Latitude,
			City:      store.Address.City,
		},
		IsPickedUpByCustomer:   req.IsPickedUpByCustomer,
		CutleryRequested:       req.CutleryRequested,
		LoyaltyCard:            req.LoyaltyCard,
		BundledOrders:          req.BundledOrders,
		RestaurantSelfDelivery: s.getRestaurantSelfDelivery(req.DeliveryFee),
	}

	//проверка для отправки продукта "ПРИБОРЫ" на кассу
	if store.Settings.SendUtensilsToPos && req.CutleryRequested {
		req.Products = append(req.Products, models2.ProductOrder{
			ID:       store.Settings.UtensilsProductID,
			Quantity: 1,
			Name:     "ПРИБОРЫ",
		})
	}

	res.Products = make([]models.OrderProduct, 0, len(req.Products))
	for _, product := range req.Products {
		res.Products = append(res.Products, s.toOrderProduct(product))
	}

	return res, nil
}

func (s *glovoService) toOrderProduct(req models2.ProductOrder) models.OrderProduct {
	if req.Quantity == 0 {
		req.Quantity = 1
	}

	res := models.OrderProduct{
		ID:                 req.ID,
		PurchasedProductID: req.PurchasedProductID,
		Name:               req.Name,
		Price: models.Price{
			Value: float64(req.Price) / 100,
		},
		PriceWithoutDiscount: models.Price{
			Value: float64(req.Price) / 100,
		},
		Quantity: req.Quantity,
	}

	res.Attributes = make([]models.ProductAttribute, 0, len(req.Attributes))
	for _, attribute := range req.Attributes {
		if attribute.Quantity == 0 {
			attribute.Quantity = 1
		}

		res.Attributes = append(res.Attributes, models.ProductAttribute{
			ID:   attribute.ID,
			Name: attribute.Name,
			Price: models.Price{
				Value: float64(attribute.Price) / 100,
			},
			Quantity: attribute.Quantity,
		})

	}

	return res
}

func (s *glovoService) parseTimeWithOffset(timeStr, layout, offsetInMinutes string) (time.Time, error) {
	minutes, err := strconv.Atoi(offsetInMinutes)
	if err != nil {
		return time.Now(), err
	}

	seconds := minutes * 60
	loc := time.FixedZone("", seconds)
	t, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}

func (s *glovoService) MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	log.Info().Msg("delivery service: glovoService")

	for _, status := range store.Glovo.PurchaseTypes.Instant {
		if status.PosStatus == posStatus.String() {
			log.Info().Msgf("[SPECIAL MATCHING] pos status: %v, glovoService status: %v", status.PosStatus, status.Status)
			return status.Status
		}
	}

	log.Info().Msgf("[DEFAULT MATCHING], pos status: %v", posStatus)

	switch posStatus {
	case models.ACCEPTED, models.COOKING_STARTED, models.WAIT_SENDING:
		return models.ACCEPTED.String()
	case models.READY_FOR_PICKUP, models.COOKING_COMPLETE, models.CLOSED:
		return models.READY_FOR_PICKUP.String()
	case models.OUT_FOR_DELIVERY:
		return models.OUT_FOR_DELIVERY.String()
	case models.PICKED_UP_BY_CUSTOMER:
		return models.PICKED_UP_BY_CUSTOMER.String()
	}
	return ""
}

func (s *glovoService) UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error {
	if aggregatorStatus == "" {
		log.Info().Msgf("aggregator status is empty, order id %s", order.OrderID)
		return nil
	}

	if store.Glovo.AdjustedPickupMinutes != 0 || (store.Glovo.ScheduledBusyMode && len(store.Glovo.ScheduledBusyModeTime) > 0) {
		switch aggregatorStatus {
		case models.ACCEPTED.String():
			apm := store.Glovo.AdjustedPickupMinutes
			if apm == 0 {
				for _, scheduleTime := range store.Glovo.ScheduledBusyModeTime {
					if s.isNowWithinRange(scheduleTime.From, scheduleTime.To) {
						apm = 10
					}
				}
			}

			preparationTime := time.Now().UTC().Add(time.Duration(apm) * time.Minute)
			log.Info().Msgf("glovo accept order preparationTime: %v for order_id: %s", preparationTime, order.OrderID)
			_, err := s.cli.AcceptOrder(ctx, order.StoreID, order.OrderID, preparationTime)
			if err != nil {
				return err
			}
		case models.READY_FOR_PICKUP.String():
			_, err := s.cli.MarkOrderAsReady(ctx, order.StoreID, order.OrderID)
			if err != nil {
				return err
			}
		case models.OUT_FOR_DELIVERY.String():
			_, err := s.cli.MarkOrderAsOutForDelivery(ctx, order.StoreID, order.OrderID)
			if err != nil {
				return err
			}
		case models.PICKED_UP_BY_CUSTOMER.String():
			_, err := s.cli.MarkOrderAsCustomerPickedUp(ctx, order.StoreID, order.OrderID)
			if err != nil {
				return err
			}
		}
	} else {
		orderID, err := strconv.ParseInt(order.OrderID, 10, 64)
		if err != nil {
			return err
		}

		if _, err = s.cli.UpdateOrderStatus(ctx, glovoModels.OrderUpdateRequest{
			ID:      orderID,
			Status:  aggregatorStatus,
			StoreID: order.StoreID,
		}); err != nil {
			return err
		}

		log.Info().Msgf("success update glovo order, order_id=%v, status=%v", orderID, aggregatorStatus)

	}

	return nil
}

func (s *glovoService) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error) {
	trID, err := s.cli.UpdateStopListByProducts(ctx, aggregatorStoreID, products, isAvailable)
	if err != nil {
		return "", err
	}

	return trID, nil
}

func (s *glovoService) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error) {
	trID, err := s.cli.UpdateStopListByProductsBulk(ctx, aggregatorStoreID, products)
	if err != nil {
		return "", err
	}

	return trID, nil
}

func (s *glovoService) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error) {
	trID, err := s.cli.UpdateStopListByAttributesBulk(ctx, aggregatorStoreID, attributes)
	if err != nil {
		return "", err
	}

	return trID, nil
}

func (s *glovoService) SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	return nil, nil
}

func (s *glovoService) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	order, ok := req.(models2.Order)
	if !ok {
		return "", errors.New("casting error")
	}

	return order.StoreID, nil
}

func (s *glovoService) GetAggregatorOrder(ctx context.Context, orderID string) (models3.Order, error) {
	return models3.Order{}, nil
}

func (s *glovoService) isNowWithinRange(from time.Time, to time.Time) bool {
	now := time.Now()

	normalize := func(t time.Time) time.Time {
		return time.Date(1970, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	}

	fromNormalized := normalize(from)
	toNormalized := normalize(to)
	nowNormalized := normalize(now)

	return nowNormalized.After(fromNormalized) && nowNormalized.Before(toNormalized)
}

func (s *glovoService) SendOrderErrorNotification(ctx context.Context, req interface{}) error {
	return nil
}

func (s *glovoService) SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error {
	return nil
}
