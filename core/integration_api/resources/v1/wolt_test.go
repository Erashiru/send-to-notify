package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mock_check "github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/mocks"
	models_check "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/wolt/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateOrderWolt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	oService := mock_check.NewMockCreationService(ctrl)
	webhook := models.OrderNotification{
		Id:        "123",
		Type:      "",
		CreatedAt: time.Time{},
		Body: models.OrderNotificationBody{
			Id:          "123",
			ResourceUrl: "",
			Status:      "CREATED",
			VenueId:     "123",
		},
	}
	bodyJS, err := json.Marshal(webhook)
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
				oService.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), models.WOLT.String(), gomock.Any(), gomock.Any()).
					Return(models_check.Order{OrderID: "123"}, nil)
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
				orderService: oService,
				Logger:       lg.Sugar(),
			}
			server.CreateOrderWolt(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestCancelOrderWolt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	woltMan := mock_check.NewWoltMockOrder(ctrl)

	webhook := models.OrderNotification{
		Id:        "",
		Type:      "",
		CreatedAt: time.Time{},
		Body: models.OrderNotificationBody{
			Id:          "123",
			ResourceUrl: "",
			Status:      "CANCELED",
			VenueId:     "123",
		},
	}
	bodyJS, err := json.Marshal(webhook)
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
			name: "cancel",
			mock: func() {
				woltMan.EXPECT().CancelOrder(gomock.Any(), webhook).Return("", nil)
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
				Logger:      lg.Sugar(),
				woltManager: woltMan,
			}
			server.CreateOrderWolt(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
