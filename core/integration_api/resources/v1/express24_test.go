package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kwaaka-team/orders-core/core/express24/models"
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

func TestReceiveOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	lg, err := zap.NewDevelopment()
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	oService := mock_check.NewMockCreationService(ctrl)

	openFl, err := os.Open("../v1/form_order/express24_test_order.json")
	if err != nil {
		t.Fatal(err)
	}
	defer openFl.Close()
	body := models.Event{}
	decoder := json.NewDecoder(openFl)
	if err = decoder.Decode(&body); err != nil {
		t.Fatal(err)
	}
	bodyJS, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/order-receive", strings.NewReader(string(bodyJS)))
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
				oService.EXPECT().CreateOrder(c, body.OrderChanged.Store.Branch.ExternalId, models.EXPRESS24.String(), *body.OrderChanged, "").
					Return(models_check.Order{ID: "123"}, nil)
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
			server.ReceiveOrder(c)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
