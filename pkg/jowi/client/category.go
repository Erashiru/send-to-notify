package client

import (
	"context"
	"fmt"
	dto2 "github.com/kwaaka-team/orders-core/pkg/jowi/client/dto"
	"github.com/pkg/errors"
)

// GetCourseCategories - getting product categories (has pagination!)
func (c *Client) GetCourseCategories(ctx context.Context, restaurantID string) (dto2.ResponseCourseCategory, error) {
	path := fmt.Sprintf("/v3/restaurants/%s/course_categories", restaurantID)

	var response dto2.ResponseCourseCategory
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
		return dto2.ResponseCourseCategory{}, err
	}

	if result.IsError() {
		return dto2.ResponseCourseCategory{}, errors.New(errResponse.Message)
	}

	// FIXME: JOWI returns 200 even if response has errors, so decided to add ErrorResponse to ResponseModel
	if response.Status != 1 {
		return dto2.ResponseCourseCategory{}, errors.New(response.Message)
	}

	return response, nil
}
