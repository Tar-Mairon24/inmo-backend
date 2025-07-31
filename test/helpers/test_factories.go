package helpers

import (
	"fmt"
	"time"

	"inmo-backend/internal/domain/models"
)

// CreateTestUser creates a test user with default values
func CreateTestUser() *models.User {
	return &models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestUserWithData creates a test user with custom data
func CreateTestUserWithData(username, email string) *models.User {
	return &models.User{
		Username:  username,
		Email:     email,
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestUsers creates multiple test users
func CreateTestUsers(count int) []*models.User {
	users := make([]*models.User, count)
	for i := 0; i < count; i++ {
		users[i] = &models.User{
			Username:  fmt.Sprintf("user%d", i+1),
			Email:     fmt.Sprintf("user%d@example.com", i+1),
			Password:  "password123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	return users
}

// CreateInvalidUser creates a user with invalid data for testing validation
func CreateInvalidUser() *models.User {
	return &models.User{
		Username: "",              // Invalid: empty username
		Email:    "invalid-email", // Invalid: bad email format
		Password: "123",           // Invalid: too short
	}
}
