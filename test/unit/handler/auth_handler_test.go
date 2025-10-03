package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/interface/api/handler"
	"inmo-backend/test/mocks/service"
	"inmo-backend/test/mocks/usecase"
)

func TestUserLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
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
	assert.Contains(t, w.Body.String(), `"data":`, loginResp.User)
	assert.Contains(t, w.Body.String(), `"message":"Login successful"`)
	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" && cookie.Value == "refresh-token" {
			found = true
			assert.Equal(t, "/", cookie.Path)
			assert.True(t, cookie.HttpOnly)
			assert.False(t, cookie.Secure)
			assert.Equal(t, 3600*24*7, cookie.MaxAge) // 1 week
			break
		}
	}
	assert.True(t, found, "refresh_token cookie should be set")
	mockAuth.AssertExpectations(t)
}

func TestUserLogin_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
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
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
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
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	// Setup mocks
	mockJWT.On("GetUserIDFromClaims", "valid-jwt-token").Return(uint(42), nil)
	mockAuth.On("Logout", uint(42)).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/logout/42", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "valid-jwt-token"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "42"}}

	handler.UserLogout(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"message":"Logout successful"`)
	mockJWT.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

func TestUserLogout_BadRequest_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	req, _ := http.NewRequest(http.MethodPost, "/logout/abc", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler.UserLogout(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Failed to parse user ID"`)
}

func TestUserLogout_BadRequest_MissingJWTToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	req, _ := http.NewRequest(http.MethodPost, "/logout/42", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "42"}}

	handler.UserLogout(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Missing JWT token"`)
}

func TestUserLogout_Unauthorized_LogoutAttempt(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	// JWT returns wrong user ID
	mockJWT.On("GetUserIDFromClaims", "valid-jwt-token").Return(uint(99), nil)

	req, _ := http.NewRequest(http.MethodPost, "/logout/42", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "valid-jwt-token"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "42"}}

	handler.UserLogout(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"Unauthorized logout attempt"`)
	mockJWT.AssertExpectations(t)
}

func TestUserLogout_Unauthorized_LogoutFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	mockJWT.On("GetUserIDFromClaims", "valid-jwt-token").Return(uint(42), nil)
	mockAuth.On("Logout", uint(42)).Return(errors.New("some logout error"))

	req, _ := http.NewRequest(http.MethodPost, "/logout/42", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "valid-jwt-token"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "42"}}

	handler.UserLogout(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"some logout error"`)
	mockJWT.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

func TestUserLogout_BadRequest_UserNotLoggedIn(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	mockJWT.On("GetUserIDFromClaims", "valid-jwt-token").Return(uint(42), nil)
	mockAuth.On("Logout", uint(42)).Return(errors.New("no token found for the given user ID, user was not logged in"))

	req, _ := http.NewRequest(http.MethodPost, "/logout/42", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "valid-jwt-token"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "42"}}

	handler.UserLogout(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Bad request"`)
	assert.Contains(t, w.Body.String(), `"message":"User was not logged in"`)
	mockJWT.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}
func TestRefreshToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)
	// Prepare mock response
	oldRefreshToken := "old-refresh-token"
	oldJwtToken := "old-jwt-token"
	newRefreshToken := "new-refresh-token"
	newJwtToken := "new-jwt-token"
	refreshTokenData := models.RefreshTokenData{
		RefreshToken: oldRefreshToken,
		JwtToken:     oldJwtToken,
	}
	newTokenData := &models.RefreshTokenData{
		RefreshToken: newRefreshToken,
		JwtToken:     newJwtToken,
	}
	mockAuth.On("RefreshToken", refreshTokenData).Return(newTokenData, nil)

	req, _ := http.NewRequest(http.MethodPost, "/refresh-token", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: oldRefreshToken})
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: oldJwtToken})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"message":"Token refreshed successfully"`)
	cookies := w.Result().Cookies()
	foundRefresh := false
	foundJwt := false
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" && cookie.Value == newRefreshToken {
			foundRefresh = true
			assert.Equal(t, "/", cookie.Path)
			assert.True(t, cookie.HttpOnly)
			assert.False(t, cookie.Secure)
			assert.Equal(t, 3600*24*7, cookie.MaxAge)
		}
		if cookie.Name == "jwt_token" && cookie.Value == newJwtToken {
			foundJwt = true
			assert.Equal(t, "/", cookie.Path)
			assert.True(t, cookie.HttpOnly)
			assert.False(t, cookie.Secure)
			assert.Equal(t, 3600*24*7, cookie.MaxAge)
		}
	}
	assert.True(t, foundRefresh, "refresh_token cookie should be set")
	assert.True(t, foundJwt, "jwt_token cookie should be set")
	mockAuth.AssertExpectations(t)
}

func TestRefreshToken_BadRequest_MissingRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	req, _ := http.NewRequest(http.MethodPost, "/refresh-token", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "jwt-token"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Missing refresh token"`)
}

