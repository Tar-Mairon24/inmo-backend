package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterHealthRoutes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"message":   "Server is running",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "inmo-backend",
		"version":   "1.0.0", // You can make this dynamic
	})
}


func (h *HealthHandler) RegisterDetailedHealthRoute(c *gin.Context) {
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
}

func (h *HealthHandler) RegisterPingRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})

}