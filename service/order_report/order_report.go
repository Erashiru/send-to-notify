package order_report

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	selector2 "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/kwaaka_3pl"
	"github.com/kwaaka-team/orders-core/service/order"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/rs/zerolog/log"
	"github.com/tealeg/xlsx"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	directDeliveryService        = "qr_menu"
	kwaakaAdminDeliveryService   = "kwaaka_admin"
	iokaPaymentSystem            = "ioka"
	kaspiPaymentSystem           = "kaspi"
	kaspiSalesCountPaymentSystem = "kaspi_salescout"
	whatsappPaymentSystem        = "whatsapp"
	byCashierPaymentSystem       = "by_cashier"
	cashPaymentSystem            = "cash"
	reportForKwaaka              = "kwaaka"
	reportForRestaurant          = "restaurant"
	noInfo                       = "no info"
	payments                     = map[string]string{"by_cashier": "By Cashier", "kaspi_salescout": "Kaspi Salescout", "ioka": "Ioka", "cash": "Cash"}
)

type OrderReport interface {
	OrderReportForRestaurant(ctx context.Context, query models.OrderReportRequest) (models.OrderReportResponse, error)
	OrderReportForRestaurantTotals(ctx context.Context, query models.OrderReportRequest) (models.OrderReportResponse, error)
	OrderReportForKwaakaTotals(ctx context.Context, query models.OrderReportRequest) (models.OrderReportResponse, error)
	OrderReportToXlsx(ctx context.Context, query models.OrderReportRequest) ([]byte, error)
	DeliveryDispatcherPrice(ctx context.Context) error
}

type OrderReportImpl struct {
	storeService      store.Service
	repository        order.Repository
	kwaaka3pl         kwaaka_3pl.Service
	cartService       order.CartService
	storeGroupService storegroup.Service
}

func NewOrderReportService(storeService store.Service, repository order.Repository, kwaaka3pl kwaaka_3pl.Service, cartService order.CartService, storeGroupService storegroup.Service) *OrderReportImpl {
	return &OrderReportImpl{
		storeService:      storeService,
		repository:        repository,
		kwaaka3pl:         kwaaka3pl,
		cartService:       cartService,
		storeGroupService: storeGroupService,
	}
}

