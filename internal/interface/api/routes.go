package api

import (
	"net/http"
	"time"

	"inmo-backend/internal/interface/api/handler"

	"github.com/gin-gonic/gin"
)

// setupUserRoutes sets up all user-related routes
func setupUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler) {
	users := rg.Group("/users")
	{
		users.GET("", userHandler.GetUsers)        // GET /api/v1/users
		users.GET("/:id", userHandler.GetUserByID) // GET /api/v1/users/:id
	}
}

func setupHealthRoutes(r *gin.Engine) {
	// Basic health check - perfect for Docker HEALTHCHECK
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "OK",
			"message":   "Server is running",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service":   "inmo-backend",
			"version":   "1.0.0", // You can make this dynamic
		})
	})

	// Detailed health check - for monitoring and frontend
	r.GET("/health/detailed", func(c *gin.Context) {
		// You can add database connectivity check here
		c.JSON(http.StatusOK, gin.H{
			"status":    "OK",
			"message":   "All services are operational",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service":   "inmo-backend",
			"version":   "1.0.0",
			"checks": gin.H{
				"database":   "connected", // TODO: Add actual DB ping
				"memory":     "ok",        // TODO: Add memory check
				"disk_space": "ok",        // TODO: Add disk space check
			},
			"uptime": "running", // TODO: Calculate actual uptime
		})
	})

	// Simple ping endpoint - minimal response for load balancers
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
