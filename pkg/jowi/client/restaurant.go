package client

import (
	"context"
	dto2 "github.com/kwaaka-team/orders-core/pkg/jowi/client/dto"
	"github.com/pkg/errors"
)

// GetRestaurants - getting integrated restaurants
func (c *Client) GetRestaurants(ctx context.Context) (dto2.ResponseListRestaurants, error) {
	path := "/v3/restaurants"

	var response dto2.ResponseListRestaurants
	var errResponse dto2.ErrorResponse

	// FIXME resty doesn't send request with body (Content-Type=application/json) if method is equal GET
	result, err := c.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetQueryParams(map[string]string{
			"api_key": c.apiKey,
			"sig":     c.sig,
		}).
		SetResult(&response).
		SetError(&errResponse).
		Get(path)
	if err != nil {
		return dto2.ResponseListRestaurants{}, err
	}

	if result.IsError() {
		return dto2.ResponseListRestaurants{}, errors.New(errResponse.Message)
	}

	// FIXME: JOWI returns 200 even if response has errors, so decided to add ErrorResponse to ResponseModel
	if response.Status != 1 {
		return dto2.ResponseListRestaurants{}, errors.New(response.Message)
	}

	return response, nil
}
