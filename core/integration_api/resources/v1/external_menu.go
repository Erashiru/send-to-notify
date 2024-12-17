package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	_ "github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/core/externalapi/utils"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/rs/zerolog/log"
	"net/http"
)

// GetMenu docs
//
//	@Tags		external
//	@Title		Get menu by store id
//	@Param		Authorization	header		string	true	"bearer"
//	@Param		restaurant_id	path		string	true	"restaurant_id"
//	@Success	200				{object}	models.Menu
//	@Failure	401				{object}	[]errors.ErrorResponse
//	@Failure	400				{object}	[]errors.ErrorResponse
//	@Failure	500				{object}	[]errors.ErrorResponse
//	@Router		/v1/menu/{restaurant_id}/composition [get]
func (server *Server) GetMenu(c *gin.Context) {
	storeID := c.Param(storePath)

	svc, ok := c.Get("service")
	if !ok {
		server.Logger.Infof("service query is empty")
		c.Set(errorKey, "service query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown delivery service",
		}})
		return
	}

	secret, ok := c.Get("client_secret")
	if !ok {
		server.Logger.Infof("client_secret query is empty")
		c.Set(errorKey, "client_secret query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown client secret",
		}})
		return
	}

	service := svc.(string)
	clientSecret := secret.(string)

	log.Trace().Msgf("GetMenu external: service: %s, store_id: %s, secret: %s", service, storeID, clientSecret)

	menu, err := server.externalMenuManager.GetMenu(c.Request.Context(), storeID, service, clientSecret)
	if err != nil {
		server.Logger.Infof("get menu error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("get menu error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	c.Writer.Header().Set("Content-Type", "application/vnd.eats.menu.composition.v2+json")
	c.JSON(http.StatusOK, menu)
}

// GetAvailability docs
//
//	@Tags		external
//	@Title		Get stop list from active aggregator menu
//	@Param		Authorization	header		string	true	"bearer"
//	@Param		restaurant_id	path		string	true	"restaurant_id"
//	@Success	200				{object}	models.StopListResponse
//	@Failure	401				{object}	[]errors.ErrorResponse
//	@Failure	400				{object}	[]errors.ErrorResponse
//	@Failure	500				{object}	[]errors.ErrorResponse
//	@Router		/v1/menu/{restaurant_id}/availability [get]
func (server *Server) GetAvailability(c *gin.Context) {
	storeID := c.Param(storePath)

	svc, ok := c.Get("service")
	if !ok {
		server.Logger.Infof("service query is empty")
		c.Set(errorKey, "service query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown delivery service",
		}})
		return
	}

	secret, ok := c.Get("client_secret")
	if !ok {
		server.Logger.Infof("client_secret query is empty")
		c.Set(errorKey, "client_secret query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown client secret",
		}})
		return
	}

	service := svc.(string)
	clientSecret := secret.(string)

	products, attributes, err := server.stopListService.GetStopListByDeliveryService(c.Request.Context(), storeID, service, clientSecret)
	if err != nil {
		server.Logger.Infof("get availability error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("get availability error: %s", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	slProducts := menuModels.ToStopListProducts(products)
	slAttributes := menuModels.ToStopListAttributes(attributes)

	stopList := menuModels.StopListResponse{
		Products:   slProducts,
		Attributes: slAttributes,
	}

	err = server.stopListService.AddYandexTransaction(c.Request.Context(), storeID, models.YANDEX.String(), slProducts, slAttributes)
	if err != nil {
		log.Err(err).Msgf("add yandex transaction for store: %s, stop list products: %+v, attributes: %+v", storeID, slProducts, slAttributes)
		//
		//server.Logger.Infof("add yandex transaction error: %s", err.Error())
		//c.Set(errorKey, fmt.Sprintf("add yandex transaction error: %s", err))
		//c.AbortWithStatusJSON(http.StatusInternalServerError, []errors.ErrorResponse{{
		//	Code:        http.StatusInternalServerError,
		//	Description: err.Error(),
		//}})
		//return
	}

	response := utils.ParseStopList(stopList)

	log.Info().Msgf("successfully send yandex stoplist for store id: %s, stop list products: %+v, attributes: %+v", storeID, slProducts, slAttributes)
	c.Writer.Header().Set("Content-Type", "application/vnd.eats.menu.availability.v2+json")
	c.JSON(http.StatusOK, response)
}

// GetRemains docs
//
//	@Tags		external
//	@Title		Get remains for retail store
//	@Param		Authorization	header		string	true	"bearer"
//	@Param		restaurant_id	path		string	true	"restaurant_id"
//	@Success	200				{object}	models.StopListResponse
//	@Failure	401				{object}	[]errors.ErrorResponse
//	@Failure	400				{object}	[]errors.ErrorResponse
//	@Failure	500				{object}	[]errors.ErrorResponse
//	@Router		/v1/nomenclature/{restaurant_id}/availability [get]
func (server *Server) GetRemains(c *gin.Context) {
	storeID := c.Param(storePath)

	svc, ok := c.Get("service")
	if !ok {
		server.Logger.Infof("service query is empty")
		c.Set(errorKey, "service query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown delivery service",
		}})
		return
	}

	secret, ok := c.Get("client_secret")
	if !ok {
		server.Logger.Infof("client_secret query is empty")
		c.Set(errorKey, "client_secret query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown client secret",
		}})
		return
	}

	service := svc.(string)
	clientSecret := secret.(string)

	products, attributes, err := server.stopListService.GetRetailRemains(c.Request.Context(), storeID, service, clientSecret)
	if err != nil {
		server.Logger.Infof("get availability error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("get availability error: %s", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	stopList := menuModels.StopListResponse{
		Products:   products,
		Attributes: attributes,
	}

	response := utils.ConvertStopListToRemains(stopList)

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}

func (server *Server) GetRetailMenu(c *gin.Context) {
	storeID := c.Param(storePath)

	svc, ok := c.Get("service")
	if !ok {
		server.Logger.Infof("service query is empty")
		c.Set(errorKey, "service query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown delivery service",
		}})
		return
	}

	secret, ok := c.Get("client_secret")
	if !ok {
		server.Logger.Infof("client_secret query is empty")
		c.Set(errorKey, "client_secret query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown client secret",
		}})
		return
	}

	service := svc.(string)
	clientSecret := secret.(string)

	menu, err := server.externalMenuManager.GetRetailMenu(c.Request.Context(), storeID, service, clientSecret)
	if err != nil {
		server.Logger.Infof("get retail menu error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("get retail menu error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.JSON(http.StatusOK, menu)
}
