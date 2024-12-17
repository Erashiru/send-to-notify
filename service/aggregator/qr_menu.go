package aggregator

import (
	"context"
	"fmt"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	models2 "github.com/kwaaka-team/orders-core/core/qrmenu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models3 "github.com/kwaaka-team/orders-core/core/wolt/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type qrMenuService struct {
}

func newQrMenuService() (*qrMenuService, error) {
	return &qrMenuService{}, nil
}

func (s *qrMenuService) OpenStore(ctx context.Context, aggregatorStoreId string) error {
	return errors.New("method not implemented")
}

func (s *qrMenuService) GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error) {
	return false, errors.New("method not implemented")
}

func (s *qrMenuService) GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error) {
	return storeModels.AggregatorSchedule{}, errors.New("method not implemented")
}

func (s *qrMenuService) IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error) {
	return store.QRMenu.IsMarketplace, nil
}

func (s *qrMenuService) GetAggregatorCreateOrderRequestBySystemRequest(req models.Order) (interface{}, error) {
	return nil, nil
}

func (s *qrMenuService) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	order, ok := req.(models2.Order)
	if !ok {
		return "", errors.New("casting error")
	}

	return order.RestaurantID, nil
}

func (s *qrMenuService) SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	order, ok := req.(models2.Order)
	if !ok {
		return nil, errors.New("casting error")
	}

	childRestaurantOrders := make(map[string]models2.Order)

	for i := range order.Items {
		product := order.Items[i]

		productRestaurantID, productID, err := splitVirtualStoreItemID(product.ProductID, "_")
		if err != nil {
			log.Err(errors.New("not valid signature, len will be 2 with _")).Msgf("orders core, splitVirtualStoreItemID error, id: %s", product.ProductID)
			continue
		}
		childProduct := product
		childProduct.ProductID = productID
		childOrder, ok := childRestaurantOrders[productRestaurantID]
		if !ok {
			child := order
			child.Items = make([]models2.Product, 0, 1)
			child.Items = append(child.Items, childProduct)
			child.RestaurantID = productRestaurantID
			child.ID = order.ID + "_" + productRestaurantID
			child.Delivery.Dispatcher = ""
			child.Delivery.ClientDeliveryPrice = 0
			child.Delivery.FullDeliveryPrice = 0
			child.Delivery.KwaakaChargedDeliveryPrice = 0
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

func (s *qrMenuService) calculateChildOrderTotalSum(order models2.Order) models2.Order {
	totalSum := float64(0)

	for _, p := range order.Items {
		totalSum = totalSum + p.Price*float64(p.Quantity)
		for _, a := range p.Attributes {
			totalSum = totalSum + a.Price*float64(a.Quantity*p.Quantity)
		}
	}

	order.TotalSum = totalSum
	return order
}

func splitVirtualStoreItemID(id string, sep string) (string, string, error) {
	result := strings.Split(id, sep)
	if len(result) != 2 && len(result) != 3 {
		return "", "", errors.New("not valid signature, len will be 2 with _")
	}

	return result[0], result[1], nil
}

func (s *qrMenuService) GetSystemCreateOrderRequestByAggregatorRequest(r interface{}, store storeModels.Store) (models.Order, error) {
	req, ok := r.(models2.Order)
	if !ok {
		return models.Order{}, errors.New("casting error")
	}
	orderType := "INSTANT"
	orderCreatedAt := time.Now().UTC()
	cookingTime := s.getCookingTime(req.Items, store.QRMenu.CookingTime)
	estimatedPickupTime := orderCreatedAt.Add(time.Duration(cookingTime) * time.Minute)

	var preorder models.PreOrder
	if !req.PreOrderTime.IsZero() {
		if req.PreOrderTime.UTC().Before(time.Now().UTC()) {
			return models.Order{}, errors.New("preorder time is before than now")
		}
		orderType = "PREORDER"
		estimatedPickupTime = req.PreOrderTime.UTC().Add(-time.Minute * time.Duration(req.Delivery.DeliveryTime))

		if req.IsPickedUpByCustomer {
			estimatedPickupTime = req.PreOrderTime.UTC()
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
			Status: "waiting",
			Time: models.TransactionTime{
				Value: models.Time{
					Time: req.PreOrderTime.UTC(),
				},
				TimeZone: store.Settings.TimeZone.TZ,
			},
		}
	}

	var phoneNumberWithPlus string

	if !strings.Contains(req.Customer.PhoneNumber, "+") {
		phoneNumberWithPlus = "+" + req.Customer.PhoneNumber
	}

	pickupCode := 100 + rand.Intn(900)
	orderID := req.ID

	sendCourier := true
	if req.IsPickedUpByCustomer {
		sendCourier = false
	}

	res := models.Order{
		RestaurantID:    req.RestaurantID,
		OrderID:         orderID,
		StoreID:         req.RestaurantID,
		Type:            orderType,
		DeliveryService: "qr_menu",
		PickUpCode:      strconv.Itoa(pickupCode),
		Status:          "NEW",
		StatusesHistory: []models.OrderStatusUpdate{
			{
				Name: "NEW",
				Time: time.Now().UTC(),
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
		PaymentMethod: req.PaymentType,
		Currency:      req.Currency,
		AllergyInfo:   req.OrderComment,
		EstimatedTotalPrice: models.Price{
			Value:        req.TotalSum,
			CurrencyCode: req.Currency,
		},
		TotalCustomerToPay: models.Price{
			Value:        req.TotalSum,
			CurrencyCode: req.Currency,
		},
		Customer: models.Customer{
			Name:                req.Customer.Name,
			PhoneNumber:         req.Customer.PhoneNumber,
			PhoneNumberWithPlus: phoneNumberWithPlus,
		},
		DeliveryAddress: models.DeliveryAddress{
			Label:        s.toLabel(req.DeliveryAddress),
			Longitude:    req.DeliveryAddress.Coordinates.Lon,
			Latitude:     req.DeliveryAddress.Coordinates.Lat,
			City:         req.DeliveryAddress.City,
			Comment:      req.DeliveryAddress.Comment,
			BuildingName: req.DeliveryAddress.BuildingName,
			Street:       req.DeliveryAddress.Street,
			Flat:         req.DeliveryAddress.Apartment,
			Porch:        req.DeliveryAddress.Entrance,
			Floor:        req.DeliveryAddress.Floor,
		},
		SpecialRequirements:         req.DeliveryAddress.Comment,
		IsPickedUpByCustomer:        true,
		Preorder:                    preorder,
		IsInstantDelivery:           store.Kwaaka3PL.IsInstantCall,
		DeliveryDispatcher:          req.Delivery.Dispatcher,
		DispatcherDeliveryTime:      req.Delivery.DeliveryTime,
		FullDeliveryPrice:           req.Delivery.FullDeliveryPrice,
		ClientDeliveryPrice:         req.Delivery.ClientDeliveryPrice,
		KwaakaChargedDeliveryPrice:  req.Delivery.KwaakaChargedDeliveryPrice,
		RestaurantPayDeliveryPrice:  s.getRestaurantPayDeliveryPrice(req.Delivery.FullDeliveryPrice, req.Delivery.ClientDeliveryPrice),
		DeliveryDropOffScheduleTime: req.Delivery.DropOffScheduleTime,
		RestaurantSelfDelivery:      true,
		SendCourier:                 sendCourier,
		PromoCode:                   req.PromoCode,
	}

	res.Products = make([]models.OrderProduct, 0, len(req.Items))
	for _, product := range req.Items {
		res.Products = append(res.Products, s.toOrderProduct(product))
	}

	if len(req.PromoCode) != 0 {
		res.AllergyInfo += " Промокод от Direct: " + req.PromoCode
	}

	return res, nil
}

func (s *qrMenuService) toOrderProduct(req models2.Product) models.OrderProduct {
	res := models.OrderProduct{
		ID:   req.ProductID,
		Name: req.Name,
		Price: models.Price{
			Value: req.Price,
		},
		Quantity: req.Quantity,
	}

	res.Attributes = make([]models.ProductAttribute, 0, len(req.Attributes))
	for _, attribute := range req.Attributes {
		res.Attributes = append(res.Attributes, models.ProductAttribute{
			ID:       attribute.AttributeID,
			Quantity: attribute.Quantity,
			Price: models.Price{
				Value: attribute.Price,
			},
			Name: attribute.Name,
		})
	}

	return res
}

func (s *qrMenuService) getRestaurantPayDeliveryPrice(fullDeliveryPrice float64, clientDeliveryPrice float64) float64 {
	diff := fullDeliveryPrice - clientDeliveryPrice

	if diff < 0 {
		return 0
	}

	return diff
}

func (s *qrMenuService) toLabel(address models2.DeliveryAddress) string {
	label := ""
	if address.City != "" {
		label = label + address.City + ", "
	}
	if address.Street != "" {
		label = label + address.Street + ", "
	}
	if address.Apartment != "" {
		label = label + "квартира " + address.Apartment + ", "
	}
	if address.Entrance != "" {
		label = label + "подъезд " + address.Entrance + ", "
	}
	if address.Floor != "" {
		label = label + "этаж " + address.Floor + ", "
	}
	if address.BuildingName != "" {
		label = label + address.BuildingName + ", "
	}
	if address.DoorBellInfo != "" {
		label = label + "домофон " + address.DoorBellInfo + ", "
	}
	if address.LocationType != "" {
		label = label + address.LocationType + ", "
	}

	label = strings.TrimSuffix(label, ", ")

	return label
}

func (s *qrMenuService) getCookingTime(items []models2.Product, storeCookingTime int32) int32 {
	var cookingTime int32
	for _, product := range items {
		if product.CookingTime < cookingTime {
			continue
		}
		cookingTime = product.CookingTime
	}
	if cookingTime == 0 {
		cookingTime = storeCookingTime
	}

	return cookingTime
}

func (s *qrMenuService) MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	return posStatus.String()
}

func (s *qrMenuService) UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error {
	return nil
}

func (s *qrMenuService) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error) {
	return "", nil
}

func (s *qrMenuService) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error) {
	return "", nil
}

func (s *qrMenuService) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error) {
	return "", nil
}

func (s *qrMenuService) GetAggregatorOrder(ctx context.Context, orderID string) (models3.Order, error) {
	return models3.Order{}, nil
}

func (s *qrMenuService) SendOrderErrorNotification(ctx context.Context, req interface{}) error {
	return nil
}

func (s *qrMenuService) SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error {
	return nil
}
