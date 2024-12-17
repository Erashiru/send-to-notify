package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/trade_group_details_request"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/trade_group_details_response"
	"net/http"
)

func (rkeeper *client) GetItemsByTradeGroup(ctx context.Context, tradeGroupId string) (trade_group_details_response.RK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  trade_group_details_response.RK7QueryResult
		request = trade_group_details_request.RK7Query{
			RK7CMD: trade_group_details_request.RK7CMD{
				CMD:     "GetRefData",
				RefName: "TRADEGROUPDETAILS",
				PROPFILTERS: trade_group_details_request.PropFilters{
					PROPFILTER: trade_group_details_request.PropFilter{
						Name:  "Parent",
						Value: tradeGroupId,
					},
				},
			},
		}
	)

	utils.Beautify("get items by trade group request body", request)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&request).
		SetResult(&result).
		Post(path)
	if err != nil {
		return trade_group_details_response.RK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return trade_group_details_response.RK7QueryResult{}, fmt.Errorf("get items by trade group response error: %v", response.Error())
	}

	utils.Beautify("pay order result body", result)

	if result.ErrorText != "" {
		return trade_group_details_response.RK7QueryResult{}, fmt.Errorf("get items by trade group result error, status: %s, error: %s", result.Status, result.ErrorText)
	}

	return result, nil
}
