package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	ext_models "github.com/kwaaka-team/orders-core/core/externalapi/models"
	mock_check "github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/mocks"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	oService := mock_check.NewMockCreationService(ctrl)

	openFl, err := os.Open("../v1/form_order/yandex_test_order.json")
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
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(bodyJS)))
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
				oService.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(models.Order{OrderID: "123"}, nil)
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
			c.Set("service", "fsagdasfds")
			c.Set("client_secret", "fsagdasfds")
			server.CreateOrder(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestCancelOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lg, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	oService := mock_check.NewMockOrderClient(ctrl)

	openFl, err := os.Open("../v1/form_order/yandex_test_order.json")
	if err != nil {
		t.Fatal(err)
	}
	defer openFl.Close()

	body := ext_models.CancelOrderRequest{}
	decoder := json.NewDecoder(openFl)
	if err = decoder.Decode(&body); err != nil {
		t.Fatal(err)
	}
	bodyJS, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(bodyJS)))
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
				oService.EXPECT().CancelOrder(gomock.Any(), body, gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
				Logger:               lg.Sugar(),
				externalOrderManager: oService,
			}
			c.Set("service", "fsagdasfds")
			c.Set("client_secret", "fsagdasfds")
			server.CancelOrder(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
