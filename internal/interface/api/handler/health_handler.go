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
	response := gin.H{
		"status":    "OK",
		"message":   "Service is healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "inmo-backend",
		"version":   "1.0.0", // You can make this dynamic
	}

	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
	} else {
		c.JSON(http.StatusOK, response)
	}
}


func (h *HealthHandler) RegisterDetailedHealthRoute(c *gin.Context) {
	response := gin.H{
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
	}

	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
	} else {
		c.JSON(http.StatusOK, response)
	}
}

func (h *HealthHandler) RegisterPingRoute(c *gin.Context) {
	response := gin.H{
		"message": "pong",
	}	

	if c.Request.Method == http.MethodHead {
		c.Status(http.StatusOK)
	} else {
		c.JSON(http.StatusOK, response)
	}
}