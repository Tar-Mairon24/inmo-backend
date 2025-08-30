package ports

import "inmo-backend/internal/domain/models"

type AuthUseCase interface {
	Login(email string, password string) (*models.LoginResponse, error)
	Logout(token string) error
}
