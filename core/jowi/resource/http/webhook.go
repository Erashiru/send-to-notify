package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/jowi/models"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Events docs
//	@Tags		jowi
//	@Summary	events
//	@Param		event	body		models.Event	true	"event"
//	@Success	200		{object}	nil
//	@Failure	400
//	@Router		/jowi/events [post]
func (server *Server) Events(c *gin.Context) {
	var event models.Event

	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Err(err).Msgf("binding json for webhookType error")
		return
	}

	if err := server.jowiManager.UpdateOrder(c.Request.Context(), event); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Err(err).Msgf("update order error")
		return
	}

	c.JSON(http.StatusOK, nil)

}
