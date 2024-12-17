package v1

import (
	"github.com/gin-gonic/gin"
	ordersCoreErr "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/service/payment/kaspi_salescout/dto"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (server *Server) secretMiddleware(secretToken string) func(c *gin.Context) {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if auth == "" {
			log.Error().Msgf("Authorization header is missing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, ordersCoreErr.ErrorResponse{
				Msg: "Authorization header is missing",
			})
			return
		}

		if auth != "Bearer "+secretToken {
			log.Error().Msgf("Authorization token not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, ordersCoreErr.ErrorResponse{
				Msg: "Authorization token not valid",
			})
			return
		}
	}
}

func (server *Server) CreateKaspiToken(c *gin.Context) {
	var request dto.CreatePaymentOrderRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	resp, err := server.SaleScoutProxeService.CreatePaymentToken(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (server *Server) CreateKaspiLink(c *gin.Context) {
	var request dto.CreatePaymentOrderRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	resp, err := server.SaleScoutProxeService.CreatePaymentLink(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (server *Server) GetKaspiStatusByID(c *gin.Context) {
	paymentID := c.Param("payment_id")
	if paymentID == "" {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: "paymentID query is missing"})
		return
	}

	resp, err := server.SaleScoutProxeService.GetPaymentStatusByID(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (server *Server) RefundKaspiPayment(c *gin.Context) {
	var request dto.RefundRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	resp, err := server.SaleScoutProxeService.RefundPayment(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusBadRequest, ordersCoreErr.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
