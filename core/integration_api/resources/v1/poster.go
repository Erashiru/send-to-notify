package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/domain/logger"
	models2 "github.com/kwaaka-team/orders-core/service/pos/models/poster"
	"io"
	"net/http"
)

const (
	CODE    = "code"
	ACCOUNT = "account"
)

// WebHookEventsHandlerPoster docs
//	@Tags		poster
//	@Title		Method for update order status and stop list
//	@Param		event	body	models.WHEvent	true	"event"
//	@Success	200
//	@Router		/poster/events [post]
func (server *Server) WebHookEventsHandlerPoster(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "poster request",
		Request: *c.Request,
	})

	// if response will not be with status 200 poster stops sending webhooks
	c.Status(http.StatusOK)

	var req models2.WHEvent

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		server.Logger.Error(logger.LoggerInfo{
			System:   "poster response error",
			Response: err,
		})
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		server.Logger.Error(logger.LoggerInfo{
			System:   "poster response error",
			Response: err,
		})
		return
	}

	if err := server.posterWebhookEvent(c, req); err != nil {
		server.Logger.Error(logger.LoggerInfo{
			System:   "poster response error",
			Response: err,
		})
		return
	}
}

// CodeReceiverHandlerPoster docs
//	@Tags		poster
//	@Title		Method for receive code
//	@Param		code	query	string	true	"code"
//	@Param		account	query	string	true	"account"
//	@Success	200
//	@Router		/poster/events [get]
func (server *Server) CodeReceiverHandlerPoster(c *gin.Context) {
	server.Logger.Info(logger.LoggerInfo{
		System:  "poster request",
		Request: *c.Request,
	})

	code, _ := c.GetQuery(CODE)
	account, _ := c.GetQuery(ACCOUNT)

	err := server.posterService.StoreAuth(c.Request.Context(), code, account)
	if err != nil {
		server.Logger.Error(logger.LoggerInfo{
			System:   "poster response error",
			Response: err,
		})
	}

	c.AbortWithStatusJSON(http.StatusOK, map[string]string{
		"code":    code,
		"account": account,
	})
}

func (server *Server) posterWebhookEvent(ctx context.Context, webhook models2.WHEvent) error {
	switch webhook.Object {
	case "stock":
		if err := server.updateStopListByBalance(ctx, webhook); err != nil {
			return err
		}
	case "incoming_order":
		if err := server.statusUpdateService.UpdateOrderStatus(ctx, webhook.AccountNumber+webhook.ObjectID, webhook.Action, ""); err != nil {
			return err
		}
	case "product", "dish":
		if err := server.updateStopList(ctx, webhook); err != nil {
			return err
		}
	default:
		server.Logger.Info(logger.LoggerInfo{
			System:   "poster response",
			Response: fmt.Sprintf("another hook %v ", webhook.Object),
		})
	}

	return nil
}

func (server *Server) updateStopListByBalance(ctx context.Context, webhook models2.WHEvent) error {
	storesItemsMap, err := server.posterService.GetStoplistByBalanceItems(ctx, webhook)
	if err != nil {
		return err
	}

	for storeID, items := range storesItemsMap {
		if items.ProductID != "" {
			if err := server.stopListService.UpdateStopListByPosProductID(ctx, items.IsAvailable, storeID, items.ProductID); err != nil {
				server.Logger.Info(logger.LoggerInfo{
					System:   "poster response",
					Response: fmt.Sprintf("updte stoplist by productID err: %v, productID: %s, isAvailable: %v, storeID: %s", err, items.ProductID, items.IsAvailable, storeID),
				})
			}
		}
		if items.AttributeID != "" {
			if err := server.stopListService.UpdateStopListByAttributeID(ctx, items.IsAvailable, storeID, items.AttributeID); err != nil {
				server.Logger.Info(logger.LoggerInfo{
					System:   "poster response",
					Response: fmt.Sprintf("updte stoplist by attributeID err: %v, attributeID: %s, isAvailable: %v, storeID: %s", err, items.AttributeID, items.IsAvailable, storeID),
				})
			}
		}
	}

	return nil
}

func (server *Server) updateStopList(ctx context.Context, webhook models2.WHEvent) error {
	storesItemsMap, err := server.posterService.GetStoplistItems(ctx, webhook)
	if err != nil {
		return err
	}

	for storeID, items := range storesItemsMap {
		if items.ProductID != "" {
			if err := server.stopListService.UpdateStopListByPosProductID(ctx, items.IsAvailable, storeID, items.ProductID); err != nil {
				server.Logger.Info(logger.LoggerInfo{
					System:   "poster response",
					Response: fmt.Sprintf("updte stoplist by productID err: %v, productID: %s, isAvailable: %v, storeID: %s", err, items.ProductID, items.IsAvailable, storeID),
				})
			}
		}
		if items.AttributeID != "" {
			if err := server.stopListService.UpdateStopListByAttributeID(ctx, items.IsAvailable, storeID, items.AttributeID); err != nil {
				server.Logger.Info(logger.LoggerInfo{
					System:   "poster response",
					Response: fmt.Sprintf("updte stoplist by attributeID err: %v, attributeID: %s, isAvailable: %v, storeID: %s", err, items.AttributeID, items.IsAvailable, storeID),
				})
			}
		}
	}

	return nil
}
