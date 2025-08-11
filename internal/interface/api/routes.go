package api

import (
	"inmo-backend/internal/interface/api/handler"

	"github.com/gin-gonic/gin"
)

// setupUserRoutes sets up all user-related routes
func setupUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	users := rg.Group("/users")
	{
		users.GET("", userHandler.GetUsers)        // GET /api/v1/users
		users.GET("/:id", userHandler.GetUserByID) // GET /api/v1/users/:id
		users.POST("", userHandler.CreateUser)     // POST /api/v1/users
		users.PUT("/:id", userHandler.UpdateUser)  // PUT /api/v1/users
		users.DELETE("/:id", userHandler.DeleteUser) // DELETE /api/v1/users/:id
		users.POST("/login", userHandler.UserLogin) // POST /api/v1/users/login
	}
}

func setupHealthRoutes(rg *gin.RouterGroup, healthHandler *handler.HealthHandler) {
	health := rg.Group("/health")
	{
		health.Match([]string{"GET", "HEAD"}, "", healthHandler.RegisterHealthRoutes)						// GET api/v1/health
		health.Match([]string{"GET", "HEAD"}, "/detailed", healthHandler.RegisterDetailedHealthRoute)		// GET api/v1/health/detailed
		health.Match([]string{"GET", "HEAD"}, "/ping", healthHandler.RegisterPingRoute)					// GET api/v1/health/ping
	}
}