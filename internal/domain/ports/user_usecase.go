package ports

import "inmo-backend/internal/domain/models"

type UserUseCase interface {
	Login(email string, password string) error
	GetAllUsers() ([]models.UserResponse, error)
	GetUserByID(id uint) (*models.UserResponse, error)
	CreateUser(user *models.User) (*models.UserResponse, error)
	UpdateUser(user *models.User) (*models.UserResponse, error)
	DeleteUser(id uint) error
}