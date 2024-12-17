package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/posist/clients/models"
)

func (p *posist) GetStopList(ctx context.Context, customerKey string) ([]models.Item, error) {
	path := fmt.Sprintf("/api/v1/online_order/menu?tabtype=delivery&customer_key=%s", customerKey)

	var (
		result []models.Item
	)

	response, err := p.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetResult(&result).
		Get(path)
	if err != nil {
		return nil, err
	}

	if response.IsError() || response.StatusCode() >= 400 {
		return nil, fmt.Errorf("get out of stock items error %v", response.Error())
	}

	return result, nil
}
