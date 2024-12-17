package aggregator

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/deliveroo"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/errors"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/express24"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/external"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/glovo"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/moysklad"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/starterapp"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/talabat"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/wolt"
	"github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/yandex"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/store"
)

type Base interface {
	ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error)
	UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error)
	BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error)
	VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error)
	BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error)
	GetMenu(ctx context.Context, extStoreId string) (models.Menu, error)
	ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error)
}

// NewManager used for integrate with aggregators as glovo, wolt, yandex etc.
func NewManager(ctx context.Context, cfg menu.Configuration, aggregatorName storeModels.AggregatorName, store storeModels.Store, mspRepo drivers.MSPositionsRepository,
	stRepo drivers.StopListTransactionRepository, storeRepo drivers.StoreRepository, menuRepo drivers.MenuRepository, restGroupMenuRepo drivers.RestaurantGroupMenuRepository, storeCli store.Client) (Base, error) {
	switch models.AggregatorName(aggregatorName) {
	case models.GLOVO:
		return glovo.NewManager(ctx, cfg)
	case models.WOLT:
		cfg.WoltConfiguration.ApiKey = store.Wolt.ApiKey
		cfg.WoltConfiguration.Username = store.Wolt.MenuUsername
		cfg.WoltConfiguration.Password = store.Wolt.MenuPassword
		return wolt.NewManager(ctx, cfg)
	case models.EMENU:
		return external.NewManager(ctx, cfg, store, aggregatorName.String())
	case models.MOYSKLAD:
		cfg.MoySkladConfiguration.Username = store.MoySklad.UserName
		cfg.MoySkladConfiguration.Password = store.MoySklad.Password
		return moysklad.NewManager(ctx, cfg, mspRepo, stRepo)
	//case models.EXPRESS24:
	//	cfg.Express24Configuration.Username = store.Express24.Username
	//	cfg.Express24Configuration.Password = store.Express24.Password
	//	return express24.NewManager(ctx, cfg)
	case models.EXPRESS24:
		cfg.Express24Configuration.Token = store.Express24.Token
		return express24.NewMenuManager(ctx, cfg, restGroupMenuRepo, storeCli)
	case models.TALABAT:
		return talabat.NewManager(ctx, cfg, store, storeRepo, menuRepo)
	case models.DELIVEROO:
		return deliveroo.NewManager(ctx, store.Deliveroo.Username, store.Deliveroo.Password, store.Deliveroo.BaseURL)
	case models.YANDEX:
		return yandex.NewManager(ctx, cfg, store, storeRepo, menuRepo)
	case models.STARTERAPP:
		return starterapp.NewStarterAppMenuManager(ctx, store.StarterApp.ApiKey, cfg, menuRepo)
	}
	return nil, errors.ErrAggregatorNotFound
}
