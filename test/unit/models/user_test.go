package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"inmo-backend/internal/domain/models"
)

func TestToUserResponse_NilUser(t *testing.T) {
	var user *models.User = nil
	resp := user.ToUserResponse()
	assert.Nil(t, resp)
}

func TestToUserResponse_ValidUser(t *testing.T) {
	now := time.Now()
	user := &models.User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}
	resp := user.ToUserResponse()
	assert.NotNil(t, resp)
	assert.Equal(t, user.ID, resp.ID)
	assert.Equal(t, user.Username, resp.Username)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.CreatedAt, resp.CreatedAt)
	assert.Equal(t, user.UpdatedAt, resp.UpdatedAt)
}
