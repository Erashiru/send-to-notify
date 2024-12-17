package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	dto2 "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/dto"
)

func (cli Client) GetStopList(ctx context.Context, objectId int) (dto2.StopListResponse, error) {

	path := "/api/v2/aggregators/Create"

	req := dto2.TaskRequest{
		TaskType: "GetStopList",
		Params: dto2.Params{
			Sync: &dto2.Sync{
				ObjectID: objectId,
				Timeout:  tokenTimeout,
			},
		},
	}

	var res dto2.StopListResponse

	resp, err := cli.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(req).
		SetResult(&res).
		Post(path)

	if err != nil {
		return dto2.StopListResponse{}, err
	}

	if resp.IsError() {
		return dto2.StopListResponse{}, fmt.Errorf("rkeeper cli err: get stoplist %s status %s", resp.Body(), resp.Status())
	}

	if res.ErrResponse.WsError.Code != "" {
		return dto2.StopListResponse{}, fmt.Errorf("rkeeper cli err: %s, desc - %s", res.ErrResponse.WsError.Code, res.ErrResponse.WsError.Desc)
	}

	utils.Beautify("rkeeper stoplist response", res)

	return res, nil
}
