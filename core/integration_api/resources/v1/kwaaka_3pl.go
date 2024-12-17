package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	err2 "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/domain/logger"
	"net/http"
)

// SetOrdersDispatcher docs
//
//	@Tags		kwaaka-admin
//	@Title		Method for set dispatcher
//	@Security	ApiKeyAuth
//	@Summary	Method for set dispatcher
//	@Param		set	dispatcher	body	models.SetKwaaka3plDispatcherRequest	true	"dispatcher"
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/kwaaka-admin/setOrdersDispatcher [post]
func (server *Server) SetOrdersDispatcher(c *gin.Context) {
	var req models.SetKwaaka3plDispatcherRequest
	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err2.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	server.Logger.Info(logger.LoggerInfo{
		System:  "kwaaka 3pl set orders dispatcher request",
		Request: req,
	})

	err := server.orderKwaaka3plService.SetKwaaka3plDispatcher(c.Request.Context(), req)
	if err != nil {
		server.Logger.Errorf("set orders kwaaka3pl dispatcher error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err2.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

// CancelOrderDispatcher docs
//
//	@Tags		kwaaka-admin
//	@Title		Method for set dispatcher
//	@Security	ApiKeyAuth
//	@Summary	Method for set dispatcher
//	@Param		order_id	path		string	true	"order_id"
//	@Failure	401			{object}	errors.ErrorResponse
//	@Failure	400			{object}	errors.ErrorResponse
//	@Failure	500			{object}	errors.ErrorResponse
//	@Router		/kwaaka-admin/cancelOrder/{order_id} [put]
func (server *Server) CancelOrderDispatcher(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "kwaaka 3pl cancel order dispatcher request",
		Request: c.Request,
	})
	id, ok := c.GetQuery("order_id")
	if !ok {
		err := fmt.Errorf("order id empty")
		server.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err2.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	err := server.orderKwaaka3plService.Cancel3plOrder(c.Request.Context(), id, models.OrderInfoForTelegramMsg{})
	if err != nil {
		server.Logger.Errorf("cancel order kwaaka3pl dispatcher error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err2.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "")
}

func (server *Server) CancelCourierSearch(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "kwaaka 3pl cancel order dispatcher request",
		Request: c.Request,
	})

	deliveryOrderId := c.Param("delivery_order_id")

	server.Logger.Info("received delivery_order_id: %s", deliveryOrderId)

	if err := server.orderKwaaka3plService.CancelCourierSearch(c.Request.Context(), deliveryOrderId); err != nil {
		c.Set(errorKey, err)
		c.JSON(http.StatusInternalServerError, err2.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	server.Logger.Info("successfully processed cancel request for delivery_order_id: %s", deliveryOrderId)

	c.AbortWithStatus(http.StatusNoContent)
}

func (server *Server) Save3plHistory(c *gin.Context) {

	deliveryOrderId := c.Param("delivery_order_id")
	newCustomerDeliveryAddress := models.SaveToHistory{}
	if err := c.BindJSON(&newCustomerDeliveryAddress); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err2.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	server.Logger.Info(logger.LoggerInfo{
		System:  "save 3pl history and remove current dispatcher service request",
		Request: fmt.Sprintf("delivery order id: %s, new delivery address: %v", deliveryOrderId, newCustomerDeliveryAddress),
	})

	if err := server.orderKwaaka3plService.Save3plHistory(c.Request.Context(), deliveryOrderId, newCustomerDeliveryAddress.DeliveryAddress, newCustomerDeliveryAddress.Customer); err != nil {
		c.Set(errorKey, err)
		c.JSON(http.StatusInternalServerError, err2.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)

}

func (server *Server) MapIIKOstatusTo3plStatus(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "map iiko delivery status to 3pl statuses",
		Request: c.Request,
	})

	var req struct {
		IikoStatus          string `json:"iiko_status"`
		CustomerPhoneNumber string `json:"customer_phone_number"`
		StoreID             string `json:"store_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err2.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	server.Logger.Errorf("received pos status: %s", req.IikoStatus)

	if err := server.orderKwaaka3plService.MapIikoStatusTo3plStatus(c.Request.Context(), req.IikoStatus, req.CustomerPhoneNumber, req.StoreID); err != nil {
		c.Set(errorKey, err)
		c.JSON(http.StatusInternalServerError, err2.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// Get3plDeliveryInfo docs
//
//	@Tags		kwaaka-admin
//	@Title		get delivery info for 3pl
//	@Security	ApiKeyAuth
//	@Param		delivery_ids	body	[]string	true	"delivery_ids"
//	@Success	200			{object}	models.GetDeliveryInfoResp
//	@Failure	401			{object}	[]errors.ErrorResponse
//	@Failure	400			{object}	[]errors.ErrorResponse
//	@Failure	500			{object}	[]errors.ErrorResponse
//	@Router		/kwaaka-admin/delivery-3pl-info [get]
func (server *Server) Get3plDeliveryInfo(c *gin.Context) {

	var req struct {
		DeliveryIDs []string `json:"delivery_ids"`
	}
	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err2.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := server.orderKwaaka3plService.ActualizeDeliveryInfoByDeliveryIDs(c.Request.Context(), req.DeliveryIDs); err != nil {
		c.Set(errorKey, err)
		c.JSON(http.StatusInternalServerError, err2.ErrorResponse{Msg: err.Error()})
		return
	}

	deliveriesInfo, err := server.orderKwaaka3plService.Get3plDeliveryInfo(c.Request.Context(), req.DeliveryIDs)
	if err != nil {
		server.Logger.Errorf("get 3pl delivery info error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err2.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, deliveriesInfo)
}

// GetOrdersByCustomerPhone
//
//	@Tags		kwaaka-admin
//	@Title		Method for get orders by customer phone
//	@Security	ApiKeyAuth
//	@Summary	Method Get orders by customer phone
//	@Param		phone			path		string	true	"phone"
//	@Param		restaurant_id	query		string	true	"restaurant_id"
//	@Param		page			query		string	true	"page"
//	@Param		limit			query		string	true	"limit"
//	@Success	200				{object}	models.GetOrdersByCustomerPhoneResponse
//	@Failure	401				{object}	errors.ErrorResponse
//	@Failure	400				{object}	errors.ErrorResponse
//	@Failure	500				{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/dispatcher/customer/phone/{phone}/orders [get]
func (s *Server) GetOrdersByCustomerPhone(c *gin.Context) {
	customerPhone := c.Param("phone")
	restaurantID, ok := c.GetQuery("restaurant_id")
	if !ok {
		err := fmt.Errorf("missing restaurant id")
		s.Logger.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, err2.ErrorResponse{Msg: err.Error()})
		return
	}

	page, limit, err := parsePaging(c)
	if err != nil {
		s.Logger.Errorf("page, limit parse error: %s", err)
		c.JSON(http.StatusBadRequest, err2.ErrorResponse{Msg: err.Error()})
		return
	}

	res, err := s.orderKwaaka3plService.GetOrdersByCustomerPhone(c.Request.Context(), models.GetOrdersByCustomerPhoneRequest{
		CustomerPhone: customerPhone,
		RestaurantID:  restaurantID,
		Pagination: models.Pagination{
			Page:  page,
			Limit: limit,
		}})
	if err != nil {
		s.Logger.Errorf("get orders by customer id error: %s", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *Server) GetCustomerByDeliveryId(c *gin.Context) {
	deliveryId := c.Param("delivery_id")

	res, storeId, err := srv.orderKwaaka3plService.GetCustomerByDeliveryId(c.Request.Context(), deliveryId)
	if err != nil {
		srv.Logger.Errorf("get order by delivery id error %s", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, struct {
		PhoneNumber string `json:"phone_number"`
		StoreId     string `json:"store_id"`
	}{
		PhoneNumber: res.PhoneNumber,
		StoreId:     storeId,
	})
}

func (srv *Server) GetOrderForTelegramByDeliveryOrderId(c *gin.Context) {

	deliveryId := c.Param("delivery_id")

	deliveryInfo, err := srv.orderKwaaka3plService.GetOrderForTelegramByDeliveryOrderId(c.Request.Context(), deliveryId)
	if err != nil {
		srv.Logger.Errorf("get order by delivery id error %s", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, deliveryInfo)
}

func (srv *Server) BulkCreate3plOrder(c *gin.Context) {

	orderID := c.Param("order_id")

	order, err := srv.orderKwaaka3plService.GetOrderByOrderID(c.Request.Context(), orderID)
	if err != nil {
		srv.Logger.Errorf("get order by order id error %s", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if order.DeliveryOrderID != "" {
		srv.Logger.Infof("save and set empty dispatcher for delivery order id: %s", order.DeliveryOrderID)
		if err := srv.orderKwaaka3plService.Save3plHistory(c.Request.Context(), order.DeliveryOrderID, models.DeliveryAddress{}, models.Customer{}); err != nil {
			srv.Logger.Errorf("save 3pl history error %s", err.Error())
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	err = srv.orderKwaaka3plService.BulkCreate3plOrder(c.Request.Context(), []models.Order{order}, true)
	if err != nil {
		srv.Logger.Errorf("bulk create 3pl error %s", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (srv *Server) InsertChangeDeliveryHistory(c *gin.Context) {

	deliveryID, ok := c.GetQuery("delivery_id")
	if !ok {
		c.JSON(http.StatusBadRequest, err2.ErrorResponse{Msg: ErrMissingDeliveryOrderID.Error()})
		return
	}

	var history models.ChangesHistory
	if err := c.BindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, err2.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := srv.orderKwaaka3plService.InsertChangeDeliveryHistory(c.Request.Context(), deliveryID, history); err != nil {
		c.JSON(http.StatusInternalServerError, err2.ErrorResponse{Msg: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
