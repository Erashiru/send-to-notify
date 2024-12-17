package v1

import (
	goErr "errors"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	paymentDto "github.com/kwaaka-team/orders-core/service/payment/whatsapp/dto"
	wppService "github.com/kwaaka-team/orders-core/service/whatsapp"
	"net/http"
)

func (srv *Server) SendNewsletter(c *gin.Context) {
	var req dto.SendNewsletterRequest
	if err := c.BindJSON(&req); err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := srv.WhatsappService.SendNewsletter(c.Request.Context(), req.RestGroupId, req.Text, req.Name); err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (srv *Server) SendMessage(c *gin.Context) {
	var req dto.SendMessage
	if err := c.BindJSON(&req); err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	srv.Logger.Infof("whatsapp will send message using next settings: store_id %s, message %s, phone_number %s", req.StoreId, req.Message, req.Phone)

	if err := srv.WhatsappService.SendMessage(c.Request.Context(), req.Phone, req.Message, req.StoreId); err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (server *Server) WhatsappWebhooks(c *gin.Context) {
	var req paymentDto.WebhookEvent
	if err := c.BindJSON(&req); err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	server.Logger.Infof("whatsapp webhook event received:\n Event type: %s, Pushname: %s, Body %s, quotedMsg: %s, time %s",
		req.EventType, req.Data.Pushname, req.Data.Body, req.Data.QuoteMsg, req.Data.Time)
	if _, err := server.paymentManager.WebhookEvent(c.Request.Context(), models.WHATSAPP, req); err != nil {
		server.Logger.Infof("whatsapp %s webhook event error: %s", req.EventType, err.Error())
		c.Set(errorKey, err)
		switch {
		case goErr.Is(err, wppService.ErrWpp):
			c.Status(http.StatusOK)
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
				Msg: err.Error(),
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
