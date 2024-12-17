package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/glovo/models"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (server *Server) secretMiddleware(secretToken string) func(c *gin.Context) {
	return func(c *gin.Context) {
		auth := c.GetHeader(Authorization)

		// If the secret is empty...
		if auth == "" {
			// If we get here, the required secret is missing
			server.Logger.Info(errAuthIsMissing)
			c.Set(errorKey, errors.ErrTokenIsNotValid)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorResponse{
				Msg: errors.ErrTokenIsNotValid.Error(),
			})
			return
		}

		// Check secret is valid
		if auth != secretToken {
			server.Logger.Info(errAuthTokenIsNotValid)
			c.Set(errorKey, errors.ErrTokenIsNotValid)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorResponse{
				Msg: errors.ErrTokenIsNotValid.Error(),
			})
			return
		}

		c.Set(Token, auth)
	}

}

// CreateOrderGlovo docs
//
//	@Tags		glovo
//	@Title		Method for create order
//	@Security	ApiKeyAuth
//	@Summary	Method create Order
//	@Param		order	body		models.Order	true	"order"
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/glovo/placeOrder [post]
func (server *Server) CreateOrderGlovo(c *gin.Context) {
	var (
		body []byte
		err  error
	)

	if server.withRedirect() {
		body = server.readBodyAndSetAgain(c)
	}

	var req models.Order

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if server.withRedirect() && !server.isTopPartner(c.Request.Context(), req.StoreID, models.GLOVO.String()) {
		server.redirectRequest(c, body)
		return
	}

	res, err := server.orderService.CreateOrder(c.Request.Context(), req.StoreID, models.GLOVO.String(), req, "")
	if err != nil {
		server.Logger.Infof("create order error: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(http.StatusOK, "")
		return
	}

	c.JSON(http.StatusOK, res.ID)
}

// CancelOrderGlovo docs
//
//	@Tags		glovo
//	@Title		Method to cancel order
//	@Security	ApiKeyAuth
//	@Summary	Method cancel Order
//	@Param		order	body		models.Order	true	"order"
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/glovo/cancelOrder [post]
func (server *Server) CancelOrderGlovo(c *gin.Context) {

	var req models.CancelOrderRequest

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	log.Info().Msgf("Glovo Order Cancel request: %+v", req)

	err := server.glovoManager.CancelOrder(c.Request.Context(), req)
	if err != nil {
		server.Logger.Infof("cancel order error: %s", err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
