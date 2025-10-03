package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
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
	c.SetCookie("refresh_token", loginResponse.RefreshToken, 3600*24*7, "/", "", false, true)
	c.SetCookie("jwt_token", loginResponse.Token, 3600*24*7, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    loginResponse.User,
		"message": "Login successful",
	})
}

func (h *AuthHandler) UserLogout(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil || userID <= 0 {
		logrus.WithError(err).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse user ID",
		})
		return
	}
	logrus.Infof("UserLogout endpoint called for user ID: %d", userID)

	var logoutData = models.LogoutData{
		UserID: uint(userID),
	}

	jwtToken, err := c.Cookie("jwt_token")
	if err != nil {
		logrus.WithError(err).Error("JWT token not found in cookies")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Missing JWT token",
		})
		return
	}
	tokenId, err := h.jwtService.GetUserIDFromClaims(jwtToken)
	if err != nil || tokenId != uint(userID) {
		logrus.WithError(err).Error("Unauthorized logout attempt")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Unauthorized logout attempt",
		})
		return
	}	

	if err := h.authUsecase.Logout(logoutData.UserID); err != nil {
		if err.Error() == "no token found for the given user ID, user was not logged in" {
			logrus.Warn("User was not logged in")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad request",
				"message": "User was not logged in",
			})
			return
		}

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
	var refreshToken string
	var jwtToken string
	var err error

	refreshToken, err = c.Cookie("refresh_token")
	if err != nil {
		logrus.WithError(err).Error("Refresh token not found in cookies")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Missing refresh token",
		})
		return
	}
	jwtToken, err = c.Cookie("jwt_token")
	if err != nil {
		logrus.WithError(err).Error("JWT token not found in cookies")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Missing JWT token",
		})
		return
	}
	refreshTokenData.RefreshToken = refreshToken
	refreshTokenData.JwtToken = jwtToken

	logrus.Info("RefreshToken endpoint called")
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
	c.SetCookie("jwt_token", newToken.JwtToken, 3600*24*7, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token refreshed successfully",
	})
}

func (h *AuthHandler) GetStatus(c *gin.Context) {
	jwtToken, err := c.Cookie("jwt_token")
	if err != nil {
		logrus.WithError(err).Error("JWT token not found in cookies")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Missing JWT token",
		})
		return
	}

	claims, err := h.jwtService.ValidateToken(jwtToken)
	if err != nil {
		logrus.WithError(err).Error("Invalid JWT token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid JWT token",
		})
		return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		logrus.WithError(err).Error("Refresh token not found in cookies")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Missing refresh token",
		})
		return
	}

	err = h.authUsecase.GetStatus(claims.ID, refreshToken)
	if err != nil {
		logrus.WithError(err).Error("Invalid refresh token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User is logged in",
	})
}
