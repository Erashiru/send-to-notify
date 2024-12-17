package client

import (
	"context"
	"fmt"
	dto2 "github.com/kwaaka-team/orders-core/pkg/jowi/client/dto"
	"github.com/pkg/errors"
	"strconv"
)

// GetCourses - getting list of products (has pagination!)
func (c *Client) GetCourses(ctx context.Context, restaurantID string) (dto2.ResponseCourse, error) {
	path := fmt.Sprintf("/v3/restaurants/%s/courses", restaurantID)
	// TODO: Need add pagination logic

	var result dto2.ResponseCourse
	var errResponse dto2.ErrorResponse

	// FIXME resty doesn't send request with body (Content-Type=application/json) if method is equal GET

	request := c.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetQueryParams(map[string]string{
			"api_key": c.apiKey,
			"sig":     c.sig,
		}).
		SetError(&errResponse)

	pages := 1

	for i := 0; i < pages; i++ {
		var response dto2.ResponseCourse

		resp, err := request.SetQueryParam("page", strconv.Itoa(i+1)).SetResult(&response).Get(path)
		if err != nil {
			return dto2.ResponseCourse{}, err
		}

		if resp.IsError() {
			return dto2.ResponseCourse{}, errors.New(errResponse.Message)
		}

		// FIXME: JOWI returns 200 even if response has errors, so decided to add ErrorResponse to ResponseModel
		if response.Status != 1 {
			return dto2.ResponseCourse{}, errors.New(response.Message)
		}

		result.Courses = append(result.Courses, response.Courses...)
		pages = response.PageCount
	}

	return result, nil
}
