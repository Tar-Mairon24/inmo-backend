package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/interface/api/handler"
)

// MockUserUseCase is a mock implementation of UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Login(email, password string) (*models.LoginResponse, error) {
	args := m.Called(email, password)

	if LoginResp, ok := args.Get(0).(*models.LoginResponse); ok {
		return LoginResp, args.Error(1)
	}
	return nil, args.Error(1)
	
}
func (m *MockUserUseCase) GetAllUsers() ([]models.UserResponse, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserResponse), args.Error(1)
}
func (m *MockUserUseCase) GetUserByID(id uint) (*models.UserResponse, error) {
	args := m.Called(id)
	if user, ok := args.Get(0).(*models.UserResponse); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserUseCase) CreateUser(user *models.User) (*models.UserResponse, error) {
	args := m.Called(user)
	if userResp, ok := args.Get(0).(*models.UserResponse); ok {
		return userResp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserUseCase) UpdateUser(user *models.User) (*models.UserResponse, error) {
	args := m.Called(user)
	if userResp, ok := args.Get(0).(*models.UserResponse); ok {
		return userResp, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserUseCase) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetUsers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	users := []models.UserResponse{
		{ID: 1, Username: "user1", Email: "user1@example.com"},
		{ID: 2, Username: "user2", Email: "user2@example.com"},
	}
	mockUseCase.On("GetAllUsers").Return(users, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users", nil)

	handler.GetUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"Users retrieved successfully"`)
	assert.Contains(t, w.Body.String(), `"count":2`)
	assert.Contains(t, w.Body.String(), `"user1@example.com"`)
	assert.Contains(t, w.Body.String(), `"user2@example.com"`)
	mockUseCase.AssertExpectations(t)
}

func TestGetUsers_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	mockUseCase.On("GetAllUsers").Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users", nil)

	handler.GetUsers(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to retrieve users"`)
	assert.Contains(t, w.Body.String(), `"message":"database error"`)
	mockUseCase.AssertExpectations(t)
}

func TestGetUsers_NoUsersFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	mockUseCase.On("GetAllUsers").Return([]models.UserResponse{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/users", nil)

	handler.GetUsers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"No users found"`)
	assert.Contains(t, w.Body.String(), `"message":"No users available in the database"`)
	mockUseCase.AssertExpectations(t)
}
func TestGetUserByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	userResp := &models.UserResponse{
		ID:       1,
		Username: "testuser",
		Email:    "testuser@example.com",
	}
	mockUseCase.On("GetUserByID", uint(1)).Return(userResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("GET", "/users/1", nil)

	handler.GetUserByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"User retrieved successfully"`)
	assert.Contains(t, w.Body.String(), `"testuser@example.com"`)
	mockUseCase.AssertExpectations(t)
}

func TestGetUserByID_InvalidIDFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	c.Request, _ = http.NewRequest("GET", "/users/abc", nil)

	handler.GetUserByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid user ID"`)
}
func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	userData := `{"username":"newuser","email":"newuser@example.com","password":"securepass"}`
	user := &models.User{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "securepass",
	}
	userResp := &models.UserResponse{
		ID:       1,
		Username: "newuser",
		Email:    "newuser@example.com",
	}
	mockUseCase.On("CreateUser", user).Return(userResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/users", bytes.NewBufferString(userData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"User created successfully"`)
	assert.Contains(t, w.Body.String(), `"newuser@example.com"`)
	mockUseCase.AssertExpectations(t)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	invalidJSON := `{"username":"baduser","email":"baduser@example.com","password":}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/users", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Failed to parse user data"`)
}

func TestCreateUser_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	userData := `{"username":"newuser","email":"newuser@example.com","password":"securepass"}`
	user := &models.User{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "securepass",
	}
	mockUseCase.On("CreateUser", user).Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/users", bytes.NewBufferString(userData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to create user"`)
	assert.Contains(t, w.Body.String(), `"message":"database error"`)
	mockUseCase.AssertExpectations(t)
}
func TestUpdateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	userData := `{"id":1,"username":"updateduser","email":"updated@example.com","password":"newpass"}`
	user := &models.User{
		ID:       1,
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "newpass",
	}
	userResp := &models.UserResponse{
		ID:       1,
		Username: "updateduser",
		Email:    "updated@example.com",
	}
	mockUseCase.On("UpdateUser", user).Return(userResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/users/1", bytes.NewBufferString(userData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"message":"User updated successfully"`)
	assert.Contains(t, w.Body.String(), `"updated@example.com"`)
	mockUseCase.AssertExpectations(t)
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	invalidJSON := `{"id":1,"username":"baduser","email":"baduser@example.com","password":}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/users/1", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid request"`)
	assert.Contains(t, w.Body.String(), `"message":"Failed to parse user data"`)
}

func TestUpdateUser_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	userData := `{"id":1,"username":"updateduser","email":"updated@example.com","password":"newpass"}`
	user := &models.User{
		ID:       1,
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "newpass",
	}
	mockUseCase.On("UpdateUser", user).Return(nil, errors.New("update failed"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/users/1", bytes.NewBufferString(userData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to update user"`)
	assert.Contains(t, w.Body.String(), `"message":"update failed"`)
	mockUseCase.AssertExpectations(t)
}
func TestDeleteUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	mockUseCase.On("DeleteUser", uint(1)).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("DELETE", "/users/1", nil)

	handler.DeleteUser(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Body.String())
	mockUseCase.AssertExpectations(t)
}

func TestDeleteUser_InvalidIDFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	c.Request, _ = http.NewRequest("DELETE", "/users/abc", nil)

	handler.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Invalid user ID"`)
	assert.Contains(t, w.Body.String(), `"message":"User ID must be a valid number"`)
}

func TestDeleteUser_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUseCase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUseCase)

	mockUseCase.On("DeleteUser", uint(2)).Return(errors.New("delete failed"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	c.Request, _ = http.NewRequest("DELETE", "/users/2", nil)

	handler.DeleteUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"Failed to delete user"`)
	assert.Contains(t, w.Body.String(), `"message":"delete failed"`)
	mockUseCase.AssertExpectations(t)
}




