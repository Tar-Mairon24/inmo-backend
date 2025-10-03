package api

import (
	"github.com/gin-gonic/gin"

	"inmo-backend/internal/domain/ports"
	"inmo-backend/internal/interface/api/handler"
	"inmo-backend/middleware"
)

func setupAuthRoutes(rg *gin.RouterGroup, authHandler *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", authHandler.UserLogin)   // POST /api/v1/auth/login
		auth.POST("/refresh-token", authHandler.RefreshToken) // POST /api/v1/auth/refresh-token
		auth.POST("/logout/:id", authHandler.UserLogout) // POST /api/v1/auth/logout
		auth.GET("/status", authHandler.GetStatus) // POST /api/v1/auth/status
	}
}

func setupUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, jwtService ports.JWTService, middleware middleware.AuthMiddlewareInterface) {
	users := rg.Group("/users")
	users.Use(middleware.JWTAuthMiddleware(jwtService))
	{
		users.GET("", userHandler.GetUsers)         // GET /api/v1/users
		users.GET("/:id", userHandler.GetUserByID)   // GET /api/v1/users/:id
		users.POST("", userHandler.CreateUser)       // POST /api/v1/users
		users.PUT("/:id", userHandler.UpdateUser)    // PUT /api/v1/users
		users.DELETE("/:id", userHandler.DeleteUser) // DELETE /api/v1/users/:id
	}
}

func setupPropertyRoutes(rg *gin.RouterGroup, propertyHandler *handler.PropertyHandler, jwtService ports.JWTService, middleware middleware.AuthMiddlewareInterface) {
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
