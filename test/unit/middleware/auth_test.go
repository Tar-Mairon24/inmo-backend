package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"inmo-backend/internal/domain/models"
	"inmo-backend/middleware"
	"inmo-backend/test/mocks/service"
)

// MockClaims represents mock JWT claims
type MockClaims struct {
	ID    string
	Email string
}

func TestJWTAuthMiddleware_MissingAuthHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockService := serviceMocks.NewMockJWTService()
	authMiddleware := middleware.NewAuthMiddleware()
	mockService.On("ValidateToken", "").Return(nil, errors.New("Missing or invalid token"))
	r.Use(authMiddleware.JWTAuthMiddleware(mockService))
	r.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing or invalid token")
}

func TestJWTAuthMiddleware_InvalidToken(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    mockService := serviceMocks.NewMockJWTService()
    
    // Mock the ValidateToken call to return an error
    mockService.On("ValidateToken", "invalidtoken").Return(nil, errors.New("Invalid token"))
    
	middleware := middleware.NewAuthMiddleware()
    // Use the actual middleware
    r.Use(middleware.JWTAuthMiddleware(mockService))
    r.GET("/test", func(c *gin.Context) {
        c.String(200, "ok")
    })

    req, _ := http.NewRequest("GET", "/test", nil)
    req.AddCookie(&http.Cookie{
        Name:  "jwt_token",
        Value: "invalidtoken",
    })
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
    assert.Contains(t, w.Body.String(), "Invalid or expired token")
    mockService.AssertExpectations(t)
}

func TestJWTAuthMiddleware_ValidToken(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    mockService := serviceMocks.NewMockJWTService()
    
    // Mock the ValidateToken call to return valid claims
    mockClaims := &models.JWTClaims{
        ID:    1,
        Email: "test@example.com",
    }
    mockService.On("ValidateToken", "validtoken").Return(mockClaims, nil)
    
	middleware := middleware.NewAuthMiddleware()
    // Use the actual middleware
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
    req.AddCookie(&http.Cookie{
        Name:  "jwt_token",
        Value: "validtoken",
    })
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.Equal(t, "ok", w.Body.String())
    mockService.AssertExpectations(t)
}
func TestGenerateRefreshToken_Success(t *testing.T) {
    m := middleware.NewAuthMiddleware()
    token, id, err := m.GenerateRefreshToken()
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    assert.NotEmpty(t, id)
    assert.NotEqual(t, token, id)
}

func TestGenerateRefreshToken_Uniqueness(t *testing.T) {
    m := middleware.NewAuthMiddleware()
    token1, id1, err1 := m.GenerateRefreshToken()
    token2, id2, err2 := m.GenerateRefreshToken()
    assert.NoError(t, err1)
    assert.NoError(t, err2)
    assert.NotEqual(t, token1, token2)
    assert.NotEqual(t, id1, id2)
}

func TestGenerateRefreshToken_TokenFormat(t *testing.T) {
    m := middleware.NewAuthMiddleware()
    token, id, err := m.GenerateRefreshToken()
    assert.NoError(t, err)
    // Token should be base64 encoded, length should be greater than 32
    assert.GreaterOrEqual(t, len(token), 44)
    // ID should be hex encoded, length should be 64 (sha256)
    assert.Equal(t, 64, len(id))
}
