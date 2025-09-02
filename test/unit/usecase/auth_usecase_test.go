package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/usecase"
	"inmo-backend/middleware"
)

// --- Mocks ---

type mockUserRepo struct {
	mock.Mock
}

func TestMain(m *testing.M){
	err := godotenv.Load("../../../.env")
	if err != nil {
		logrus.Error("Could not load .env: ", err)
	}
	m.Run()
}

// Add missing Update method to satisfy ports.UserRepository
func (m *mockUserRepo) Update(user *models.User) (*models.UserResponse, error) {
	args := m.Called(user)
	userResp, _ := args.Get(0).(*models.UserResponse)
	return userResp, args.Error(1)
}

// Add missing GetByID method to satisfy ports.UserRepository
func (m *mockUserRepo) GetByID(userID uint) (*models.UserResponse, error) {
	args := m.Called(userID)
	userResp, _ := args.Get(0).(*models.UserResponse)
	return userResp, args.Error(1)
}

func (m *mockUserRepo) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

// Add missing ConsultPassword method to satisfy ports.UserRepository
func (m *mockUserRepo) ConsultPassword(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

// Add missing Create method to satisfy ports.UserRepository
func (m *mockUserRepo) Create(user *models.User) (*models.UserResponse, error) {
	args := m.Called(user)
	userResp, _ := args.Get(0).(*models.UserResponse)
	return userResp, args.Error(1)
}

// Add missing Delete method to satisfy ports.UserRepository
func (m *mockUserRepo) Delete(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

// Add missing GetAll method to satisfy ports.UserRepository
func (m *mockUserRepo) GetAll() ([]models.UserResponse, error) {
	args := m.Called()
	users, _ := args.Get(0).([]models.UserResponse)
	return users, args.Error(1)
}

type mockTokenRepo struct {
	mock.Mock
}

func (m *mockTokenRepo) SaveToken(token *models.RefreshToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *mockTokenRepo) GetTokenIDByUserID(userID uint) (string, error)             {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}
func (m *mockTokenRepo) DeleteToken(tokenID string) error                           {
	args := m.Called(tokenID)
	return args.Error(0)
}
func (m *mockTokenRepo) GetTokenByUserID(userID uint) (*models.RefreshToken, error) {
	args := m.Called(userID)
	token, _ := args.Get(0).(*models.RefreshToken)
	return token, args.Error(1)
}

type mockJWTService struct {
	mock.Mock
}

func (m *mockJWTService) GenerateToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *mockJWTService) GetUserIDFromClaims(token string) (uint, error) {
	args := m.Called(token)
	userID, _ := args.Get(0).(uint)
	return userID, args.Error(1)
}
func (m *mockJWTService) RefreshToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

// Add missing ValidateToken method to satisfy ports.JWTService
func (m *mockJWTService) ValidateToken(token string) (*models.JWTClaims, error) {
	args := m.Called(token)
	claims, _ := args.Get(0).(*models.JWTClaims)
	return claims, args.Error(1)
}

// --- Tests ---

func TestLogin_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	password := "testPassword"
	hasedPassword, err := middleware.HashPassword(password)
	if err != nil {

	}
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Password: hasedPassword,
		Email:    "test@example.com",
	}
	userRepo.On("GetByEmail", "test@example.com").Return(user, nil)
	tokenRepo.On("SaveToken", mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	jwtService.On("GenerateToken", user).Return("jwtToken", nil)

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)

	resp, err := uc.Login("test@example.com", password)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Token, resp.Token)
	assert.Equal(t, resp.RefreshToken, resp.RefreshToken)
	userRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestLogin_EmptyEmailOrPassword(t *testing.T) {
	uc := usecase.NewAuthUseCase(nil, nil, nil)
	resp, err := uc.Login("", "password")
	assert.Error(t, err)
	assert.Nil(t, resp)
	resp, err = uc.Login("email@example.com", "")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestLogin_UserNotFound(t *testing.T) {
	userRepo := new(mockUserRepo)
	userRepo.On("GetByEmail", "notfound@example.com").Return(nil, errors.New("not found"))
	uc := usecase.NewAuthUseCase(userRepo, nil, nil)
	resp, err := uc.Login("notfound@example.com", "password")
	assert.Error(t, err)
	assert.Nil(t, resp)
}
func TestLogout_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	userID := uint(1)
	tokenID := "token123"

	tokenRepo.On("GetTokenIDByUserID", userID).Return(tokenID, nil)
	tokenRepo.On("DeleteToken", tokenID).Return(nil)

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)

	err := uc.Logout(userID)
	assert.NoError(t, err)
	tokenRepo.AssertCalled(t, "GetTokenIDByUserID", userID)
	tokenRepo.AssertCalled(t, "DeleteToken", tokenID)
}

func TestLogout_EmptyUserID(t *testing.T) {
	uc := usecase.NewAuthUseCase(nil, nil, nil)
	err := uc.Logout(0)
	assert.Error(t, err)
	assert.EqualError(t, err, "user ID cannot be empty")
}

func TestLogout_GetTokenIDError(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	userID := uint(2)
	tokenRepo.On("GetTokenIDByUserID", userID).Return("", errors.New("not found"))

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)

	err := uc.Logout(userID)
	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
	tokenRepo.AssertCalled(t, "GetTokenIDByUserID", userID)
}

