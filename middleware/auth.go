package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/ports"
)

func JWTAuthMiddleware(jwtService ports.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			logrus.WithError(err).Warn("Token validation failed")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("user_email", claims.Email)
		c.Set("user_claims", claims)

		c.Next()
	}
}
