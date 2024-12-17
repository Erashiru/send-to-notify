package http

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients"
	dto2 "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/dto"
	"testing"
	"time"
)

func TestClient_CreateOrder(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *clients.Config
		arg     dto2.CreateOrderRequest
		wantErr bool
	}{
		{
			name: "auth failed",
			cfg: &clients.Config{
				Protocol: "http",
				BaseURL:  "https://ws.ucs.ru/wsserverlp",
				ApiKey:   "QFX46lKB23c=9sz5zp1zixKRXN8/nf4bBflQ4NykfWfAr09r32G0DkWUa0p3nc29rtuvU356R28gQMxchi/gFygxg2cVwHBhwQR9lDlieUAm3V5zn37rPXjb2MQfwo/nMZbmThfnoxHqlAZK/SQXaF2ybcCpLlLDAjAQGkxBeD9vlfn+QU54ub2/nK15+cqYKBFrpefVrCCLu8ZCTqB2GsyOxrMroxqQT2nGXVhhl4fD",
			},
			arg: dto2.CreateOrderRequest{
				TaskType: "CreateOrder",
				Params: dto2.CreateOrderParam{
					Async: dto2.Sync{
						ObjectID: 553620001,
						Timeout:  10,
					},
					Order: dto2.Order{
						OriginalOrderId: uuid.New().String(),
						Customer: &dto2.PersonInfo{
							Name:  "Alisher",
							Phone: "77777777",
						},
						Payment: dto2.Payment{
							Type: "cash",
						},
						Delivery: dto2.Delivery{
							ExpectedTime: time.Now().Add(30 * time.Minute),
							Address: &dto2.DeliveryAddress{
								FullAddress: "Astana Hub",
							},
						},
						Products: []dto2.CreateOrderProduct{
							{
								Id:       "1000059",
								Name:     "Coca Cola",
								Quantity: 1,
							},
						},
						Comment: "не готовить",
						Price: &dto2.Price{
							Total: 500,
						},
						PersonsQuantity: 2,
					},
				},
			},
			wantErr: true,
		},
		//{
		//	name: "different object ID",
		//	cfg: &clients.Config{
		//		Protocol: "http",
		//		BaseURL:  "https://ws.ucs.ru/wsserverlp",
		//		ApiKey:   "QFX46lKB23c=9sz5zp1zixKRXN8/nf4bBflQ4NykfWfAr09r32G0DkWUa0p3nc29rtuvU356R28gQMxchi/gFygxg2cVwHBhwQR9lDlieUAm3V5zn37rPXjb2MQfwo/nMZbmThfnoxHqlAZK/SQXaF2ybcCpLlLDAjAQGkxBeD9vlfn+QU54ub2/nK15+cqYKBFrpefVrCCLu8ZCTqB2GsyOxrMroxqQT2nGXVhhl4fD",
		//	},
		//	arg:     dto.CreateOrderRequest{},
		//	wantErr: true,
		//},
		//{
		//	name: "OK",
		//	cfg: &clients.Config{
		//		Protocol: "http",
		//		BaseURL:  "https://ws.ucs.ru/wsserverlp",
		//		ApiKey:   "QFX46lKB23c=9sz5zp1zixKRXN8/nf4bBflQ4NykfWfAr09r32G0DkWUa0p3nc29rtuvU356R28gQMxchi/gFygxg2cVwHBhwQR9lDlieUAm3V5zn37rPXjb2MQfwo/nMZbmThfnoxHqlAZK/SQXaF2ybcCpLlLDAjAQGkxBeD9vlfn+QU54ub2/nK15+cqYKBFrpefVrCCLu8ZCTqB2GsyOxrMroxqQT2nGXVhhl4fD",
		//	},
		//	arg:     dto.CreateOrderRequest{},
		//	wantErr: false,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cli, err := NewClient(tt.cfg)
			if err != nil {
				t.Error(err)
				return
			}

			response, err := cli.CreateOrder(context.Background(), tt.arg.Params.Async.ObjectID, tt.arg.Params.Order)
			if err != nil {
				if tt.wantErr {
					t.Log(err)
					return
				}
				t.Error(err)
				return
			}

			t.Log(response)

			order, err := cli.CreateOrderTask(context.Background(), response.ResponseCommon.TaskGUID)
			if err != nil {
				t.Error(err)
				return
			}

			fmt.Println(order)

			getOrderTask, err := cli.GetOrder(context.TODO(), order.TaskResponse.Order.OrderGuid, tt.arg.Params.Async.ObjectID)
			if err != nil {
				t.Error(err)
				return
			}

			getOrder, err := cli.GetOrderTask(context.TODO(), getOrderTask.ResponseCommon.TaskGUID)
			if err != nil {
				t.Error(err)
			}

			fmt.Println(getOrder)

			cancelOrderTask, err := cli.CancelOrder(context.TODO(), tt.arg.Params.Async.ObjectID, order.TaskResponse.Order.OrderGuid)
			if err != nil {
				t.Error(err)
			}

			fmt.Println(cancelOrderTask)

			cancelOrder, err := cli.CancelOrderTask(context.TODO(), cancelOrderTask.ResponseCommon.TaskGUID)
			if err != nil {
				t.Error(err)
			}

			fmt.Println(cancelOrder)

		})
	}

}
