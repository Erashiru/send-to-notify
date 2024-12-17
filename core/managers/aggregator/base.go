package aggregator

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"time"
)

type BaseAggregatorManager interface {
	AcceptOrder(ctx context.Context, orderID string, pickUpTime *time.Time) error
	AcceptSelfDeliveryOrder(ctx context.Context, orderID string, deliveryTime *time.Time) error
	RejectOrder(ctx context.Context, orderID, reason string) error
	MarkOrder(ctx context.Context, orderID string) error
	ConfirmPreOrder(ctx context.Context, orderID string) error
	UpdateOrderStatus(ctx context.Context, req selector.OrderStatusUpdate) error
	DeliveredOrder(ctx context.Context, orderID string) error
}

func NewAggregatorManager(delivery string, globalConfig config.Configuration, store coreStoreModels.Store) (BaseAggregatorManager, error) {
	switch delivery {
	case models.WOLT.String():
		return NewWoltManager(globalConfig, store)
	case models.GLOVO.String():
		return NewGlovoManager(globalConfig)
	case models.EMENU.String():
		return NewExternalManager(store, delivery)
	case models.DELIVEROO.String():
		return NewManager(store.Deliveroo.Username, store.Deliveroo.Password, store.Deliveroo.BaseURL)
	case models.TALABAT.String():
		return NewTalabatManager(globalConfig.TalabatConfiguration.MiddlewareBaseURL, store.Talabat.Username, store.Talabat.Password)
	case models.KWAAKA_ADMIN.String():
		return NewKwaakaAdminManager()
	case models.QRMENU.String():
		return NewQRMenuManager()
	case models.STARTERAPP.String():
		return NewStarterAppManager(globalConfig, store.StarterApp.ApiKey)
	}

	return nil, errors.New("aggregator manager doesn't initialize")
}
