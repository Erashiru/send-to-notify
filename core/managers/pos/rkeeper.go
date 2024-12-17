package pos

import (
	"context"
	"fmt"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"strconv"
	"time"

	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/custom"
	MenuUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	rkeeperClient "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite"
	rkeeperConf "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients"
	rkeeperDto "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/dto"

	"github.com/rs/zerolog/log"
)

type RKeeperManager struct {
	rkeeperCli         rkeeperConf.RKeeper
	menu               coreMenuModels.Menu
	aggregatorMenu     coreMenuModels.Menu
	globalConfig       config.Configuration
	productsMap        map[string]coreMenuModels.Product
	attributesMap      map[string]coreMenuModels.Attribute
	attributeGroupsMap map[string]coreMenuModels.AttributeGroup
}

func NewRKeeperManager(globalConfig config.Configuration, menu coreMenuModels.Menu, aggregatorMenu coreMenuModels.Menu, store coreStoreModels.Store) (RKeeperManager, error) {
	var apyKey string
	switch {
	case store.RKeeper.ApiKey != "":
		apyKey = store.RKeeper.ApiKey
	default:
		apyKey = globalConfig.RKeeperApiKey
	}

	rkeeperClient, err := rkeeperClient.NewRKeeperClient(&rkeeperConf.Config{
		Protocol: "http",
		BaseURL:  globalConfig.RKeeperBaseURL,
		ApiKey:   apyKey,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize RKeeper Client.")
		return RKeeperManager{}, err
	}

	return RKeeperManager{
		rkeeperCli:     rkeeperClient,
		menu:           menu,
		aggregatorMenu: aggregatorMenu,
		globalConfig:   globalConfig,

		// Preprocess menu parts to "Map[Intance ID] = Instance"
		productsMap:        MenuUtils.ProductsMap(menu),
		attributesMap:      MenuUtils.AtributesMap(menu),
		attributeGroupsMap: MenuUtils.AtributeGroupsMap(menu),
	}, nil
}

func (manager RKeeperManager) GetOrderStatus(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, error) {
	return "", nil
}

func (manager RKeeperManager) CancelOrder(ctx context.Context, order models.Order, cancelReason, paymentStrategy string, store coreStoreModels.Store) error {
	return errs.ErrUnsupportedMethod
}

func (manager RKeeperManager) sendOrder(ctx context.Context, order any, store coreStoreModels.Store) (any, error) {
	var errs custom.Error
	posOrder, ok := order.(rkeeperDto.CreateOrderRequest)

	if !ok {
		return "", validator.ErrCastingPos
	}

	log.Info().Msgf("rkeeper Request Body: %+v", posOrder)

	createResponse, err := manager.rkeeperCli.CreateOrder(ctx, posOrder.Params.Async.ObjectID, posOrder.Params.Order)
	if err != nil {
		log.Err(err).Msg("rkeeper create order error")
		errs.Append(err, validator.ErrIgnoringPos)
		return "", errs
	}

	return createResponse, nil
}

func (manager RKeeperManager) constructPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	var (
		taker          = "courier"
		expiditionType = "delivery"
	)

	switch order.IsPickedUpByCustomer {
	case true:
		taker = "customer"
		expiditionType = "pickup"
	}

	aggregatorProducts, aggregatorAttributes := ActiveMenuPositions(ctx, manager.aggregatorMenu)
	orderComment, _, _ := ConstructOrderComments(ctx, order, store)

	rkeeperOrder := rkeeperDto.Order{
		OriginalOrderId: order.OrderID,
		Customer: &rkeeperDto.PersonInfo{
			Name:  order.Customer.Name,
			Phone: order.Customer.PhoneNumber,
		},
		Delivery: rkeeperDto.Delivery{
			ExpectedTime: order.EstimatedPickupTime.Value.Time.Add(time.Duration(store.Settings.TimeZone.UTCOffset) * time.Hour), // :TODO
			Address: &rkeeperDto.DeliveryAddress{
				FullAddress: order.DeliveryAddress.Label,
			},
		},
		ExpeditionType: expiditionType,
		Pickup: rkeeperDto.PickUp{
			Courier: &rkeeperDto.PersonInfo{
				Name:  order.Courier.Name,
				Phone: order.Courier.PhoneNumber,
			},
			ExpectedTime: order.EstimatedPickupTime.Value.Time.Add(time.Duration(store.Settings.TimeZone.UTCOffset) * time.Hour), // :TODO
			Taker:        taker,
		},
		Comment:         orderComment,
		PersonsQuantity: order.Persons,
	}

	if order.PosPaymentInfo.OrderType != "" {
		orderType, err := strconv.Atoi(order.PosPaymentInfo.OrderType)
		if err != nil {
			log.Err(err).Msgf("order type is incorrect to convert to int")
		} else {
			rkeeperOrder.OrderTypeCode = orderType
		}
	}

	if !store.ExternalPosIntegrationSettings.PayOrderIsOn {
		rkeeperOrder.Price = &rkeeperDto.Price{
			Total: int(order.EstimatedTotalPrice.Value),
		}
	}

	var products = make([]rkeeperDto.CreateOrderProduct, 0)
	for _, product := range order.Products {
		var ingredients = make([]rkeeperDto.CreateOrderIngredient, 0, len(product.Attributes))
		var ingridientsIds = make(map[string]string)

		productPosID, exist := aggregatorProducts[product.ID]
		if exist {
			product.ID = productPosID
		} else {
			log.Warn().Msgf("Product with ID %s %s not matched", product.ID, product.Name)
		}

		posProduct, ok := manager.productsMap[product.ID]
		if ok {
			product.ID = posProduct.ProductID
		}

		for _, attribute := range product.Attributes {
			attributePosID, exist_ := aggregatorAttributes[attribute.ID]
			if exist_ {
				attribute.ID = attributePosID
			} else {
				log.Warn().Msgf("Attribute with ID %s %s not matched", attribute.ID, attribute.Name)
			}

			// Actual attribute
			_, ok := manager.attributesMap[attribute.ID]
			if ok {
				ingridientsIds[attribute.ID] = attribute.ID
				ingredients = append(ingredients, rkeeperDto.CreateOrderIngredient{
					Id:       attribute.ID,
					Name:     attribute.Name,
					Quantity: attribute.Quantity,
				})
				continue
			}

			// Product as Attribute
			productAsAttr, ok := manager.productsMap[attribute.ID]
			if ok {
				products = append(products, rkeeperDto.CreateOrderProduct{
					Id:       productAsAttr.ProductID,
					Name:     attribute.Name,
					Quantity: attribute.Quantity * product.Quantity,
				})
			}
		}

		// Add default ingridients to product
		if posProduct.MenuDefaultAttributes != nil && len(posProduct.MenuDefaultAttributes) > 0 {
			for _, defaultAttribute := range posProduct.MenuDefaultAttributes {

				defaultAttributePosID, exist__ := aggregatorAttributes[defaultAttribute.ExtID]
				if exist__ {
					defaultAttribute.ExtID = defaultAttributePosID
				}

				_, ok_ := ingridientsIds[defaultAttribute.ExtID]

				if ok_ {
					continue
				}

				price := "0"
				amount := 1
				if defaultAttribute.DefaultAmount > 0 {
					amount = defaultAttribute.DefaultAmount
				}

				defaultIngridient := rkeeperDto.CreateOrderIngredient{
					Id:       defaultAttribute.ExtID,
					Quantity: amount,
					Price:    price,
				}

				ingridientsIds[defaultAttribute.ExtID] = defaultAttribute.ExtID
				ingredients = append(ingredients, defaultIngridient)
			}
		}

		products = append(products, rkeeperDto.CreateOrderProduct{
			Id:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			Ingredients: ingredients,
		})
	}

	//switch order.PaymentMethod {
	//case "CASH":
	//	rkeeperOrder.Payment.Type = "cash"
	//case "DELAYED":
	//	rkeeperOrder.Payment.Type = "online"
	//}

	rkeeperOrder.Products = products

	return rkeeperOrder, order, nil
}

