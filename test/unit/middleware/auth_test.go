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
	assert.Contains(t, w.Body.String(), "Authorization header required")
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
	assert.Contains(t, w.Body.String(), "Invalid authorization header format")
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockService := &mockJWTService{
		validateFunc: func(token string) (*MockClaims, error) {
			return nil, errors.New("invalid token")
		},
	}
	r.Use(middleware.JWTAuthMiddleware(mockService))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid or expired token")
}

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockClaims := &MockClaims{ID: "123", Email: "test@example.com"}
	mockService := &mockJWTService{
		validateFunc: func(token string) (*MockClaims, error) {
			if token == "validtoken" {
				return mockClaims, nil
			}
			return nil, errors.New("invalid token")
		},
	}
	r.Use(func(c *gin.Context) {
		// Adapt to ports.JWTService interface
		type jwtServiceAdapter struct{ *mockJWTService }
		c.Set("jwtService", &jwtServiceAdapter{mockService})
		middleware.JWTAuthMiddleware(mockService)(c)
	})
	r.GET("/test", func(c *gin.Context) {
		id, _ := c.Get("user_id")
		email, _ := c.Get("user_email")
		c.JSON(200, gin.H{"id": id, "email": email})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer validtoken")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "test@example.com")
	assert.Contains(t, w.Body.String(), "123")
}
