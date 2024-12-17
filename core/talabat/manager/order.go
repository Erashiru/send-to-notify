package manager

import (
	"context"
	"fmt"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/talabat/models"
	"github.com/kwaaka-team/orders-core/domain/logger"
	orderCli "github.com/kwaaka-team/orders-core/pkg/order"
	storeCli "github.com/kwaaka-team/orders-core/pkg/store"
	"go.uber.org/zap"
)

type Order interface {
	CancelOrder(ctx context.Context, req models.CancelOrderRequest) error
	//CreateOrder(ctx context.Context, remoteID string, order models.CreateOrderRequest) (models.CreateOrderResponse, error)
}

type orderImplementation struct {
	orderCli orderCli.Client
	storeCli storeCli.Client
	logger   *zap.SugaredLogger
}

func NewOrder(orderCli orderCli.Client, storeCli storeCli.Client, logger *zap.SugaredLogger) Order {
	return &orderImplementation{
		orderCli: orderCli,
		storeCli: storeCli,
		logger:   logger,
	}
}

func (man *orderImplementation) CancelOrder(ctx context.Context, req models.CancelOrderRequest) error {
	if req.RemoteID == "" || req.RemoteOrderID == "" || req.Status == "" {
		man.logger.Error(logger.LoggerInfo{
			System:   "talabat cancel order response error",
			Response: fmt.Sprintf("invalid input, RemoteID: %s, RemoteOrderID: %s, status: %s", req.RemoteID, req.RemoteOrderID, req.Status),
		})
		return fmt.Errorf("invalid input")
	}

	err := man.orderCli.CancelOrderInPos(ctx, coreModels.CancelOrderInPos{
		OrderID:         req.RemoteOrderID,
		DeliveryService: "talabat",
		CancelReason: coreModels.CancelReason{
			Reason: req.Message,
		},
		PaymentStrategy: "PAY_NOTHING",
	})

	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "talabat cancel order response error",
			Response: err,
		})
		return err
	}

	return nil
}

//func (man *orderImplementation) CreateOrder(ctx context.Context, remoteID string, order models.CreateOrderRequest) (models.CreateOrderResponse, error) {
//	store, err := man.storeCli.FindStore(ctx, storeModels.StoreSelector{
//		TalabatRemoteBranchId: remoteID,
//		DeliveryService:       "talabat",
//	})
//	if err != nil {
//		man.logger.Error(logger.LoggerInfo{
//			System:   "talabat create order response error",
//			Response: err,
//		})
//		return models.CreateOrderResponse{}, err
//	}
//	req, err := toOrderRequest(order, store)
//	if err != nil {
//		man.logger.Error(logger.LoggerInfo{
//			System:   "talabat create order response error",
//			Response: err,
//		})
//		return models.CreateOrderResponse{}, err
//	}
//	req.PosType = store.PosType
//	req.RestaurantID = store.ID
//
//	req.LogMessages.FromDelivery, err = json.Marshal(order)
//	if err != nil {
//		man.logger.Error(logger.LoggerInfo{
//			System:   "talabat create order response error",
//			Response: err,
//		})
//		return models.CreateOrderResponse{}, err
//	}
//
//	res, err := man.orderCli.CreateOrder(ctx, req)
//	if err != nil {
//		man.logger.Error(logger.LoggerInfo{
//			System:   "talabat create order response error",
//			Response: err,
//		})
//		return models.CreateOrderResponse{}, err
//	}
//
//	return models.CreateOrderResponse{
//		RemoteResponse: models.RemoteResponse{
//			RemoteOrderId: res.OrderID,
//		},
//	}, nil
//}

