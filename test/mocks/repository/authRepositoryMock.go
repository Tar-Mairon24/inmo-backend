package repositoryMock

import (
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
)

type mockTokenRepo struct {
	mock.Mock
}

func NewMockTokenRepo() *mockTokenRepo {
	return &mockTokenRepo{}
}

func (m *mockTokenRepo) SaveToken(token *models.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockTokenRepo) GetTokenIDByUserID(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *mockTokenRepo) DeleteToken(tokenID string) error {
	args := m.Called(tokenID)
	return args.Error(0)
}

func (m *mockTokenRepo) GetTokenByUserID(userID uint) (*models.RefreshToken, error) {
	args := m.Called(userID)
	token, _ := args.Get(0).(*models.RefreshToken)
	return token, args.Error(1)
}