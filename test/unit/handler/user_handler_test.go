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

func (m *MockUserUseCase) Login(email, password string) error {
	args := m.Called(email, password)
	return args.Error(0)
}
func (m *MockUserUseCase) GetAllUsers() ([]models.UserResponse, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.UserResponse), args.Error(1)
}
func (m *MockUserUseCase) GetUserByID(id uint) (*models.UserResponse, error) { 
	args:= m.Called(id)
	if user, ok := args.Get(0).(*models.UserResponse); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockUserUseCase) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUserUseCase) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUserUseCase) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	loginData := `{"email":"test@example.com","password":"password123"}`
	mockUsecase.On("Login", "test@example.com", "password123").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(loginData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UserLogin(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login successful")
	mockUsecase.AssertExpectations(t)
}

func TestUserLogin_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	invalidJSON := `{"email": "test@example.com", "password": 123}` // password should be string

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UserLogin(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to parse login data")
}

func TestUserLogin_LoginFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	loginData := `{"email":"test@example.com","password":"wrongpass"}`
	mockUsecase.On("Login", "test@example.com", "wrongpass").Return(errors.New("invalid credentials"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(loginData))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UserLogin(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid email or password")
	mockUsecase.AssertExpectations(t)
}
func TestGetUsers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	users := []models.UserResponse{
		{ID: 1, Email: "user1@example.com"},
		{ID: 2, Email: "user2@example.com"},
	}
	mockUsecase.On("GetAllUsers").Return(users, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/users", nil)

	handler.GetUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Users retrieved successfully")
	assert.Contains(t, w.Body.String(), "user1@example.com")
	assert.Contains(t, w.Body.String(), "user2@example.com")
	mockUsecase.AssertExpectations(t)
}

func TestGetUsers_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	mockUsecase.On("GetAllUsers").Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/api/v1/users", nil)

	handler.GetUsers(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to retrieve users")
	assert.Contains(t, w.Body.String(), "database error")
	mockUsecase.AssertExpectations(t)
}
func TestGetUserByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	userResp := &models.UserResponse{ID: 1, Email: "user1@example.com"}
	mockUsecase.On("GetUserByID", uint(1)).Return(userResp, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("GET", "/api/v1/users/1", nil)

	handler.GetUserByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User retrieved successfully")
	assert.Contains(t, w.Body.String(), "user1@example.com")
	mockUsecase.AssertExpectations(t)
}

func TestGetUserByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	c.Request, _ = http.NewRequest("GET", "/api/v1/users/abc", nil)

	handler.GetUserByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid user ID")
	assert.Contains(t, w.Body.String(), "User ID must be a valid number")
}

func TestGetUserByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	mockUsecase.On("GetUserByID", uint(99)).Return((*models.UserResponse)(nil), errors.New("user not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "99"}}
	c.Request, _ = http.NewRequest("GET", "/api/v1/users/99", nil)

	handler.GetUserByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
	assert.Contains(t, w.Body.String(), "user not found")
	mockUsecase.AssertExpectations(t)
}
func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	user := models.User{
		Email:    "newuser@example.com",
		Password: "securepassword",
	}
	userJSON := `{"email":"newuser@example.com","password":"securepassword"}`

	mockUsecase.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
		argUser := args.Get(0).(*models.User)
		assert.Equal(t, user.Email, argUser.Email)
		assert.Equal(t, user.Password, argUser.Password)
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(userJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")
	assert.Contains(t, w.Body.String(), "newuser@example.com")
	mockUsecase.AssertExpectations(t)
}

func TestCreateUser_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	invalidJSON := `{"email": "baduser@example.com", "password": 123}` // password should be string

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to parse user data")
}

func TestCreateUser_CreateUserFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	userJSON := `{"email":"failuser@example.com","password":"failpass"}`
	mockUsecase.On("CreateUser", mock.AnythingOfType("*models.User")).Return(errors.New("db error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(userJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create user")
	assert.Contains(t, w.Body.String(), "db error")
	mockUsecase.AssertExpectations(t)
}
func TestUpdateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	user := models.User{
		ID:       1,
		Email:    "updateuser@example.com",
		Password: "newpassword",
	}
	userJSON := `{"id":1,"email":"updateuser@example.com","password":"newpassword"}`

	mockUsecase.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
		argUser := args.Get(0).(*models.User)
		assert.Equal(t, user.ID, argUser.ID)
		assert.Equal(t, user.Email, argUser.Email)
		assert.Equal(t, user.Password, argUser.Password)
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBufferString(userJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User updated successfully")
	assert.Contains(t, w.Body.String(), "updateuser@example.com")
	mockUsecase.AssertExpectations(t)
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	invalidJSON := `{"id":1,"email":123,"password":"pass"}` // email should be string

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/v1/users/1", bytes.NewBufferString(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to parse user data")
}

func TestUpdateUser_UpdateFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	userJSON := `{"id":2,"email":"failupdate@example.com","password":"failpass"}`
	mockUsecase.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(errors.New("update error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("PUT", "/api/v1/users/2", bytes.NewBufferString(userJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to update user")
	assert.Contains(t, w.Body.String(), "update error")
	mockUsecase.AssertExpectations(t)
}
func TestDeleteUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	mockUsecase.On("DeleteUser", uint(1)).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request, _ = http.NewRequest("DELETE", "/api/v1/users/1", nil)

	handler.DeleteUser(c)

	t.Logf("Actual Status Code: %d", w.Code)
    t.Logf("Response Body: '%s'", w.Body.String())
    t.Logf("Response Headers: %v", w.Header())
	t.Logf("🔍 Mock calls: %v", mockUsecase.Calls)  // Debug mock calls
    t.Logf("🔍 Actual Status Code: %d", w.Code)
    t.Logf("🔍 Response Body: '%s'", w.Body.String())

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
	mockUsecase.AssertExpectations(t)
}

func TestDeleteUser_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	c.Request, _ = http.NewRequest("DELETE", "/api/v1/users/abc", nil)

	handler.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid user ID")
	assert.Contains(t, w.Body.String(), "User ID must be a valid number")
}

func TestDeleteUser_DeleteFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockUsecase := new(MockUserUseCase)
	handler := handler.NewUserHandler(mockUsecase)

	mockUsecase.On("DeleteUser", uint(2)).Return(errors.New("delete error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	c.Request, _ = http.NewRequest("DELETE", "/api/v1/users/2", nil)

	handler.DeleteUser(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to delete user")
	assert.Contains(t, w.Body.String(), "delete error")
	mockUsecase.AssertExpectations(t)
}




