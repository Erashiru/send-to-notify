package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	errResp "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/kwaaka-team/orders-core/service/promo_code/dto"
	"github.com/pkg/errors"
	"net/http"
)

// CreatePromoCode docs
//
//	@Tags		qr-menu/promo-code
//	@Title		Method for create promo code
//	@Security	ApiKeyAuth
//	@Summary	Method create PromoCode
//	@Param		request	body		models.PromoCode	true	"request"
//	@Success	200
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/promo-code [post]
func (server *Server) CreatePromoCode(c *gin.Context) {

	server.Logger.Info("start server to create promo code")

	var request models.PromoCode
	if err := c.BindJSON(&request); err != nil {
		server.Logger.Infof("bind error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("bind error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errResp.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	if err := server.PromoCode.Create(c.Request.Context(), request); err != nil {
		server.Logger.Infof("create promo code error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("create promo code error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errResp.ErrorResponse{
			Code:        http.StatusInternalServerError,
			Description: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

// UpdatePromoCode docs
//
//	@Tags		qr-menu/promo-code
//	@Title		Method for update promo code
//	@Security	ApiKeyAuth
//	@Summary	Method update PromoCode
//	@Param		request	body		models.UpdatePromoCode	true	"request"
//	@Success	200
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/promo-code [put]
func (server *Server) UpdatePromoCode(c *gin.Context) {

	server.Logger.Info("start server to update promo code")

	var request models.UpdatePromoCode
	if err := c.BindJSON(&request); err != nil {
		server.Logger.Infof("bind error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("bind error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, errResp.ErrorResponse{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	err := server.PromoCode.Update(c.Request.Context(), request)
	if err != nil {
		server.Logger.Infof("update promo code error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("update promo code error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errResp.ErrorResponse{
			Code:        http.StatusInternalServerError,
			Description: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

// GetPromoCodeByID docs
//
//	@Tags		qr-menu/promo-code
//	@Title		Method for get promo code by id
//	@Security	ApiKeyAuth
//	@Summary	Method get PromoCodeById
//	@Param		promo-code-id	path		string	true	"promo-code-id"
//	@Success	200		{object}	models.PromoCode
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/promo-code/id/:promo-code-id [get]
func (server *Server) GetPromoCodeByID(c *gin.Context) {

	server.Logger.Infof("server start to get promo code by id")

	promoCodeID := c.Param("promo-code-id")

	promoCode, err := server.PromoCode.GetByID(c.Request.Context(), promoCodeID)
	if err != nil {
		server.Logger.Infof("get promo code by id: %v error: %s", promoCodeID, err.Error())
		c.Set(errorKey, fmt.Sprintf("get promo code by id: %v error: %s", promoCodeID, err.Error()))
		switch {
		case errors.Is(err, dto.ErrPromoCodeNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, errResp.ErrorResponse{Msg: err.Error()})
		case errors.Is(err, dto.ErrInvalidPromoCodeID):
			c.AbortWithStatusJSON(http.StatusBadRequest, errResp.ErrorResponse{Msg: err.Error()})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, errResp.ErrorResponse{Msg: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, promoCode)
}

// GetAvailablePromoCodeByCode docs
//
//	@Tags		qr-menu/promo-code
//	@Title		Method for get available promo code by code
//	@Security	ApiKeyAuth
//	@Summary	Method get GetAvailablePromoCodeByCode
//	@Param		promo-code	path		string	true	"promo-code"
//	@Success	200		{object}	models.PromoCode
//	@Failure	404		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/promo-code/code/:promo-code [get]
func (server *Server) GetAvailablePromoCodeByCode(c *gin.Context) {

	server.Logger.Infof("start to get available promo code by code")

	promoCodeValue := c.Param("promo-code")

	promoCode, err := server.PromoCode.GetAvailablePromoCodeByCode(c.Request.Context(), promoCodeValue)
	if err != nil {
		server.Logger.Infof("get available promo code by code: %v error: %s", promoCodeValue, err.Error())
		c.Set(errorKey, fmt.Sprintf("get available promo code by code: %v error: %s", promoCodeValue, err.Error()))
		switch {
		case errors.Is(err, dto.ErrPromoCodeNotFound):
			c.AbortWithStatusJSON(http.StatusNotFound, errResp.ErrorResponse{Msg: err.Error()})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, errResp.ErrorResponse{Msg: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, promoCode)
}

// GetPromoCodesByRestaurantID docs
//
//	@Tags		qr-menu/promo-code
//	@Title		Method for get promo codes by restaurant id
//	@Security	ApiKeyAuth
//	@Summary	Method get PromoCodesByRestaurantID
//	@Param		restaurant-id	path		string	true	"restaurant-id"
//	@Param		page			query		string	true	"page"
//	@Param		limit			query		string	true	"limit"
//	@Success	200		{object}	[]models.PromoCode
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/promo-code/restaurant/:restaurant-id [get]
func (server *Server) GetPromoCodesByRestaurantID(c *gin.Context) {

	server.Logger.Infof("server start to get all promo codes for restaurant")

	restaurantID := c.Param("restaurant-id")

	page, limit, err := parsePaging(c)
	if err != nil {
		server.Logger.Infof("get promo code for restaurant id: %v  parse page error: %s", restaurantID, err.Error())
		c.Set(errorKey, fmt.Sprintf("get promo code for restaurant id: %v  parse page error: %s", restaurantID, err.Error()))
		return
	}
	pagination := selector.Pagination{
		Page:  page - 1,
		Limit: limit,
	}

	promoCodes, err := server.PromoCode.GetPromoCodesByRestaurantId(c.Request.Context(), restaurantID, pagination)
	if err != nil {
		server.Logger.Infof("get promo code for restaurant id: %v error: %s", restaurantID, err.Error())
		c.Set(errorKey, fmt.Sprintf("get promo code for restaurant id: %v error: %s", restaurantID, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errResp.ErrorResponse{
			Code:        http.StatusInternalServerError,
			Description: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, promoCodes)
}

// ValidatePromoCodeForUser docs
//
//	@Tags		qr-menu/promo-code
//	@Title		Method for validate promo code for user
//	@Security	ApiKeyAuth
//	@Summary	Method ValidatePromoCodeForUser
//	@Param		request	body	models.ValidateUserPromoCode	true	"request"
//	@Success	200		{object}	models.ValidateUserPromoCodeResponse
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/qr-menu/promo-code/validate-promo-code [post]
func (server *Server) ValidatePromoCodeForUser(c *gin.Context) {

	var userPromoCode models.ValidateUserPromoCode
	if err := c.BindJSON(&userPromoCode); err != nil {
		server.Logger.Infof("bind error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("bind error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, errResp.ErrorResponse{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	exist, comment, totalPrice, salePrice, products, err := server.PromoCode.ValidatePromoCodeForUser(c.Request.Context(), userPromoCode)
	if err != nil {
		server.Logger.Infof("validate promo code for user error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("validate promo code for user error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errResp.ErrorResponse{
			Code:        http.StatusInternalServerError,
			Description: err.Error(),
		})
		return
	}

	result := models.ValidateUserPromoCodeResponse{
		Exist:      exist,
		Comment:    comment,
		TotalPrice: totalPrice,
		SalePrice:  salePrice,
		Products:   products,
	}

	c.JSON(http.StatusOK, result)
}

// GetPromoCodeByCodeAndRestaurantId docs
//
//	@Tags		qr-menu/promo-code
//	@Title		Method for get promo code by code and restaurant id
//	@Security	ApiKeyAuth
//	@Summary	Method get GetPromoCodeByCodeAndRestaurantId
//	@Param		restaurant-id	path		string	true	"restaurant-id"
//	@Param		promo-code		path		string	true	"promo-code"
//	@Success	200		{object}	models.PromoCode
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/promo-code/:promo-code/restaurant/:restaurant-id [get]
func (server *Server) GetPromoCodeByCodeAndRestaurantId(c *gin.Context) {

	server.Logger.Infof("server start to get promo code by promo code value and restaurant id")

	promoCodeValue := c.Param("promo-code")
	restaurantID := c.Param("restaurant-id")

	promoCode, err := server.PromoCode.GetPromoCodeByCodeAndRestaurantId(c.Request.Context(), promoCodeValue, restaurantID)
	if err != nil {
		server.Logger.Infof("get promo code: %v for restaurant id: %v error: %s", promoCode, restaurantID, err.Error())
		c.Set(errorKey, fmt.Sprintf("get promo code: %v for restaurant id: %v error: %s", promoCode, restaurantID, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, errResp.ErrorResponse{
			Code:        http.StatusInternalServerError,
			Description: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, promoCode)
}

// AddUserPromoCodeUsageTimeToDB docs
//
//	@Tags		qr-menu/promo-code
//	@Title		Method for add user promo code usage time to database
//	@Security	ApiKeyAuth
//	@Summary	Method get AddUserPromoCodeUsageTimeToDB
//	@Param		request	body	models.UserPromoCodeUsageTimeRequest	true	"request"
//	@Success	200
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/qr-menu/promo-code/add-usage-time [post]
func (server *Server) AddUserPromoCodeUsageTimeToDB(c *gin.Context) {

	server.Logger.Infof("server start to add user promo code usage time to database")

	var request models.UserPromoCodeUsageTimeRequest
	if err := c.BindJSON(&request); err != nil {
		server.Logger.Infof("bind error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("bind error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, errResp.ErrorResponse{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		})
		return
	}

	err := server.PromoCode.AddUserPromoCodeUsageTimeToDB(c.Request.Context(), request.UserId, request.PromoCodeValue, request.RestaurantId)
	if err != nil {
		server.Logger.Infof("add user promo code usage time to db error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("add user promo code usage time to db error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, []errResp.ErrorResponse{{
			Code:        http.StatusInternalServerError,
			Description: err.Error(),
		}})
		return
	}

	c.Status(http.StatusOK)
}
