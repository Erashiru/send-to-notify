package client

import (
	"context"
	"fmt"
	dto2 "github.com/kwaaka-team/orders-core/pkg/jowi/client/dto"
	"github.com/pkg/errors"
)

// GetStopList - getting stop list of products (count = 0 - in stop list, count > 0 - has limited amount)
func (c *Client) GetStopList(ctx context.Context, restaurantID string) (dto2.ResponseStopList, error) {
	path := fmt.Sprintf("/v3/restaurants/%s/course_counts", restaurantID)
	// FIXME JOWI docs has mistakes restayrabts...

	var response dto2.ResponseStopList
	var errResponse dto2.ErrorResponse

	// FIXME resty doesn't send request with body (Content-Type=application/json) if method is equal GET
	result, err := c.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetQueryParams(map[string]string{
			"api_key": c.apiKey,
			"sig":     c.sig,
			// TODO: restaurant_id as param?
		}).
		SetResult(&response).
		SetError(&errResponse).
		Get(path)
	if err != nil {
		return dto2.ResponseStopList{}, err
	}

	if result.IsError() {
		return dto2.ResponseStopList{}, errors.New(errResponse.Message)
	}

	// FIXME: JOWI returns 200 even if response has errors, so decided to add ErrorResponse to ResponseModel
	if response.Status != 1 {
		return dto2.ResponseStopList{}, errors.New(response.Message)
	}

	return response, nil
}
