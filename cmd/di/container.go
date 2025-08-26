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
	authHandler 		*handler.AuthHandler
}

func NewContainer() *Container {
	logrus.Info("Initializing DI container")

	container := &Container{}

	db.Init()
	container.SqlDB = db.GetSqlDB()

	if container.SqlDB == nil {
		logrus.Fatal("Failed to initialize database connection")
	}

	// repos
	container.userRepo = repository.NewUserRepository(container.SqlDB)
	container.propertyRepo = repository.NewPropertyRepository(container.SqlDB)

	// services
	container.jwtService = service.NewJWTService(container.userRepo)

	// usecases
	container.userUsecase = usecase.NewUserUseCase(container.userRepo, container.jwtService)
	container.propertyUsecase = usecase.NewPropertyUseCase(container.propertyRepo)

	// handlers
	container.userHandler = handler.NewUserHandler(container.userUsecase)
	container.propertyHandler = handler.NewPropertyHandler(container.propertyUsecase)
	container.authHandler = handler.NewAuthHandler(container.jwtService, container.userUsecase)
	container.healthHandler = handler.NewHealthHandler()

	logrus.Info("DI container initialized successfully")
	return container
}

type Handlers struct {
	PropertyHandler 	*handler.PropertyHandler
	UserHandler   		*handler.UserHandler
	HealthHandler 		*handler.HealthHandler
	AuthHandler 		*handler.AuthHandler
}

func (c *Container) GetHandlers() *Handlers {
	return &Handlers{
		PropertyHandler: c.propertyHandler,
		UserHandler:  c.userHandler,
		HealthHandler: c.healthHandler,
		AuthHandler: c.authHandler,
	}
}

type Services struct {
	JwtService 			ports.JWTService
}

func (c *Container) GetServices() Services{
	return Services{
		JwtService: c.jwtService,
	}
}