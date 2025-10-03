package repositoryMock

import (
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
)

type MockUserRepository struct {
	mock.Mock
}

func NewMockUserRepo() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) Create(user *models.User) (*models.UserResponse, error) {
	args := m.Called(user)
	if userResponse, ok := args.Get(0).(*models.UserResponse); ok {
		return userResponse, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) GetByID(id uint) (*models.UserResponse, error) {
	args := m.Called(id)
	if user, ok := args.Get(0).(*models.UserResponse); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) ConsultPassword(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) GetAll() ([]models.UserResponse, error) {
	args := m.Called()
	if users, ok := args.Get(0).([]models.UserResponse); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) Update(user *models.User) (*models.UserResponse, error) {
	args := m.Called(user)
	if userResponse, ok := args.Get(0).(*models.UserResponse); ok {
		return userResponse, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserRepository) Delete(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}
