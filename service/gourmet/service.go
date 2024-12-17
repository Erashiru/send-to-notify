package gourmet

import (
	"context"
	"github.com/google/martian/v3/log"
	"github.com/kwaaka-team/orders-core/core/models/gourmet"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/iiko/clients"
	iikoConf "github.com/kwaaka-team/orders-core/pkg/iiko/clients"
	IIKOClient "github.com/kwaaka-team/orders-core/pkg/iiko/clients/http"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"time"
)

type ServiceImpl struct {
	storeService store.Service
	iikoBaseUrl  string
}

func NewServiceImpl(storeService store.Service, iikoBaseUrl string) (*ServiceImpl, error) {
	if storeService == nil {
		return nil, errors.New("storeService  is nil")
	}
	return &ServiceImpl{
		storeService: storeService,
		iikoBaseUrl:  iikoBaseUrl,
	}, nil
}

func (s *ServiceImpl) GetRestaurantOrders(ctx context.Context, restaurantID, tableID, orderID string) (gourmet.GourmetGetOrdersResponse, error) {
	store, err := s.storeService.GetByID(ctx, restaurantID)
	if err != nil {
		return gourmet.GourmetGetOrdersResponse{}, err
	}

	iikoCli, err := s.getIikoClient(store)
	if err != nil {
		return gourmet.GourmetGetOrdersResponse{}, err
	}

	if tableID != "" {
		tableOrders, err := iikoCli.GetOrdersByTables(ctx, iikoModels.OrdersByTablesRequest{
			OrganizationIds: []string{store.IikoCloud.OrganizationID},
			TableIds:        []string{tableID},
			Statuses:        []string{"New", "Bill"},
		})
		if err != nil {
			return gourmet.GourmetGetOrdersResponse{}, err
		}

		return gourmet.GourmetGetOrdersResponse{
			Orders: s.fromIikoOrderToGourmetOrders(tableOrders.Orders, tableID),
		}, nil
	}

	orders, err := iikoCli.GetOrdersByIDs(ctx, iikoModels.GetOrdersByIDsRequest{
		OrganizationIds: []string{store.IikoCloud.OrganizationID},
		OrderIds:        []string{orderID},
	})
	if err != nil {
		return gourmet.GourmetGetOrdersResponse{}, err
	}

	return gourmet.GourmetGetOrdersResponse{
		Orders: s.fromIikoOrderToGourmetOrders(orders.Orders, ""),
	}, nil
}

func (s *ServiceImpl) GetRestaurantTables(ctx context.Context, restaurantID string) (gourmet.GourmetGetTablesResponse, error) {
	store, err := s.storeService.GetByID(ctx, restaurantID)
	if err != nil {
		return gourmet.GourmetGetTablesResponse{}, err
	}

	iikoCli, err := s.getIikoClient(store)
	if err != nil {
		return gourmet.GourmetGetTablesResponse{}, err
	}

	//terminalGroups, err := iikoCli.GetTerminalGroups(ctx, store.IikoCloud.OrganizationID)
	//if err != nil {
	//	return "", err
	//}
	//
	//terminalGroupsIDs := make([]string, 0, len(terminalGroups.TerminalGroups))
	//for _, terminalGroup := range terminalGroups.TerminalGroups {
	//	for _, item := range terminalGroup.Items {
	//		terminalGroupsIDs = append(terminalGroupsIDs, item.ID)
	//	}
	//}

	tablesResp, err := iikoCli.GetTables(ctx, iikoModels.TableRequest{
		TerminalGroupIds: []string{store.IikoCloud.TerminalID},
	})
	if err != nil {
		return gourmet.GourmetGetTablesResponse{}, err
	}

	tabels := make([]gourmet.GourmetTabels, 0, len(tablesResp.RestaurantSections))

	for _, section := range tablesResp.RestaurantSections {
		for _, table := range section.Tables {
			tabels = append(tabels, gourmet.GourmetTabels{
				Id:              table.Id,
				Number:          table.Number,
				Name:            table.Name,
				SeatingCapacity: table.SeatingCapacity,
				SectionId:       section.Id,
				SectionName:     section.Name,
			})
		}
	}

	return gourmet.GourmetGetTablesResponse{Tables: tabels}, nil
}

