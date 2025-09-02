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
	mockJWT := new(MockJWTService)
	mockAuth := new(MockAuthUseCase)
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	user := &models.User{ID: 1, Username: "testuser", Email: "test@example.com"}
	userResp := &models.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email}
	loginResp := &models.LoginResponse{
		User:         userResp,
		Token:        "jwt-token",
		RefreshToken: "refresh-token",
	}
	mockAuth.On("Login", "test@example.com", "password123").Return(loginResp, nil)

	body := []byte(`{"email":"test@example.com","password":"password123"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UserLogin(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"data":"jwt-token"`)
	assert.Contains(t, w.Body.String(), `"message":"Login successful"`)
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" && cookie.Value == "refresh-token" {
			found = true
			break
		}
	}
	assert.True(t, found, "refresh_token cookie should be set")
	mockAuth.AssertExpectations(t)
}

func TestUserLogin_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := new(MockJWTService)
	mockAuth := new(MockAuthUseCase)
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	body := []byte(`{"email":123,"password":true}`) // invalid types
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UserLogin(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Failed to parse login data"`)
}

func TestUserLogin_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := new(MockJWTService)
	mockAuth := new(MockAuthUseCase)
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	mockAuth.On("Login", "test@example.com", "wrongpassword").Return(nil, errors.New("invalid credentials"))

	body := []byte(`{"email":"test@example.com","password":"wrongpassword"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UserLogin(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"invalid credentials"`)
	mockAuth.AssertExpectations(t)
}
func TestUserLogout_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := new(MockJWTService)
	mockAuth := new(MockAuthUseCase)
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	mockAuth.On("Logout", uint(1)).Return(nil)

	body := []byte(`{"user_id":1}`)
	req, _ := http.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UserLogout(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"message":"Logout successful"`)
	mockAuth.AssertExpectations(t)
}

func TestUserLogout_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := new(MockJWTService)
	mockAuth := new(MockAuthUseCase)
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	body := []byte(`{"user_id":"not-a-number"}`) // invalid type
	req, _ := http.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UserLogout(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Failed to parse logout data"`)
}

func TestUserLogout_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := new(MockJWTService)
	mockAuth := new(MockAuthUseCase)
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	mockAuth.On("Logout", uint(2)).Return(errors.New("logout failed"))

	body := []byte(`{"user_id":2}`)
	req, _ := http.NewRequest(http.MethodPost, "/logout", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UserLogout(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"logout failed"`)
	mockAuth.AssertExpectations(t)
}
func TestRefreshToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := new(MockJWTService)
	mockAuth := new(MockAuthUseCase)
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	input := models.RefreshTokenData{
		RefreshToken: "old-refresh-token",
		JwtToken:     "old-jwt-token",
	}
	output := &models.RefreshTokenData{
		RefreshToken: "new-refresh-token",
		JwtToken:     "new-jwt-token",
	}
	mockAuth.On("RefreshToken", input).Return(output, nil)

	body := []byte(`{"refresh_token":"old-refresh-token","jwt_token":"old-jwt-token"}`)
	req, _ := http.NewRequest(http.MethodPost, "/refresh-token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"data":"new-jwt-token"`)
	assert.Contains(t, w.Body.String(), `"message":"Token refreshed successfully"`)
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" && cookie.Value == "new-refresh-token" {
			found = true
			break
		}
	}
	assert.True(t, found, "refresh_token cookie should be set")
	mockAuth.AssertExpectations(t)
}

func TestRefreshToken_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := new(MockJWTService)
	mockAuth := new(MockAuthUseCase)
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	body := []byte(`{"refresh_token":123,"jwt_token":true}`) // invalid types
	req, _ := http.NewRequest(http.MethodPost, "/refresh-token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Failed to parse refresh token data"`)
}

func TestRefreshToken_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := new(MockJWTService)
	mockAuth := new(MockAuthUseCase)
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	input := models.RefreshTokenData{
		RefreshToken: "bad-refresh-token",
		JwtToken:     "bad-jwt-token",
	}
	mockAuth.On("RefreshToken", input).Return(nil, errors.New("refresh failed"))

	body := []byte(`{"refresh_token":"bad-refresh-token","jwt_token":"bad-jwt-token"}`)
	req, _ := http.NewRequest(http.MethodPost, "/refresh-token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"refresh failed"`)
	mockAuth.AssertExpectations(t)
}
