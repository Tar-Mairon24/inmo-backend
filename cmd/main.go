package cmd

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"inmo-backend/internal/infrastructure/db"
	"inmo-backend/internal/interface/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables or defaults")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)

	db.Init()

	r := api.SetupRouter()

	port := os.Getenv("SERVER_PORT")

	logrus.Infof("Starting server on port %s", port)
	if err := r.Run(":" + port); err != nil {
		logrus.WithError(err).Fatal("Failed to start server")
	}
}
