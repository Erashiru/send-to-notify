package clients

import (
	"context"

	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

type Config struct {
	Protocol string
	BaseURL  string
	BaseAPI  string
	Insecure bool

	ApiLogin string
}

type IIKO interface {
	Auth(ctx context.Context) error
	Close(ctx context.Context)

	GetOrganizations(ctx context.Context) ([]models.Info, error)
	GetMenu(ctx context.Context, organizationID string) (models.GetMenuResponse, error)
	GetExternalMenu(ctx context.Context, organizationID, externalMenuID, priceCategoryId string) (models.GetExternalMenuResponse, error)

	CreateDeliveryOrder(ctx context.Context, req models.CreateDeliveryRequest) (models.CreateDeliveryResponse, error)
	RetrieveDeliveryOrder(ctx context.Context, organizationID, orderID string) (models.RetrieveOrder, error)
	CancelDeliveryOrder(ctx context.Context, organizationID, orderID, removalTypeId string) (models.CorID, error)
	UpdateOrderProblem(ctx context.Context, problem models.UpdateOrderProblem) error
	CloseOrder(ctx context.Context, posOrderId, organizationId string) error

	GetWebhookSetting(ctx context.Context, organizationID string) (models.GetWebhookSettingResponse, error)
	UpdateWebhookSetting(ctx context.Context, request models.UpdateWebhookRequest) (models.CorID, error)

	GetStopList(ctx context.Context, req models.StopListRequest) (models.StopListResponse, error)
	GetOrderTypes(ctx context.Context, req models.OrderTypesRequest) (models.OrderTypesResponse, error)
	CreateTableOrder(ctx context.Context, req models.CreateDeliveryRequest) (models.CreateDeliveryResponse, error)
	AddOrderItem(ctx context.Context, req models.OrderItem) (models.OrderItemResponse, error)

	GetTables(ctx context.Context, req models.TableRequest) (models.TableResponse, error)
	GetOrdersByTables(ctx context.Context, req models.OrdersByTablesRequest) (models.OrdersByTablesResponse, error)
	GetOrdersByIDs(ctx context.Context, req models.GetOrdersByIDsRequest) (models.OrdersByTablesResponse, error)
	GetCombos(ctx context.Context, req models.GetCombosRequest) (models.GetCombosResponse, error)
	IsAlive(ctx context.Context, req models.IsAliveRequest) (models.IsAliveResponse, error)
	AwakeTerminal(ctx context.Context, req models.IsAliveRequest) (models.AwakeResponse, error)
	GetCustomerInfo(ctx context.Context, req models.GetCustomerInfoRequest) (models.GetCustomerInfoResponse, error)
	GetTerminalGroups(ctx context.Context, organizationID string) (models.TerminalGroupsResponse, error)
	GetCustomerTransactions(ctx context.Context, req models.GetTransactionInfoReq) (models.GetTransactionInfoResp, error)
	GetDiscounts(ctx context.Context, organizationID string) (models.StoreDiscountsResponse, error)
	AddOrdersPayment(ctx context.Context, req models.ChangePaymentReq) (string, error)
	SendNotification(ctx context.Context, notificationInfo models.SendNotificationRequest) error
	CloseTableOrder(ctx context.Context, req models.CloseTableOrderReq) (string, error)
	GetCommandStatus(ctx context.Context, req models.GetCommandStatusReq) error
}
