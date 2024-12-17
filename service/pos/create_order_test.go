package pos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/config"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/mocks"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"testing"
	"time"
)

type TestCases struct {
	Case []Case
}

type Case struct {
	CaseName        string `json:"case_name"`
	CaseDescription string `json:"description"`
	Order           models.Order
	AggregatorMenu  coreMenuModels.Menu
	PosMenu         coreMenuModels.Menu
	Store           coreStoreModels.Store
	iikoExpectValue iikoModels.CreateDeliveryRequest
	iikoReturnValue iikoModels.CreateDeliveryResponse
}

func TestCreateOrderV2(t *testing.T) {
	// set environment
	os.Setenv("SECRET_ENV", "StageEnvs")
	os.Setenv("REGION", "eu-west-1")
	os.Setenv("SENTRY", "StageSentry")

	// create context
	ctx := context.Background()

	// mocks
	ctrl := gomock.NewController(t)

	iikoCli := mocks.NewMockIIKO(ctrl)
	menuCli := mocks.NewMockClient(ctrl)

	// menu entity
	productExtID := uuid.New().String()
	productPosID := uuid.New().String()
	productProductID := uuid.New().String()

	aggregatorMenuID := primitive.NewObjectID().Hex()
	posMenuID := primitive.NewObjectID().Hex()

	// store entity
	storeID := primitive.NewObjectID().Hex()

	// glovo config entity
	externalStoreID := "kwaaka123"
	cashPaymentTypeID := uuid.New().String()
	cardPaymentTypeID := uuid.New().String()
	cashPaymentTypeKind := "Cash"
	cardPaymentTypeKind := "Card"

	// store settings entity
	timeZoneAlmaty := "Asia/Almaty"
	timeZoneUTCOffset := 6

	iikoCloudOrganizationID := uuid.New().String()
	iikoCloudTerminalID := uuid.New().String()
	iikoCloudApiKey := uuid.New().String()

	// pos settings for iiko
	transportToFrontTimeout := 180

	// order entity
	orderUniqueId := primitive.NewObjectID().Hex()
	orderCode := uuid.New().String()

	orderPaymentMethod := "CASH"
	pickUpCode := "777"
	orderTime := time.Now().UTC()
	estimatedPickUpTime := time.Now().UTC().Add(30 * time.Minute)

	totalOrder := 500.0

	// order glovo entity
	glovoOrderId := "12345678" // must be number in glovo

	completeBeforeDate := estimatedPickUpTime.Add(time.Duration(timeZoneUTCOffset) * time.Hour)

	completeBeforeComment := completeBeforeDate.Format("15:04:05")
	//completeBefore := completeBeforeDate.Format("2006-01-02 15:04:05.000")

	// pos order entiry
	posOrderID := uuid.New().String()

	testCases := TestCases{
		Case: []Case{
			{
				CaseName:        "glovo_order_case_1",
				CaseDescription: "Creating order with marketplace=true(glovo couriers), payment method cash, without order types",
				iikoExpectValue: iikoModels.CreateDeliveryRequest{
					OrganizationID:  iikoCloudOrganizationID,
					TerminalGroupID: iikoCloudTerminalID,
					CreateOrderSettings: &iikoModels.CreateOrderSettings{
						TransportToFrontTimeout: transportToFrontTimeout,
					},
					Order: &iikoModels.Order{
						Phone:            "+77771111111",
						OrderServiceType: models.ORDER_DELIVERY_CLIENT,
						Payments: []iikoModels.Payment{
							{
								PaymentTypeKind: cashPaymentTypeKind,
								Sum:             500,
								PaymentTypeID:   cashPaymentTypeID,
							},
						},
						Items: []iikoModels.Item{
							{
								ProductId: productProductID,
								Price:     &totalOrder,
								Type:      "Product",
								Amount:    1,
							},
						},
						Comment: fmt.Sprintf("Код заказа: %s\nАллергия: Нет\nАдрес: Нет\nПриготовить к: %s\nКомментарий: Нет\nДоставка: Glovo\nТип оплаты: Наличный Glovo", pickUpCode, completeBeforeComment),
						Customer: &iikoModels.Customer{
							Comment: "Номер курьера: Нет\nАллергия: Нет\nКомментарий: Нет\nСтоловые приборы: Не нужны",
							Gender:  "NotSpecified",
						},
					},
				},
				iikoReturnValue: iikoModels.CreateDeliveryResponse{
					CorrelationID: "123",
					OrderInfo: &iikoModels.OrderInfo{
						ID: posOrderID,
					},
				},
				Order: models.Order{
					ID:              orderUniqueId,
					Type:            models.ORDER_TYPE_INSTANT, // type of order
					DeliveryService: models.GLOVO.String(),     // delivery service
					PosType:         models.IIKO.String(),      // pos type
					RestaurantID:    storeID,                   // restaurant id
					OrderID:         glovoOrderId,              // aggregator order id
					OrderCode:       orderCode,                 // order code
					PickUpCode:      pickUpCode,                // pick up code
					IsMarketplace:   true,                      // delivery type
					PaymentMethod:   orderPaymentMethod,        // payment method (CASH, DELAYED)
					OrderTime: models.TransactionTime{
						Value: models.Time{
							Time: orderTime, // order time
						},
					},
					EstimatedPickupTime: models.TransactionTime{
						Value: models.Time{
							Time: estimatedPickUpTime, // estimated pickup time
						},
					},
					EstimatedTotalPrice: models.Price{
						Value:        500, // order total price
						CurrencyCode: "KZT",
					},
					HasServiceFee: false, // has service fee
					ServiceFeeSum: 0,     // service fee sum
					PosPaymentInfo: models.PosPaymentInfo{
						PaymentTypeID:   cashPaymentTypeID,   // cash payment type id
						PaymentTypeKind: cashPaymentTypeKind, // cash payment type kind
					},
					Products: []models.OrderProduct{
						{
							ID:       productExtID,
							Quantity: 1,
							Name:     "Burger Test",
							Price: models.Price{
								Value:        500,
								CurrencyCode: "KZT",
							},
						},
					},
				},
				Store: coreStoreModels.Store{
					ID:      storeID,
					Name:    "store IIKO, deliveryService glovo, marketplace true, without any order types",
					MenuID:  posMenuID,
					PosType: models.IIKO.String(),
					Glovo: coreStoreModels.StoreGlovoConfig{
						StoreID:       []string{externalStoreID},
						SendToPos:     true,
						IsMarketplace: true, // delivery type (true - aggregator delivery or pickup, false - restaurant delivery)
						PaymentTypes: coreStoreModels.DeliveryServicePaymentType{
							CASH: coreStoreModels.PaymentType{
								PaymentTypeID:   cashPaymentTypeID,   // payment type id
								PaymentTypeKind: cashPaymentTypeKind, // kind of payment
							},
							DELAYED: coreStoreModels.PaymentType{
								PaymentTypeID:   cardPaymentTypeID,   // payment type id
								PaymentTypeKind: cardPaymentTypeKind, // kind of payment
							},
						},
					},
					IikoCloud: coreStoreModels.StoreIikoConfig{
						OrganizationID: iikoCloudOrganizationID,
						TerminalID:     iikoCloudTerminalID,
						Key:            iikoCloudApiKey,
					},
					Menus: coreStoreModels.StoreDSMenus{
						{
							ID:        aggregatorMenuID,
							Name:      "test aggregator menu",
							Delivery:  models.GLOVO.String(),
							IsActive:  true,
							IsDeleted: false,
						},
					},
					Settings: coreStoreModels.Settings{
						TimeZone: coreStoreModels.TimeZone{
							UTCOffset: 6,
							TZ:        timeZoneAlmaty,
						},
					},
				},
				AggregatorMenu: coreMenuModels.Menu{
					ID:   aggregatorMenuID,
					Name: "aggregator menu",
					Products: coreMenuModels.Products{
						{
							ExtID: productExtID,
							PosID: productPosID,
							Name: []coreMenuModels.LanguageDescription{
								{
									Value: "Burger Test",
								},
							},
							Price: []coreMenuModels.Price{
								{
									Value:        500,
									CurrencyCode: "KZT",
								},
							},
						},
					},
				},
				PosMenu: coreMenuModels.Menu{
					ID:   posMenuID,
					Name: "pos menu",
					Products: coreMenuModels.Products{
						{
							ExtID:     productPosID,
							ProductID: productProductID,
							Name: []coreMenuModels.LanguageDescription{
								{
									Value: "Burger Test",
								},
							},
							Price: []coreMenuModels.Price{
								{
									Value:        500,
									CurrencyCode: "KZT",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases.Case {
		menuCli.EXPECT().GetStorePromos(ctx, gomock.Any()).Return([]dto.PromoDiscount{}, nil)

		iikoCli.EXPECT().CreateDeliveryOrder(ctx, gomock.Any()).Do(func(ctx context.Context, arg iikoModels.CreateDeliveryRequest) {
			if testCase.iikoExpectValue.TerminalGroupID != arg.TerminalGroupID {
				t.Errorf("terminal group id, want=%s, got=%s", testCase.iikoExpectValue.TerminalGroupID, arg.TerminalGroupID)
			}

			if testCase.iikoExpectValue.OrganizationID != arg.OrganizationID {
				t.Errorf("organization id, want=%s, got=%s", testCase.iikoExpectValue.OrganizationID, arg.OrganizationID)
			}

			if testCase.iikoExpectValue.CreateOrderSettings.TransportToFrontTimeout != arg.CreateOrderSettings.TransportToFrontTimeout {
				t.Errorf("transport to front timeout, want=%d, got=%d", testCase.iikoExpectValue.CreateOrderSettings.TransportToFrontTimeout, arg.CreateOrderSettings.TransportToFrontTimeout)
			}

			if testCase.iikoExpectValue.Order.ID != arg.Order.ID {
				t.Errorf("order id, want=%s, got=%s", testCase.iikoExpectValue.Order.ID, arg.Order.ID)
			}

			if testCase.iikoExpectValue.Order.Phone != arg.Order.Phone {
				t.Errorf("order phone, want=%s, got=%s", testCase.iikoExpectValue.Order.Phone, arg.Order.Phone)
			}

			if testCase.iikoExpectValue.Order.OrderServiceType != arg.Order.OrderServiceType {
				t.Errorf("order service type, want=%s, got=%s", testCase.iikoExpectValue.Order.OrderServiceType, arg.Order.OrderServiceType)
			}

			if testCase.iikoExpectValue.Order.Payments[0].PaymentTypeKind != arg.Order.Payments[0].PaymentTypeKind {
				t.Errorf("order payment type kind, want=%s, got=%s", testCase.iikoExpectValue.Order.Payments[0].PaymentTypeKind, arg.Order.Payments[0].PaymentTypeKind)
			}

			if testCase.iikoExpectValue.Order.Payments[0].Sum != arg.Order.Payments[0].Sum {
				t.Errorf("order payment sum, want=%d, got=%d", testCase.iikoExpectValue.Order.Payments[0].Sum, arg.Order.Payments[0].Sum)
			}

			if testCase.iikoExpectValue.Order.Payments[0].PaymentTypeID != arg.Order.Payments[0].PaymentTypeID {
				t.Errorf("order payment type id, want=%s, got=%s", testCase.iikoExpectValue.Order.Payments[0].PaymentTypeID, arg.Order.Payments[0].PaymentTypeID)
			}

			if testCase.iikoExpectValue.Order.Items[0].ProductId != arg.Order.Items[0].ProductId {
				t.Errorf("order item product id, want=%s, got=%s", testCase.iikoExpectValue.Order.Items[0].ProductId, arg.Order.Items[0].ProductId)
			}

			if *testCase.iikoExpectValue.Order.Items[0].Price != *arg.Order.Items[0].Price {
				t.Errorf("order item price, want=%v, got=%v", *testCase.iikoExpectValue.Order.Items[0].Price, *arg.Order.Items[0].Price)
			}

			if testCase.iikoExpectValue.Order.Items[0].Type != arg.Order.Items[0].Type {
				t.Errorf("order item type, want=%s, got=%s", testCase.iikoExpectValue.Order.Items[0].Type, arg.Order.Items[0].Type)
			}

			if testCase.iikoExpectValue.Order.Comment != arg.Order.Comment {
				t.Errorf("order comment, want=%s, got=%s", testCase.iikoExpectValue.Order.Comment, arg.Order.Comment)
			}

			if os.Getenv("ENVIRONMENT") != "github" {
				expectedFile, err := json.MarshalIndent(testCase.iikoExpectValue, "", " ")
				if err != nil {
					t.Errorf("marshal expected iiko order v2 error: %v", err)
				}

				if err := os.WriteFile(fmt.Sprintf("./test_jsons/expected_%s_v2.json", testCase.CaseName), expectedFile, 0644); err != nil {
					t.Errorf("write expected iiko order v2 error: %v", err)
				}

				returnedFile, err := json.MarshalIndent(testCase.iikoExpectValue, "", " ")
				if err != nil {
					t.Errorf("marshal returned iiko order v2 error: %v", err)
				}

				if err = os.WriteFile(fmt.Sprintf("./test_jsons/returned_%s_v2.json", testCase.CaseName), returnedFile, 0644); err != nil {
					t.Errorf("write returned iiko order v2 error: %v", err)
				}
			} else {
				log.Info().Msgf("ENVIRONMENT=github")
			}
		}).Return(testCase.iikoReturnValue, nil)

		bps := &BasePosService{}

		iikoSvc := iikoService{
			BasePosService:          bps,
			iikoClient:              iikoCli,
			transportToFrontTimeout: transportToFrontTimeout,
		}

		// create order v2
		res, err := iikoSvc.CreateOrder(ctx, testCase.Order, config.Configuration{}, testCase.Store, testCase.PosMenu, menuCli, testCase.AggregatorMenu, nil, nil, nil)
		if err != nil {
			t.Errorf("create order v2 error: %v", err)
			return
		}

		_ = res

		menuCli.EXPECT().GetStorePromos(ctx, gomock.Any()).Return([]dto.PromoDiscount{}, nil)

		iikoCli.EXPECT().CreateDeliveryOrder(ctx, gomock.Any()).Do(func(ctx context.Context, arg iikoModels.CreateDeliveryRequest) {
			if testCase.iikoExpectValue.TerminalGroupID != arg.TerminalGroupID {
				t.Errorf("terminal group id, want=%s, got=%s", testCase.iikoExpectValue.TerminalGroupID, arg.TerminalGroupID)
			}

			if testCase.iikoExpectValue.OrganizationID != arg.OrganizationID {
				t.Errorf("organization id, want=%s, got=%s", testCase.iikoExpectValue.OrganizationID, arg.OrganizationID)
			}

			if testCase.iikoExpectValue.CreateOrderSettings.TransportToFrontTimeout != arg.CreateOrderSettings.TransportToFrontTimeout {
				t.Errorf("transport to front timeout, want=%d, got=%d", testCase.iikoExpectValue.CreateOrderSettings.TransportToFrontTimeout, arg.CreateOrderSettings.TransportToFrontTimeout)
			}

			if testCase.iikoExpectValue.Order.ID != arg.Order.ID {
				t.Errorf("order id, want=%s, got=%s", testCase.iikoExpectValue.Order.ID, arg.Order.ID)
			}

			if testCase.iikoExpectValue.Order.Phone != arg.Order.Phone {
				t.Errorf("order phone, want=%s, got=%s", testCase.iikoExpectValue.Order.Phone, arg.Order.Phone)
			}

			if testCase.iikoExpectValue.Order.OrderServiceType != arg.Order.OrderServiceType {
				t.Errorf("order service type, want=%s, got=%s", testCase.iikoExpectValue.Order.OrderServiceType, arg.Order.OrderServiceType)
			}

			if testCase.iikoExpectValue.Order.Payments[0].PaymentTypeKind != arg.Order.Payments[0].PaymentTypeKind {
				t.Errorf("order payment type kind, want=%s, got=%s", testCase.iikoExpectValue.Order.Payments[0].PaymentTypeKind, arg.Order.Payments[0].PaymentTypeKind)
			}

			if testCase.iikoExpectValue.Order.Payments[0].Sum != arg.Order.Payments[0].Sum {
				t.Errorf("order payment sum, want=%d, got=%d", testCase.iikoExpectValue.Order.Payments[0].Sum, arg.Order.Payments[0].Sum)
			}

			if testCase.iikoExpectValue.Order.Payments[0].PaymentTypeID != arg.Order.Payments[0].PaymentTypeID {
				t.Errorf("order payment type id, want=%s, got=%s", testCase.iikoExpectValue.Order.Payments[0].PaymentTypeID, arg.Order.Payments[0].PaymentTypeID)
			}

			if testCase.iikoExpectValue.Order.Items[0].ProductId != arg.Order.Items[0].ProductId {
				t.Errorf("order item product id, want=%s, got=%s", testCase.iikoExpectValue.Order.Items[0].ProductId, arg.Order.Items[0].ProductId)
			}

			if *testCase.iikoExpectValue.Order.Items[0].Price != *arg.Order.Items[0].Price {
				t.Errorf("order item price, want=%v, got=%v", *testCase.iikoExpectValue.Order.Items[0].Price, *arg.Order.Items[0].Price)
			}

			if testCase.iikoExpectValue.Order.Items[0].Type != arg.Order.Items[0].Type {
				t.Errorf("order item type, want=%s, got=%s", testCase.iikoExpectValue.Order.Items[0].Type, arg.Order.Items[0].Type)
			}

			if testCase.iikoExpectValue.Order.Comment != arg.Order.Comment {
				t.Errorf("order comment, want=%s, got=%s", testCase.iikoExpectValue.Order.Comment, arg.Order.Comment)
			}

			if os.Getenv("ENVIRONMENT") != "github" {
				expectedFile, err := json.MarshalIndent(testCase.iikoExpectValue, "", " ")
				if err != nil {
					t.Errorf("marshal expected iiko order v2 error: %v", err)
				}

				if err := os.WriteFile(fmt.Sprintf("./test_jsons/expected_%s.json", testCase.CaseName), expectedFile, 0644); err != nil {
					t.Errorf("write expected iiko order v2 error: %v", err)
				}

				returnedFile, err := json.MarshalIndent(testCase.iikoExpectValue, "", " ")
				if err != nil {
					t.Errorf("marshal returned iiko order v2 error: %v", err)
				}

				if err = os.WriteFile(fmt.Sprintf("./test_jsons/returned_%s.json", testCase.CaseName), returnedFile, 0644); err != nil {
					t.Errorf("write returned iiko order v2 error: %v", err)
				}
			} else {
				log.Info().Msgf("ENVIRONMENT=github")
			}

		}).Return(testCase.iikoReturnValue, nil)

		res2, err := iikoSvc.CreateOrder(ctx, testCase.Order, config.Configuration{}, testCase.Store, testCase.PosMenu, menuCli, testCase.AggregatorMenu, nil, nil, nil)
		if err != nil {
			t.Errorf("create order error: %v", err)
			return
		}

		_ = res2
	}

}
