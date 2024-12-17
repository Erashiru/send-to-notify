package http

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	ginAdapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/jowi/config"
	"github.com/kwaaka-team/orders-core/core/jowi/managers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var ginLambda *ginAdapter.GinLambda

type Server struct {
	Router      *gin.Engine
	conf        config.Configuration
	jowiManager managers.JowiManager
}

func NewServer(jowiManager managers.JowiManager, conf config.Configuration) *Server {
	server := &Server{
		Router:      gin.Default(),
		conf:        conf,
		jowiManager: jowiManager,
	}

	ginLambda = ginAdapter.New(server.Router)
	server.Register(server.Router)

	server.Router.RedirectTrailingSlash = true
	server.Router.RedirectFixedPath = true
	server.Router.HandleMethodNotAllowed = true

	return server
}

func (server *Server) Register(engine *gin.Engine) {

	jowi := engine.Group("/jowi")

	swagger := jowi.Group("/documentation/jowi/swagger")

	{
		swagger.GET("/*any", ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.URL("doc.json"),
		))
	}

	{
		jowi.POST("/events", server.Events)
	}

	engine.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(405, gin.H{
			"description": "Unsupported method",
			"code":        405,
		})
	})

	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"description": "Route not found",
			"code":        404,
		})
	})

}

func (server *Server) GinProxy(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Request ginproxy: %v\n", req)
	return ginLambda.Proxy(req)
}
