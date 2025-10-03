package usecaseMocks

import (
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
)

type MockAuthUseCase struct {
	mock.Mock
}

func NewMockAuthUseCase() *MockAuthUseCase {
	return &MockAuthUseCase{}
}

func (m *MockAuthUseCase) Login(email, password string) (*models.LoginResponse, error) {
	args := m.Called(email, password)

	if LoginResp, ok := args.Get(0).(*models.LoginResponse); ok {
		return LoginResp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthUseCase) Logout(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthUseCase) RefreshToken(data models.RefreshTokenData) (*models.RefreshTokenData, error) {
	args := m.Called(data)
	if resp, ok := args.Get(0).(*models.RefreshTokenData); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthUseCase) GetStatus(userID uint, refreshToken string) error {
	args := m.Called(userID, refreshToken)
	return args.Error(0)
}