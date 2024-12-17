package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/jowi/models"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func (server *Server) JowiEvents(c *gin.Context) {
	var event models.Event

	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Err(err).Msgf("binding json for webhookType error")
		return
	}

	if err := server.statusUpdateService.UpdateOrderStatus(c.Request.Context(), event.Data.OrderId, strconv.Itoa(event.Status), ""); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Err(err).Msgf("update order error")
		return
	}

	c.Status(http.StatusNoContent)
}