func (or *OrderReportImpl) OrderReportForRestaurant(ctx context.Context, query models.OrderReportRequest) (models.OrderReportResponse, error) {
	var (
		err     error
		reports []models.OrderReport
		wg      sync.WaitGroup
		mut     sync.Mutex
	)

	storeMap, err := or.getStoreMapByQuery(ctx, query)
	if err != nil {
		return models.OrderReportResponse{}, err
	}

	restGrByRestMap, storeIDs, err := or.getMapStoreGroupByStores(ctx, storeMap)
	if err != nil {
		return models.OrderReportResponse{}, err
	}

	orders, totalOrdersCount, err := or.repository.GetAllOrders(ctx, selector.EmptyOrderSearch().
		SetRestaurants(storeIDs).
		SetDeliveryService("qr_menu").
		SetDeliveryDispatcher(query.DeliveryDispatcher).
		SetSearchReport(query.Search).
		SetOrderTimeFrom(query.StartDate).
		SetOrderTimeTo(query.EndDate).
		SetPage(query.Pagination.Page).
		SetLimit(query.Pagination.Limit).
		SetSorting("order_time.value", -1))
	if err != nil {
		return models.OrderReportResponse{}, err
	}

	kwaakaAdminOrders, kwaakaAdminTotalOrdersCount, err := or.repository.GetAllOrders(ctx, selector.EmptyOrderSearch().
		SetRestaurants(storeIDs).
		SetDeliveryService("kwaaka_admin").
		SetDeliveryDispatcher(query.DeliveryDispatcher).
		SetSearchReport(query.Search).
		SetOrderTimeFrom(query.StartDate).
		SetOrderTimeTo(query.EndDate).
		SetPage(query.Pagination.Page).
		SetLimit(query.Pagination.Limit).
		SetSorting("order_time.value", -1))
	if err != nil {
		return models.OrderReportResponse{}, err
	}

	orders = append(orders, kwaakaAdminOrders...)

	totalOrdersCount += kwaakaAdminTotalOrdersCount

	for _, ordr := range orders {
		wg.Add(1)
		go func(order models.Order) {
			defer wg.Done()
			if order.IsTestOrder {
				return
			}

			paymentSystem, err := or.GetPaymentSystem(ctx, order.DeliveryService, order.OrderID)
			if err != nil {
				return
			}

			if query.PaymentSystem != "" && query.PaymentSystem != paymentSystem {
				return
			}

			var paymentSystemBeautified = paymentSystem
			if val, ok := payments[paymentSystem]; ok {
				paymentSystemBeautified = val
			}

			deliveryStatus, err := or.GetOrderDeliveryStatus(ctx, order.OrderID)
			if err != nil {
				deliveryStatus = noInfo
			}

			driverDeliveryPrice, err := or.GetDeliveryPrice(ctx, order.DeliveryOrderID)
			if err != nil {
				driverDeliveryPrice = 0
			}

			deliveryHistoricalDeliveries, newDeliveryStatus, err := or.GetHistoricalDeliveriesInfo(ctx, order, deliveryStatus, driverDeliveryPrice, order.History3plDeliveryInfo, order.Canceled3PlDeliveryInfo)
			if err != nil {
				log.Error().Msgf("error occured: %s", err.Error())
			}

			if newDeliveryStatus != "" && deliveryStatus == "no info" {
				deliveryStatus = newDeliveryStatus
			}

			store, err := or.storeService.GetByID(ctx, order.RestaurantID)
			if err != nil {
				return
			}

			formattedTime := order.OrderTime.Value.Time.Add(time.Duration(store.Settings.TimeZone.UTCOffset) * time.Hour).Format("02-01-2006 15:04:05")
			if store.Settings.TimeZone.TZ == "Asia/Almaty" {
				formattedTime = order.OrderTime.Value.Time.Add(5 * time.Hour).Format("02-01-2006 15:04:05")
			}
			if order.Type == models.ORDER_TYPE_PREORDER || !order.Preorder.Time.Value.IsZero() {
				switch store.Settings.TimeZone.TZ {
				case "Asia/Almaty":
					formattedTime = order.Preorder.Time.Value.Add(5 * time.Hour).Format("02-01-2006 15:04:05")
				default:
					formattedTime = order.Preorder.Time.Value.Add(time.Duration(store.Settings.TimeZone.UTCOffset) * time.Hour).Format("02-01-2006 15:04:05")
				}
			}

			orderType := or.GetOrderType(order.Type, order.SendCourier)

			source := or.GetSource(order.DeliveryService)

			calculatedDeliveryHistoryPricesSum := or.getCalculatedDeliveryHistoryPricesSUM(deliveryHistoricalDeliveries)

			actualDeliveryHistoryPricesSum := or.getActualDeliveryHistoryPricesSUM(deliveryHistoricalDeliveries)

			dispatcherDeliveryHistoryPricesSum := or.getDispatcherDeliveryHistoryPricesSUM(deliveryHistoricalDeliveries)

			report := models.OrderReport{
				OrderID:                 order.OrderID,
				RestaurantID:            order.RestaurantID,
				RestaurantName:          order.RestaurantName,
				RestaurantGroupID:       restGrByRestMap[order.RestaurantID].ID,
				RestaurantGroupName:     restGrByRestMap[order.RestaurantID].Name,
				Source:                  source,
				OrderTime:               formattedTime,
				EstimatedPickupTime:     order.EstimatedPickupTime.Value.Time,
				OrderType:               orderType,
				OrderStatus:             order.Status,
				DeliveryStatus:          deliveryStatus,
				DeliveryStatusHistory:   or.getDeliveryHistoryStatuses(deliveryHistoricalDeliveries),
				DeliveryOrderHistoryIDs: or.getDeliveryHistoryIDs(deliveryHistoricalDeliveries),

				DeliveryAddress:   order.DeliveryAddress.Label,
				RestaurantAddress: fmt.Sprintf("%s, %s", storeMap[order.RestaurantID].Address.City, storeMap[order.RestaurantID].Address.Street),

				DeliveryOrderProviderHistory: or.getDeliveryHistoryProviders(deliveryHistoricalDeliveries),
				DeliveryDispatcher:           order.DeliveryDispatcher,
				OrderComment:                 order.AllergyInfo,
				PaymentSystem:                paymentSystemBeautified,
				EstimatedTotalPrice:          order.EstimatedTotalPrice.Value,
				TotalOrderPrice:              order.EstimatedTotalPrice.Value + order.ClientDeliveryPrice,
				CustomerName:                 order.Customer.Name,
				CustomerPhone:                order.Customer.PhoneNumber,
				SendCourier:                  order.SendCourier,
				Products:                     order.Products,
				Numbers: models.OrderReportNumbers{
					RestaurantIncome:  or.getIncome(reportForRestaurant, paymentSystem, order.EstimatedTotalPrice.Value, order.ClientDeliveryPrice, calculatedDeliveryHistoryPricesSum, dispatcherDeliveryHistoryPricesSum),
					KwaakaIncome:      or.getIncome(reportForKwaaka, paymentSystem, order.EstimatedTotalPrice.Value, order.ClientDeliveryPrice, calculatedDeliveryHistoryPricesSum, dispatcherDeliveryHistoryPricesSum),
					BalanceKwaaka:     or.getBalance(reportForRestaurant, paymentSystem, order.EstimatedTotalPrice.Value, order.ClientDeliveryPrice, calculatedDeliveryHistoryPricesSum),
					BalanceRestaurant: or.getBalance(reportForKwaaka, paymentSystem, order.EstimatedTotalPrice.Value, order.ClientDeliveryPrice, calculatedDeliveryHistoryPricesSum),
					// Прогнозируемые цены
					ProjectedDeliveryHistoryPrices:    or.getProjectedDeliveryHistoryPrices(deliveryHistoricalDeliveries),
					ProjectedDeliveryHistoryPricesSUM: or.getProjectedDeliveryHistoryPricesSum(deliveryHistoricalDeliveries),
					// Фактические
					ActualDeliveryHistoryPrices:    or.getActualDeliveryHistoryPrices(deliveryHistoricalDeliveries),
					ActualDeliveryHistoryPricesSUM: actualDeliveryHistoryPricesSum,
					// Расчетные
					CalculatedDeliveryHistoryPrices:    or.getCalculatedDeliveryHistoryPrices(deliveryHistoricalDeliveries),
					CalculatedDeliveryHistoryPricesSUM: calculatedDeliveryHistoryPricesSum,

					ClientDeliveryPrice:        order.ClientDeliveryPrice,
					KwaakaChargedDeliveryPrice: math.Ceil(order.KwaakaChargedDeliveryPrice),
					BankBalance:                or.getBankBalance(paymentSystem, order.EstimatedTotalPrice.Value, order.ClientDeliveryPrice),
					DeliveryBalance:            dispatcherDeliveryHistoryPricesSum,
				},
			}

			mut.Lock()
			reports = append(reports, report)
			mut.Unlock()
		}(ordr)
	}

	wg.Wait()

	return models.OrderReportResponse{
		OrdersReport:     reports,
		TotalOrdersCount: totalOrdersCount,
	}, nil
}

func (or *OrderReportImpl) OrderReportForRestaurantTotals(ctx context.Context, query models.OrderReportRequest) (models.OrderReportResponse, error) {
	var (
		storeIDs         []string
		ordersTotalPrice float64
		totalIncome      float64
		totalBalance     float64
		wg               sync.WaitGroup
		mut              sync.Mutex
	)

	storeMap, err := or.getStoreMapByQuery(ctx, query)
	if err != nil {
		return models.OrderReportResponse{}, err
	}

	for _, v := range storeMap {
		storeIDs = append(storeIDs, v.ID)
	}

	orders, totalOrdersCount, err := or.repository.GetAllOrders(ctx, selector.EmptyOrderSearch().
		SetRestaurants(storeIDs).
		SetDeliveryService(directDeliveryService).
		SetDeliveryDispatcher(query.DeliveryDispatcher).
		SetOrderTimeFrom(query.StartDate).
		SetOrderTimeTo(query.EndDate))
	if err != nil {
		return models.OrderReportResponse{}, err
	}

	for _, ordr := range orders {

		wg.Add(1)
		go func(order models.Order) {
			defer wg.Done()

			if order.DeliveryOrderID == "" || or.isFailedOrder(order.Status) || order.IsTestOrder {
				return
			}

			paymentSystem, err := or.GetPaymentSystem(ctx, order.DeliveryService, order.OrderID)
			if err != nil {
				return
			}

			driverDeliveryPrice, err := or.GetDeliveryPrice(ctx, order.DeliveryOrderID)
			if err != nil {
				driverDeliveryPrice = 0
			}

			if driverDeliveryPrice == 0 {
				return
			}

			orderTotalPrice := order.EstimatedTotalPrice.Value + order.ClientDeliveryPrice
			income := or.getIncome(reportForRestaurant, paymentSystem, order.EstimatedTotalPrice.Value, order.ClientDeliveryPrice, 0, 0)
			balance := or.getBalance(reportForRestaurant, paymentSystem, order.EstimatedTotalPrice.Value, order.ClientDeliveryPrice, 0)

			mut.Lock()
			defer mut.Unlock()
			ordersTotalPrice += orderTotalPrice
			totalIncome += income
			totalBalance += balance
		}(ordr)
	}

	wg.Wait()

	return models.OrderReportResponse{
		OrdersTotalPrice: ordersTotalPrice,
		TotalIncome:      totalIncome,
		TotalBalance:     totalBalance,
		TotalOrdersCount: totalOrdersCount,
	}, nil
}

