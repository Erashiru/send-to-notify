package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
)

type Aggregator interface {
	UpdateStoreStatus(ctx context.Context, responses []models.StoreManagementResponse)
	OpenStore(ctx context.Context, aggregatorStoreID, systemStoreID string) (models.StoreManagementResponse, error)
	CloseStore(ctx context.Context, aggregatorStoreID, systemStoreID string) (models.StoreManagementResponse, error)
}
