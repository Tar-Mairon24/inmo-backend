package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/interface/api/handler"
)

// Mock for PropertyUseCase
type mockPropertyUseCase struct {
	mock.Mock
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

func TestGetProperties_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	properties := []models.PropertyResponse{
		{ID: 1, Title: "Prop1"},
		{ID: 2, Title: "Prop2"},
	}
	mockUC.On("GetAllProperties").Return(properties, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetProperties(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetProperties_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	mockUC.On("GetAllProperties").Return([]models.PropertyResponse{}, errors.New("db error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetProperties(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetProperties_EmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	mockUC.On("GetAllProperties").Return([]models.PropertyResponse{}, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetProperties(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetPropertyByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	expected := &models.PropertyResponse{ID: 1, Title: "Prop1"}
	mockUC.On("GetPropertyByID", uint(1)).Return(expected, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetPropertyByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPropertyByID_NegativeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "-5"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetPropertyByID_ErrorFromUsecase(t *testing.T) {
	var property *models.PropertyResponse = nil

	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	mockUC.On("GetPropertyByID", uint(2)).Return(property, errors.New("db error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetPropertyByID_NotFound(t *testing.T) {
	var property *models.PropertyResponse = nil
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	mockUC.On("GetPropertyByID", uint(3)).Return(property, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "3"}}

	h.GetPropertyByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}
func TestCreateProperty_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	expected := &models.PropertyResponse{ID: 1, Title: "New Property"}
	mockUC.On("CreateProperty", mock.AnythingOfType("*models.Property")).Return(expected, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"New Property"}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateProperty_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{invalid json}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateProperty_UsecaseError(t *testing.T) {
	var property *models.PropertyResponse = nil
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	mockUC.On("CreateProperty", mock.AnythingOfType("*models.Property")).Return(property, errors.New("db error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/properties", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"title":"New Property"}`))

	h.CreateProperty(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}
func TestUpdateProperty_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	expected := &models.PropertyResponse{ID: 1, Title: "Updated Property"}
	mockUC.On("UpdateProperty", mock.AnythingOfType("*models.Property")).Return(expected, nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("PUT", "/properties/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"id":1,"title":"Updated Property"}`))

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestUpdateProperty_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("PUT", "/properties/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{invalid json}`))

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateProperty_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	expected := &models.PropertyResponse{ID: 1, Title: "Updated Property"}
	mockUC.On("UpdateProperty", mock.AnythingOfType("*models.Property")).Return(expected, errors.New("db error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("PUT", "/properties/1", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("Content-Type", "application/json")
	c.Request.Body = io.NopCloser(strings.NewReader(`{"id":1,"title":"Updated Property"}`))

	h.UpdateProperty(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}
func TestDeleteProperty_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	mockUC.On("DeleteProperty", uint(1)).Return(nil)

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DeleteProperty(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeleteProperty_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h.DeleteProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteProperty_NegativeID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "-10"}}

	h.DeleteProperty(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteProperty_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUC := new(mockPropertyUseCase)
	mockUC.On("DeleteProperty", uint(2)).Return(errors.New("db error"))

	h := handler.NewPropertyHandler(mockUC)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}

	h.DeleteProperty(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}


