package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/integration_api/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"net/http"
)

// ExternalPosIntegrationStopList docs
//
//	@Tags		integration-pos
//	@Title		Method for update stoplist
//	@Security	ApiKeyAuth
//	@Summary	Method for update stoplist
//	@Param		stoplist	body		models.StopListRequest	true	"stoplist"
//	@Failure	401			{object}	errors.ErrorResponse
//	@Failure	400			{object}	errors.ErrorResponse
//	@Failure	500			{object}	errors.ErrorResponse
//	@Router		/integration/pos/stoplist [post]
func (server *Server) ExternalPosIntegrationStopList(c *gin.Context) {
	token := c.GetHeader(Authorization)
	if token == "" {
		server.Logger.Infof("Authorization header is missing")
		c.Set(errorKey, "Authorization header is missing")
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Authorization header is missing",
		})
		return
	}

	var request models.StopListRequest

	if err := c.BindJSON(&request); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	utils.Beautify("stop list request body", request)

	if err := server.externalPosIntegrationManager.UpdateStopList(c.Request.Context(), token, request, "ctmax"); err != nil {
		server.Logger.Infof("update stoplist error: %v", err)
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	server.Logger.Infof("successful stoplist update")

	c.Status(http.StatusNoContent)
}

// ExternalPosIntegrationUpdateOrder docs
//
//	@Tags		integration-pos
//	@Title		Method for update order
//	@Security	ApiKeyAuth
//	@Summary	Method for update order
//	@Param		order	body		models.UpdateOrderStatusRequest	true	"order"
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/integration/pos/order [post]
func (server *Server) ExternalPosIntegrationUpdateOrder(c *gin.Context) {
	token := c.GetHeader(Authorization)
	if token == "" {
		server.Logger.Infof("Authorization header is missing")
		c.Set(errorKey, "Authorization header is missing")
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Authorization header is missing",
		})
		return
	}

	var request models.UpdateOrderStatusRequest

	if err := c.BindJSON(&request); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	utils.Beautify("update order status request body", request)

	if err := server.externalPosIntegrationManager.UpdateOrderStatus(c.Request.Context(), token, request.RestaurantID, request.OrderID, request.Status, request.Reason); err != nil {
		server.Logger.Infof("update order error: %v", err)
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	server.Logger.Infof("successful order status update")
	c.Status(http.StatusNoContent)
}

// ExternalPosIntegrationGetOrders docs
//
//	@Tags		integration-pos
//	@Title		Method for update stoplist
//	@Security	ApiKeyAuth
//	@Summary	Method for update stoplist
//	@Param		restaurant_id	path		string	true	"restaurant_id"
//	@Success	200				{object}	models.GetOrdersResponse
//	@Failure	401				{object}	errors.ErrorResponse
//	@Failure	400				{object}	errors.ErrorResponse
//	@Failure	500				{object}	errors.ErrorResponse
//	@Router		/integration/pos/{restaurant_id}/orders [get]
func (server *Server) ExternalPosIntegrationGetOrders(c *gin.Context) {
	restaurantID := c.Param("restaurant_id")

	if restaurantID == "" {
		server.Logger.Infof("restaurant_id query is empty")
		c.Set(errorKey, "restaurant_id query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "restaurant_id param is empty",
		})
		return
	}

	token := c.GetHeader(Authorization)
	if token == "" {
		server.Logger.Infof("Authorization header is missing")
		c.Set(errorKey, "Authorization header is missing")
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  "Authorization header is missing",
		})
		return
	}

	status, _ := c.GetQuery("status")

	server.Logger.Infof("request data; query restaurant_id=%s, header Authorization={hidden}", restaurantID)

	ordersResponse, err := server.externalPosIntegrationManager.GetOrders(c.Request.Context(), token, restaurantID, status, "", "")
	if err != nil {
		server.Logger.Infof("get orders error: %v", err)
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Code: http.StatusBadRequest,
			Msg:  err.Error(),
		})
		return
	}

	utils.Beautify("integration get orders response", ordersResponse)

	c.JSON(http.StatusOK, ordersResponse)
}
