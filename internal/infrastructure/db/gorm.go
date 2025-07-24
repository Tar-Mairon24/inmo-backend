package db

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"inmo-backend/internal/domain/models"
)

var DB *gorm.DB

func Init() {
	// Get database configuration from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Build DSN from environment variables
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	logrus.Infof("Connecting to database at %s:%s/%s", dbHost, dbPort, dbName)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}

	DB = database
	logrus.Info("Successfully connected to database")

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to auto-migrate database")
	} else {
		logrus.Info("Database auto-migration completed successfully")
	}
	logrus.Info("Database initialized successfully")
}
