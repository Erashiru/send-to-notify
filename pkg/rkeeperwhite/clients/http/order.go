package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	dto2 "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/dto"
	"github.com/pkg/errors"
	"log"
)

func (cli Client) PayOrder(ctx context.Context, objectID, amount int, orderId, currency string) (dto2.SyncResponse, error) {
	path := "/api/v2/aggregators/Create"

	var (
		response dto2.SyncResponse
		body     = dto2.PayOrderRequest{
			TaskType: dto2.PayOrder.String(),
			Params: dto2.PayOrderParams{
				Async: dto2.PayOrderAsync{
					ObjectId: objectID,
					Timeout:  59,
				},
				OrderGuid: orderId,
				Payments: []dto2.PayOrderPayment{
					{
						Amount:   amount,
						Currency: currency,
					},
				},
			},
		}
	)

	utils.Beautify("rkeeper pay order request body", body)

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&body).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto2.SyncResponse{}, err
	}

	if resp.IsError() {
		return dto2.SyncResponse{}, fmt.Errorf("status: %s, body: %s, error: %+v", resp.Status(), string(resp.Body()), resp.Error())
	}

	utils.Beautify("rkeeper pay order response body", response)

	return response, nil
}

func (cli Client) CreateOrder(ctx context.Context, objectID int, order dto2.Order) (dto2.SyncResponse, error) {
	path := "/api/v2/aggregators/Create"

	var (
		response dto2.SyncResponse
		body     = dto2.CreateOrderRequest{
			TaskType: dto2.CreateOrder.String(),
			Params: dto2.CreateOrderParam{
				Async: dto2.Sync{
					ObjectID: objectID,
					Timeout:  180,
				},
				Order: order,
			},
		}
	)

	utils.Beautify("rkeeper create order request body", body)

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&body).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto2.SyncResponse{}, err
	}

	if resp.IsError() {
		return dto2.SyncResponse{}, errors.New(resp.Status() + " " + string(resp.Body()))
	}

	utils.Beautify("rkeeper create order response body", response)

	return response, nil
}

func (cl Client) CancelOrder(ctx context.Context, objectID int, orderGUID string) (dto2.SyncResponse, error) {
	path := "/api/v2/aggregators/Create"

	var (
		req = dto2.CancelOrderRequest{
			TaskType: dto2.CancelOrder.String(),
			Params: dto2.CancelOrderParam{
				Async: dto2.Sync{
					ObjectID: objectID,
					Timeout:  120,
				},
				OrderGuid: orderGUID,
			},
		}
		response dto2.SyncResponse
	)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(req).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto2.SyncResponse{}, err
	}

	if resp.IsError() {
		return dto2.SyncResponse{}, errors.New(resp.Status() + " " + string(resp.Body()))
	}

	return response, nil
}

func (cli Client) GetOrder(ctx context.Context, orderGUID string, objectID int) (dto2.SyncResponse, error) {
	path := "/api/v2/aggregators/Create"

	req := dto2.GetOrderRequest{
		TaskType: dto2.GetOrder.String(),
		Params: dto2.Params{
			OrderGUID: orderGUID,
			Async: &dto2.Sync{
				ObjectID: objectID,
				Timeout:  120,
			},
		},
	}

	b, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	var result dto2.SyncResponse

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetResult(&result).
		Post(path)

	if err != nil {
		return dto2.SyncResponse{}, err
	}

	if resp.IsError() {
		return dto2.SyncResponse{}, fmt.Errorf("rkeeper cli err: get order response %s status %s", string(resp.Body()), resp.Status())
	}

	return result, nil
}
