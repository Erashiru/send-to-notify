package pos

import (
	"context"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/config"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
	"strings"
)

type ctMaxService struct {
	*BasePosService
}

func newCTMaxService(bps *BasePosService) (*ctMaxService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "pos integration service constructor error")
	}
	return &ctMaxService{bps}, nil
}

func (*ctMaxService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {
	switch posStatus {
	case "ACCEPTED":
		return models.ACCEPTED, nil
	case "COOKING_STARTED":
		return models.COOKING_STARTED, nil
	case "COOKING_COMPLETE":
		return models.COOKING_COMPLETE, nil
	case "CLOSED":
		return models.CLOSED, nil
	}

	return 0, models.StatusIsNotExist
}

func (*ctMaxService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	return coreMenuModels.Menu{}, ErrUnsupportedMethod
}

func (*ctMaxService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration,
	store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	if !strings.Contains(order.Customer.PhoneNumber, "+") {
		order.Customer.PhoneNumber = "+77771111111"
	}

	var serviceFee float64

	for index, product := range order.Products {
		for position, attribute := range product.Attributes {
			order.Products[index].Attributes[position].Quantity = order.Products[index].Attributes[position].Quantity * product.Quantity
			serviceFee += attribute.Price.Value * float64(product.Quantity) * float64(attribute.Quantity)
		}
	}

	if serviceFee != 0 {
		order.HasServiceFee = true
		order.ServiceFeeSum = serviceFee
		order.EstimatedTotalPrice.Value = order.EstimatedTotalPrice.Value - serviceFee
	}

	order = setPosOrderId(order, uuid.New().String())

	order, promosMap, giftMap, promoWithPercentMap, err := getPromosMap(ctx, order, menuClient)
	if err != nil {
		return order, err
	}

	order = applyOrderDiscount(ctx, order, promosMap, giftMap, promoWithPercentMap)
	order.EstimatedTotalPrice.Value -= order.PartnerDiscountsProducts.Value

	return order, nil
}

func (p *ctMaxService) GetStopList(ctx context.Context) (result coreMenuModels.StopListItems, err error) {
	return coreMenuModels.StopListItems{}, errors.New("unsupported method")
}

func (p *ctMaxService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	return nil
}

func (s *ctMaxService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (s *ctMaxService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (s *ctMaxService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
