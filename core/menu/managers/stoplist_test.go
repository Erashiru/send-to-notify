package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"testing"
)

func Test_mnm_getMessage(t *testing.T) {

	tests := []struct {
		name      string
		store     storeModels.Store
		stoplists []string
		trx       models.StopListTransaction
	}{
		{
			name: "#1",
			store: storeModels.Store{
				Address: storeModels.StoreAddress{
					City: "Astana",
				},
				Name: "Rest1",
			},
			trx: models.StopListTransaction{
				StoreID: "1",
				Transactions: []models.TransactionData{
					{
						ID:       "1",
						StoreID:  "12",
						Delivery: "glovo",
						Status:   setStatus(models.SUCCESS.String()),
						Products: []models.StopListProduct{
							{
								ExtID:       "1",
								Name:        "товар1",
								IsAvailable: true,
							},
							{
								ExtID:       "2",
								Name:        "товар2",
								IsAvailable: false,
							},
							{
								ExtID:       "3",
								Name:        "товар3",
								IsAvailable: false,
							},
						},
						Attributes: []models.StopListAttribute{
							{
								AttributeID:   "11",
								AttributeName: "аттрибут1",
								IsAvailable:   true,
							},
							{
								AttributeID:   "21",
								AttributeName: "аттрибут2",
								IsAvailable:   false,
							},
							{
								AttributeID:   "31",
								AttributeName: "аттрибут3",
								IsAvailable:   false,
							},
						},
					},
					{
						ID:       "2",
						StoreID:  "15",
						Delivery: "glovo",
						Status:   setStatus(models.NOT_PROCESSED.String()),
						Products: []models.StopListProduct{
							{
								ExtID:       "1",
								Name:        "товар1",
								IsAvailable: true,
							},
							{
								ExtID:       "2",
								Name:        "товар2",
								IsAvailable: false,
							},
							{
								ExtID:       "3",
								Name:        "товар3",
								IsAvailable: false,
							},
						},
						Attributes: []models.StopListAttribute{
							{
								AttributeID:   "11",
								AttributeName: "аттрибут1",
								IsAvailable:   true,
							},
							{
								AttributeID:   "21",
								AttributeName: "аттрибут2",
								IsAvailable:   false,
							},
							{
								AttributeID:   "31",
								AttributeName: "аттрибут3",
								IsAvailable:   false,
							},
						},
					},
					{
						ID:       "2",
						StoreID:  "13",
						Delivery: "wolt",
						Status:   setStatus(models.SUCCESS.String()),
						Products: []models.StopListProduct{
							{
								ExtID:       "1",
								Name:        "товар1",
								IsAvailable: true,
							},
							{
								ExtID:       "2",
								Name:        "товар2",
								IsAvailable: false,
							},
							{
								ExtID:       "3",
								Name:        "товар3",
								IsAvailable: false,
							},
						},
						Attributes: []models.StopListAttribute{
							{
								AttributeID:   "11",
								AttributeName: "аттрибут1",
								IsAvailable:   true,
							},
							{
								AttributeID:   "21",
								AttributeName: "аттрибут2",
								IsAvailable:   false,
							},
							{
								AttributeID:   "31",
								AttributeName: "аттрибут3",
								IsAvailable:   false,
							},
						},
					},
					{
						ID:       "2",
						StoreID:  "14",
						Delivery: "wolt",
						Status:   setStatus(models.NOT_PROCESSED.String()),
						Products: []models.StopListProduct{
							{
								ExtID:       "1",
								Name:        "товар1",
								IsAvailable: true,
							},
							{
								ExtID:       "2",
								Name:        "товар2",
								IsAvailable: false,
							},
							{
								ExtID:       "3",
								Name:        "товар3",
								IsAvailable: false,
							},
						},
						Attributes: []models.StopListAttribute{
							{
								AttributeID:   "11",
								AttributeName: "аттрибут1",
								IsAvailable:   true,
							},
							{
								AttributeID:   "21",
								AttributeName: "аттрибут2",
								IsAvailable:   false,
							},
							{
								AttributeID:   "31",
								AttributeName: "аттрибут3",
								IsAvailable:   false,
							},
						},
					},
				},
			},
			stoplists: []string{
				"Товар1", "Товар2", "Аттрибут1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := getMessage(context.Background(), tt.store, tt.trx, tt.stoplists)
			t.Log(msg)
		})
	}
}
