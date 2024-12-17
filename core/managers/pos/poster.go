package pos

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/config"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/custom"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	posterClient "github.com/kwaaka-team/orders-core/pkg/poster"
	posterConf "github.com/kwaaka-team/orders-core/pkg/poster/clients"
	posterModels "github.com/kwaaka-team/orders-core/pkg/poster/clients/models"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type PosterManager struct {
	posterCli      posterConf.Poster
	menu           coreMenuModels.Menu
	aggregatorMenu coreMenuModels.Menu
	globalConfig   config.Configuration
}

func (manager PosterManager) CancelOrder(ctx context.Context, order models.Order, cancelReason, paymentStrategy string, store coreStoreModels.Store) error {
	return errs.ErrUnsupportedMethod
}

func (p PosterManager) sendOrder(ctx context.Context, order any, store coreStoreModels.Store) (any, error) {
	var errs custom.Error

	posOrder, ok := order.(posterModels.CreateOrderRequest)

	if !ok {
		return "", validator.ErrCastingPos
	}

	utils.Beautify("Poster Request Body", posOrder)

	createResponse, err := p.posterCli.CreateOrder(ctx, posOrder)
	if err != nil {
		log.Err(err).Msg("poster error")
		errs.Append(err, validator.ErrIgnoringPos)
		return "", errs
	}
	return createResponse, nil
}

func (p PosterManager) constructPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	layout := "2006-01-02 15:04:05"
	orderComment, _, _ := ConstructOrderComments(ctx, order, store)
	utcOffset := time.Duration(int64(store.Settings.TimeZone.UTCOffset)) * time.Hour
	deliveryTime := order.EstimatedPickupTime.Value.Add(utcOffset)

	posterOrder := posterModels.CreateOrderRequest{
		SpotID: store.Poster.SpotId,
		Phone:  order.Customer.PhoneNumber,
		ClientAddress: posterModels.CreateOrderAddressRequest{
			Address1: order.DeliveryAddress.Label,
		},
		Comment:      orderComment,
		DeliveryTime: deliveryTime.Format(layout),
	}

	if order.DeliveryAddress.Latitude != 0 && order.DeliveryAddress.Longitude != 0 {
		posterOrder.ClientAddress.Longitude = fmt.Sprintf("%v", order.DeliveryAddress.Longitude)
		posterOrder.ClientAddress.Latitude = fmt.Sprintf("%v", order.DeliveryAddress.Latitude)
	}

	// fills only for prepaid orders
	if order.PaymentMethod == "DELAYED" {
		posterOrder.Payment = &posterModels.CreateOrderPaymentRequest{
			Type:     1,
			Sum:      int(order.EstimatedTotalPrice.Value * 100),
			Currency: store.Settings.Currency,
		}
	}

	if store.Poster.IgnorePaymentType {
		posterOrder.Payment = &posterModels.CreateOrderPaymentRequest{
			Type:     int(store.Poster.PaymentType),
			Sum:      int(order.EstimatedTotalPrice.Value * 100),
			Currency: store.Settings.Currency,
		}
	}

	switch order.IsPickedUpByCustomer {
	case true:
		posterOrder.ServiceMode = 2
	case false:
		posterOrder.ServiceMode = 3
		posterOrder.DeliveryPrice = int(order.DeliveryFee.Value * 100)
	}

	var posterProducts = make([]posterModels.CreateOrderProductRequest, 0)

	for _, product := range order.Products {
		productID, err := strconv.Atoi(product.ID)
		if err != nil {
			return nil, models.Order{}, err
		}
		attributesPrice := 0.0
		productAttributes := make([]posterModels.CreateOrderModificationRequest, 0, len(product.Attributes))
		for _, attribute := range product.Attributes {
			attributeId, err := strconv.Atoi(attribute.ID)
			if err != nil {
				return nil, models.Order{}, err
			}
			posterAttribute := posterModels.CreateOrderModificationRequest{
				M: attributeId,
				A: attribute.Quantity,
			}
			productAttributes = append(productAttributes, posterAttribute)
			attributesPrice = attributesPrice + attribute.Price.Value*float64(attribute.Quantity)
		}
		posterProduct := posterModels.CreateOrderProductRequest{
			ProductID:     productID,
			ModificatorID: 0,
			Modifications: productAttributes,
			Count:         strconv.Itoa(product.Quantity),
			Price:         int(product.Price.Value+attributesPrice) * 100,
		}
		posterProducts = append(posterProducts, posterProduct)
	}

	posterOrder.Products = posterProducts
	return posterOrder, order, nil
}

func (p PosterManager) CreateOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, error) {
	posOrder, _, err := p.constructPosOrder(ctx, order, store)
	if err != nil {
		return order, validator.ErrCastingPos
	}

	posOrder, ok := posOrder.(posterModels.CreateOrderRequest)
	if !ok {
		return order, validator.ErrCastingPos
	}

	response, err := p.sendOrder(ctx, posOrder, coreStoreModels.Store{})
	responseOrder, ok := response.(posterModels.CreateOrderResponse)
	if ok {
		order.PosOrderID = store.Poster.AccountNumberString + strconv.Itoa(responseOrder.Response.IncomingOrderID)
		order.CreationResult = models.CreationResult{
			Message: responseOrder.Message,
			OrderInfo: models.OrderInfo{
				CreationStatus: strconv.Itoa(responseOrder.Response.Status),
				OrganizationID: strconv.Itoa(responseOrder.Response.SpotId),
			},
			ErrorDescription: responseOrder.ErrorResponse.Message,
		}
	}
	if err != nil {
		return order, err
	}

	return order, nil
}

func (p PosterManager) GetOrderStatus(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, error) {
	return "", ErrUnsupportedMethod
}

func NewPosterManager(globalConfig config.Configuration, menu coreMenuModels.Menu, aggregatorMenu coreMenuModels.Menu, store coreStoreModels.Store) (PosterManager, error) {
	client, err := posterClient.NewClient(&posterConf.Config{
		Protocol: "http",
		BaseURL:  globalConfig.PosterConfiguration.BaseURL,
		Token:    store.Poster.Token,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize Poster Client.")
		return PosterManager{}, err
	}

	return PosterManager{
		posterCli:      client,
		menu:           menu,
		aggregatorMenu: aggregatorMenu,
		globalConfig:   globalConfig,
	}, nil
}

func (p PosterManager) UpdateOrderProblem(ctx context.Context, organizationID, posOrderID string) error {
	return ErrUnsupportedMethod
}
