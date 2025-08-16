package usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/usecase"
)

// MockPropertyRepository implements ports.PropertyRepository for testing
type MockPropertyRepository struct {
	mock.Mock
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



func TestPropertyUseCase_GetAllProperties(t *testing.T) {
	t.Run("should return all properties successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		expectedProperties := []models.PropertyResponse{
			{ID: 1, Address: "123 Main St", Price: 100000},
			{ID: 2, Address: "456 Oak Ave", Price: 200000},
		}
		
		mockRepo.On("GetAll").Return(expectedProperties, nil)
		
		// Act
		result, err := propertyUseCase.GetAllProperties()
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProperties, result)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return empty slice when no properties exist", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		expectedProperties := []models.PropertyResponse{}
		
		mockRepo.On("GetAll").Return(expectedProperties, nil)
		
		// Act
		result, err := propertyUseCase.GetAllProperties()
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProperties, result)
		assert.Len(t, result, 0)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		expectedError := errors.New("database connection failed")
		
		mockRepo.On("GetAll").Return([]models.PropertyResponse(nil), expectedError)
		
		// Act
		result, err := propertyUseCase.GetAllProperties()
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
func TestPropertyUseCase_GetPropertyByID(t *testing.T) {
	t.Run("should return property successfully when found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		expectedProperty := &models.PropertyResponse{
			ID:      1,
			Address: "123 Main St",
			Price:   100000,
		}
		
		mockRepo.On("GetByID", uint(1)).Return(expectedProperty, nil)
		
		// Act
		result, err := propertyUseCase.GetPropertyByID(1)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedProperty, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when property not found", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		mockRepo.On("GetByID", uint(999)).Return((*models.PropertyResponse)(nil), nil)
		
		// Act
		result, err := propertyUseCase.GetPropertyByID(999)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "property not found", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		expectedError := errors.New("database connection failed")
		
		mockRepo.On("GetByID", uint(1)).Return((*models.PropertyResponse)(nil), expectedError)
		
		// Act
		result, err := propertyUseCase.GetPropertyByID(1)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
func TestPropertyUseCase_CreateProperty(t *testing.T) {
	t.Run("should create property successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			Address: "123 Main St",
			Price:   100000,
		}
		
		expectedResponse := &models.PropertyResponse{
			ID:      1,
			Address: "123 Main St",
			Price:   100000,
		}
		
		mockRepo.On("Create", inputProperty).Return(expectedResponse, nil)
		
		// Act
		result, err := propertyUseCase.CreateProperty(inputProperty)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when property is nil", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		// Act
		result, err := propertyUseCase.CreateProperty(nil)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "property cannot be nil", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error when address is empty", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			Address: "",
			Price:   100000,
		}
		
		// Act
		result, err := propertyUseCase.CreateProperty(inputProperty)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "address cannot be empty", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error when price is zero", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			Address: "123 Main St",
			Price:   0,
		}
		
		// Act
		result, err := propertyUseCase.CreateProperty(inputProperty)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "price must be greater than zero", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error when price is negative", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			Address: "123 Main St",
			Price:   -1000,
		}
		
		// Act
		result, err := propertyUseCase.CreateProperty(inputProperty)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "price must be greater than zero", err.Error())
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			Address: "123 Main St",
			Price:   100000,
		}
		
		expectedError := errors.New("database connection failed")
		
		mockRepo.On("Create", inputProperty).Return((*models.PropertyResponse)(nil), expectedError)
		
		// Act
		result, err := propertyUseCase.CreateProperty(inputProperty)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
func TestPropertyUseCase_UpdateProperty(t *testing.T) {
	t.Run("should update property successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			ID:      1,
			Address: "123 Updated St",
			Price:   150000,
		}
		
		expectedResponse := &models.PropertyResponse{
			ID:      1,
			Address: "123 Updated St",
			Price:   150000,
		}
		
		mockRepo.On("Update", inputProperty).Return(expectedResponse, nil)
		
		// Act
		result, err := propertyUseCase.UpdateProperty(inputProperty)
		
		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when property is nil", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		// Act
		result, err := propertyUseCase.UpdateProperty(nil)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "property cannot be nil", err.Error())
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("should return error when property ID is zero", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			ID:      0,
			Address: "123 Main St",
			Price:   100000,
		}
		
		// Act
		result, err := propertyUseCase.UpdateProperty(inputProperty)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "property ID must be provided", err.Error())
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("should return error when address is empty", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			ID:      1,
			Address: "",
			Price:   100000,
		}
		
		// Act
		result, err := propertyUseCase.UpdateProperty(inputProperty)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "address cannot be empty", err.Error())
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("should return error when price is zero", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			ID:      1,
			Address: "123 Main St",
			Price:   0,
		}
		
		// Act
		result, err := propertyUseCase.UpdateProperty(inputProperty)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "price must be greater than zero", err.Error())
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("should return error when price is negative", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			ID:      1,
			Address: "123 Main St",
			Price:   -1000,
		}
		
		// Act
		result, err := propertyUseCase.UpdateProperty(inputProperty)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "price must be greater than zero", err.Error())
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		inputProperty := &models.Property{
			ID:      1,
			Address: "123 Main St",
			Price:   100000,
		}
		
		expectedError := errors.New("database connection failed")
		
		mockRepo.On("Update", inputProperty).Return((*models.PropertyResponse)(nil), expectedError)
		
		// Act
		result, err := propertyUseCase.UpdateProperty(inputProperty)
		
		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
func TestPropertyUseCase_DeleteProperty(t *testing.T) {
	t.Run("should delete property successfully", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		mockRepo.On("Delete", uint(1)).Return(nil)
		
		// Act
		err := propertyUseCase.DeleteProperty(1)
		
		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when property ID is zero", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		// Act
		err := propertyUseCase.DeleteProperty(0)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, "property ID must be provided", err.Error())
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Arrange
		mockRepo := new(MockPropertyRepository)
		propertyUseCase := usecase.NewPropertyUseCase(mockRepo)
		
		expectedError := errors.New("database connection failed")
		
		mockRepo.On("Delete", uint(1)).Return(expectedError)
		
		// Act
		err := propertyUseCase.DeleteProperty(1)
		
		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}




