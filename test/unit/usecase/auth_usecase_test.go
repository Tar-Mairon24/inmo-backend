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
)
	

type MockJWTService struct {
	mock.Mock
}

type MockTokenRepository struct {
	mock.Mock
}

// Implement missing method to satisfy ports.TokenRepository
func (m *MockTokenRepository) GetTokenIDByUserID(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenRepository) GetTokenByUserID(userID uint) (*models.RefreshToken, error) {
	args := m.Called(userID)
	if token, ok := args.Get(0).(*models.RefreshToken); ok {
		return token, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestMain(m *testing.M){
	err := godotenv.Load("../../../.env")
	if err != nil {
		logrus.Error("Could not load .env: ", err)
	}
	m.Run()
}

func (m *MockTokenRepository) SaveToken(token *models.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockTokenRepository) GetIDByToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockTokenRepository) DeleteToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

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

// Implement missing method to satisfy ports.JWTService
func (m *MockJWTService) GetUserIDFromClaims(claims string) (uint, error) {
	args := m.Called(claims)
	return args.Get(0).(uint), args.Error(1)
}

func TestAuthUseCase_Login_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)
	mockJWT := new(MockJWTService)

	user := &models.User{
		ID:       123,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)
	mockJWT.On("GenerateToken", user).Return("jwt-token", nil)
	mockTokenRepo.On("SaveToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil)

	uc := usecase.NewAuthUseCase(mockRepo, mockTokenRepo, mockJWT)

	// Act
	resp, err := uc.Login("test@example.com", "password")

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "jwt-token", resp.Token)
	assert.Equal(t, user.Username, resp.User.Username)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
	mockTokenRepo.AssertExpectations(t)
}

func TestAuthUseCase_Login_EmptyEmailOrPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)
	mockJWT := new(MockJWTService)
	uc := usecase.NewAuthUseCase(mockRepo, mockTokenRepo, mockJWT)

	resp, err := uc.Login("", "password")
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "email and password cannot be empty", err.Error())

	resp, err = uc.Login("test@example.com", "")
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "email and password cannot be empty", err.Error())
}

func TestAuthUseCase_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)
	mockJWT := new(MockJWTService)
	mockRepo.On("GetByEmail", "notfound@example.com").Return(nil, errors.New("not found"))
	uc := usecase.NewAuthUseCase(mockRepo, mockTokenRepo, mockJWT)

	resp, err := uc.Login("notfound@example.com", "password")
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestAuthUseCase_Login_PasswordVerificationFailed(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)
	mockJWT := new(MockJWTService)
	user := &models.User{
		ID:       123,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)

	uc := usecase.NewAuthUseCase(mockRepo, mockTokenRepo, mockJWT)

	resp, err := uc.Login("test@example.com", "wrongpassword")
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "invalid password", err.Error())
}

func TestAuthUseCase_Login_RefreshTokenGenerationFailed(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)
	mockJWT := new(MockJWTService)
	user := &models.User{
		ID:       123,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)

	uc := usecase.NewAuthUseCase(mockRepo, mockTokenRepo, mockJWT)

	resp, err := uc.Login("test@example.com", "password")
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "refresh token error", err.Error())
}

func TestAuthUseCase_Login_SaveTokenFailed(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)
	mockJWT := new(MockJWTService)
	user := &models.User{
		ID:       123,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)
	mockTokenRepo.On("SaveToken", mock.AnythingOfType("*models.RefreshToken")).Return(errors.New("save token error"))

	uc := usecase.NewAuthUseCase(mockRepo, mockTokenRepo, mockJWT)

	resp, err := uc.Login("test@example.com", "password")
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "save token error", err.Error())
}

func TestAuthUseCase_Login_GenerateTokenFailed(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)
	mockJWT := new(MockJWTService)
	user := &models.User{
		ID:       123,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)
	mockTokenRepo.On("SaveToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockJWT.On("GenerateToken", user).Return("", errors.New("jwt error"))


	uc := usecase.NewAuthUseCase(mockRepo, mockTokenRepo, mockJWT)

	resp, err := uc.Login("test@example.com", "password")
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, "jwt error", err.Error())
}




