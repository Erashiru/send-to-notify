package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	_ "github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/core/externalapi/resources/http/v1/dto"
	"github.com/kwaaka-team/orders-core/core/externalapi/utils"
	"net/http"
)

// GetRestaurants docs
//
//		@Tags		external
//		@Title		Get restaurants by delivery service
//		@Param		Authorization	header		string	true	"bearer"
//	    @Param      service         query       string  false    "delivery service"
//	    @Param      client_secret   query       string  false    "client secret"
//		@Success	200				{object}	models.GetStoreResponse
//		@Failure	401				{object}	[]errors.ErrorResponse
//		@Failure	400				{object}	[]errors.ErrorResponse
//		@Failure	500				{object}	[]errors.ErrorResponse
//		@Router		/v1/restaurants [get]
func (server *Server) GetRestaurants(c *gin.Context) {
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

	stores, err := server.externalMenuManager.GetStores(c.Request.Context(), service, clientSecret)

	if err != nil {
		server.Logger.Infof("get stores error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("get stores error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	c.JSON(http.StatusOK, stores)
}

// GetPromos docs
//
//	@Tags		external
//	@Title		Get promos by store id
//	@Param		Authorization	header		string	true	"bearer"
//	@Param		restaurant_id	path		string	true	"restaurant_id"
//	@Success	200				{object}	models.Promo
//	@Failure	401				{object}	[]errors.ErrorResponse
//	@Failure	400				{object}	[]errors.ErrorResponse
//	@Failure	500				{object}	[]errors.ErrorResponse
//	@Router		/v1/menu/{restaurant_id}/promos [get]
func (server *Server) GetPromos(c *gin.Context) {
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

	promos, err := server.externalMenuManager.GetPromos(c.Request.Context(), storeID, service, clientSecret)
	if err != nil {
		server.Logger.Infof("get promos error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("get promos error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: err.Error(),
		}})
		return
	}

	c.JSON(http.StatusOK, promos)
}

// CreateToken
//
//	@Tags		external
//	@Title		Create access token
//	@Summary	Authenticate in Kwaaka system
//	@Success	200	{object}	dto.SuccessResponse
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/security/oauth/token [post]
func (server *Server) CreateToken(c *gin.Context) {

	var req dto.AuthenticateData

	// Serialize request body
	if err := c.Bind(&req); err != nil {
		server.Logger.Infof("bind error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("bind error: %s", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Description: err.Error(),
			Code:        http.StatusBadRequest,
		}})
		return
	}

	authClient, err := server.externalAuthManager.FindByIDAndSecret(c.Request.Context(), req.ClientID, req.ClientSecret)

	if err != nil {
		server.Logger.Infof("find secret error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("find secret error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusNotFound, []errors.ErrorResponse{{
			Description: "Client not found",
			Code:        http.StatusNotFound,
		}})
		return
	}

	jwtService := utils.JWTAuthService(server.Config.AppSecret)
	accessToken := jwtService.GenerateJWTToken(authClient.ClientID, authClient.ClientSecret, authClient.Service)

	if accessToken == "" {
		server.Logger.Infof("unexpected error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("undexpected error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Description: "Unexpected error.",
			Code:        http.StatusBadRequest,
		}})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		AccessToken: accessToken,
	})
}
