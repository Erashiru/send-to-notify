package v1

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

const (
	lambdaFunctionNameEnv    = "AWS_LAMBDA_FUNCTION_NAME"
	lambdaIntegrationApiName = "integration-api"
	withRedirectEnv          = "WITH_REDIRECT"
)

func (server *Server) readBodyAndSetAgain(c *gin.Context) []byte {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}

	server.Logger.Infof("read glovo c.Request.Body and set to body, len body=%d", len(body))

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return body
}

func (server *Server) withRedirect() bool {
	withRedirect := os.Getenv(withRedirectEnv)

	server.Logger.Infof("env WITH_REDIRECT=%s", withRedirect)

	return withRedirect == "true"
}

func (server *Server) isTopPartner(ctx context.Context, externalStoreId, deliveryService string) bool {
	store, err := server.storeService.GetByExternalIdAndDeliveryService(ctx, externalStoreId, deliveryService)
	if err != nil {
		server.Logger.Errorf("get store by external id=%s and delivery service=%s error=%v", externalStoreId, deliveryService, err)
		return false
	}

	storeGroup, err := server.storeGroupService.GetStoreGroupByStoreID(ctx, store.ID)
	if err != nil {
		server.Logger.Errorf("get store group by store id=%s error=%v", store.ID, err)
		return false
	}

	if server.storeGroupService.IsTopPartner(storeGroup) {
		server.Logger.Infof("restaurant with external store id=%s is top partner=true", externalStoreId)
		return true
	}

	server.Logger.Infof("restaurant with external store id=%s is top partner=false", externalStoreId)

	return false
}

func (server *Server) redirectRequest(c *gin.Context, body []byte) {
	redirectURL := os.Getenv("REDIRECT_URL")

	server.Logger.Infof("env REDIRECT_URL=%s", redirectURL)

	if redirectURL != "" {
		server.redirectMiddleware(c, body, redirectURL)
	}
}

func (server *Server) redirectMiddleware(c *gin.Context, body []byte, target string) {
	url := target + c.Request.URL.Path

	server.Logger.Infof("redirect url=%s", url)

	client := http.Client{}

	request, err := http.NewRequest(c.Request.Method, url, bytes.NewBuffer(body))
	if err != nil {
		server.Logger.Infof("prepare http request error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, "")
		return
	}

	server.Logger.Infof("Set HEADERS to redirect request")
	for key, values := range c.Request.Header {
		for _, value := range values {
			request.Header.Add(key, value)
			server.Logger.Infof("%s=%s", key, value)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		server.Logger.Infof("send http request error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, "")
		return
	}
	defer response.Body.Close()

	server.Logger.Infof("success send http request to %s", url)

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		server.Logger.Infof("read response body error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, "")
		return
	}

	for key, values := range response.Header {
		c.Writer.Header()[key] = values
	}

	c.Writer.WriteHeader(response.StatusCode)
	c.Writer.Write(respBody)
}
