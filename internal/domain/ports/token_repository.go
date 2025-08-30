package ports

import "inmo-backend/internal/domain/models"

type TokenRepository interface {
	SaveToken(token *models.RefreshToken) error
	DeleteToken(tokenID string) error
	GetIDByToken(token string) (string, error)
}
 