func (or *OrderReportImpl) OrderReportForKwaakaTotals(ctx context.Context, query models.OrderReportRequest) (models.OrderReportResponse, error) {
	var (
		storeIDs                      []string
		ordersTotalPrice              float64
		totalIncome                   float64
		totalBalance                  float64
		totalDeliveryBalanceForKwaaka float64
		wg                            sync.WaitGroup
		mut                           sync.Mutex
	)

	storeMap, err := or.getStoreMapByQuery(ctx, query)
	if err != nil {
		return models.OrderReportResponse{}, err
	}

	for _, v := range storeMap {
		storeIDs = append(storeIDs, v.ID)
	}

	orders, totalOrdersCount, err := or.repository.GetAllOrders(ctx, selector.EmptyOrderSearch().
		SetRestaurants(storeIDs).
		SetDeliveryService(directDeliveryService).
		SetDeliveryDispatcher(query.DeliveryDispatcher).
		SetOrderTimeFrom(query.StartDate).
		SetOrderTimeTo(query.EndDate))
	if err != nil {
		return models.OrderReportResponse{}, err
	}

	for _, ordr := range orders {

		wg.Add(1)
		go func(order models.Order) {
			defer wg.Done()

			if order.DeliveryOrderID == "" || or.isFailedOrder(order.Status) || order.IsTestOrder {
				return
			}

			paymentSystem, err := or.GetPaymentSystem(ctx, order.DeliveryService, order.OrderID)
			if err != nil {
				return
			}

			driverDeliveryPrice, err := or.GetDeliveryPrice(ctx, order.DeliveryOrderID)
			if err != nil {
				driverDeliveryPrice = 0
			}

			if driverDeliveryPrice == 0 {
				return
			}

			orderTotalPrice := order.EstimatedTotalPrice.Value + order.ClientDeliveryPrice
			income := or.getIncome(reportForKwaaka, paymentSystem, order.EstimatedTotalPrice.Value, order.ClientDeliveryPrice, 0, 0)
			balance := or.getBalance(reportForKwaaka, paymentSystem, order.EstimatedTotalPrice.Value, order.ClientDeliveryPrice, 0)

			mut.Lock()
			defer mut.Unlock()
			ordersTotalPrice += orderTotalPrice
			totalIncome += income
			totalBalance += balance
			totalDeliveryBalanceForKwaaka += driverDeliveryPrice
		}(ordr)
	}

	wg.Wait()

	return models.OrderReportResponse{
		OrdersTotalPrice:              ordersTotalPrice,
		TotalIncome:                   totalIncome,
		TotalBalance:                  totalBalance,
		TotalDeliveryBalanceForKwaaka: totalDeliveryBalanceForKwaaka,
		TotalOrdersCount:              totalOrdersCount,
	}, nil
}

func (or *OrderReportImpl) getStoreMapByQuery(ctx context.Context, query models.OrderReportRequest) (map[string]storeModels.Store, error) {
	if query.RestaurantIDs == nil && query.RestaurantGroupIds == nil {
		return nil, fmt.Errorf("error no restaurant and restaurant group in query")
	}

	storeMap := make(map[string]storeModels.Store)

	if query.RestaurantGroupIds != nil && len(query.RestaurantGroupIds) != 0 {
		stores, err := or.getRestsByRestGroupIDs(ctx, query.RestaurantGroupIds)
		if err != nil {
			return nil, err
		}

		for i := range stores {
			if _, ok := storeMap[stores[i].ID]; !ok {
				storeMap[stores[i].ID] = stores[i]
			}
		}
	}

	if query.RestaurantIDs != nil && len(query.RestaurantIDs) != 0 {
		stores, err := or.storeService.GetStoresBySelectorFilter(ctx, selector2.NewEmptyStoreSearch().SetStoreIDs(query.RestaurantIDs))
		if err != nil {
			return nil, err
		}

		for i := range stores {
			if _, ok := storeMap[stores[i].ID]; !ok {
				storeMap[stores[i].ID] = stores[i]
			}
		}
	}

	return storeMap, nil
}

func (or *OrderReportImpl) getRestsByRestGroupIDs(ctx context.Context, restGroups []string) ([]storeModels.Store, error) {
	var rests []storeModels.Store
	for i := range restGroups {
		stores, err := or.storeService.GetRestaurantsByGroupId(ctx, selector2.Pagination{Page: 0, Limit: 0}, restGroups[i])
		if err != nil {
			return nil, err
		}
		rests = append(rests, stores...)
	}

	return rests, nil
}

func (or *OrderReportImpl) getMapStoreGroupByStores(ctx context.Context, storesMap map[string]storeModels.Store) (map[string]storeModels.StoreGroup, []string, error) {
	var (
		storeIDs        []string
		restGrByRestMap = make(map[string]storeModels.StoreGroup, len(storesMap))
	)

	for _, st := range storesMap {
		storeIDs = append(storeIDs, st.ID)

		group, err := or.storeGroupService.GetStoreGroupByID(ctx, st.RestaurantGroupID)
		if err != nil {
			return nil, nil, err
		}
		restGrByRestMap[st.ID] = group
	}

	return restGrByRestMap, storeIDs, nil
}

