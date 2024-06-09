package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/vsitnev/sync-manager/docs"
	"github.com/vsitnev/sync-manager/internal/service"
)

func NewRouter(handler *gin.Engine, services *service.Services) {
	handler.Use(gin.Recovery())

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handler.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "sync-manager",
			"status":  "ok",
		})
	})

	v1 := handler.Group("/api/v1") // <--- you can set auth middleware here
	{
		newMessageRoutes(v1.Group("/messages"), services.Message)
	}
}
