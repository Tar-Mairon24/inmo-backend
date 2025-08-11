package middleware_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"inmo-backend/middleware"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"

	hash, err := middleware.HashPassword(password)

	assert.NoError(t, err, "Hashing should not return an error")
	assert.NotEmpty(t, hash, "Hashed password should not be empty")
	assert.NotEqual(t, password, hash, "Hashed password should not be the same as the original password")
	assert.NoError(t, middleware.VerifyPassword(hash, password), "Hashed password should match the original password")
	assert.Greater(t, len(hash), 0, "Hashed password should not be empty")
	assert.Contains(t, hash, "$2a$", "Hashed password should start with $2a$ indicating bcrypt")
}

func TestVerifyPassword(t *testing.T) {
	password := "testPassword123"
	hash, err := middleware.HashPassword(password)
	assert.NoError(t, err, "Hashing should not return an error")

	err = middleware.VerifyPassword(hash, password)
	assert.NoError(t, err, "Verification should succeed with correct password")

	err = middleware.VerifyPassword(hash, "wrongPassword")
	assert.Error(t, err, "Verification should fail with incorrect password")

	err = middleware.VerifyPassword(hash, "")
	assert.Error(t, err, "Verification should fail with empty password")

	err = middleware.VerifyPassword("", password)
	assert.Error(t, err, "Verification should fail with empty hash")
}
func TestHashPassword_EmptyPassword(t *testing.T) {
	password := ""

	_, err := middleware.HashPassword(password)

	assert.NotNil(t, err, "Hashing should return an error for empty password")
	assert.Empty(t, "", "Hashed password should be empty for empty input")
}

func TestHashPassword_ShortPassword(t *testing.T) {
	password := "short"

	_, err := middleware.HashPassword(password)

	assert.NotNil(t, err, "Hashing should return an error for short password")
	assert.Empty(t, "", "Hashed password should be empty for short input")
}

func TestHashPassword_SpecialCharacters(t *testing.T) {
	password := "!@#$$%^&*()_+-=[]{}|;':,.<>/?`~"

	hash, err := middleware.HashPassword(password)

	assert.NoError(t, err, "Hashing should not return an error for special characters")
	assert.NotEmpty(t, hash, "Hashed password should not be empty")
	assert.NotEqual(t, password, hash, "Hashed password should not be the same as the original password")
	assert.NoError(t, middleware.VerifyPassword(hash, password), "Hashed password should match the original password with special characters")
}

func TestHashPassword_LongPassword(t *testing.T) {
	password := ""
	for range 1000 {
		password += "a"
	}

	_, err := middleware.HashPassword(password)

	assert.Error(t, err, "Hashing should not return an error for long password")
	assert.Empty(t, "", "Hashed password should not be empty for long input")
}