func (or *OrderReportImpl) getIncome(reportFor, paymentSystem string, estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum, dispatcherDeliveryHistoryPricesSum float64) float64 {
	switch reportFor {

	case reportForRestaurant:
		switch paymentSystem {
		case iokaPaymentSystem:
			return or.incomeForRestaurantIOKA(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum)
		case kaspiPaymentSystem, whatsappPaymentSystem, byCashierPaymentSystem, cashPaymentSystem:
			return or.incomeForRestaurantKaspiOrByCashierOrCash(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum)
		case kaspiSalesCountPaymentSystem:
			return or.incomeForRestaurantKaspiSalesCount(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum)
		}

	case reportForKwaaka:
		switch paymentSystem {
		case iokaPaymentSystem:
			return or.incomeForKwaakaIOKA(estimatedTotalPrice, calculatedDeliveryHistoryPricesSum, dispatcherDeliveryHistoryPricesSum)
		case kaspiPaymentSystem, whatsappPaymentSystem, byCashierPaymentSystem, cashPaymentSystem:
			return or.incomeForKwaakaKaspiOrByCashierOrCash(estimatedTotalPrice, calculatedDeliveryHistoryPricesSum, dispatcherDeliveryHistoryPricesSum)
		case kaspiSalesCountPaymentSystem:
			return or.incomeForKwaakaKaspiSalesCount(estimatedTotalPrice, calculatedDeliveryHistoryPricesSum, dispatcherDeliveryHistoryPricesSum)
		}
	}

	return 0
}

func (or *OrderReportImpl) incomeForRestaurantIOKA(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum float64) float64 {
	return estimatedTotalPrice - math.Ceil(estimatedTotalPrice*0.05) - math.Ceil(0.029*clientDeliveryPrice) - (calculatedDeliveryHistoryPricesSum - clientDeliveryPrice)
}

func (or *OrderReportImpl) incomeForRestaurantKaspiOrByCashierOrCash(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum float64) float64 {
	return estimatedTotalPrice - math.Ceil(estimatedTotalPrice*0.03) - (calculatedDeliveryHistoryPricesSum - clientDeliveryPrice)
}

func (or *OrderReportImpl) incomeForRestaurantKaspiSalesCount(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum float64) float64 {
	return estimatedTotalPrice - math.Ceil(estimatedTotalPrice*0.05) - math.Ceil(0.025*clientDeliveryPrice) - (calculatedDeliveryHistoryPricesSum - clientDeliveryPrice)
}

func (or *OrderReportImpl) incomeForKwaakaIOKA(estimatedTotalPrice, calculatedDeliveryHistoryPricesSum, dispatcherDeliveryHistoryPricesSum float64) float64 {
	return math.Ceil(0.021*estimatedTotalPrice) + calculatedDeliveryHistoryPricesSum - dispatcherDeliveryHistoryPricesSum
}

func (or *OrderReportImpl) incomeForKwaakaKaspiOrByCashierOrCash(estimatedTotalPrice, calculatedDeliveryHistoryPricesSum, dispatcherDeliveryHistoryPricesSum float64) float64 {
	return math.Ceil(0.03*estimatedTotalPrice) + calculatedDeliveryHistoryPricesSum - dispatcherDeliveryHistoryPricesSum
}

func (or *OrderReportImpl) incomeForKwaakaKaspiSalesCount(estimatedTotalPrice, calculatedDeliveryHistoryPricesSum, dispatcherDeliveryHistoryPricesSum float64) float64 {
	return math.Ceil(0.025*estimatedTotalPrice) + calculatedDeliveryHistoryPricesSum - dispatcherDeliveryHistoryPricesSum
}

func (or *OrderReportImpl) getBalance(reportFor, paymentSystem string, estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum float64) float64 {
	switch reportFor {

	case reportForKwaaka:
		switch paymentSystem {
		case iokaPaymentSystem:
			return or.balanceForKwaakaIOKA(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum)
		case kaspiPaymentSystem, whatsappPaymentSystem, byCashierPaymentSystem, cashPaymentSystem:
			return or.balanceForKwaakaKaspi(estimatedTotalPrice, calculatedDeliveryHistoryPricesSum)
		case kaspiSalesCountPaymentSystem:
			return or.balanceForKwaakaKaspiSalesCount(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum)
		}

	case reportForRestaurant:
		switch paymentSystem {
		case iokaPaymentSystem:
			return or.balanceForRestaurantIOKA(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum)
		case kaspiPaymentSystem, whatsappPaymentSystem, byCashierPaymentSystem, cashPaymentSystem:
			return or.balanceForRestaurantKaspi(estimatedTotalPrice, calculatedDeliveryHistoryPricesSum)
		case kaspiSalesCountPaymentSystem:
			return or.balanceForRestaurantKaspiSalesCount(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum)
		}
	}

	return 0
}

func (or *OrderReportImpl) balanceForKwaakaIOKA(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum float64) float64 {
	return -estimatedTotalPrice + math.Ceil(estimatedTotalPrice*0.05) + math.Ceil(0.029*clientDeliveryPrice) + (calculatedDeliveryHistoryPricesSum - clientDeliveryPrice)
}

func (or *OrderReportImpl) balanceForKwaakaKaspi(estimatedTotalPrice, calculatedDeliveryHistoryPricesSum float64) float64 {
	return math.Ceil(+0.03*estimatedTotalPrice) + calculatedDeliveryHistoryPricesSum
}

func (or *OrderReportImpl) balanceForKwaakaKaspiSalesCount(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum float64) float64 {
	return -estimatedTotalPrice + math.Ceil(estimatedTotalPrice*0.05) + math.Ceil(0.025*clientDeliveryPrice) + (calculatedDeliveryHistoryPricesSum - clientDeliveryPrice)
}

func (or *OrderReportImpl) balanceForRestaurantIOKA(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum float64) float64 {
	return estimatedTotalPrice - math.Ceil(estimatedTotalPrice*0.05) - math.Ceil(0.029*clientDeliveryPrice) - (calculatedDeliveryHistoryPricesSum - clientDeliveryPrice)
}

func (or *OrderReportImpl) balanceForRestaurantKaspi(estimatedTotalPrice, calculatedDeliveryHistoryPricesSum float64) float64 {
	return math.Ceil(-0.03*estimatedTotalPrice) - calculatedDeliveryHistoryPricesSum
}

func (or *OrderReportImpl) balanceForRestaurantKaspiSalesCount(estimatedTotalPrice, clientDeliveryPrice, calculatedDeliveryHistoryPricesSum float64) float64 {
	return estimatedTotalPrice - math.Ceil(estimatedTotalPrice*0.05) - math.Ceil(0.025*clientDeliveryPrice) - (calculatedDeliveryHistoryPricesSum - clientDeliveryPrice)
}

func (or *OrderReportImpl) getBankBalance(paymentSystem string, estimatedTotalPrice, clientDeliveryPrice float64) float64 {
	switch paymentSystem {
	case iokaPaymentSystem:
		return or.getIokaBalanceFieldForKwaaka(estimatedTotalPrice, clientDeliveryPrice)
	case kaspiSalesCountPaymentSystem:
		return or.getKaspiSalesCountBalanceFieldForKwaaka(estimatedTotalPrice, clientDeliveryPrice)
	}
	return 0
}

