package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"inmo-backend/test/helpers"
)

type UserAPIE2ESuite struct {
	suite.Suite
	router  *gin.Engine
	cleanup func()
}

func (suite *UserAPIE2ESuite) SetupSuite() {
	suite.router, suite.cleanup = helpers.SetupTestServer()
}

func (suite *UserAPIE2ESuite) TearDownSuite() {
	suite.cleanup()
}

func (suite *UserAPIE2ESuite) SetupTest() {
	helpers.CleanDatabase()
}

func (suite *UserAPIE2ESuite) TestGetUsers() {
	userID := helpers.CreateUserViaAPI(suite.router, "testuser", "test@example.com")
	assert.NotZero(suite.T(), userID)

	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	assert.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		suite.T().Logf("Expected status 200, got %d. Response: %s", w.Code, w.Body.String())
	}

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Users retrieved successfully", response["message"])
	assert.Equal(suite.T(), float64(1), response["count"])
}

func (suite *UserAPIE2ESuite) TestCreateUser() {
	user := map[string]any{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}

	jsonBytes, err := json.Marshal(user)
	assert.NoError(suite.T(), err)

	req, err := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBytes))
	assert.NoError(suite.T(), err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		suite.T().Logf("Expected status 201, got %d. Response: %s", w.Code, w.Body.String())
	}

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "User created successfully", response["message"])

	assert.Contains(suite.T(), response, "data")
}

func (suite *UserAPIE2ESuite) TestDeleteUser() {
	userID := helpers.CreateUserViaAPI(suite.router, "testuser", "test@example.com")
	assert.NotZero(suite.T(), userID)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%d", userID), nil)
	assert.NoError(suite.T(), err)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		suite.T().Logf("Expected status 204, got %d. Response: %s", w.Code, w.Body.String())
	}

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)

	req, err = http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", userID), nil)
	assert.NoError(suite.T(), err)

	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func TestUserAPIE2E(t *testing.T) {
	suite.Run(t, new(UserAPIE2ESuite))
}
