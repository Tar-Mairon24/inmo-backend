package usecase_test

import (
	"errors"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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

type MockJWTService struct {
	mock.Mock
}

func TestMain(m *testing.M){
	err := godotenv.Load("../../../.env")
	if err != nil {
		logrus.Error("Could not load .env: ", err)
	}
	m.Run()
}

// GenerateToken mocks the token generation for a user
func (m *MockJWTService) GenerateToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (*models.JWTClaims, error) {
	args := m.Called(token)
	if claims, ok := args.Get(0).(*models.JWTClaims); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockJWTService) RefreshToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
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
func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}


func TestUserUseCase_Login(t *testing.T) {
	t.Run("should return error when email is empty", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		result, err := uc.Login("", "password")

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "email and password cannot be empty", err.Error())
	})

	t.Run("should return error when password is empty", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		result, err := uc.Login("test@example.com", "")

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "email and password cannot be empty", err.Error())
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		mockRepo.On("GetByEmail", "test@example.com").Return(nil, errors.New("user not found"))

		result, err := uc.Login("test@example.com", "password")

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when password verification fails", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		hashedPassword, _ := middleware.HashPassword("correctpassword")
		user := &models.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)

		result, err := uc.Login("test@example.com", "wrongpassword")

		assert.Nil(t, result)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when token generation fails", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		hashedPassword, _ := middleware.HashPassword("password")
		user := &models.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)
		mockJWT.On("GenerateToken", user).Return("", errors.New("token generation failed"))

		result, err := uc.Login("test@example.com", "password")

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "token generation failed", err.Error())
		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("should login successfully", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		hashedPassword, _ := middleware.HashPassword("password")
		user := &models.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)
		mockJWT.On("GenerateToken", user).Return("mock-token", nil)

		result, err := uc.Login("test@example.com", "password")

		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "mock-token", result.Token)
		assert.Equal(t, user.ToUserResponse(), result.User)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
}
func TestUserUseCase_GetAllUsers(t *testing.T) {
	t.Run("should return all users successfully", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		users := []models.UserResponse{
			{ID: 1, Username: "user1", Email: "user1@example.com"},
			{ID: 2, Username: "user2", Email: "user2@example.com"},
		}

		mockRepo.On("GetAll").Return(users, nil)

		result, err := uc.GetAllUsers()

		assert.NoError(t, err)
		assert.Equal(t, users, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		mockRepo.On("GetAll").Return(nil, errors.New("repo error"))

		result, err := uc.GetAllUsers()

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "repo error", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
func TestUserUseCase_GetUserByID(t *testing.T) {
	t.Run("should return user response when user exists", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		expectedUser := &models.UserResponse{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockRepo.On("GetByID", uint(1)).Return(expectedUser, nil)

		result, err := uc.GetUserByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		mockRepo.On("GetByID", uint(2)).Return(nil, errors.New("user not found"))

		result, err := uc.GetUserByID(2)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
func TestUserUseCase_CreateUser(t *testing.T) {
	t.Run("should return error when password is empty", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		user := &models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "",
		}

		result, err := uc.CreateUser(user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "password cannot be empty", err.Error())
	})

	t.Run("should return error when username is empty", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		user := &models.User{
			Username: "",
			Email:    "test@example.com",
			Password: "password",
		}

		result, err := uc.CreateUser(user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "username cannot be empty", err.Error())
	})

	t.Run("should return error when email is empty", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		user := &models.User{
			Username: "testuser",
			Email:    "",
			Password: "password",
		}

		result, err := uc.CreateUser(user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "email cannot be empty", err.Error())
	})

	t.Run("should create user successfully", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		user := &models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password",
		}

		// Use real hash for password
		hashedPassword, _ := middleware.HashPassword("password")
		expectedUser := &models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: hashedPassword,
		}
		expectedResponse := &models.UserResponse{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockRepo.On("Create", mock.MatchedBy(func(u *models.User) bool {
			return u.Username == expectedUser.Username &&
				u.Email == expectedUser.Email &&
				u.Password != "" && u.Password != "password"
		})).Return(expectedResponse, nil)

		result, err := uc.CreateUser(user)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockRepo.AssertExpectations(t)
	})
}
func TestUserUseCase_UpdateUser(t *testing.T) {
	t.Run("should update user successfully", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		user := &models.User{
			ID:       1,
			Username: "updateduser",
			Email:    "updated@example.com",
			Password: "newpassword",
		}
		expectedResponse := &models.UserResponse{
			ID:       1,
			Username: "updateduser",
			Email:    "updated@example.com",
		}

		mockRepo.On("Update", user).Return(expectedResponse, nil)

		result, err := uc.UpdateUser(user)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		user := &models.User{
			ID:       2,
			Username: "failuser",
			Email:    "fail@example.com",
			Password: "password",
		}

		mockRepo.On("Update", user).Return(nil, errors.New("update failed"))

		result, err := uc.UpdateUser(user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "update failed", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
func TestUserUseCase_DeleteUser(t *testing.T) {
	t.Run("should delete user successfully", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		mockRepo.On("Delete", uint(1)).Return(nil)

		err := uc.DeleteUser(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		mockRepo := &MockUserRepository{}
		mockJWT := &MockJWTService{}
		uc := usecase.NewUserUseCase(mockRepo, mockJWT)

		mockRepo.On("Delete", uint(2)).Return(errors.New("delete failed"))

		err := uc.DeleteUser(2)

		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
		mockRepo.AssertExpectations(t)
	})
}





