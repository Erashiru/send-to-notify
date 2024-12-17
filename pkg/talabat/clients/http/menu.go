package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/talabat/models"
	"github.com/pkg/errors"
)

func (c *Client) GetRequestStatus(ctx context.Context, requestID string) (models.GetRequestStatusResponse, error) {
	path := fmt.Sprintf("/api/1/Menu/RequestStatus/%s", requestID)

	var (
		response    models.GetRequestStatusResponse
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		Get(path)

	if err != nil {
		return models.GetRequestStatusResponse{}, err
	}

	if resp.IsError() {
		return models.GetRequestStatusResponse{}, errors.New(fmt.Sprintf("httpStatus: %s, type: %s, title: %s, detail: %s, instance: %s", resp.Status(), errResponse.Type, errResponse.Title, errResponse.Detail, errResponse.Instance))
	}

	return response, nil
}

func (c *Client) CreateNewMenu(ctx context.Context, req models.CreateNewMenuRequest) error {
	path := fmt.Sprintf("/api/1/Menu/CreateMenu/%s/%s", req.RestaurantID, req.RequestID)

	var errResponse models.CreateNewMenuErrorResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(req.Menu).
		Post(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(fmt.Sprintf("httpStatus: %s, requestId: %s, message: %s", resp.Status(), errResponse.RequestID, errResponse.Message))
	}

	return nil
}

func (c *Client) UpdateItemsAvailability(ctx context.Context, req models.UpdateItemsAvailabilityRequest) error {
	path := fmt.Sprintf("/api/1/MenuItem/ItemAvailability/%s/%s", req.RestaurantID, req.RequestID)

	var (
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(req).
		Post(path)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(fmt.Sprintf("httpStatus: %s, type: %s, title: %s, detail: %s, instance: %s", resp.Status(), errResponse.Type, errResponse.Title, errResponse.Detail, errResponse.Instance))
	}

	return nil
}
