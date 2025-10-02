package mocks

import (
	"inmo-backend/internal/domain/models"

	"github.com/stretchr/testify/mock"
)

type MockUserUseCase struct {
	mock.Mock
}

func NewMockUserUseCase() *MockUserUseCase {
	return &MockUserUseCase{}
}

func (m *MockUserUseCase) GetAllUsers() ([]models.UserResponse, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserResponse), args.Error(1)
}
func (m *MockUserUseCase) GetUserByID(id uint) (*models.UserResponse, error) {
	args := m.Called(id)
	if user, ok := args.Get(0).(*models.UserResponse); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserUseCase) CreateUser(user *models.User) (*models.UserResponse, error) {
	args := m.Called(user)
	if userResp, ok := args.Get(0).(*models.UserResponse); ok {
		return userResp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserUseCase) UpdateUser(user *models.User) (*models.UserResponse, error) {
	args := m.Called(user)
	if userResp, ok := args.Get(0).(*models.UserResponse); ok {
		return userResp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserUseCase) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
