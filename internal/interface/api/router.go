package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Add basic health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	// TODO: Add your API routes here
	// Example: api/v1 route group
	v1 := r.Group("/api/v1")
	{
		// Add user routes here later
		v1.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Users endpoint - TODO: implement",
			})
		})
	}

	return r
}
