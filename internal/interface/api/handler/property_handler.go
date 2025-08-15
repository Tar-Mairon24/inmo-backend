package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
)

type PropertyHandler struct {
	propertyUsecase ports.PropertyUseCase
}

func NewPropertyHandler(propertyUsecase ports.PropertyUseCase) *PropertyHandler {
	return &PropertyHandler{
		propertyUsecase: propertyUsecase,
	}
}

func (h *PropertyHandler) GetProperties(c *gin.Context) {
	logrus.Info("GetProperties endpoint called")

	properties, err := h.propertyUsecase.GetAllProperties()
	if err != nil {
		logrus.WithError(err).Error("Failed to get properties")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve properties",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, properties)
}

func (h *PropertyHandler) GetPropertyByID(c *gin.Context) {
	logrus.Info("GetPropertyByID endpoint called")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		logrus.WithError(err).Error("Invalid property ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid property ID",
			"message": "Property ID must be a positive integer",
		})
		return
	}

	property, err := h.propertyUsecase.GetPropertyByID(uint(id))
	if err != nil {
		logrus.WithError(err).Error("Failed to get property by ID")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve property",
			"message": err.Error(),
		})
		return
	}

	if property == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Property not found",
			"message": "No property found with the given ID",
		})
		return
	}

	c.JSON(http.StatusOK, property)
}

func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	logrus.Info("CreateProperty endpoint called")

	var property models.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": "Please provide valid property data",
		})
		return
	}

	newProperty, err := h.propertyUsecase.CreateProperty(&property)
	if err != nil {
		logrus.WithError(err).Error("Failed to create property")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create property",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, newProperty)
}

func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	logrus.Info("UpdateProperty endpoint called")

	var property models.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		logrus.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"message": "Please provide valid property data",
		})
		return
	}

	updatedProperty, err := h.propertyUsecase.UpdateProperty(&property)
	if err != nil {
		logrus.WithError(err).Error("Failed to update property")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update property",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedProperty)
}

func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	logrus.Info("DeleteProperty endpoint called")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		logrus.WithError(err).Error("Invalid property ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid property ID",
			"message": "Property ID must be a positive integer",
		})
		return
	}

	err = h.propertyUsecase.DeleteProperty(uint(id))
	if err != nil {
		logrus.WithError(err).Error("Failed to delete property")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete property",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}