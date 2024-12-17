package notifier

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models"
	storecoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
)

type Notifier interface {
	Notify(ctx context.Context, status string, order models.Order, storeGroup storecoreModels.StoreGroup, store storecoreModels.Store) error
}
