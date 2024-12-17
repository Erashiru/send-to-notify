package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/kwaaka_admin/models"
	"github.com/kwaaka-team/orders-core/core/managers/telegram"
	models2 "github.com/kwaaka-team/orders-core/core/menu/models"
	orderModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/service/iiko/resources/http/v1/detector"
	"github.com/kwaaka-team/orders-core/core/service/iiko/resources/http/v1/dto"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"net/http"
	"strconv"
)

func (s *Server) SendTelegramMessage(c *gin.Context) {
	var req models.SendTelegramMessageRequest
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		return
	}

	if req.Message == "" {
		c.Set(errorKey, "message cannot be empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Msg: "message cannot be empty",
		})
		return
	}
	if err := s.TelegramService.SendMessageToQueue(telegram.NotificationType(req.NotificationType), orderModels.Order{}, storeModels.Store{}, "", req.Message, "", models2.Product{}); err != nil {
		s.Logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (s *Server) Send3plErrorMsg(c *gin.Context) {
	var req models.SendTelegramMessageRequest
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		return
	}
	if req.Message == "" {
		c.Set(errorKey, "message cannot be empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Msg: "message cannot be empty",
		})
		return
	}
	order, err := s.orderInfoSharingService.GetOrderByID(c.Request.Context(), req.OrderID)
	if err != nil {
		c.Set(errorKey, "cannot find order for send tg message")
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Msg: "cannot find order for send tg message",
		})
		return
	}
	store, err := s.storeService.GetByID(c.Request.Context(), order.RestaurantID)
	if err != nil {
		c.Set(errorKey, "cannot find store for send tg message")
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Msg: "cannot find store for send tg message",
		})
		return
	}
	if err := s.TelegramService.SendMessageToQueue(telegram.ThirdPartyError, order, store, "", req.Message, "", models2.Product{}); err != nil {
		s.Logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (s *Server) SendCompensationMessage(c *gin.Context) {
	var req models.SendCompensationMessageRequest
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		return
	}

	order, err := s.orderInfoSharingService.GetOrder(c.Request.Context(), req.OrderID)
	if err != nil {
		s.Logger.Errorf("couldn't find order to send compensation notification: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	store, err := s.storeService.GetByID(c.Request.Context(), order.RestaurantID)
	if err != nil {
		s.Logger.Errorf("couldn't find store to send compensation notification: %s", err.Error())
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	if err = s.TelegramService.SendMessageToQueue(telegram.Compensation, order, store, req.CompensationID, strconv.Itoa(req.CompensationNumber), req.CompensationText, models2.Product{}); err != nil {
		s.Logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
