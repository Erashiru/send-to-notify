package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/posist/clients/models"
)

func (p *posist) GetMenu(ctx context.Context, customerKey, tabId string) (models.Menu, error) {
	path := fmt.Sprintf("/api/v1/online_order/standard_menu?customer_key=%s&tabId=%s", customerKey, tabId)

	var (
		result models.Menu
	)

	response, err := p.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetResult(&result).
		Get(path)
	if err != nil {
		return models.Menu{}, err
	}

	if response.IsError() || response.StatusCode() >= 400 {
		return models.Menu{}, fmt.Errorf("get menu error %v", response.Error())
	}

	return result, nil
}