func (s *ServiceImpl) CreatePayment(ctx context.Context, restaurantID, orderID, paymentTypeId, paymentTypeKind string, isPaid bool) (gourmet.PaymentChangeResponse, error) {
	if !isPaid {
		return gourmet.PaymentChangeResponse{}, errors.New("payment is not paid")
	}

	if paymentTypeId == "" || paymentTypeKind == "" {
		return gourmet.PaymentChangeResponse{}, errors.New("payment type id or kind is empty")
	}

	store, err := s.storeService.GetByID(ctx, restaurantID)
	if err != nil {
		return gourmet.PaymentChangeResponse{}, err
	}

	iikoCli, err := s.getIikoClient(store)
	if err != nil {
		return gourmet.PaymentChangeResponse{}, err
	}

	order, err := iikoCli.GetOrdersByIDs(ctx, iikoModels.GetOrdersByIDsRequest{
		OrganizationIds: []string{store.IikoCloud.OrganizationID},
		OrderIds:        []string{orderID},
	})
	if err != nil {
		return gourmet.PaymentChangeResponse{}, err
	}

	if len(order.Orders) == 0 {
		return gourmet.PaymentChangeResponse{}, errors.New("no orders found")
	}

	sum := s.calculateTotalToPay(order.Orders[0])

	if sum > 0 {
		corID, err := iikoCli.AddOrdersPayment(ctx, iikoModels.ChangePaymentReq{
			OrganizationId: store.IikoCloud.OrganizationID,
			OrderId:        orderID,
			Payments: []iikoModels.PaymentReq{
				{
					PaymentTypeKind: paymentTypeKind,
					PaymentTypeId:   paymentTypeId,
					Sum:             sum,
					//IsProcessedExternally: true,
					//IsFiscalizedExternally: true,
				},
			},
		})
		if err != nil {
			return gourmet.PaymentChangeResponse{}, err
		}

		log.Infof("AddOrdersPayment, orderID: %s, corID: %s", orderID, corID)
		time.Sleep(2 * time.Second)
		if err := iikoCli.GetCommandStatus(ctx, iikoModels.GetCommandStatusReq{
			OrganizationId: store.IikoCloud.OrganizationID,
			CorrelationId:  corID,
		}); err != nil {
			return gourmet.PaymentChangeResponse{}, err
		}
	}

	corID, err := iikoCli.CloseTableOrder(ctx, iikoModels.CloseTableOrderReq{
		OrganizationId: store.IikoCloud.OrganizationID,
		OrderId:        orderID,
	})

	if err != nil {
		return gourmet.PaymentChangeResponse{}, err
	}

	log.Infof("CloseTableOrder, orderID: %s, corID: %s", orderID, corID)

	time.Sleep(time.Second)

	if err := iikoCli.GetCommandStatus(ctx, iikoModels.GetCommandStatusReq{
		OrganizationId: store.IikoCloud.OrganizationID,
		CorrelationId:  corID,
	}); err != nil {
		return gourmet.PaymentChangeResponse{}, err
	}

	tableID := ""
	if len(order.Orders[0].Order.TableIds) != 0 {
		tableID = order.Orders[0].Order.TableIds[0]
	}

	return gourmet.PaymentChangeResponse{
		OrderId: orderID,
		TableId: tableID,
		Status:  "closed",
		Message: "check status and totalToPay in get order",
	}, nil
}

func (s *ServiceImpl) getIikoClient(store coreStoreModels.Store) (clients.IIKO, error) {
	iikoClient, err := IIKOClient.New(&iikoConf.Config{
		Protocol: "http",
		BaseURL:  s.iikoBaseUrl,
		ApiLogin: store.IikoCloud.Key,
	})
	if err != nil {
		return nil, err
	}

	return iikoClient, nil
}

func (s *ServiceImpl) fromIikoOrderToGourmetOrders(iikoOrders []iikoModels.Orders, tableID string) []gourmet.GourmetOrder {
	res := make([]gourmet.GourmetOrder, 0, len(iikoOrders))
	for _, order := range iikoOrders {
		res = append(res, gourmet.GourmetOrder{
			Id:              order.ID,
			TableID:         tableID,
			Customer:        order.Order.Customer,
			Waiter:          order.Order.Waiter,
			Phone:           order.Order.Phone,
			Status:          order.Order.Status,
			Sum:             order.Order.Sum,
			TotalToPay:      s.calculateTotalToPay(order),
			WhenCreated:     order.Order.WhenCreated,
			WhenBillPrinted: order.Order.WhenBillPrinted,
			WhenClosed:      order.Order.WhenClosed,
			Items:           order.Order.Items,
			Discounts:       order.Order.Discounts,
			Payments:        order.Order.Payments,
		})
	}

	return res
}

func (s *ServiceImpl) calculateTotalToPay(order iikoModels.Orders) float64 {
	//totalSum := order.Order.Sum
	//paidSum := float64(0)
	//for _, payment := range order.Order.Payments {
	//	paidSum = paidSum + payment.Sum
	//}
	//totalToPay := totalSum - paidSum
	totalToPay := order.Order.Sum - order.Order.ProcessedPaymentsSum

	return totalToPay
}