//func toOrderRequest(req models.CreateOrderRequest, store coreStoreModels.Store) (coreOrderModels.Order, error) {
//	estimatedTotalPriceValue, err := strconv.ParseFloat(req.Price.GrandTotal, 64)
//	if err != nil {
//		return coreOrderModels.Order{}, err
//	}
//
//	var customerCashPaymentAmount float64
//	if req.Price.PayRestaurant != "" {
//		amount, err := strconv.ParseFloat(req.Price.PayRestaurant, 64)
//		if err != nil {
//			return coreOrderModels.Order{}, err
//		}
//		customerCashPaymentAmount = amount
//	}
//
//	orderProducts, err := toProducts(req.Products)
//	if err != nil {
//		return coreOrderModels.Order{}, err
//	}
//
//	order := coreOrderModels.Order{
//		OrderID:             req.Token,
//		OrderCode:           req.Code,
//		DeliveryService:     "talabat",
//		SpecialRequirements: req.Comments.CustomerComment,
//		AllergyInfo:         req.Comments.CustomerComment,
//		Customer: coreOrderModels.Customer{
//			Name:        req.Customer.FirstName + " " + req.Customer.LastName,
//			PhoneNumber: req.Customer.MobilePhone,
//		},
//		DeliveryAddress: coreOrderModels.DeliveryAddress{
//			Longitude: req.Delivery.Address.Longitude,
//			Latitude:  req.Delivery.Address.Latitude,
//			Label:     toLabel(req.Delivery.Address),
//		},
//		EstimatedPickupTime: coreOrderModels.TransactionTime{
//			Value:    coreOrderModels.Time{Time: toEstimatedPickUpTime(req.Delivery, store.Talabat.AdjustedPickupMinutes)},
//			TimeZone: store.Settings.TimeZone.TZ,
//		},
//		PickUpCode:       req.ShortCode,
//		UtcOffsetMinutes: strconv.Itoa(int(store.Settings.TimeZone.UTCOffset)),
//		EstimatedTotalPrice: coreOrderModels.Price{
//			Value:        estimatedTotalPriceValue,
//			CurrencyCode: store.Settings.Currency,
//		},
//		DeliveryFee: coreOrderModels.Price{
//			Value:        toDeliveryFeeValue(req.Price.DeliveryFees),
//			CurrencyCode: store.Settings.Currency,
//		},
//		CustomerCashPaymentAmount: coreOrderModels.Price{
//			Value:        customerCashPaymentAmount,
//			CurrencyCode: store.Settings.Currency,
//		},
//		Products: orderProducts,
//		OrderTime: coreOrderModels.TransactionTime{
//			Value:    coreOrderModels.Time{Time: req.CreatedAt.UTC()},
//			TimeZone: store.Settings.TimeZone.TZ,
//		},
//		StatusesHistory: []coreOrderModels.OrderStatusUpdate{
//			{
//				Name: "NEW",
//				Time: time.Now().UTC(),
//			},
//		},
//	}
//
//	if req.ExpeditionType == "pickup" {
//		order.IsPickedUpByCustomer = true
//		order.PickUpCode = req.Pickup.PickupCode
//		order.EstimatedPickupTime = coreOrderModels.TransactionTime{
//			Value:    coreOrderModels.Time{Time: toEstimatedPickUpTimeCustomerPickUp(req.Pickup.PickupTime, store.Talabat.AdjustedPickupMinutes)},
//			TimeZone: store.Settings.TimeZone.TZ,
//		}
//	}
//
//	switch req.Payment.Status {
//	case "paid":
//		order.PaymentMethod = "DELAYED"
//	case "pending":
//		order.PaymentMethod = "CASH"
//	}
//
//	switch req.PreOrder {
//	case true:
//		order.Type = "PREORDER"
//		order.Preorder = coreOrderModels.PreOrder{
//			Time: coreOrderModels.TransactionTime{
//				Value:    order.EstimatedPickupTime.Value,
//				TimeZone: store.Settings.TimeZone.TZ,
//			},
//		}
//	case false:
//		order.Type = "INSTANT"
//	}
//
//	return order, nil
//}

//func toProducts(req []models.Product) ([]coreOrderModels.OrderProduct, error) {
//	res := make([]coreOrderModels.OrderProduct, 0, len(req))
//	for _, p := range req {
//		product, err := toProduct(p)
//		if err != nil {
//			return nil, err
//		}
//		res = append(res, product)
//	}
//	return res, nil
//}

