package aggregator

import (
	externalApiModels "github.com/kwaaka-team/orders-core/core/externalapi/models"
	externalApiUtils "github.com/kwaaka-team/orders-core/core/externalapi/utils"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"regexp"
	"strconv"
	"time"
)

type externalService struct {
}

func newExternalService() (*externalService, error) {
	return &externalService{}, nil
}

func (s *externalService) mapSystemStatusToAggregatorStatus(aggregatorName models.Aggregator, order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	for _, external := range store.ExternalConfig {
		if external.Type != aggregatorName.String() {
			continue
		}
		for _, status := range external.PurchaseTypes.Instant {
			if status.PosStatus != posStatus.String() {
				continue
			}
			return status.Status
		}
	}

	log.Info().Msgf("[DEFAULT MATCHING], pos status: %v", posStatus)
	return ""
}

func (s *externalService) getRestaurantSelfDelivery(discriminator string) bool {
	if discriminator == "marketplace" {
		return true
	}

	return false
}

func (s *externalService) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	order, ok := req.(externalApiModels.Order)
	if !ok {
		return "", errors.New("casting error")
	}

	return order.RestaurantId, nil
}

func (s *externalService) splitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	order, ok := req.(externalApiModels.Order)
	if !ok {
		return nil, errors.New("casting error")
	}

	childRestaurantOrders := make(map[string]externalApiModels.Order)

	for i := range order.Items {
		product := order.Items[i]

		productRestaurantID, productID, err := splitVirtualStoreItemID(product.Id, "_")
		if err != nil {
			log.Err(errors.New("not valid signature, len will be 2 with _")).Msgf("orders core, splitVirtualStoreItemID error, id: %s", product.Id)
			continue
		}
		childProduct := product
		childProduct.Id = productID
		childOrder, ok := childRestaurantOrders[productRestaurantID]
		if !ok {
			child := order
			child.Items = make([]externalApiModels.OrderItem, 0, 1)
			child.Items = append(child.Items, childProduct)
			child.RestaurantId = child.RestaurantId + "_" + productRestaurantID
			child.EatsId = order.EatsId + "_" + productRestaurantID

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

func (s *externalService) calculateChildOrderTotalSum(order externalApiModels.Order) externalApiModels.Order {
	totalSum := 0

	for _, p := range order.Items {
		totalSum = totalSum + p.Price*int(p.Quantity)
	}

	order.PaymentInfo.ItemsCost = totalSum
	return order
}

func (s *externalService) getSystemCreateOrderRequestByAggregatorRequest(req interface{}, store storeModels.Store, deliveryServiceName string) (models.Order, error) {

	o, ok := req.(externalApiModels.Order)
	if !ok {
		return models.Order{}, errors.New("casting error")
	}

	var (
		deliveryDate     time.Time
		err              error
		pickUpByCustomer bool
		isMarketplace    bool
	)

	switch o.Discriminator {
	case "pickup":
		pickUpByCustomer = true
		deliveryDate, err = time.Parse(time.RFC3339Nano, o.DeliveryInfo.ClientArrivementDate)
		if err != nil {
			return models.Order{}, errors.New("pickup invalid date type")
		}
	case "yandex":
		pickUpByCustomer = true
		deliveryDate, err = time.Parse(time.RFC3339Nano, o.DeliveryInfo.DeliveryDate)
		if err != nil {
			return models.Order{}, errors.New("yandex invalid date type")
		}
	case "marketplace":
		isMarketplace = true
		deliveryDate, err = time.Parse(time.RFC3339Nano, o.DeliveryInfo.MarketPlaceDeliveryDate)
		if err != nil {
			return models.Order{}, errors.New("marketplace invalid date type")
		}
	}

	var address models.DeliveryAddress
	if o.DeliveryInfo.DeliveryAddress != nil {
		address, err = s.addressByFields(o.DeliveryInfo.DeliveryAddress)
		if err != nil {
			address = models.DeliveryAddress{
				Label: o.DeliveryInfo.DeliveryAddress.Full,
				City:  store.Address.City,
			}
		}
	}

	order := models.Order{
		Type:            "INSTANT",
		Discriminator:   o.Discriminator,
		DeliveryService: deliveryServiceName,
		PosType:         store.PosType,
		RestaurantID:    store.ID,
		OrderID:         o.EatsId,
		OrderCode:       o.EatsId,
		StoreID:         o.RestaurantId,
		Status:          "NEW",
		StatusesHistory: []models.OrderStatusUpdate{
			{
				Name: "NEW",
				Time: time.Now().UTC(),
			},
		},
		OrderTime: models.TransactionTime{
			Value:     models.Time{Time: time.Now().UTC()},
			TimeZone:  store.Settings.TimeZone.TZ,
			UTCOffset: strconv.Itoa(int(store.Settings.TimeZone.UTCOffset)),
		},
		EstimatedPickupTime: models.TransactionTime{
			Value:     models.Time{Time: deliveryDate.UTC()},
			TimeZone:  store.Settings.TimeZone.TZ,
			UTCOffset: strconv.Itoa(int(store.Settings.TimeZone.UTCOffset)),
		},
		Currency:            store.Settings.Currency,
		SpecialRequirements: o.Comment,
		AllergyInfo:         o.Comment,
		PickUpCode:          o.EatsId,
		EstimatedTotalPrice: models.Price{
			Value:        float64(o.PaymentInfo.ItemsCost),
			CurrencyCode: store.Settings.Currency,
		},
		TotalCustomerToPay: models.Price{
			Value:        float64(o.PaymentInfo.ItemsCost),
			CurrencyCode: store.Settings.Currency,
		},
		Customer: models.Customer{
			Name:        o.DeliveryInfo.ClientName,
			PhoneNumber: s.phoneNumber(o.DeliveryInfo.PhoneNumber),
		},
		Persons:                o.Persons,
		IsPickedUpByCustomer:   pickUpByCustomer,
		CutleryRequested:       true,
		DeliveryAddress:        address,
		IsMarketplace:          isMarketplace,
		RestaurantSelfDelivery: s.getRestaurantSelfDelivery(o.Discriminator),
	}

	for _, promo := range o.Promos {
		order.Promos = append(order.Promos, models.Promo{
			Discount: promo.Discount,
			Type:     promo.Type,
		})
	}

	switch o.PaymentInfo.PaymentType {
	case "CASH":
		order.PaymentMethod = "CASH"
	default:
		order.PaymentMethod = "DELAYED"
	}

	//проверка для отправки продукта "ПРИБОРЫ" на кассу
	if store.Settings.SendUtensilsToPos && o.Persons > 0 {
		o.Items = append(o.Items, externalApiModels.OrderItem{
			Id:       store.Settings.UtensilsProductID,
			Quantity: float64(o.Persons),
			Name:     "ПРИБОРЫ",
		})
	}

	order.Products = make([]models.OrderProduct, 0, len(o.Items))
	for index, item := range o.Items {
		product := models.OrderProduct{
			ID:                 item.Id,
			PurchasedProductID: item.Id,
			Name:               item.Name,
			Quantity:           int(item.Quantity),
		}

		for _, promo := range item.Promos {
			product.Promos = append(product.Promos, models.Promo{
				Discount: promo.Discount,
				Type:     promo.Type,
			})
		}

		for _, modifier := range item.Modifications {
			product.Attributes = append(product.Attributes, models.ProductAttribute{
				ID:       modifier.Id,
				Name:     modifier.Name,
				Quantity: modifier.Quantity,
				Price: models.Price{
					Value:        float64(modifier.Price),
					CurrencyCode: store.Settings.Currency,
				},
			})
		}

		if len(order.Promos) > index {
			product.Price.Value = float64(item.Price - (order.Promos[index].Discount / product.Quantity))
			product.Price.CurrencyCode = store.Settings.Currency
		} else {
			product.Price.Value = float64(item.Price)
			product.Price.CurrencyCode = store.Settings.Currency
		}

		order.Products = append(order.Products, product)
	}

	return order, nil
}

func (s *externalService) addressByFields(address *externalApiModels.MarketPlaceDeliveryAddress) (models.DeliveryAddress, error) {
	yAddrss, err := externalApiUtils.ExtractYandexAddressInfo(address.Full)
	if err != nil {
		return models.DeliveryAddress{}, err
	}

	lat, err := strconv.ParseFloat(address.Latitude, 64)
	if err != nil {
		return models.DeliveryAddress{}, err
	}

	lon, err := strconv.ParseFloat(address.Longitude, 64)
	if err != nil {
		return models.DeliveryAddress{}, err
	}

	return models.DeliveryAddress{
		Label:       address.Full,
		Latitude:    lat,
		Longitude:   lon,
		City:        yAddrss.City,
		Street:      yAddrss.Street,
		HouseNumber: yAddrss.HouseNumber,
		Entrance:    yAddrss.Entrance,
		Intercom:    yAddrss.Intercom,
		Floor:       yAddrss.Floor,
	}, nil
}

func (s *externalService) phoneNumber(phoneNumber string) string {
	if len(phoneNumber) <= 15 {
		return phoneNumber
	}

	re := regexp.MustCompile(`^\+\d+`)
	res := re.FindString(phoneNumber)

	return res
}
