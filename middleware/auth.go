package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/ports"
)

type AuthMiddlewareInterface interface {
	JWTAuthMiddleware(jwtService ports.JWTService) gin.HandlerFunc
	GenerateRefreshToken() (string, string, error)
}

type authMiddleware struct{}

func NewAuthMiddleware() AuthMiddlewareInterface {
	return &authMiddleware{}
}

func (m *authMiddleware) JWTAuthMiddleware(jwtService ports.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		if cookieToken, err := c.Cookie("jwt_token"); err == nil {
			token = cookieToken
		} else {
			logrus.Warn("JWT token not found in cookies")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Missing or invalid token",
			})
			c.Abort()
			return
		}

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

func (m *authMiddleware) GenerateRefreshToken() (string, string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", "", err
	}
	token := base64.StdEncoding.EncodeToString(b)
	hash := sha256.Sum256([]byte(token))
	id := hex.EncodeToString(hash[:])
	return token, id, nil
}
