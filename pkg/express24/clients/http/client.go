package http

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/core/express24/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/express24/clients"
	dto2 "github.com/kwaaka-team/orders-core/pkg/express24/clients/dto"
	"github.com/pkg/errors"
	"time"
)

type Client struct {
	restyClient        *resty.Client
	BaseUrl            string
	Username, Password string
	quit               chan struct{}
}

func (c *Client) UpdateProducts(ctx context.Context, req dto2.UpdateProductsRequest) (dto2.UpdateProductsResponse, error) {
	var (
		response     dto2.UpdateProductsResponse
		errResponses dto2.ProductsError
	)

	path := "/api/external/update-products"

	utils.Beautify("Stoplist request send", req)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponses).
		SetResult(&response).
		Put(path)

	if err != nil {
		return dto2.UpdateProductsResponse{}, err
	}

	if resp.IsError() {
		return dto2.UpdateProductsResponse{}, fmt.Errorf("express24 cli: %s", resp.Error())
	}

	utils.Beautify("Stoplist request send result", response)

	return response, nil
}

func (c *Client) UpdateOffers(ctx context.Context, req dto2.UpdateOffersRequest) (dto2.UpdateOffersResponse, error) {
	var (
		response     dto2.UpdateOffersResponse
		errResponses dto2.ProductsError
	)
	path := "/api/external/update-offers"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponses).
		SetResult(&response).
		Put(path)

	if err != nil {
		return dto2.UpdateOffersResponse{}, err
	}

	if resp.IsError() {
		return dto2.UpdateOffersResponse{}, fmt.Errorf("express24 cli: %s", resp.Error())
	}

	return response, nil
}

func (c *Client) GetBranches(ctx context.Context) (dto2.GetBranchesResponse, error) {
	var (
		response     dto2.GetBranchesResponse
		errResponses []dto2.BranchesError
	)

	path := "/api/external/branches"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResponses).
		SetResult(&response).
		Get(path)

	if err != nil {
		return dto2.GetBranchesResponse{}, err
	}

	if resp.IsError() {
		return dto2.GetBranchesResponse{}, fmt.Errorf("express24 cli: %s", err.Error())
	}

	return response, nil
}

func (c *Client) UpdateBranches(ctx context.Context, req dto2.UpdateBranchesRequest) (dto2.UpdateBranchesResponse, error) {
	var (
		response     dto2.UpdateBranchesResponse
		errResponses []dto2.UpdateBranchesError
	)

	path := "/api/external/branches"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponses).
		SetResult(&response).
		Put(path)

	if err != nil {
		return dto2.UpdateBranchesResponse{}, err
	}

	if resp.IsError() {
		return dto2.UpdateBranchesResponse{}, fmt.Errorf("express24 cli: %s", err.Error())
	}

	return response, nil
}
func NewClient(cfg *clients.Config) (clients.Express24, error) {

	if cfg.Username == "" || cfg.Password == "" {
		return nil, errors.New("username or password is not provided")
	}
	if cfg.BaseURL == "" {
		return nil, errors.New("base URL could not be empty")
	}

	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
		}).
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime)

	cl := &Client{
		restyClient: client,
		Username:    cfg.Username,
		Password:    cfg.Password,
		BaseUrl:     cfg.BaseURL,
		quit:        make(chan struct{}),
	}

	if err := cl.Auth(context.Background()); err != nil {
		return nil, ErrAuth
	}

	ticker := time.NewTicker(tokenTimeout)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := cl.Auth(context.Background()); err != nil {
					return
				}
			case <-cl.quit:
				ticker.Stop()
				return
			}
		}
	}()

	return cl, nil
}
