package aggregator

import (
	"encoding/json"
	models2 "github.com/kwaaka-team/orders-core/core/glovo/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGlovoGetSystemCreateOrderRequestByAggregatorRequest(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(t.Name())
		}
	}()

	store := storeModels.Store{}
	store.Settings.TimeZone.TZ = "Asia/Almaty"

	fileData, err := os.Open("glovo_test_data.json")
	if err != nil {
		t.Fatal(err)
	}

	actual := models2.Order{}
	jsonParser := json.NewDecoder(fileData)
	if err = jsonParser.Decode(&actual); err != nil {
		t.Fatal(err)
	}
	deliveryFee := 10000
	expected := models2.Order{
		OrderID:       "1000000001",
		StoreID:       "2000000002",
		OrderTime:     "2023-12-09 15:42:46",
		PaymentMethod: "CASH",
		Currency:      "KGS",
		Customer: models2.Customer{
			Name:        "TestCustomer",
			PhoneNumber: "+71234567890",
			Hash:        "hash1",
		},
		OrderCode:                 "order_code_1",
		EstimatedTotalPrice:       31900,
		DeliveryFee:               &deliveryFee,
		MinimumBasketSurcharge:    3000,
		UtcOffsetMinutes:          "360",
		CustomerCashPaymentAmount: 100000,
		Products: []models2.ProductOrder{
			{
				Quantity:           1,
				Price:              31900,
				Name:               "Роллы \"Кани темпура\"",
				ID:                 "00000000-0000-0000-0000-000000000000",
				PurchasedProductID: "1233456789",
				Attributes:         make([]models2.AttributeOrder, 0),
			},
		},
		DeliveryAddress: models2.DeliveryAddress{
			Label:     "address 1",
			Latitude:  11.0000001234,
			Longitude: 22.000000012345,
		},
		BundledOrders:                  make([]string, 0),
		PickUpCode:                     "123",
		IsPickedUpByCustomer:           false,
		CutleryRequested:               false,
		PartnerDiscountedProductsTotal: 31900,
		TotalCustomerToPay:             44900,
		DiscountedProductsTotal:        31900,
	}

	assert.Equal(t, expected, actual)
}
