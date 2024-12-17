package http

import (
	"context"
	"fmt"
	dto2 "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/dto"
	"github.com/rs/zerolog/log"
)

func (cli Client) GetMenu(ctx context.Context, objectId int) (dto2.MenuResponse, error) {

	path := "/api/v2/aggregators/Create"

	req := dto2.TaskRequest{
		TaskType: "GetMenu",
		Params: dto2.Params{
			Sync: &dto2.Sync{
				ObjectID: objectId,
				Timeout:  tokenTimeout,
			},
		},
	}

	var res dto2.MenuResponse

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(req).
		SetResult(&res).
		Post(path)

	if err != nil {
		return dto2.MenuResponse{}, err
	}

	if resp.IsError() {
		return dto2.MenuResponse{}, fmt.Errorf("rkeeper cli err: get menu %s status %s", resp.Body(), resp.Status())
	}

	if res.ErrResponse.WsError.Code != "" {
		return dto2.MenuResponse{}, fmt.Errorf("rkeeper cli err: code - %s, desc - %s", res.ErrResponse.WsError.Code, res.ErrResponse.WsError.Desc)
	}

	return res, nil

}

func (cli Client) UpdateMenu(ctx context.Context, objectID int) (dto2.SyncResponse, error) {
	path := "/api/v2/aggregators/Create"

	req := dto2.TaskRequest{
		TaskType: "UpdateMenu",
		Params: dto2.Params{
			Sync: &dto2.Sync{
				ObjectID: objectID,
				Timeout:  tokenTimeout,
			},
		},
	}

	var res dto2.SyncResponse

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(req).
		SetResult(&res).
		Post(path)

	if err != nil {
		return dto2.SyncResponse{}, err
	}

	if resp.IsError() {
		return dto2.SyncResponse{}, fmt.Errorf("rkeeper cli err: update menu %s status %s", resp.Body(), resp.Status())
	}

	if res.Error.WsError.Code != "" {
		return dto2.SyncResponse{}, fmt.Errorf("rkeeper cli err: code - %s, desc - %s", res.Error.WsError.Code, res.Error.WsError.Desc)
	}

	return res, nil
}

func (cli Client) GetMenuByParams(ctx context.Context, objectID, priceTypeID int) (dto2.MenuResponse, error) {
	path := "/api/v2/aggregators/Create"

	req := dto2.TaskRequest{
		TaskType: "GetMenuByParams",
		Params: dto2.Params{
			Sync: &dto2.Sync{
				ObjectID: objectID,
				Timeout:  tokenTimeout,
			},
			PriceTypeID:         priceTypeID,
			FilterByKassPresets: false,
		},
	}

	var res dto2.MenuResponse

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(req).
		SetResult(&res).
		Post(path)
	if err != nil {
		return dto2.MenuResponse{}, err
	}

	log.Info().Msgf("rkeeper cli get menu by params request url: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		log.Err(err).Msgf("rkeeper cli err: get menu by params error: %+v, status code: %d, body: %s", resp.Error(), resp.StatusCode(), string(resp.Body()))
		return dto2.MenuResponse{}, fmt.Errorf("rkeeper cli err: get menu by params error: %+v, status code: %d, body: %s", resp.Error(), resp.StatusCode(), string(resp.Body()))
	}

	if res.ErrResponse.WsError.Code != "" {
		log.Err(err).Msgf("rkeeper cli err: code - %s, desc - %s", res.ErrResponse.WsError.Code, res.ErrResponse.WsError.Desc)
		return dto2.MenuResponse{}, fmt.Errorf("rkeeper cli err: code - %s, desc - %s", res.ErrResponse.WsError.Code, res.ErrResponse.WsError.Desc)
	}

	log.Info().Msgf("rkeeper cli get menu by params status code: %d, products: %d, categroies: %d", resp.StatusCode(), len(res.TaskResponse.Menu.Products), len(res.TaskResponse.Menu.Categories))

	return res, nil
}
