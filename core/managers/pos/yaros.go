package pos

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	"github.com/kwaaka-team/orders-core/service/pos"

	yarosClient "github.com/kwaaka-team/orders-core/pkg/yaros"
	yarosConf "github.com/kwaaka-team/orders-core/pkg/yaros/clients"
	yarosModels "github.com/kwaaka-team/orders-core/pkg/yaros/models"

	"github.com/rs/zerolog/log"
	"strconv"
)

type YarosManager struct {
	yarosCli       yarosConf.Yaros
	menu           coreMenuModels.Menu
	aggregatorMenu coreMenuModels.Menu
	globalConfig   config.Configuration
	sqsCli         notifyQueue.SQSInterface
}

func NewYarosManager(globalConfig config.Configuration, menu coreMenuModels.Menu, aggregatorMenu coreMenuModels.Menu, store coreStoreModels.Store) (YarosManager, error) {
	client, err := yarosClient.NewClient(&yarosConf.Config{
		Protocol: "http",
		BaseURL:  globalConfig.YarosConfiguration.BaseURL,
		Username: store.Yaros.Username,
		Password: store.Yaros.Password,
		RestID:   store.Yaros.StoreId,
	})
	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Yaros Client.")
		return YarosManager{}, err
	}

	sqsCli := notifyQueue.NewSQS(sqs.NewFromConfig(globalConfig.AwsConfig))

	return YarosManager{
		yarosCli:       client,
		menu:           menu,
		aggregatorMenu: aggregatorMenu,
		globalConfig:   globalConfig,
		sqsCli:         sqsCli,
	}, nil
}

func (manager YarosManager) sendOrder(ctx context.Context, order any, store coreStoreModels.Store) (any, error) {
	var errs errors.Error

	posOrder, ok := order.(yarosModels.OrderRequest)
	if !ok {
		return "", validator.ErrCastingPos
	}

	utils.Beautify("Yaros Request Body", posOrder)

	createResponse, err := manager.yarosCli.CreateOrder(ctx, store.Yaros.StoreId, posOrder)
	if err != nil {
		log.Err(err).Msg("yaros error")
		errs.Append(err, validator.ErrIgnoringPos)
		return "", errs
	}
	return createResponse, nil
}

func (manager YarosManager) constructPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	timestamp := order.OrderTime.Value.Unix()

	var payMethod string
	switch order.PaymentMethod {
	case "DELAYED":
		payMethod = "visa"
	case "CASH":
		payMethod = "cash"
	default:
		//log error or add other pay methods

		return nil, models.Order{}, nil
	}

	var yarosProducts = make([]yarosModels.OrderItem, 0, len(order.Products))

	for _, product := range order.Products {
		yarosProduct := yarosModels.OrderItem{
			ProductId: product.ID,
			Quantity:  strconv.Itoa(product.Quantity),
			Price:     strconv.Itoa(int(product.Price.Value)),
			Amount:    strconv.Itoa(product.Quantity * int(product.Price.Value)),
		}
		yarosProducts = append(yarosProducts, yarosProduct)
	}

	yarosOrder := yarosModels.OrderRequest{
		Orders: []yarosModels.PosOrder{
			{
				Id:         order.OrderID,
				Type:       "delivery",
				InfoSystem: order.DeliveryService,
				Date:       strconv.Itoa(int(timestamp)),
				Change:     strconv.Itoa(int(order.MinimumBasketSurcharge.Value)),
				Total:      strconv.Itoa(int(order.EstimatedTotalPrice.Value)),
				Status:     "created", //надо узнать
				User: yarosModels.OrderUser{
					Name:  order.Customer.Name,
					Phone: order.Customer.PhoneNumber,
				},
				Address:   order.DeliveryAddress.Label,
				Comment:   order.AllergyInfo,
				PayMethod: payMethod,
				Items:     yarosProducts,
			},
		},
	}

	if store.Yaros.Department != "" && len(yarosOrder.Orders) > 0 {
		yarosOrder.Orders[0].Department = store.Yaros.Department
	}

	return yarosOrder, order, nil
}

func (manager YarosManager) CreateOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, error) {
	posOrder, _, err := manager.constructPosOrder(ctx, order, store)

	if err != nil {
		return order, err
	}

	response, err := manager.sendOrder(ctx, posOrder, store)
	if err != nil {
		log.Err(err).Msgf("couldn't create order in YAROS pos (retry)")
		log.Info().Msgf("yaros error case: run RETRY (retry)")
		if err = manager.orderRetry(ctx, order); err != nil {
			log.Err(err).Msgf("(CreateOrder - yarosService) error")
			return order, err
		}
		return order, pos.ErrRetry
	}

	responseOrder, ok := response.(yarosModels.OrderResponse)
	if ok {
		order.CreationResult = models.CreationResult{
			Message: responseOrder.Message,
			OrderInfo: models.OrderInfo{
				CreationStatus: responseOrder.Status,
			},
		}
	}

	return order, nil
}

func (manager YarosManager) CancelOrder(ctx context.Context, order models.Order, cancelReason, paymentStrategy string, store coreStoreModels.Store) error {
	return errors.ErrUnsupportedMethod
}

func (manager YarosManager) GetOrderStatus(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, error) {
	return "", errors.ErrUnsupportedMethod
}

func (manager YarosManager) orderRetry(ctx context.Context, order models.Order) error {
	if err := manager.sqsCli.SendSQSMessage(ctx, manager.globalConfig.RetryConfiguration.QueueName, fmt.Sprintf("yaros_%s", order.OrderID)); err != nil {
		log.Trace().Err(err).Msgf("SendSQSMessage creation-timeout error: %s", order.OrderID)
		return err
	}
	return nil
}

func (manager YarosManager) UpdateOrderProblem(ctx context.Context, organizationID, posOrderID string) error {
	return ErrUnsupportedMethod
}
