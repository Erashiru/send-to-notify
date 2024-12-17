package pos

import (
	//MenuModels "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/config"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	MenuUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	palomaClient "github.com/kwaaka-team/orders-core/pkg/paloma"
	palomaConf "github.com/kwaaka-team/orders-core/pkg/paloma/clients"
	palomaModels "github.com/kwaaka-team/orders-core/pkg/paloma/clients/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"
)

type PalomaManager struct {
	palomaCli          palomaConf.Paloma
	menu               coreMenuModels.Menu
	aggregatorMenu     coreMenuModels.Menu
	globalConfig       config.Configuration
	productsMap        map[string]coreMenuModels.Product
	attributesMap      map[string]coreMenuModels.Attribute
	attributeGroupsMap map[string]coreMenuModels.AttributeGroup
}

func (manager PalomaManager) CancelOrder(ctx context.Context, order models.Order, cancelReason, paymentStrategy string, store coreStoreModels.Store) error {
	return errs.ErrUnsupportedMethod
}

func (manager PalomaManager) sendOrder(ctx context.Context, order any, store coreStoreModels.Store) (any, error) {
	//var errs custom.Error
	posOrder, ok := order.(palomaModels.Order)

	if !ok {
		return "", validator.ErrCastingPos
	}

	utils.Beautify("paloma request body", posOrder)

	createResponse, err := manager.palomaCli.CreateOrder(ctx, store.Paloma.PointID, posOrder)
	if err != nil {
		log.Err(err).Msg("paloma create order error in PalomaManager")
		//errs.Append(err, validator.ErrIgnoringPos)
		return "", err
	}

	return createResponse, nil
}

func (manager PalomaManager) constructPosOrder(ctx context.Context, req models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	aggregatorProducts, aggregatorAttributes := ActiveMenuPositions(ctx, manager.aggregatorMenu)

	orderComments, _, _ := ConstructOrderComments(ctx, req, store)

	order := palomaModels.Order{
		OrderId:        req.ID,
		Date:           req.EstimatedPickupTime.Value.Time.Add(time.Duration(store.Settings.TimeZone.UTCOffset) * time.Hour).String(),
		Name:           req.Customer.Name,
		Phone:          "+77777777777",
		Address:        req.DeliveryAddress.Label,
		CoordinateLong: strconv.Itoa(int(req.DeliveryAddress.Longitude)),
		CoordinateLat:  strconv.Itoa(int(req.DeliveryAddress.Latitude)),
		Comment:        orderComments,
		PersonAmount:   req.Persons,
		TotalPrice:     int(req.EstimatedTotalPrice.Value),
		DiscountAmount: 0,
	}

	switch req.PaymentMethod {
	case "CASH":
		order.IsCash = true
		order.IsPayed = false
	case "DELAYED":
		order.IsCash = false
		order.IsPayed = true
	default:
		order.IsCash = false
		order.IsPayed = true
	}

	switch req.IsPickedUpByCustomer {
	case false:
		order.DeliveryType = 1
	default:
		order.DeliveryType = 2
	}

	var items = make([]palomaModels.OrderItem, 0, len(req.Products))

	for _, product := range req.Products {
		productPosID, exist := aggregatorProducts[product.ID]
		switch exist {
		case true:
			product.ID = productPosID
		default:
			log.Warn().Msgf("product with ID %s %s not matched", product.ID, product.Name)
			return nil, req, fmt.Errorf("product with id=%s, name=%s is not matched", product.ID, product.Name)
		}

		posProduct, ok := manager.productsMap[product.ID]
		if !ok {
			return nil, req, errors.Wrap(errs.ErrProductNotFound, fmt.Sprintf("PRODUCT NOT FOUND IN POS MENU, ID %s, NAME %s", product.ID, product.Name))
		}

		product.ID = posProduct.ProductID

		productID, err := strconv.Atoi(product.ID)
		if err != nil {
			// TODO: err?
			return nil, req, err
		}

		item := palomaModels.OrderItem{
			ObjectId: productID,
			Name:     product.Name,
			Price:    int(product.Price.Value),
			Count:    product.Quantity,
		}

		var itemPriceForCombo int

		for _, defaultAttribute := range posProduct.MenuDefaultAttributes {
			attributePosID, hasAttribute := aggregatorAttributes[defaultAttribute.ExtID]
			switch hasAttribute {
			case true:
				defaultAttribute.ExtID = attributePosID
			default:
				log.Warn().Msgf("default attribute with ID %s %s not matched", defaultAttribute.ExtID, defaultAttribute.Name)
				return nil, req, fmt.Errorf("default attribute with id=%s, name=%s is not matched", defaultAttribute.ExtID, defaultAttribute.Name)
			}

			attributeID, err := strconv.Atoi(defaultAttribute.ExtID)
			if err != nil {
				return nil, req, err
			}

			if defaultAttribute.DefaultAmount == 0 {
				defaultAttribute.DefaultAmount = 1
			}

			switch posProduct.IsCombo {
			case true:
				item.ComplexItems = append(item.ComplexItems, palomaModels.OrderComplexItem{
					ObjectId: attributeID,
					Name:     defaultAttribute.Name,
					Count:    defaultAttribute.DefaultAmount * product.Quantity,
					Price:    0, // TODO: price is 0?
				})
			default:
				item.Modifications = append(item.Modifications, palomaModels.OrderModification{
					ObjectId: attributeID,
					Name:     defaultAttribute.Name,
					Count:    defaultAttribute.DefaultAmount * product.Quantity,
					Price:    0, // TODO: price is 0?
				})
			}
		}

		for _, attribute := range product.Attributes {
			attributePosID, hasAttribute := aggregatorAttributes[attribute.ID]
			switch hasAttribute {
			case true:
				attribute.ID = attributePosID
			default:
				log.Warn().Msgf("attribute with ID %s %s not matched", attribute.ID, attribute.Name)
				return nil, req, fmt.Errorf("attribute with id=%s, name=%s is not matched", attribute.ID, attribute.Name)
			}

			// Actual attribute
			_, hasAttributeInPos := manager.attributesMap[attribute.ID]
			if !hasAttributeInPos {
				return nil, req, fmt.Errorf("ATTRIBUTE NOT FOUND IN POS MENU, ID %s, NAME %s", attribute.ID, attribute.Name)
			}

			attributeID, err := strconv.Atoi(attribute.ID)
			if err != nil {
				// TODO: err?
				return nil, req, err
			}

			switch posProduct.IsCombo {
			case true:
				item.ComplexItems = append(item.ComplexItems, palomaModels.OrderComplexItem{
					ObjectId: attributeID,
					Name:     attribute.Name,
					Count:    attribute.Quantity,
					Price:    int(attribute.Price.Value),
				})
				itemPriceForCombo += int(attribute.Price.Value) * attribute.Quantity
			default:
				item.Modifications = append(item.Modifications, palomaModels.OrderModification{
					ObjectId: attributeID,
					Name:     attribute.Name,
					Count:    attribute.Quantity,
					Price:    int(attribute.Price.Value),
				})
			}
		}

		if posProduct.IsCombo {
			item.Price = item.Price + itemPriceForCombo
		}

		items = append(items, item)
	}

	order.OrderItems = items

	return order, req, nil
}

