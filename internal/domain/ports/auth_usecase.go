package ports

import "inmo-backend/internal/domain/models"

type AuthUseCase interface {
	Login(email string, password string) (*models.LoginResponse, error)
	Logout(id uint) error
	RefreshToken(data models.RefreshTokenData) (*models.RefreshTokenData, error)
	GetStatus(userID uint, refreshToken string) error
}
