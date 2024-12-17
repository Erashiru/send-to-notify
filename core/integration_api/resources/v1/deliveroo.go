package v1

import (
	"github.com/gin-gonic/gin"
	models2 "github.com/kwaaka-team/orders-core/core/deliveroo/models"
	"github.com/kwaaka-team/orders-core/core/errors"
	"net/http"
)

// OrderEventDeliveroo docs
//	@Tags		deliveroo
//	@Summary	order event deliveroo
//	@Param		order_event	body		models2.OrderEvent	true	"order_event"
//	@Success	200			{object}	string
//	@Failure	400			{object}	errors.ErrorResponse
//	@Router		/v1/deliveroo/order-events [post]
func (server *Server) OrderEventDeliveroo(c *gin.Context) {

	var req models2.OrderEvent

	if err := c.BindJSON(&req); err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}
	res, err := server.deliverooManager.OrderEvent(c.Request.Context(), req)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// MenuEventDeliveroo docs
//	@Tags		deliveroo
//	@Summary	menu event deliveroo
//	@Param		menu_event	body	models2.MenuEvent	true	"menu_event"
//	@Success	200
//	@Failure	400	{object}	errors.ErrorResponse
//	@Router		/v1/deliveroo/menu-events [post]
func (server *Server) MenuEventDeliveroo(c *gin.Context) {
	var req models2.MenuEvent

	if err := c.BindJSON(&req); err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	err := server.deliverooManager.MenuEvent(c.Request.Context(), req)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
