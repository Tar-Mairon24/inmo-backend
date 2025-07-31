package repository_test

import (
    "testing"

    "github.com/stretchr/testify/suite"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    "inmo-backend/internal/domain/models"
    "inmo-backend/internal/infrastructure/db"
    "inmo-backend/internal/infrastructure/repository"
    "inmo-backend/test/helpers"
)

type UserRepositoryIntegrationSuite struct {
    suite.Suite
    db   *gorm.DB
    repo *repository.UserRepositoryImpl
}

func (suite *UserRepositoryIntegrationSuite) SetupSuite() {
    // Create in-memory SQLite database
    database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    suite.Require().NoError(err)

    // Auto-migrate
    err = database.AutoMigrate(&models.User{})
    suite.Require().NoError(err)

    // Set global DB
    db.DB = database
    suite.db = database
    suite.repo = repository.NewUserRepository().(*repository.UserRepositoryImpl)
}

func (suite *UserRepositoryIntegrationSuite) SetupTest() {
    // Clean database before each test
    suite.db.Exec("DELETE FROM users")
}

func (suite *UserRepositoryIntegrationSuite) TestCRUDOperations() {
    // Create
    user := helpers.CreateTestUser()
    err := suite.repo.Create(user)
    suite.NoError(err)
    suite.NotZero(user.ID)

    // Read
    retrievedUser, err := suite.repo.GetByID(user.ID)
    suite.NoError(err)
    suite.Equal(user.Username, retrievedUser.Username)

    // Update
    user.Username = "updated_user"
    err = suite.repo.Update(user)
    suite.NoError(err)

    updatedUser, err := suite.repo.GetByID(user.ID)
    suite.NoError(err)
    suite.Equal("updated_user", updatedUser.Username)

    // Delete
    err = suite.repo.Delete(user.ID)
    suite.NoError(err)

    _, err = suite.repo.GetByID(user.ID)
    suite.Error(err)
}

func (suite *UserRepositoryIntegrationSuite) TestGetAll() {
    // Create multiple users
    users := helpers.CreateTestUsers(3)
    for _, user := range users {
        err := suite.repo.Create(user)
        suite.NoError(err)
    }

    // Get all
    allUsers, err := suite.repo.GetAll()
    suite.NoError(err)
    suite.Len(allUsers, 3)
}

func TestUserRepositoryIntegration(t *testing.T) {
    suite.Run(t, new(UserRepositoryIntegrationSuite))
}