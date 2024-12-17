package pos

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/domain/foodband"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	posIntegrationClient "github.com/kwaaka-team/orders-core/pkg/posintegration"
	posIntegrationConf "github.com/kwaaka-team/orders-core/pkg/posintegration/clients"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type foodbandService struct {
	*BasePosService
	posIntegrationCli posIntegrationConf.FOODBAND
	storeID           string
}

func newFoodbandService(bps *BasePosService, createOrderUrl, cancelOrderUrl, apiToken, storeID string, retryMaxCount int) (*foodbandService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "foodbandService constructor error")
	}

	posIntegrationCli, err := posIntegrationClient.NewClient(&posIntegrationConf.Config{
		CreateOrderUrl: createOrderUrl,
		CancelOrderUrl: cancelOrderUrl,
		ApiToken:       apiToken,
		RetryMaxCount:  retryMaxCount,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize PosIntegration Client")
		return nil, err
	}

	return &foodbandService{
		bps,
		posIntegrationCli,
		storeID,
	}, nil
}

func (s *foodbandService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {
	return coreMenuModels.StopListItems{}, ErrUnsupportedMethod
}

func (s *foodbandService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	return coreMenuModels.Menu{}, ErrUnsupportedMethod
}

func (s *foodbandService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {
	previousStatusPriority, _ := getFoodBandStatusPriority(currentSystemStatus)
	externalStatusPriority, status := getFoodBandStatusPriority(posStatus)

	if externalStatusPriority == 0 {
		return 0, models.StatusIsNotExist
	}

	if previousStatusPriority > externalStatusPriority {
		return status, models.InvalidStatusPriority
	}
	return status, nil
}

func getFoodBandStatusPriority(status string) (int, models.PosStatus) {
	switch status {
	case "ACCEPTED":
		return 1, models.ACCEPTED
	case "COOKING_STARTED":
		return 2, models.COOKING_STARTED
	case "COOKING_COMPLETE":
		return 3, models.COOKING_COMPLETE
	case "READY_FOR_PICKUP":
		return 4, models.READY_FOR_PICKUP
	case "OUT_FOR_DELIVERY":
		return 5, models.OUT_FOR_DELIVERY
	case "PICKED_UP_BY_CUSTOMER":
		return 6, models.PICKED_UP_BY_CUSTOMER
	case "DELIVERED":
		return 6, models.DELIVERED
	case "CLOSED":
		return 7, models.CLOSED
	default:
		return 0, 0
	}
}

func (s *foodbandService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration, store coreStoreModels.Store,
	menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	order, err := prepareAnOrder(ctx, order, store, menu, aggregatorMenu, menuClient)
	if err != nil {
		log.Trace().Err(err).Msg("")
		return order, err
	}

	posOrder, _, err := s.toPosOrder(ctx, order, store)
	if err != nil {
		log.Trace().Err(validator.ErrCastingPos).Msg("")
		return order, err
	}

	utils.Beautify("FOODBAND Request Body", posOrder)

	order, err = s.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	retryCount, err := s.posIntegrationCli.CreateOrder(ctx, posOrder)
	if err != nil {
		return models.Order{
			PosOrderID: posOrder.Order.ID,
			CreationResult: models.CreationResult{
				OrderInfo: models.OrderInfo{
					ID:             posOrder.Order.ID,
					OrganizationID: posOrder.StoreID,
					CreationStatus: "ERROR",
				},
				ErrorDescription: err.Error(),
			},
			RetryCount: retryCount,
			IsRetry:    retryCount != 0,
		}, errors.Wrap(err, "foodband create order error")
	}

	order = setPosOrderId(order, posOrder.Order.ID)

	order.CreationResult = models.CreationResult{
		OrderInfo: models.OrderInfo{
			ID:             posOrder.Order.ID,
			OrganizationID: posOrder.StoreID,
		},
	}
	order.RetryCount = retryCount
	order.IsRetry = retryCount != 0

	utils.Beautify("finished order model result", order)

	return order, nil
}

func (s *foodbandService) toPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (foodband.CreateOrderRequest, models.Order, error) {
	orderComment := s.constructOrderComment(ctx, order, store)
	products, err := s.toPosProducts(order)
	if err != nil {
		return foodband.CreateOrderRequest{}, order, err
	}
	payment, err := s.toPosPayment(order)
	if err != nil {
		return foodband.CreateOrderRequest{}, order, err
	}
	deliveryProviderType := s.toPosDeliveryProviderType(order)
	orderFoodBand := foodband.Order{
		ID:              uuid.New().String(),
		Type:            order.Type,
		Code:            order.OrderCode,
		PickUpCode:      order.PickUpCode,
		CompleteBefore:  s.toPosCompleteBeforeDate(store, order),
		Phone:           s.getPhone(order.Customer.PhoneNumber),
		DeliveryService: order.DeliveryService,
		DeliveryPoint:   s.getDeliveryPoint(order.DeliveryAddress, deliveryProviderType),
		Comment:         orderComment,
		Customer: foodband.Customer{
			Name: order.Customer.Name,
		},
		Courier: foodband.Courier{
			Name:        order.Courier.Name,
			PhoneNumber: order.Courier.PhoneNumber,
		},
		Products: products,
		Payments: []foodband.Payment{payment},
		//DeliveryFee:          order.DeliveryFee.Value,
		DeliveryFee:          0,
		DeliveryProviderType: deliveryProviderType,
	}

	return foodband.CreateOrderRequest{
		StoreID: s.storeID,
		Order:   orderFoodBand,
	}, order, nil
}

func (s *foodbandService) toPosDeliveryProviderType(order models.Order) string {
	if order.IsPickedUpByCustomer {
		return models.FOODBAND_CUSTOMER_PICKUP
	}
	if order.RestaurantSelfDelivery {
		return models.FOODBAND_DELIVERY_RESTAURANT
	}
	return models.FOODBAND_DELIVERY_AGGREGATOR
}

func (s *foodbandService) toPosProducts(order models.Order) ([]foodband.Product, error) {
	var items []foodband.Product

	for _, product := range order.Products {
		var modifiers []foodband.Attribute

		for _, attribute := range product.Attributes {
			itemModifier := foodband.Attribute{
				ID:       attribute.ID,
				Quantity: attribute.Quantity,
				Price:    attribute.Price.Value,
			}

			modifiers = append(modifiers, itemModifier)
		}
		orderItem := foodband.Product{
			ID:         product.ID,
			Price:      product.Price.Value,
			Quantity:   product.Quantity,
			Attributes: modifiers,
		}

		items = append(items, orderItem)
	}

	return items, nil
}

func (s *foodbandService) constructOrderComment(ctx context.Context, order models.Order, store coreStoreModels.Store) string {
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

func (s *foodbandService) toPosPayment(order models.Order) (foodband.Payment, error) {
	var paymentType string
	switch order.PaymentMethod {
	case models.PAYMENT_METHOD_DELAYED:
		paymentType = models.PAYMENT_METHOD_CARD
	case models.PAYMENT_METHOD_CASH:
		paymentType = models.PAYMENT_METHOD_CASH
	default:
		log.Info().Msgf("Payment Method: %v", order.PaymentMethod)
		return foodband.Payment{}, fmt.Errorf("invalid payment type")
	}

	return foodband.Payment{
		PaymentTypeKind: paymentType,
		Sum:             order.EstimatedTotalPrice.Value - order.PartnerDiscountsProducts.Value,
	}, nil
}

func (s *foodbandService) toPosCompleteBeforeDate(store coreStoreModels.Store, order models.Order) string {
	completeBeforeDate := order.EstimatedPickupTime.Value.Time
	if completeBeforeDate.IsZero() {
		completeBeforeDate = time.Now().UTC().Add(time.Hour)
	}
	completeBeforeDate = completeBeforeDate.Add(time.Duration(store.Settings.TimeZone.UTCOffset)*time.Hour - 1*time.Minute)

	completeBefore := completeBeforeDate.Format("2006-01-02 15:04:05.000")

	return completeBefore
}

func (s *foodbandService) getPhone(customerPhone string) string {
	if !strings.Contains(customerPhone, "+") {
		customerPhone = "+77771111111"
	}
	return customerPhone
}

func (s *foodbandService) getDeliveryPoint(deliveryAddress models.DeliveryAddress, deliveryProviderType string) foodband.DeliveryPoint {
	if deliveryProviderType == models.FOODBAND_DELIVERY_RESTAURANT {
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

func (s *foodbandService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	err := s.posIntegrationCli.CancelOrder(ctx, foodband.CancelOrderRequest{
		StoreID:         store.ID,
		OrderID:         order.PosOrderID,
		DeliveryService: order.DeliveryService,
		CancelReason:    order.CancelReason.Reason,
		PaymentStrategy: order.PaymentStrategy,
	})
	if err != nil {
		log.Err(err).Msg("FOODBAND cancel order error")
		return err
	}

	return nil
}

func (s *foodbandService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (s *foodbandService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (s *foodbandService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
