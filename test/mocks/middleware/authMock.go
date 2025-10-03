package middlewareMock

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/ports"
)

type MockMiddleware struct {
	mock.Mock
}

func NewMockMiddleware() *MockMiddleware {
	return &MockMiddleware{}
}

func (m *MockMiddleware) JWTAuthMiddleware(jwtService ports.JWTService) gin.HandlerFunc {
	args := m.Called(jwtService)
	return args.Get(0).(gin.HandlerFunc)
}

func (m *MockMiddleware) GenerateRefreshToken() (string, string, error) {
	args := m.Called()
	return args.String(0), args.String(1), args.Error(2)
}
