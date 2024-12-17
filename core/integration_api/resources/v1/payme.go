package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"github.com/kwaaka-team/orders-core/service/payment/payme/dto"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
)

// PaymeWebhookEvent docs
//	@Tags		payme
//	@Title		Method for process webhooks
//	@Summary	Method for process webhooks
//	@Param		event	body		dto.WebhookEvent	true	"event"
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/payme/webhooks [post]
func (server *Server) PaymeWebhookEvent(c *gin.Context) {
	server.Logger.Info(zap.Any("Payme event request body", *c.Request))

	var req dto.WebhookEvent

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Jsonrpc: "2.0",
			Error: dto.Error{
				Message: err.Error(),
				Code:    -31008,
			},
		})
		return
	}

	resp, err := server.paymentManager.WebhookEvent(c.Request.Context(), models.PAYME, req)
	if err != nil {
		server.Logger.Errorf("payme %s webhook event error: %s", req.Method, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Jsonrpc: "2.0",
			Error: dto.Error{
				Message: err.Error(),
				Code:    -31008,
			},
			ID: req.ID,
		})
		return
	}

	response, ok := resp.(dto.WebhookResultResponse)
	if !ok {
		server.Logger.Infof("payme %s webhook event cast response error: %s", req.Method, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Jsonrpc: "2.0",
			Error: dto.Error{
				Message: errors.New("casting response error").Error(),
				Code:    -31008,
			},
			ID: req.ID,
		})
		return
	}

	response.ID = req.ID

	c.JSON(http.StatusOK, response)
}
