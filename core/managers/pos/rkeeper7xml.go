package pos

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/config"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/managers/validator"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	rkeeperXMLCli "github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml"
	rkeeperXMLConf "github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/create_order_request"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/create_order_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/save_order_request"
	storeClient "github.com/kwaaka-team/orders-core/pkg/store"
	storeModels "github.com/kwaaka-team/orders-core/pkg/store/dto"

	"strconv"
)

type RKeeper7XMLManager struct {
	rkeeper7XMLCli rkeeperXMLConf.RKeeper7
	storeCli       storeClient.Client
}

func NewRKeeper7XMLManager(globalConfig config.Configuration, store coreStoreModels.Store, storeCli storeClient.Client) (*RKeeper7XMLManager, error) {
	cli, err := rkeeperXMLCli.NewClient(&rkeeperXMLConf.Config{
		Protocol:               "http",
		BaseURL:                store.RKeeper7XML.Domain,
		Token:                  store.RKeeper7XML.Token,
		Username:               store.RKeeper7XML.Username,
		Password:               store.RKeeper7XML.Password,
		UCSUsername:            store.RKeeper7XML.UCSUsername,
		UCSPassword:            store.RKeeper7XML.UCSPassword,
		LicenseBaseURL:         globalConfig.RKeeper7XMLConfiguration.LicenseBaseURL,
		Anchor:                 store.RKeeper7XML.Anchor,
		ObjectID:               store.RKeeper7XML.ObjectID,
		StationID:              store.RKeeper7XML.StationID,
		StationCode:            store.RKeeper7XML.StationCode,
		LicenseInstanceGUID:    store.RKeeper7XML.LicenseInstanceGUID,
		ChildItems:             store.RKeeper7XML.ChildItems,
		ClassificatorItemIdent: store.RKeeper7XML.ClassificatorItemIdent,
		ClassificatorPropMask:  store.RKeeper7XML.ClassificatorPropMask,
		MenuItemsPropMask:      store.RKeeper7XML.MenuItemsPropMask,
		PropFilter:             store.RKeeper7XML.PropFilter,
		Cashier:                store.RKeeper7XML.Cashier,
	})
	if err != nil {
		return nil, err
	}

	_, err = cli.SetLicense(context.Background())
	if err != nil {
		return nil, fmt.Errorf("set license error: %s", err.Error())
	}

	return &RKeeper7XMLManager{
		rkeeper7XMLCli: cli,
		storeCli:       storeCli,
	}, nil
}
func (manager *RKeeper7XMLManager) SendOrder(ctx context.Context, order any, store coreStoreModels.Store) (any, error) {
	posOrder, ok := order.(create_order_request.RK7Query)
	if !ok {
		return order, fmt.Errorf("create order request cast error")
	}

	response, err := manager.rkeeper7XMLCli.CreateOrder(ctx, posOrder.RK7CMD.Order.Table.Code, posOrder.RK7CMD.Order.Station.ID, posOrder.RK7CMD.Order.PersistentComment, posOrder.RK7CMD.Order.OrderType.Code)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (manager *RKeeper7XMLManager) ConstructPosOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (any, models.Order, error) {
	if order.TableID == "" {
		order.TableID = store.RKeeper7XML.DefaultTable
	}

	if store.RKeeper7XML.StationID == "" {
		store.RKeeper7XML.StationID = "1"
	}

	clientComment := ""
	if order.SpecialRequirements != "" || order.AllergyInfo != "" {
		clientComment += "\nКомментарий клиента: " + order.SpecialRequirements + order.AllergyInfo
	}

	request := create_order_request.RK7Query{
		RK7CMD: create_order_request.RK7CMD{
			Order: create_order_request.CreateOrderRequest{
				PersistentComment: "Код заказа: " + order.PickUpCode + clientComment,
				Table: create_order_request.Table{
					Code: order.TableID, // TableID
				},
				OrderType: create_order_request.OrderType{
					Code: store.RKeeper7XML.OrderTypeCode,
				},
				Station: create_order_request.Station{
					ID: store.RKeeper7XML.StationID,
				},
			},
		},
	}

	return request, order, nil
}

func (manager *RKeeper7XMLManager) setDeliveryTypeToOrder(ctx context.Context, order models.Order) error {
	return manager.rkeeper7XMLCli.SetDeliveryTypeToOrder(ctx, order.PosOrderID, order.PosPaymentInfo.OrderType)
}

func (manager *RKeeper7XMLManager) saveOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, coreStoreModels.Store, error) {
	dishes := make(save_order_request.Dishes, 0, len(order.Products))

	for _, product := range order.Products {
		dish := save_order_request.Dish{
			ID:       product.ID,
			Price:    strconv.Itoa(int(product.Price.Value)),
			Quantity: strconv.Itoa(product.Quantity),
		}

		for _, attribute := range product.Attributes {
			dish.Modi = append(dish.Modi, save_order_request.Modi{
				ID:    attribute.ID,
				Count: strconv.Itoa(attribute.Quantity),
				Price: strconv.Itoa(int(attribute.Price.Value)),
			})
		}

		dishes = append(dishes, dish)
	}

	_, err := manager.rkeeper7XMLCli.SaveOrder(ctx, order.PosOrderID, strconv.Itoa(store.RKeeper7XML.SeqNumber), order.PosPaymentInfo.PaymentTypeID, store.RKeeper7XML.PrepayReasonId, dishes, strconv.Itoa(int(order.EstimatedTotalPrice.Value)), store.RKeeper7XML.IsLifeTimeLicence)
	if err != nil {
		return order, store, err
	}

	return order, store, nil
}