func (manager RKeeperManager) CreateOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, error) {
	posOrder, _, err := manager.constructPosOrder(ctx, order, store)

	if err != nil {
		return order, validator.ErrCastingPos
	}

	rkeeperOrder, ok := posOrder.(rkeeperDto.Order)

	if !ok {
		return order, validator.ErrCastingPos
	}

	req := rkeeperDto.CreateOrderRequest{
		Params: rkeeperDto.CreateOrderParam{
			Async: rkeeperDto.Sync{
				ObjectID: store.RKeeper.ObjectId,
			},
			Order: rkeeperOrder,
		},
	}

	response, err := manager.sendOrder(ctx, req, coreStoreModels.Store{})
	if err != nil {
		return order, err
	}

	syncResponse, ok := response.(rkeeperDto.SyncResponse)
	if !ok {
		log.Warn().Msgf("Cant serialize response %v", response)
		return order, nil
	}

	// createOrderResult, err := manager.rkeeperCli.CreateOrderTask(ctx, syncResponse.ResponseCommon.TaskGUID)
	// if err != nil {
	// 	return order, err
	// }

	// utils.Beautify("Get creation result response", createOrderResult)

	// Fill Rkeeper order ID
	// order.PosOrderID = createOrderResult.TaskResponse.Order.OrderGuid
	order.CreationResult = models.CreationResult{
		Message: syncResponse.TaskResponse.Status,
		OrderInfo: models.OrderInfo{
			ID:             syncResponse.ResponseCommon.TaskGUID,
			OrganizationID: strconv.Itoa(syncResponse.ResponseCommon.ObjectID),
			CreationStatus: syncResponse.TaskResponse.Status,
		},
		ErrorDescription: syncResponse.Error.WsError.Desc + " " + syncResponse.Error.AgentError.Desc,
	}

	if syncResponse.Error.WsError.Desc != "" || syncResponse.Error.AgentError.Desc != "" {
		order.Status = "FAILED"
		return order, err
	}

	if store.ExternalPosIntegrationSettings.PayOrderIsOn {
		time.Sleep(120 * time.Second)

		createOrderTask, err := manager.rkeeperCli.CreateOrderTask(ctx, syncResponse.ResponseCommon.TaskGUID)
		if err != nil {
			return order, err
		}

		if createOrderTask.TaskResponse.Order.OrderGuid == "" {
			order.CreationResult.ErrorDescription = createOrderTask.Error.WsError.Desc + " " + createOrderTask.Error.AgentError.Desc
			return order, fmt.Errorf(createOrderTask.Error.WsError.Desc + " " + createOrderTask.Error.AgentError.Desc)
		}

		time.Sleep(2 * time.Second)

		payOrder, err := manager.rkeeperCli.PayOrder(ctx, store.RKeeper.ObjectId, int(order.EstimatedTotalPrice.Value), createOrderTask.TaskResponse.Order.OrderGuid, order.PosPaymentInfo.PaymentTypeID)
		if err != nil {
			return order, err
		}

		utils.Beautify("pay order result", payOrder)
	}

	return order, nil
}

func (manager RKeeperManager) UpdateOrderProblem(ctx context.Context, organizationID, posOrderID string) error {
	return ErrUnsupportedMethod
}
