package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
)

type UserHandler struct {
	userUsecase ports.UserUseCase
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userUsecase ports.UserUseCase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) UserLogin(c *gin.Context) {
	var loginData = models.UserLoginData{}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		logrus.WithError(err).Error("Invalid login data")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse login data",
		})
		return
	}

	if err := h.userUsecase.Login(loginData.Email, loginData.Password); err != nil {
		logrus.WithError(err).Error("Login failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Invalid email or password",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}

// GetUsers handles GET /api/v1/users
func (h *UserHandler) GetUsers(c *gin.Context) {
	logrus.Info("GetUsers endpoint called")

	// Use the actual usecase to get users from database
	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		logrus.WithError(err).Error("Failed to get users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve users",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    users,
		"message": "Users retrieved successfully",
		"count":   len(users),
	})
}

// GetUserByID handles GET /api/v1/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	logrus.Infof("GetUserByID endpoint called with ID: %s", userIDStr)

	// Convert string ID to uint
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Invalid user ID format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	// Use the actual usecase to get user from database
	user, err := h.userUsecase.GetUserByID(uint(userID))
	if err != nil {
		logrus.WithError(err).Error("Failed to get user")
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "User retrieved successfully",
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse user data",
		})
		return
	}

	if err := h.userUsecase.CreateUser(&user); err != nil {
		logrus.WithError(err).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    user,
		"message": "User created successfully",
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": "Failed to parse user data",
		})
		return
	}

	if err := h.userUsecase.UpdateUser(&user); err != nil {
		logrus.WithError(err).Error("Failed to update user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    user,
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	logrus.Infof("DeleteUser endpoint called with ID: %s", userIDStr)

	// Convert string ID to uint
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		logrus.WithError(err).Error("Invalid user ID format")
		logrus.Info("🔍 DEBUG: Returning 400 (Bad Request)")  // ← Debug
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"message": "User ID must be a valid number",
		})
		return
	}

	if err := h.userUsecase.DeleteUser(uint(userID)); err != nil {
		logrus.WithError(err).Error("Failed to delete user")
		logrus.Info("🔍 DEBUG: Returning 500 (Internal Server Error)")  // ← Debug
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"message": err.Error(),
		})
		return
	}

	logrus.Info("🔍 DEBUG: Returning 204 (No Content)")
	c.JSON(http.StatusNoContent, nil)
}


