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
	var loginData = models.LoginData{}
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
	logrus.Infof("User %s logged in successfully", loginResponse.User.Username)
	c.SetCookie("refresh_token", loginResponse.RefreshToken, 3600*24*7, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    loginResponse.Token,
		"message": "Login successful",
	})
}

func (h *AuthHandler) UserLogout(c *gin.Context) {
	var logoutData = models.LogoutData{}
	if err := c.ShouldBindJSON(&logoutData); err != nil {
		logrus.WithError(err).Error("Invalid logout data")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse logout data",
		})
		return
	}

	if err := h.authUsecase.Logout(logoutData.UserID); err != nil {
		logrus.WithError(err).Error("Logout failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
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

	newToken, err := h.authUsecase.RefreshToken(refreshTokenData)
	if err != nil {
		logrus.WithError(err).Error("Failed to refresh token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": err.Error(),
		})
		return
	}

	c.SetCookie("refresh_token", newToken.RefreshToken, 3600*24*7, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    newToken.JwtToken,
		"message": "Token refreshed successfully",
	})
}
