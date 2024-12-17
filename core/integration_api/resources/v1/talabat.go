package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/talabat/models"
	"github.com/kwaaka-team/orders-core/domain/logger"
	"net/http"
)

// CancelOrderTalabat docs
//	@Tags		talabat
//	@Summary	cancel order talabat
//	@Param		cancel_order_request	body	models.CancelOrderRequest	true	"cancel_order_request"
//	@Param		remoteId				path	string						true	"remoteId"
//	@Param		remoteOrderId			path	string						true	"remoteOrderOd"
//	@Success	200
//	@Failure	400
//	@Router		/v1/talabat/remoteId/{remoteId}/remoteOrder/{remoteOrderId}/posOrderStatus [post]
func (server *Server) CancelOrderTalabat(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "talabat cancel order request",
		Request: c.Request,
	})

	var req models.CancelOrderRequest

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		server.Logger.Error(logger.LoggerInfo{
			System:   "talabat order cancel response error",
			Response: err,
			Status:   http.StatusBadRequest,
		})
		return
	}
	req.RemoteID = c.Param(remoteId)
	req.RemoteOrderID = c.Param(remoteOrderId)

	if err := server.talabatOrderManager.CancelOrder(c.Request.Context(), req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

// CreateOrderTalabat docs
//	@Tags		talabat
//	@Summary	create order talabat
//	@Param		create_order_request	body		models.CreateOrderRequest	true	"create_order_request"
//	@Param		remoteId				path		string						true	"remoteId"
//	@Success	200						{object}	models.CreateOrderResponse
//	@Failure	400
//	@Router		/v1/talabat/order/{remoteId} [post]
func (server *Server) CreateOrderTalabat(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "talabat create order request",
		Request: c.Request,
	})

	var req models.CreateOrderRequest

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		server.Logger.Error(logger.LoggerInfo{
			System:   "talabat order create response error",
			Response: err,
			Status:   http.StatusBadRequest,
		})
		return
	}

	remoteID := c.Param(remoteId)

	resp, err := server.orderService.CreateOrder(c.Request.Context(), remoteID, "talabat", req, "")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	server.Logger.Info(logger.LoggerInfo{
		System:   "talabat create order response",
		Response: resp,
		Status:   http.StatusOK,
	})

	c.JSON(http.StatusOK, resp)
}

// MenuUploadCallbackTalabat docs
//	@Tags		talabat
//	@Summary	menu upload callback talabat
//	@Param		menu_upload_callback_request	body	models.MenuUploadCallbackRequest	true	"menu_upload_callback_request"
//	@Success	200
//	@Failure	200
//	@Router		/v1/talabat/requestResult [post]
func (server *Server) MenuUploadCallbackTalabat(c *gin.Context) {
	var req models.MenuUploadCallbackRequest

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Error(logger.LoggerInfo{
			System:   "talabat menu upload callback response error",
			Response: err,
		})
		c.Status(http.StatusOK)
		return
	}

	if err := server.talabatMenuManager.UpdateMenuUploadTransaction(c.Request.Context(), req); err != nil {
		server.Logger.Error(logger.LoggerInfo{
			System:   "talabat menu upload callback response error",
			Response: err,
		})
		c.Status(http.StatusOK)
		return
	}

	c.Status(http.StatusOK)
}
