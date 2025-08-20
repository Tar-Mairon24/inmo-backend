package api

import (
	"github.com/gin-gonic/gin"

	"inmo-backend/internal/domain/ports"
	"inmo-backend/internal/interface/api/handler"
	"inmo-backend/middleware"
)

func setupUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, jwtService ports.JWTService) {
	users := rg.Group("/users")
	{
		users.POST("/login", userHandler.UserLogin) // POST /api/v1/users/login
	}

	protected := users.Group("")
	protected.Use(middleware.JWTAuthMiddleware(jwtService))
	{
		protected.GET("", userHandler.GetUsers)         // GET /api/v1/users
		protected.GET("/:id", userHandler.GetUserByID)   // GET /api/v1/users/:id
		protected.POST("", userHandler.CreateUser)       // POST /api/v1/users
		protected.PUT("/:id", userHandler.UpdateUser)    // PUT /api/v1/users
		protected.DELETE("/:id", userHandler.DeleteUser) // DELETE /api/v1/users/:id
	}
}

func setupPropertyRoutes(rg *gin.RouterGroup, propertyHandler *handler.PropertyHandler, jwtService ports.JWTService) {
	properties := rg.Group("/properties")
	properties.Use(middleware.JWTAuthMiddleware(jwtService))
	{
		properties.GET("", propertyHandler.GetProperties)         // GET /api/v1/properties
		properties.GET("/:id", propertyHandler.GetPropertyByID)   // GET /api/v1/properties/:id
		properties.POST("", propertyHandler.CreateProperty)       // POST /api/v1/properties
		properties.PUT("/:id", propertyHandler.UpdateProperty)    // PUT /api/v1/properties/:id
		properties.DELETE("/:id", propertyHandler.DeleteProperty) // DELETE /api/v1/properties/:id
	}
}

func setupHealthRoutes(rg *gin.RouterGroup, healthHandler *handler.HealthHandler) {
	health := rg.Group("/health")
	{
		health.Match([]string{"GET", "HEAD"}, "", healthHandler.RegisterHealthRoutes)                 // GET api/v1/health
		health.Match([]string{"GET", "HEAD"}, "/detailed", healthHandler.RegisterDetailedHealthRoute) // GET api/v1/health/detailed
		health.Match([]string{"GET", "HEAD"}, "/ping", healthHandler.RegisterPingRoute)               // GET api/v1/health/ping
	}
}
