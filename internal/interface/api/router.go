package api

import (
	"github.com/gin-gonic/gin"

	"inmo-backend/cmd/di"
)

func SetupRouter(handlers *di.Handlers, services di.Services) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		setupHealthRoutes(v1, handlers.HealthHandler)
		setupUserRoutes(v1, handlers.UserHandler, services.JwtService)
		setupPropertyRoutes(v1, handlers.PropertyHandler, services.JwtService)
	}

	return r
}
