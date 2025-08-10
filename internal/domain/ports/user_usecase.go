package ports

import "inmo-backend/internal/domain/models"

type UserUseCaseInterfaceImpl interface {
	Login(email string, password string) error
	GetAllUsers() ([]models.UserResponse, error)
	GetUserByID(id uint) (*models.UserResponse, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}