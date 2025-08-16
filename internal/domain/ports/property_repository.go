package ports

import "inmo-backend/internal/domain/models"

type PropertyRepository interface {
	GetAll() ([]models.PropertyResponse, error)
	GetByID(id uint) (*models.PropertyResponse, error)
	Create(property *models.Property) (*models.PropertyResponse, error)
	Update(property *models.Property) (*models.PropertyResponse, error)
	Delete(id uint) error
}