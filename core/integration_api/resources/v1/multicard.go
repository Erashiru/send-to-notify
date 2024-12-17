package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"github.com/kwaaka-team/orders-core/service/payment/multicard/dto"
	"go.uber.org/zap"
	"net/http"
)

func (server *Server) MulticardWebhook(c *gin.Context) {
	server.Logger.Info(zap.Any("Multricard event request body", *c.Request))

	var req dto.Webhook
	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if _, err := server.paymentManager.WebhookEvent(c.Request.Context(), models.MultiCard, req); err != nil {
		server.Logger.Infof("multicard webhook event error for number %s: %s", req.Phone, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "")
}
