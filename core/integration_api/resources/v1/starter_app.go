package v1

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/starter_app/models"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
)

func (s *Server) StarterAppAuth(c *gin.Context) {
	apiToken := c.GetHeader(Authorization)
	if apiToken == "" {
		s.Logger.Error(errAuthIsMissing)
		c.Set(errorKey, errors.ErrTokenIsNotValid)
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorResponse{
			Msg: errors.ErrTokenIsNotValid.Error(),
		})
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		s.Logger.Error("Error reading request body", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	type request struct {
		ShopID int `json:"shopId"`
	}

	var req request

	if err := json.Unmarshal(body, &req); err != nil {
		s.Logger.Error(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	starterStoreID := strconv.Itoa(req.ShopID)

	store, err := s.storeService.GetByExternalIdAndDeliveryService(c.Request.Context(), starterStoreID, models.STARTERAPP.String())
	if err != nil {
		s.Logger.Error("starter app middleware error: %s", err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if apiToken != store.StarterApp.ApiKey {
		s.Logger.Error(errBindBody, errors.ErrTokenIsNotValid)
		c.Set(errorKey, errors.ErrTokenIsNotValid)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: errors.ErrTokenIsNotValid.Error()})
		return
	}
}

func (s *Server) CreateOrderStarterApp(c *gin.Context) {

	log.Info().Msgf("create order starter app")

	var body []byte

	if s.withRedirect() {
		body = s.readBodyAndSetAgain(c)
	}

	var orderDto models.OrderDto

	if err := c.BindJSON(&orderDto); err != nil {
		s.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	order := orderDto.FromDto()

	log.Info().Msgf("create order starter app bind json ok")

	if s.withRedirect() && !s.isTopPartner(c.Request.Context(), order.ShopId, models.STARTERAPP.String()) {
		s.redirectRequest(c, body)
		return
	}

	log.Info().Msgf("create order starter no redirect")

	switch models.StarterAppStatus(order.Status) {
	case models.Checked:
		res, err := s.orderService.CreateOrder(c.Request.Context(), order.ShopId, models.STARTERAPP.String(), order, "")
		if err != nil {
			s.Logger.Infof("create order error: %s", err.Error())
			c.Set(errorKey, err)
			c.JSON(http.StatusOK, "")
			return
		}

		c.JSON(http.StatusOK, res.OrderID)
		return

	case models.Canceled:
		if err := s.starterAppOrderManager.CancelOrder(c.Request.Context(), order); err != nil {
			s.Logger.Infof("cancel order error: %s", err.Error())
			c.Set(errorKey, err)
			c.JSON(http.StatusOK, "")
		}
	}

	c.JSON(http.StatusOK, "")
}
