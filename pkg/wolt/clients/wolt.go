package clients

import (
	"context"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/wolt/models"
	"github.com/kwaaka-team/orders-core/core/wolt/models_v2"
	dto2 "github.com/kwaaka-team/orders-core/pkg/wolt/clients/dto"
)

type Config struct {
	Protocol           string
	BaseURL            string
	Insecure           bool
	ApiKey             string
	StoreID            string
	Username, Password string
}

type Wolt interface {
	BulkUpdate(ctx context.Context, storeID string, products dto2.UpdateProducts) (string, error)
	BulkAttribute(ctx context.Context, storeID string, attributes dto2.UpdateAttributes) (string, error)
	AcceptOrder(ctx context.Context, order dto2.AcceptOrderRequest) error
	AcceptSelfDeliveryOrder(ctx context.Context, order dto2.AcceptSelfDeliveryOrderOrderRequest) error
	RejectOrder(ctx context.Context, order dto2.RejectOrderRequest) error
	MarkOrder(ctx context.Context, orderID string) error
	ConfirmPreOrder(ctx context.Context, orderID string) error
	DeliveredOrder(ctx context.Context, orderID string) error
	GetOrder(ctx context.Context, orderID string) (models.Order, error)
	GetOrderByV2(ctx context.Context, orderID string) (models_v2.Order, error)
	UploadMenu(ctx context.Context, menu dto2.Menu, storeID string) error
	ManageStore(ctx context.Context, storeStatus dto2.IsStoreOpen) error
	UpdateStopListByProducts(ctx context.Context, storeId string, products []menuModels.Product, isAvailable bool) error
	UpdateStopListByProductsBulk(ctx context.Context, storeId string, products []menuModels.Product) error
	GetMenu(ctx context.Context, storeID string) (dto2.Menu, error)
	GetStoreStatus(ctx context.Context, venueId string) (dto2.StoreStatusResponse, error)
	UpdateMenuItemInventory(ctx context.Context, storeID string, woltInventory dto2.WoltInventory) error
}