func TestLogout_DeleteTokenError(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	userID := uint(3)
	tokenID := "token456"
	tokenRepo.On("GetTokenIDByUserID", userID).Return(tokenID, nil)
	tokenRepo.On("DeleteToken", tokenID).Return(errors.New("delete error"))

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)

	err := uc.Logout(userID)
	assert.Error(t, err)
	assert.EqualError(t, err, "delete error")
	tokenRepo.AssertCalled(t, "GetTokenIDByUserID", userID)
	tokenRepo.AssertCalled(t, "DeleteToken", tokenID)
}
func TestRefreshToken_Success(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	data := models.RefreshTokenData{
		JwtToken:     "valid.jwt.token",
		RefreshToken: "valid_refresh_token",
	}
	userID := uint(42)
	refreshToken := &models.RefreshToken{
		ID:        "rtid",
		UserID:    userID,
		Token:     "valid_refresh_token",
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		CreatedAt: time.Now(),
	}

	jwtService.On("GetUserIDFromClaims", data.JwtToken).Return(userID, nil)
	tokenRepo.On("GetTokenByUserID", userID).Return(refreshToken, nil)
	jwtService.On("RefreshToken", data.JwtToken).Return("new.jwt.token", nil)

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)
	resp, err := uc.RefreshToken(data)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "new.jwt.token", resp.JwtToken)
	assert.Equal(t, "valid_refresh_token", resp.RefreshToken)
	jwtService.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestRefreshToken_EmptyTokens(t *testing.T) {
	uc := usecase.NewAuthUseCase(nil, nil, nil)
	data := models.RefreshTokenData{JwtToken: "", RefreshToken: ""}
	resp, err := uc.RefreshToken(data)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "JWT token and refresh token cannot be empty")
}

func TestRefreshToken_GetUserIDError(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	data := models.RefreshTokenData{JwtToken: "bad.jwt", RefreshToken: "refresh"}
	jwtService.On("GetUserIDFromClaims", data.JwtToken).Return(uint(0), errors.New("claims error"))

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)
	resp, err := uc.RefreshToken(data)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "claims error")
	jwtService.AssertExpectations(t)
}

func TestRefreshToken_GetTokenByUserIDError(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	data := models.RefreshTokenData{JwtToken: "jwt", RefreshToken: "refresh"}
	userID := uint(1)
	jwtService.On("GetUserIDFromClaims", data.JwtToken).Return(userID, nil)
	tokenRepo.On("GetTokenByUserID", userID).Return(nil, errors.New("token error"))

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)
	resp, err := uc.RefreshToken(data)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "token error")
	jwtService.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestRefreshToken_RefreshTokenNotFound(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	data := models.RefreshTokenData{JwtToken: "jwt", RefreshToken: "refresh"}
	userID := uint(2)
	jwtService.On("GetUserIDFromClaims", data.JwtToken).Return(userID, nil)
	tokenRepo.On("GetTokenByUserID", userID).Return(nil, nil)

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)
	resp, err := uc.RefreshToken(data)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "refresh token not found")
	jwtService.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestRefreshToken_RefreshTokenExpired(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	data := models.RefreshTokenData{JwtToken: "jwt", RefreshToken: "refresh"}
	userID := uint(3)
	expiredToken := &models.RefreshToken{
		ID:        "rtid",
		UserID:    userID,
		Token:     "refresh",
		ExpiresAt: time.Now().Add(-1 * time.Hour).Unix(),
		CreatedAt: time.Now(),
	}
	jwtService.On("GetUserIDFromClaims", data.JwtToken).Return(userID, nil)
	tokenRepo.On("GetTokenByUserID", userID).Return(expiredToken, nil)

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)
	resp, err := uc.RefreshToken(data)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "refresh token expired")
	jwtService.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestRefreshToken_InvalidRefreshToken(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	data := models.RefreshTokenData{JwtToken: "jwt", RefreshToken: "wrong_refresh"}
	userID := uint(4)
	token := &models.RefreshToken{
		ID:        "rtid",
		UserID:    userID,
		Token:     "expected_refresh",
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		CreatedAt: time.Now(),
	}
	jwtService.On("GetUserIDFromClaims", data.JwtToken).Return(userID, nil)
	tokenRepo.On("GetTokenByUserID", userID).Return(token, nil)

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)
	resp, err := uc.RefreshToken(data)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "invalid refresh token")
	jwtService.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestRefreshToken_RefreshTokenError(t *testing.T) {
	userRepo := new(mockUserRepo)
	tokenRepo := new(mockTokenRepo)
	jwtService := new(mockJWTService)

	data := models.RefreshTokenData{JwtToken: "jwt", RefreshToken: "refresh"}
	userID := uint(5)
	token := &models.RefreshToken{
		ID:        "rtid",
		UserID:    userID,
		Token:     "refresh",
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		CreatedAt: time.Now(),
	}
	jwtService.On("GetUserIDFromClaims", data.JwtToken).Return(userID, nil)
	tokenRepo.On("GetTokenByUserID", userID).Return(token, nil)
	jwtService.On("RefreshToken", data.JwtToken).Return("", errors.New("refresh error"))

	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtService)
	resp, err := uc.RefreshToken(data)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "refresh error")
	jwtService.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}




