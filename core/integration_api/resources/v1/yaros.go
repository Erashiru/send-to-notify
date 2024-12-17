package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	yarosModels "github.com/kwaaka-team/orders-core/core/yaros/models"
	"github.com/kwaaka-team/orders-core/core/yaros/resources/http/v1/dto"
	"github.com/rs/zerolog/log"
	"net/http"
)

func YarosSecretMiddleware(secretToken string) func(c *gin.Context) {
	return func(c *gin.Context) {
		auth := c.GetHeader(Authorization)

		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorResponse{
				Msg: errors.ErrTokenIsNotValid.Error(),
			})
			return
		}

		if auth != secretToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorResponse{
				Msg: errors.ErrTokenIsNotValid.Error(),
			})
			return
		}

		c.Set(Token, auth)
	}
}

// OrderUpdateYaros docs
//	@Tags		yaros
//	@Summary	order update yaros
//	@Security	ApiKeyAuth
//	@Summary	order update yaros
//	@Param		order_update_request_body	body	yarosModels.OrderUpdateRequestBody	true	"order_update_request_body"
//	@Success	200
//	@Failure	400	{object}	dto.ErrorResponse
//	@Router		/v1/yaros/order-update [patch]
func (server *Server) OrderUpdateYaros(c *gin.Context) {

	var req yarosModels.OrderUpdateRequestBody

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}
	log.Info().Msgf("Yaros updating order status: %v", req.PosOrderID)

	if err := server.statusUpdateService.UpdateOrderStatus(c.Request.Context(), req.PosOrderID, req.Status, req.ErrorDescription); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Msg: err.Error(),
		})
		log.Err(err).Msgf("Yaros updating order status error: %v", req.PosOrderID)
		return
	}

	c.Status(http.StatusOK)
}

// StoplistUpdateYaros docs
//	@Tags		yaros
//	@Summary	stoplist update yaros
//	@Param		stoplist_update_request	body	yarosModels.StoplistUpdateRequest	true	"stoplist_update_request"
//	@Success	200
//	@Failure	400	{object}	dto.ErrorResponse
//	@Router		/v1/yaros/stoplist-update [patch]
func (server *Server) StoplistUpdateYaros(c *gin.Context) {
	var req yarosModels.StoplistUpdateRequest

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}
	if err := server.stopListService.ActualizeStoplistbyYarosStoreID(c.Request.Context(), req.StoreId); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
