package aggregator

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	models2 "github.com/kwaaka-team/orders-core/core/kwaaka_admin/models"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models3 "github.com/kwaaka-team/orders-core/core/wolt/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	CognitoUserPoolId = "eu-west-1_bUR412LMx"
	OperatorId        = "custom:operator_id"
)

type kwaakaAdminService struct {
	cognitoService *cognitoidentityprovider.CognitoIdentityProvider
}

// PAY ATTENTION: using cognito in places where there was no declaration may cause panic
func newKwaakaAdminService(cognito *cognitoidentityprovider.CognitoIdentityProvider) (*kwaakaAdminService, error) {

	return &kwaakaAdminService{
		cognitoService: cognito,
	}, nil
}

func (s *kwaakaAdminService) OpenStore(ctx context.Context, aggregatorStoreId string) error {
	return errors.New("method not implemented")
}

func (s *kwaakaAdminService) GetStoreStatus(ctx context.Context, aggregatorStoreId string) (bool, error) {
	return false, errors.New("method not implemented")
}

func (s *kwaakaAdminService) GetStoreSchedule(ctx context.Context, aggregatorStoreId string) (storeModels.AggregatorSchedule, error) {
	return storeModels.AggregatorSchedule{}, errors.New("method not implemented")
}

func (s *kwaakaAdminService) IsMarketPlace(restaurantSelfDelivery bool, store storeModels.Store) (bool, error) {
	return store.KwaakaAdmin.IsMarketPlace, nil
}

func (s *kwaakaAdminService) SplitVirtualStoreOrder(req interface{}, store storeModels.Store) ([]interface{}, error) {
	return nil, nil
}

func (s *kwaakaAdminService) GetStoreIDFromAggregatorOrderRequest(req interface{}) (string, error) {
	order, ok := req.(models2.Order)
	if !ok {
		return "", errors.New("casting error")
	}

	return order.RestaurantID, nil
}

func (s *kwaakaAdminService) GetSystemCreateOrderRequestByAggregatorRequest(r interface{}, store storeModels.Store) (models.Order, error) {
	req, ok := r.(models2.Order)
	if !ok {
		return models.Order{}, errors.New("casting error")
	}
	orderType := "INSTANT"
	orderCreatedAt := time.Now().UTC()
	cookingTime := s.getCookingTime(req.Items, store.KwaakaAdmin.CookingTime)
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

	pickupCode := 100 + rand.Intn(900)

	var phoneNumberWithPlus string

	if !strings.Contains(req.Customer.PhoneNumber, "+") {
		phoneNumberWithPlus = "+" + req.Customer.PhoneNumber
	}

	sendCourier := true
	if req.IsPickedUpByCustomer {
		sendCourier = false
	}

	res := models.Order{
		RestaurantID:    req.RestaurantID,
		OrderID:         req.ID,
		StoreID:         req.RestaurantID,
		Type:            orderType,
		DeliveryService: models.KWAAKA_ADMIN.String(),
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
		PaymentMethod: "DELAYED",
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
			Email:               req.Customer.Email,
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
		SpecialRequirements:  req.DeliveryAddress.Comment,
		IsPickedUpByCustomer: true,
		Preorder:             preorder,
		PosPaymentInfo: models.PosPaymentInfo{
			PaymentTypeID:   req.PaymentType.ID,
			PaymentTypeKind: req.PaymentType.Kind,
		},
		IsMarketplace:               false,
		DeliveryDispatcher:          req.Delivery.Dispatcher,
		IsInstantDelivery:           store.Kwaaka3PL.IsInstantCall,
		DispatcherDeliveryTime:      req.Delivery.DeliveryTime,
		FullDeliveryPrice:           req.Delivery.FullDeliveryPrice,
		ClientDeliveryPrice:         req.Delivery.ClientDeliveryPrice,
		KwaakaChargedDeliveryPrice:  req.Delivery.KwaakaChargedDeliveryPrice,
		RestaurantPayDeliveryPrice:  s.getRestaurantPayDeliveryPrice(req.Delivery.FullDeliveryPrice, req.Delivery.ClientDeliveryPrice),
		DeliveryDropOffScheduleTime: req.Delivery.DropOffScheduleTime,
		RestaurantSelfDelivery:      true,
		SendCourier:                 sendCourier,
		OperatorName:                req.OperatorName,
	}

	res.Products = make([]models.OrderProduct, 0, len(req.Items))
	for _, product := range req.Items {
		res.Products = append(res.Products, s.toOrderProduct(product))
	}

	if req.Discount.Type != "" {
		res.Promos = append(res.Promos, models.Promo{
			Type:           req.Discount.Type,
			Discount:       req.Discount.Percent,
			IikoDiscountId: req.Discount.IikoDiscountId,
		})
	}

	log.Info().Msgf("finding operator id ....")
	if awsOperatorId, err := s.getOperatorIdFromCognito(req.OperatorName); err == nil {
		res.OperatorID = awsOperatorId
	} else {
		log.Trace().Msgf("Could not get identity of operator from username %s", err.Error())
		res.OperatorID = req.OperatorID
	}

	return res, nil
}

