package v1

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	ginAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	salescout_proxy "github.com/kwaaka-team/orders-core/core/salescout_proxy/service"
)

var ginLambda *ginAdapter.GinLambda

type Server struct {
	Router                *gin.Engine
	SaleScoutProxeService *salescout_proxy.Service
	StarterAppToken       string
}

func NewServer(saleScoutProxeService *salescout_proxy.Service, authToken string) *Server {
	server := &Server{
		Router:                gin.Default(),
		SaleScoutProxeService: saleScoutProxeService,
		StarterAppToken:       authToken,
	}

	ginLambda = ginAdapter.New(server.Router)
	server.register(server.Router)

	server.Router.RedirectTrailingSlash = true
	server.Router.RedirectFixedPath = true
	server.Router.HandleMethodNotAllowed = true

	return server
}

func (server *Server) GinProxy(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request ginproxy: %v\n", req)
	return ginLambda.Proxy(req)
}
