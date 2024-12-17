package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/qrmenu/models"
	"github.com/kwaaka-team/orders-core/domain/logger"
	"github.com/kwaaka-team/orders-core/service/payment/ioka/dto"
	models2 "github.com/kwaaka-team/orders-core/service/payment/models"
	"net/http"
	"net/url"
	"strconv"
)

// CreateOrderQRmenu docs
//
//	@Tags		qrmenu
//	@Title		Method for create order
//	@Security	ApiKeyAuth
//	@Summary	Method create Order
//	@Param		order	body		models.Order	true	"order"
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/qr-menu/placeOrder [post]
func (server *Server) CreateOrderQRmenu(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "qr_menu create order request",
		Request: c.Request,
	})

	var req models.Order

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	server.Logger.Infof("starting create qr_menu order: ID: %s, rest id: %s, customer name: %s, payment type: %s, dispatcher: %s, client delivery price: %f, full delivery price: %f, kwaaka charged delivery price: %f",
		req.ID, req.RestaurantID, req.Customer.Name, req.PaymentType, req.Delivery.Dispatcher, req.Delivery.ClientDeliveryPrice, req.Delivery.FullDeliveryPrice, req.Delivery.KwaakaChargedDeliveryPrice)

	res, err := server.orderService.CreateOrder(c.Request.Context(), req.RestaurantID, "qr_menu", req, "")
	if err != nil {
		server.Logger.Errorf("create order error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if res.IsInstantDelivery && res.DeliveryDispatcher != "" {
		if err = server.orderKwaaka3plService.Instant3plOrder(c.Request.Context(), res); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
				Msg: fmt.Sprintf("order created but Instant 3pl did not: %s", err.Error()),
			})
		}
	}

	c.JSON(http.StatusOK, res.ID)
}

// OpenApplePaySessionByQRmenu docs
//
//	@Tags		qrmenu
//	@Title		Method for open apple pay session
//	@Security	ApiKeyAuth
//	@Summary	Method open apple pay session
//	@Param		request	body		dto.ApplePaySessionOpenRequest	true	"request"
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/qrmenu/applePay/session [post]
func (server *Server) OpenApplePaySessionByQRmenu(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "qrmenu open apple pay session request",
		Request: c.Request,
	})

	var req dto.ApplePaySessionOpenRequest

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	res, err := server.paymentManager.OpenApplePaySession(c.Request.Context(), models2.ApplePaySessionOpenRequest{
		PaymentSystem: models2.IOKA,
		OrderID:       req.OrderID,
		DomainName:    req.DomainName,
		Url:           req.Url,
		Platform:      req.Platform,
	})
	if err != nil {
		server.Logger.Errorf("open apple pay session: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (server *Server) CreateApplePayPayment(c *gin.Context) {
	var req models2.ApplePayPayment

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	paymentOrderID := c.Param("payment_order_id")

	err := server.paymentManager.CreatePaymentByApplePayAndSavePaymentDetails(c.Request.Context(), models2.ApplePayPayment{
		PaymentSystem: models2.IOKA,
		ToolType:      req.ToolType,
		ApplePay:      req.ApplePay,
		OrderID:       paymentOrderID,
	})
	if err != nil {
		server.Logger.Errorf("create apple pay payment: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "")
}

func (s *Server) GetTwoGisReviewLink(ctx *gin.Context) {
	s.Logger.Info(logger.LoggerInfo{
		System:  "get two gis link",
		Request: ctx.Request,
	})

	restID := ctx.Param("restaurant_id")
	if restID == "" {
		ctx.Set(errorKey, errors.ErrorResponse{Msg: "restaurant id is empty"})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: "restaurant id is empty"})
		return
	}

	link, err := s.storeService.GetTwoGisReviewLink(ctx.Request.Context(), restID)
	if err != nil {
		ctx.Set(errorKey, err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, link)
}

// NotifyAllUnpaidCustomers
//
// This handler is called by notify_unpaid_customers cron
func (s *Server) NotifyAllUnpaidCustomers(c *gin.Context) {
	ctx := c.Request.Context()
	var errArray errors.Error

	minutesBeforeCheckStr := c.Query("minutes_before_check")
	notificationCountStr := c.Query("notification_count")
	minutesBeforeCheck, err := strconv.Atoi(minutesBeforeCheckStr)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}
	notificationCount, err := strconv.Atoi(notificationCountStr)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	payments, err := s.paymentManager.GetUnpaidPayments(ctx, minutesBeforeCheck)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var message string
	for _, payment := range payments {
		message = url.QueryEscape(fmt.Sprintf("💳 Пожалуйста, оплатите вашу корзину.\n\n%s\n\n💲 Неоплаченная сумма - %d", payment.CheckoutURL, payment.Amount/100))
		if err = s.WhatsappService.SendMessage(ctx, payment.CustomerPhoneNumber, message, ""); err != nil {
			s.Logger.Errorf("error sending whatsapp message to unpaid customer: %v", err)
			errArray.Append(err)
		}
		if err = s.paymentManager.SetNotificationCount(ctx, payment.ExternalID, notificationCount); err != nil {
			s.Logger.Errorf("error setting notification count: %v", err)
			errArray.Append(err)
		}
	}
	paymentIDs := make([]string, 0, len(payments))
	if len(payments) > 0 {
		for _, payment := range payments {
			paymentIDs = append(paymentIDs, payment.ExternalID)
		}
	}

	if errArray.ErrorOrNil() != nil {
		c.JSON(http.StatusInternalServerError, errArray)
		return
	}

	c.JSON(http.StatusOK, paymentIDs)
}