func (s *kwaakaAdminService) getRestaurantPayDeliveryPrice(fullDeliveryPrice float64, clientDeliveryPrice float64) float64 {
	diff := fullDeliveryPrice - clientDeliveryPrice

	if diff < 0 {
		return 0
	}

	return diff
}

func (s *kwaakaAdminService) getOperatorIdFromCognito(username string) (string, error) {
	log.Info().Msg("starting operation to get operator id from aws")

	input := cognitoidentityprovider.AdminGetUserInput{}

	input.SetUserPoolId(CognitoUserPoolId).SetUsername(username)

	output, err := s.cognitoService.AdminGetUser(&input)
	if err != nil {
		return "", err
	}

	for _, att := range output.UserAttributes {
		if *att.Name == OperatorId {
			log.Info().Msgf("operator id from aws cognito: %s", *att.Value)
			return *att.Value, nil
		}
	}

	return "", errors.New("cannot find operator id in aws")
}

func (s *kwaakaAdminService) getCookingTime(items []models2.Item, storeCookingTime int32) int32 {
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

func (s *kwaakaAdminService) toLabel(address models2.DeliveryAddress) string {
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

func (s *kwaakaAdminService) toOrderProduct(req models2.Item) models.OrderProduct {
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

func (s *kwaakaAdminService) MapSystemStatusToAggregatorStatus(order models.Order, posStatus models.PosStatus, store storeModels.Store) string {
	return posStatus.String()
}

func (s *kwaakaAdminService) UpdateOrderInAggregator(ctx context.Context, order models.Order, store storeModels.Store, aggregatorStatus string) error {
	return nil
}

func (s *kwaakaAdminService) UpdateStopListByProducts(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isAvailable bool) (string, error) {
	return "", nil
}

func (s *kwaakaAdminService) UpdateStopListByProductsBulk(ctx context.Context, aggregatorStoreID string, products []menuModels.Product, isSendRemains bool) (string, error) {
	return "", nil
}

func (s *kwaakaAdminService) UpdateStopListByAttributes(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute, isAvailable bool) (string, error) {
	return "", nil
}

func (s *kwaakaAdminService) UpdateStopListByAttributesBulk(ctx context.Context, aggregatorStoreID string, attributes []menuModels.Attribute) (string, error) {
	return "", nil
}

func (s *kwaakaAdminService) GetAggregatorOrder(ctx context.Context, orderID string) (models3.Order, error) {
	return models3.Order{}, nil
}

func (s *kwaakaAdminService) SendOrderErrorNotification(ctx context.Context, req interface{}) error {
	return nil
}

func (s *kwaakaAdminService) SendStopListUpdateNotification(ctx context.Context, aggregatorStoreID string) error {
	return nil
}
