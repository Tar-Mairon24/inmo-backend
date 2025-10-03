package repositoryMock

import (
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
)

type MockPropertyRepository struct {
	mock.Mock
}

func NewMockPropertyRepository() *MockPropertyRepository {
	return &MockPropertyRepository{}
}

func (m *MockPropertyRepository) GetAll() ([]models.PropertyResponse, error) {
	args := m.Called()
	if properties, ok := args.Get(0).([]models.PropertyResponse); ok {
		return properties, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyRepository) GetByID(id uint) (*models.PropertyResponse, error) {
	args := m.Called(id)
	if property, ok := args.Get(0).(*models.PropertyResponse); ok {
		return property, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyRepository) Create(property *models.Property) (*models.PropertyResponse, error) {
	args := m.Called(property)
	if propertyResponse, ok := args.Get(0).(*models.PropertyResponse); ok {
		return propertyResponse, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyRepository) Update(property *models.Property) (*models.PropertyResponse, error) {
	args := m.Called(property)
	if propertyResponse, ok := args.Get(0).(*models.PropertyResponse); ok {
		return propertyResponse, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPropertyRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
