package usecaseMocks

import (
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
)

type mockPropertyUseCase struct {
	mock.Mock
}

func NewMockPropertyUseCase() *mockPropertyUseCase {
	return &mockPropertyUseCase{}
}

func (m *mockPropertyUseCase) GetAllProperties() ([]models.PropertyResponse, error) {
	args := m.Called()
	return args.Get(0).([]models.PropertyResponse), args.Error(1)
}
func (m *mockPropertyUseCase) GetPropertyByID(id uint) (*models.PropertyResponse, error) {
	args := m.Called(id)
	return args.Get(0).(*models.PropertyResponse), args.Error(1)
}
func (m *mockPropertyUseCase) CreateProperty(p *models.Property) (*models.PropertyResponse, error) {
	args := m.Called(p)
	return args.Get(0).(*models.PropertyResponse), args.Error(1)
}
func (m *mockPropertyUseCase) UpdateProperty(p *models.Property) (*models.PropertyResponse, error) {
	args := m.Called(p)
	return args.Get(0).(*models.PropertyResponse), args.Error(1)
}
func (m *mockPropertyUseCase) DeleteProperty(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
