package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/create_order_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_menu_by_category"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_menu_modifiers_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_groups"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_schema_details"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_schemas"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_order_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/pay_order_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/product_price_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/save_order_request"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/save_order_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/trade_group_details_response"
)

type Config struct {
	Protocol               string
	BaseURL                string
	Insecure               bool
	Username               string
	Password               string
	UCSUsername            string
	UCSPassword            string
	Token                  string
	LicenseBaseURL         string
	Anchor                 string
	ObjectID               string
	StationID              string
	StationCode            string
	LicenseInstanceGUID    string
	ChildItems             int
	ClassificatorItemIdent int
	ClassificatorPropMask  string
	MenuItemsPropMask      string
	PropFilter             string
	Cashier                string
}

type RKeeper7 interface {
	GetMenuItems(ctx context.Context) (models.MenuRK7QueryResult, error)
	GetMenuModifiers(ctx context.Context) (get_menu_modifiers_response.RK7QueryResult, error)
	GetOrderMenu(ctx context.Context) (models.OrderMenuRK7QueryResult, error)
	CreateOrder(ctx context.Context, table, stationID, comment, orderTypeCode string) (create_order_response.CreateOrderResponse, error)
	SetLicense(ctx context.Context) (models.LicenseResponse, error)
	SaveOrder(ctx context.Context, visitID, seqNumber, paymentId, paymentReasonId string, dishes save_order_request.Dishes, amount string, isLifeLicence bool) (save_order_response.RK7QueryResult, error)
	GetOrder(ctx context.Context, visitID string) (get_order_response.RK7QueryResult, error)
	GetMenuModifierSchemaDetails(ctx context.Context) (get_modifier_schema_details.RK7QueryResult, error)
	GetMenuModifierGroups(ctx context.Context) (get_modifier_groups.RK7QueryResult, error)
	PayOrder(ctx context.Context, orderGUID, paymentID, sum, stationCode string) (pay_order_response.RK7QueryResult, error)
	GetMenuModifierSchemas(ctx context.Context) (get_modifier_schemas.RK7QueryResult, error)
	GetMenuByCategory(ctx context.Context) (get_menu_by_category.RK7QueryResult, error)
	SetDeliveryTypeToOrder(ctx context.Context, visitId, orderType string) error
	GetItemsByTradeGroup(ctx context.Context, tradeGroupId string) (trade_group_details_response.RK7QueryResult, error)
	GetDeliveryPriceByPriceTypeId(ctx context.Context, priceTypeId string) (product_price_response.RK7QueryResult, error)
	GetSeqNumber(ctx context.Context) (models.GetSeqNumberRK7QueryResult, error)
}
