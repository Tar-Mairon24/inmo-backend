package helpers

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"inmo-backend/internal/infrastructure/db"
)

// SetupMockDB creates a mock database for testing
func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// Create GORM DB with mock
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	// Store original DB and set mock DB
	originalDB := db.DB
	db.DB = gormDB

	// Return cleanup function
	cleanup := func() {
		db.DB = originalDB
		mockDB.Close()
	}

	return gormDB, mock, cleanup
}
