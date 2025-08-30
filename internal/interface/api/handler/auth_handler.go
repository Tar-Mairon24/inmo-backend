package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/ports"
	"inmo-backend/internal/domain/models"
)

type AuthHandler struct {
	jwtService   ports.JWTService
	authUsecase  ports.AuthUseCase
}

func NewAuthHandler(jwtService ports.JWTService, authUsecase ports.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		jwtService:  jwtService,
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) UserLogin(c *gin.Context) {
	var loginData = models.UserLoginData{}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		logrus.WithError(err).Error("Invalid login data")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse login data",
		})
		return
	}

	loginResponse, err := h.authUsecase.Login(loginData.Email, loginData.Password)
	if err != nil {
		logrus.WithError(err).Error("Login failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": err.Error(),
		})
		return
	}
	logrus.Infof("User %s logged in succesfully", loginResponse.User.Username)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    loginResponse,
		"message": "Login successful",
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var refreshTokenData = models.RefreshTokenData{}
	if err := c.ShouldBindJSON(&refreshTokenData); err != nil {
		logrus.WithError(err).Error("Invalid refresh token data")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse refresh token data",
		})
		return
	}

	newToken, err := h.jwtService.RefreshToken(refreshTokenData.Token)
	if err != nil {
		logrus.WithError(err).Error("Failed to refresh token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    newToken,
		"message": "Token refreshed successfully",
	})
}
