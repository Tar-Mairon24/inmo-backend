package di

import (
	"database/sql"
	
	"inmo-backend/internal/infrastructure/db"
	"inmo-backend/internal/infrastructure/repository"
	"inmo-backend/internal/domain/ports"
	"inmo-backend/internal/usecase"
	"inmo-backend/internal/interface/api/handler"

	"github.com/sirupsen/logrus"
)

type Container struct {
	SqlDB      		*sql.DB
	userRepo   		ports.UserRepository
	userUsecase 	*usecase.UserUseCase
	userHandler 	*handler.UserHandler
	healthHandler 	*handler.HealthHandler
}

func NewContainer() *Container {
	logrus.Info("Initializing DI container")

	container := &Container{}

	db.Init()
	container.SqlDB = db.GetSqlDB()

	if container.SqlDB == nil {
		logrus.Fatal("Failed to initialize database connection")
	}

	container.userRepo = repository.NewUserRepository(container.SqlDB)
	container.userUsecase = usecase.NewUserUseCase(container.userRepo)
	container.userHandler = handler.NewUserHandler(container.userUsecase)
	container.healthHandler = handler.NewHealthHandler()

	logrus.Info("DI container initialized successfully")
	return container
}

type Handlers struct {
	UserHandler   *handler.UserHandler
	HealthHandler *handler.HealthHandler
}

func (c *Container) GetHandlers() *Handlers {
	return &Handlers{
		UserHandler:  c.userHandler,
		HealthHandler: c.healthHandler,
	}
}