func (manager PalomaManager) CreateOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, error) {
	posOrder, _, err := manager.constructPosOrder(ctx, order, store)

	if err != nil {
		return order, err
	}

	palomaOrder, ok := posOrder.(palomaModels.Order)

	if !ok {
		return order, validator.ErrCastingPos
	}

	response, err := manager.sendOrder(ctx, palomaOrder, store)
	if err != nil {
		return order, err
	}

	orderResponse, ok := response.(palomaModels.OrderResponse)
	if !ok {
		log.Warn().Msgf("Cant serialize response %v", response)
		return order, nil
	}

	order.CreationResult = models.CreationResult{
		Message: orderResponse.Status,
		OrderInfo: models.OrderInfo{
			ID:             orderResponse.OrderId,
			OrganizationID: store.Paloma.PointID,
			CreationStatus: orderResponse.Status,
		},
	}

	return order, nil
}

func (manager PalomaManager) GetOrderStatus(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, error) {
	response, err := manager.palomaCli.GetOrderStatus(ctx, order.ID)
	if err != nil {
		return "", err
	}

	return response.Status, nil
}

func NewPalomaManager(globalConfig config.Configuration, menu coreMenuModels.Menu, aggregatorMenu coreMenuModels.Menu, store coreStoreModels.Store) (PalomaManager, error) {
	palomaClient, err := palomaClient.New(&palomaConf.Config{
		Protocol: "http",
		BaseURL:  globalConfig.PalomaConfiguration.BaseURL,
		Class:    globalConfig.PalomaConfiguration.Class,
		ApiKey:   store.Paloma.ApiKey,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Paloma Client.")
		return PalomaManager{}, err
	}

	return PalomaManager{
		palomaCli:      palomaClient,
		menu:           menu,
		aggregatorMenu: aggregatorMenu,
		globalConfig:   globalConfig,

		// Preprocess menu parts to "Map[Intance ID] = Instance"
		productsMap:        MenuUtils.ProductsMap(menu),
		attributesMap:      MenuUtils.AtributesMap(menu),
		attributeGroupsMap: MenuUtils.AtributeGroupsMap(menu),
	}, nil
}

func (manager PalomaManager) UpdateOrderProblem(ctx context.Context, organizationID, posOrderID string) error {
	return ErrUnsupportedMethod
}
