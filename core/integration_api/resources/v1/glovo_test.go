package v1

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kwaaka-team/orders-core/core/glovo/models"
	mock_check "github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/mocks"
	models_check "github.com/kwaaka-team/orders-core/core/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCreateOrderGlovo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	oService := mock_check.NewMockCreationService(ctrl)

	openFl, err := os.Open("form_order/glovo_test_order.json")
	if err != nil {
		t.Fatal(err)
	}
	defer openFl.Close()

	body := models.Order{}
	decoder := json.NewDecoder(openFl)
	if err = decoder.Decode(&body); err != nil {
		t.Fatal(err)
	}
	bodyJS, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/placeOrder", strings.NewReader(string(bodyJS)))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cases := []struct {
		name           string
		mock           func()
		expectedStatus int
	}{
		{
			name: "ok",
			mock: func() {
				oService.EXPECT().CreateOrder(context.Background(), "test-ansar", gomock.Any(), gomock.Any(), gomock.Any()).
					Return(models_check.Order{ID: "12343214"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mock != nil {
				tc.mock()
			}
			server := Server{
				Logger:       lg.Sugar(),
				orderService: oService,
			}
			server.CreateOrderGlovo(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestCancelOrderGlovo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	glovoManager := mock_check.NewMockOrder(ctrl)

	body := models.CancelOrderRequest{
		OrderID:         "ALLE3ZWVY",
		StoreID:         "test-ansar",
		CancelReason:    "Customer changed his mind",
		PaymentStrategy: "CASH",
	}
	bodyJS, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/cancelOrder", strings.NewReader(string(bodyJS)))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	cases := []struct {
		name           string
		mock           func()
		expectedStatus int
	}{
		{
			name: "cancelled",
			mock: func() {
				glovoManager.EXPECT().CancelOrder(context.Background(), body).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mock != nil {
				tc.mock()
			}
			server := Server{
				Logger:       lg.Sugar(),
				glovoManager: glovoManager,
			}
			server.CancelOrderGlovo(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
