package aggregator

import (
	"context"
	expressModels "github.com/kwaaka-team/orders-core/core/express24/models"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models3 "github.com/kwaaka-team/orders-core/core/wolt/models"
	expressConf "github.com/kwaaka-team/orders-core/pkg/express24"
	expressCli "github.com/kwaaka-team/orders-core/pkg/express24/clients"
	"github.com/kwaaka-team/orders-core/pkg/express24/clients/dto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type express24Service struct {
	deliveryServiceName models.Aggregator
	cli                 expressCli.Express24
}

func newExpress24Service(baseUrl string, store storeModels.Store) (*express24Service, error) {

	cli, err := expressConf.NewExpress24Client(&expressCli.Config{
		Protocol: "http",
		BaseURL:  baseUrl,
		Username: store.Express24.Username,
		Password: store.Express24.Password,
	})
	if err != nil {
		return nil, errors.Wrap(constructorError, err.Error())
	}

	return &express24Service{
		models.EXPRESS24, cli,
	}, nil
}

func (s *express24Service) MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	return ""
}

func (s *express24Service) OpenStore(ctx context.Context, aggregatorStoreId string) error {
	return errors.New("method not implemented")
}

func (s *express24Service) GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error) {
	return false, errors.New("method not implemented")
}

func (s *express24Service) GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error) {
	return storeModels.AggregatorSchedule{}, errors.New("method not implemented")
}

func (s *express24Service) UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error {
	return nil
}

func (s *express24Service) IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error) {
	return store.Express24.IsMarketplace, nil
}

func (s *express24Service) SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	return nil, nil
}

func (s *express24Service) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	order, ok := req.(expressModels.Order)
	if !ok {
		return "", errors.New("casting error")
	}

	return order.Store.Branch.ExternalId, nil
}

func (s *express24Service) GetSystemCreateOrderRequestByAggregatorRequest(req interface{}, store storeModels.Store) (models.Order, error) {

	order, ok := req.(expressModels.Order)
	if !ok {
		return models.Order{}, errors.New("casting error")
	}

	return toOrderRequest(order, store), nil
}

