package pos

import (
	"context"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"

	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	burgerKingClient "github.com/kwaaka-team/orders-core/pkg/burgerking"
	burgerKingConfig "github.com/kwaaka-team/orders-core/pkg/burgerking/clients"
	burgerKingModels "github.com/kwaaka-team/orders-core/pkg/burgerking/clients/models"
	"github.com/rs/zerolog/log"
)

type BurgerKingManager struct {
	bkClient     burgerKingConfig.BK
	bkOffersRepo drivers.BKOfferRepository
}

func NewBKManager(globalConfig config.Configuration, bkOffersRepo drivers.BKOfferRepository) (*BurgerKingManager, error) {
	// Initialize new Burger King client
	bkClient, err := burgerKingClient.NewBKClient(&burgerKingConfig.Config{
		Protocol: "http",
		Address:  globalConfig.BurgerKingConfiguration.BaseURL,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize BK Manager.")
		return nil, err
	}

	return &BurgerKingManager{
		bkClient:     bkClient,
		bkOffersRepo: bkOffersRepo,
	}, nil
}

func (manager BurgerKingManager) CancelOrder(ctx context.Context, order models.Order, cancelReason, paymentStrategy string, store coreStoreModels.Store) error {
	err := manager.bkClient.CancelOrder(ctx, burgerKingModels.CancelOrderRequest{
		OrderID:         order.OrderID,
		StoreID:         order.StoreID,
		CancelReason:    cancelReason,
		PaymentStrategy: paymentStrategy,
	})

	if err != nil {
		log.Err(err).Msg("Burger King cancel order error")
		return err
	}

	log.Err(err).Msg("Burger King cancel order success")

	return nil
}
func (manager BurgerKingManager) GetOrderStatus(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, error) {
	return "", nil
}

func (manager BurgerKingManager) constructPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	// Construct Burger King order from DB order model
	products := make([]burgerKingModels.Product, 0, len(order.Products))

	offers, err := manager.bkOffersRepo.GetActiveOffers(ctx)

	if err != nil {
		return nil, order, err
	}

	offersMap := make(map[string]int, len(offers))
	for _, offer := range offers {
		offersMap[offer.ProductID] = offer.FinalPrice
	}

	for _, p := range order.Products {
		attributes := make([]burgerKingModels.Attribute, 0, len(p.Attributes))

		if finalPrice, ok := offersMap[p.ID]; ok {
			p.Price.Value = float64(finalPrice)
		}

		for _, a := range p.Attributes {
			attributes = append(attributes, burgerKingModels.Attribute{
				ID:       a.ID,
				Quantity: a.Quantity,
				Price:    int(a.Price.Value),
				Name:     a.Name,
			})
		}

		products = append(products, burgerKingModels.Product{
			ID:                 p.ID,
			Name:               p.Name,
			Quantity:           p.Quantity,
			Price:              int(p.Price.Value),
			PurchasedProductID: p.PurchasedProductID,
			Attributes:         attributes,
		})
	}

	bkOrder := burgerKingModels.Order{
		StoreID:             order.StoreID,
		OrderID:             order.OrderID,
		OrderCode:           order.OrderCode,
		PickUpCode:          order.PickUpCode,
		UTCOffsetMinutes:    order.UtcOffsetMinutes,
		Products:            products,
		EstimatedPickupTime: order.EstimatedPickupTime.Value.String(),
	}
	return bkOrder, order, nil
}

func (manager BurgerKingManager) sendOrder(ctx context.Context, order any, store coreStoreModels.Store) (any, error) {
	// Send order to Burger King
	posOrder, ok := order.(burgerKingModels.Order)

	if !ok {
		return "", validator.ErrCastingPos
	}

	utils.Beautify("Burger king request body", posOrder)

	response, err := manager.bkClient.SendOrder(ctx, posOrder)

	if err != nil {
		utils.Beautify("Burger king Err response", err)
		return "", err
	}

	utils.Beautify("Burger king response", response)

	return response.Message, nil
}

func (manager BurgerKingManager) CreateOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, error) {
	posOrder, _, err := manager.constructPosOrder(ctx, order, store)

	if err != nil {
		return order, validator.ErrCastingPos
	}

	posOrder, ok := posOrder.(burgerKingModels.Order)

	if !ok {
		return order, validator.ErrCastingPos
	}

	response, err := manager.sendOrder(ctx, posOrder, coreStoreModels.Store{})

	message, ok := response.(string)
	if ok {
		order.CreationResult = models.CreationResult{
			Message: message,
		}
	}

	if err != nil {
		return order, err
	}

	return order, nil
}

func (manager BurgerKingManager) UpdateOrderProblem(ctx context.Context, organizationID, posOrderID string) error {
	return ErrUnsupportedMethod
}
