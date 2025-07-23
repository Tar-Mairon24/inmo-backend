package usecase

import (
	"inmo-backend/internal/domain/models"
	"inmo-backend/internal/domain/ports"
)

type UserUseCase struct {
	repo ports.UserRepository
}

func NewUserUseCase(repo ports.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) GetAllUsers() ([]models.User, error) {
	return uc.repo.GetAll()
}

func (uc *UserUseCase) GetUserByID(id uint) (*models.User, error) {
	return uc.repo.GetByID(id)
}