func toOrderRequest(req expressModels.Order, store storeModels.Store) models.Order {

	var orderType string
	if req.Status == "pre_order" {
		orderType = "PREORDER"
	} else {
		orderType = "INSTANT"
	}

	var pickUpTime = req.CreatedAt.UTC().Add(time.Duration(store.Express24.AdjustedPickupMinutes) * time.Minute)

	var paymentMethod string

	if req.Payment.Id == "cash" {
		paymentMethod = "CASH"
	} else {
		paymentMethod = "DELAYED"
	}

	res := models.Order{
		OrderID:         strconv.Itoa(req.Id),
		StoreID:         req.Store.Branch.ExternalId,
		OrderCode:       strconv.Itoa(req.Id),
		RestaurantID:    store.ID,
		PosType:         store.PosType,
		Type:            orderType,
		DeliveryService: expressModels.EXPRESS24.String(),
		Status:          expressModels.STATUS_NEW.String(),
		StatusesHistory: []models.OrderStatusUpdate{
			{
				Name: expressModels.STATUS_NEW.String(),
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
		PickUpCode: strconv.Itoa(req.Id),
		Preorder: models.PreOrder{ // TODO: check if it's correct
			Time: models.TransactionTime{
				Value:    models.Time{Time: req.CreatedAt.UTC()},
				TimeZone: store.Settings.TimeZone.TZ,
			},
			Status: req.Status,
		},
		UtcOffsetMinutes:    strconv.Itoa(int(store.Settings.TimeZone.UTCOffset)),
		Currency:            store.Settings.Currency,
		SpecialRequirements: req.OrderComment,
		EstimatedTotalPrice: models.Price{
			Value:        req.TotalPrice - req.Delivery.Price,
			CurrencyCode: store.Settings.Currency,
		},
		TotalCustomerToPay: models.Price{
			Value:        req.TotalPrice,
			CurrencyCode: store.Settings.Currency,
		},
		DeliveryFee: models.Price{
			Value:        req.Delivery.Price,
			CurrencyCode: store.Settings.Currency,
		},
		Customer: models.Customer{
			Name:        "Express24",
			PhoneNumber: "+998000000000",
		},
		DeliveryAddress: models.DeliveryAddress{
			Label:     req.Delivery.Address.Text,
			Latitude:  req.Delivery.Address.Lat,
			Longitude: req.Delivery.Address.Lon,
		},
		PaymentMethod:        paymentMethod,
		IsPickedUpByCustomer: false,
	}

	res.Products = make([]models.OrderProduct, 0, len(req.Products))
	for _, product := range req.Products {
		res.Products = append(res.Products, toOrderProduct(product, store))
	}
	return res
}

func toOrderProduct(reqProduct expressModels.Product, store storeModels.Store) models.OrderProduct {
	res := models.OrderProduct{
		ID:   reqProduct.ExternalId,
		Name: reqProduct.Name,
		Price: models.Price{
			Value:        reqProduct.Price,
			CurrencyCode: store.Settings.Currency,
		},
		Quantity: reqProduct.Qty,
	}

	res.Attributes = make([]models.ProductAttribute, 0)
	for _, attributeGroup := range reqProduct.Params {
		for _, attribute := range attributeGroup.Options {
			res.Attributes = append(res.Attributes, models.ProductAttribute{
				ID:      attribute.ExternalId,
				GroupID: "",
				Name:    attribute.Name,
				Price: models.Price{
					Value:        attribute.Price,
					CurrencyCode: store.Settings.Currency,
				},
				Quantity: 1,
			})
		}
	}

	return res
}

func (s *express24Service) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error) {
	branch, err := strconv.Atoi(aggregatorStoreID)
	if err != nil {
		return "", err
	}

	resp, err := s.cli.UpdateProducts(ctx, dto.UpdateProductsRequest{
		Data: dto.UpdateProductData{
			Branches: []int{branch},
			Products: s.toProducts(products, isAvailable),
		},
	})

	if err != nil {
		return "", err
	}

	if resp.Failed != nil {
		for _, fail := range resp.Failed {
			log.Info().Msgf("PRODUCT ID: %s, ERROR MESSAGE: %s", fail.ExternalId, fail.Message)
		}
	}
	return "", nil
}

func (s *express24Service) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error) {
	branch, err := strconv.Atoi(aggregatorStoreID)
	if err != nil {
		return "", err
	}

	resp, err := s.cli.UpdateProducts(ctx, dto.UpdateProductsRequest{
		Data: dto.UpdateProductData{
			Branches: []int{branch},
			Products: s.toProductsBulk(products),
		},
	})

	if err != nil {
		return "", err
	}

	if resp.Failed != nil {
		for _, fail := range resp.Failed {
			log.Info().Msgf("PRODUCT ID: %s, ERROR MESSAGE: %s", fail.ExternalId, fail.Message)
		}
	}
	return "", nil
}

func (s *express24Service) toProducts(req menuModels.Products, isAvailable bool) []dto.Product {
	products := make([]dto.Product, 0, len(req))

	for i := range req {
		var price int
		if len(req[i].Price) > 0 {
			price = int(req[i].Price[0].Value)
		}

		products = append(products, dto.Product{
			ExternalId:  req[i].ExtID,
			Quantity:    menuModels.BASEQUANTITY,
			IsAvailable: s.toIsAvailable(isAvailable),
			Price:       price,
		})
	}
	return products
}

func (s *express24Service) toProductsBulk(req menuModels.Products) []dto.Product {
	products := make([]dto.Product, 0, len(req))

	for i := range req {
		var price int
		if len(req[i].Price) > 0 {
			price = int(req[i].Price[0].Value)
		}

		products = append(products, dto.Product{
			ExternalId:  req[i].ExtID,
			Quantity:    menuModels.BASEQUANTITY,
			IsAvailable: s.toIsAvailable(req[i].IsAvailable),
			Price:       price,
		})
	}
	return products
}

func (s *express24Service) toIsAvailable(isAvailable bool) int {
	if isAvailable {
		return 1
	}
	return 0
}

func (s *express24Service) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error) {
	return "", nil
}

func (s *express24Service) GetAggregatorOrder(ctx context.Context, orderID string) (models3.Order, error) {
	return models3.Order{}, nil
}

func (s *express24Service) SendOrderErrorNotification(ctx context.Context, req interface{}) error {
	return nil
}

func (s *express24Service) SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error {
	return nil
}
