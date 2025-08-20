package di

import (
	"database/sql"

	"github.com/sirupsen/logrus"

	"inmo-backend/internal/domain/ports"
	"inmo-backend/internal/infrastructure/db"
	"inmo-backend/internal/infrastructure/repository"
	"inmo-backend/internal/infrastructure/service"
	"inmo-backend/internal/interface/api/handler"
	"inmo-backend/internal/usecase"
)

type Container struct {
	SqlDB      			*sql.DB
	userRepo   			ports.UserRepository
	propertyRepo    	ports.PropertyRepository
	userUsecase 		ports.UserUseCase
	propertyUsecase  	ports.PropertyUseCase
	jwtService 			ports.JWTService
	userHandler 		*handler.UserHandler
	propertyHandler 	*handler.PropertyHandler
	healthHandler 		*handler.HealthHandler
}

func NewContainer() *Container {
	logrus.Info("Initializing DI container")

	container := &Container{}

	db.Init()
	container.SqlDB = db.GetSqlDB()

	if container.SqlDB == nil {
		logrus.Fatal("Failed to initialize database connection")
	}

	//repos
	container.userRepo = repository.NewUserRepository(container.SqlDB)
	container.propertyRepo = repository.NewPropertyRepository(container.SqlDB)

	//services
	container.jwtService = service.NewJWTService(container.userRepo)

	//usecases
	container.userUsecase = usecase.NewUserUseCase(container.userRepo, container.jwtService)
	container.propertyUsecase = usecase.NewPropertyUseCase(container.propertyRepo)

	//handlers
	container.userHandler = handler.NewUserHandler(container.userUsecase)
	container.propertyHandler = handler.NewPropertyHandler(container.propertyUsecase)
	container.healthHandler = handler.NewHealthHandler()

	logrus.Info("DI container initialized successfully")
	return container
}

type Handlers struct {
	PropertyHandler 	*handler.PropertyHandler
	UserHandler   		*handler.UserHandler
	HealthHandler 		*handler.HealthHandler
}

func (c *Container) GetHandlers() *Handlers {
	return &Handlers{
		PropertyHandler: c.propertyHandler,
		UserHandler:  c.userHandler,
		HealthHandler: c.healthHandler,
	}
}