package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/usecase"
	"inmo-backend/middleware"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
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
func (m *MockUserRepository) GetAll() ([]models.UserResponse, error) {
	args := m.Called()
	return args.Get(0).([]models.UserResponse), args.Error(1)
}
func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUserRepository) Delete(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}
func (m *MockUserRepository) GetByEmail(email string) (*models.UserResponse, error) {
	args := m.Called(email)
	user, _ := args.Get(0).(*models.UserResponse)
	return user, args.Error(1)
}

func TestUserUsecase_CreateUser(t *testing.T) {
	type testCase struct {
		name           string
		user           *models.User
		mockError      error
		wantError      bool
		shouldCallRepo bool
	}

	tests := []testCase{
		{
			name: "valid user",
			user: &models.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "testpassword",
			},
			wantError:      false,
			shouldCallRepo: true,
		},
		{
			name: "empty username",
			user: &models.User{
				Username: "",
				Email:    "test@example.com",
				Password: "testpassword",
			},
			wantError:      true,
			shouldCallRepo: false,
		},
		{
			name: "empty email",
			user: &models.User{
				Username: "testuser",
				Email:    "",
				Password: "testpassword",
			},
			wantError:      true,
			shouldCallRepo: false,
		},
		{
			name: "empty password",
			user: &models.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "",
			},
			wantError:      true,
			shouldCallRepo: false,
		},
		{
			name: "repository error",
			user: &models.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "testpassword",
			},
			mockError:      assert.AnError, // Simulate an error from the repository
			wantError:      true,
			shouldCallRepo: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			usecase := usecase.NewUserUseCase(mockRepo)

			if tc.shouldCallRepo {
				mockRepo.On("Create", tc.user).Return(tc.mockError).Once()
			}

			err := usecase.CreateUser(tc.user)

			if tc.wantError {
				assert.Error(t, err, "Expected error for case: %s", tc.name)
			} else {
				assert.NoError(t, err, "Expected no error for case: %s", tc.name)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestUserUsecase_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUseCase(mockRepo)

		email := "test@example.com"
		password := "mypassword123"

		hash, err := middleware.HashPassword(password)
		require.NoError(t, err)

		mockRepo.On("ConsultPassword", email).Return(hash, nil)

		err = uc.Login(email, password)
		assert.NoError(t, err)
	})

	t.Run("wrong password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUseCase(mockRepo)

		email := "test@example.com"
		correctPassword := "mypassword123"
		wrongPassword := "wrongpassword"

		// Hash of correct password
		hash, err := middleware.HashPassword(correctPassword)
		require.NoError(t, err)

		mockRepo.On("ConsultPassword", email).Return(hash, nil)

		// Try to login with wrong password
		err = uc.Login(email, wrongPassword)
		assert.Error(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUseCase(mockRepo)

		mockRepo.On("ConsultPassword", "notfound@test.com").Return("", errors.New("user not found"))

		err := uc.Login("notfound@test.com", "anypassword")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("hashing error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		uc := usecase.NewUserUseCase(mockRepo)

		mockRepo.On("ConsultPassword", "hashingerror@test.com").Return("", errors.New("hashing error"))

		err := uc.Login("hashingerror@test.com", "anyPassword")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "hashing error")
	})
}

func TestUserUseCase_GetAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUseCase(mockRepo)

	expectedUsers := []models.UserResponse{
		{ID: 1, Username: "user1", Email: "user1@email.com"},
		{ID: 2, Username: "user2", Email: "user2@email.com"},
	}
	mockRepo.On("GetAll").Return(expectedUsers, nil)
	users, err := uc.GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUseCase(mockRepo)

	expectedUser := &models.UserResponse{ID: 1, Username: "user1", Email: "user1@email.com"}
	mockRepo.On("GetByID", uint(1)).Return(expectedUser, nil)
	user, err := uc.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_GetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUseCase(mockRepo)

	mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("user not found"))
	user, err := uc.GetUserByID(999)
	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUseCase(mockRepo)

	userToUpdate := &models.User{
		ID:       1,
		Username: "updatedUser",
		Email:    "update@update.com",
		Password: "newpassword",
	}
	mockRepo.On("Update", userToUpdate).Return(nil)
	err := uc.UpdateUser(userToUpdate)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUseCase(mockRepo)

	userID := uint(1)
	mockRepo.On("Delete", userID).Return(nil)
	err := uc.DeleteUser(userID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_DeleteUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := usecase.NewUserUseCase(mockRepo)

	userID := uint(999)
	mockRepo.On("Delete", userID).Return(errors.New("user not found"))
	err := uc.DeleteUser(userID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)
}
