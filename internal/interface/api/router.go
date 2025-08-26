package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"inmo-backend/cmd/di"
)

func SetupRouter(handlers *di.Handlers, services di.Services) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")
	{
		setupHealthRoutes(v1, handlers.HealthHandler)
		setupAuthRoutes(v1, handlers.AuthHandler)
		setupUserRoutes(v1, handlers.UserHandler, services.JwtService)
		setupPropertyRoutes(v1, handlers.PropertyHandler, services.JwtService)
	}

	return r
}
