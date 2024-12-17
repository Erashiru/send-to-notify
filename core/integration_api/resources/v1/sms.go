package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	"github.com/kwaaka-team/orders-core/core/models"
	"net/http"
)

func (srv *Server) SendVerificationCodeBySms(c *gin.Context) {
	var req models.SendUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if code, restGroupName, err := srv.SmsService.SendVerificationCode(c.Request.Context(), req.PhoneNum, req.RestaurantGroupId); err != nil {
		srv.Logger.Infof("error in sms service %s, error with number %s in restaurant group %s", err, req.PhoneNum, req.RestaurantGroupId)

		if ok, _, _ := srv.SmsService.IsSmsServiceError(err); ok {
			srv.Logger.Infof("attempting to send WhatsApp message for number: %s", req.PhoneNum)

			if err = srv.WhatsappService.SendMessage(c.Request.Context(), req.PhoneNum, fmt.Sprintf("%s: код подтверждения в %s", code, restGroupName), ""); err != nil {
				srv.Logger.Error("error sending wpp message: %s", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Message: fmt.Sprintf("wpp and sms services could not send message %s", err)})
				return
			}

			c.Status(http.StatusNoContent)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
