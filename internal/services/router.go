package services

import (
	"medods-service/internal/handlers"

	"github.com/gin-gonic/gin"
)





func SetupRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/auth")
	{
		api.POST("/create", handlers.CreateUserTokens)
		api.POST("/refresh", handlers.RefreshTokens)
	}

	return router

}