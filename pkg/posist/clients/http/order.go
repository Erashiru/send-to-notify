package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/posist/clients/models"
)

func (p *posist) CreateOrder(ctx context.Context, customerKey string, order models.Order) error {
	path := fmt.Sprintf("/api/v1/online_order/push?customer_key=%s", customerKey)

	response, err := p.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&order).
		Post(path)
	if err != nil {
		return err
	}

	if response.IsError() || response.StatusCode() >= 400 {
		return fmt.Errorf("create order error %v", response.Error())
	}

	return nil
}

func (p *posist) GetOrderStatus(ctx context.Context, orderId string) (models.OrderStatusResponse, error) {
	path := fmt.Sprintf("/api/v1/online_order/order_status?customer_key=%s&order_id=%s", p.customerKey, orderId)

	var (
		result models.OrderStatusResponse
	)

	response, err := p.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetResult(&result).
		Get(path)
	if err != nil {
		return models.OrderStatusResponse{}, err
	}

	if response.IsError() || response.StatusCode() >= 400 {
		return models.OrderStatusResponse{}, fmt.Errorf("get menu error %v", response.Error())
	}

	return result, nil
}
