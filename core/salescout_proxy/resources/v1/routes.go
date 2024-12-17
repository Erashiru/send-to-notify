package v1

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (server *Server) register(engine *gin.Engine) {
	engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Phone", "Usergroup", "token"},
	}))

	v1 := engine.Group("/api")
	v1.Use(server.secretMiddleware(server.StarterAppToken))

	kaspiApi := v1.Group("/kaspi-api")
	{
		kaspiApi.POST("/create-token", server.CreateKaspiToken)
		kaspiApi.POST("/create-link", server.CreateKaspiLink)
		kaspiApi.GET(fmt.Sprintf("/status/:%s", "payment_id"), server.GetKaspiStatusByID)
		kaspiApi.POST("/refund", server.RefundKaspiPayment)
	}
}
