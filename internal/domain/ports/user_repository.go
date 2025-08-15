package ports

import "inmo-backend/internal/domain/models"

type UserRepository interface {
	GetAll() ([]models.UserResponse, error)
	GetByID(id uint) (*models.UserResponse, error)
	GetByEmail(email string) (*models.UserResponse, error)
	ConsultPassword(email string) (string, error)
	Create(user *models.User) (*models.UserResponse, error)
	Update(user *models.User) (*models.UserResponse, error)
	Delete(id uint) error
}
