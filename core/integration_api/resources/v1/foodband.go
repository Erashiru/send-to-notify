package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/foodband/models"
	"github.com/kwaaka-team/orders-core/core/foodband/resources/http/v1/dto"
	"github.com/kwaaka-team/orders-core/domain/logger"
	"net/http"
)

func (server *Server) Auth(c *gin.Context) {
	apiToken := c.GetHeader(Authorization)
	storeId := c.Param("store_id")
	if apiToken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if storeId == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{
			Details: "store id is empty",
		})
		return
	}

	_, restaurantId, err := server.foodBandStoreManager.GetApiTokenStores(c.Request.Context(), apiToken, storeId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		return
	}

	if restaurantId != "" {
		c.Set("restaurant_id", restaurantId)
		return
	}

	c.AbortWithStatus(http.StatusUnauthorized)
}

// GetStoresFoodBand docs
//	@Tags		foodband
//	@Title		Method for getting available stores
//	@Param		Authorization	header		string	true	"token"
//	@Success	200				{object}	[]dto.Store
//	@Failure	400				{object}	dto.ErrorResponse
//	@Failure	500				{object}	dto.ErrorResponse
//	@Router		/v1/stores [get]
func (server *Server) GetStoresFoodBand(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "foodband request",
		Request: *c.Request,
	})

	apiToken := c.GetHeader(Authorization)

	stores, _, err := server.foodBandStoreManager.GetApiTokenStores(c.Request.Context(), apiToken, "")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stores)
}

// UploadMenuFoodBand docs
//	@Tags		foodband
//	@Title		Method for uploading menu
//	@Param		Authorization		header		string					true	"token"
//	@Param		menu				body		dto.MenuUploadRequest	true	"menu"
//	@Param		store_id			path		string					true	"store_id"
//	@Param		delivery_service	path		string					true	"delivery_service"
//	@Success	200					{object}	dto.MenuUploadResponse
//	@Failure	400					{object}	dto.ErrorResponse
//	@Failure	500					{object}	dto.ErrorResponse
//	@Router		/v1/stores/{store_id}/menus/{delivery_service} [post]
func (server *Server) UploadMenuFoodBand(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "foodband request",
		Request: *c.Request,
	})

	storeId := c.Param("store_id")
	deliveryService := c.Param("delivery_service")

	var req dto.MenuUploadRequest

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		server.Logger.Error(logger.LoggerInfo{
			System: "foodband response error",
			Response: dto.ErrorResponse{
				Details: err.Error(),
			},
		})
		return
	}

	transactionID, err := server.foodBandMenuManager.UploadMenu(c.Request.Context(), models.UploadMenuReq{
		StoreID:         storeId,
		DeliveryService: deliveryService,
		MenuURL:         req.MenuURL,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.MenuUploadResponse{
		TransactionID: transactionID,
	})
}

