package cmd

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/infrastructure/db"
	"inmo-backend/internal/interface/api"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables or defaults")
	}

	// Configure logging
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

	// Initialize database
	db.Init()

	// Setup router
	r := api.SetupRouter()

	// Get server port from environment variable
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
}
