package pos

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/config"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/kwaaka-team/orders-core/service/error_solutions"
	"github.com/pkg/errors"
)

type kwaakaPosService struct {
	*BasePosService
}

func newKwaakaPosService(bps *BasePosService) (*kwaakaPosService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "kwaaka pos constructor error")
	}

	return &kwaakaPosService{bps}, nil
}

func (ks *kwaakaPosService) CreateOrder(ctx context.Context, order models.Order, globalConfig config.Configuration,
	store coreStoreModels.Store, menu coreMenuModels.Menu, menuClient menuCore.Client, aggregatorMenu coreMenuModels.Menu,
	storeCli storeClient.Client, errSolution error_solutions.Service, notifyQueue notifyQueue.SQSInterface) (models.Order, error) {
	return order, nil
}

func (ks *kwaakaPosService) IsAliveStatus(ctx context.Context, store coreStoreModels.Store) (bool, error) {
	return true, nil
}

func (ks *kwaakaPosService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {
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
	return 0, nil
}

func (ks *kwaakaPosService) GetOrderStatus(ctx context.Context, order models.Order) (string, error) {
	return order.Status, nil
}

func (ks *kwaakaPosService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {
	return nil, errors.New("method is temporarily out of service")
}

func (ks *kwaakaPosService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	return coreMenuModels.Menu{}, errors.New("method is temporarily out of service")
}

func (ks *kwaakaPosService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	return nil
}

func (s *kwaakaPosService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (s *kwaakaPosService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	return items, nil
}

func (s *kwaakaPosService) CloseOrder(ctx context.Context, posOrderId string) error {
	return nil
}