// ManageAggregatorStoreFoodBand docs
//	@Tags		foodband
//	@Title		Method for managing store status in aggregator
//	@Param		Authorization		header	string								true	"token"
//	@Param		is_open				body	dto.ManageAggregatorStoreRequest	true	"is_open"
//	@Param		store_id			path	string								true	"store_id"
//	@Param		delivery_service	path	string								true	"delivery_service"
//	@Success	200
//	@Failure	400	{object}	dto.ErrorResponse
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/v1/stores/{store_id}/manage/{delivery_service} [post]
func (server *Server) ManageAggregatorStoreFoodBand(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "foodband request",
		Request: *c.Request,
	})

	storeId := c.Param("store_id")
	deliveryService := c.Param("delivery_service")

	var req dto.ManageAggregatorStoreRequest

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		server.Logger.Error(logger.LoggerInfo{
			System: "foodband response error",
			Response: dto.ErrorResponse{
				Details: err.Error(),
			},
		})
		return
	}

	err := server.foodBandStoreManager.ManageStoreInAggregator(c.Request.Context(), models.ManageAggregatorStoreRequest{
		DeliveryService:       deliveryService,
		IsOpen:                *req.IsOpen,
		PosIntegrationStoreID: storeId,
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

// GetMenuUploadStatusFoodBand docs
//	@Tags		foodband
//	@Title		Method for getting menu upload status
//	@Param		Authorization		header		string	true	"token"
//	@Param		store_id			path		string	true	"store_id"
//	@Param		delivery_service	path		string	true	"delivery_service"
//	@Param		transaction_id		path		string	true	"transaction_id"
//	@Success	200					{object}	dto.MenuUploadStatusResponse
//	@Failure	404					{object}	dto.ErrorResponse
//	@Failure	401
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/v1/stores/{store_id}/menus/{delivery_service}/{transaction_id} [get]
func (server *Server) GetMenuUploadStatusFoodBand(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "foodband request",
		Request: *c.Request,
	})

	resp, err := server.foodBandMenuManager.GetMenuUploadStatus(c.Request.Context(), models.GetMenuUploadStatusReq{
		StoreID:         c.Param("store_id"),
		DeliveryService: c.Param("delivery_service"),
		TransactionID:   c.Param("transaction_id"),
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.FromMenuUploadTransaction(resp))
}

// UpdateOrderStatusFoodBand docs
//	@Tags		foodband
//	@Title		Method for updating order status
//	@Param		Authorization	header	string							true	"token"
//	@Param		status			body	dto.UpdateOrderStatusRequest	true	"status"
//	@Param		store_id		path	string							true	"store_id"
//	@Param		order_id		path	string							true	"order_id"
//	@Success	200
//	@Failure	401
//	@Failure	400	{object}	dto.ErrorResponse
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/v1/stores/{store_id}/orders/{order_id} [post]
func (server *Server) UpdateOrderStatusFoodBand(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "foodband request",
		Request: *c.Request,
	})

	var req dto.UpdateOrderStatusRequest
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		server.Logger.Error(logger.LoggerInfo{
			System: "foodband response error",
			Response: dto.ErrorResponse{
				Details: err.Error(),
			},
		})
		return
	}

	err := server.foodBandOrderManager.UpdateOrderStatus(c.Request.Context(), models.UpdateOrderStatusReq{
		Status:  req.Status,
		OrderID: c.Param("order_id"),
		StoreID: c.Param("store_id"),
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

// StopListProductFoodBand docs
//	@Tags		foodband
//	@Title		Method for updating products availability
//	@Param		Authorization	header		string							true	"token"
//	@Param		item			body		dto.UpdateItemStoplistRequest	true	"item"
//	@Param		store_id		path		string							true	"store_id"
//	@Param		product_id		path		string							true	"product_id"
//	@Success	200				{object}	dto.UpdateItemStoplistRespone
//	@Failure	400				{object}	dto.ErrorResponse
//	@Failure	401
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/v1/stores/{store_id}/products/{product_id} [post]
func (server *Server) StopListProductFoodBand(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "foodband request",
		Request: *c.Request,
	})

	productId := c.Param("product_id")
	restaurantId, ok := c.Get("restaurant_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: "invalid product id",
		})
		server.Logger.Error(logger.LoggerInfo{
			System: "foodband response error",
			Response: dto.ErrorResponse{
				Details: "invalid product id",
			},
		})
		return
	}

	var req dto.UpdateItemStoplistRequest

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		server.Logger.Error(logger.LoggerInfo{
			System: "foodband response error",
			Response: dto.ErrorResponse{
				Details: err.Error(),
			},
		})
		return
	}

	storeID := restaurantId.(string)
	if err := server.stopListService.UpdateStopListByPosProductID(c.Request.Context(), req.Available, storeID, productId); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateItemStoplistRespone{
		ID:        productId,
		Available: req.Available,
		Price:     req.Price,
	})
}

// StopListAttributeFoodBand docs
//	@Tags		foodband
//	@Title		Method for updating attributes availability
//	@Param		Authorization	header		string							true	"token"
//	@Param		item			body		dto.UpdateItemStoplistRequest	true	"item"
//	@Param		store_id		path		string							true	"store_id"
//	@Param		attribute_id	path		string							true	"attribute_id"
//	@Success	200				{object}	dto.UpdateItemStoplistRespone
//	@Failure	400				{object}	dto.ErrorResponse
//	@Failure	401
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/v1/stores/{store_id}/attributes/{attribute_id} [post]
func (server *Server) StopListAttributeFoodBand(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "foodband request",
		Request: *c.Request,
	})

	attributeId := c.Param("attribute_id")
	restaurantId, ok := c.Get("restaurant_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: "invalid product id",
		})
		server.Logger.Error(logger.LoggerInfo{
			System: "foodband response error",
			Response: dto.ErrorResponse{
				Details: "invalid product id",
			},
		})
		return
	}
	var req dto.UpdateItemStoplistRequest

	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		server.Logger.Error(logger.LoggerInfo{
			System: "foodband response error",
			Response: dto.ErrorResponse{
				Details: err.Error(),
			},
		})
		return
	}

	storeID := restaurantId.(string)
	if err := server.stopListService.UpdateStopListByAttributeID(c.Request.Context(), req.Available, storeID, attributeId); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateItemStoplistRespone{
		ID:        attributeId,
		Available: req.Available,
		Price:     req.Price,
	})
}
