package models

import "time"

type OrderReportRequest struct {
	RestaurantIDs      []string  `json:"restaurant_ids"`
	RestaurantGroupIds []string  `json:"restaurant_group_ids"`
	DeliveryDispatcher string    `json:"delivery_dispatcher"`
	Search             string    `json:"search"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	PaymentSystem      string    `json:"payment_system"`
	Pagination
}

type OrderReportResponse struct {
	OrdersReport                  []OrderReport `json:"orders_report"`
	OrdersTotalPrice              float64       `json:"orders_total_price"`
	TotalIncome                   float64       `json:"total_income"`
	TotalBalance                  float64       `json:"total_balance"`
	TotalDeliveryBalanceForKwaaka float64       `json:"total_delivery_balance_for_kwaaka"`
	TotalOrdersCount              int           `json:"total_orders_count"`
}

type OrderReport struct {
	OrderID                      string             `json:"order_id"`
	RestaurantID                 string             `json:"restaurant_id"`
	RestaurantName               string             `json:"restaurant_name"`
	RestaurantGroupID            string             `json:"restaurant_group_id"`
	RestaurantGroupName          string             `json:"restaurant_group_name"`
	Source                       string             `json:"source"`
	OrderTime                    string             `json:"order_time"`
	EstimatedPickupTime          time.Time          `json:"estimated_pickup_time"`
	OrderType                    string             `json:"order_type"`
	OrderStatus                  string             `json:"order_status"`
	DeliveryStatus               string             `json:"delivery_status"`
	DeliveryStatusHistory        string             `json:"delivery_status_history"`
	DeliveryOrderHistoryIDs      string             `json:"delivery_order_history_ids"`
	DeliveryAddress              string             `json:"delivery_address"`
	RestaurantAddress            string             `json:"restaurant_address"`
	DeliveryOrderProviderHistory string             `json:"delivery_order_provider_history"`
	DeliveryDispatcher           string             `json:"delivery_dispatcher"`
	OrderComment                 string             `json:"order_comment"`
	PaymentSystem                string             `json:"payment_system"`
	EstimatedTotalPrice          float64            `json:"estimated_total_price"`
	TotalOrderPrice              float64            `json:"total_order_price"` // Общая стоимость заказа
	CustomerName                 string             `json:"customer_name"`
	CustomerPhone                string             `json:"customer_phone"`
	SendCourier                  bool               `json:"send_courier"`
	ReportReady                  bool               `json:"report_ready"`
	Products                     []OrderProduct     `json:"products"`
	Numbers                      OrderReportNumbers `json:"for_restaurant"`
}

type OrderReportNumbers struct {
	RestaurantIncome  float64 `json:"restaurant_income"` // Итоговый Заработок ресторана
	KwaakaIncome      float64 `json:"income"`            // Итоговый Заработок Kwaaka
	BalanceKwaaka     float64 `json:"balance"`           // Баланс взаиморасчетов с Kwaaka
	BalanceRestaurant float64 `json:"balance_restaurant"`

	ProjectedDeliveryHistoryPrices    string  `json:"projected_delivery_history_prices"`     // Стоимость доставок (прогнозируемые)
	ProjectedDeliveryHistoryPricesSUM float64 `json:"projected_delivery_history_prices_sum"` // Сумма стоимости доставок (прогнозируемые - СУММА)

	ActualDeliveryHistoryPrices    string  `json:"actual_delivery_history_prices"`     // Стоимость доставок (фактические)
	ActualDeliveryHistoryPricesSUM float64 `json:"actual_delivery_history_prices_sum"` // Сумма стоимости доставок (фактических - Сумма)

	CalculatedDeliveryHistoryPrices    string  `json:"calculated_delivery_history_prices"`     // Стоимость доставок (расчетные)
	CalculatedDeliveryHistoryPricesSUM float64 `json:"calculated_delivery_history_prices_sum"` // Сумма стоимости доставок (расчетных - Сумма)

	RestaurantDeliveryPrice float64 `json:"restaurant_delivery_price"` // Стоимость доставки для ресторана (фактическая)
	ClientDeliveryPrice     float64 `json:"client_delivery_price"`     // Стоимость доставки для клиента (фактическая)

	KwaakaChargedDeliveryPrice float64 `json:"kwaaka_charged_delivery_price"` // Kwaaka Charge (markup)
	BankBalance                float64 `json:"bank_balance"`                  //Баланс взаиморасчетов с банком
	DeliveryBalance            float64 `json:"delivery_balance"`              // Баланс взаиморасчетов с провайдером
}
