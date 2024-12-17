package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/service/iiko/models/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
)

func (c *Client) CreateDeliveryOrder(ctx context.Context, createDeliveryBody models.CreateDeliveryRequest) (models.CreateDeliveryResponse, error) {
	path := "/api/1/deliveries/create"

	var (
		response    models.CreateDeliveryResponse
		errResponse models.ErrorResponse2
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(createDeliveryBody).
		Post(path)
	if err != nil {
		utils.Beautify("iiko create delivery order error", err)
		utils.Beautify("iiko create delivery order errorResponse body", errResponse)
		log.Err(err).Msgf("iiko create delivery order request error: %v", err)
		return models.CreateDeliveryResponse{}, errors.New(errResponse.Err)
	}

	if resp.IsError() {
		utils.Beautify("iiko create delivery order errorResponse body", errResponse)
		log.Info().Msgf("iiko create delivery order error response body: %+v", resp)
		return models.CreateDeliveryResponse{}, errors.Wrap(errors.New(errResponse.Description), errResponse.Error())
	}

	return response, nil
}

func (c *Client) RetrieveDeliveryOrder(ctx context.Context, organizationID, orderID string) (models.RetrieveOrder, error) {
	path := "/api/1/deliveries/by_id"

	var (
		response    models.RetrieveDeliveryResponse
		errResponse models.ErrorResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&response).
		SetBody(models.RetrieveDeliveryRequest{
			OrganizationId: organizationID,
			OrderIds:       []string{orderID},
		}).
		Post(path)
	if err != nil {
		return models.RetrieveOrder{}, errResponse
	}

	if resp.IsError() {
		return models.RetrieveOrder{}, fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	if len(response.Orders) == 0 {
		return models.RetrieveOrder{}, fmt.Errorf("order is not exist in POS system")
	}

	return response.Orders[0], nil
}

func (c *Client) AddOrderItem(ctx context.Context, req models.OrderItem) (models.OrderItemResponse, error) {
	path := "/api/1/deliveries/add_items"

	var (
		response models.OrderItemResponse
	)

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&response).
		SetResult(&response).
		SetBody(req).
		Post(path)
	if err != nil {
		return response, err
	}

	utils.Beautify("request AddOrderItem", req)
	utils.Beautify("client response AddOrderItem", resp)

	if resp.IsError() {
		return response, fmt.Errorf("%s %+v", response.ErrorDescription, response)
	}
	utils.Beautify("response AddOrderItem", response)

	return response, nil
}

func (c *Client) CancelDeliveryOrder(ctx context.Context, organizationID, orderID, removalTypeId string) (models.CorID, error) {

	path := "/api/1/deliveries/cancel"

	var (
		errResponse models.ErrorResponse
		corID       models.CorID
	)

	body := models.CancelDeliveryResponse{
		OrganizationId: organizationID,
		OrderId:        orderID,
	}

	if removalTypeId != "" {
		body.RemovalTypeId = removalTypeId
	}

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&corID).
		SetBody(body).
		Post(path)
	if err != nil {
		utils.Beautify("iiko cancel delivery order error", err)
		utils.Beautify("iiko cancel delivery order errorResponse body", errResponse)
		log.Err(err).Msgf("iiko cancel delivery order request error: %v", err)
		return models.CorID{}, err
	}

	if resp.IsError() {
		utils.Beautify("iiko cancel delivery order response", resp)
		utils.Beautify("iiko cancel delivery order errorResponse body", errResponse)
		log.Info().Msgf("iiko cancel delivery order error response body: %v", resp)
		return models.CorID{}, fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return corID, nil
}

func (c *Client) UpdateOrderProblem(ctx context.Context, problem models.UpdateOrderProblem) error {

	path := "/api/1/deliveries/update_order_problem"

	var errResponse models.ErrorResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(problem).
		Post(path)
	if err != nil {
		log.Err(err).Msgf("error: iiko UpdateOrderProblem for pos_order_id: %s", problem.OrderId)
		return err
	}

	log.Info().Msgf("%+v request iiko UpdateOrderProblem, URL: %s\n", resp.Request.Body, resp.Request.URL)
	log.Info().Msgf("%+v response iiko UpdateOrderProblem\n", string(resp.Body()))

	if resp.IsError() {
		log.Info().Msgf("iiko UpdateOrderProblem response error: %s", resp.Error())
		return fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return nil
}

func (c *Client) CloseOrder(ctx context.Context, posOrderId, organizationId string) error {
	path := "/api/1/deliveries/close"

	var errResponse models.ErrorResponse

	body := models.CloseOrderRequest{OrderId: posOrderId, OrganizationId: organizationId}

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetBody(body).
		Post(path)
	if err != nil {
		log.Err(err).Msgf("iiko close order error, posOrderId: %s", posOrderId)
		return err
	}

	if resp.IsError() {
		log.Info().Msgf("iiko close order error, posOrderId: %s", posOrderId)
		return fmt.Errorf("%s %w", errResponse.Description, ErrResponse)
	}

	return nil
}
