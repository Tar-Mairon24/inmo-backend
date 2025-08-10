package api

import (
	"inmo-backend/cmd/di"

	"github.com/gin-gonic/gin"
)

func SetupRouter(handlers *di.Handlers) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		setupHealthRoutes(v1, handlers.HealthHandler)
		setupUserRoutes(v1, handlers.UserHandler)
	}

	return r
}
