package ports

import "inmo-backend/internal/domain/models"

type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id uint) (*models.User, error)
}
