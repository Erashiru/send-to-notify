package http

import (
	"context"
	"fmt"
	dto2 "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/dto"
)

func (cli Client) CreateOrderTask(ctx context.Context, taskGUID string) (dto2.CreateOrderTaskResponse, error) {
	path := "/api/v2/aggregators/Create"

	req := dto2.TaskRequest{
		TaskType: dto2.Task.String(),
		Params: dto2.Params{
			TaskGUID: taskGUID,
		},
	}

	var result dto2.CreateOrderTaskResponse

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetResult(&result).
		Post(path)

	if err != nil {
		return dto2.CreateOrderTaskResponse{}, err
	}

	if resp.IsError() {
		return dto2.CreateOrderTaskResponse{}, fmt.Errorf("rkeeper cli err: create order task response %s status %s", string(resp.Body()), resp.Status())
	}

	return result, nil
}

func (cli Client) GetOrderTask(ctx context.Context, taskGUID string) (dto2.GetOrderTaskResponse, error) {
	path := "/api/v2/aggregators/Create"

	req := dto2.TaskRequest{
		TaskType: dto2.Task.String(),
		Params: dto2.Params{
			TaskGUID: taskGUID,
		},
	}

	var result dto2.GetOrderTaskResponse

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetResult(&result).
		Post(path)

	if err != nil {
		return dto2.GetOrderTaskResponse{}, err
	}

	if resp.IsError() {
		return dto2.GetOrderTaskResponse{}, fmt.Errorf("rkeeper cli err: get order task response %s status %s", string(resp.Body()), resp.Status())
	}

	return result, nil
}

func (cli Client) CancelOrderTask(ctx context.Context, taskGUID string) (dto2.CancelOrderResponse, error) {
	path := "/api/v2/aggregators/Create"

	req := dto2.TaskRequest{
		TaskType: dto2.CancelOrder.String(),
		Params: dto2.Params{
			TaskGUID: taskGUID,
		},
	}

	var result dto2.CancelOrderResponse

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetResult(&result).
		Post(path)

	if err != nil {
		return dto2.CancelOrderResponse{}, err
	}

	if resp.IsError() {
		return dto2.CancelOrderResponse{}, fmt.Errorf("rkeeper cli err: get order task response %s status %s", string(resp.Body()), resp.Status())
	}

	return result, nil
}
