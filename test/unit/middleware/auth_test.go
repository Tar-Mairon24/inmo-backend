package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"inmo-backend/internal/domain/models"
	"inmo-backend/middleware"
)

// MockClaims represents mock JWT claims
type MockClaims struct {
	ID    string
	Email string
}

type mockJWTService struct {
	validateFunc func(token string) (*MockClaims, error)
}

// Implement GetUserIDFromClaims to satisfy ports.JWTService
func (m *mockJWTService) GetUserIDFromClaims(claimsStr string) (uint, error) {
	parsedID, err := strconv.ParseUint(claimsStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(parsedID), nil
}

// Implement ValidateToken to satisfy ports.JWTService
func (m *mockJWTService) ValidateToken(token string) (jwtResponse *models.JWTClaims, err error) {
	claims, err := m.validateFunc(token)
	if err != nil {
		return nil, err
	}
	idUint := uint(0)
	if claims != nil {
		parsedID, parseErr := strconv.ParseUint(claims.ID, 10, 32)
		if parseErr == nil {
			idUint = uint(parsedID)
		}
	}
	return &models.JWTClaims{
		ID:    idUint,
		Email: claims.Email,
	}, nil
}

// Implement GenerateToken to satisfy ports.JWTService
func (m *mockJWTService) GenerateToken(user *models.User) (string, error) {
	// For testing, just return a dummy token
	return "generatedtoken", nil
}

// Implement RefreshToken to satisfy ports.JWTService
func (m *mockJWTService) RefreshToken(token string) (string, error) {
	// For testing, just return a dummy refreshed token
	return "refreshedtoken", nil
}

func TestJWTAuthMiddleware_MissingAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockService := &mockJWTService{}
	r.Use(middleware.JWTAuthMiddleware(mockService))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing or invalid token")
}

func TestJwtAuthMiddleware_MissingCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockService := &mockJWTService{}
	r.Use(middleware.JWTAuthMiddleware(mockService))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing or invalid token")
}

func TestJWTAuthMiddleware_InvalidAuthHeaderFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockService := &mockJWTService{}
	r.Use(middleware.JWTAuthMiddleware(mockService))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidTokenFormat")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing or invalid token")
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockService := &mockJWTService{
		validateFunc: func(token string) (*MockClaims, error) {
			return nil, errors.New("Invalid token or expired token")
		},
	}
	r.Use(middleware.JWTAuthMiddleware(mockService))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: "invalidtoken",
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or expired token")
}

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockService := &mockJWTService{
		validateFunc: func(token string) (*MockClaims, error) {
			return &MockClaims{
				ID:    "1",
				Email: "test@example.com",
			}, nil
		},
	}
	r.Use(middleware.JWTAuthMiddleware(mockService))
	r.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, uint(1), userID)

		userEmail, exists := c.Get("user_email")
		assert.True(t, exists)
		assert.Equal(t, "test@example.com", userEmail)
		c.String(200, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: "validtoken",
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}
