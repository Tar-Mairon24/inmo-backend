package helpers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/infrastructure/db"
	"inmo-backend/internal/interface/api"
)

// SetupTestServer creates a test server with in-memory database
func SetupTestServer() (*gin.Engine, func()) {
	gin.SetMode(gin.TestMode)

	// Create in-memory SQLite database for testing
	database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to create test database: " + err.Error())
	}

	// Auto-migrate the schema
	err = database.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to migrate test database: " + err.Error())
	}

	// Store original DB and set test DB
	originalDB := db.DB
	db.DB = database

	// Setup router with real dependencies
	router := api.SetupRouter()

	// Return cleanup function
	cleanup := func() {
		db.DB = originalDB
		sqlDB, _ := database.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
	}

	return router, cleanup
}

// CleanDatabase clears all data from test database
func CleanDatabase() {
	if db.DB == nil {
		return
	}

	// Delete all users (and other tables when you add them)
	db.DB.Exec("DELETE FROM users")

	// Reset auto-increment
	db.DB.Exec("DELETE FROM sqlite_sequence WHERE name='users'")
}
