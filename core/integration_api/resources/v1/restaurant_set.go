package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/models"
	"net/http"
)

func (srv *Server) CreateRestaurantSet(c *gin.Context) {
	var req models.RestaurantSet
	if err := c.BindJSON(&req); err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	srv.Logger.Infof("getting request to create restaurant set with next params %v", req)

	id, err := srv.restaurantSetService.CreateRestaurantSet(c.Request.Context(), req)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, id)
}

func (srv *Server) GetRestaurantSetById(c *gin.Context) {
	id := c.Param("restaurant_set_id")

	srv.Logger.Infof("getting request to find restaurant_set with id %s", id)

	res, err := srv.restaurantSetService.GetRestaurantSetById(c.Request.Context(), id)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *Server) GetRestaurantSetWithRestGroup(c *gin.Context) {
	id := c.Param("restaurant_set_id")

	srv.Logger.Infof("getting request to find restaurant_set and rest_group with id %s", id)

	res, err := srv.restaurantSetService.GetRestaurantSetInfoWithRestGroup(c.Request.Context(), id)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (srv *Server) GetRestaurantSetByDomainName(c *gin.Context) {
	domainName := c.Query("domain_name")
	if domainName == "" {
		c.Set(errorKey, errGetQuery)
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: errGetQuery + " :missing query param"})
		return
	}

	srv.Logger.Infof("getting request to find restaurant_set and rest_group with domain name %s", domainName)

	res, err := srv.restaurantSetService.GetRestaurantSetByDomainName(c.Request.Context(), domainName)
	if err != nil {
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
