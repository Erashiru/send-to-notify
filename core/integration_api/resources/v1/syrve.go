package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/service/iiko/resources/http/v1/dto"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"go.uber.org/zap"
	"net/http"
)

// EventSyrve docs
//	@Tags		syrve
//	@Title		get events from syrve pos
//	@Accept		json
//	@Produce	json
//	@Param		req	body	models.WebhookEvent	true	"events"
//	@Success	200
//	@Failure	401	{object}	[]dto.ErrorResponse
//	@Failure	400	{object}	[]dto.ErrorResponse
//	@Failure	500	{object}	[]dto.ErrorResponse
//	@Router		/syrve/events [post]
func (server *Server) EventSyrve(c *gin.Context) {

	token := c.GetString(StoreTokenKey)

	if token == "" {
		server.Logger.Info(errAuthTokenIsNotValid)
		c.Set(errorKey, errAuthTokenIsNotValid)
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code: http.StatusUnauthorized,
		})
		return
	}

	var req models.WebhookEvents

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		return
	}

	server.Logger.Info(zap.Any("Syrve event request body", req))

	rsp, err := server.iikoManager.WebhookEvent(c.Request.Context(), token, req, models.SYRVE)
	if err != nil {
		server.Logger.Infof("webhook event error: %s", err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		return
	}

	server.Logger.Info(zap.Any("Syrve event response body", rsp))

	c.JSON(http.StatusOK, models.WebhookEventResponse{
		Details: rsp,
	})
}
