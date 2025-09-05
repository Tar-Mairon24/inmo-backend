package ports

import "inmo-backend/internal/domain/models"

type TokenRepository interface {
	SaveToken(token *models.RefreshToken) error
	DeleteToken(tokenID string) error
	GetTokenIDByUserID(userID uint) (string, error)
	GetTokenByUserID(userID uint) (*models.RefreshToken, error)
}
 