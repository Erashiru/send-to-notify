package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	"github.com/kwaaka-team/orders-core/core/models"
	"net/http"
	"strconv"
)

func (server *Server) SetMarkUpToAggregatorMenu(c *gin.Context) {
	var req dto.RequestSetMarkUpToAggregatorMenu

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf(errBindBody, err.Error()))
		return
	}

	store, err := server.storeService.GetByID(c.Request.Context(), req.StoreId)
	if err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if !store.VerifyMenuOwnership(req.MenuId) {
		server.Logger.Error("menu doesn't belong to the restaurant")
		c.Set(errorKey, "menu doesn't belong to the restaurant")
		c.AbortWithStatusJSON(http.StatusBadRequest, "menu doesn't belong to the restaurant")
		return
	}

	if store.Settings.PriceSource != models.POSPriceSource {
		server.Logger.Error("restaurant.settings.price_source field is not equal POS")
		c.Set(errorKey, "restaurant.settings.price_source field is not equal POS")
		c.AbortWithStatusJSON(http.StatusBadRequest, "restaurant.settings.price_source field is not equal POS")
		return
	}

	markupPercent := store.GetMenuMarkupPercent(req.MenuId)
	if markupPercent == 0 {
		server.Logger.Error("restaurant.menus.markup_percent is 0")
		c.Set(errorKey, "restaurant.menus.markup_percent is 0")
		c.AbortWithStatusJSON(http.StatusBadRequest, "restaurant.menus.markup_percent is 0")
		return
	}

	if err = server.menuService.SetMarkupToAggregatorMenu(c.Request.Context(), req.MenuId, store.MenuID, store.Settings.Currency, markupPercent); err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (server *Server) GetWoltMenuInCsv(c *gin.Context) {
	menuId, ok := c.GetQuery("menu_id")
	if !ok {
		server.Logger.Error("menu_id query is missing")
		c.Set(errorKey, "menu_id query is missing")
		c.AbortWithStatusJSON(http.StatusBadRequest, "menu_id query is missing")
		return
	}

	needUploadImage, ok := c.GetQuery("need_upload_image")
	if !ok {
		server.Logger.Error("need_upload_image query is missing")
		c.Set(errorKey, "need_upload_image query is missing")
		c.AbortWithStatusJSON(http.StatusBadRequest, "need_upload_image query is missing")
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		server.Logger.Error("file upload is missing or invalid")
		c.Set(errorKey, "file upload is missing or invalid")
		c.AbortWithStatusJSON(http.StatusBadRequest, "file upload is missing or invalid")
		return
	}
	defer file.Close()

	url, err := server.menuService.ConvertMenuToWoltCsv(c.Request.Context(), menuId, needUploadImage == "true", file)
	if err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.String(http.StatusOK, url)
}

func (server *Server) AutoUpdateAggregatorMenu(c *gin.Context) {
	var req dto.RequestAutoUpdateAggregatorMenu

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf(errBindBody, err.Error()))
		return
	}

	store, err := server.storeService.GetByID(c.Request.Context(), req.StoreID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	err = server.menuService.AutoUpdateAggregatorMenu(c.Request.Context(), store, req.AggregatorMenuId)
	if err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (server *Server) GenerateNewAggregatorMenu(c *gin.Context) {
	var req dto.RequestGenerateNewAggregatorMenu

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf(errBindBody, err.Error()))
		return
	}

	store, err := server.storeService.GetByID(c.Request.Context(), req.StoreID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	id, err := server.menuService.GenerateAggregatorMenuFromPosMenu(c.Request.Context(), store, req.AggregatorMenuId, req.Delivery)
	if err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = server.storeService.AddMenuObjectToMenus(c.Request.Context(), store.ID, req.Delivery, id); err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, id)
}

func (server *Server) GenerateAggregatorMenuFromPosMenu(c *gin.Context) {
	var req dto.RequestGenerateAggregatorMenuFromPos

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, fmt.Sprintf(errBindBody, err.Error()))
		return
	}

	store, err := server.storeService.GetByID(c.Request.Context(), req.StoreId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, nil)
		return
	}

	id, err := server.menuService.GenerateNewAggregatorMenuFromPosMenu(c.Request.Context(), store, req.Delivery)
	if err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if err = server.storeService.AddMenuObjectToMenus(c.Request.Context(), store.ID, req.Delivery, id); err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, id)
}

func (server *Server) CreateDelivery3plCron(c *gin.Context) {
	defaultTime := 10

	callTime := c.Query("call_time")
	timing, err := strconv.Atoi(callTime)
	if err != nil {
		timing = defaultTime
	}

	orders, err := server.orderCronService.Get3plOrdersWithoutDriver(c.Request.Context(), int64(timing))
	if err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if len(orders) > 0 {
		server.Logger.Infof("orders that need 3pl: %v", orders)
	}

	err = server.orderKwaaka3plService.BulkCreate3plOrder(c.Request.Context(), orders, false)
	if err != nil {
		server.Logger.Error(err)
		c.Set(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
