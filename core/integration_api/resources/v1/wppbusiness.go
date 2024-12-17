package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/models"
	"net/http"
)

func (srv *Server) SendVerificationCode(c *gin.Context) {
	var req models.SendUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := srv.WppBusinessService.SendVerificationCode(c.Request.Context(), req.PhoneNum, req.RestaurantGroupId); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
