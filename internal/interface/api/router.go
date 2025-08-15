package api

import (
	"github.com/gin-gonic/gin"

	"inmo-backend/cmd/di"
)

func SetupRouter(handlers *di.Handlers) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		setupHealthRoutes(v1, handlers.HealthHandler)
		setupUserRoutes(v1, handlers.UserHandler)
		setupPropertyRoutes(v1, handlers.PropertyHandler)
	}

	return r
}
