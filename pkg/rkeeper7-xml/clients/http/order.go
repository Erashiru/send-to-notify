package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/create_order_request"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/create_order_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_order_request"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_order_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/pay_order_request"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/pay_order_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/save_order_request"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/save_order_response"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/set_delivery_type"
	"net/http"
)

func (rkeeper *client) SetDeliveryTypeToOrder(ctx context.Context, visitId, orderType string) error {
	path := "/rk7api/v0/xmlinterface.xml"

	req := set_delivery_type.RK7Query{
		RK7CMD: set_delivery_type.RK7CMD{
			CMD: "DeliveryUpdateStatus",
			Order: set_delivery_type.Order{
				OrderIdent: "256",
				Visit:      visitId,
			},
			OrderType: set_delivery_type.OrderType{
				ID: orderType,
			},
			ExtSource: set_delivery_type.ExtSource{
				Source: "31",
			},
		},
	}

	utils.Beautify("request body", req)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		Post(path)
	if err != nil {
		return err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return fmt.Errorf("set delivery type to order response error: %v", response.Error())
	}

	return nil
}

func (rkeeper *client) CreateOrder(ctx context.Context, table, stationID, comment, orderTypeCode string) (create_order_response.CreateOrderResponse, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result create_order_response.CreateOrderResponse
		req    = create_order_request.RK7Query{
			RK7CMD: create_order_request.RK7CMD{
				CMD: "CreateOrder",
				Order: create_order_request.CreateOrderRequest{
					PersistentComment: comment,
					Table: create_order_request.Table{
						Code: table,
					},
					OrderType: create_order_request.OrderType{
						Code: orderTypeCode,
					},
					Station: create_order_request.Station{
						ID: stationID,
					},
					GuestType: create_order_request.GuestType{
						ID: "1",
					},
					Guests: create_order_request.Guests{
						Guest: create_order_request.Guest{
							GuestLabel: "1",
						},
					},
				},
			},
		}
	)

	utils.Beautify("request body", req)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetResult(&result).
		Post(path)
	if err != nil {
		return create_order_response.CreateOrderResponse{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return create_order_response.CreateOrderResponse{}, fmt.Errorf("create order response error: %v", response.Error())
	}

	if result.ErrorText != "" {
		return create_order_response.CreateOrderResponse{}, fmt.Errorf("create order result error, status: %s, error: %s", result.Status, result.ErrorText)
	}

	return result, nil
}

func (rkeeper *client) SaveOrder(ctx context.Context, visitID, seqNumber, paymentId, prepayReasonId string, dishes save_order_request.Dishes, amount string, isLifeLicence bool) (save_order_response.RK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  save_order_response.RK7QueryResult
		request = save_order_request.RK7Query{
			RK7CMD: save_order_request.RK7CMD{
				CMD: "SaveOrder",
				Order: save_order_request.Order{
					Visit:      visitID,
					OrderIdent: "256",
				},
				Session: save_order_request.Session{
					Station: save_order_request.Station{
						ID: rkeeper.stationID,
					},
					Dish: dishes.ConvertToCorrectFormat(),
				},
			},
		}
	)

	if !isLifeLicence {
		request.RK7CMD.LicenseInfo = &save_order_request.LicenseInfo{
			Anchor:       "6:" + rkeeper.anchor + "#" + rkeeper.objectID + "/17",
			LicenseToken: rkeeper.licenseToken,
			LicenseInstance: &save_order_request.LicenseInstance{
				Guid:      rkeeper.licenseInstanceGUID,
				SeqNumber: seqNumber,
			},
		}
		request.RK7CMD.Session.Prepay = &save_order_request.Prepay{
			ID:       paymentId,
			Amount:   amount + "00",
			Promised: "1",
			Reason: &save_order_request.Reason{
				ID: prepayReasonId,
			},
		}
	}

	utils.Beautify("save order request body", request)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&request).
		SetResult(&result).
		Post(path)
	if err != nil {
		return save_order_response.RK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return save_order_response.RK7QueryResult{}, fmt.Errorf("save order response error: %v", response.Error())
	}

	utils.Beautify("save order result body", result)

	if result.ErrorText != "" {
		return save_order_response.RK7QueryResult{}, fmt.Errorf("save order result error, status: %s, error: %s", result.Status, result.ErrorText)
	}

	return result, nil
}

func (rkeeper *client) PayOrder(ctx context.Context, orderGUID, paymentID, sum, stationCode string) (pay_order_response.RK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  pay_order_response.RK7QueryResult
		request = pay_order_request.RK7Query{
			RK7CMD: pay_order_request.RK7CMD{
				CMD: "PayOrder",
				Order: pay_order_request.Order{
					Guid: orderGUID,
				},
				Cashier: pay_order_request.Cashier{
					Code: rkeeper.cashier,
				},
				Station: pay_order_request.Station{
					Code: stationCode,
				},
				Payment: pay_order_request.Payment{
					ID:     paymentID,
					Amount: sum + "00",
				},
			},
		}
	)

	utils.Beautify("pay order request body", request)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&request).
		SetResult(&result).
		Post(path)
	if err != nil {
		return pay_order_response.RK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return pay_order_response.RK7QueryResult{}, fmt.Errorf("pay order response error: %v", response.Error())
	}

	utils.Beautify("pay order result body", result)

	if result.ErrorText != "" {
		return pay_order_response.RK7QueryResult{}, fmt.Errorf("pay order result error, status: %s, error: %s", result.Status, result.ErrorText)
	}

	return result, nil
}

func (rkeeper *client) GetOrder(ctx context.Context, visitID string) (get_order_response.RK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  get_order_response.RK7QueryResult
		request = get_order_request.RK7Query{
			RK7Command: get_order_request.RK7Command{
				CMD: "GetOrder",
				Order: get_order_request.Order{
					Visit:      visitID,
					OrderIdent: "256",
				},
			},
		}
	)

	utils.Beautify("get order request body", request)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&request).
		SetResult(&result).
		Post(path)
	if err != nil {
		return get_order_response.RK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return get_order_response.RK7QueryResult{}, fmt.Errorf("get order response error: %v", response.Error())
	}

	utils.Beautify("get order result body", result)

	if result.CommandResult.ErrorText != "" {
		return get_order_response.RK7QueryResult{}, fmt.Errorf("get order result error, status: %s, error: %s", result.CommandResult.Status, result.CommandResult.ErrorText)
	}

	return result, nil
}
