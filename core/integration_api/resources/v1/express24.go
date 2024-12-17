package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/express24/models"
	"net/http"
	"strings"
)

func (server *Server) express24SecretMiddleware(secretToken string) func(c *gin.Context) {
	return func(c *gin.Context) {
		auth := c.GetHeader(Authorization)
		// If the secret is empty...
		if auth == "" {
			// If we get here, the required secret is missingc
			server.Logger.Info(errAuthIsMissing)
			c.Set(errorKey, errors.ErrTokenIsNotValid)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorResponse{
				Msg: errors.ErrTokenIsNotValid.Error(),
			})
			return
		}

		arr := strings.Split(auth, " ")
		if len(arr) != 2 {
			server.Logger.Info(errAuthTokenIsNotValid)
			c.Set(errorKey, errors.ErrTokenIsNotValid)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorResponse{
				Msg: errors.ErrTokenIsNotValid.Error(),
			})
			return
		}

		// Check secret is valid
		if arr[1] != secretToken {
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

// ReceiveOrder docs
//	@Tags		express24
//	@Summary	receive order
//	@Param		event	body		models.Event	true	"event"
//	@Success	200		{object}	string
//	@Failure	400		{object}	errors.ErrorResponse
//	@Router		/v1/express24/order-receive [post]
func (server *Server) ReceiveOrder(c *gin.Context) {
	var (
		body []byte
		err  error
	)

	if server.withRedirect() {
		body = server.readBodyAndSetAgain(c)
	}

	var req models.Event

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof("bind error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("bind error: %s", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if req.OrderChanged == nil {
		server.Logger.Infof("order changed body is nil")
		c.Set(errorKey, "express24 order changed body is nil")
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: "express24 order changed body is nil",
		})
		return
	}

	if req.OrderChanged.Status == "not_paid" {
		server.Logger.Infof("received not paid order, id %d", req.OrderChanged.Id)
		c.Set(errorKey, fmt.Sprintf("not paid order %d", req.OrderChanged.Id))
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: "not paid order",
		})
		return
	}

	if server.withRedirect() && !server.isTopPartner(c.Request.Context(), req.OrderChanged.Store.Branch.ExternalId, models.EXPRESS24.String()) {
		server.redirectRequest(c, body)
		return
	}

	res, err := server.orderService.CreateOrder(c, req.OrderChanged.Store.Branch.ExternalId, models.EXPRESS24.String(), *req.OrderChanged, "")
	if err != nil {
		server.Logger.Infof("create order error: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(http.StatusOK, "")
		return
	}

	c.JSON(http.StatusOK, res.OrderID)
}
