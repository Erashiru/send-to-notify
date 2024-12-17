package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/service/payment/wooppay/dto"
	"go.uber.org/zap"
	"net/http"
)

// WoopPayWebhookEvent docs
//
//	@Tags		wooppay
//	@Title		Method for process webhooks
//	@Summary	Method for process webhooks
//	@Param		event	body		dto.WebhookEvent	true	"event"
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/wooppay/webhooks [post]
func (server *Server) WoopPayWebhookEvent(c *gin.Context) {
	server.Logger.Info(zap.Any("Wooppay event request body", *c.Request))

	var req dto.WebhookEventRequest

	if err := c.Bind(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.WebhookEventResponse{
			Data: 1,
		})
		return
	}

	err := server.paymentManager.WoopPayWebhookEvent(c.Request.Context(), req)
	if err != nil {
		server.Logger.Errorf("Wooppay webhook event error: %s", err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.WebhookEventResponse{
			Data: 1,
		})
		return
	}

	c.JSON(http.StatusOK, dto.WebhookEventResponse{
		Data: 1,
	})
}
