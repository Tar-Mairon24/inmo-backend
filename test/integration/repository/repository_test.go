package repository_test

import (
	"log"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"inmo-backend/internal/infrastructure/db"
	"inmo-backend/internal/infrastructure/repository"
	"inmo-backend/test/helpers"
)

type UserRepositoryIntegrationSuite struct {
    suite.Suite
    db   *gorm.DB
    repo *repository.UserRepositoryImpl
    router *gin.Engine
    cleanup func()
}

func (suite *UserRepositoryIntegrationSuite) SetupSuite() {
	suite.router, suite.cleanup = helpers.SetupTestServer()
    suite.repo = repository.NewUserRepository().(*repository.UserRepositoryImpl)
    suite.db = db.DB 
    log.Println("Setting up UserRepositoryIntegrationSuite")
}

func (suite *UserRepositoryIntegrationSuite) TearDownSuite() {
	suite.cleanup()
}

func (suite *UserRepositoryIntegrationSuite) SetupTest() {
	helpers.CleanDatabase()
    logrus.Info("Cleaning database before each test")
}

func (suite *UserRepositoryIntegrationSuite) TestCRUDOperations() {
    user := helpers.CreateTestUser()
    err := suite.repo.Create(user)
    suite.NoError(err)
    suite.NotZero(user.ID)

    retrievedUser, err := suite.repo.GetByID(user.ID)
    suite.NoError(err)
    suite.Equal(user.Username, retrievedUser.Username)

    user.Username = "updated_user"
    err = suite.repo.Update(user)
    suite.NoError(err)

    updatedUser, err := suite.repo.GetByID(user.ID)
    suite.NoError(err)
    suite.Equal("updated_user", updatedUser.Username)

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