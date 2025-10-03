package serviceMocks

import (
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
)

type MockJWTService struct {
	mock.Mock
}

func NewMockJWTService() *MockJWTService {
	return &MockJWTService{}
}

func (m *MockJWTService) GenerateToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) RefreshToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*models.JWTClaims, error) {
	args := m.Called(token)
	if claims, ok := args.Get(0).(*models.JWTClaims); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockJWTService) GetUserIDFromClaims(claims string) (uint, error) {
	args := m.Called(claims)
	return args.Get(0).(uint), args.Error(1)
}
