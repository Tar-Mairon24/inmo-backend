package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/interface/api/handler"
)


// MockJWTService is a mock implementation of JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

// Add missing RefreshToken method to satisfy ports.JWTService interface
func (m *MockJWTService) RefreshToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

// Add missing ValidateToken method to satisfy ports.JWTService interface
func (m *MockJWTService) ValidateToken(token string) (*models.JWTClaims, error) {
	args := m.Called(token)
	if claims, ok := args.Get(0).(*models.JWTClaims); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}

// Add missing GetUserIDFromClaims method to satisfy ports.JWTService interface
func (m *MockJWTService) GetUserIDFromClaims(claims string) (uint, error) {
	args := m.Called(claims)
	return args.Get(0).(uint), args.Error(1)
}

// MockAuthUseCase is a mock implementation of AuthUseCase
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Login(email, password string) (*models.LoginResponse, error) {
	args := m.Called(email, password)

	if LoginResp, ok := args.Get(0).(*models.LoginResponse); ok {
		return LoginResp, args.Error(1)
	}
	return nil, args.Error(1)
	
}
func (m *MockAuthUseCase) GetAllUsers() ([]models.UserResponse, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserResponse), args.Error(1)
}
func (m *MockAuthUseCase) GetUserByID(id uint) (*models.UserResponse, error) {
	args := m.Called(id)
	if user, ok := args.Get(0).(*models.UserResponse); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAuthUseCase) CreateUser(user *models.User) (*models.UserResponse, error) {
	args := m.Called(user)
	if userResp, ok := args.Get(0).(*models.UserResponse); ok {
		return userResp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAuthUseCase) UpdateUser(user *models.User) (*models.UserResponse, error) {
	args := m.Called(user)
	if userResp, ok := args.Get(0).(*models.UserResponse); ok {
		return userResp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAuthUseCase) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Add missing Logout method to satisfy ports.AuthUseCase interface
func (m *MockAuthUseCase) Logout(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

// Add missing RefreshToken method to satisfy ports.AuthUseCase interface
func (m *MockAuthUseCase) RefreshToken(data models.RefreshTokenData) (*models.RefreshTokenData, error) {
	args := m.Called(data)
	if resp, ok := args.Get(0).(*models.RefreshTokenData); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}


func TestUserLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockAuthUseCase)
	mockJWTService := new(MockJWTService)
	handler := handler.NewAuthHandler(mockJWTService, mockUseCase)

	loginData := `{"email":"test@example.com","password":"password123"}`
	loginResp := &models.LoginResponse{
		User: &models.UserResponse{
			Username: "testuser",
			Email:    "test@example.com",
		},
		Token: "token123",
	}
	mockUseCase.On("Login", "test@example.com", "password123").Return(loginResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(loginData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UserLogin(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Login successful"`)
	assert.Contains(t, w.Body.String(), `"token":"token123"`)
	mockUseCase.AssertExpectations(t)
}

func TestUserLogin_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockAuthUseCase)
	mockJWTService := new(MockJWTService)
	handler := handler.NewAuthHandler(mockJWTService, mockUseCase)

	invalidJSON := `{"email": "test@example.com", "password":}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UserLogin(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Failed to parse login data"`)
}

func TestUserLogin_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockAuthUseCase)
	mockJWTService := new(MockJWTService)
	handler := handler.NewAuthHandler(mockJWTService, mockUseCase)

	loginData := `{"email":"test@example.com","password":"wrongpass"}`
	mockUseCase.On("Login", "test@example.com", "wrongpass").Return(nil, errors.New("invalid credentials"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(loginData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UserLogin(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"invalid credentials"`)
	mockUseCase.AssertExpectations(t)
}
func TestRefreshToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockAuthUseCase)
	mockJWTService := new(MockJWTService)
	handler := handler.NewAuthHandler(mockJWTService, mockUseCase)

	refreshData := `{"token":"oldtoken123"}`
	mockJWTService.On("RefreshToken", "oldtoken123").Return("newtoken456", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/refresh", bytes.NewBufferString(refreshData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Token refreshed successfully"`)
	assert.Contains(t, w.Body.String(), `"data":"newtoken456"`)
	mockJWTService.AssertExpectations(t)
}

func TestRefreshToken_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockAuthUseCase)
	mockJWTService := new(MockJWTService)
	handler := handler.NewAuthHandler(mockJWTService, mockUseCase)

	invalidJSON := `{"token":}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/refresh", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Failed to parse refresh token data"`)
}

func TestRefreshToken_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockAuthUseCase)
	mockJWTService := new(MockJWTService)
	handler := handler.NewAuthHandler(mockJWTService, mockUseCase)

	refreshData := `{"token":"expiredtoken"}`
	mockJWTService.On("RefreshToken", "expiredtoken").Return("", errors.New("token expired"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/refresh", bytes.NewBufferString(refreshData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"token expired"`)
	mockJWTService.AssertExpectations(t)
}

