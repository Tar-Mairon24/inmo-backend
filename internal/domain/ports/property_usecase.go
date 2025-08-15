package ports

import "inmo-backend/internal/domain/models"

type PropertyUseCase interface {
	GetAllProperties() ([]models.PropertyResponse, error)
	GetPropertyByID(id uint) (*models.PropertyResponse, error)
	CreateProperty(property *models.Property) (*models.PropertyResponse, error)
	UpdateProperty(property *models.Property) (*models.PropertyResponse, error)
	DeleteProperty(id uint) error
}