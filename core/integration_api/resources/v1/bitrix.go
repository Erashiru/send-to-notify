package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	"net/http"
)

// SendMessageToCustomerFromBitrixLead
//
//	@Tags		bitrix
//	@Title		Method for send message to whatsapp from bitrix lead
//	@Security	ApiKeyAuth
//	@Summary	Method send message to whatsapp from bitrix lead
//	@Param		request	body	dto.BitrixEventDataRequest	true	"request"
//	@Success	204
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/api/bitrix/event [post]
func (server *Server) SendMessageToCustomerFromBitrixLead(c *gin.Context) {
	var event dto.BitrixEventDataRequest

	if err := c.ShouldBind(&event); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := server.bitrixService.SendMessageToCustomerFromBitrixLead(c.Request.Context(), event.DataFieldsID); err != nil {
		server.Logger.Errorf("send message from bitrix customer error: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
