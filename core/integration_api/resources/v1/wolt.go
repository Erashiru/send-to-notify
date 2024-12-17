package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/wolt/models"
	"net/http"
)

// CreateOrderWolt docs
//	@Tags		wolt
//	@Title		Method for create order
//	@Security	ApiKeyAuth
//	@Summary	Method create Order
//	@Param		order	body		models.Order	true	"order"
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/wolt/placeOrder [post]
func (server *Server) CreateOrderWolt(c *gin.Context) {
	var body []byte

	if server.withRedirect() {
		body = server.readBodyAndSetAgain(c)
	}

	var webhook models.OrderNotification

	if err := c.BindJSON(&webhook); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if server.withRedirect() && !server.isTopPartner(c.Request.Context(), webhook.Body.VenueId, models.WOLT.String()) {
		server.redirectRequest(c, body)
		return
	}

	if models.WoltStatus(webhook.Body.Status) == models.Created {
		res, err := server.orderService.CreateOrder(c.Request.Context(), webhook.Body.VenueId, models.WOLT.String(), webhook, "")
		if err != nil {
			server.Logger.Infof("create order error: %s", err.Error())
			c.Set(errorKey, err)
			c.JSON(http.StatusOK, "")
			return
		}

		c.JSON(http.StatusOK, res.OrderID)
		return
	}

	if webhook.Body.Status == models.Canceled.String() {
		res, err := server.woltManager.CancelOrder(c.Request.Context(), webhook)
		if err != nil {
			server.Logger.Infof("create order error: %s", err.Error())
			c.Set(errorKey, err)
			c.JSON(http.StatusOK, "")
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	c.JSON(http.StatusOK, "")
}
