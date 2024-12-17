package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/service/iiko/resources/http/v1/detector"
	"github.com/kwaaka-team/orders-core/core/service/iiko/resources/http/v1/dto"
	"github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func (server *Server) iikoProductsSecretMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		auth := c.GetHeader(Authorization)
		splitToken := strings.Split(auth, Bearer)

		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ErrorResponse{
				Msg: errors.ErrTokenIsNotValid.Error(),
			})
			return
		}

		c.Set(StoreTokenKey, strings.TrimSpace(splitToken[1]))
	}
}

// EventIIKO docs
//
//	@Tags		iiko
//	@Title		get events from iiko pos
//	@Accept		json
//	@Produce	json
//	@Param		req	body	models.WebhookEvent	true	"events"
//	@Success	200
//	@Failure	401	{object}	[]dto.ErrorResponse
//	@Failure	400	{object}	[]dto.ErrorResponse
//	@Failure	500	{object}	[]dto.ErrorResponse
//	@Router		/iiko/events [post]
func (server *Server) EventIIKO(c *gin.Context) {

	token := c.GetString(StoreTokenKey)

	if token == "" {
		server.Logger.Info(errAuthTokenIsNotValid)
		c.Set(errorKey, errAuthTokenIsNotValid)
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code: http.StatusUnauthorized,
		})
		return
	}

	var req models.WebhookEvents

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		return
	}

	server.Logger.Info(zap.Any("IIKO event request body", req))

	rsp, err := server.iikoManager.WebhookEvent(c.Request.Context(), token, req, models.IIKO)
	if err != nil {
		server.Logger.Infof("webhook event error: %s", err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(detector.ErrorHandler(err))
		return
	}

	server.Logger.Info(zap.Any("IIKO event response body", rsp))

	c.JSON(http.StatusOK, models.WebhookEventResponse{
		Details: rsp,
	})
}

func (srv *Server) GetCustomerDiscount(c *gin.Context) {
	storeId := c.Param("store_id")

	var req struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.BindJSON(&req); err != nil {
		srv.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		return
	}

	resp, err := srv.iikoManager.GetCustomerDiscounts(c.Request.Context(), storeId, req.PhoneNumber)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(detector.ErrorHandler(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (srv *Server) GetDiscountHistory(c *gin.Context) {
	storeId := c.Param("store_id")

	var req struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.BindJSON(&req); err != nil {
		srv.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:  http.StatusBadRequest,
			Error: err,
		})
		return
	}

	resp, err := srv.iikoManager.GetDiscountHistory(c.Request.Context(), storeId, req.PhoneNumber)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(detector.ErrorHandler(err))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (srv *Server) GetDiscountsForStore(c *gin.Context) {
	storeID := c.Param("store_id")

	res, err := srv.iikoManager.GetDiscountsForStore(c.Request.Context(), storeID)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(detector.ErrorHandler(err))
		return
	}
	c.JSON(http.StatusOK, res)
}