func (or *OrderReportImpl) getIokaBalanceFieldForKwaaka(estimatedTotalPrice, clientDeliveryPrice float64) float64 {
	return (estimatedTotalPrice + clientDeliveryPrice) * 0.029
}

func (or *OrderReportImpl) getKaspiSalesCountBalanceFieldForKwaaka(estimatedTotalPrice, clientDeliveryPrice float64) float64 {
	return (estimatedTotalPrice + clientDeliveryPrice) * 0.025
}

func (or *OrderReportImpl) isFailedOrder(status string) bool {
	if status == string(models.STATUS_CANCELLED) || status == string(models.STATUS_DELIVERED) || status == string(models.STATUS_FAILED) || status == string(models.STATUS_CANCELLED_BY_POS_SYSTEM) || status == string(models.STATUS_CANCELLED_BY_DELIVERY_SERVICE) || status == string(models.STATUS_SKIPPED) {
		return true
	}
	return false
}

func (or *OrderReportImpl) reportToXlsx(reportsRest []models.OrderReport, totalsForRest, totalsForKwaaka models.OrderReportResponse) ([]byte, error) {

	layout := "02-01-2006 15:04:05"
	sort.SliceStable(reportsRest, func(i, j int) bool {
		timeI, errI := time.Parse(layout, reportsRest[i].OrderTime)
		timeJ, errJ := time.Parse(layout, reportsRest[j].OrderTime)
		if errI != nil || errJ != nil {
			log.Error().Msgf("error while parsing time: %v, %v", errI, errJ)
			return false
		}
		return timeI.Before(timeJ)
	})

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("OrderReportForRest")
	if err != nil {
		return nil, err
	}

	headersRest := []string{
		"ID заказа",
		"Название ресторана",
		"Название ресторан группы",
		"Агрегатор",
		"Время заказа",
		"Тип заказа",
		"Статус заказа",
		"Статус доставки",
		"Статусы доставок (история)",
		"ID Доставок (история)",

		"Адрес доставки",
		"Адрес ресторана",

		"Провайдер",
		"Провайдеры (история)",
		"Комментарии",
		"Способ оплаты",
		"Стоимость товаров",
		"Общая стоимость заказа",
		"Имя клиента",
		"Номер клиента",
		"Продукты",

		"Баланс взаиморасчетов с рестораном",
		"Итоговый Заработок ресторана",
		"Баланс взаиморасчетов с Kwaaka",
		"Итоговый Заработок Kwaaka",

		"Стоимость доставок (прогнозируемые / история)", //
		"Сумма стоимости доставок (прогнозируемые / история)",

		"Стоимость доставок (фактические / история)",
		"Сумма стоимости доставок (фактические / история)",

		"Стоимость доставок (расчетные / история)",
		"Сумма стоимости доставок (расчетные / история)",

		"Стоимость доставки для клиента",
		"Kwaaka Charge (markup)",
		"Баланс взаиморасчетов с банком",
		"Баланс взаиморасчетов с провайдером"}

	mainRowRest := sheet.AddRow()
	for _, header := range headersRest {
		mainRowRest.AddCell().Value = header
	}

	for _, report := range reportsRest {
		row := sheet.AddRow()

		row.AddCell().Value = report.OrderID
		row.AddCell().Value = report.RestaurantName
		row.AddCell().Value = report.RestaurantGroupName
		row.AddCell().Value = report.Source
		row.AddCell().Value = report.OrderTime
		row.AddCell().Value = report.OrderType
		row.AddCell().Value = report.OrderStatus
		row.AddCell().Value = report.DeliveryStatus
		row.AddCell().Value = report.DeliveryStatusHistory
		row.AddCell().Value = report.DeliveryOrderHistoryIDs

		row.AddCell().Value = report.DeliveryAddress
		row.AddCell().Value = report.RestaurantAddress

		row.AddCell().Value = report.DeliveryDispatcher
		row.AddCell().Value = report.DeliveryOrderProviderHistory
		row.AddCell().Value = report.OrderComment
		row.AddCell().Value = report.PaymentSystem
		row.AddCell().Value = strconv.FormatFloat(report.EstimatedTotalPrice, 'f', -1, 64)
		row.AddCell().Value = strconv.FormatFloat(report.TotalOrderPrice, 'f', -1, 64)
		row.AddCell().Value = report.CustomerName
		row.AddCell().Value = report.CustomerPhone

		var productStr string
		for _, product := range report.Products {
			for _, attribute := range product.Attributes {
				productStr = fmt.Sprintf("%s, цена: %.2f, количество: %d\n Аттрибуты:\n %s, цена: %.2f, количество: %d\n", product.Name, product.Price.Value, product.Quantity, attribute.Name, attribute.Price.Value, attribute.Quantity)
			}
		}
		row.AddCell().Value = productStr                                                                    // 21
		row.AddCell().Value = strconv.FormatFloat(math.Ceil(report.Numbers.BalanceRestaurant), 'f', -1, 64) // Баланс взаиморасчетов с рестораном
		row.AddCell().Value = strconv.FormatFloat(math.Ceil(report.Numbers.RestaurantIncome), 'f', -1, 64)  // Итоговый Заработок ресторана
		row.AddCell().Value = strconv.FormatFloat(math.Ceil(report.Numbers.BalanceKwaaka), 'f', -1, 64)     // Баланс взаиморасчетов с Kwaaka
		row.AddCell().Value = strconv.FormatFloat(math.Ceil(report.Numbers.KwaakaIncome), 'f', -1, 64)      // Итоговый Заработок Kwaaka

		row.AddCell().Value = report.Numbers.ProjectedDeliveryHistoryPrices // Стоимость доставок (прогнозируемые) (история)
		row.AddCell().Value = strconv.FormatFloat(report.Numbers.ProjectedDeliveryHistoryPricesSUM, 'f', -1, 64)

		row.AddCell().Value = report.Numbers.ActualDeliveryHistoryPrices
		row.AddCell().Value = strconv.FormatFloat(report.Numbers.ActualDeliveryHistoryPricesSUM, 'f', -1, 64)

		row.AddCell().Value = report.Numbers.CalculatedDeliveryHistoryPrices
		row.AddCell().Value = strconv.FormatFloat(report.Numbers.CalculatedDeliveryHistoryPricesSUM, 'f', -1, 64) // Сумма стоимости доставок (расчетные / история)

		row.AddCell().Value = strconv.FormatFloat(math.Ceil(report.Numbers.ClientDeliveryPrice), 'f', -1, 64) // Стоимость доставки для клиента
		row.AddCell().Value = strconv.FormatFloat(report.Numbers.KwaakaChargedDeliveryPrice, 'f', -1, 64)
		row.AddCell().Value = strconv.FormatFloat(math.Ceil(report.Numbers.BankBalance), 'f', -1, 64)
		row.AddCell().Value = strconv.FormatFloat(math.Ceil(report.Numbers.DeliveryBalance), 'f', -1, 64)
	}

	sheet.AddRow()
	totalsRest := []string{"Общее количество заказов", "Общая стоимость заказов за выбранный период времени", "Итоговый заработок ресторана за выбранный период времени", "Баланс взаиморасчетов с Kwaaka за выбранный период времени"}
	totalsRowRest := sheet.AddRow()
	for _, total := range totalsRest {
		totalsRowRest.AddCell().Value = total
	}

	rowRest := sheet.AddRow()
	rowRest.AddCell().Value = strconv.Itoa(totalsForRest.TotalOrdersCount)
	rowRest.AddCell().Value = strconv.FormatFloat(math.Ceil(totalsForRest.OrdersTotalPrice), 'f', -1, 64)
	rowRest.AddCell().Value = strconv.FormatFloat(math.Ceil(totalsForRest.TotalIncome), 'f', -1, 64)
	rowRest.AddCell().Value = strconv.FormatFloat(math.Ceil(totalsForRest.TotalBalance), 'f', -1, 64)

	sheet.AddRow()
	totalsKwaaka := []string{"Общее количество заказов", "Общая стоимость заказов за выбранный период времени", "Итоговый заработок Kwaaka за выбранный период времени", "Баланс взаиморасчетов с рестораном за выбранный период времени", "Баланс взаиморасчетов с сервисами доставки за выбранный период времени"}
	totalsRow := sheet.AddRow()
	for _, total := range totalsKwaaka {
		totalsRow.AddCell().Value = total
	}

	rowKwaaka := sheet.AddRow()
	rowKwaaka.AddCell().Value = strconv.Itoa(totalsForKwaaka.TotalOrdersCount)
	rowKwaaka.AddCell().Value = strconv.FormatFloat(math.Ceil(totalsForKwaaka.OrdersTotalPrice), 'f', -1, 64)
	rowKwaaka.AddCell().Value = strconv.FormatFloat(math.Ceil(totalsForKwaaka.TotalIncome), 'f', -1, 64)
	rowKwaaka.AddCell().Value = strconv.FormatFloat(math.Ceil(totalsForKwaaka.TotalBalance), 'f', -1, 64)
	rowKwaaka.AddCell().Value = strconv.FormatFloat(math.Ceil(totalsForKwaaka.TotalDeliveryBalanceForKwaaka), 'f', -1, 64)

	var b bytes.Buffer
	if err := file.Write(&b); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (or *OrderReportImpl) OrderReportToXlsx(ctx context.Context, query models.OrderReportRequest) ([]byte, error) {
	restaurantReport, err := or.OrderReportForRestaurant(ctx, query)
	if err != nil {
		return nil, err
	}

	totalsForRest, err := or.OrderReportForRestaurantTotals(ctx, query)
	if err != nil {
		return nil, err
	}

	totalsForKwaaka, err := or.OrderReportForKwaakaTotals(ctx, query)
	if err != nil {
		return nil, err
	}

	file, err := or.reportToXlsx(restaurantReport.OrdersReport, totalsForRest, totalsForKwaaka)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (or *OrderReportImpl) DeliveryDispatcherPrice(ctx context.Context) error {
	log.Info().Msgf("delivery dispatcher price")
	orders, _, err := or.repository.GetAllOrders(ctx, selector.EmptyOrderSearch().
		SetOrderTimeFrom(time.Now().AddDate(0, 0, -2)).
		SetDeliveryServices([]string{directDeliveryService, kwaakaAdminDeliveryService}))
	if err != nil {
		log.Err(err).Msgf("get order from db error")
		return err
	}

	log.Info().Msgf("delivery dispatcher price getting orders from db, orders count: %d", len(orders))

	var deliveryIDs []string
	for _, order := range orders {
		for _, canceledOrder := range order.History3plDeliveryInfo {
			deliveryIDs = append(deliveryIDs, canceledOrder.DeliveryOrderID)
		}
		deliveryIDs = append(deliveryIDs, order.DeliveryOrderID)
	}

	deliveryDispatcherPrices, err := or.kwaaka3pl.GetDeliveryDispatcherPrices(ctx, deliveryIDs)
	if err != nil {
		log.Err(err).Msgf("get delivery dispathcer prices error")
		return err
	}

	log.Info().Msgf("delivery dispatcher price getting prices from 3pl, count: %d, prices: %+v", len(deliveryDispatcherPrices.DeliveryPrices), deliveryDispatcherPrices)

	deliveryIDToDeliveryPriceMap := make(map[string]float64, len(deliveryDispatcherPrices.DeliveryPrices))
	for _, deliveryPrice := range deliveryDispatcherPrices.DeliveryPrices {
		deliveryIDToDeliveryPriceMap[deliveryPrice.DeliveryID] = deliveryPrice.Price
	}

	for _, order := range orders {
		price, ok := deliveryIDToDeliveryPriceMap[order.DeliveryOrderID]
		if !ok {
			continue
		}

		if order.DeliveryDispatcherPrice == 0 {
			if err := or.repository.SetDeliveryDispatcherPrice(ctx, order.ID, price); err != nil {
				log.Err(err).Msgf("set delivery dispatcher price error, order id: %s", order.OrderID)
				continue
			}
		}
		log.Info().Msgf("update delivery dispatcher price success, order id: %s", order.OrderID)
	}

	for _, order := range orders {
		for _, canceledOrder := range order.History3plDeliveryInfo {
			canceledDeliveryPrice, ok := deliveryIDToDeliveryPriceMap[canceledOrder.DeliveryOrderID]
			if !ok {
				continue
			}
			if canceledOrder.DeliveryDispatcherPrice == 0 {
				if err := or.repository.SetCancelledDeliveryDispatcherPrice(ctx, order.ID, canceledDeliveryPrice); err != nil {
					log.Err(err).Msgf("set canceled delivery dispatcher price error, order id: %s", order.OrderID)
					continue
				}
			}
			log.Info().Msgf("update canceled delivery dispatcher price success, order id: %s", order.OrderID)
		}
	}

	return nil
}

func (or *OrderReportImpl) GetPaymentSystem(ctx context.Context, deliveryService, orderID string) (string, error) {
	var paymentSystem string
	switch deliveryService {
	case models.QRMENU.String():
		cart, err := or.cartService.GetQRMenuCartByID(ctx, orderID)
		if err != nil {
			cart = models.Cart{PaymentSystem: noInfo}
		}
		paymentSystem = cart.PaymentSystem
		if paymentSystem == "" {
			cart, err := or.cartService.GetOldQRMenuCartByID(ctx, orderID)
			if err != nil {
				cart = models.OldCart{PaymentSystem: noInfo}
			}
			paymentSystem = cart.PaymentSystem
		}
	case models.KWAAKA_ADMIN.String():
		cart, err := or.cartService.GetKwaakaAdminCartByOrderID(ctx, orderID)
		if err != nil {
			cart = models.Cart{PaymentSystem: noInfo}
		}
		paymentSystem = cart.PaymentType
	}
	return paymentSystem, nil
}

func (or *OrderReportImpl) GetOrderType(orderType string, sendCourier bool) string {
	switch {
	case orderType == models.ORDER_TYPE_INSTANT && sendCourier:
		orderType = "Доставка"
	case orderType == models.ORDER_TYPE_INSTANT && !sendCourier:
		orderType = "Самовывоз"
	case orderType == models.ORDER_TYPE_PREORDER && sendCourier:
		orderType = "Предзаказ & Доставка"
	case orderType == models.ORDER_TYPE_PREORDER && !sendCourier:
		orderType = "Предзаказ & Самовывоз"
	}
	return orderType
}

func (or *OrderReportImpl) GetSource(deliveryService string) string {
	switch deliveryService {
	case models.QRMENU.String():
		return "Kwaaka Direct"
	case models.KWAAKA_ADMIN.String():
		return "Kwaaka Call-Center"
	}
	return deliveryService
}

func (or *OrderReportImpl) GetOrderDeliveryStatus(ctx context.Context, orderID string) (string, error) {
	deliveryInfo, err := or.kwaaka3pl.GetDeliveryInfoForReport(ctx, orderID)
	if err != nil {
		return "", err
	}
	if deliveryInfo.Statuses != nil && len(deliveryInfo.Statuses) != 0 {
		return deliveryInfo.Statuses[len(deliveryInfo.Statuses)-1].Status, nil
	}
	return "", err
}

func (or *OrderReportImpl) GetDeliveryPrice(ctx context.Context, deliveryOrderID string) (float64, error) {
	var (
		deliveryPrice float64
		err           error
	)
	switch len([]rune(deliveryOrderID)) {
	case 1, 2, 3, 4, 5:
		deliveryPrice, err = or.kwaaka3pl.GetDeliveryPriceFromPQ(ctx, deliveryOrderID)
		if err != nil {
			return 0, err
		}
	default:
		deliveryPrice, err = or.kwaaka3pl.GetDeliveryPrice(ctx, deliveryOrderID)
		if err != nil {
			return 0, err
		}
	}
	return deliveryPrice, nil
}

func (or *OrderReportImpl) GetHistoricalDeliveriesInfo(ctx context.Context, order models.Order, deliveryStatus string, currentDeliveryDispatcherPrice float64,
	history3plDelivery []models.History3plDelivery, canceled3plDelivery []models.Cancelled3PLDelivery) ([]models.DeliveryHistory, string, error) {

	var deliveryHistoryInfos = make([]models.DeliveryHistory, 0, len(history3plDelivery)+len(canceled3plDelivery)+1)

	switch {
	case deliveryStatus == "no info" && len(history3plDelivery) == 0 && len(canceled3plDelivery) == 0:
		return deliveryHistoryInfos, "", nil
	}

	if len(history3plDelivery) != 0 && history3plDelivery != nil {
		for _, history := range history3plDelivery {
			status, err := or.kwaaka3pl.GetDeliveryStatus(ctx, history.DeliveryOrderID)
			if err != nil {
				log.Error().Msgf("couldn't get status by delivery order ID")
				continue
			}
			deliveryPrice, err := or.GetDeliveryPrice(ctx, history.DeliveryOrderID)
			if err != nil {
				log.Error().Msgf("couldn't get delivery price by delivery order id")
				continue
			}

			deliveryHistoryInfos = append(deliveryHistoryInfos, models.DeliveryHistory{
				DeliveryOrderID:            history.DeliveryOrderID,
				DeliveryDispatcher:         history.DeliveryDispatcher,
				DeliveryDispatcherPrice:    math.Ceil(deliveryPrice),
				FullDeliveryPrice:          history.FullDeliveryPrice,
				RestaurantPayDeliveryPrice: history.RestaurantPayDeliveryPrice,
				KwaakaChargedDeliveryPrice: history.KwaakaChargedDeliveryPrice,
				Status:                     status,
			})
		}
	}
	if len(canceled3plDelivery) != 0 && canceled3plDelivery != nil {
		for _, canceled := range canceled3plDelivery {
			status, err := or.kwaaka3pl.GetDeliveryStatus(ctx, canceled.DeliveryOrderID)
			if err != nil {
				log.Error().Msgf("couldn't get status by delivery order ID")
				continue
			}

			deliveryPrice, err := or.GetDeliveryPrice(ctx, canceled.DeliveryOrderID)
			if err != nil {
				log.Error().Msgf("couldn't get delivery price by delivery order id")
				continue
			}

			deliveryHistoryInfos = append(deliveryHistoryInfos, models.DeliveryHistory{
				DeliveryOrderID:            canceled.DeliveryOrderID,
				DeliveryDispatcher:         canceled.DeliveryDispatcher,
				DeliveryDispatcherPrice:    math.Ceil(deliveryPrice),
				FullDeliveryPrice:          canceled.FullDeliveryPrice,
				RestaurantPayDeliveryPrice: canceled.RestaurantPayDeliveryPrice,
				KwaakaChargedDeliveryPrice: canceled.KwaakaChargedDeliveryPrice,
				Status:                     status,
			})
		}
	}

	if order.FullDeliveryPrice != 0 && deliveryStatus != "no info" {
		deliveryHistoryInfos = append(deliveryHistoryInfos, models.DeliveryHistory{
			DeliveryOrderID:            order.DeliveryOrderID,
			DeliveryDispatcher:         order.DeliveryDispatcher,
			DeliveryDispatcherPrice:    currentDeliveryDispatcherPrice,
			FullDeliveryPrice:          order.FullDeliveryPrice,
			RestaurantPayDeliveryPrice: order.RestaurantPayDeliveryPrice,
			KwaakaChargedDeliveryPrice: order.KwaakaChargedDeliveryPrice,
			Status:                     deliveryStatus,
		})
	}

	var newDeliveryStatus string

	if len(deliveryHistoryInfos) != 0 && deliveryHistoryInfos != nil && deliveryStatus == "no info" {
		for i := len(deliveryHistoryInfos) - 1; i >= 0; i-- {
			if deliveryHistoryInfos[i].Status != "" {
				newDeliveryStatus = deliveryHistoryInfos[i].Status
				break
			}
		}
	}

	return deliveryHistoryInfos, newDeliveryStatus, nil
}

func (or *OrderReportImpl) getProjectedDeliveryHistoryPrices(deliveryHistoryInfos []models.DeliveryHistory) string {
	if len(deliveryHistoryInfos) == 0 {
		return ""
	}
	var predictableDeliveryPrices = make([]string, 0, len(deliveryHistoryInfos))

	for _, delivery := range deliveryHistoryInfos {
		predictableDeliveryPrices = append(predictableDeliveryPrices, strconv.FormatFloat(delivery.FullDeliveryPrice, 'f', -1, 64))
	}
	return strings.Join(predictableDeliveryPrices, ", ")
}

func (or *OrderReportImpl) getProjectedDeliveryHistoryPricesSum(deliveryHistoryInfos []models.DeliveryHistory) float64 {
	if len(deliveryHistoryInfos) == 0 {
		return 0
	}
	var sum float64

	for _, delivery := range deliveryHistoryInfos {
		sum += delivery.FullDeliveryPrice
	}
	return math.Ceil(sum)
}

func (or *OrderReportImpl) getActualDeliveryHistoryPrices(deliveryHistoryInfos []models.DeliveryHistory) string {
	if len(deliveryHistoryInfos) == 0 {
		return ""
	}
	var actualDeliveryPrices = make([]string, 0, len(deliveryHistoryInfos))

	for _, delivery := range deliveryHistoryInfos {
		if delivery.DeliveryDispatcherPrice == 0 {
			actualDeliveryPrices = append(actualDeliveryPrices, strconv.FormatFloat(delivery.DeliveryDispatcherPrice, 'f', -1, 64))
			continue
		}
		actualDeliveryPrices = append(actualDeliveryPrices, strconv.FormatFloat(math.Ceil(delivery.DeliveryDispatcherPrice+delivery.KwaakaChargedDeliveryPrice), 'f', -1, 64))
	}
	return strings.Join(actualDeliveryPrices, ", ")
}

func (or *OrderReportImpl) getActualDeliveryHistoryPricesSUM(deliveryHistoryInfos []models.DeliveryHistory) float64 {
	if len(deliveryHistoryInfos) == 0 {
		return 0
	}
	var sum float64

	for _, delivery := range deliveryHistoryInfos {
		if delivery.DeliveryDispatcherPrice == 0 {
			continue
		}
		sum += delivery.DeliveryDispatcherPrice + delivery.KwaakaChargedDeliveryPrice
	}
	return math.Ceil(sum)
}

func (or *OrderReportImpl) getCalculatedDeliveryHistoryPrices(deliveryHistoryInfos []models.DeliveryHistory) string {
	if len(deliveryHistoryInfos) == 0 {
		return ""
	}
	var actualDeliveryPrices = make([]string, 0, len(deliveryHistoryInfos))

	for _, delivery := range deliveryHistoryInfos {
		switch {
		case delivery.DeliveryDispatcherPrice == 0:
			actualDeliveryPrices = append(actualDeliveryPrices, strconv.FormatFloat(math.Ceil(delivery.DeliveryDispatcherPrice), 'f', -1, 64))
			continue
		case delivery.FullDeliveryPrice > delivery.DeliveryDispatcherPrice+delivery.KwaakaChargedDeliveryPrice:
			actualDeliveryPrices = append(actualDeliveryPrices, strconv.FormatFloat(math.Ceil(delivery.FullDeliveryPrice), 'f', -1, 64))
		case delivery.FullDeliveryPrice < delivery.DeliveryDispatcherPrice+delivery.KwaakaChargedDeliveryPrice:
			actualDeliveryPrices = append(actualDeliveryPrices, strconv.FormatFloat(math.Ceil(delivery.DeliveryDispatcherPrice+delivery.KwaakaChargedDeliveryPrice), 'f', -1, 64))
		}
	}
	return strings.Join(actualDeliveryPrices, ", ")
}

func (or *OrderReportImpl) getCalculatedDeliveryHistoryPricesSUM(deliveryHistoryInfos []models.DeliveryHistory) float64 {
	if len(deliveryHistoryInfos) == 0 {
		return 0
	}
	var sum float64

	for _, delivery := range deliveryHistoryInfos {
		switch {
		case delivery.DeliveryDispatcherPrice == 0:
			continue
		case delivery.FullDeliveryPrice > delivery.DeliveryDispatcherPrice+delivery.KwaakaChargedDeliveryPrice:
			sum += delivery.FullDeliveryPrice
		case delivery.FullDeliveryPrice < delivery.DeliveryDispatcherPrice+delivery.KwaakaChargedDeliveryPrice:
			sum += delivery.DeliveryDispatcherPrice + delivery.KwaakaChargedDeliveryPrice
		}
	}
	return math.Ceil(sum)
}

func (or *OrderReportImpl) getDispatcherDeliveryHistoryPricesSUM(deliveryHistoryInfos []models.DeliveryHistory) float64 {
	if len(deliveryHistoryInfos) == 0 {
		return 0
	}
	var sum float64

	for _, delivery := range deliveryHistoryInfos {
		sum += delivery.DeliveryDispatcherPrice
	}
	return math.Ceil(sum)
}

func (or *OrderReportImpl) getDeliveryHistoryStatuses(deliveryHistoryInfos []models.DeliveryHistory) string {
	statuses := make([]string, 0, len(deliveryHistoryInfos))
	for _, deliveryInfo := range deliveryHistoryInfos {
		statuses = append(statuses, deliveryInfo.Status)
	}
	return strings.Trim(strings.Join(statuses, ", "), ", ")
}

func (or *OrderReportImpl) getDeliveryHistoryIDs(deliveryHistoryInfos []models.DeliveryHistory) string {
	ids := make([]string, 0, len(deliveryHistoryInfos))
	for _, deliveryInfo := range deliveryHistoryInfos {
		ids = append(ids, deliveryInfo.DeliveryOrderID)
	}
	return strings.Join(ids, ", ")
}

func (or *OrderReportImpl) getDeliveryHistoryProviders(deliveryHistoryInfos []models.DeliveryHistory) string {
	providers := make([]string, 0, len(deliveryHistoryInfos))
	for _, deliveryInfo := range deliveryHistoryInfos {
		providers = append(providers, deliveryInfo.DeliveryDispatcher)
	}
	return strings.Join(providers, ", ")
}
