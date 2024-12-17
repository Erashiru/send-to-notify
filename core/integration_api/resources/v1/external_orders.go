package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/externalapi/resources/http/v1/dto"
	orderModels "github.com/kwaaka-team/orders-core/pkg/order/dto"
	"net/http"
	"sync"
	"time"
)

// GetOrders docs
//	@Tags		orders
//	@Summary	get orders
//	@Param		Authorization	header		string					true	"bearer"
//	@Param		orders_request	body		dto.GetOrdersRequest	true	"orders_request"
//	@Success	200				{object}	dto.GetOrdersResponse
//	@Failure	400				{object}	[]errors.ErrorResponse
//	@Failure	400				{object}	[]errors.ErrorResponse
//	@Router		/v1/orders [post]
func (server *Server) GetOrders(c *gin.Context) {
	var (
		err      error
		timeFrom time.Time
		timeTo   time.Time
	)

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

	clientId, ok := c.Get("client_id")
	if !ok {
		server.Logger.Infof("client_secret query is empty")
		c.Set(errorKey, "client_secret query is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "unknown client id",
		}})
		return
	}

	service := svc.(string)
	clientID := clientId.(string)

	var req dto.GetOrdersRequest

	if err = c.BindJSON(&req); err != nil {
		server.Logger.Infof("bind error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("bind error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, []errors.ErrorResponse{{
			Code:        http.StatusBadRequest,
			Description: "wrong body",
		}})
		return
	}

	if req.StartDate != "" {
		timeFrom, err = time.Parse("2006-01-02 15:04:05", req.StartDate)
		if err != nil {
			server.Logger.Infof("time parse start_date error: %s", err.Error())
			c.Set(errorKey, fmt.Sprintf("time parse start_date error: %s", err.Error()))
			c.AbortWithStatusJSON(http.StatusNotFound, []errors.ErrorResponse{{
				Description: "start date invalid format",
				Code:        http.StatusNotFound,
			}})
			return
		}
	}

	if req.EndDate != "" {
		timeTo, err = time.Parse("2006-01-02 15:04:05", req.EndDate)
		if err != nil {
			server.Logger.Infof("time parse end_date error: %s", err.Error())
			c.Set(errorKey, fmt.Sprintf("time parse end_date error: %s", err.Error()))
			c.AbortWithStatusJSON(http.StatusNotFound, []errors.ErrorResponse{{
				Description: "end date invalid format",
				Code:        http.StatusNotFound,
			}})
			return
		}
	}

	authClient, err := server.externalAuthManager.FindByID(c.Request.Context(), clientID)
	if err != nil {
		server.Logger.Infof("client not found error: %s", err.Error())
		c.Set(errorKey, fmt.Sprintf("client not found error: %s", err.Error()))
		c.AbortWithStatusJSON(http.StatusNotFound, []errors.ErrorResponse{{
			Description: "client not found",
			Code:        http.StatusNotFound,
		}})
		return
	}

	isBelongs := make(map[string]struct{})

	for _, id := range authClient.ExternalStoreIDs {
		isBelongs[id] = struct{}{}
	}

	for _, id := range req.RestaurantIDs {
		if _, exist := isBelongs[id]; !exist {
			server.Logger.Infof("restaurant id %s doesn't belongs to auth client", id)
			c.Set(errorKey, fmt.Sprintf("restaurant id %s doesn't belongs to auth client", id))
			c.AbortWithStatusJSON(http.StatusNotFound, []errors.ErrorResponse{{
				Description: fmt.Sprintf("restaurant id %s doesn't belongs to auth client", id),
				Code:        http.StatusNotFound,
			}})
			return
		}
	}

	result := dto.GetOrdersResponse{
		Infos: make([]dto.OrderInfo, len(req.RestaurantIDs)),
	}

	var wg sync.WaitGroup
	var mtx sync.Mutex

	for position, id := range req.RestaurantIDs {
		wg.Add(1)

		go func(index int, externalStoreID string) {
			defer wg.Done()

			store, err := server.externalMenuManager.FindStore(c.Request.Context(), externalStoreID, service)
			if err != nil {
				server.Logger.Infof("find store error: %s", err.Error())
				return
			}

			orders, total, err := server.externalOrderManager.GetOrders(c.Request.Context(), orderModels.OrderSelector{
				StoreID:       store.ID,
				OrderTimeFrom: timeFrom,
				OrderTimeTo:   timeTo,
			})
			if err != nil {
				server.Logger.Infof(fmt.Sprintf("get orders error: %s", err.Error()))
				return
			}

			mtx.Lock()
			result.Infos[index] = dto.OrderInfo{
				RestaurantID: externalStoreID,
				Orders:       orders,
				Total:        total,
			}
			defer mtx.Unlock()
		}(position, id)

	}

	wg.Wait()

	c.AbortWithStatusJSON(http.StatusOK, result)
}
