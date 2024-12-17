package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/express24_v2/clients"
	"github.com/kwaaka-team/orders-core/pkg/express24_v2/clients/dto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
)

type Client struct {
	restyClient *resty.Client
	BaseUrl     string
	token       string
}

func NewClient(cfg *clients.Config) (clients.Express24V2, error) {

	if cfg.Token == "" {
		return nil, errors.New("token is not provided")
	}
	if cfg.BaseURL == "" {
		return nil, errors.New("base URL could not be empty")
	}

	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
			authHeader:        "Bearer " + cfg.Token,
		}).
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime)

	cl := &Client{
		restyClient: client,
		BaseUrl:     cfg.BaseURL,
	}

	return cl, nil
}

func (c *Client) GetCategories(ctx context.Context) ([]dto.Category, error) {
	var (
		response    []dto.Category
		errResponse dto.ErrorResponse
	)

	path := "external/api/v2/categories"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResponse).
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (dto.CreateCategoryResponse, error) {
	var (
		response    dto.CreateCategoryResponse
		errResponse dto.ErrorResponse
	)

	path := "external/api/v2/categories"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto.CreateCategoryResponse{}, err
	}

	if resp.IsError() {
		return dto.CreateCategoryResponse{}, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) GetSubCategories(ctx context.Context, categoryID string) ([]dto.Category, error) {
	var (
		response    []dto.Category
		errResponse dto.ErrorResponse
	)

	path := fmt.Sprintf("external/api/v2/categories/%s/subs", categoryID)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResponse).
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) CreateSubCategory(ctx context.Context, req dto.CreateCategoryRequest) (dto.CreateCategoryResponse, error) {
	var (
		response    dto.CreateCategoryResponse
		errResponse dto.ErrorResponse
	)

	path := fmt.Sprintf("external/api/v2/categories/%s/subs", req.CategoryID)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto.CreateCategoryResponse{}, err
	}

	if resp.IsError() {
		return dto.CreateCategoryResponse{}, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.CreateProductResponse, error) {
	var (
		response    dto.CreateProductResponse
		errResponse dto.ErrorResponse
	)

	path := "external/api/v2/products"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto.CreateProductResponse{}, err
	}

	if resp.IsError() {
		return dto.CreateProductResponse{}, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) UpdateProduct(ctx context.Context, req dto.UpdateProductRequest) (dto.UpdateProductResponse, error) {
	var (
		response    dto.UpdateProductResponse
		errResponse dto.ErrorResponse
	)

	path := fmt.Sprintf("external/api/v2/products/%s", req.ID)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Patch(path)

	if err != nil {
		return dto.UpdateProductResponse{}, err
	}

	if resp.IsError() {
		return dto.UpdateProductResponse{}, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) GetCategoryProducts(ctx context.Context, categoryID string) ([]dto.GetCategoryProductsResponse, error) {
	var (
		response    []dto.GetCategoryProductsResponse
		errResponse dto.ErrorResponse
	)

	path := fmt.Sprintf("external/api/v2/categories/%s/products?isSub=0", categoryID)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResponse).
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) GetSubCategoryProducts(ctx context.Context, subCategoryID string) ([]dto.GetCategoryProductsResponse, error) {
	var (
		response    []dto.GetCategoryProductsResponse
		errResponse dto.ErrorResponse
	)

	path := fmt.Sprintf("external/api/v2/categories/%s/products?isSub=1", subCategoryID)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResponse).
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) GetProductsAttributeGroups(ctx context.Context, productID string) (dto.GetAttributeGroupResponse, error) {
	var (
		response    dto.GetAttributeGroupResponse
		errResponse dto.ErrorResponse
	)

	path := fmt.Sprintf("external/api/v2/products/%s/mods", productID)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResponse).
		SetResult(&response).
		Get(path)

	if err != nil {
		return dto.GetAttributeGroupResponse{}, err
	}

	if resp.IsError() {
		return dto.GetAttributeGroupResponse{}, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) GetAttributeGroups(ctx context.Context) ([]dto.AttributeGroup, error) {
	var (
		response    []dto.AttributeGroup
		errResponse dto.ErrorResponse
	)

	path := "external/api/v2/modifiers"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResponse).
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) CreateAttributeGroup(ctx context.Context, req dto.CreateProductsAttributeGroupRequest) (dto.AttributeGroup, error) {
	var (
		response    dto.AttributeGroup
		errResponse dto.ErrorResponse
	)

	path := "external/api/v2/modifiers"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto.AttributeGroup{}, err
	}

	if resp.IsError() {
		return dto.AttributeGroup{}, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) CreateProductsAttributeGroup(ctx context.Context, req dto.CreateProductsAttributeGroupRequest) (dto.AttributeGroup, error) {
	var (
		response    dto.AttributeGroup
		errResponse dto.ErrorResponse
	)

	path := fmt.Sprintf("external/api/v2/products/%s/mods", req.ProductID)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto.AttributeGroup{}, err
	}

	if resp.IsError() {
		return dto.AttributeGroup{}, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) CreateAttributeGroupsItem(ctx context.Context, req dto.CreateAttributeGroupsItemRequest) (dto.AttributeGroupItem, error) {
	var (
		response    dto.AttributeGroupItem
		errResponse dto.ErrorResponse
	)

	path := fmt.Sprintf("external/api/v2/modifiers/%s/items", req.AttributeGroupID)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto.AttributeGroupItem{}, err
	}

	if resp.IsError() {
		return dto.AttributeGroupItem{}, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) GetAttributeGroupsItems(ctx context.Context, attributeGroupID string) ([]dto.GetAttributeGroupsItemResponse, error) {
	var (
		response    []dto.GetAttributeGroupsItemResponse
		errResponse dto.ErrorResponse
	)

	path := fmt.Sprintf("external/api/v2/modifiers/%s/items", attributeGroupID)

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResponse).
		SetResult(&response).
		Get(path)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return response, nil
}

func (c *Client) StopListByProducts(ctx context.Context, req dto.StopListByProductsRequest) error {
	var (
		response    dto.UpdateProductResponse
		errResponse dto.ErrorResponse
	)

	path := "external/api/v2/stop-list"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Patch(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) StopListBulk(ctx context.Context, req dto.StopListBulkRequest) error {
	var (
		response    dto.UpdateProductResponse
		errResponse dto.ErrorResponse
	)

	path := "external/api/v2/stop-list"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Patch(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) StopListByAttributes(ctx context.Context, req dto.StopListByAttributesRequest) error {
	var (
		response    dto.UpdateProductResponse
		errResponse dto.ErrorResponse
	)

	path := "external/api/v2/stop-list"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Patch(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return nil
}

func (c *Client) SyncMenu(ctx context.Context, req dto.MenuSyncReq) error {
	var (
		response    dto.MenuSyncResp
		errResponse dto.ErrorResponse
	)

	requestBody, err := json.Marshal(req)
	if err != nil {
		log.Info().Msgf("failed to marshal request body: %s", err)
	}

	err = os.WriteFile("request_body.json", requestBody, 0644)
	if err != nil {
		log.Info().Msgf("failed to write request body to file: %s", err)
	}

	path := "external/api/v2/menu"

	resp, err := c.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResponse).
		SetResult(&response).
		Post(path)
	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("express24 cli: %s", string(resp.Body()))
	}

	return nil
}
