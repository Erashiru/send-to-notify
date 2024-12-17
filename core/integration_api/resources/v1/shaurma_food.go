package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/models"
	"net/http"
)

// GetOrdersShaurmaFood
//
//	@Tags		shaurma-food
//	@Title		Method for get orders for shaurma food
//	@Security	ApiKeyAuth
//	@Summary	Method Get orders for shaurma food
//	@Param		delivery_service	query		string	true	"delivery_service"
//	@Param		restaurant_id		query		string	true	"restaurant_id"
//	@Param		page				query		string	true	"page"
//	@Param		limit				query		string	true	"limit"
//	@Param		status				query		string	true	"status"
//	@Param		date_from			query		string	true	"date_from"
//	@Param		date_to				query		string	true	"date_to"
//	@Success	200					{object}	models.ShaurmaFoodOrdersInfoResponse
//	@Failure	401					{object}	errors.ErrorResponse
//	@Failure	400					{object}	errors.ErrorResponse
//	@Failure	500					{object}	errors.ErrorResponse
//	@Router		/v1/shaurma-food/orders [get]
func (s *Server) GetOrdersShaurmaFood(c *gin.Context) {
	restaurantID, _ := c.GetQuery("restaurant_id")

	page, limit, err := parsePaging(c)
	if err != nil {
		s.Logger.Errorf("page, limit parse error: %s", err)
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}
	deliveryService, _ := c.GetQuery("delivery_service")
	status, _ := c.GetQuery("status")
	dateFrom, _ := c.GetQuery("date_from")
	dateTo, _ := c.GetQuery("date_to")

	orders, total, err := s.orderInfoSharingService.GetShaurmaFoodOrders(c.Request.Context(), restaurantID, deliveryService, status, dateFrom, dateTo, page, limit)
	if err != nil {
		s.Logger.Errorf("get orders by customer id error: %s", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, models.ShaurmaFoodOrdersInfoResponse{
		Orders:          orders,
		PagesTotalCount: total,
	})
}

func (s *Server) SetExternalMenuPricesToAggregatorMenus(c *gin.Context) {
	type request struct {
		RestaurantID                                string `json:"restaurant_id"`
		ApiKey                                      string `json:"api_key"`
		OrganizationID                              string `json:"organization_id"`
		TerminalID                                  string `json:"terminal_id"`
		ExternalMenuID                              string `json:"external_menu_id"`
		PriceCategory                               string `json:"price_category_id"`
		IgnoreExternalMenuProductsWithZeroNullPrice bool   `json:"ignore_external_menu_products_with_zero_null_price"`
	}

	var req request

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if req.RestaurantID == "" {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: "restaurant id is empty"})
		return
	}

	if err := s.shaurmaFoodService.SetAggregatorsMenuPricesFromExternalMenu(c.Request.Context(), req.RestaurantID, req.ApiKey, req.OrganizationID, req.TerminalID, req.ExternalMenuID, req.PriceCategory, req.IgnoreExternalMenuProductsWithZeroNullPrice); err != nil {
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) UpdateImagesAndDescriptionsInAggregatorMenus(c *gin.Context) {
	restaurantID := c.Param("restaurant_id")
	aggregator := c.Param("aggregator")

	if err := s.shaurmaFoodService.UpdateProductsFromMainMenu(c.Request.Context(), restaurantID, aggregator); err != nil {
		c.Set(errorKey, err)
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
