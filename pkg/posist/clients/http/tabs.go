package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/posist/clients/models"
)

func (p *posist) GetTabs(ctx context.Context, customerKey string) ([]models.Tab, error) {
	path := fmt.Sprintf("/api/v1/online_order/deployment_tabs?customer_key=%s", customerKey)

	var (
		result []models.Tab
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
		return nil, fmt.Errorf("get tabs error %v", response.Error())
	}

	return result, nil
}
