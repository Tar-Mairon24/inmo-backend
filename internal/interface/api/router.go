package api

import (
	"inmo-backend/internal/interface/api/handler"
	"inmo-backend/internal/infrastructure/repository"
	"inmo-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the Gin router
func SetupRouter() *gin.Engine {
	r := gin.Default()

	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	setupHealthRoutes(r)

	v1 := r.Group("/api/v1")
	{
		setupUserRoutes(v1, userHandler)
	}

	return r
}
