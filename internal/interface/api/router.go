package api

import (
	"inmo-backend/internal/interface/api/handler"
	"inmo-backend/internal/infrastructure/repository"
	"inmo-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)
	healthHandler := handler.NewHealthHandler()


	v1 := r.Group("/api/v1")
	{
		setupHealthRoutes(v1, healthHandler)
		setupUserRoutes(v1, userHandler)
	}

	return r
}
