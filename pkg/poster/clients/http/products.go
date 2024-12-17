package http

import (
	"context"
	"fmt"
	models2 "github.com/kwaaka-team/orders-core/pkg/poster/clients/models"
)

func (c *Client) GetProduct(ctx context.Context, productId string) (models2.GetProductsResponseBody, error) {
	path := "/api/menu.getProduct"

	var (
		response    models2.GetProductResponse
		errResponse models2.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetQueryParam("product_id", productId).
		Get(path)
	if err != nil {
		return models2.GetProductsResponseBody{}, fmt.Errorf("%v + %v", err, response)
	}

	if resp.IsError() {
		return models2.GetProductsResponseBody{}, errResponse
	}

	return response.Response, nil
}
func (c *Client) GetProducts(ctx context.Context) (models2.GetProductsResponse, error) {
	path := "/api/menu.getProducts"

	var (
		response    models2.GetProductsResponse
		errResponse models2.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		//SetQueryParam("type", "products").
		Get(path)
	if err != nil {
		return response, fmt.Errorf("%v + %v", err, response)
	}

	if resp.IsError() {
		return response, errResponse
	}

	if response.Message != "" {
		return response, fmt.Errorf("get products error: %s, status: %d", response.Message, response.Code)
	}

	return response, nil
}
