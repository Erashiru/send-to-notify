package pos

import (
	"context"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"strconv"

	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/custom"
	pkg "github.com/kwaaka-team/orders-core/pkg/jowi"
	jowiClient "github.com/kwaaka-team/orders-core/pkg/jowi/client"
	jowiDto "github.com/kwaaka-team/orders-core/pkg/jowi/client/dto"
	"github.com/rs/zerolog/log"
)

type JowiManager struct {
	jowiClient pkg.Jowi
}

func NewJOWIManager(globalConfig config.Configuration) (*JowiManager, error) {
	jowiCli, err := jowiClient.New(pkg.Config{
		ApiKey:    globalConfig.JowiConfiguration.ApiKey,
		ApiSecret: globalConfig.JowiConfiguration.ApiSecret,
		BaseURL:   globalConfig.JowiConfiguration.BaseURL,
		Protocol:  "http",
	})

	if err != nil {
		log.Trace().Err(err).Msg("Can not initialize Jowi Manager")
		return nil, err
	}

	return &JowiManager{
		jowiClient: jowiCli,
	}, nil
}

func (manager JowiManager) CancelOrder(ctx context.Context, order models.Order, cancelReason, paymentStrategy string, store coreStoreModels.Store) error {
	return errs.ErrUnsupportedMethod
}

func (manager JowiManager) GetOrderStatus(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, error) {
	return "", nil
}

func (manager JowiManager) sendOrder(ctx context.Context, order any, store coreStoreModels.Store) (any, error) {
	var errs custom.Error
	posOrder, ok := order.(jowiDto.RequestCreateOrder)

	if !ok {
		return "", validator.ErrCastingPos
	}

	log.Info().Msgf("Jowi Request Body: %+v", posOrder)

	createResponse, err := manager.jowiClient.CreateOrder(ctx, posOrder)
	if err != nil {
		log.Err(err).Msg("Jowi create order error")
		errs.Append(err, validator.ErrIgnoringPos)
		return "", errs
	}

	return createResponse, nil
}

func (manager JowiManager) constructPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	jowiOrder := jowiDto.RequestCreateOrder{
		RestaurantID: store.Token,
		Order: jowiDto.RequestCreateOrderBody{
			Address:     order.DeliveryAddress.Label,
			Phone:       order.Customer.PhoneNumber,
			Contact:     order.Customer.Name,
			AmountOrder: strconv.Itoa(int(order.EstimatedTotalPrice.Value)),
		},
	}

	if order.Persons != 0 {
		jowiOrder.Order.PeopleCount = order.Persons
	}

	var courses = make([]jowiDto.RequestCreateOrderCourse, 0)
	for _, product := range order.Products {
		courses = append(courses, jowiDto.RequestCreateOrderCourse{
			CourseId: product.ID,
			Count:    product.Quantity,
			Price:    int(product.Price.Value),
		})

		for _, attribute := range product.Attributes {
			courses = append(courses, jowiDto.RequestCreateOrderCourse{
				CourseId: attribute.ID,
				Count:    attribute.Quantity,
				Price:    int(attribute.Price.Value),
			})
		}
	}

	switch order.PaymentMethod {
	case "CASH":
		jowiOrder.Order.PaymentType = 0
		jowiOrder.Order.PaymentMethod = 0
	case "DELAYED":
		jowiOrder.Order.PaymentType = 1
		jowiOrder.Order.PaymentMethod = 1
	}

	var isMarketplace bool
	switch order.DeliveryService {
	case "glovo":
		isMarketplace = store.Glovo.IsMarketplace
	case "wolt":
		isMarketplace = store.Wolt.IsMarketplace
	case "chocofood":
		isMarketplace = store.Chocofood.IsMarketplace
	case "qr_menu":
		isMarketplace = store.QRMenu.IsMarketplace
	default:
		for _, deliveryService := range store.ExternalConfig {
			if deliveryService.Type == order.DeliveryService {
				isMarketplace = deliveryService.IsMarketplace
			}
		}
	}

	jowiOrder.Order.OrderType = 1
	if isMarketplace {
		jowiOrder.Order.OrderType = 0
	}

	jowiOrder.Order.Courses = courses

	return jowiOrder, order, nil
}

func (manager JowiManager) CreateOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, error) {
	posOrder, _, err := manager.constructPosOrder(ctx, order, store)

	if err != nil {
		return order, validator.ErrCastingPos
	}

	posOrder, ok := posOrder.(jowiDto.RequestCreateOrder)

	if !ok {
		return order, validator.ErrCastingPos
	}

	response, err := manager.sendOrder(ctx, posOrder, coreStoreModels.Store{})

	responseOrder, ok := response.(jowiDto.ResponseOrder)
	if ok {
		order.PosOrderID = responseOrder.Order.Id
		order.CreationResult = models.CreationResult{
			Message: responseOrder.Message,
			OrderInfo: models.OrderInfo{
				CreationStatus: strconv.Itoa(responseOrder.Order.Status),
				OrganizationID: responseOrder.Order.RestaurantId,
			},
			ErrorDescription: responseOrder.ErrorResponse.Message,
		}
	}

	if err != nil {
		return order, err
	}

	return order, nil
}

func (manager JowiManager) UpdateOrderProblem(ctx context.Context, organizationID, posOrderID string) error {
	return ErrUnsupportedMethod
}
