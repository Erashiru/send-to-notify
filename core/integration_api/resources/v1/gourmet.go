package v1

import (
	"github.com/gin-gonic/gin"
	ordersCoreErr "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/models/gourmet"
	"net/http"
)

func (s *Server) GourmetGetTables(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	if restaurantId == "" {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: "restaurantId is missing"})
		return
	}

	resp, err := s.gourmetService.GetRestaurantTables(c.Request.Context(), restaurantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) GourmetGetOrders(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	if restaurantId == "" {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: "restaurantId is missing"})
		return
	}

	tableID := c.Query("tableId")
	orderID := c.Query("orderId")

	if tableID == "" && orderID == "" {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: "tableId or orderId should be filled"})
		return
	}

	if tableID != "" && orderID != "" {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: "only one of tableId or orderId should be filled"})
		return
	}

	resp, err := s.gourmetService.GetRestaurantOrders(c.Request.Context(), restaurantId, tableID, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (s *Server) GourmetPay(c *gin.Context) {
	restaurantId := c.Param("restaurantId")
	if restaurantId == "" {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: "restaurantId is missing"})
		return
	}

	orderID := c.Param("orderId")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: "orderId should be filled"})
		return
	}

	var req gourmet.PaymentChangeRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	resp, err := s.gourmetService.CreatePayment(c.Request.Context(), restaurantId, orderID, req.PaymentTypeId, req.PaymentTypeKind, req.IsPaid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
