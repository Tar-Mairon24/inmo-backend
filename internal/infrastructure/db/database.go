package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"inmo-backend/internal/domain/models"
)

var (
	DB *gorm.DB
	SqlDB *sql.DB
)

func GetSqlDB() *sql.DB{
	logrus.Info("Getting SQL DB connection")
	if SqlDB == nil {
		logrus.Error("SQL DB connection is not initialized")
	}
	return SqlDB
}

func GetDB() *gorm.DB {
	logrus.Info("Getting GORM DB connection")
	if DB == nil {
		logrus.Error("GORM DB connection is not initialized")
	}
	return DB
}

func Init() {	
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	logrus.Infof("Connecting to database at %s:%s/%s", dbHost, dbPort, dbName)

	// Open GORM connection
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}
	DB = database
	logrus.Info("Successfully connected to database")

	// Open Squirrel SQL connection
	SqlDB, err = database.DB()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to get underlying sql.DB from GORM")
	}

	if SqlDB == nil {
		logrus.Error("SQL DB connection is nil")
	}
	logrus.Info("Successfully obtained SQL DB connection")

	err = DB.AutoMigrate(&models.User{}, &models.Property{})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to auto-migrate database")
	} else {
		logrus.Info("Database auto-migration completed successfully")
	}
	logrus.Info("Database initialized successfully")
}
