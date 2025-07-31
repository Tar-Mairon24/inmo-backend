package fixtures

import "inmo-backend/internal/domain/models"

var TestUsers = []models.User{
    {
        ID:       1,
        Username: "john_doe",
        Email:    "john@example.com",
        Password: "password123",
    },
    {
        ID:       2,
        Username: "jane_smith",
        Email:    "jane@example.com",
        Password: "password456",
    },
}

var InvalidUser = models.User{
    Username: "", // Invalid: empty username
    Email:    "invalid-email", // Invalid: bad email format
    Password: "123", // Invalid: too short
}