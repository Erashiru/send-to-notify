package v1

import (
	goErrors "errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/service/order"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (server *Server) UpdateOrderStatusByPosTypes(c *gin.Context) {
	var req dto.UpdateOrderStatusCronRequest

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	for _, posType := range req.PosTypes {
		if err := server.orderCronService.ActualizeOrdersStatusByPosType(c.Request.Context(), posType); err != nil {
			log.Err(err).Msgf("manual update order status update by cron, pos_type=%s", posType)
			continue
		}

		log.Info().Msgf("success manual order status update by cron, pos_type=%s", posType)
	}

	c.Status(http.StatusNoContent)
}

func (server *Server) UpdateStopListByPosTypes(c *gin.Context) {
	var req dto.UpdateStopListCronRequest

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Info(errBindBody)
		c.Set(errorKey, fmt.Errorf("[UPDATE STOPLIST CRON] %w", err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if err := server.stopListService.ActualizeStopListByPosTypes(c.Request.Context(), req.PosTypes); err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (server *Server) UpdateStopListBySection(c *gin.Context) {
	var req dto.UpdateStopListBySectionCronRequest

	if err := c.Bind(&req); err != nil {
		server.Logger.Info(errBindBody)
		c.Set(errorKey, fmt.Errorf("[UPDATE STOPLIST CRON] %w", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if req.IsAvailable == nil {
		server.Logger.Info("is_available is empty")
		c.Set(errorKey, fmt.Errorf("[UPDATE STOPLIST CRON] is_available is empty"))
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: "is_available is empty",
		})
		return
	}

	deliveryToSectionIDs := map[string][]string{
		models2.WOLT.String():   req.WoltSectionIDs,
		models2.GLOVO.String():  req.GlovoSectionIDs,
		models2.YANDEX.String(): req.YandexSectionIDs,
	}

	log.Info().Msg("[UPDATE STOPLIST CRON]")

	if err := server.stopListService.UpdateStopListBySectionID(c, *req.IsAvailable, req.RestaurantGroupID, deliveryToSectionIDs); err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (server *Server) GetOrderStat(c *gin.Context) {
	interval, ok := c.GetQuery("interval")
	if !ok {
		server.Logger.Error(fmt.Errorf("missing interval param"))
		c.Set(errorKey, fmt.Errorf("missing interval param"))
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Errorf("missing interval param"))
		return
	}

	log.Info().Msg("[ORDER STAT CRON START]")

	result, err := server.orderCronService.GetOrderStat(c.Request.Context(), interval)
	if err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		status := http.StatusInternalServerError
		if goErrors.Is(err, order.ErrInvalidInterval) {
			status = http.StatusBadRequest
		}
		c.AbortWithStatusJSON(status, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (server *Server) SendDeferOrders(c *gin.Context) {

	if err := server.orderCronService.SendDeferStatusSubmission(c.Request.Context()); err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (server *Server) NoDispatcherMessage(c *gin.Context) {
	if err := server.orderKwaaka3plService.NoDispatcherMessage(c.Request.Context()); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (server *Server) PerformerLookupMoreThan15Minute(c *gin.Context) {

	if err := server.orderKwaaka3plService.PerformerLookupMoreThan15Minute(c.Request.Context()); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
