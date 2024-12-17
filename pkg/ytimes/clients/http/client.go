package http

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg"
	"github.com/kwaaka-team/orders-core/pkg/ytimes/clients"
	"github.com/kwaaka-team/orders-core/pkg/ytimes/clients/models"
)

type clientImpl struct {
	restyCli *resty.Client
	token    string
}

func NewClient(cfg clients.Config) *clientImpl {
	restyClient := resty.New().
		SetBaseURL(cfg.BaseUrl).
		SetHeaders(map[string]string{
			pkg.ContentTypeHeader: pkg.JsonType,
			pkg.AuthHeader:        cfg.Token,
		})

	return &clientImpl{
		restyCli: restyClient,
		token:    cfg.Token,
	}
}

func (c clientImpl) GetPoints(ctx context.Context) (models.PointInfo, error) {
	path := "/ex/shop/list"

	var (
		result models.PointInfo
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetResult(&result).
		Get(path)

	if err != nil {
		return models.PointInfo{}, err
	}

	if resp.IsError() {
		return models.PointInfo{}, fmt.Errorf("get points error: %v", resp.Error())
	}

	return result, nil
}

func (c clientImpl) CreateOrder(ctx context.Context, req models.Order) (models.CreateOrderResponse, error) {
	path := "/ex/order/save"

	var (
		result models.CreateOrderResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetBody(&req).
		SetResult(&result).
		Post(path)

	if err != nil {
		return models.CreateOrderResponse{}, err
	}

	if resp.IsError() {
		return models.CreateOrderResponse{}, fmt.Errorf("create order error: %v", resp.Error())
	}

	return result, nil
}

func (c clientImpl) GetMenu(ctx context.Context, pointGuid string) (models.Menu, error) {
	path := "/ex/menu/item/list"

	var (
		result models.Menu
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetQueryParam("shopGuid", pointGuid).
		SetResult(&result).
		Get(path)

	if err != nil {
		return models.Menu{}, err
	}

	if resp.IsError() {
		return models.Menu{}, fmt.Errorf("get menu error: %v", resp.Error())
	}

	return result, nil
}

func (c clientImpl) GetSupplementList(ctx context.Context, pointGuid string) (models.SupplementList, error) {
	path := "/ex/menu/supplement/list"

	var (
		result models.SupplementList
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		SetQueryParam("shopGuid", pointGuid).
		SetResult(&result).
		Get(path)

	if err != nil {
		return models.SupplementList{}, err
	}

	if resp.IsError() {
		return models.SupplementList{}, fmt.Errorf("get supplement list error: %v", resp.Error())
	}

	return result, nil
}
