package ports

import "inmo-backend/internal/domain/models"

type JWTService interface {
    GenerateToken(user *models.User) (string, error)
    ValidateToken(tokenString string) (*models.JWTClaims, error)
    RefreshToken(tokenString string) (string, error)
    GetUserIDFromClaims(tokenString string) (uint, error)
}