func (manager *RKeeper7XMLManager) CreateOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) (models.Order, error) {
	defer func(store *coreStoreModels.Store) {
		if err := manager.storeCli.Update(ctx, storeModels.UpdateStore{
			ID: &store.ID,
			RKeeper7XML: &storeModels.UpdateStoreRKeeper7XMLConfig{
				SeqNumber: &store.RKeeper7XML.SeqNumber,
			},
		}); err != nil {
			return
		}
	}(&store)

	posOrder, _, err := manager.ConstructPosOrder(ctx, order, store)

	if err != nil {
		return order, err
	}

	rkeeperXMLOrder, ok := posOrder.(create_order_request.RK7Query)

	if !ok {
		return order, validator.ErrCastingPos
	}

	response, err := manager.SendOrder(ctx, rkeeperXMLOrder, store)
	if err != nil {
		return order, err
	}

	orderResponse, ok := response.(create_order_response.CreateOrderResponse)
	if !ok {
		return order, fmt.Errorf("cast create order response error")
	}

	order.PosOrderID = orderResponse.VisitID
	order.CreationResult = models.CreationResult{
		Message: orderResponse.Status,
		OrderInfo: models.OrderInfo{
			ID:             orderResponse.VisitID,
			OrganizationID: store.RKeeper7XML.ObjectID,
			CreationStatus: orderResponse.Status,
		},
	}
	order.PosPaymentInfo.OrderType = store.RKeeper7XML.OrderType

	if err = manager.setDeliveryTypeToOrder(ctx, order); err != nil {
		return order, err
	}

	order, store, err = manager.saveOrder(ctx, order, store)
	if err != nil {
		store.RKeeper7XML.SeqNumber = store.RKeeper7XML.SeqNumber + 1
		return order, err
	}

	store.RKeeper7XML.SeqNumber = store.RKeeper7XML.SeqNumber + 1

	return order, nil
}

func (manager *RKeeper7XMLManager) GetOrderStatus(ctx context.Context, order models.Order, store coreStoreModels.Store) (string, error) {
	response, err := manager.rkeeper7XMLCli.GetOrder(ctx, order.PosOrderID)
	if err != nil {
		return "", err
	}

	return response.CommandResult.Order.Finished, nil
}

func (manager *RKeeper7XMLManager) CancelOrder(ctx context.Context, order models.Order, cancelReason, paymentStrategy string, store coreStoreModels.Store) error {
	return errors.ErrUnsupportedMethod
}

func (manager *RKeeper7XMLManager) UpdateOrderProblem(ctx context.Context, organizationID, posOrderID string) error {
	return ErrUnsupportedMethod
}
