package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/service/payment/ioka/dto"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"go.uber.org/zap"
	"net/http"
)

// IokaWebhookEvent docs
//	@Tags		ioka
//	@Title		Method for process webhooks
//	@Summary	Method for process webhooks
//	@Param		event	body		dto.WebhookEvent	true	"event"
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/ioka/webhooks [post]
func (server *Server) IokaWebhookEvent(c *gin.Context) {
	server.Logger.Info(zap.Any("Ioka event request body", *c.Request))

	var req dto.WebhookEvent

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if _, err := server.paymentManager.WebhookEvent(c.Request.Context(), models.IOKA, req); err != nil {
		server.Logger.Infof("ioka %s webhook event error: %s", req.Event, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "")
}
