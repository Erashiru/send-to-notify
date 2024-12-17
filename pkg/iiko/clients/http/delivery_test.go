package http

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/iiko/clients"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestClient_AddOrderItem(t *testing.T) {
	type fields struct {
		cli    *resty.Client
		apiKey string
		quit   chan struct{}
	}
	type args struct {
		ctx context.Context
		req models.OrderItem
	}
	cl, _ := New(&clients.Config{
		Protocol: "http",
		ApiLogin: "a99c4655-acb",
		BaseURL:  "https://api-ru.iiko.services",
	})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.OrderItemResponse
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				cli:    resty.New(),
				apiKey: "a99c4655-acb",
				quit:   make(chan struct{}),
			},
			args: args{
				ctx: context.Background(),
				req: models.OrderItem{
					OrganizationId: "08a9a229-e78a-4e2b-a351-9ab243801299",
					OrderId:        "8ab04ddb-d89b-46bc-8801-f0953c8da4f7",
					Items: []models.Item{
						models.Item{
							ProductId: "d7666134-4437-4a59-a503-afe1a0ac5700",
							Type:      "Product",
							Amount:    2,
						},
					},
				},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := cl.AddOrderItem(tt.args.ctx, tt.args.req)
			log.Info().Msgf("resp %v %v", got, err)
		})
	}
}
