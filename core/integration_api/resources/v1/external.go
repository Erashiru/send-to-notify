package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/core/externalapi/resources/http/v1/dto"
	"github.com/kwaaka-team/orders-core/core/externalapi/utils"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func (server *Server) authorizeJWT(appSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.Split(authHeader, "Bearer ")

		if len(tokenString) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedResponse{
				Reason: "Access token has been expired. You should request a new one",
			})
			return
		}

		jwtService := utils.JWTAuthService(appSecret)
		token, claims, err := jwtService.ValidateToken(tokenString[1])

		if err != nil {
			log.Err(err).Msg("Couldnt validate access token.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedResponse{
				Reason: "Access token has been expired. You should request a new one",
			})
			return
		}

		if !token.Valid {
			log.Err(err).Msg("Token is invalid.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.UnauthorizedResponse{
				Reason: "Access token has been expired. You should request a new one",
			})
			return
		}

		c.Set("client_id", claims.ClientID)
		c.Set("client_secret", claims.ClientSecret)
		c.Set("service", claims.Service)
	}
}

// CreateOrder docs
//
//	@Tags		external
//	@Title		create order
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string			true	"bearer"
//	@Param		order			body		models.Order	true	"order"
//	@Success	200				{object}	models.CreationResult
//	@Failure	401				{object}	[]errors.ErrorResponse
//	@Failure	400				{object}	[]errors.ErrorResponse
//	@Failure	500				{object}	[]errors.ErrorResponse
//	@Router		/v1/order [post]
func (server *Server) CreateOrder(c *gin.Context) {
	var req models.Order

	svc, ok := c.Get("service")
	if !ok {
		server.Logger.Infof("service query is empty")
		c.Set(errorKey, "service query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown delivery service",
		}})
		return
	}

	secret, ok := c.Get("client_secret")
	if !ok {
		server.Logger.Infof("client_secret query is empty")
		c.Set(errorKey, "client_secret query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown client secret",
		}})
		return
	}

	service := svc.(string)
	clientSecret := secret.(string)

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof("bind error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("bind error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	order, err := server.orderService.CreateOrder(c.Request.Context(), req.RestaurantId, service, req, clientSecret)
	if err != nil {
		server.Logger.Infof("create order error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("create order error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	res := models.CreationResult{
		Result:  "OK",
		OrderId: order.ID,
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}

// GetOrder docs
//
//	@Tags		external
//	@Title		get order
//	@Param		Authorization	header		string	true	"bearer"
//	@Param		order_id		path		string	true	"order_id"
//	@Success	200				{object}	models.Order
//	@Failure	401				{object}	[]errors.ErrorResponse
//	@Failure	400				{object}	[]errors.ErrorResponse
//	@Failure	500				{object}	[]errors.ErrorResponse
//	@Router		/v1/order/{order_id} [get]
func (server *Server) GetOrder(c *gin.Context) {
	orderID := c.Param(orderPath)

	svc, ok := c.Get("service")
	if !ok {
		server.Logger.Infof("service query is empty")
		c.Set(errorKey, "service query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown delivery service",
		}})
		return
	}

	service := svc.(string)

	order, err := server.externalOrderManager.GetOrder(c.Request.Context(), orderID, service)
	if err != nil {
		server.Logger.Infof("get order error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("get order error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	c.Writer.Header().Set("Content-Type", "application/vnd.eats.order.v2+json")
	c.AbortWithStatusJSON(http.StatusOK, order)
}

// CancelOrder docs
//
//	@Tags		external
//	@Title		cancel order
//	@Param		Authorization	header	string	true	"bearer"
//	@Param		order_id		path	string	true	"order_id"
//	@Success	200
//	@Failure	401	{object}	[]errors.ErrorResponse
//	@Failure	400	{object}	[]errors.ErrorResponse
//	@Failure	500	{object}	[]errors.ErrorResponse
//	@Router		/v1/order/{order_id} [delete]
func (server *Server) CancelOrder(c *gin.Context) {
	orderID := c.Param(orderPath)

	svc, ok := c.Get("service")
	if !ok {
		server.Logger.Infof("service query is empty")
		c.Set(errorKey, "service query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown delivery service",
		}})
		return
	}

	secret, ok := c.Get("client_secret")
	if !ok {
		server.Logger.Infof("client_secret query is empty")
		c.Set(errorKey, "client_secret query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown client secret",
		}})
		return
	}

	service := svc.(string)
	clientSecret := secret.(string)

	var req models.CancelOrderRequest
	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof("cancel order error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("cancel order error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	if err := server.externalOrderManager.CancelOrder(c.Request.Context(), req, orderID, service, clientSecret); err != nil {
		server.Logger.Infof("cancel order error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("cancel order error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error()},
		})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// UpdateOrder docs
//
//	@Tags		external
//	@Title		update order
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string			true	"bearer"
//	@Param		order_id		path		string			true	"order_id"
//	@Param		order			body		models.Order	true	"order"
//	@Success	200				{object}	models.Order
//	@Failure	401				{object}	[]errors.ErrorResponse
//	@Failure	400				{object}	[]errors.ErrorResponse
//	@Failure	500				{object}	[]errors.ErrorResponse
//	@Router		/v1/order/{order_id} [put]
func (server *Server) UpdateOrder(c *gin.Context) {
	orderID := c.Param(orderPath)

	svc, ok := c.Get("service")
	if !ok {
		server.Logger.Infof("service query is empty")
		c.Set(errorKey, "service query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown delivery service",
		}})
		return
	}

	secret, ok := c.Get("client_secret")
	if !ok {
		server.Logger.Infof("client_secret query is empty")
		c.Set(errorKey, "client_secret query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown client secret",
		}})
		return
	}

	service := svc.(string)
	clientSecret := secret.(string)

	var req models.Order

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof("bind error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("bind error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	err := server.externalOrderManager.UpdateOrder(c.Request.Context(), req, orderID, service, clientSecret)
	if err != nil {
		server.Logger.Infof("update order error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("update order error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"result": "OK",
	})
}

// GetOrderStatus docs
//
//	@Tags		external
//	@Title		get order status
//	@Param		Authorization	header		string	true	"bearer"
//	@Param		order_id		path		string	true	"order_id"
//	@Success	200				{object}	models.OrderStatusResponse
//	@Failure	401				{object}	[]errors.ErrorResponse
//	@Failure	400				{object}	[]errors.ErrorResponse
//	@Failure	500				{object}	[]errors.ErrorResponse
//	@Router		/v1/order/{order_id}/status [get]
func (server *Server) GetOrderStatus(c *gin.Context) {
	orderID := c.Param(orderPath)

	svc, ok := c.Get("service")
	if !ok {
		server.Logger.Infof("service query is empty")
		c.Set(errorKey, "service query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown delivery service",
		}})
		return
	}

	service := svc.(string)

	status, err := server.externalOrderManager.GetOrderStatus(c.Request.Context(), orderID, service)
	if err != nil {
		server.Logger.Infof("get order status error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("get order status error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error()},
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, status)
}