//func toProduct(req models.Product) (coreOrderModels.OrderProduct, error) {
//	productUnitPrice, err := strconv.ParseFloat(req.UnitPrice, 64)
//	if err != nil {
//		return coreOrderModels.OrderProduct{}, err
//	}
//
//	productQuantity, err := strconv.Atoi(req.Quantity)
//	if err != nil {
//		return coreOrderModels.OrderProduct{}, err
//	}
//
//	res := coreOrderModels.OrderProduct{
//		ID:   req.RemoteCode,
//		Name: req.Name,
//		Price: coreOrderModels.Price{
//			Value: productUnitPrice,
//		},
//		Quantity: productQuantity,
//	}
//	res.Attributes = make([]coreOrderModels.ProductAttribute, 0, len(req.SelectedToppings))
//	for _, t := range req.SelectedToppings {
//		price, err := strconv.ParseFloat(t.Price, 64)
//		if err != nil {
//			return coreOrderModels.OrderProduct{}, err
//		}
//		res.Attributes = append(res.Attributes, coreOrderModels.ProductAttribute{
//			ID:       t.RemoteCode,
//			Quantity: t.Quantity,
//			Price: coreOrderModels.Price{
//				Value: price,
//			},
//			Name: t.Name,
//		})
//	}
//	return res, nil
//}

//func toDeliveryFeeValue(req []models.DeliveryFees) float64 {
//	var res float64
//	for _, f := range req {
//		res = res + f.Value
//	}
//	return res
//}
//
//func toEstimatedPickUpTimeCustomerPickUp(reqTime string, adjuctedPickUpMinutes int) time.Time {
//	layout := "2006-01-02 15:04:05"
//	if reqTime != "" {
//		estimatedPickUpTime, err := time.Parse(layout, reqTime)
//		if err != nil {
//			estimatedPickUpTime = time.Now().Add(time.Duration(adjuctedPickUpMinutes) * time.Minute)
//		}
//		return estimatedPickUpTime
//	}
//	return time.Now().Add(time.Duration(adjuctedPickUpMinutes) * time.Minute)
//}
//
//func toEstimatedPickUpTime(delivery models.Delivery, adjuctedPickUpMinutes int) time.Time {
//	layout := "2006-01-02 15:04:05"
//
//	if delivery.ExpectedDeliveryTime != "" {
//		estimatedPickUpTime, err := time.Parse(layout, delivery.ExpectedDeliveryTime)
//		if err != nil {
//			estimatedPickUpTime = time.Now().Add(time.Duration(adjuctedPickUpMinutes) * time.Minute)
//		}
//		return estimatedPickUpTime
//	}
//
//	if delivery.RiderPickupTime != "" {
//		estimatedPickUpTime, err := time.Parse(layout, delivery.RiderPickupTime)
//		if err != nil {
//			estimatedPickUpTime = time.Now().Add(time.Duration(adjuctedPickUpMinutes) * time.Minute)
//		}
//		return estimatedPickUpTime
//	}
//
//	return time.Now().Add(time.Duration(adjuctedPickUpMinutes) * time.Minute)
//}
//
//func toLabel(address models.Address) string {
//	res := ""
//	if address.City != "" {
//		res = res + address.City + ","
//	}
//
//	if address.Street != "" {
//		res = res + address.Street + ","
//	}
//
//	if address.Number != "" {
//		res = res + address.Number + ","
//	}
//	if address.DeliveryArea != "" {
//		res = res + address.DeliveryArea + ","
//	}
//
//	if address.Entrance != "" {
//		res = res + address.Entrance + ","
//	}
//
//	if address.Intercom != "" {
//		res = res + "intercom: " + address.Intercom + ","
//	}
//
//	if address.FlatNumber != "" {
//		res = res + address.FlatNumber + ","
//	}
//
//	if address.Building != "" {
//		res = res + address.Building + ","
//	}
//
//	if address.DeliveryInstructions != "" {
//		res = res + address.DeliveryInstructions + ","
//	}
//
//	if len(res) > 0 && res[len(res)-1:] == "," {
//		res = res[:len(res)-1]
//	}
//
//	return res
//}
