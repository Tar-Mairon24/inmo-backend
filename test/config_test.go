package test

import (
    "os"

    "github.com/joho/godotenv"
)

func init() {
    // Load test environment
    if err := godotenv.Load("../.env.test"); err != nil {
        // Use default test values if no .env.test file
        os.Setenv("DB_HOST", "localhost")
        os.Setenv("DB_PORT", "3306")
        os.Setenv("DB_USER", "test_user")
        os.Setenv("DB_PASSWORD", "test_password")
        os.Setenv("DB_NAME", "test_db")
        os.Setenv("GIN_MODE", "test")
    }
}