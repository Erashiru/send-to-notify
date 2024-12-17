package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/service/iiko/resources/http/v1/detector"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"net/http"
	"strconv"
)

// CreatePaymentOrder docs
//
//	@Tags		kwaaka-admin
//	@Title		Method for create payment order
//	@Summary	Method for create payment order
//	@Param		request	body		models.PaymentOrder	true	"request"
//	@Success	200		{object}	models.PaymentOrder
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/payment-order [post]
func (server *Server) CreatePaymentOrder(c *gin.Context) {

	server.Logger.Info("create payment order request")

	var request models.PaymentOrder
	if err := c.BindJSON(&request); err != nil {
		server.Logger.Errorf("create payment order error: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	paymentOrder, err := server.paymentManager.CreatePaymentOrder(c.Request.Context(), request)
	if err != nil {
		server.Logger.Errorf("create payment order error: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	c.JSON(http.StatusOK, paymentOrder)

}

// UpdatePaymentOrderStatus docs
//
//	@Tags		kwaaka-admin
//	@Title		Method for update payment order
//	@Summary	Method for update payment order
//	@Param		request	body	models.UpdatePaymentOrderStatus	true	"request"
//	@Success	200
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/payment-order [put]
func (server *Server) UpdatePaymentOrderStatus(c *gin.Context) {

	server.Logger.Info("update payment order request")

	var request models.UpdatePaymentOrderStatus
	if err := c.BindJSON(&request); err != nil {
		server.Logger.Errorf("update payment order status: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	server.Logger.Infof("update payment order status request: cart id: %s, status: %s", request.CartID, request.Status)

	if err := server.paymentManager.UpdatePaymentOrderByOrderID(c.Request.Context(), request.CartID, request.Status); err != nil {
		server.Logger.Errorf("update payment order status: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	c.Status(http.StatusOK)
}

// RefundPayment
//
//		@Tags		kwaaka-admin
//		@Title		Method for refunding
//		@Summary	Method for refunding amount to customer
//		@Param		order_id	path   string	 true	"order id for finding payment order"
//	    @Param      amount      query  string    true   "amount has to be without added 00 at the end of amount (ioka format)"
//		@Param		reason	    query  string	 false	"optional to get additional data for refund"
//		@Success	200 {object}    models.RefundResponse
//		@Failure	401	{object}	errors.ErrorResponse
//		@Failure	400	{object}	errors.ErrorResponse
//		@Failure	500	{object}	errors.ErrorResponse
//		@Router		/v1/kwaaka-admin/refund/:order_id [post]
func (server *Server) RefundPayment(c *gin.Context) {
	server.Logger.Info("refund payment order request")

	orderID := c.Param("order_id")

	amount, ok := c.GetQuery("amount")
	if !ok {
		c.Set(errorKey, fmt.Errorf("amount not received for refund"))
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "amount has to be specified for refund",
		})
		return
	}
	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	reason, ok := c.GetQuery("reason")
	if ok {
		server.Logger.Infof("reason for refund payment received: %s", reason)
	}

	paymentOrder, refundResponse, err, amountErr := server.paymentManager.RefundToCustomer(c.Request.Context(), orderID, reason, amountInt)
	if err != nil {
		server.Logger.Errorf("refund payment: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}
	if amountErr != nil {
		server.Logger.Errorf("refund payment unavailable for the amount less than 100: %s", amountErr.Error())
		c.Set(errorKey, amountErr)
		c.JSON(http.StatusBadRequest, gin.H{
			"description": "refund amount has to be greater than 100 or equal",
		})
		return
	}

	order, err := server.orderInfoSharingService.GetOrder(c.Request.Context(), paymentOrder.OrderID)
	if err != nil {
		server.Logger.Errorf("couldn't find order to send notification of refund: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	store, err := server.storeService.GetByID(c.Request.Context(), paymentOrder.RestaurantID)
	if err != nil {
		server.Logger.Errorf("couldn't find store to send notification of refund: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	err = server.TelegramService.SendMessageToQueue(telegram.Refund, order, store, "", amount, reason, models2.Product{})
	if err != nil {
		server.Logger.Errorf("couldn't send refund notification to telegram chat: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	c.JSON(http.StatusOK, refundResponse)
}

// GetRefund
//
//	@Tags		kwaaka-admin
//	@Title		Method for get refund
//	@Summary	Method for get refund
//	@Param		order_id	path   string	 true	"order id for finding refund"
//	@Success	200 {object}    models.Refund
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/get-refund/:order_id [get]
func (srv *Server) GetRefund(c *gin.Context) {
	orderID := c.Param("order_id")

	res, err := srv.paymentManager.GetRefund(c.Request.Context(), orderID)
	if err != nil {
		srv.Logger.Errorf("couldn't get refund by order id: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

// CreatePaymentLink
//
//	@Tags     kwaaka-admin
//	@Title    Method for creating payment link for kaspi salescout
//	@Summary  Method for creating payment link for kaspi salescout
//	@Param    order_id  path   string   true  "order id for finding payment order"
//	@Success  200 {object}    models.PaymentOrder
//	@Failure  400  {object}  errors.ErrorResponse
//	@Failure  500  {object}  errors.ErrorResponse
//	@Router    /v1/kwaaka-admin/payment-order/create-payment-link/:order_id [post]
func (server *Server) CreatePaymentLink(c *gin.Context) {
	orderId := c.Param("order_id")

	paymentOrder, err := server.paymentManager.CreatePaymentLinkForCustomerToPay(c.Request.Context(), orderId)
	if err != nil {
		server.Logger.Errorf("create payment link error: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	c.JSON(http.StatusOK, paymentOrder)
}
