package managers

import (
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"testing"
)

func Test_isNecessaryUpdateOrderStatus(t *testing.T) {
	type args struct {
		order models.Order
		store coreStoreModels.Store
	}
	var tests = []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				order: models.Order{
					Type:            "INSTANT",
					DeliveryService: "wolt",
					Status:          "NEW",
					StatusesHistory: []models.OrderStatusUpdate{
						models.OrderStatusUpdate{
							Name: "NEW",
						},
						models.OrderStatusUpdate{
							Name: "COOKING_COMPLETE",
						},
					},
				},
				store: coreStoreModels.Store{
					Wolt: coreStoreModels.StoreWoltConfig{
						PurchaseTypes: coreStoreModels.PurchaseTypes{
							Instant: []coreStoreModels.Status{

								coreStoreModels.Status{
									PosStatus: "COOKING_COMPLETE",
									Status:    "",
								},
								coreStoreModels.Status{
									PosStatus: "CLOSED",
									Status:    "",
								},
								coreStoreModels.Status{
									PosStatus: "ON_WAY",
									Status:    "Ready",
								},
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "test2",
			args: args{
				order: models.Order{
					Type:            "INSTANT",
					DeliveryService: "wolt",
					Status:          "CLOSED",
					StatusesHistory: []models.OrderStatusUpdate{
						models.OrderStatusUpdate{
							Name: "NEW",
						},
						models.OrderStatusUpdate{
							Name: "COOKING_COMPLETE",
						},
					},
				},
				store: coreStoreModels.Store{
					Wolt: coreStoreModels.StoreWoltConfig{
						PurchaseTypes: coreStoreModels.PurchaseTypes{
							Instant: []coreStoreModels.Status{
								coreStoreModels.Status{
									PosStatus: "COOKING_COMPLETE",
									Status:    "",
								},
								coreStoreModels.Status{
									PosStatus: "ON_WAY",
									Status:    "Ready",
								},
							},
							Preorder: []coreStoreModels.Status{},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "test3",
			args: args{
				order: models.Order{
					Type:            "INSTANT",
					DeliveryService: "wolt",
					Status:          "CLOSED",
					StatusesHistory: []models.OrderStatusUpdate{
						models.OrderStatusUpdate{
							Name: "NEW",
						},
						models.OrderStatusUpdate{
							Name: "COOKING_COMPLETE",
						},
						/*models.OrderStatusUpdate{
							Name: "ON_WAY",
						},*/
					},
				},
				store: coreStoreModels.Store{
					Wolt: coreStoreModels.StoreWoltConfig{
						PurchaseTypes: coreStoreModels.PurchaseTypes{
							Instant: []coreStoreModels.Status{
								{
									PosStatus: "ON_WAY",
									Status:    "Ready",
								},
							},
							Preorder: []coreStoreModels.Status{},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "test3",
			args: args{
				order: models.Order{
					Type:            "INSTANT",
					DeliveryService: "wolt",
					Status:          "WAIT_COOKING",
					StatusesHistory: []models.OrderStatusUpdate{
						{
							Name: "ACCEPTED",
						},
						{
							Name: "COOKING_COMPLETE",
						},
					},
				},
				store: coreStoreModels.Store{
					Wolt: coreStoreModels.StoreWoltConfig{
						PurchaseTypes: coreStoreModels.PurchaseTypes{
							Instant: []coreStoreModels.Status{
								coreStoreModels.Status{
									PosStatus: "READY_FOR_COOKING",
									Status:    "",
								},
								coreStoreModels.Status{
									PosStatus: "ON_WAY",
									Status:    "Ready",
								},
							},
							Preorder: []coreStoreModels.Status{},
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNecessaryUpdateOrderStatus(tt.args.order, tt.args.store)
			if got != tt.want {
				t.Errorf("isNecessaryUpdateOrderStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
