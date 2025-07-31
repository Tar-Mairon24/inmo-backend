package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// CreateUserViaAPI creates a user through the API and returns the user ID
func CreateUserViaAPI(router *gin.Engine, username, email string) uint {
	user := map[string]interface{}{
		"username": username,
		"email":    email,
		"password": "password123",
	}

	jsonBytes, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		panic(fmt.Sprintf("Failed to create user via API. Status: %d, Body: %s", w.Code, w.Body.String()))
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Extract user ID from response data
	if data, ok := response["data"].(map[string]interface{}); ok {
		if id, ok := data["id"].(float64); ok {
			return uint(id)
		}
	}

	panic("Could not extract user ID from API response")
}

// GetUserViaAPI gets a user by ID through the API
func GetUserViaAPI(router *gin.Engine, userID uint) map[string]interface{} {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/users/%d", userID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	return response
}

// UpdateUserViaAPI updates a user through the API
func UpdateUserViaAPI(router *gin.Engine, userID uint, updates map[string]interface{}) *httptest.ResponseRecorder {
	jsonBytes, _ := json.Marshal(updates)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/v1/users/%d", userID), bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// DeleteUserViaAPI deletes a user through the API
func DeleteUserViaAPI(router *gin.Engine, userID uint) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/v1/users/%d", userID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
