package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/moysklad/clients"
	models2 "github.com/kwaaka-team/orders-core/pkg/moysklad/models"
	"github.com/pkg/errors"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	retriesNumber   = 3
	retriesWaitTime = 1 * time.Second
)

const (
	acceptHeader      = "Accept"
	authHeader        = "Authorization"
	contentTypeHeader = "Content-Type"
	jsonType          = "application/json;charset=utf-8"
	organizationURL   = "https://online.moysklad.ru/api/remap/1.2/entity/organization/"
	organizationType  = "organization"
)

type Client struct {
	ApiKey             string
	restyClient        *resty.Client
	StoreID            string
	BaseUrl            string
	Username, Password string
}

func NewClient(cfg *clients.Config) (clients.MoySklad, error) {
	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
		}).
		SetBasicAuth(cfg.Username, cfg.Password).
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime)

	cl := Client{
		restyClient: client,
		BaseUrl:     cfg.BaseURL,
	}
	return cl, nil
}

func (cli Client) GetOrders(ctx context.Context) (models2.Order, error) {
	path := "/api/remap/1.2/entity/customerorder"

	var response models2.Order
	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetResult(&response).
		Get(path)
	if err != nil {
		return models2.Order{}, err
	}

	if resp.IsError() {
		return models2.Order{}, errors.New(resp.Status() + " " + string(resp.Body()))
	}

	return response, nil
}

func (cli Client) GetMenu(ctx context.Context, param models2.GetMenuRequest) (models2.Menu, error) {

	var (
		response models2.Menu
		qParams  = make(map[string]string)
		path     = "/api/remap/1.2/entity/assortment"
	)

	qParams["limit"] = param.Limit
	qParams["offset"] = param.Offset

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		SetQueryParams(qParams).
		EnableTrace().
		SetResult(&response).
		Get(path)
	if err != nil {
		return models2.Menu{}, err
	}

	if resp.IsError() {
		return models2.Menu{}, errors.New(resp.Status() + " " + string(resp.Body()))
	}

	return response, nil
}

func (cli Client) CreateSupplierOrder(ctx context.Context, request models2.SupplierOrder) (string, error) {

	path := "/api/remap/1.2/entity/purchaseorder"
	type Response struct {
		ID string `json:"id"`
	}
	var result Response

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetResult(&result).
		SetBody(request).
		Post(path)
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", errors.New(resp.Status() + " " + string(resp.Body()))
	}
	return result.ID, nil
}

func (cli Client) AddProductSupplier(ctx context.Context, position models2.Position) (string, error) {

	path := fmt.Sprintf("/api/remap/1.2/entity/purchaseorder/%s/positions", position.OrderID)

	type Response struct {
		ID string `json:"id"`
	}
	var result []Response

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetResult(&result).
		SetBody(position).
		Post(path)
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", errors.New(resp.Status() + " " + string(resp.Body()))
	}
	return result[0].ID, nil
}

func (cli Client) DeleteProductSupplier(ctx context.Context, position models2.Position) error {

	path := fmt.Sprintf("/api/remap/1.2/entity/purchaseorder/%s/positions/%s", position.OrderID, position.ProductID)

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		Delete(path)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(resp.Status() + " " + string(resp.Body()))
	}
	return nil
}