func TestRefreshToken_BadRequest_MissingJWTToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	req, _ := http.NewRequest(http.MethodPost, "/refresh-token", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-token"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Missing JWT token"`)
}

func TestRefreshToken_Unauthorized_RefreshFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	refreshTokenData := models.RefreshTokenData{
		RefreshToken: "refresh-token",
		JwtToken:     "jwt-token",
	}
	mockAuth.On("RefreshToken", refreshTokenData).Return(nil, errors.New("invalid refresh token"))

	req, _ := http.NewRequest(http.MethodPost, "/refresh-token", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refresh-token"})
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "jwt-token"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.RefreshToken(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"invalid refresh token"`)
	mockAuth.AssertExpectations(t)
}
func TestGetStatus_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	jwtToken := "valid-jwt-token"
	refreshToken := "valid-refresh-token"
	claims := &models.JWTClaims{ID: 123}

	mockJWT.On("ValidateToken", jwtToken).Return(claims, nil)
	mockAuth.On("GetStatus", claims.ID, refreshToken).Return(nil)

	req, _ := http.NewRequest(http.MethodGet, "/status", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: jwtToken})
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetStatus(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), `"message":"User is logged in"`)
	mockJWT.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

func TestGetStatus_Unauthorized_MissingJWTToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	req, _ := http.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetStatus(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"Missing JWT token"`)
}

func TestGetStatus_Unauthorized_InvalidJWTToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	jwtToken := "invalid-jwt-token"
	mockJWT.On("ValidateToken", jwtToken).Return(nil, errors.New("invalid token"))

	req, _ := http.NewRequest(http.MethodGet, "/status", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: jwtToken})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetStatus(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"Invalid JWT token"`)
	mockJWT.AssertExpectations(t)
}

func TestGetStatus_Unauthorized_MissingRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	jwtToken := "valid-jwt-token"
	claims := &models.JWTClaims{ID: 123}
	mockJWT.On("ValidateToken", jwtToken).Return(claims, nil)

	req, _ := http.NewRequest(http.MethodGet, "/status", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: jwtToken})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetStatus(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"Missing refresh token"`)
	mockJWT.AssertExpectations(t)
}

func TestGetStatus_Unauthorized_InvalidRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockJWT := serviceMocks.NewMockJWTService()
	mockAuth := usecaseMocks.NewMockAuthUseCase()
	handler := handler.NewAuthHandler(mockJWT, mockAuth)

	jwtToken := "valid-jwt-token"
	refreshToken := "invalid-refresh-token"
	claims := &models.JWTClaims{ID: 123}
	mockJWT.On("ValidateToken", jwtToken).Return(claims, nil)
	mockAuth.On("GetStatus", claims.ID, refreshToken).Return(errors.New("invalid refresh token"))

	req, _ := http.NewRequest(http.MethodGet, "/status", nil)
	req.AddCookie(&http.Cookie{Name: "jwt_token", Value: jwtToken})
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: refreshToken})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetStatus(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Unauthorized"`)
	assert.Contains(t, w.Body.String(), `"message":"Invalid refresh token"`)
	mockJWT.AssertExpectations(t)
	mockAuth.AssertExpectations(t)
}

