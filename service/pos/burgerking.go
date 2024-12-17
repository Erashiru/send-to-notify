package pos

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/database/drivers"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	burgerKingClient "github.com/kwaaka-team/orders-core/pkg/burgerking"
	burgerKingConfig "github.com/kwaaka-team/orders-core/pkg/burgerking/clients"
	burgerKingModels "github.com/kwaaka-team/orders-core/pkg/burgerking/clients/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type burgerKingService struct {
	*BasePosService
	bkClient          burgerKingConfig.BK
	bkOfferRepository drivers.BKOfferRepository
}

func newBurgerKingService(bps *BasePosService, address string, bkOfferRepository drivers.BKOfferRepository) (*burgerKingService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "burgerKingService constructor error")
	}

	if bkOfferRepository == nil {
		return nil, errors.Wrap(constructorError, "bkOfferRepository is nil")
	}

	client, err := burgerKingClient.NewBKClient(&burgerKingConfig.Config{
		Protocol: "http",
		Address:  address,
	})

	if err != nil {
		log.Trace().Err(err).Msg("Cant initialize BurgerKing Client.")
		return nil, err
	}

	return &burgerKingService{
		bkClient:          client,
		BasePosService:    bps,
		bkOfferRepository: bkOfferRepository,
	}, nil
}

func (s *burgerKingService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {

	switch posStatus {
	case models.ACCEPTED.String(), models.COOKING_STARTED.String(), models.WAIT_SENDING.String():
		return models.ACCEPTED, nil
	case models.READY_FOR_PICKUP.String(), models.COOKING_COMPLETE.String(), models.CLOSED.String():
		return models.READY_FOR_PICKUP, nil
	case models.OUT_FOR_DELIVERY.String():
		return models.OUT_FOR_DELIVERY, nil
	case models.PICKED_UP_BY_CUSTOMER.String():
		return models.PICKED_UP_BY_CUSTOMER, nil
	}
	return 0, models.StatusIsNotExist
}

func (s *burgerKingService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration, store coreStoreModels.Store,
	menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	posOrder, _, err := s.constructPosOrder(ctx, order)

	if err != nil {
		return order, err
	}

	order, err = s.SetPosRequestBodyToOrder(order, posOrder)
	if err != nil {
		return order, err
	}

	response, err := s.sendOrder(ctx, posOrder)

	if err != nil {
		order.CreationResult.OrderInfo.CreationStatus = "FAILED"
		order.FailReason.Code = BK_ORDER_FAILED
		order.FailReason.Message = response.Message
		return order, err
	}

	order = setPosOrderId(order, order.OrderID)

	order.CreationResult = models.CreationResult{
		Message: response.Message,
		OrderInfo: models.OrderInfo{
			Timestamp:      order.CreatedAt.Unix(),
			CreationStatus: "SUCCESS",
		},
	}

	utils.Beautify("successfully send order to burger king, response (struct)", response)
	utils.Beautify("finished order model result", order)

	return order, nil

}

func (s *burgerKingService) constructPosOrder(ctx context.Context, order models.Order) (burgerKingModels.Order, models.Order, error) {
	products := make([]burgerKingModels.Product, 0, len(order.Products))

	offers, err := s.bkOfferRepository.GetActiveOffers(ctx)

	if err != nil {
		order.FailReason.Code = OFFERS_MISSED_CODE
		order.FailReason.Message = OFFERS_MISSED
		return burgerKingModels.Order{}, order, err
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

func (s *burgerKingService) sendOrder(ctx context.Context, posOrder burgerKingModels.Order) (burgerKingModels.OrderResponse, error) {
	utils.Beautify("Burger king request body", posOrder)

	response, err := s.bkClient.SendOrder(ctx, posOrder)

	if err != nil {
		utils.Beautify("Burger king Err response", err)
		return response, err
	}

	utils.Beautify("Burger king response", response)

	return response, nil
}

func (s *burgerKingService) GetStopList(ctx context.Context) (result coreMenuModels.StopListItems, err error) {
	return coreMenuModels.StopListItems{}, errors.New("unsupported method")
}

func (s *burgerKingService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	return coreMenuModels.Menu{}, ErrUnsupportedMethod
}

func (s *burgerKingService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	if err := s.bkClient.CancelOrder(ctx, burgerKingModels.CancelOrderRequest{
		OrderID:         order.OrderID,
		StoreID:         order.StoreID,
		CancelReason:    order.CancelReason.Reason,
		PaymentStrategy: order.PaymentStrategy,
	}); err != nil {
		log.Err(err).Msg("Burger King cancel order error")
		return err
	}

	return nil
}

func (s *burgerKingService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (s *burgerKingService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (s *burgerKingService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
