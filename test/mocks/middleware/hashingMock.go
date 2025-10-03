package middlewareMock	

import (
	"github.com/stretchr/testify/mock"
)

type MockHashing struct {
	mock.Mock
}

func NewMockHashing() *MockHashing {
	return &MockHashing{}
}

func (m *MockHashing) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHashing) VerifyPassword(databasePassword string, password string) error {
	args := m.Called(databasePassword, password)
	return args.Error(0)
}