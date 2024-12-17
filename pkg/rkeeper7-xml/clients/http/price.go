package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/product_price_request"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/product_price_response"
	"net/http"
)

func (rkeeper *client) GetDeliveryPriceByPriceTypeId(ctx context.Context, priceTypeId string) (product_price_response.RK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  product_price_response.RK7QueryResult
		request = product_price_request.RK7Query{
			RK7CMD: product_price_request.RK7CMD{
				CMD:     "GetRefData",
				RefName: "PRICES",
				PROPFILTERS: product_price_request.PropFilters{
					PROPFILTER: product_price_request.PropFilter{
						Name:  "PriceType",
						Value: priceTypeId,
					},
				},
			},
		}
	)

	utils.Beautify("get product prices by price type request body", request)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&request).
		SetResult(&result).
		Post(path)
	if err != nil {
		return product_price_response.RK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return product_price_response.RK7QueryResult{}, fmt.Errorf("get product prices by price type response error: %v", response.Error())
	}

	utils.Beautify("get product prices by price result body", result)

	if result.ErrorText != "" {
		return product_price_response.RK7QueryResult{}, fmt.Errorf("get product prices by price type result error, status: %s, error: %s", result.Status, result.ErrorText)
	}

	return result, nil
}
