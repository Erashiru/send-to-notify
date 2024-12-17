package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func (server *Server) ytimesSecretMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		auth := c.GetHeader(Authorization)
		splitToken := strings.Split(auth, " ")

		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorResponse{
				Msg: errors.ErrTokenIsNotValid.Error(),
			})
			return
		}

		c.Set(StoreTokenKey, strings.TrimSpace(splitToken[1]))
	}
}

func (server *Server) YTimesUpdateOrderStatus(c *gin.Context) {
	server.Logger.Info(zap.Any("ytimes update order status event request body", *c.Request))

	var (
		req dto.YTimesUpdateOrderStatusBody
	)

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if err := server.statusUpdateService.UpdateOrderStatus(c.Request.Context(), req.Guid, req.Status, req.StatusMessage); err != nil {
		server.Logger.Infof("update order status error: %s", err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.String(http.StatusOK, "OK")
}

func (server *Server) YTimesMenuUpdates(c *gin.Context) {
	server.Logger.Info(zap.Any("ytimes menu updates event request body", *c.Request))

	var (
		req  dto.YTimesWebhookRequest
		body dto.YTimesMenuUpdatesRequestBody
	)

	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		server.Logger.Infof("unmarshal menu updates body error: %s", err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	storeId, ok := c.Get(StoreTokenKey)
	if !ok {
		server.Logger.Info("Authorization bearer token error")
		c.Set(errorKey, fmt.Errorf("authorization bearer token not valid"))
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: "authorization bearer token not valid",
		})
		return
	}

	_ = storeId

	c.String(http.StatusOK, "OK")
}
