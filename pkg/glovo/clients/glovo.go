package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/glovo/models"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	dto2 "github.com/kwaaka-team/orders-core/pkg/glovo/clients/dto"
	"time"
)

type Config struct {
	Protocol string
	BaseURL  string
	Insecure bool
	ApiKey   string
	StoreID  string
}

type Glovo interface {
	UploadMenu(ctx context.Context, req dto2.UploadMenuRequest) (dto2.UploadMenuResponse, error)
	BulkUpdate(ctx context.Context, storeId string, request dto2.BulkUpdateRequest) (string, error)
	VerifyMenu(ctx context.Context, storeId, trxId string) (dto2.UploadMenuResponse, error)
	ValidateMenu(ctx context.Context, req dto2.ValidateMenuRequest) (dto2.ValidateMenuResponse, error)
	UpdateOrderStatus(ctx context.Context, order dto2.OrderUpdateRequest) (dto2.OrderUpdateResponse, error)
	ModifyOrderProduct(ctx context.Context, order models.ModifyOrderProductRequest) (*models.Order, error)
	ModifyProduct(ctx context.Context, req dto2.ProductModifyRequest) (dto2.ProductModifyResponse, error)
	ModifyAttribute(ctx context.Context, req dto2.AttributeModifyRequest) (dto2.AttributeModifyResponse, error)
	OpenStore(ctx context.Context, req dto2.StoreManageRequest) error
	CloseStore(ctx context.Context, req dto2.StoreManageRequest) error
	StoreStatus(ctx context.Context, storeID string) (dto2.StoreStatusResponse, error)
	UpdateStopListByProducts(ctx context.Context, storeId string, products []menuModels.Product, isAvailable bool) (string, error)
	UpdateStopListByProductsBulk(ctx context.Context, storeId string, products []menuModels.Product) (string, error)
	UpdateStopListByAttributesBulk(ctx context.Context, storeId string, attributes []menuModels.Attribute) (string, error)
	GetStoreSchedule(ctx context.Context, storeID string) (dto2.StoreScheduleResponse, error)
	GetStoreStatus(ctx context.Context, storeID string) (models.StoreStatus, error)
	GetBusyMode(ctx context.Context, storeId string) (int, error)
	CreateBusyMode(ctx context.Context, storeId string, additionalPreparationTimeInMinutes int) error
	DeleteBusyMode(ctx context.Context, storeId string) error

	AcceptOrder(ctx context.Context, storeId, orderId string, adjustedPickUpTime time.Time) (dto2.Response, error)
	MarkOrderAsReady(ctx context.Context, storeId, orderId string) (dto2.Response, error)
	MarkOrderAsOutForDelivery(ctx context.Context, storeId, orderId string) (dto2.Response, error)
	MarkOrderAsCustomerPickedUp(ctx context.Context, storeId, orderId string) (dto2.Response, error)
}
