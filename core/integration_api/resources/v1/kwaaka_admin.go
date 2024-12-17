package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	integrationApiModels "github.com/kwaaka-team/orders-core/core/integration_api/resources/v1/dto"
	"github.com/kwaaka-team/orders-core/core/kwaaka_admin/models"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/domain/logger"
	"net/http"
)

func (server *Server) validPosForUpsertMenu(posType string) bool {
	valid := map[string]struct{}{
		"ytimes":   {},
		"tillypad": {},
		"posist":   {},
	}

	if _, ok := valid[posType]; ok {
		return true
	}

	return false
}

func (server *Server) UpsertMenu(c *gin.Context) {
	var req integrationApiModels.UpsertMenuRequest

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Infof(errBindBody, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	store, err := server.storeService.GetByID(c.Request.Context(), req.StoreId)
	if err != nil {
		server.Logger.Infof("restaurant with id %s not found: %s", req.StoreId, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if !server.validPosForUpsertMenu(store.PosType) {
		server.Logger.Infof("upsert menu for %s pos system not allowed", store.PosType)
		c.Set(errorKey, fmt.Errorf("upsert menu for %s pos system not allowed", store.PosType))
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: fmt.Sprintf("upsert menu for %s pos system not allowed", store.PosType),
		})
		return
	}

	posService, err := server.posService.GetPosService(coreModels.Pos(store.PosType), store)
	if err != nil {
		server.Logger.Infof("get %s pos service: %s", store.PosType, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	var systemMenu coreMenuModels.Menu

	if store.MenuID != "" {
		menu, err := server.menuService.FindById(c.Request.Context(), store.MenuID)
		if err != nil {
			server.Logger.Infof("get system menu with id %s: %s", store.MenuID, err.Error())
			c.Set(errorKey, err)
			c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
				Msg: err.Error(),
			})
			return
		}

		systemMenu = *menu
	}

	posMenu, err := posService.GetMenu(c.Request.Context(), store, systemMenu)
	if err != nil {
		server.Logger.Infof("get menu from pos: %s", err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	id, err := server.menuService.UpsertMenu(c.Request.Context(), store, systemMenu, posMenu)
	if err != nil {
		server.Logger.Infof("upsert menu error: %s", err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.String(http.StatusOK, id)
}

// CreateOrderKwaakaAdmin docs
//
//	@Tags		kwaaka-admin
//	@Title		Method for create order
//	@Security	ApiKeyAuth
//	@Summary	Method for create order
//	@Param		create	order		body	models.Order	true	"order"
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/kwaaka-admin/placeOrder [post]
func (server *Server) CreateOrderKwaakaAdmin(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "kwaaka_admin create order request",
		Request: c.Request,
	})

	var req models.Order

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	server.Logger.Infof("starting create kwaaka_admin order: ID: %s, rest id: %s, customer name: %s, payment type: %s, operator id: %s, operator name: %s, dispatcher: %s, client delivery price: %f, full delivery price: %f, kwaaka charged delivery price: %f",
		req.ID, req.RestaurantID, req.Customer.Name, req.PaymentType, req.OperatorID, req.OperatorName, req.Delivery.Dispatcher, req.Delivery.ClientDeliveryPrice, req.Delivery.FullDeliveryPrice, req.Delivery.KwaakaChargedDeliveryPrice)

	res, err := server.orderService.CreateOrder(c.Request.Context(), req.RestaurantID, "kwaaka_admin", req, "")
	if err != nil {
		server.Logger.Errorf("create order error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if res.IsInstantDelivery && res.DeliveryDispatcher != "" {
		if err = server.orderKwaaka3plService.Instant3plOrder(c.Request.Context(), res); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
				Msg: fmt.Sprintf("order created but Instant 3pl did not: %s", err.Error()),
			})
		}
	}

	c.JSON(http.StatusOK, res.ID)
}

// KwaakaAdminStopListByProductID docs
//
//	@Tags		kwaaka-admin
//	@Title		Method for update stoplist by product
//	@Security	ApiKeyAuth
//	@Summary	Method for update stoplist by product
//	@Param		stoplist	body		models.StopListByProductIDRequest	true	"stoplist"
//	@Failure	401			{object}	errors.ErrorResponse
//	@Failure	400			{object}	errors.ErrorResponse
//	@Failure	500			{object}	errors.ErrorResponse
//	@Router		/kwaaka-admin/stoplist/product [post]
func (server *Server) KwaakaAdminStopListByProductID(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "kwaaka_admin update stoplist by product id request",
		Request: c.Request,
	})

	var req models.StopListByProductIDRequest

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if err := server.stopListService.UpdateStopListByPosProductID(c.Request.Context(), req.IsAvailabe, req.StoreID, req.ProductID); err != nil {
		server.Logger.Errorf("update stoplist by product id error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

// KwaakaAdminStopListByAttributeID docs
//
//	@Tags		kwaaka-admin
//	@Title		Method for update stoplist by attribute
//	@Security	ApiKeyAuth
//	@Summary	Method for update stoplist by attribute
//	@Param		stoplist	body		models.StopListByAttributeIDRequest	true	"stoplist"
//	@Failure	401			{object}	errors.ErrorResponse
//	@Failure	400			{object}	errors.ErrorResponse
//	@Failure	500			{object}	errors.ErrorResponse
//	@Router		/kwaaka-admin/stoplist/attribute [post]
func (server *Server) KwaakaAdminStopListByAttributeID(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "kwaaka_admin update stoplist by attribute id request",
		Request: c.Request,
	})

	var req models.StopListByAttributeIDRequest

	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if err := server.stopListService.UpdateStopListByAttributeID(c.Request.Context(), req.IsAvailabe, req.StoreID, req.AttributeID); err != nil {
		server.Logger.Errorf("update stoplist by attribute id error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (server *Server) GetStoresInRestaurantGroupByQuery(c *gin.Context) {
	restGroupId := c.Param("restaurant_group_id")
	query := c.Query("query")
	LegalEntities := c.QueryArray("legal_entity")

	stores, err := server.storeService.GetStoresInRestGroupByName(c, restGroupId, query, LegalEntities)
	if err != nil {
		server.Logger.Errorf(errorKey, err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stores)
}

// CancelOrderKwaakaAdmin
//
//	@Tags		kwaaka-admin
//	@Title		Method for cancel order by KwaakaAdmin
//	@Security	ApiKeyAuth
//	@Summary	Method for cancel order by KwaakaAdmin
//	@Param		order_id	param	string	true	order_id
//	@Success	204
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	400	{object}	errors.ErrorResponse	e
//	@Router		v1/kwaaka-admin/cancelOrder/{order_id} [delete]
func (server *Server) CancelOrderKwaakaAdmin(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "kwaaka_admin cancel order request",
		Request: c.Request,
	})

	orderID := c.Param("order_id")

	if err := server.orderCancellationService.CancelOrderByAggregator(c.Request.Context(), orderID, "kwaaka_admin"); err != nil {
		server.Logger.Errorf("kwaaka admin cancel order error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	if err := server.orderKwaaka3plService.Cancel3plOrder(c.Request.Context(), orderID, coreModels.OrderInfoForTelegramMsg{}); err != nil {
		server.Logger.Errorf("kwaaka admin cancel dispatcher error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// UpdateKwaakaAdminBusyMode godoc
//
//	@Summary		update kwaaka_admin busy mode status and value
//	@Description	update kwaaka_admin busy mode status and value
//	@Tags			settings
//	@Accept			json
//	@Security		ApiTokenAuth
//	@Param			request	body	[]dto.BusyModeRequest	true	"request"
//	@Success		204
//	@Failure		400 {object} detector.ErrorResponse
//	@Failure		500 {object} detector.ErrorResponse
//	@Router			/v1/kwaaka-admin/busy-mode [post]
func (server *Server) UpdateKwaakaAdminBusyMode(c *gin.Context) {
	var req []integrationApiModels.BusyModeRequest
	if err := c.BindJSON(&req); err != nil {
		server.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}
	err := server.storeService.UpdateKwaakaAdminBusyMode(c.Request.Context(), req)
	if err != nil {
		server.Logger.Errorf("kwaaka admin update busy mode error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{
			Msg: err.Error(),
		})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}

// SetActualSeqNumber
//
//	@Tags		store
//	@Title		Method for rkeeper7xml set actual seq number
//	@Security	ApiKeyAuth
//	@Summary	Method rkeeper7xml set actual seq number
//	@Param		restaurant_id	path	string	true	"restaurant_id"
//	@Success	204
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/set-seq-number/{restaurant_id} [get]
func (server *Server) SetActualSeqNumber(c *gin.Context) {
	restID := c.Param("restaurant_id")

	store, err := server.storeService.GetByID(c.Request.Context(), restID)
	if err != nil {
		server.Logger.Infof("restaurant with id %s not found: %s", restID, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	posService, err := server.posService.GetPosService(coreModels.Pos(store.PosType), store)
	if err != nil {
		server.Logger.Infof("get %s pos service: %s", store.PosType, err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	seqNumber, err := posService.GetSeqNumber(c.Request.Context())
	if err != nil {
		server.Logger.Infof("get seq number error: %s", err.Error())
		c.Set(errorKey, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := server.storeService.SetActualRkeeper7xmlSeqNumber(c.Request.Context(), store.ID, seqNumber); err != nil {
		server.Logger.Errorf("kwaaka admin set seq number error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// OrdersReportForRestaurant
//
//	@Tags		store
//	@Title		Method for order report for restaurant
//	@Security	ApiKeyAuth
//	@Summary	Method order report for restaurant
//	@Param		request	body	coreModels.OrderReportRequest	true	"request"
//	@Param		page	query		string				true	"page"
//	@Param		limit	query		string				true	"limit"
//	@Success	200 {object}	coreModels.OrderReportResponse
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/order-report/restaurant [post]
func (s *Server) OrdersReportForRestaurant(c *gin.Context) {
	var req coreModels.OrderReportRequest
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	page, limit, err := parsePaging(c)
	if err != nil {
		s.Logger.Errorf("page parsing error: %s", err.Error())
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	req.Pagination.Page = page
	req.Pagination.Limit = limit

	res, err := s.orderReport.OrderReportForRestaurant(c.Request.Context(), req)
	if err != nil {
		s.Logger.Errorf("order report for restaurant error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// OrderReportForRestaurantTotals
//
//	@Tags		store
//	@Title		Method for order report get totals for restaurant
//	@Security	ApiKeyAuth
//	@Summary	Method order report get totals for restaurant
//	@Param		request	body	coreModels.OrderReportRequest	true	"request"
//	@Success	200 {object}	coreModels.OrderReportResponse
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/order-report/restaurant/totals [post]
func (s *Server) OrderReportForRestaurantTotals(c *gin.Context) {
	var req coreModels.OrderReportRequest
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	res, err := s.orderReport.OrderReportForRestaurantTotals(c.Request.Context(), req)
	if err != nil {
		s.Logger.Errorf("order report for kwaaka error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// OrderReportForKwaakaTotals
//
//	@Tags		store
//	@Title		Method for order report get totals for kwaaka
//	@Security	ApiKeyAuth
//	@Summary	Method order report get totals for kwaaka
//	@Param		request	body	coreModels.OrderReportRequest	true	"request"
//	@Success	200 {object}	coreModels.OrderReportResponse
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/order-report/kwaaka/totals [post]
func (s *Server) OrderReportForKwaakaTotals(c *gin.Context) {
	var req coreModels.OrderReportRequest
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	res, err := s.orderReport.OrderReportForKwaakaTotals(c.Request.Context(), req)
	if err != nil {
		s.Logger.Errorf("order report for kwaaka error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// OrderReportToXlsx
//
//	@Tags		store
//	@Title		Method for order report generate to xlsx
//	@Security	ApiKeyAuth
//	@Summary	Method order generate to xlsx
//	@Param		request	body	coreModels.OrderReportRequest	true	"request"
//	@Success	200 {object}	[]byte
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/order-report/xlsx [post]
func (s *Server) OrderReportToXlsx(c *gin.Context) {
	var req coreModels.OrderReportRequest
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}
	res, err := s.orderReport.OrderReportToXlsx(c.Request.Context(), req)
	if err != nil {
		s.Logger.Errorf("order report to xlsx error: %s", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=OrderReport.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", res)
}

func (s *Server) DeliveryDispatcherPrice(c *gin.Context) {
	if err := s.orderReport.DeliveryDispatcherPrice(c.Request.Context()); err != nil {
		s.Logger.Errorf("delivery dispatcher price error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
