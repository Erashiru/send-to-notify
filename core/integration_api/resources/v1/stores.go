package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/kwaaka_admin/models"
	"github.com/kwaaka-team/orders-core/core/service/iiko/resources/http/v1/detector"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"net/http"
	"strconv"
)

func (srv *Server) GetRestaurantsByGroupId(c *gin.Context) {
	groupId := c.Param("restaurant_group_id")

	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		srv.Logger.Info(err.Error())
		srv.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		srv.Logger.Info(err.Error())
		srv.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	pagination := selector.Pagination{
		Page:  page,
		Limit: limit,
	}

	stores, err := srv.storeService.GetRestaurantsByGroupId(c, pagination, groupId)
	if err != nil {
		srv.Logger.Info(err.Error())
		srv.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stores)
}

// CreatePolygon docs
//
//	//@Tags		kwaaka-admin
//	@Title		Method for create polygon  restaurant
//	@Security	ApiKeyAuth
//	@Summary	Method for create restaurant polygon
//	@Param		create	polygon	body	models.PolygonRequest	true	"polygon"
//	@Success	200
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/polygon [post]
func (srv *Server) CreatePolygon(c *gin.Context) {

	var request models.PolygonRequest
	if err := c.BindJSON(&request); err != nil {
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	err := srv.storeService.CreatePolygon(c.Request.Context(), request)
	if err != nil {
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	c.Status(http.StatusOK)
}

// UpdatePolygon docs
//
//	//@Tags		kwaaka-admin
//	@Title		Method for update polygon  restaurant
//	@Security	ApiKeyAuth
//	@Summary	Method for update polygonrestaurant
//	@Param		update	polygon	body	models.PolygonRequest	true	"polygon"
//	@Success	200
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/polygon [put]
func (srv *Server) UpdatePolygon(c *gin.Context) {

	var request models.PolygonRequest
	if err := c.BindJSON(&request); err != nil {
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	if err := srv.storeService.UpdatePolygon(c.Request.Context(), request); err != nil {
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	c.Status(http.StatusOK)
}

// GetPolygonByRestaurantID docs
//
//	//@Tags		kwaaka-admin
//	@Title		Method for get polygon by id for restaurant
//	@Security	ApiKeyAuth
//	@Summary	Method for get polygon by id for restaurant
//	@Param		restaurant_id	true	"request"
//	@Success	200{object}		models.GetPolygonResponse
//	@Failure	401				{object}	errors.ErrorResponse
//	@Failure	400				{object}	errors.ErrorResponse
//	@Failure	500				{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/polygon [get]
func (srv *Server) GetPolygonByRestaurantID(c *gin.Context) {
	restaurantID := c.Param("restaurant_id")
	if restaurantID == "" {
		c.JSON(http.StatusBadRequest, "restaurant id is empty")
	}

	polygon, err := srv.storeService.GetPolygonByRestaurantID(c.Request.Context(), restaurantID)
	if err != nil {
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	c.JSON(http.StatusOK, polygon)
}

// CreateStorePhoneEmail docs
//
//	@Tags		kwaaka-admin
//	@Title		Method for create phone and email for restaurant
//	@Security	ApiKeyAuth
//	@Summary	Method for create phone and email
//	@Param		request			models.StorePhoneEmail	true	"request"
//	@Param		restaurant-id	param					true	"restaurant-id"
//	@Success	200
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/{restaurant-id}/customer-data [post]
func (srv *Server) CreateStorePhoneEmail(c *gin.Context) {

	var request models.StorePhoneEmail
	if err := c.BindJSON(&request); err != nil {
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	restaurantID := c.Param("restaurant-id")

	if err := srv.storeService.CreateStorePhoneEmail(c.Request.Context(), restaurantID, request); err != nil {
		c.Set(errorKey, err)
		c.JSON(detector.ErrorHandler(err))
		return
	}

	c.Status(http.StatusOK)
}
